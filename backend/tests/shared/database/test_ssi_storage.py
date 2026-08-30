import pytest

from shared.database.ssi_models import (
    DIDRecord,
    VerifiableCredentialRecord,
    StatusListRecord,
    WalletCredentialRecord,
    DIDCommMessageRecord,
    CredentialFlowRecord,
    CHAPIRecord,
    CredentialManifestRecord,
    RevocationEventRecord,
)
from shared.ssi_storage import (
    DIDStorage,
    VCStorage,
    StatusListStorage,
    WalletCredentialStorage,
    DIDCommMessageStorage,
    CredentialFlowStorage,
    CHAPIStorage,
    CredentialManifestStorage,
    RevocationEventStorage,
)


class TestDIDStorage:
    @pytest.mark.asyncio
    async def test_save_and_get(self, db_session):
        storage = DIDStorage(db_session)
        await storage.save("did:key:alice", "key", {"id": "did:key:alice"})
        doc = await storage.get("did:key:alice")
        assert doc is not None
        assert doc["id"] == "did:key:alice"

    @pytest.mark.asyncio
    async def test_get_nonexistent(self, db_session):
        storage = DIDStorage(db_session)
        assert await storage.get("did:key:nobody") is None

    @pytest.mark.asyncio
    async def test_delete(self, db_session):
        storage = DIDStorage(db_session)
        await storage.save("did:key:alice", "key", {"id": "did:key:alice"})
        assert await storage.delete("did:key:alice") is True
        assert await storage.delete("did:key:nobody") is False
        assert await storage.get("did:key:alice") is None


class TestVCStorage:
    @pytest.mark.asyncio
    async def test_save_and_get(self, db_session):
        storage = VCStorage(db_session)
        await storage.save("vc-1", "did:key:issuer", "did:key:subject", "IdentityCredential", {"id": "vc-1"})
        doc = await storage.get("vc-1")
        assert doc is not None
        assert doc["id"] == "vc-1"

    @pytest.mark.asyncio
    async def test_list_by_subject(self, db_session):
        storage = VCStorage(db_session)
        await storage.save("vc-1", "issuer", "did:key:alice", "Type1", {"id": "vc-1"})
        await storage.save("vc-2", "issuer", "did:key:alice", "Type2", {"id": "vc-2"})
        await storage.save("vc-3", "issuer", "did:key:bob", "Type1", {"id": "vc-3"})
        alice_vcs = await storage.list_by_subject("did:key:alice")
        assert len(alice_vcs) == 2


class TestStatusListStorage:
    @pytest.mark.asyncio
    async def test_save_and_get(self, db_session):
        storage = StatusListStorage(db_session)
        await storage.save("list-1", "revocation", "AAAA")
        record = await storage.get("list-1")
        assert record is not None
        assert record.list_id == "list-1"
        assert record.bitstring == "AAAA"

    @pytest.mark.asyncio
    async def test_update(self, db_session):
        storage = StatusListStorage(db_session)
        await storage.save("list-1", "revocation", "AAAA")
        assert await storage.update_bitstring("list-1", "BBBB") is True
        record = await storage.get("list-1")
        assert record.bitstring == "BBBB"
        assert await storage.update_bitstring("nonexistent", "CCCC") is False

    @pytest.mark.asyncio
    async def test_delete(self, db_session):
        storage = StatusListStorage(db_session)
        await storage.save("list-1", "revocation", "AAAA")
        assert await storage.delete("list-1") is True
        assert await storage.delete("nonexistent") is False


class TestWalletCredentialStorage:
    @pytest.mark.asyncio
    async def test_save_and_get(self, db_session):
        storage = WalletCredentialStorage(db_session)
        await storage.save("did:key:alice", "vc-1", {"id": "vc-1"}, "issuer", "IdentityCredential")
        doc = await storage.get("did:key:alice", "vc-1")
        assert doc is not None

    @pytest.mark.asyncio
    async def test_list_by_wallet(self, db_session):
        storage = WalletCredentialStorage(db_session)
        await storage.save("did:key:alice", "vc-1", {"id": "vc-1"}, "issuer", "A")
        await storage.save("did:key:alice", "vc-2", {"id": "vc-2"}, "issuer", "B")
        await storage.save("did:key:bob", "vc-3", {"id": "vc-3"}, "issuer", "A")
        alice = await storage.list_by_wallet("did:key:alice")
        assert len(alice) == 2

    @pytest.mark.asyncio
    async def test_delete(self, db_session):
        storage = WalletCredentialStorage(db_session)
        await storage.save("did:key:alice", "vc-1", {"id": "vc-1"}, "issuer", "A")
        assert await storage.delete("did:key:alice", "vc-1") is True
        assert await storage.get("did:key:alice", "vc-1") is None

    @pytest.mark.asyncio
    async def test_search(self, db_session):
        storage = WalletCredentialStorage(db_session)
        await storage.save("did:key:alice", "vc-1", {"id": "vc-1"}, "did:key:issuer1", "IdentityCredential")
        await storage.save("did:key:alice", "vc-2", {"id": "vc-2"}, "did:key:issuer2", "DriversLicense")
        results = await storage.search("did:key:alice", "Identity")
        assert len(results) == 1


class TestDIDCommMessageStorage:
    @pytest.mark.asyncio
    async def test_save_and_inbox(self, db_session):
        storage = DIDCommMessageStorage(db_session)
        await storage.save("msg-1", "did:key:alice", "did:key:bob", "basicmessage", {"text": "hello"})
        inbox = await storage.get_inbox("did:key:bob")
        assert len(inbox) == 1
        assert inbox[0].message_id == "msg-1"

    @pytest.mark.asyncio
    async def test_sent(self, db_session):
        storage = DIDCommMessageStorage(db_session)
        await storage.save("msg-1", "did:key:alice", "did:key:bob", "basicmessage", {"text": "hello"})
        sent = await storage.get_sent("did:key:alice")
        assert len(sent) == 1

    @pytest.mark.asyncio
    async def test_mark_read(self, db_session):
        storage = DIDCommMessageStorage(db_session)
        await storage.save("msg-1", "alice", "bob", "basicmessage", {"text": "hi"})
        assert await storage.mark_read("msg-1") is True
        inbox = await storage.get_inbox("bob")
        assert inbox[0].is_read is True

    @pytest.mark.asyncio
    async def test_delete(self, db_session):
        storage = DIDCommMessageStorage(db_session)
        await storage.save("msg-1", "alice", "bob", "basicmessage", {"text": "hi"})
        assert await storage.delete("msg-1") is True
        assert await storage.delete("nonexistent") is False


class TestRevocationEventStorage:
    @pytest.mark.asyncio
    async def test_save_and_history(self, db_session):
        storage = RevocationEventStorage(db_session)
        await storage.save("evt-1", "vc-1", "did:key:alice", "credential.revoked", "Lost")
        await storage.save("evt-2", "vc-2", "did:key:bob", "credential.suspended", "Review")
        history = await storage.get_history()
        assert len(history) == 2

    @pytest.mark.asyncio
    async def test_filter_by_credential(self, db_session):
        storage = RevocationEventStorage(db_session)
        await storage.save("evt-1", "vc-1", "did:key:alice", "credential.revoked")
        await storage.save("evt-2", "vc-1", "did:key:alice", "credential.suspended")
        await storage.save("evt-3", "vc-2", "did:key:bob", "credential.revoked")
        history = await storage.get_history(credential_id="vc-1")
        assert len(history) == 2
