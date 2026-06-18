"""E2E integration test: full SSI flow across all modules."""

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

from services.did import DIDManager, DIDMethod
from services.didcomm import DIDCommMessage, DIDCommMessenger
from services.status_list import StatusListManager
from services.vc import VCStatus
from services.vc.issuer import VCIssuer
from services.vc.verifier import VCVerifier
from services.vp import VPIssuer, VerifiablePresentation


@pytest.fixture
def issuer_did():
    mgr = DIDManager()
    doc = mgr.create(DIDMethod.KEY)
    return doc.id


@pytest.fixture
def holder_did():
    mgr = DIDManager()
    doc = mgr.create(DIDMethod.SNISID, identifier="citizen-001")
    return doc.id


def test_full_ssi_flow(issuer_did, holder_did):
    """
    Complete SSI flow:
    DID -> VC Issuance -> StatusList2021 -> VP -> DIDComm
    """
    # 1. Create issuer with StatusList2021-backed revocation
    issuer = VCIssuer(issuer_id=issuer_did)
    verifier = VCVerifier(trusted_issuers=[issuer_did])
    status_mgr = StatusListManager(issuer_id=issuer_did)
    vp_issuer = VPIssuer()

    # 2. Issue a VC
    vc = issuer.issue_identity_credential(
        subject_id=holder_did,
        national_id="SN-E2E-001",
        first_name="E2E",
        last_name="Test",
        date_of_birth="1990-01-01",
        gender="female",
        nationality="HTI",
    )
    assert vc.id.startswith("urn:uuid:")
    assert vc.proof is not None
    assert vc.credentialStatus is not None
    assert vc.credentialStatus["type"] == "StatusList2021Entry"
    assert issuer.get_credential_status(vc.id) == VCStatus.ACTIVE

    # 3. Verify the VC
    verifier.register_status(vc.id, VCStatus.ACTIVE)
    result = verifier.verify_credential(vc)
    assert result.valid is True, f"VC verification failed: {result.errors}"

    # 4. Revoke via StatusList2021
    assert issuer.revoke_credential(vc.id) is True
    assert issuer.get_credential_status(vc.id) == VCStatus.REVOKED

    # 5. Verify revocation is detected
    verifier.register_status(vc.id, VCStatus.REVOKED)
    result = verifier.verify_credential(vc)
    assert result.valid is False
    assert any("revoked" in e for e in result.errors)

    # 6. Create a VP with an unrevoked VC (re-issue since revoked)
    vc2 = issuer.issue_identity_credential(
        subject_id=holder_did,
        national_id="SN-E2E-002",
        first_name="E2E",
        last_name="Test",
        date_of_birth="1990-01-01",
        gender="female",
        nationality="HTI",
    )
    vp = vp_issuer.create_presentation(
        holder_did=holder_did,
        verifiable_credentials=[vc2.model_dump()],
        verification_method=f"{holder_did}#key-1",
    )
    assert vp.proof is not None
    assert vp.proof["proofPurpose"] == "authentication"

    # 7. Verify the VP
    assert vp_issuer.verify_presentation(vp) is True
    assert vp.holder == holder_did
    assert len(vp.verifiable_credential) == 1

    # 8. Tampered VP fails
    vp.holder = "did:key:attacker"
    assert vp_issuer.verify_presentation(vp) is False

    # 9. DIDComm transport
    messenger = DIDCommMessenger()
    ping = messenger.create_trust_ping(
        from_did=holder_did, to_did=issuer_did
    )
    packed = messenger.send(ping, holder_did, issuer_did)

    # DIDComm with encrypted payload
    assert "ciphertext" in packed or "signatures" in packed

    received = messenger.receive(packed)
    assert received.id == ping.id
    assert received.type == "https://didcomm.org/trust-ping/2.0/ping"


def test_status_list_vc_integration(issuer_did, holder_did):
    """Verify StatusList2021 entry matches VC credentialStatus."""
    issuer = VCIssuer(issuer_id=issuer_did)
    vc = issuer.issue_identity_credential(
        subject_id=holder_did,
        national_id="SN-INT-001",
        first_name="Integration",
        last_name="Test",
        date_of_birth="1995-05-15",
        gender="male",
        nationality="HTI",
    )
    cs = vc.credentialStatus
    assert cs["type"] == "StatusList2021Entry"
    assert "statusListIndex" in cs
    assert "statusListCredential" in cs
    assert cs["statusPurpose"] == "revocation"
    assert cs["statusListIndex"].isdigit()


def test_vp_with_verified_vc(issuer_did, holder_did):
    """VP that passes verifier checks."""
    issuer = VCIssuer(issuer_id=issuer_did)
    verifier = VCVerifier(trusted_issuers=[issuer_did])
    vp_issuer = VPIssuer()

    vc = issuer.issue_identity_credential(
        subject_id=holder_did,
        national_id="SN-VP-001",
        first_name="VP",
        last_name="Test",
        date_of_birth="2000-12-25",
        gender="other",
        nationality="HTI",
    )
    verifier.register_status(vc.id, VCStatus.ACTIVE)
    vc_result = verifier.verify_credential(vc)
    assert vc_result.valid is True

    vp = vp_issuer.create_presentation(
        holder_did=holder_did,
        verifiable_credentials=[vc.model_dump()],
    )
    assert vp_issuer.verify_presentation(vp) is True

    vp_dict = vp.to_dict()
    restored = VerifiablePresentation.from_dict(vp_dict)
    assert restored.holder == holder_did
    assert len(restored.verifiable_credential) == 1
    assert restored.proof is not None


def test_did_vc_vp_presentation_flow(issuer_did, holder_did):
    """End-to-end: resolve DID -> issue VC -> create VP -> verify."""
    from services.did import resolve_did

    # Resolve both DIDs
    issuer_doc = resolve_did(issuer_did)
    assert issuer_doc.id == issuer_did
    assert len(issuer_doc.verification_method) >= 1

    holder_doc = resolve_did(holder_did)
    assert holder_doc.id == holder_did

    # Issue VC from resolved issuer
    issuer = VCIssuer(issuer_id=issuer_did)
    vc = issuer.issue_identity_credential(
        subject_id=holder_did,
        national_id="SN-E2E-FULL-001",
        first_name="Full",
        last_name="Flow",
        date_of_birth="1988-07-20",
        gender="male",
        nationality="HTI",
    )
    assert vc.issuer == issuer_did

    # Create VP
    vp_issuer = VPIssuer()
    vp = vp_issuer.create_presentation(
        holder_did=holder_did,
        verifiable_credentials=[vc.model_dump()],
    )
    assert vp_issuer.verify_presentation(vp) is True

    # DIDComm the VP
    messenger = DIDCommMessenger()
    vp_message = DIDCommMessage(
        id="vp-transfer-1",
        type="https://didcomm.org/verifiable-presentation/1.0/propose-presentation",
        body={"presentation": vp.to_dict()},
        from_did=holder_did,
        to_did=issuer_did,
    )
    packed = messenger.send(vp_message, holder_did, issuer_did)
    received = messenger.receive(packed)

    assert received.body["presentation"]["holder"] == holder_did
    received_vp = VerifiablePresentation.from_dict(received.body["presentation"])
    assert vp_issuer.verify_presentation(received_vp) is True
