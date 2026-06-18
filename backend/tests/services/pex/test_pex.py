import pytest

from services.pex import (
    InputDescriptor,
    PresentationDefinition,
    match_descriptor,
    evaluate_definition,
    filter_credentials_by_definition,
    _evaluate_constraint,
    _resolve_json_path,
)


@pytest.fixture
def alice_id_credential():
    return {
        "issuer": "did:key:issuer1",
        "issuanceDate": "2026-01-01T00:00:00Z",
        "type": ["VerifiableCredential", "IdentityCredential"],
        "credentialSubject": {
            "id": "did:key:alice",
            "givenName": "Alice",
            "familyName": "Smith",
            "birthDate": "1990-01-01",
            "nationality": "US",
        },
        "credentialStatus": {"id": "https://status.example.com/1", "status": "Active"},
    }


@pytest.fixture
def alice_driver_license():
    return {
        "issuer": "did:key:dmv",
        "issuanceDate": "2026-03-01T00:00:00Z",
        "type": ["VerifiableCredential", "DriversLicenseCredential"],
        "credentialSubject": {
            "id": "did:key:alice",
            "licenseNumber": "DL-12345",
            "licenseClass": "C",
            "state": "California",
        },
    }


class TestResolveJsonPath:
    def test_simple_path(self):
        doc = {"a": {"b": {"c": 42}}}
        assert _resolve_json_path(doc, "$.a.b.c") == 42

    def test_credential_subject_path(self):
        doc = {"credentialSubject": {"id": "did:key:alice", "name": "Alice"}}
        assert _resolve_json_path(doc, "$.credentialSubject.id") == "did:key:alice"
        assert _resolve_json_path(doc, "$.credentialSubject.name") == "Alice"

    def test_nonexistent_path(self):
        doc = {"a": 1}
        assert _resolve_json_path(doc, "$.b.c") is None

    def test_array_index(self):
        doc = {"items": [10, 20, 30]}
        assert _resolve_json_path(doc, "$.items.1") == 20
        assert _resolve_json_path(doc, "$.items.5") is None


class TestEvaluateConstraint:
    def test_fields_match(self, alice_id_credential):
        constraint = {
            "fields": [
                {
                    "path": ["$.credentialSubject.givenName"],
                    "filter": {"const": "Alice"},
                }
            ]
        }
        ok, errors = _evaluate_constraint(alice_id_credential, constraint)
        assert ok
        assert errors == []

    def test_fields_mismatch(self, alice_id_credential):
        constraint = {
            "fields": [
                {
                    "path": ["$.credentialSubject.givenName"],
                    "filter": {"const": "Bob"},
                }
            ]
        }
        ok, errors = _evaluate_constraint(alice_id_credential, constraint)
        assert not ok
        assert len(errors) > 0

    def test_fields_optional_missing(self, alice_id_credential):
        constraint = {
            "fields": [
                {"path": ["$.credentialSubject.nonexistent"], "optional": True},
                {"path": ["$.credentialSubject.givenName"], "filter": {"const": "Alice"}},
            ]
        }
        ok, errors = _evaluate_constraint(alice_id_credential, constraint)
        assert ok

    def test_subject_is_issuer_mismatch(self, alice_id_credential):
        constraint = {"subject_is_issuer": True}
        ok, _ = _evaluate_constraint(alice_id_credential, constraint)
        assert not ok

    def test_is_holder_matches(self, alice_id_credential):
        constraint = {"is_holder": ["did:key:alice"]}
        ok, _ = _evaluate_constraint(alice_id_credential, constraint)
        assert ok

    def test_is_holder_mismatch(self, alice_id_credential):
        constraint = {"is_holder": ["did:key:bob"]}
        ok, _ = _evaluate_constraint(alice_id_credential, constraint)
        assert not ok


class TestInputDescriptor:
    def test_from_dict_with_schema(self):
        data = {
            "id": "id_1",
            "name": "Identity Credential",
            "purpose": "Verify identity",
            "schema": [{"uri": "https://www.w3.org/ns/IdentityCredential"}],
            "constraints": {"fields": []},
        }
        desc = InputDescriptor(id=data["id"], name=data["name"], purpose=data["purpose"],
                                schema_uris=["https://www.w3.org/ns/IdentityCredential"],
                                constraints=data["constraints"])
        assert desc.id == "id_1"
        assert desc.name == "Identity Credential"


class TestMatchDescriptor:
    def test_match_by_schema(self, alice_id_credential):
        desc = InputDescriptor(id="id_1", schema_uris=["IdentityCredential"])
        result = match_descriptor(desc, alice_id_credential)
        assert result.matched

    def test_no_match_wrong_schema(self, alice_id_credential):
        desc = InputDescriptor(id="id_1", schema_uris=["DriversLicense"])
        result = match_descriptor(desc, alice_id_credential)
        assert not result.matched

    def test_match_with_constraints(self, alice_id_credential):
        desc = InputDescriptor(
            id="id_1",
            constraints={"fields": [{"path": ["$.credentialSubject.givenName"], "filter": {"const": "Alice"}}]},
        )
        result = match_descriptor(desc, alice_id_credential)
        assert result.matched

    def test_no_match_wrong_constraint(self, alice_id_credential):
        desc = InputDescriptor(
            id="id_1",
            constraints={"fields": [{"path": ["$.credentialSubject.givenName"], "filter": {"const": "Bob"}}]},
        )
        result = match_descriptor(desc, alice_id_credential)
        assert not result.matched


class TestEvaluateDefinition:
    def test_all_descriptors_match(self, alice_id_credential, alice_driver_license):
        pd = PresentationDefinition(
            id="pd-1",
            input_descriptors=[
                InputDescriptor(id="id_1", schema_uris=["IdentityCredential"]),
                InputDescriptor(id="id_2", schema_uris=["DriversLicenseCredential"]),
            ],
        )
        result = evaluate_definition(pd, [alice_id_credential, alice_driver_license])
        assert result.valid
        assert len(result.matches) == 2
        assert all(m.matched for m in result.matches)

    def test_missing_descriptor(self, alice_id_credential):
        pd = PresentationDefinition(
            id="pd-1",
            input_descriptors=[
                InputDescriptor(id="id_1", schema_uris=["IdentityCredential"]),
                InputDescriptor(id="id_2", schema_uris=["PassportCredential"]),
            ],
        )
        result = evaluate_definition(pd, [alice_id_credential])
        assert not result.valid

    def test_empty_credentials(self):
        pd = PresentationDefinition(
            id="pd-1",
            input_descriptors=[InputDescriptor(id="id_1", schema_uris=["IdentityCredential"])],
        )
        result = evaluate_definition(pd, [])
        assert not result.valid

    def test_submission_requirements_all(self, alice_id_credential, alice_driver_license):
        pd = PresentationDefinition(
            id="pd-1",
            input_descriptors=[
                InputDescriptor(id="id_1", schema_uris=["IdentityCredential"], group=["A"]),
                InputDescriptor(id="id_2", schema_uris=["DriversLicenseCredential"], group=["A"]),
            ],
            submission_requirements=[{"name": "All IDs", "rule": "all", "from": "A"}],
        )
        result = evaluate_definition(pd, [alice_id_credential, alice_driver_license])
        assert result.valid


class TestFilterCredentials:
    def test_filters_matching(self, alice_id_credential, alice_driver_license):
        pd = PresentationDefinition(
            id="pd-1",
            input_descriptors=[InputDescriptor(id="id_1", schema_uris=["IdentityCredential"])],
        )
        filtered = filter_credentials_by_definition(pd, [alice_id_credential, alice_driver_license])
        assert len(filtered) == 1
        assert "IdentityCredential" in filtered[0]["type"]
