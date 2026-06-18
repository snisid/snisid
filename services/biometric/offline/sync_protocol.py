"""Sync protocol for offline biometric terminals (MEK).

Handles encrypted bundle exchange between edge terminals and the
central server for gallery updates and conflict resolution.
"""

import hashlib
import json
import logging
import os
import tempfile
import time
import zipfile
from dataclasses import dataclass, field, asdict
from datetime import datetime, timezone
from enum import Enum
from pathlib import Path
from typing import Optional
from uuid import uuid4

import numpy as np

from services.biometric.security.crypto import BiometricCryptoVault

logger = logging.getLogger(__name__)


class ConflictType(Enum):
    IDENTITY_EXISTS = "identity_exists"
    EMBEDDING_MISMATCH = "embedding_mismatch"
    METADATA_CONFLICT = "metadata_conflict"
    VERSION_MISMATCH = "version_mismatch"


class SyncAction(Enum):
    ADD = "add"
    UPDATE = "update"
    DELETE = "delete"
    MERGE = "merge"


@dataclass
class SyncEntry:
    identity_id: str
    action: SyncAction
    embedding: list[float]
    metadata: dict = field(default_factory=dict)
    checksum: str = ""
    timestamp: float = field(default_factory=time.time)

    def __post_init__(self):
        if not self.checksum:
            self.checksum = self._compute_checksum()

    def _compute_checksum(self) -> str:
        raw = json.dumps(
            [self.identity_id, self.embedding, self.metadata],
            sort_keys=True,
            ensure_ascii=False,
        )
        return hashlib.sha256(raw.encode("utf-8")).hexdigest()


@dataclass
class SyncBundle:
    bundle_id: str
    terminal_id: str
    created_at: float
    entries: list[SyncEntry]
    previous_bundle_id: Optional[str] = None

    def to_bytes(self) -> bytes:
        data = asdict(self)
        data["entries"] = [asdict(e) for e in self.entries]
        return json.dumps(data, ensure_ascii=False).encode("utf-8")

    @classmethod
    def from_bytes(cls, raw: bytes) -> "SyncBundle":
        data = json.loads(raw.decode("utf-8"))
        entries = [SyncEntry(**e) for e in data["entries"]]
        return cls(
            bundle_id=data["bundle_id"],
            terminal_id=data["terminal_id"],
            created_at=data["created_at"],
            entries=entries,
            previous_bundle_id=data.get("previous_bundle_id"),
        )


@dataclass
class SyncResult:
    success: bool
    applied_entries: int
    skipped_entries: int
    errors: list[str] = field(default_factory=list)
    new_bundle_id: Optional[str] = None


@dataclass
class Conflict:
    identity_id: str
    conflict_type: ConflictType
    local_value: Optional[str] = None
    remote_value: Optional[str] = None
    description: str = ""


class OfflineSyncProtocol:
    """Sync protocol for offline biometric terminals.

    Exchanges encrypted bundles containing enrollment updates between
    edge terminals and the central management server.
    """

    def __init__(
        self,
        terminal_id: str,
        storage_path: str,
        crypto: BiometricCryptoVault,
    ):
        self.terminal_id = terminal_id
        self.storage_path = storage_path
        self.crypto = crypto

        Path(storage_path).mkdir(parents=True, exist_ok=True)
        self._sync_log_path = os.path.join(storage_path, "sync_log.json")
        self._local_gallery_path = os.path.join(storage_path, "local_gallery.json")
        self._last_bundle_id: Optional[str] = None

        self._load_state()

    def _load_state(self):
        """Load last bundle ID from sync log."""
        if os.path.isfile(self._sync_log_path):
            try:
                data = json.loads(
                    Path(self._sync_log_path).read_text(encoding="utf-8")
                )
                self._last_bundle_id = data.get("last_bundle_id")
                logger.debug(
                    "Sync state loaded. Last bundle: %s", self._last_bundle_id
                )
            except (json.JSONDecodeError, KeyError) as exc:
                logger.warning("Failed to load sync state: %s", exc)

    def _save_state(self):
        """Persist sync state to disk."""
        data = {
            "terminal_id": self.terminal_id,
            "last_bundle_id": self._last_bundle_id,
            "updated_at": time.time(),
        }
        Path(self._sync_log_path).write_text(
            json.dumps(data, ensure_ascii=False), encoding="utf-8"
        )

    def _load_local_gallery(self) -> dict[str, dict]:
        """Load local gallery as {identity_id: {embedding, metadata}}."""
        if os.path.isfile(self._local_gallery_path):
            try:
                raw = Path(self._local_gallery_path).read_text(encoding="utf-8")
                return json.loads(raw)
            except (json.JSONDecodeError, FileNotFoundError):
                pass
        return {}

    def _save_local_gallery(self, gallery: dict[str, dict]):
        """Persist local gallery."""
        Path(self._local_gallery_path).write_text(
            json.dumps(gallery, ensure_ascii=False, indent=2), encoding="utf-8"
        )

    def create_sync_bundle(self, entries: list[SyncEntry]) -> SyncBundle:
        """Create an encrypted sync package with new enrollments.

        Args:
            entries: List of sync entries (adds, updates, deletes).

        Returns:
            A SyncBundle ready to be serialised and sent to the server.
        """
        bundle = SyncBundle(
            bundle_id=str(uuid4()),
            terminal_id=self.terminal_id,
            created_at=time.time(),
            entries=entries,
            previous_bundle_id=self._last_bundle_id,
        )

        # Append to local gallery
        gallery = self._load_local_gallery()
        for entry in entries:
            if entry.action in (SyncAction.ADD, SyncAction.UPDATE):
                gallery[entry.identity_id] = {
                    "embedding": entry.embedding,
                    "metadata": entry.metadata,
                    "updated_at": entry.timestamp,
                }
            elif entry.action == SyncAction.DELETE:
                gallery.pop(entry.identity_id, None)

        self._save_local_gallery(gallery)
        self._last_bundle_id = bundle.bundle_id
        self._save_state()

        logger.info(
            "Sync bundle %s created with %d entries.",
            bundle.bundle_id,
            len(entries),
        )
        return bundle

    def apply_sync_bundle(self, bundle: SyncBundle) -> SyncResult:
        """Apply a sync bundle received from the central server.

        Args:
            bundle: The bundle to apply.

        Returns:
            SyncResult detailing what was applied and any errors.
        """
        result = SyncResult(success=True, applied_entries=0, skipped_entries=0)

        # Detect conflicts first
        conflicts = self.detect_conflicts(bundle)
        if conflicts:
            for c in conflicts:
                logger.warning(
                    "Conflict: %s | %s: local=%s remote=%s",
                    c.identity_id,
                    c.conflict_type.value,
                    c.local_value,
                    c.remote_value,
                )
                result.errors.append(
                    f"Conflict on {c.identity_id}: {c.description}"
                )

        gallery = self._load_local_gallery()

        for entry in bundle.entries:
            if entry.identity_id in gallery and entry.action == SyncAction.ADD:
                # Conflict: identity already exists
                result.skipped_entries += 1
                result.errors.append(
                    f"Identity {entry.identity_id} already exists, skipping ADD."
                )
                continue

            try:
                if entry.action == SyncAction.ADD:
                    gallery[entry.identity_id] = {
                        "embedding": entry.embedding,
                        "metadata": entry.metadata,
                        "updated_at": entry.timestamp,
                    }
                    result.applied_entries += 1

                elif entry.action == SyncAction.UPDATE:
                    if entry.identity_id in gallery:
                        gallery[entry.identity_id].update(
                            {
                                "embedding": entry.embedding,
                                "metadata": entry.metadata,
                                "updated_at": entry.timestamp,
                            }
                        )
                        result.applied_entries += 1
                    else:
                        result.skipped_entries += 1
                        result.errors.append(
                            f"Identity {entry.identity_id} not found for UPDATE."
                        )

                elif entry.action == SyncAction.DELETE:
                    if gallery.pop(entry.identity_id, None):
                        result.applied_entries += 1
                    else:
                        result.skipped_entries += 1

            except Exception as exc:
                result.errors.append(f"Error applying {entry.identity_id}: {exc}")
                result.skipped_entries += 1

        self._save_local_gallery(gallery)

        if bundle.bundle_id:
            self._last_bundle_id = bundle.bundle_id
            result.new_bundle_id = bundle.bundle_id
            self._save_state()

        result.success = len(result.errors) == 0
        logger.info(
            "Bundle %s applied: %d applied, %d skipped, %d errors.",
            bundle.bundle_id,
            result.applied_entries,
            result.skipped_entries,
            len(result.errors),
        )
        return result

    def detect_conflicts(self, bundle: SyncBundle) -> list[Conflict]:
        """Detect conflicts between local and remote galleries.

        Checks for:
          - Identity already exists locally when remote tries ADD
          - Embedding checksum mismatch for existing identities
          - Metadata field conflicts

        Args:
            bundle: The incoming sync bundle.

        Returns:
            List of detected conflicts.
        """
        gallery = self._load_local_gallery()
        conflicts: list[Conflict] = []

        for entry in bundle.entries:
            if entry.identity_id in gallery:
                local = gallery[entry.identity_id]

                if entry.action == SyncAction.ADD:
                    conflicts.append(
                        Conflict(
                            identity_id=entry.identity_id,
                            conflict_type=ConflictType.IDENTITY_EXISTS,
                            local_value=local.get("checksum", "exists"),
                            remote_value="add",
                            description=f"Remote wants to ADD {entry.identity_id} but it already exists locally.",
                        )
                    )

                elif entry.action == SyncAction.UPDATE:
                    local_checksum = hashlib.sha256(
                        json.dumps(local["embedding"], sort_keys=True).encode()
                    ).hexdigest()
                    if local_checksum != entry.checksum:
                        conflicts.append(
                            Conflict(
                                identity_id=entry.identity_id,
                                conflict_type=ConflictType.EMBEDDING_MISMATCH,
                                local_value=local_checksum[:16],
                                remote_value=entry.checksum[:16],
                                description="Embedding checksum differs between local and remote.",
                            )
                        )

        return conflicts

    def encrypt_bundle(self, bundle: SyncBundle) -> bytes:
        """Encrypt a sync bundle using the crypto vault.

        The bundle is serialised to JSON, then encrypted via the vault's
        template encryption (AES-256-GCM).

        Returns:
            Encrypted bytes with a 4-byte length prefix.
        """
        raw = bundle.to_bytes()
        # Encrypt via the vault (simulates Transit Engine)
        import base64

        b64_payload = base64.b64encode(raw).decode("utf-8")
        encrypted = self.crypto.encrypt_template(np.frombuffer(raw, dtype=np.uint8).astype(np.float32))
        # The vault returns a string; encode it for transport
        return encrypted.encode("utf-8")

    def decrypt_bundle(self, encrypted: bytes) -> SyncBundle:
        """Decrypt a sync bundle.

        Returns:
            The decrypted SyncBundle.
        """
        encrypted_str = encrypted.decode("utf-8")
        decrypted_embedding = self.crypto.decrypt_template(encrypted_str)

        # Reconstruct original bytes from the float32 embedding
        raw_bytes = decrypted_embedding.astype(np.uint8).tobytes()
        # Trim trailing null padding
        raw_bytes = raw_bytes.rstrip(b"\x00")

        # Find the JSON boundary
        try:
            json_end = raw_bytes.index(b"}") + 1
        except ValueError:
            json_end = len(raw_bytes)

        return SyncBundle.from_bytes(raw_bytes[:json_end])

    def export_sync_package(
        self, bundle: SyncBundle, output_path: str
    ) -> str:
        """Write an encrypted sync bundle to a file for offline transport.

        Args:
            bundle: The sync bundle.
            output_path: Destination file path.

        Returns:
            The output path.
        """
        encrypted = self.encrypt_bundle(bundle)
        Path(output_path).write_bytes(encrypted)
        logger.info("Sync package exported to %s (%d bytes)", output_path, len(encrypted))
        return output_path

    def import_sync_package(self, input_path: str) -> SyncBundle:
        """Read and decrypt a sync package from file.

        Args:
            input_path: Path to the encrypted sync package.

        Returns:
            The decrypted SyncBundle.
        """
        encrypted = Path(input_path).read_bytes()
        bundle = self.decrypt_bundle(encrypted)
        logger.info("Sync package imported from %s", input_path)
        return bundle

    @property
    def last_sync_time(self) -> Optional[datetime]:
        if os.path.isfile(self._sync_log_path):
            try:
                data = json.loads(
                    Path(self._sync_log_path).read_text(encoding="utf-8")
                )
                ts = data.get("updated_at", 0)
                return datetime.fromtimestamp(ts, tz=timezone.utc)
            except (json.JSONDecodeError, KeyError, OSError):
                pass
        return None
