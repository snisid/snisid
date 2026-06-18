import pytest

from services.chapi import (
    CHAPIMediator,
    CHAPIStoreRequest,
    CHAPIGetRequest,
    CHAPIResponse,
)
from services.wallet import Wallet


@pytest.fixture
def wallet():
    w = Wallet()
    w.store({
        "type": ["VerifiableCredential", "SNISIDIdentityCredential"],
        "issuer": "did:key:issuer",
        "credentialSubject": {"national_id": "SN-001", "first_name": "Alice"},
    })
    w.store({
        "type": ["VerifiableCredential", "DiplomaCredential"],
        "issuer": "did:key:uni",
        "credentialSubject": {"degree": "BSc", "name": "Alice"},
    })
    return w


@pytest.fixture
def mediator(wallet):
    return CHAPIMediator(wallet=wallet)


class TestCHAPIModels:
    def test_store_request_to_dict(self):
        req = CHAPIStoreRequest(credential={"type": ["VerifiableCredential"]})
        d = req.to_dict()
        assert d["protocol"] == "vc"
        assert "credential_id" in d

    def test_get_request_to_dict(self):
        req = CHAPIGetRequest(query=[{"type": "VerifiableCredential"}])
        d = req.to_dict()
        assert d["protocol"] == "vc"
        assert len(d["query"]) == 1

    def test_response_with_data(self):
        resp = CHAPIResponse(data=[{"key": "value"}])
        d = resp.to_dict()
        assert d["data"] == [{"key": "value"}]
        assert "error" not in d

    def test_response_with_error(self):
        resp = CHAPIResponse(error="something went wrong")
        d = resp.to_dict()
        assert d["error"] == "something went wrong"
        assert "data" not in d


class TestCHAPIMediator:
    def test_handler_registration(self, mediator):
        resp = mediator.handler_registration("https://example.com")
        assert resp.data["handler"] == "snisid-wallet"
        assert resp.data["origin"] == "https://example.com"
        assert "store" in resp.data["capabilities"]

    def test_handle_store(self, mediator):
        req = CHAPIStoreRequest(credential={
            "type": ["VerifiableCredential"],
            "credentialSubject": {"name": "New User"},
        })
        resp = mediator.handle_store(req)
        assert resp.error is None
        assert resp.data["status"] == "stored"

    def test_handle_get_all(self, mediator):
        req = CHAPIGetRequest(query=[{"type": "VerifiableCredential"}])
        resp = mediator.handle_get(req)
        assert resp.error is None
        vp = resp.data
        assert vp["holder"] == mediator.wallet.did
        assert len(vp["verifiableCredential"]) == 2

    def test_handle_get_by_type(self, mediator):
        req = CHAPIGetRequest(query=[{"type": "SNISIDIdentityCredential"}])
        resp = mediator.handle_get(req)
        assert resp.error is None
        assert len(resp.data["verifiableCredential"]) == 1

    def test_handle_get_with_frame(self, mediator):
        req = CHAPIGetRequest(query=[{
            "type": "VerifiableCredential",
            "credentialFrame": {"issuer": "did:key:issuer"},
        }])
        resp = mediator.handle_get(req)
        assert resp.error is None
        vcs = resp.data["verifiableCredential"]
        assert len(vcs) == 1
        assert vcs[0]["issuer"] == "did:key:issuer"

    def test_handle_get_no_match(self, mediator):
        req = CHAPIGetRequest(query=[{
            "type": "NonexistentCredential",
        }])
        resp = mediator.handle_get(req)
        assert resp.error is None
        assert len(resp.data["verifiableCredential"]) == 0

    def test_handle_store_increments_count(self, mediator):
        before = mediator.wallet.count()
        mediator.handle_store(CHAPIStoreRequest(credential={
            "type": ["VerifiableCredential"],
            "credentialSubject": {"name": "Another"},
        }))
        assert mediator.wallet.count() == before + 1


class TestCHAPIMediatorFrameMatching:
    def test_matches_frame_exact(self, mediator):
        cred = {"type": "VerifiableCredential", "issuer": "did:key:test"}
        assert mediator._matches_frame(cred, {"issuer": "did:key:test"}) is True
        assert mediator._matches_frame(cred, {"issuer": "did:key:other"}) is False

    def test_matches_frame_nested(self, mediator):
        cred = {"credentialSubject": {"name": "Alice", "age": 30}}
        assert mediator._matches_frame(cred, {"credentialSubject": {"name": "Alice"}}) is True
        assert mediator._matches_frame(cred, {"credentialSubject": {"name": "Bob"}}) is False

    def test_matches_frame_empty(self, mediator):
        assert mediator._matches_frame({"any": "value"}, {}) is True

    def test_matches_frame_list(self, mediator):
        cred = {"type": ["A", "B", "C"]}
        assert mediator._matches_frame(cred, {"type": ["A", "B"]}) is True
        assert mediator._matches_frame(cred, {"type": ["X", "Y"]}) is False
