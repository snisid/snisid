import pytest
import numpy as np
import json
import os
from unittest.mock import patch, MagicMock, PropertyMock
from pathlib import Path


class TestGalleryManager:
    @pytest.fixture
    def manager(self, tmp_path):
        from services.biometric.offline.gallery_manager import (
            OfflineGalleryManager,
        )

        return OfflineGalleryManager(str(tmp_path / "galleries"))

    def test_create_gallery_creates_file(self, manager, tmp_path):
        identities = [{"id": "user_001"}, {"id": "user_002"}]
        embeddings = np.random.rand(2, 512).astype(np.float32)
        gallery_path = manager.create_gallery("test_gal", identities, embeddings)
        assert os.path.isfile(gallery_path)
        assert gallery_path.endswith(".gallery")

    def test_create_gallery_mismatched_counts_raises(self, manager):
        identities = [{"id": "user_001"}]
        embeddings = np.random.rand(3, 512).astype(np.float32)
        with pytest.raises(ValueError, match="must match"):
            manager.create_gallery("bad_gal", identities, embeddings)

    def test_load_gallery_returns_faiss_index(self, manager):
        import faiss

        identities = [{"id": "user_001"}]
        embeddings = np.random.rand(1, 512).astype(np.float32)
        manager.create_gallery("load_test", identities, embeddings)
        index = manager.load_gallery("load_test")
        assert isinstance(index, faiss.IndexFlatIP)
        assert index.ntotal == 1

    def test_load_nonexistent_gallery_raises(self, manager):
        with pytest.raises(FileNotFoundError):
            manager.load_gallery("nonexistent")

    def test_list_galleries_returns_manifests(self, manager):
        identities = [{"id": "u1"}, {"id": "u2"}]
        embeddings = np.random.rand(2, 512).astype(np.float32)
        manager.create_gallery("list_test", identities, embeddings)
        manifests = manager.list_galleries()
        assert len(manifests) >= 1
        assert manifests[0].name == "list_test"
        assert manifests[0].identity_count == 2

    def test_encrypt_gallery_validates_key_length(self, manager):
        with pytest.raises(ValueError, match="32 bytes"):
            manager.encrypt_gallery("dummy.gallery", b"short_key")

    def test_encrypt_gallery_creates_enc_file(self, manager, tmp_path):
        from services.biometric.offline.gallery_manager import (
            OfflineGalleryManager,
        )

        dummy_path = str(tmp_path / "dummy.gallery")
        Path(dummy_path).write_text("fake gallery content")
        enc_path = manager.encrypt_gallery(dummy_path, b"a" * 32)
        assert os.path.isfile(enc_path)
        assert enc_path.endswith(".enc")

    def test_decrypt_gallery_roundtrip(self, manager, tmp_path):
        from services.biometric.offline.gallery_manager import (
            OfflineGalleryManager,
        )

        original_content = b"test gallery content for roundtrip"
        gallery_path = str(tmp_path / "roundtrip.gallery")
        Path(gallery_path).write_bytes(original_content)

        key = b"k" * 32
        enc_path = manager.encrypt_gallery(gallery_path, key)
        dec_path = manager.decrypt_gallery(enc_path, key)
        assert Path(dec_path).read_bytes() == original_content

    def test_gallery_contains_manifest_with_checksum(self, manager):
        identities = [{"id": "u1"}]
        embeddings = np.random.rand(1, 512).astype(np.float32)
        gallery_path = manager.create_gallery("checksum_test", identities, embeddings)
        import zipfile

        with zipfile.ZipFile(gallery_path, "r") as zf:
            meta = json.loads(zf.read("checksum_test_meta.json"))
        assert "manifest" in meta
        assert "checksum" in meta["manifest"]
        assert len(meta["manifest"]["checksum"]) == 64


class TestOfflineMatchingService:
    @pytest.fixture
    def npu_mock(self):
        mock = MagicMock()
        mock.infer.return_value = np.random.rand(1, 512).astype(np.float32)
        return mock

    @pytest.fixture
    def service(self, npu_mock, tmp_path):
        with patch(
            "services.biometric.offline.matching_service._lazy_import_faiss"
        ) as mock_faiss_import:
            import faiss

            mock_faiss_import.return_value = faiss
            from services.biometric.offline.matching_service import (
                OfflineMatchingService,
            )

            return OfflineMatchingService(
                npu_mock,
                faiss_index_path=str(tmp_path / "test_index"),
                dimension=512,
            )

    def test_enroll_adds_identity(self, service):
        image = np.random.randint(0, 256, (112, 112, 3), dtype=np.uint8)
        identity_id = service.enroll("user_001", image)
        assert identity_id == "user_001"
        assert service.size == 1

    def test_enroll_multiple_identities(self, service):
        for i in range(5):
            image = np.random.randint(0, 256, (112, 112, 3), dtype=np.uint8)
            service.enroll(f"user_{i:03d}", image)
        assert service.size == 5

    def test_match_1_n_returns_results(self, service):
        image = np.random.randint(0, 256, (112, 112, 3), dtype=np.uint8)
        service.enroll("target_user", image)
        results = service.match_1_n(image, top_k=3)
        assert len(results) >= 1
        assert results[0].identity_id == "target_user"

    def test_match_1_n_empty_gallery_returns_empty(self, service):
        image = np.random.randint(0, 256, (112, 112, 3), dtype=np.uint8)
        results = service.match_1_n(image)
        assert results == []

    def test_match_1_1_verification(self, service):
        image = np.random.randint(0, 256, (112, 112, 3), dtype=np.uint8)
        service.enroll("verify_user", image)
        result = service.match_1_1(image, "verify_user")
        assert result.identity_id == "verify_user"
        assert isinstance(result.confidence, float)

    def test_match_1_1_nonexistent_identity_raises(self, service):
        image = np.random.randint(0, 256, (112, 112, 3), dtype=np.uint8)
        with pytest.raises(ValueError, match="not found"):
            service.match_1_1(image, "nonexistent")

    def test_remove_identity(self, service):
        image = np.random.randint(0, 256, (112, 112, 3), dtype=np.uint8)
        service.enroll("remove_me", image)
        assert service.size == 1
        removed = service.remove_identity("remove_me")
        assert removed is True
        assert service.size == 0

    def test_remove_nonexistent_identity_returns_false(self, service):
        assert service.remove_identity("not_there") is False

    def test_embedding_extraction_normalizes_input(self, service, npu_mock):
        image = np.random.randint(0, 256, (112, 112, 3), dtype=np.uint8)
        embedding = service._extract_embedding(image)
        assert np.isclose(np.linalg.norm(embedding), 1.0, atol=1e-5)


class TestOfflineSyncProtocol:
    @pytest.fixture
    def crypto(self):
        from services.biometric.security.crypto import BiometricCryptoVault

        return BiometricCryptoVault()

    @pytest.fixture
    def protocol(self, tmp_path, crypto):
        from services.biometric.offline.sync_protocol import (
            OfflineSyncProtocol,
        )

        return OfflineSyncProtocol(
            terminal_id="term_001",
            storage_path=str(tmp_path / "sync_storage"),
            crypto=crypto,
        )

    def test_create_sync_bundle_has_uuid(self, protocol):
        from services.biometric.offline.sync_protocol import SyncEntry, SyncAction

        entry = SyncEntry(
            identity_id="user_001",
            action=SyncAction.ADD,
            embedding=[0.1] * 512,
        )
        bundle = protocol.create_sync_bundle([entry])
        assert bundle.bundle_id is not None
        assert len(bundle.bundle_id) > 0

    def test_create_sync_bundle_stores_locally(self, protocol):
        from services.biometric.offline.sync_protocol import SyncEntry, SyncAction

        entry = SyncEntry(
            identity_id="user_001",
            action=SyncAction.ADD,
            embedding=[0.5] * 512,
            metadata={"name": "Alice"},
        )
        bundle = protocol.create_sync_bundle([entry])
        assert os.path.isfile(protocol._local_gallery_path)
        gallery = json.loads(
            Path(protocol._local_gallery_path).read_text(encoding="utf-8")
        )
        assert "user_001" in gallery

    def test_apply_sync_bundle_adds_entries(self, protocol):
        from services.biometric.offline.sync_protocol import SyncEntry, SyncAction, SyncBundle
        from uuid import uuid4

        entry = SyncEntry(
            identity_id="remote_user",
            action=SyncAction.ADD,
            embedding=[0.3] * 512,
        )
        bundle = SyncBundle(
            bundle_id=str(uuid4()),
            terminal_id="remote_term",
            created_at=1000.0,
            entries=[entry],
        )
        result = protocol.apply_sync_bundle(bundle)
        assert result.applied_entries == 1
        assert result.success is True

    def test_apply_sync_bundle_detects_duplicate_add(self, protocol):
        from services.biometric.offline.sync_protocol import SyncEntry, SyncAction, SyncBundle
        from uuid import uuid4

        entry = SyncEntry(
            identity_id="dup_user",
            action=SyncAction.ADD,
            embedding=[0.3] * 512,
        )
        bundle1 = SyncBundle(
            bundle_id=str(uuid4()),
            terminal_id="remote_term",
            created_at=1000.0,
            entries=[entry],
        )
        protocol.apply_sync_bundle(bundle1)

        bundle2 = SyncBundle(
            bundle_id=str(uuid4()),
            terminal_id="remote_term",
            created_at=1001.0,
            entries=[entry],
        )
        result = protocol.apply_sync_bundle(bundle2)
        assert result.skipped_entries == 1
        assert len(result.errors) >= 1

    def test_apply_sync_bundle_update(self, protocol):
        from services.biometric.offline.sync_protocol import SyncEntry, SyncAction, SyncBundle
        from uuid import uuid4

        add_entry = SyncEntry(
            identity_id="update_user",
            action=SyncAction.ADD,
            embedding=[0.3] * 512,
        )
        add_bundle = SyncBundle(
            bundle_id=str(uuid4()),
            terminal_id="remote_term",
            created_at=1000.0,
            entries=[add_entry],
        )
        protocol.apply_sync_bundle(add_bundle)

        upd_entry = SyncEntry(
            identity_id="update_user",
            action=SyncAction.UPDATE,
            embedding=[0.9] * 512,
            metadata={"role": "updated"},
        )
        upd_bundle = SyncBundle(
            bundle_id=str(uuid4()),
            terminal_id="remote_term",
            created_at=1001.0,
            entries=[upd_entry],
        )
        result = protocol.apply_sync_bundle(upd_bundle)
        assert result.applied_entries == 1

    def test_apply_sync_bundle_delete(self, protocol):
        from services.biometric.offline.sync_protocol import SyncEntry, SyncAction, SyncBundle
        from uuid import uuid4

        add_entry = SyncEntry(
            identity_id="delete_me",
            action=SyncAction.ADD,
            embedding=[0.3] * 512,
        )
        add_bundle = SyncBundle(
            bundle_id=str(uuid4()),
            terminal_id="remote_term",
            created_at=1000.0,
            entries=[add_entry],
        )
        protocol.apply_sync_bundle(add_bundle)

        del_entry = SyncEntry(
            identity_id="delete_me",
            action=SyncAction.DELETE,
            embedding=[0.3] * 512,
        )
        del_bundle = SyncBundle(
            bundle_id=str(uuid4()),
            terminal_id="remote_term",
            created_at=1002.0,
            entries=[del_entry],
        )
        result = protocol.apply_sync_bundle(del_bundle)
        assert result.applied_entries == 1

    def test_detect_conflict_on_existing_identity(self, protocol):
        from services.biometric.offline.sync_protocol import SyncEntry, SyncAction, SyncBundle
        from uuid import uuid4

        add_entry = SyncEntry(
            identity_id="conflict_user",
            action=SyncAction.ADD,
            embedding=[0.3] * 512,
        )
        add_bundle = SyncBundle(
            bundle_id=str(uuid4()),
            terminal_id="local_term",
            created_at=1000.0,
            entries=[add_entry],
        )
        protocol.apply_sync_bundle(add_bundle)

        conflict_bundle = SyncBundle(
            bundle_id=str(uuid4()),
            terminal_id="remote_term",
            created_at=1001.0,
            entries=[add_entry],
        )
        conflicts = protocol.detect_conflicts(conflict_bundle)
        assert len(conflicts) >= 1
        assert conflicts[0].conflict_type.name == "IDENTITY_EXISTS"

    def test_encrypt_decrypt_bundle_roundtrip(self, protocol):
        from services.biometric.offline.sync_protocol import SyncEntry, SyncAction

        entry = SyncEntry(
            identity_id="crypto_user",
            action=SyncAction.ADD,
            embedding=[0.4] * 512,
        )
        bundle = protocol.create_sync_bundle([entry])
        encrypted = protocol.encrypt_bundle(bundle)
        assert isinstance(encrypted, bytes)

        decrypted = protocol.decrypt_bundle(encrypted)
        assert decrypted.bundle_id == bundle.bundle_id
        assert len(decrypted.entries) == 1
        assert decrypted.entries[0].identity_id == "crypto_user"

    def test_export_import_sync_package(self, protocol, tmp_path):
        from services.biometric.offline.sync_protocol import SyncEntry, SyncAction

        entry = SyncEntry(
            identity_id="export_user",
            action=SyncAction.ADD,
            embedding=[0.2] * 512,
        )
        bundle = protocol.create_sync_bundle([entry])
        pkg_path = str(tmp_path / "sync_package.bin")
        protocol.export_sync_package(bundle, pkg_path)
        assert os.path.isfile(pkg_path)

        imported = protocol.import_sync_package(pkg_path)
        assert imported.bundle_id == bundle.bundle_id
        assert imported.entries[0].identity_id == "export_user"

    def test_last_sync_time_none_when_no_sync(self, protocol):
        assert protocol.last_sync_time is None

    def test_last_sync_time_after_sync(self, protocol):
        from services.biometric.offline.sync_protocol import SyncEntry, SyncAction

        entry = SyncEntry(
            identity_id="time_test",
            action=SyncAction.ADD,
            embedding=[0.5] * 512,
        )
        protocol.create_sync_bundle([entry])
        assert protocol.last_sync_time is not None
