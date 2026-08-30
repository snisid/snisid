from __future__ import annotations

import pytest
from httpx import AsyncClient, ASGITransport

from fastapi import FastAPI
from services.bio_adn.api import router as bio_adn_router

app = FastAPI()
app.include_router(bio_adn_router)


@pytest.fixture
def client():
    return AsyncClient(transport=ASGITransport(app=app), base_url="http://test")


@pytest.mark.asyncio
async def test_ndis_stats(client):
    resp = await client.get("/v1/bio-adn/ndis/stats")
    assert resp.status_code == 200
    data = resp.json()
    assert "total_bio_con" in data
    assert data["total_bio_con"] == 5000


@pytest.mark.asyncio
async def test_ndis_hits(client):
    resp = await client.get("/v1/bio-adn/ndis/hits")
    assert resp.status_code == 200


@pytest.mark.asyncio
async def test_ndis_hits_filtered(client):
    resp = await client.get("/v1/bio-adn/ndis/hits?sdis=SDIS-OUEST&match_type=FULL_MATCH")
    assert resp.status_code == 200


@pytest.mark.asyncio
async def test_ndis_reports(client):
    resp = await client.get("/v1/bio-adn/ndis/reports")
    assert resp.status_code == 200


@pytest.mark.asyncio
async def test_ndis_generate_report(client):
    resp = await client.post("/v1/bio-adn/ndis/reports/generate?report_type=STATS")
    assert resp.status_code == 202
    data = resp.json()
    assert data["status"] == "generating"


@pytest.mark.asyncio
async def test_ndis_generate_report_invalid_type(client):
    resp = await client.post("/v1/bio-adn/ndis/reports/generate?report_type=INVALID")
    assert resp.status_code == 202


@pytest.mark.asyncio
async def test_interpol_submit(client):
    payload = {
        "sample_ids": ["SAMPLE-001", "SAMPLE-002"],
        "reason": "disaster_victim",
        "case_number": "DIS-2026-001",
    }
    resp = await client.post("/v1/bio-adn/ndis/interpol/submit", json=payload)
    assert resp.status_code == 202
    data = resp.json()
    assert data["submitted_samples"] == 2
    assert data["status"] == "submitted"


@pytest.mark.asyncio
async def test_interpol_submit_invalid_reason(client):
    payload = {
        "sample_ids": ["SAMPLE-001"],
        "reason": "invalid_reason",
    }
    resp = await client.post("/v1/bio-adn/ndis/interpol/submit", json=payload)
    assert resp.status_code == 422
