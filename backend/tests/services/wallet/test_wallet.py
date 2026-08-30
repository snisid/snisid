import pytest

from services.wallet import Wallet, CredentialRecord


@pytest.fixture
def wallet():
    return Wallet()


@pytest.fixture
def sample_vc():
    return {
        "@context": ["https://www.w3.org/ns/credentials/v2"],
        "id": "http://example.gov/credentials/1",
        "type": ["VerifiableCredential", "SNISIDIdentityCredential"],
        "issuer": "did:key:issuer-001",
        "issuanceDate": "2026-01-01T00:00:00Z",
        "credentialSubject": {
            "national_id": "SN-CITIZEN-001",
            "first_name": "Alice",
            "last_name": "Citizen",
            "date_of_birth": "1990-01-15",
        },
    }


class TestWalletCore:
    def test_wallet_has_did(self, wallet):
        assert wallet.did.startswith("did:")
        assert len(wallet.did) > 10

    def test_did_document(self, wallet):
        doc = wallet.did_document
        assert doc["id"] == wallet.did
        assert "verification_method" in doc

    def test_store_returns_record(self, wallet, sample_vc):
        record = wallet.store(sample_vc, label="My ID")
        assert record.id is not None
        assert record.credential["id"] == "http://example.gov/credentials/1"
        assert record.label == "My ID"
        assert record.issuer_did == "did:key:issuer-001"

    def test_get_credential(self, wallet, sample_vc):
        record = wallet.store(sample_vc)
        retrieved = wallet.get(record.id)
        assert retrieved is not None
        assert retrieved.id == record.id

    def test_get_nonexistent(self, wallet):
        assert wallet.get("nonexistent") is None

    def test_list_all(self, wallet, sample_vc):
        wallet.store(sample_vc, label="VC 1")
        wallet.store(sample_vc, label="VC 2")
        assert len(wallet.list()) == 2

    def test_list_by_type(self, wallet, sample_vc):
        wallet.store(sample_vc)
        wallet.store({
            "type": ["VerifiableCredential", "DiplomaCredential"],
            "credentialSubject": {"degree": "BSc"},
        }, label="Diploma")
        assert len(wallet.list("SNISIDIdentityCredential")) == 1
        assert len(wallet.list("DiplomaCredential")) == 1
        assert len(wallet.list()) == 2

    def test_delete_credential(self, wallet, sample_vc):
        record = wallet.store(sample_vc)
        assert wallet.count() == 1
        assert wallet.delete(record.id) is True
        assert wallet.count() == 0
        assert wallet.delete(record.id) is False

    def test_count(self, wallet, sample_vc):
        assert wallet.count() == 0
        wallet.store(sample_vc)
        assert wallet.count() == 1


class TestWalletPresentations:
    def test_create_presentation_all(self, wallet, sample_vc):
        wallet.store(sample_vc)
        vp = wallet.create_presentation()
        assert vp.holder == wallet.did
        assert len(vp.verifiable_credential) == 1

    def test_create_presentation_by_ids(self, wallet, sample_vc):
        r1 = wallet.store(sample_vc)
        r2 = wallet.store({
            "type": ["VerifiableCredential"],
            "credentialSubject": {"name": "Bob"},
        })
        vp = wallet.create_presentation(credential_ids=[r1.id])
        assert len(vp.verifiable_credential) == 1

    def test_create_presentation_by_type(self, wallet, sample_vc):
        wallet.store(sample_vc)
        wallet.store({
            "type": ["VerifiableCredential", "DiplomaCredential"],
            "credentialSubject": {"degree": "BSc"},
        })
        vp = wallet.create_presentation(credential_type="SNISIDIdentityCredential")
        assert len(vp.verifiable_credential) == 1

    def test_create_empty_presentation(self, wallet):
        vp = wallet.create_presentation()
        assert vp.holder == wallet.did
        assert len(vp.verifiable_credential) == 0

    def test_verify_own_presentation(self, wallet, sample_vc):
        wallet.store(sample_vc)
        vp = wallet.create_presentation()
        assert wallet.verify_presentation(vp) is True

    def test_verify_tampered_presentation(self, wallet, sample_vc):
        wallet.store(sample_vc)
        vp = wallet.create_presentation()
        vp.holder = "did:key:attacker"
        assert wallet.verify_presentation(vp) is False


class TestWalletSearch:
    def test_search_by_label(self, wallet, sample_vc):
        wallet.store(sample_vc, label="National ID Card")
        results = wallet.search("national")
        assert len(results) == 1

    def test_search_by_subject(self, wallet, sample_vc):
        wallet.store(sample_vc)
        results = wallet.search("alice")
        assert len(results) == 1

    def test_search_by_issuer(self, wallet, sample_vc):
        wallet.store(sample_vc)
        results = wallet.search("issuer-001")
        assert len(results) == 1

    def test_search_no_match(self, wallet, sample_vc):
        wallet.store(sample_vc)
        results = wallet.search("nonexistent")
        assert len(results) == 0

    def test_get_by_issuer(self, wallet, sample_vc):
        wallet.store(sample_vc)
        wallet.store({
            "type": ["VerifiableCredential"],
            "issuer": "did:key:issuer-002",
            "credentialSubject": {"name": "Bob"},
        })
        assert len(wallet.get_credentials_by_issuer("did:key:issuer-001")) == 1
        assert len(wallet.get_credentials_by_issuer("did:key:issuer-002")) == 1
        assert len(wallet.get_credentials_by_issuer("did:key:unknown")) == 0


class TestWalletDIDComm:
    def test_send_via_didcomm(self, wallet, sample_vc):
        record = wallet.store(sample_vc)
        packed = wallet.send_via_didcomm(record.id, "did:key:recipient")
        assert "ciphertext" in packed or "signatures" in packed

    def test_send_nonexistent_credential(self, wallet):
        with pytest.raises(ValueError, match="Credential not found"):
            wallet.send_via_didcomm("nonexistent", "did:key:recipient")

    def test_receive_via_didcomm(self, wallet, sample_vc):
        record = wallet.store(sample_vc)
        packed = wallet.send_via_didcomm(record.id, "did:key:recipient")

        receiver = Wallet()
        stored = receiver.receive_via_didcomm(packed, label="Received via DIDComm")
        assert stored is not None
        assert receiver.count() == 1
        assert stored.label == "Received via DIDComm"
        assert stored.credential["credentialSubject"]["national_id"] == "SN-CITIZEN-001"

    def test_receive_empty_message(self, wallet):
        messenger = __import__("services.didcomm", fromlist=["DIDCommMessenger"]).DIDCommMessenger()
        msg = __import__("services.didcomm", fromlist=["DIDCommMessage"]).DIDCommMessage(
            id="empty-msg",
            type="https://didcomm.org/test/1.0/test",
            body={},
            from_did="did:key:sender",
            to_did=wallet.did,
        )
        packed = messenger.send(msg, "did:key:sender", wallet.did)
        assert wallet.receive_via_didcomm(packed) is None

    def test_send_message(self, wallet):
        packed = wallet.send_message(
            message_type="https://didcomm.org/trust-ping/2.0/ping",
            body={"comment": "hello"},
            to_did="did:key:friend",
        )
        assert "ciphertext" in packed or "signatures" in packed

    def test_roundtrip_multiple_credentials(self, wallet):
        r1 = wallet.store({"type": ["VerifiableCredential"], "credentialSubject": {"name": "A"}})
        r2 = wallet.store({"type": ["VerifiableCredential"], "credentialSubject": {"name": "B"}})

        packed1 = wallet.send_via_didcomm(r1.id, "did:key:bob")
        packed2 = wallet.send_via_didcomm(r2.id, "did:key:bob")

        receiver = Wallet()
        received1 = receiver.receive_via_didcomm(packed1)
        received2 = receiver.receive_via_didcomm(packed2)
        assert receiver.count() == 2


class TestWalletExportImport:
    def test_export_all(self, wallet, sample_vc):
        wallet.store(sample_vc, label="Export Test")
        exported = wallet.export_all()
        assert len(exported) == 1
        assert exported[0]["label"] == "Export Test"

    def test_import_credentials(self, wallet):
        data = [
            {
                "id": "imported-1",
                "credential": {"type": ["VerifiableCredential"], "credentialSubject": {"name": "Import"}},
                "issuer_did": "did:key:issuer",
                "credential_type": ["VerifiableCredential"],
                "issued_at": "2026-01-01T00:00:00Z",
                "label": "Imported VC",
            }
        ]
        wallet.import_credentials(data)
        assert wallet.count() == 1
        record = wallet.get("imported-1")
        assert record is not None
        assert record.label == "Imported VC"
