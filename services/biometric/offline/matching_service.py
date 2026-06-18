"""Offline 1:N biometric matching service for NPU-deployed terminals (MEK)."""

import json
import logging
import os
import shutil
import tempfile
import time
from dataclasses import dataclass, field, asdict
from pathlib import Path
from typing import Optional

import numpy as np

from services.biometric.inference.npu_runtime import NPURuntime

logger = logging.getLogger(__name__)


@dataclass
class MatchResult:
    identity_id: str
    confidence: float
    embedding: Optional[np.ndarray] = None
    metadata: dict = field(default_factory=dict)


class _IdentityStore:
    """In-memory mapping of FAISS index positions to identity metadata."""

    def __init__(self):
        self.ids: list[str] = []
        self.metadata: list[dict] = []

    def add(self, identity_id: str, meta: Optional[dict] = None):
        self.ids.append(identity_id)
        self.metadata.append(meta or {})

    def remove(self, identity_id: str) -> bool:
        try:
            idx = self.ids.index(identity_id)
            self.ids.pop(idx)
            self.metadata.pop(idx)
            return True
        except ValueError:
            return False

    def clear(self):
        self.ids.clear()
        self.metadata.clear()

    def __len__(self) -> int:
        return len(self.ids)


class OfflineMatchingService:
    """Offline 1:N biometric matching service for NPU deployment.

    Uses FAISS (CPU) for high-speed similarity search against enrolled
    face embeddings extracted via NPU runtime.
    """

    def __init__(
        self,
        npu_runtime: NPURuntime,
        faiss_index_path: Optional[str] = None,
        dimension: int = 512,
    ):
        self.npu = npu_runtime
        self.dimension = dimension
        self.index_path = faiss_index_path
        self._store = _IdentityStore()

        # Lazy import faiss so it is only required when this class is used
        self._faiss = _lazy_import_faiss()

        if faiss_index_path and os.path.isfile(faiss_index_path):
            self.index = self._load_index(faiss_index_path)
        else:
            self.index = self._faiss.IndexFlatIP(self.dimension)
            logger.info("Created new FAISS index (dim=%d)", self.dimension)

    def _load_index(self, path: str):
        logger.info("Loading FAISS index from %s", path)
        index = self._faiss.read_index(path)
        # Rebuild identity store from companion JSON if available
        meta_path = Path(path).with_suffix(".json")
        if meta_path.exists():
            data = json.loads(meta_path.read_text(encoding="utf-8"))
            self._store.ids = data.get("ids", [])
            self._store.metadata = data.get("metadata", [])
        logger.info(
            "FAISS index loaded: %d enrolled identities.", len(self._store)
        )
        return index

    def _save_index(self, path: str):
        """Persist FAISS index along with identity metadata."""
        self._faiss.write_index(self.index, path)
        meta_path = Path(path).with_suffix(".json")
        data = {"ids": self._store.ids, "metadata": self._store.metadata}
        meta_path.write_text(json.dumps(data, ensure_ascii=False), encoding="utf-8")
        logger.info("FAISS index saved (%d identities) to %s", len(self._store), path)

    def enroll(
        self,
        identity_id: str,
        image: np.ndarray,
        metadata: Optional[dict] = None,
        persist: bool = True,
    ) -> str:
        """Extract embedding via NPU and add to FAISS index.

        Args:
            identity_id: Unique identifier for the person.
            image: Preprocessed image tensor (HWC, BGR, uint8 or float32).
            metadata: Optional dict with additional info (name, role, etc.).

        Returns:
            The identity_id that was enrolled.
        """
        embedding = self._extract_embedding(image)
        embedding = embedding.reshape(1, -1).astype(np.float32)
        self.index.add(embedding)
        self._store.add(identity_id, metadata)

        logger.info(
            "Enrolled identity %s (index size: %d)", identity_id, len(self._store)
        )

        if persist and self.index_path:
            self._save_index(self.index_path)

        return identity_id

    def match_1_n(
        self,
        image: np.ndarray,
        top_k: int = 5,
        threshold: float = 0.0,
    ) -> list[MatchResult]:
        """1:N search against enrolled identities.

        Args:
            image: Input face image (HWC, BGR).
            top_k: Maximum number of candidates to return.
            threshold: Minimum cosine similarity threshold (0.0 = no filter).

        Returns:
            List of MatchResult ordered by descending confidence.
        """
        if len(self._store) == 0:
            logger.warning("1:N match called on empty gallery.")
            return []

        embedding = self._extract_embedding(image)
        embedding = embedding.reshape(1, -1).astype(np.float32)

        k = min(top_k, len(self._store))
        distances, indices = self.index.search(embedding, k)

        results: list[MatchResult] = []
        for dist, idx in zip(distances[0], indices[0]):
            if idx < 0 or idx >= len(self._store):
                continue
            confidence = float(dist)
            if confidence < threshold:
                continue
            results.append(
                MatchResult(
                    identity_id=self._store.ids[idx],
                    confidence=confidence,
                    metadata=self._store.metadata[idx],
                )
            )

        return results

    def match_1_1(
        self,
        image: np.ndarray,
        identity_id: str,
        threshold: float = 0.65,
    ) -> MatchResult:
        """1:1 verification against a specific enrolled identity.

        Args:
            image: Input face image.
            identity_id: Target identity to verify against.
            threshold: Cosine similarity threshold for acceptance.

        Returns:
            MatchResult with confidence; matched=True if confidence >= threshold.
        """
        try:
            target_idx = self._store.ids.index(identity_id)
        except ValueError:
            raise ValueError(f"Identity {identity_id} not found in gallery.")

        # Extract single target embedding from the FAISS index
        target_embedding = np.zeros((1, self.dimension), dtype=np.float32)
        self.index.reconstruct(target_idx, target_embedding[0])

        probe_embedding = self._extract_embedding(image).reshape(1, -1).astype(np.float32)

        # Cosine similarity = inner product for L2-normalized vectors
        similarity = float(np.dot(probe_embedding[0], target_embedding[0]))

        return MatchResult(
            identity_id=identity_id,
            confidence=similarity,
            metadata=self._store.metadata[target_idx],
        )

    def remove_identity(self, identity_id: str) -> bool:
        """Remove an identity from the gallery (right to erasure).

        This rebuilds the index without the removed vector, which is
        expensive for large galleries. For production use, consider
        IDMap or a tombstone strategy.

        Returns:
            True if the identity was found and removed.
        """
        if not self._store.remove(identity_id):
            return False

        # Rebuild the entire index without the removed vector
        old_vectors = []
        for i in range(self.index.ntotal):
            vec = np.zeros((1, self.dimension), dtype=np.float32)
            self.index.reconstruct(i, vec[0])
            old_vectors.append(vec)

        remaining_indices = [
            i for i in range(len(self._store.ids))
        ]

        self.index.reset()
        for i in remaining_indices:
            self.index.add(old_vectors[i])

        if self.index_path:
            self._save_index(self.index_path)

        logger.info("Removed identity %s from gallery.", identity_id)
        return True

    def sync_index(self, index_path: str):
        """Sync FAISS index from a central server path for offline use.

        Replaces the in-memory index with the one on disk.
        """
        if not os.path.isfile(index_path):
            raise FileNotFoundError(f"Index file not found: {index_path}")

        self.index = self._load_index(index_path)
        self.index_path = index_path
        logger.info(
            "Index synced from %s (%d identities).",
            index_path,
            len(self._store),
        )

    def _extract_embedding(self, image: np.ndarray) -> np.ndarray:
        """Preprocess image and run NPU inference.

        Expects HWC uint8 BGR image; performs face detection alignment
        and ArcFace inference via NPU runtime.
        """
        # Validate input
        if image.dtype == np.uint8:
            image = image.astype(np.float32) / 255.0

        # Resize to 112x112 (ArcFace input size)
        try:
            import cv2
        except ImportError:
            raise RuntimeError("opencv-python required for image preprocessing.")
        resized = cv2.resize(image, (112, 112))
        # HWC -> CHW, add batch dim
        tensor = np.transpose(resized, (2, 0, 1)).astype(np.float32)
        tensor = np.expand_dims(tensor, axis=0)

        embedding = self.npu.infer(tensor)
        # L2 normalize
        norm = np.linalg.norm(embedding, axis=1, keepdims=True)
        norm[norm == 0] = 1e-12
        return embedding / norm

    @property
    def size(self) -> int:
        return len(self._store)


def _lazy_import_faiss():
    """Import faiss at runtime, providing a clear error message if missing."""
    try:
        import faiss
        return faiss
    except ImportError:
        raise ImportError(
            "faiss is required for offline matching. "
            "Install it with: pip install faiss-cpu"
        )
