"""Offline gallery management for biometric terminals.

Manages the lifecycle of compressed, encrypted FAISS galleries
suitable for deployment on edge NPU terminals (MEK).
"""

import hashlib
import json
import logging
import os
import shutil
import tempfile
import time
import zipfile
from dataclasses import dataclass
from pathlib import Path
from typing import Optional
from urllib.parse import urljoin

import numpy as np

logger = logging.getLogger(__name__)


@dataclass
class GalleryManifest:
    name: str
    version: str
    created_at: float
    identity_count: int
    checksum: str
    encrypted: bool = False


class OfflineGalleryManager:
    """Manage biometric galleries for offline terminals.

    Galleries are stored as compressed, optionally encrypted FAISS index
    bundles that can be synced from a central server.
    """

    def __init__(self, gallery_dir: str):
        self.gallery_dir = gallery_dir
        Path(gallery_dir).mkdir(parents=True, exist_ok=True)
        logger.info("Gallery manager initialized at %s", gallery_dir)

    def create_gallery(
        self,
        name: str,
        identities: list[dict],
        embeddings: np.ndarray,
        version: str = "1.0.0",
    ) -> str:
        """Create a compressed, encrypted gallery bundle from identities.

        Args:
            name: Gallery name (e.g. "region_ouest_2026").
            identities: List of dicts with keys "id" and optional "metadata".
            embeddings: (N, D) numpy array of L2-normalized embeddings.
            version: Semantic version string.

        Returns:
            Path to the created gallery bundle (.gallery file, a zip archive).
        """
        import faiss

        if len(identities) != embeddings.shape[0]:
            raise ValueError(
                f"Number of identities ({len(identities)}) must match "
                f"embeddings rows ({embeddings.shape[0]})."
            )

        dimension = embeddings.shape[1]
        index = faiss.IndexFlatIP(dimension)
        index.add(embeddings.astype(np.float32))

        bundle_dir = Path(tempfile.mkdtemp())
        try:
            index_path = bundle_dir / f"{name}.index"
            faiss.write_index(index, str(index_path))

            manifest = GalleryManifest(
                name=name,
                version=version,
                created_at=time.time(),
                identity_count=len(identities),
                checksum=self._compute_checksum(str(index_path)),
            )

            meta_path = bundle_dir / f"{name}_meta.json"
            meta_data = {
                "manifest": {
                    "name": manifest.name,
                    "version": manifest.version,
                    "created_at": manifest.created_at,
                    "identity_count": manifest.identity_count,
                    "checksum": manifest.checksum,
                },
                "identities": identities,
            }
            meta_path.write_text(json.dumps(meta_data, ensure_ascii=False), encoding="utf-8")

            gallery_path = os.path.join(self.gallery_dir, f"{name}.gallery")
            with zipfile.ZipFile(gallery_path, "w", zipfile.ZIP_DEFLATED) as zf:
                zf.write(str(index_path), arcname=f"{name}.index")
                zf.write(str(meta_path), arcname=f"{name}_meta.json")

            logger.info(
                "Gallery '%s' created with %d identities at %s",
                name,
                len(identities),
                gallery_path,
            )
            return gallery_path

        finally:
            shutil.rmtree(bundle_dir, ignore_errors=True)

    def load_gallery(self, name: str) -> "FAISSIndex":
        """Load gallery into memory for matching.

        Args:
            name: Gallery name (without extension).

        Returns:
            FAISS IndexFlatIP loaded into memory.
        """
        import faiss

        gallery_path = os.path.join(self.gallery_dir, f"{name}.gallery")
        if not os.path.isfile(gallery_path):
            raise FileNotFoundError(f"Gallery not found: {gallery_path}")

        bundle_dir = Path(tempfile.mkdtemp())
        try:
            with zipfile.ZipFile(gallery_path, "r") as zf:
                zf.extractall(str(bundle_dir))

            index_path = bundle_dir / f"{name}.index"
            if not index_path.exists():
                raise RuntimeError(f"Index file missing in gallery bundle: {name}")

            index = faiss.read_index(str(index_path))
            logger.info("Gallery '%s' loaded (%d vectors).", name, index.ntotal)
            return index

        finally:
            shutil.rmtree(bundle_dir, ignore_errors=True)

    def encrypt_gallery(self, gallery_path: str, key: bytes) -> str:
        """Encrypt gallery with AES-256-GCM for secure offline storage.

        Produces a `.gallery.enc` file.

        Args:
            gallery_path: Path to existing .gallery file.
            key: 32-byte AES-256 key.

        Returns:
            Path to the encrypted gallery file.
        """
        from cryptography.hazmat.primitives.ciphers.aead import AESGCM

        if len(key) != 32:
            raise ValueError("AES-256 key must be exactly 32 bytes.")

        with open(gallery_path, "rb") as f:
            plaintext = f.read()

        aesgcm = AESGCM(key)
        nonce = os.urandom(12)
        ciphertext = aesgcm.encrypt(nonce, plaintext, None)

        enc_path = gallery_path + ".enc"
        with open(enc_path, "wb") as f:
            f.write(nonce)
            f.write(ciphertext)

        logger.info("Gallery encrypted: %s", enc_path)
        return enc_path

    def decrypt_gallery(self, encrypted_path: str, key: bytes) -> str:
        """Decrypt gallery for loading.

        Writes a temporary .gallery file and returns its path.

        Args:
            encrypted_path: Path to .gallery.enc file.
            key: 32-byte AES-256 key.

        Returns:
            Path to the decrypted temporary file (caller should clean up).
        """
        from cryptography.hazmat.primitives.ciphers.aead import AESGCM

        if len(key) != 32:
            raise ValueError("AES-256 key must be exactly 32 bytes.")

        with open(encrypted_path, "rb") as f:
            nonce = f.read(12)
            ciphertext = f.read()

        aesgcm = AESGCM(key)
        plaintext = aesgcm.decrypt(nonce, ciphertext, None)

        out_path = encrypted_path.replace(".enc", "")
        with open(out_path, "wb") as f:
            f.write(plaintext)

        logger.info("Gallery decrypted: %s", out_path)
        return out_path

    def sync_from_server(
        self,
        api_url: str,
        api_key: str,
        region: str,
        timeout: int = 60,
    ) -> list[str]:
        """Download latest gallery bundle(s) from central server.

        Args:
            api_url: Base URL of the central sync API.
            api_key: Bearer token for authentication.
            region: Terminal region code (e.g. "region_ouest").
            timeout: HTTP request timeout in seconds.

        Returns:
            List of downloaded gallery file paths.
        """
        import requests

        headers = {
            "Authorization": f"Bearer {api_key}",
            "X-Region": region,
        }

        sync_url = urljoin(api_url.rstrip("/") + "/", "api/v1/sync/galleries")
        logger.info("Syncing galleries from %s (region=%s)", sync_url, region)

        resp = requests.get(sync_url, headers=headers, timeout=timeout)
        resp.raise_for_status()
        payload = resp.json()

        downloaded: list[str] = []
        for gallery_info in payload.get("galleries", []):
            name = gallery_info["name"]
            download_url = gallery_info["download_url"]

            gal_resp = requests.get(download_url, headers=headers, timeout=timeout)
            gal_resp.raise_for_status()

            gal_path = os.path.join(self.gallery_dir, f"{name}.gallery")
            with open(gal_path, "wb") as f:
                f.write(gal_resp.content)

            downloaded.append(gal_path)
            logger.info("Downloaded gallery: %s", gal_path)

        return downloaded

    def list_galleries(self) -> list[GalleryManifest]:
        """List all available gallery bundles in the gallery directory."""
        results: list[GalleryManifest] = []
        for fpath in Path(self.gallery_dir).glob("*.gallery"):
            try:
                with zipfile.ZipFile(str(fpath), "r") as zf:
                    meta_name = fpath.stem + "_meta.json"
                    if meta_name not in zf.namelist():
                        continue
                    meta_data = json.loads(zf.read(meta_name))
                    manifest = GalleryManifest(**meta_data["manifest"])
                    results.append(manifest)
            except Exception as exc:
                logger.warning("Failed to read manifest from %s: %s", fpath, exc)
                continue
        return sorted(results, key=lambda m: m.created_at, reverse=True)

    @staticmethod
    def _compute_checksum(file_path: str, algorithm: str = "sha256") -> str:
        h = hashlib.new(algorithm)
        with open(file_path, "rb") as f:
            for chunk in iter(lambda: f.read(65536), b""):
                h.update(chunk)
        return h.hexdigest()
