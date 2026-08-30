import pytest

from services.siopv2 import (
    InputDescriptor,
    PresentationDefinition,
    PresentationSubmission,
    SelfIssuedIDToken,
    OIDC4VPRequest,
    OIDC4VPResponse,
    SIOPWallet,
    OIDC4VPVerifier,
)


@pytest.fixture
def wallet():
    w = SIOPWallet()
    w.store_credential({
        "@context": ["https://www.w3.org/ns/credentials/v2"],
        "id": "http://example.gov/credentials/1",
        "type": ["VerifiableCredential", "SNISIDIdentityCredential"],
        "issuer": "did:key:issuer",
        "credentialSubject": {"national_id": "SN-001"},
    })
    return w


@pytest.fixture
def verifier():
    return OIDC4VPVerifier()


@pytest.fixture
def presentation_definition():
    desc = InputDescriptor(
        id="identity-credential",
        name="SNISID Identity Credential",
        purpose="Verify identity",
    )
    return PresentationDefinition(
        id="pd-1",
        name="Identity Check",
        input_descriptors=[desc],
    )


class TestPresentationDefinition:
    def test_input_descriptor_to_dict(self):
        desc = InputDescriptor(id="test", name="Test", purpose="Testing")
        d = desc.to_dict()
        assert d["id"] == "test"
        assert "schema" in d
        assert "constraints" in d

    def test_presentation_definition_to_dict(self, presentation_definition):
        d = presentation_definition.to_dict()
        assert d["id"] == "pd-1"
        assert len(d["input_descriptors"]) == 1

    def test_presentation_submission_to_dict(self):
        sub = PresentationSubmission(
            definition_id="pd-1",
            descriptor_map=[{"id": "test", "format": "vp", "path": "$.vp_token[0]"}],
        )
        d = sub.to_dict()
        assert d["definition_id"] == "pd-1"


class TestSelfIssuedIDToken:
    def test_sign_and_verify(self):
        token = SelfIssuedIDToken(
            sub="did:key:holder",
            sub_jwk={"kty": "oct", "k": "dev-siop-key"},
            iss="did:key:holder",
            nonce="abc123",
        )
        signed = token.sign()
        assert len(signed.split(".")) == 3

        verified = SelfIssuedIDToken.verify(signed)
        assert verified is not None
        assert verified.sub == "did:key:holder"
        assert verified.nonce == "abc123"

    def test_verify_tampered_token(self):
        token = SelfIssuedIDToken(
            sub="did:key:holder",
            sub_jwk={},
            iss="did:key:holder",
            nonce="abc",
        ).sign()
        parts = token.split(".")
        tampered = parts[0] + "." + parts[1] + "." + "AAAAAAAAAA"
        assert SelfIssuedIDToken.verify(tampered) is None

    def test_verify_bad_format(self):
        assert SelfIssuedIDToken.verify("not-a-jwt") is None

    def test_expires_in_future(self):
        import time
        token = SelfIssuedIDToken(
            sub="did:key:h",
            sub_jwk={},
            iss="did:key:h",
            nonce="x",
            exp=int(time.time()) + 3600,
        )
        assert token.exp > token.iat


class TestOIDC4VPRequest:
    def test_to_dict_minimal(self):
        req = OIDC4VPRequest(client_id="did:key:verifier")
        d = req.to_dict()
        assert d["client_id"] == "did:key:verifier"
        assert d["scope"] == "openid"
        assert "nonce" in d
        assert "state" in d

    def test_to_dict_with_presentation_definition(self, presentation_definition):
        req = OIDC4VPRequest(
            client_id="did:key:verifier",
            presentation_definition=presentation_definition,
        )
        d = req.to_dict()
        assert "presentation_definition" in d
        assert d["presentation_definition"]["id"] == "pd-1"


class TestSIOPWallet:
    def test_wallet_has_did(self, wallet):
        assert wallet.did.startswith("did:")

    def test_store_credential(self, wallet):
        assert len(wallet._credentials) == 1

    def test_respond_without_definition(self, wallet):
        req = OIDC4VPRequest(client_id="did:key:verifier")
        response = wallet.respond_to_request(req)
        assert response.id_token is not None
        assert response.vp_token is None

    def test_respond_with_presentation_definition(self, wallet, presentation_definition):
        req = OIDC4VPRequest(
            client_id="did:key:verifier",
            presentation_definition=presentation_definition,
        )
        response = wallet.respond_to_request(req)
        assert response.id_token is not None
        assert response.vp_token is not None
        assert len(response.vp_token) == 1
        assert response.presentation_submission is not None
        assert response.presentation_submission.definition_id == "pd-1"

    def test_respond_to_dict(self, wallet, presentation_definition):
        req = OIDC4VPRequest(
            client_id="did:key:verifier",
            presentation_definition=presentation_definition,
            state="mystate",
        )
        response = wallet.respond_to_request(req)
        d = response.to_dict()
        assert "id_token" in d
        assert "vp_token" in d
        assert "presentation_submission" in d
        assert d["state"] == "mystate"


class TestOIDC4VPVerifier:
    def test_verify_valid_response(self, wallet, verifier, presentation_definition):
        req = OIDC4VPRequest(
            client_id="did:key:verifier",
            presentation_definition=presentation_definition,
        )
        response = wallet.respond_to_request(req)
        result = verifier.verify_response(response, expected_nonce=req.nonce)
        assert result["valid"] is True
        assert len(result["errors"]) == 0

    def test_verify_tampered_id_token(self, wallet, verifier, presentation_definition):
        req = OIDC4VPRequest(
            client_id="did:key:verifier",
            presentation_definition=presentation_definition,
        )
        response = wallet.respond_to_request(req)
        response.id_token = response.id_token + "tampered"
        result = verifier.verify_response(response, expected_nonce=req.nonce)
        assert result["valid"] is False
        assert any("id_token" in e for e in result["errors"]) or not result["valid"]

    def test_verify_wrong_nonce(self, wallet, verifier, presentation_definition):
        req = OIDC4VPRequest(
            client_id="did:key:verifier",
            presentation_definition=presentation_definition,
        )
        response = wallet.respond_to_request(req)
        result = verifier.verify_response(response, expected_nonce="wrong-nonce")
        assert result["valid"] is False
        assert any("nonce" in e for e in result["errors"])

    def test_verify_state_mismatch(self, wallet, verifier):
        req = OIDC4VPRequest(client_id="did:key:verifier")
        response = wallet.respond_to_request(req)
        result = verifier.verify_response(response, expected_nonce=req.nonce, expected_state="wrong")
        assert result["valid"] is False

    def test_verify_full_flow_with_credentials(self, wallet, verifier, presentation_definition):
        req = OIDC4VPRequest(
            client_id="did:key:verifier",
            presentation_definition=presentation_definition,
        )
        response = wallet.respond_to_request(req)
        result = verifier.verify_response(response, expected_nonce=req.nonce, expected_state=req.state)
        assert result["valid"] is True
        assert result.get("vp_count") == 1
        assert result["did"] == wallet.did

    def test_verify_with_no_credentials(self, verifier):
        empty_wallet = SIOPWallet()
        desc = InputDescriptor(id="id", name="Test", purpose="Test")
        pd = PresentationDefinition(id="pd-empty", input_descriptors=[desc])
        req = OIDC4VPRequest(
            client_id="did:key:verifier",
            presentation_definition=pd,
        )
        response = empty_wallet.respond_to_request(req)
        result = verifier.verify_response(response, expected_nonce=req.nonce)
        assert result["valid"] is True
        assert result.get("vp_count", 0) == 0


class TestCrossWallet:
    def test_different_wallet_keys_fail(self):
        wallet_a = SIOPWallet(wallet_key="key-a")
        wallet_a.store_credential({
            "type": ["VerifiableCredential"],
            "credentialSubject": {"name": "Alice"},
        })
        verifier_b = OIDC4VPVerifier(wallet_key="key-b")
        desc = InputDescriptor(id="id", name="N", purpose="P")
        pd = PresentationDefinition(id="pd-x", input_descriptors=[desc])
        req = OIDC4VPRequest(client_id="did:key:v", presentation_definition=pd)
        response = wallet_a.respond_to_request(req)
        result = verifier_b.verify_response(response, expected_nonce=req.nonce)
        assert result["valid"] is False
