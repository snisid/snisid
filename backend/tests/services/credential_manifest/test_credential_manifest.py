import pytest

from services.credential_manifest import (
    CredentialManifest,
    CredentialApplication,
    CredentialResponse,
    OutputDescriptor,
    ManifestManager,
)
from services.vc.issuer import VCIssuer


@pytest.fixture
def issuer():
    return VCIssuer(issuer_id="did:key:manifest-issuer")


@pytest.fixture
def manager(issuer):
    return ManifestManager(issuer)


class TestOutputDescriptor:
    def test_to_dict(self):
        od = OutputDescriptor(id="vc1", schema="https://example.com/schema", name="Test VC")
        d = od.to_dict()
        assert d["id"] == "vc1"
        assert d["name"] == "Test VC"


class TestCredentialManifest:
    def test_to_dict(self):
        od = OutputDescriptor(id="id1", schema="https://w3id.org/credentials/v2")
        manifest = CredentialManifest(
            id="manifest-1",
            issuer="did:key:issuer",
            name="National ID",
            output_descriptors=[od],
        )
        d = manifest.to_dict()
        assert d["id"] == "manifest-1"
        assert d["issuer"] == "did:key:issuer"
        assert len(d["output_descriptors"]) == 1

    def test_to_dict_with_presentation_definition(self):
        od = OutputDescriptor(id="id1", schema="https://w3id.org/credentials/v2")
        pd = {"id": "pd-1", "input_descriptors": []}
        manifest = CredentialManifest(
            id="manifest-2",
            issuer="did:key:issuer",
            output_descriptors=[od],
            presentation_definition=pd,
        )
        d = manifest.to_dict()
        assert d["presentation_definition"]["id"] == "pd-1"


class TestCredentialApplication:
    def test_to_dict(self):
        app = CredentialApplication(
            id="app-1",
            manifest_id="manifest-1",
            applicant="did:key:alice",
            claims={"national_id": "SN-001"},
        )
        d = app.to_dict()
        assert d["id"] == "app-1"
        assert d["claims"]["national_id"] == "SN-001"


class TestCredentialResponse:
    def test_to_dict_with_credentials(self):
        resp = CredentialResponse(
            id="resp-1",
            manifest_id="manifest-1",
            applicant="did:key:bob",
            credentials=[{"type": ["VerifiableCredential"]}],
        )
        d = resp.to_dict()
        assert len(d["credentials"]) == 1

    def test_to_dict_with_error(self):
        resp = CredentialResponse(
            id="resp-2",
            manifest_id="manifest-1",
            applicant="did:key:bob",
            error="Manifest not found",
        )
        d = resp.to_dict()
        assert d["error"] == "Manifest not found"
        assert "credentials" not in d


class TestManifestManager:
    def test_create_manifest(self, manager):
        manifest = manager.create_manifest(
            issuer_did="did:key:issuer",
            name="National ID",
            description="Get your national ID",
        )
        assert manifest.id is not None
        assert manager.get_manifest(manifest.id) is manifest

    def test_list_manifests(self, manager):
        manager.create_manifest(issuer_did="did:key:a", name="A")
        manager.create_manifest(issuer_did="did:key:b", name="B")
        assert len(manager.list_manifests()) == 2

    def test_list_manifests_by_issuer(self, manager):
        manager.create_manifest(issuer_did="did:key:x", name="X")
        manager.create_manifest(issuer_did="did:key:y", name="Y")
        assert len(manager.list_manifests(issuer_did="did:key:x")) == 1

    def test_submit_application_issues_vc(self, manager):
        manifest = manager.create_manifest(issuer_did="did:key:issuer", name="ID")
        app = CredentialApplication(
            id=str(uuid4()),
            manifest_id=manifest.id,
            applicant="did:key:applicant",
            claims={
                "national_id": "SN-MANIFEST-001",
                "first_name": "Manifest",
                "last_name": "Test",
                "date_of_birth": "1990-01-01",
                "gender": "female",
                "nationality": "HTI",
            },
        )
        response = manager.submit_application(app)
        assert response.error is None
        assert response.credentials is not None
        assert len(response.credentials) == 1
        assert response.credentials[0]["credentialSubject"]["national_id"] == "SN-MANIFEST-001"

    def test_submit_application_unknown_manifest(self, manager):
        app = CredentialApplication(
            id=str(uuid4()),
            manifest_id="nonexistent",
            applicant="did:key:alice",
        )
        response = manager.submit_application(app)
        assert response.error == "Manifest not found"

    def test_get_application(self, manager):
        manifest = manager.create_manifest(issuer_did="did:key:issuer", name="Test")
        app = CredentialApplication(
            id="app-specific",
            manifest_id=manifest.id,
            applicant="did:key:bob",
        )
        manager.submit_application(app)
        apps = manager.get_applications(manifest_id=manifest.id)
        assert len(apps) == 1
        assert apps[0].id == "app-specific"

    def test_register_manifest(self, manager):
        od = OutputDescriptor(id="custom", schema="https://example.com/schema")
        manifest = CredentialManifest(
            id="custom-manifest",
            issuer="did:key:custom",
            output_descriptors=[od],
        )
        manager.register(manifest)
        assert manager.get_manifest("custom-manifest") is manifest

    def test_get_response(self, manager):
        manifest = manager.create_manifest(issuer_did="did:key:issuer", name="Test")
        app = CredentialApplication(
            id=str(uuid4()),
            manifest_id=manifest.id,
            applicant="did:key:bob",
            claims={"national_id": "SN-RESP", "first_name": "Resp", "last_name": "Test",
                    "date_of_birth": "2000-01-01", "gender": "male", "nationality": "HTI"},
        )
        response = manager.submit_application(app)
        retrieved = manager.get_response(response.id)
        assert retrieved is not None
        assert retrieved.id == response.id


def uuid4():
    import uuid
    return str(uuid.uuid4())
