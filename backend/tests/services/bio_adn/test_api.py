import pytest


class TestSubmitDNAProfile:

    async def test_submit_valid(self, client, valid_str_profile_submit):
        response = await client.post("/v1/bio-adn/dna/profiles", json=valid_str_profile_submit)
        assert response.status_code == 201
        data = response.json()
        assert data["accepted"] is True
        assert "sample_id" in data

    async def test_submit_missing_required_loci(self, client, valid_str_profile_submit):
        payload = {**valid_str_profile_submit}
        payload["loci_data"] = {"CSF1PO": {"value1": "10", "value2": "12"}}
        response = await client.post("/v1/bio-adn/dna/profiles", json=payload)
        assert response.status_code == 422

    async def test_submit_quality_below_threshold(self, client, valid_str_profile_submit):
        payload = {**valid_str_profile_submit, "specimen_number": "FSC-2026-LOW-Q", "quality_score": 0.55}
        response = await client.post("/v1/bio-adn/dna/profiles", json=payload)
        assert response.status_code == 422
        detail = response.json()["detail"]
        assert detail["code"] == "BIO-001"

    async def test_submit_unknown_index_type(self, client, valid_str_profile_submit):
        payload = {**valid_str_profile_submit, "index_type": "BIO-UNKNOWN"}
        response = await client.post("/v1/bio-adn/dna/profiles", json=payload)
        assert response.status_code == 422

    async def test_submit_empty_specimen(self, client, valid_str_profile_submit):
        payload = {**valid_str_profile_submit, "specimen_number": ""}
        response = await client.post("/v1/bio-adn/dna/profiles", json=payload)
        assert response.status_code == 422

    async def test_submit_bio_con_exact_threshold(self, client, valid_str_profile_submit):
        payload = {**valid_str_profile_submit, "index_type": "BIO-CON", "quality_score": 0.95}
        response = await client.post("/v1/bio-adn/dna/profiles", json=payload)
        assert response.status_code == 201

    async def test_submit_bio_rni_minimum(self, client, valid_str_profile_submit):
        payload = {**valid_str_profile_submit, "index_type": "BIO-RNI", "quality_score": 0.50}
        response = await client.post("/v1/bio-adn/dna/profiles", json=payload)
        assert response.status_code == 201


class TestDNAHit:

    async def test_get_hit(self, client):
        response = await client.get("/v1/bio-adn/dna/hits/hit-001")
        assert response.status_code == 200
        assert response.json()["hit_id"] == "hit-001"

    async def test_get_hit_empty_id(self, client):
        response = await client.get("/v1/bio-adn/dna/hits/")
        assert response.status_code in (404, 405)

    async def test_get_hit_not_found(self, client):
        response = await client.get("/v1/bio-adn/dna/hits/    ")
        assert response.status_code in (200, 404)


class TestExpungeProfile:

    async def test_expunge_valid(self, client):
        response = await client.post(
            "/v1/bio-adn/dna/profiles/sample-001/expunge?court_order_ref=COURT-2026-001&reason=acquittement&officer_niu=NIU-001"
        )
        assert response.status_code == 202
        data = response.json()
        assert data["success"] is True

    async def test_expunge_missing_court_order(self, client):
        response = await client.post("/v1/bio-adn/dna/profiles/sample-001/expunge")
        assert response.status_code == 422


class TestWantedPersons:

    async def test_create_wanted(self, client, valid_wanted_request):
        response = await client.post("/v1/bio-adn/persons/wanted", json=valid_wanted_request)
        assert response.status_code == 201
        data = response.json()
        assert "record_id" in data
        assert data["status"] == "ACTIVE"

    async def test_create_wanted_missing_charges(self, client, valid_wanted_request):
        payload = {**valid_wanted_request, "charges": []}
        response = await client.post("/v1/bio-adn/persons/wanted", json=payload)
        assert response.status_code == 422

    async def test_query_wanted(self, client):
        response = await client.get("/v1/bio-adn/persons/wanted/query?last_name=Pierre")
        assert response.status_code == 200
        assert "results" in response.json()

    async def test_get_wanted(self, client):
        response = await client.get("/v1/bio-adn/persons/wanted/WP-001")
        assert response.status_code == 200

    async def test_update_wanted_status(self, client):
        response = await client.patch("/v1/bio-adn/persons/wanted/WP-001/status?status=CLEARED")
        assert response.status_code == 200
        assert response.json()["status"] == "CLEARED"


class TestMissingPersons:

    async def test_create_missing(self, client):
        response = await client.post("/v1/bio-adn/persons/missing", json={
            "last_name": "Dupont",
            "first_name": "Marie",
            "category": "CHILD",
            "missing_date": "2026-06-01",
            "missing_location": "Pétion-Ville",
            "entering_agency": "PNH-PAP",
        })
        assert response.status_code == 201

    async def test_create_missing_invalid_category(self, client):
        response = await client.post("/v1/bio-adn/persons/missing", json={
            "last_name": "Dupont",
            "first_name": "Marie",
            "category": "INVALID",
            "missing_date": "2026-06-01",
            "missing_location": "Pétion-Ville",
            "entering_agency": "PNH-PAP",
        })
        assert response.status_code == 422

    async def test_query_missing(self, client):
        response = await client.get("/v1/bio-adn/persons/missing/query?last_name=Dupont")
        assert response.status_code == 200


class TestStolenProperty:

    async def test_report_vehicle(self, client, valid_stolen_vehicle_request):
        response = await client.post("/v1/bio-adn/property/vehicles", json=valid_stolen_vehicle_request)
        assert response.status_code == 201

    async def test_report_vehicle_invalid_vin(self, client, valid_stolen_vehicle_request):
        payload = {**valid_stolen_vehicle_request, "vin": "SHORT"}
        response = await client.post("/v1/bio-adn/property/vehicles", json=payload)
        assert response.status_code == 422

    async def test_report_firearm(self, client):
        response = await client.post("/v1/bio-adn/property/firearms", json={
            "serial_number": "SN-12345",
            "theft_date": "2026-06-01",
            "entering_agency": "PNH-PAP",
        })
        assert response.status_code == 201

    async def test_report_document(self, client):
        response = await client.post("/v1/bio-adn/property/documents", json={
            "document_type": "PASSPORT",
            "report_date": "2026-06-01",
            "theft_type": "STOLEN",
        })
        assert response.status_code == 201

    async def test_report_document_invalid_type(self, client):
        response = await client.post("/v1/bio-adn/property/documents", json={
            "document_type": "INVALID",
            "report_date": "2026-06-01",
            "theft_type": "STOLEN",
        })
        assert response.status_code == 422

    async def test_report_vessel(self, client):
        response = await client.post("/v1/bio-adn/property/vessels", json={
            "theft_location": "Port-au-Prince",
            "theft_date": "2026-06-01",
        })
        assert response.status_code == 201

    async def test_query_property(self, client):
        response = await client.get("/v1/bio-adn/property/query?plate_number=AA-1234")
        assert response.status_code == 200


class TestVehicleRecovery:

    async def test_recover_vehicle(self, client):
        response = await client.patch("/v1/bio-adn/property/vehicles/V-001/recover", json={
            "recovered_location": "Pétion-Ville",
            "recovering_agency": "PNH-PAP",
            "notes": "Véhicule retrouvé stationné",
        })
        assert response.status_code == 200
        assert response.json()["status"] == "RECOVERED"

    async def test_recover_missing_id(self, client):
        response = await client.patch("/v1/bio-adn/property/vehicles//recover", json={
            "recovered_location": "Delmas",
            "recovering_agency": "PNH",
        })
        assert response.status_code in (404, 405)


class TestStolenArticles:

    async def test_report_article(self, client):
        response = await client.post("/v1/bio-adn/property/articles", json={
            "category": "CATTLE",
            "description": "Zebu marque PNH",
            "serial_number": "EAR-12345",
            "estimated_value": 150000,
            "theft_date": "2026-06-10",
            "theft_location": "Arcahaie",
            "entering_agency": "PNH-ARC",
        })
        assert response.status_code == 201
        data = response.json()
        assert data["category"] == "CATTLE"
        assert "record_id" in data

    async def test_report_article_missing_description(self, client):
        response = await client.post("/v1/bio-adn/property/articles", json={
            "category": "JEWELRY",
            "theft_date": "2026-06-01",
            "theft_location": "PAP",
            "entering_agency": "PNH",
        })
        assert response.status_code == 422

    async def test_query_articles(self, client):
        response = await client.get("/v1/bio-adn/property/articles/query?category=CATTLE")
        assert response.status_code == 200


class TestStolenSecurities:

    async def test_report_security(self, client):
        response = await client.post("/v1/bio-adn/property/securities", json={
            "security_type": "CHEQUE",
            "issuer": "BRH",
            "security_number": "CHQ-001",
            "face_value": 500000,
            "theft_date": "2026-06-05",
            "theft_location": "PAP",
            "entering_agency": "PNH-PAP",
        })
        assert response.status_code == 201
        data = response.json()
        assert data["security_type"] == "CHEQUE"
        assert "record_id" in data

    async def test_report_security_missing_number(self, client):
        response = await client.post("/v1/bio-adn/property/securities", json={
            "security_type": "BOND",
            "issuer": "BRH",
            "theft_date": "2026-06-01",
            "theft_location": "PAP",
            "entering_agency": "PNH",
        })
        assert response.status_code == 422

    async def test_query_securities(self, client):
        response = await client.get("/v1/bio-adn/property/securities/query?security_type=CHEQUE")
        assert response.status_code == 200


class TestArmCrimeSceneHit:

    async def test_report_arm_hit(self, client):
        response = await client.post(
            "/v1/bio-adn/property/firearms/F-001/crime-scene-hit"
            "?crime_scene_ref=SC-2026-001"
            "&case_number=DCPJ-2026-001"
        )
        assert response.status_code == 201
        assert response.json()["status"] == "dispatched"


class TestONIDocumentRevoke:

    async def test_oni_document_revoke(self, client):
        response = await client.post("/v1/bio-adn/property/documents/oni-revoke", json={
            "document_type": "PASSPORT",
            "document_number": "HT-123456",
            "revocation_reason": "perte",
            "revoked_by": "ONI-AGENT-001",
            "revoked_at": "2026-06-11T10:00:00Z",
        })
        assert response.status_code == 200
        assert response.json()["success"] is True


class TestLAPI:

    async def test_plate_query(self, client):
        response = await client.get("/v1/bio-adn/lapi/plate/AA-1234")
        assert response.status_code == 200
        data = response.json()
        assert data["hit_found"] is False
        assert "query_id" in data

    async def test_vin_query(self, client):
        response = await client.get("/v1/bio-adn/lapi/vin/VIN-1234567890")
        assert response.status_code == 200
        assert response.json()["hit_found"] is False


class TestHealth:

    async def test_health(self, client):
        response = await client.get("/v1/bio-adn/health")
        assert response.status_code in (200, 404)
