import base64
import json

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

from services.sd_jwt import SDJWTBuilder, SDJWTIssuer, SDJWTVerifier
from services.sd_jwt.api import router as sd_jwt_router


@pytest.fixture
def issuer():
    return SDJWTIssuer(issuer_id="https://snisid.ht/issuer")


@pytest.fixture
def verifier():
    return SDJWTVerifier(trusted_issuers={"https://snisid.ht/issuer"})


class TestSDJWTIssuer:
    def test_issue_returns_jwt_and_disclosures(self, issuer):
        sd_jwt, disclosures = issuer.issue(
            subject="user-1",
            disclosed_claims={"name": "Alice"},
            sd_claims={"email": "alice@example.com", "age": 30},
        )
        assert sd_jwt.count(".") == 2
        assert len(disclosures) == 2
        assert sd_jwt.startswith("eyJ")

    def test_issue_no_sd_claims(self, issuer):
        sd_jwt, disclosures = issuer.issue(
            subject="user-1",
            disclosed_claims={"name": "Alice"},
        )
        assert len(disclosures) == 0
        assert sd_jwt.count(".") == 2

    def test_issue_has_sd_in_payload(self, issuer):
        sd_jwt, disclosures = issuer.issue(
            subject="user-1",
            disclosed_claims={"name": "Alice"},
            sd_claims={"email": "alice@example.com"},
        )
        payload_b64 = sd_jwt.split(".")[1]
        import base64, json
        payload = json.loads(base64.urlsafe_b64decode(payload_b64 + "=="))
        assert "_sd" in payload
        assert len(payload["_sd"]) == 1
        assert payload["name"] == "Alice"
        assert "email" not in payload

    def test_issue_with_expiration(self, issuer):
        sd_jwt, _ = issuer.issue(
            subject="user-1",
            disclosed_claims={},
            sd_claims={"data": "secret"},
            expiration_seconds=0,
        )
        payload_b64 = sd_jwt.split(".")[1]
        import base64, json
        payload = json.loads(base64.urlsafe_b64decode(payload_b64 + "=="))
        assert payload["exp"] == payload["iat"]


class TestSDJWTVerifier:
    def test_verify_valid_sd_jwt(self, issuer, verifier):
        sd_jwt, disclosures = issuer.issue(
            subject="user-1",
            disclosed_claims={"name": "Alice"},
            sd_claims={"email": "alice@example.com"},
        )
        result = verifier.verify(sd_jwt, disclosures)
        assert result["email"] == "alice@example.com"
        assert result["name"] == "Alice"

    def test_verify_partial_disclosure(self, issuer, verifier):
        sd_jwt, all_disclosures = issuer.issue(
            subject="user-1",
            disclosed_claims={"name": "Alice"},
            sd_claims={"email": "alice@example.com", "age": 30},
        )
        import base64, json
        partial = [
            d for d in all_disclosures
            if json.loads(base64.urlsafe_b64decode(d + "=="))[1] == "email"
        ]
        result = verifier.verify(sd_jwt, partial)
        assert result["email"] == "alice@example.com"
        assert "age" not in result

    def test_verify_expired(self, issuer, verifier):
        sd_jwt, disclosures = issuer.issue(
            subject="user-1",
            disclosed_claims={},
            sd_claims={"data": "x"},
            expiration_seconds=-1,
        )
        with pytest.raises(ValueError, match="has expired"):
            verifier.verify(sd_jwt, disclosures)

    def test_verify_invalid_signature(self, verifier):
        with pytest.raises(ValueError, match="Invalid signature"):
            verifier.verify(
                "eyJhbGciOiJIUzI1NiIsInR5cCI6InNkK2p3dCJ9.eyJf"
                "c2QiOlsibHNNd3BzdkhOdmkzQy1ISTlwV0JJUT09Il0sIm"
                "5hbWUiOiJCb2IifQ.fakesig",
                ["W1wic2FsdC0xXCIsXCJlbWFpbFwiLFwiYm9iQGV4YW1wbGUuY29tXCJd"],
            )

    def test_verify_requires_claims(self, issuer, verifier):
        sd_jwt, disclosures = issuer.issue(
            subject="user-1",
            disclosed_claims={},
            sd_claims={"email": "alice@example.com"},
        )
        with pytest.raises(ValueError, match="Required claims not disclosed"):
            verifier.verify(
                sd_jwt,
                [],  # disclose nothing
                required_claims=["email"],
            )

    def test_verify_unknown_issuer(self):
        verifier = SDJWTVerifier(trusted_issuers={"other"})
        issuer = SDJWTIssuer(issuer_id="https://snisid.ht/issuer")
        sd_jwt, disclosures = issuer.issue(
            subject="user-1",
            disclosed_claims={},
            sd_claims={"email": "alice@example.com"},
        )
        with pytest.raises(ValueError, match="Untrusted issuer"):
            verifier.verify(sd_jwt, disclosures)

    def test_verify_wrong_header_type(self, verifier):
        with pytest.raises(ValueError, match="Not an SD-JWT"):
            verifier.verify("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjMifQ.fake", [])

    def test_verify_invalid_format(self, verifier):
        with pytest.raises(ValueError, match="Invalid SD-JWT format"):
            verifier.verify("not.a.jwt.payload", [])

    def test_verify_fake_disclosure(self, issuer, verifier):
        sd_jwt, disclosures = issuer.issue(
            subject="user-1",
            disclosed_claims={"name": "Alice"},
            sd_claims={"email": "alice@example.com"},
        )
        import base64, json
        fake_disclosure = base64.urlsafe_b64encode(
            json.dumps(["salt-fake", "fake-claim", "fake-value"]).encode()
        ).rstrip(b"=").decode()
        with pytest.raises(ValueError, match="Disclosure not in _sd"):
            disclosures_with_fake = disclosures + [fake_disclosure]
            verifier.verify(sd_jwt, disclosures_with_fake)


class TestSDJWTBuilder:
    def test_create_presentation_selects_claims(self, issuer):
        sd_jwt, all_disclosures = issuer.issue(
            subject="user-1",
            disclosed_claims={"name": "Alice"},
            sd_claims={"email": "alice@example.com", "phone": "123-456-7890"},
        )
        sd_jwt_presented, selected = SDJWTBuilder.create_presentation(
            sd_jwt, all_disclosures, disclose=["email"]
        )
        assert sd_jwt_presented == sd_jwt
        assert len(selected) == 1

        import base64, json
        decoded = json.loads(base64.urlsafe_b64decode(selected[0] + "=="))
        assert decoded[1] == "email"

    def test_create_presentation_empty_selection(self, issuer):
        sd_jwt, all_disclosures = issuer.issue(
            subject="user-1",
            disclosed_claims={},
            sd_claims={"a": "1", "b": "2"},
        )
        _, selected = SDJWTBuilder.create_presentation(
            sd_jwt, all_disclosures, disclose=[]
        )
        assert len(selected) == 0

    def test_create_presentation_select_all(self, issuer):
        sd_jwt, all_disclosures = issuer.issue(
            subject="user-1",
            disclosed_claims={},
            sd_claims={"a": "1", "b": "2"},
        )
        _, selected = SDJWTBuilder.create_presentation(
            sd_jwt, all_disclosures, disclose=["a", "b"]
        )
        assert len(selected) == 2


class TestSDJWTApi:
    @pytest.fixture
    def client(self):
        app = FastAPI()
        app.include_router(sd_jwt_router)
        return TestClient(app)

    def test_issue_endpoint(self, client):
        resp = client.post("/sd-jwt/issue", json={
            "subject": "user-1",
            "disclosed_claims": {"name": "Alice"},
            "sd_claims": {"email": "alice@example.com"},
        })
        assert resp.status_code == 200
        data = resp.json()
        assert "sd_jwt" in data
        assert len(data["disclosures"]) == 1

    def test_verify_endpoint(self, client):
        issue_resp = client.post("/sd-jwt/issue", json={
            "subject": "user-1",
            "disclosed_claims": {"name": "Alice"},
            "sd_claims": {"email": "alice@example.com"},
        })
        data = issue_resp.json()
        resp = client.post("/sd-jwt/verify", json={
            "sd_jwt": data["sd_jwt"],
            "disclosures": data["disclosures"],
        })
        assert resp.status_code == 200
        assert resp.json()["claims"]["email"] == "alice@example.com"

    def test_verify_endpoint_invalid(self, client):
        resp = client.post("/sd-jwt/verify", json={
            "sd_jwt": "invalid.jwt.here",
            "disclosures": [],
        })
        assert resp.status_code == 400

    def test_present_endpoint(self, client):
        issue_resp = client.post("/sd-jwt/issue", json={
            "subject": "user-1",
            "disclosed_claims": {},
            "sd_claims": {"email": "alice@example.com", "phone": "123456"},
        })
        data = issue_resp.json()
        resp = client.post("/sd-jwt/present", json={
            "sd_jwt": data["sd_jwt"],
            "all_disclosures": data["disclosures"],
            "disclose": ["email"],
        })
        assert resp.status_code == 200
        assert len(resp.json()["disclosures"]) == 1
