import pytest

from services.credential_flow import CredentialFlow
from services.vc import VCStatus, VerifiableCredential
from services.vc.issuer import VCIssuer
from services.vc.verifier import VCVerifier


@pytest.fixture
def issuer_did():
    return "did:key:issuer123"


@pytest.fixture
def holder_did():
    return "did:snisid:mainnet:citizen-001"


@pytest.fixture
def flow(issuer_did):
    issuer = VCIssuer(issuer_id=issuer_did)
    verifier = VCVerifier(trusted_issuers=[issuer_did])
    return CredentialFlow(issuer=issuer, verifier=verifier)


class TestCredentialOffer:
    def test_create_offer(self, flow, issuer_did, holder_did):
        offer = flow.create_offer(issuer_did=issuer_did, holder_did=holder_did)
        assert offer.offer_id is not None
        assert offer.issuer_did == issuer_did
        assert offer.holder_did == holder_did
        assert offer.status == "pending"
        assert "national_id" in offer.claims_requested

    def test_offer_to_dict(self, flow, issuer_did, holder_did):
        offer = flow.create_offer(issuer_did=issuer_did, holder_did=holder_did)
        d = offer.to_dict()
        assert d["issuer_did"] == issuer_did
        assert d["holder_did"] == holder_did

    def test_send_offer_returns_packed_message(self, flow, issuer_did, holder_did):
        offer = flow.create_offer(issuer_did=issuer_did, holder_did=holder_did)
        packed = flow.send_offer(offer)
        assert "ciphertext" in packed or "signatures" in packed


class TestCredentialRequest:
    def test_receive_request(self, flow, issuer_did, holder_did):
        offer = flow.create_offer(issuer_did=issuer_did, holder_did=holder_did)
        body = flow.build_request(offer.offer_id, holder_did)
        packed = flow.pack_and_send_request(body, holder_did, issuer_did)
        req = flow.receive_request(packed)
        assert req.offer_id == offer.offer_id
        assert req.claims["national_id"] == "SN-CITIZEN-001"

    def test_issue_from_request_returns_vc(self, flow, issuer_did, holder_did):
        offer = flow.create_offer(issuer_did=issuer_did, holder_did=holder_did)
        body = flow.build_request(offer.offer_id, holder_did)
        packed = flow.pack_and_send_request(body, holder_did, issuer_did)
        req = flow.receive_request(packed)
        vc_data = flow.issue_from_request(req)
        assert vc_data is not None
        assert vc_data["credentialSubject"]["national_id"] == "SN-CITIZEN-001"
        assert offer.status == "fulfilled"

    def test_issue_from_request_invalid_offer(self, flow, issuer_did, holder_did):
        body = flow.build_request("nonexistent", holder_did)
        packed = flow.pack_and_send_request(body, holder_did, issuer_did)
        req = flow.receive_request(packed)
        assert flow.issue_from_request(req) is None

    def test_full_flow(self, flow, issuer_did, holder_did):
        offer = flow.create_offer(issuer_did=issuer_did, holder_did=holder_did)
        body = flow.build_request(offer.offer_id, holder_did)
        packed = flow.pack_and_send_request(body, holder_did, issuer_did)
        req = flow.receive_request(packed)
        vc_data = flow.issue_from_request(req)
        assert vc_data is not None

        packed_vc = flow.send_credential(vc_data, issuer_did, holder_did)
        assert "ciphertext" in packed_vc or "signatures" in packed_vc

        received = flow.receive_credential(packed_vc)
        assert received is not None
        assert received["credentialSubject"]["national_id"] == "SN-CITIZEN-001"

    def test_issued_credential_verifiable(self, flow, issuer_did, holder_did):
        offer = flow.create_offer(issuer_did=issuer_did, holder_did=holder_did)
        body = flow.build_request(offer.offer_id, holder_did)
        packed = flow.pack_and_send_request(body, holder_did, issuer_did)
        req = flow.receive_request(packed)
        vc_data = flow.issue_from_request(req)

        vc = VerifiableCredential(**vc_data)
        flow._verifier.register_status(vc.id, VCStatus.ACTIVE)
        result = flow._verifier.verify_credential(vc)
        assert result.valid is True
