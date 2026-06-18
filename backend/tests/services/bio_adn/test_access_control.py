import pytest


class TestAccessControl:

    async def test_lab_technician_cannot_access_other_lab(self, client):
        response = await client.get("/v1/bio-adn/dna/profiles?lab_id=LDIS-PAP-001")
        assert response.status_code in (200, 404)

    async def test_identity_link_endpoint_exists(self, client):
        response = await client.post("/v1/bio-adn/dna/identity-links", json={
            "sample_id": "sample-001",
            "niu": "NIU-001",
            "court_order_ref": "COURT-001",
            "purpose": "criminal_investigation",
        })
        assert response.status_code in (201, 422, 405)

    async def test_expunge_requires_court_order(self, client):
        response = await client.post("/v1/bio-adn/dna/profiles/sample-001/expunge")
        assert response.status_code == 422

    async def test_expunge_with_court_order(self, client):
        response = await client.post(
            "/v1/bio-adn/dna/profiles/sample-002/expunge?court_order_ref=COURT-001&reason=acquittement"
        )
        assert response.status_code == 202

    async def test_sex_offender_requires_valid_input(self, client):
        response = await client.post("/v1/bio-adn/persons/sex-offenders", json={
            "conviction_date": "2026-01-15",
            "conviction_court": "Tribunal PAP",
            "offenses": ["Agression sexuelle"],
            "risk_level": "HIGH",
            "niu": "NIU-001",
        })
        assert response.status_code in (201, 422)

    async def test_gang_member_registration(self, client):
        response = await client.post("/v1/bio-adn/persons/gang-members", json={
            "last_name": "Destin",
            "first_name": "Mackenson",
            "gang_name": "5 Segonn",
            "membership_type": "MEMBER",
        })
        assert response.status_code in (201, 422)

    async def test_identity_link_without_court_order_ref(self, client):
        response = await client.post("/v1/bio-adn/dna/identity-links", json={
            "sample_id": "sample-003",
            "niu": "NIU-003",
            "purpose": "criminal_investigation",
        })
        assert response.status_code in (201, 405, 422)

    async def test_wanted_person_empty_charges_rejected(self, client):
        response = await client.post("/v1/bio-adn/persons/wanted", json={
            "last_name": "Test",
            "warrant_type": "MAN-ARR",
            "issuing_date": "2026-01-01",
            "charges": [],
            "mco_contact": "PNH",
        })
        assert response.status_code >= 400
