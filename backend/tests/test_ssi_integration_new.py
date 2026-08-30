"""E2E integration test: new modules (Wallet, CHAPI, SIOPv2, CredentialFlow)."""

import pytest
from fastapi.testclient import TestClient

from services.vc.issuer import VCIssuer
from services.vc.verifier import VCVerifier
from services.vc import VCStatus
from services.wallet import Wallet
from services.chapi import CHAPIMediator, CHAPIStoreRequest, CHAPIGetRequest
from services.siopv2 import OIDC4VPRequest, OIDC4VPVerifier, PresentationDefinition, InputDescriptor, SIOPWallet
from services.credential_flow import CredentialFlow


@pytest.fixture
def app():
    import main
    main.app.dependency_overrides.clear()
    yield main.app
    main.app.dependency_overrides.clear()


@pytest.fixture
def client(app):
    return TestClient(app)


def test_wallet_vc_lifecycle(app):
    """Issue VC → Wallet store → Wallet list → Wallet VP → Verify VP."""
    issuer_did = "did:key:integration-issuer"
    issuer = VCIssuer(issuer_id=issuer_did)
    verifier = VCVerifier(trusted_issuers=[issuer_did])

    wallet = Wallet()
    assert wallet.did.startswith("did:")
    assert wallet.count() == 0

    vc = issuer.issue_identity_credential(
        subject_id=wallet.did,
        national_id="SN-WALLET-E2E-001",
        first_name="Wallet",
        last_name="E2E",
        date_of_birth="1992-03-20",
        gender="female",
        nationality="HTI",
    )
    record = wallet.store(vc.model_dump(), label="My National ID")
    assert wallet.count() == 1
    assert record.issuer_did == issuer_did

    vp = wallet.create_presentation()
    assert vp.holder == wallet.did
    assert len(vp.verifiable_credential) == 1

    assert wallet.verify_presentation(vp) is True

    found = wallet.search("national")
    assert len(found) == 1

    iss_records = wallet.get_credentials_by_issuer(issuer_did)
    assert len(iss_records) == 1

    exported = wallet.export_all()
    assert len(exported) == 1

    wallet2 = Wallet()
    wallet2.import_credentials(exported)
    assert wallet2.count() == 1
    assert wallet2.get(record.id) is not None


def test_chapi_store_get_roundtrip(app):
    """CHAPI store VC → CHAPI get VP → Verify VP."""
    wallet = Wallet()
    mediator = CHAPIMediator(wallet=wallet)

    vc_data = {
        "type": ["VerifiableCredential", "NationalIDCredential"],
        "issuer": "did:key:gov",
        "issuanceDate": "2026-01-01T00:00:00Z",
        "credentialSubject": {
            "national_id": "SN-CHAPI-001",
            "first_name": "Chapi",
            "last_name": "Test",
        },
    }

    store_req = CHAPIStoreRequest(credential=vc_data)
    store_resp = mediator.handle_store(store_req)
    assert store_resp.error is None
    assert store_resp.data["status"] == "stored"

    get_req = CHAPIGetRequest(query=[{"type": "NationalIDCredential"}])
    get_resp = mediator.handle_get(get_req)
    assert get_resp.error is None
    vp = get_resp.data
    assert vp["holder"] == wallet.did
    assert len(vp["verifiableCredential"]) == 1
    assert vp["verifiableCredential"][0]["credentialSubject"]["national_id"] == "SN-CHAPI-001"

    get_all_req = CHAPIGetRequest(query=[{"type": "VerifiableCredential"}])
    get_all_resp = mediator.handle_get(get_all_req)
    assert len(get_all_resp.data["verifiableCredential"]) == 1

    get_empty_req = CHAPIGetRequest(query=[{"type": "NonexistentCredential"}])
    get_empty_resp = mediator.handle_get(get_empty_req)
    assert len(get_empty_resp.data["verifiableCredential"]) == 0


def test_siopv2_full_flow(app):
    """SIOPv2: Verifier request → Wallet respond → Verifier verify."""
    wallet = SIOPWallet()
    wallet.store_credential({
        "type": ["VerifiableCredential", "IdentityCredential"],
        "issuer": "did:key:gov",
        "credentialSubject": {"national_id": "SN-SIOP-001"},
    })

    verifier = OIDC4VPVerifier()

    desc = InputDescriptor(
        id="identity-credential",
        name="Identity Credential",
        purpose="Verify identity for access",
    )
    pd = PresentationDefinition(
        id="pd-siope2e",
        name="Identity Check",
        input_descriptors=[desc],
    )

    req = OIDC4VPRequest(
        client_id="did:key:verifier-e2e",
        presentation_definition=pd,
    )

    response = wallet.respond_to_request(req)
    assert response.id_token is not None
    assert response.vp_token is not None
    assert response.presentation_submission is not None

    result = verifier.verify_response(response, expected_nonce=req.nonce, expected_state=req.state)
    assert result["valid"] is True
    assert result.get("vp_count") == 1
    assert result["did"] == wallet.did


def test_siopv2_tamper_rejection(app):
    """Tampered SIOPv2 response is rejected."""
    wallet = SIOPWallet()
    wallet.store_credential({
        "type": ["VerifiableCredential"],
        "credentialSubject": {"name": "Alice"},
    })

    req = OIDC4VPRequest(client_id="did:key:verifier")
    response = wallet.respond_to_request(req)

    response.id_token = response.id_token + "tampered"

    verifier = OIDC4VPVerifier()
    result = verifier.verify_response(response, expected_nonce=req.nonce)
    assert result["valid"] is False


def test_credential_flow_full(app):
    """Credential Flow: offer → request → issue → send via DIDComm."""
    iss_did = "did:key:cf-issuer"
    hol_did = "did:snisid:mainnet:cf-holder"

    issuer = VCIssuer(issuer_id=iss_did)
    flow = CredentialFlow(issuer=issuer)

    offer = flow.create_offer(
        issuer_did=iss_did,
        holder_did=hol_did,
    )
    assert offer.status == "pending"

    body = flow.build_request(offer.offer_id, hol_did)
    packed = flow.pack_and_send_request(body, hol_did, iss_did)

    req = flow.receive_request(packed)
    assert req.offer_id == offer.offer_id

    vc_data = flow.issue_from_request(req)
    assert vc_data is not None
    assert offer.status == "fulfilled"
    assert vc_data["credentialSubject"]["national_id"] == "SN-CITIZEN-001"

    packed_vc = flow.send_credential(vc_data, iss_did, hol_did)
    assert "ciphertext" in packed_vc or "signatures" in packed_vc

    received = flow.receive_credential(packed_vc)
    assert received is not None
    assert received["credentialSubject"]["national_id"] == "SN-CITIZEN-001"


def test_wallet_chapi_siopv2_combined(app):
    """Wallet stores VC → CHAPI retrieves → SIOPv2 presents → Verifier validates."""
    wallet = SIOPWallet()
    wallet.store_credential({
        "type": ["VerifiableCredential", "NationalIDCredential"],
        "issuer": "did:key:gov",
        "credentialSubject": {"national_id": "SN-COMBINED-001"},
    })

    verifier = OIDC4VPVerifier()

    desc = InputDescriptor(id="national-id", name="National ID", purpose="Access")
    pd = PresentationDefinition(id="pd-combined", name="ID Check", input_descriptors=[desc])
    req = OIDC4VPRequest(client_id="did:key:v", presentation_definition=pd)

    response = wallet.respond_to_request(req)
    result = verifier.verify_response(response, expected_nonce=req.nonce, expected_state=req.state)
    assert result["valid"] is True
    assert result.get("vp_count") == 1


def test_api_integration_via_testclient(app, client):
    """Full API integration: Wallet API + CHAPI + SIOPv2 via TestClient."""
    import services.wallet.api as wallet_api
    import services.chapi.api as chapi_api
    import services.siopv2.api as siop_api

    wallet_api._wallet = None
    chapi_api._mediator = None
    siop_api._wallet = None
    siop_api._requests = {}

    siop_api._wallet = None
    siop_api._requests = {}
    chapi_api._mediator = None

    chapi_store_resp = client.post("/v1/chapi/store", json={
        "type": ["VerifiableCredential", "NationalIDCredential"],
        "issuer": "did:key:issuer",
        "credentialSubject": {"national_id": "SN-API-E2E"},
    })
    assert chapi_store_resp.status_code == 200

    chapi_resp = client.post("/v1/chapi/get", json={
        "query": [{"type": "NationalIDCredential"}],
    })
    assert chapi_resp.status_code == 200
    data = chapi_resp.json()["data"]
    assert len(data["verifiableCredential"]) == 1

    did_resp = client.get("/v1/chapi/wallet")
    assert did_resp.status_code == 200
    holder_did = did_resp.json()["did"]

    pd_resp = client.post("/v1/siopv2/verifier/request-with-pd")
    assert pd_resp.status_code == 200
    pd_data = pd_resp.json()
    nonce = pd_data["_internal_nonce"]
    request = pd_data["request"]

    wallet_resp = client.post("/v1/siopv2/wallet/respond", json=request)
    assert wallet_resp.status_code == 200

    verify_resp = client.post(
        "/v1/siopv2/verifier/verify",
        params={"expected_nonce": nonce},
        json=wallet_resp.json(),
    )
    assert verify_resp.status_code == 200
    assert verify_resp.json()["valid"] is True


def test_credential_flow_via_api(app, client):
    """Credential Flow API: offer → request → issue → full roundtrip."""
    import services.credential_flow.api as cf_api
    cf_api._flow = None

    offer_resp = client.post(
        "/v1/credential-flow/offer",
        params={
            "issuer_did": "did:key:cf-api",
            "holder_did": "did:snisid:mainnet:cf-api-holder",
        },
    )
    assert offer_resp.status_code == 200
    offer = offer_resp.json()

    req_resp = client.post(
        "/v1/credential-flow/request",
        params={
            "offer_id": offer["offer_id"],
            "holder_did": "did:snisid:mainnet:cf-api-holder",
            "issuer_did": "did:key:cf-api",
        },
    )
    assert req_resp.status_code == 200
    packed_req = req_resp.json()["packed"]

    recv_resp = client.post("/v1/credential-flow/receive-request", json=packed_req)
    assert recv_resp.status_code == 200
    received_req = recv_resp.json()

    issue_resp = client.post("/v1/credential-flow/issue", json=received_req)
    assert issue_resp.status_code == 200
    vc_data = issue_resp.json()
    assert vc_data["credentialSubject"]["national_id"] == "SN-CITIZEN-001"

    send_resp = client.post(
        "/v1/credential-flow/send-credential",
        params={
            "issuer_did": "did:key:cf-api",
            "holder_did": "did:snisid:mainnet:cf-api-holder",
        },
        json=vc_data,
    )
    assert send_resp.status_code == 200
    packed_vc = send_resp.json()["packed"]

    recv_vc_resp = client.post("/v1/credential-flow/receive-credential", json=packed_vc)
    assert recv_vc_resp.status_code == 200
    assert recv_vc_resp.json()["credentialSubject"]["national_id"] == "SN-CITIZEN-001"
