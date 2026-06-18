from __future__ import annotations

import pytest
from httpx import AsyncClient, ASGITransport

from services.bio_adn.models import LabEquipmentCreate, StaffTrainingCreate, UploadRequest
from fastapi import FastAPI

from services.bio_adn.api import router as bio_adn_router

app = FastAPI()
app.include_router(bio_adn_router)


@pytest.fixture
def client():
    return AsyncClient(transport=ASGITransport(app=app), base_url="http://test")


@pytest.mark.asyncio
async def test_submit_dna_duplicate_specimen(client):
    payload = {
        "specimen_number": "DUP-001",
        "index_type": "BIO-CON",
        "loci_data": {
            "CSF1PO": {"value1": 10, "value2": 12},
            "D3S1358": {"value1": 15, "value2": 17},
            "D5S818": {"value1": 11, "value2": 13},
            "D7S820": {"value1": 8, "value2": 10},
            "D8S1179": {"value1": 13, "value2": 14},
            "D13S317": {"value1": 9, "value2": 11},
            "D16S539": {"value1": 11, "value2": 12},
            "D18S51": {"value1": 14, "value2": 16},
            "D21S11": {"value1": 29, "value2": 30},
            "FGA": {"value1": 21, "value2": 23},
            "TH01": {"value1": 7, "value2": 9},
            "TPOX": {"value1": 8, "value2": 11},
            "vWA": {"value1": 16, "value2": 17},
            "D1S1656": {"value1": 13, "value2": 15},
            "D2S441": {"value1": 10, "value2": 11},
            "D2S1338": {"value1": 19, "value2": 23},
            "D10S1248": {"value1": 13, "value2": 15},
            "D12S391": {"value1": 18, "value2": 20},
            "D19S433": {"value1": 13, "value2": 14},
            "D22S1045": {"value1": 15, "value2": 16},
        },
        "quality_score": 0.95,
        "collected_date": "2026-06-01",
        "correlation_id": "corr-001",
    }
    resp1 = await client.post("/v1/bio-adn/dna/profiles", json=payload)
    assert resp1.status_code == 201

    resp2 = await client.post("/v1/bio-adn/dna/profiles", json=payload)
    assert resp2.status_code == 409
    assert "déjà soumis" in resp2.json()["detail"]["message"]


@pytest.mark.asyncio
async def test_register_equipment(client):
    payload = {
        "lab_code": "LDIS-PAP-001",
        "equipment_name": "Séquenceur capillaire ABI 3500",
        "model": "3500xL",
        "serial_number": "ABI-2026-001",
        "role": "Analyse STR",
    }
    resp = await client.post("/v1/bio-adn/lab/equipment", json=payload)
    assert resp.status_code == 201
    data = resp.json()
    assert data["status"] == "ACTIVE"
    assert data["equipment_name"] == "Séquenceur capillaire ABI 3500"


@pytest.mark.asyncio
async def test_list_equipment(client):
    resp = await client.get("/v1/bio-adn/lab/equipment")
    assert resp.status_code == 200


@pytest.mark.asyncio
async def test_get_equipment(client):
    resp = await client.get("/v1/bio-adn/lab/equipment/some-id")
    assert resp.status_code == 200


@pytest.mark.asyncio
async def test_update_calibration(client):
    payload = {"calibration_date": "2026-06-01", "calibration_due": "2027-06-01"}
    resp = await client.patch("/v1/bio-adn/lab/equipment/some-id/calibration", json=payload)
    assert resp.status_code == 200


@pytest.mark.asyncio
async def test_record_training(client):
    payload = {
        "staff_niu": "HTI-12345",
        "training_name": "STR Analysis",
        "training_code": "STR-40H",
        "duration_hours": 40,
        "completed_date": "2026-06-01",
        "issued_by": "Direction SNISID",
        "frequency": "INITIALE",
    }
    resp = await client.post("/v1/bio-adn/lab/training", json=payload)
    assert resp.status_code == 201
    data = resp.json()
    assert data["training_code"] == "STR-40H"


@pytest.mark.asyncio
async def test_list_training(client):
    resp = await client.get("/v1/bio-adn/lab/training")
    assert resp.status_code == 200


@pytest.mark.asyncio
async def test_get_training(client):
    resp = await client.get("/v1/bio-adn/lab/training/some-id")
    assert resp.status_code == 200


@pytest.mark.asyncio
async def test_ldis_upload(client):
    payload = {
        "level": "ldis-to-sdis",
        "lab_code": "LDIS-PAP-001",
        "date_from": "2026-06-01",
        "date_to": "2026-06-09",
        "operator_niu": "HTI-XXXXXXXXXX",
    }
    resp = await client.post("/v1/bio-adn/lab/upload", json=payload)
    assert resp.status_code == 202
    data = resp.json()
    assert data["success"] is True


@pytest.mark.asyncio
async def test_list_labs(client):
    resp = await client.get("/v1/bio-adn/lab/labs")
    assert resp.status_code == 200
    data = resp.json()
    assert data["total"] == 6
    assert any(lab["code"] == "LDIS-PAP-001" for lab in data["results"])
