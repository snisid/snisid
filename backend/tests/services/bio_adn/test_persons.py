from __future__ import annotations

from datetime import datetime, timezone
from uuid import uuid4

import pytest
from httpx import AsyncClient, ASGITransport

from services.bio_adn.api import router as bio_adn_router
from services.bio_adn.models import (
    ArticleCategory,
    SecurityType,
    InterpolNoticeType,
    ProtectionOrderType,
    SupervisionType,
    SexOffenderRiskLevel,
    MissingCategory,
    WarrantType,
    DangerLevel,
)
from fastapi import FastAPI

app = FastAPI()
app.include_router(bio_adn_router)


@pytest.fixture
def client():
    return AsyncClient(transport=ASGITransport(app=app), base_url="http://test")


@pytest.mark.asyncio
async def test_create_wanted_person_mandatory_warrant_number(client):
    payload = {
        "charges": ["vol à main armée"],
        "issuing_date": "2026-06-01",
        "warrant_type": WarrantType.ARREST.value,
        "mco_contact": "mco@pnh.ht",
    }
    resp = await client.post("/v1/bio-adn/persons/wanted", json=payload)
    assert resp.status_code == 422


@pytest.mark.asyncio
async def test_create_wanted_person_success(client):
    payload = {
        "last_name": "Dupont",
        "first_name": "Jean",
        "charges": ["vol"],
        "issuing_date": "2026-06-01",
        "warrant_type": WarrantType.ARREST.value,
        "warrant_number": "MAND-2026-1234",
        "mco_contact": "mco@pnh.ht",
        "entering_officer": "INSP Sorel",
    }
    resp = await client.post("/v1/bio-adn/persons/wanted", json=payload)
    assert resp.status_code == 201
    data = resp.json()
    assert data["status"] == "ACTIVE"
    assert data["mco_contact"] == "mco@pnh.ht"


@pytest.mark.asyncio
async def test_create_foreign_fugitive_success(client):
    payload = {
        "interpol_notice_number": "RED-2026-00123",
        "notice_type": InterpolNoticeType.RED.value,
        "last_name": "Garcia",
        "nationality": "COL",
        "charges": ["trafic de stupéfiants"],
        "issuing_country": "Colombie",
        "entering_agency": "DCPJ",
    }
    resp = await client.post("/v1/bio-adn/persons/foreign-fugitives", json=payload)
    assert resp.status_code == 201
    data = resp.json()
    assert data["interpol_notice_number"] == "RED-2026-00123"
    assert data["status"] == "ACTIVE"


@pytest.mark.asyncio
async def test_query_foreign_fugitives(client):
    resp = await client.get("/v1/bio-adn/persons/foreign-fugitives/query")
    assert resp.status_code == 200


@pytest.mark.asyncio
async def test_get_foreign_fugitive(client):
    resp = await client.get("/v1/bio-adn/persons/foreign-fugitives/some-id")
    assert resp.status_code == 200


@pytest.mark.asyncio
async def test_create_unidentified_person_success(client):
    payload = {
        "discovery_date": "2026-06-01",
        "discovery_location": "Rue du Centre, Port-au-Prince",
        "discovery_department": "OUEST",
        "estimated_age_min": 25,
        "estimated_age_max": 35,
        "gender": "M",
        "entering_agency": "PNH",
    }
    resp = await client.post("/v1/bio-adn/persons/unidentified", json=payload)
    assert resp.status_code == 201
    data = resp.json()
    assert data["status"] == "ACTIVE"


@pytest.mark.asyncio
async def test_query_unidentified_persons(client):
    resp = await client.get("/v1/bio-adn/persons/unidentified/query")
    assert resp.status_code == 200


@pytest.mark.asyncio
async def test_get_unidentified_person(client):
    resp = await client.get("/v1/bio-adn/persons/unidentified/some-id")
    assert resp.status_code == 200


@pytest.mark.asyncio
async def test_create_terrorism_watch_success(client):
    payload = {
        "last_name": "Mohamed",
        "threat_type": "RADICALISATION",
        "entering_agency": "DCPJ",
        "approved_by_director": "DIR DCPJ",
        "approved_by_pg": "PG Cayard",
    }
    resp = await client.post("/v1/bio-adn/persons/terrorism", json=payload)
    assert resp.status_code == 201
    data = resp.json()
    assert data["status"] == "ACTIVE"


@pytest.mark.asyncio
async def test_query_terrorism_watches(client):
    resp = await client.get("/v1/bio-adn/persons/terrorism/query")
    assert resp.status_code == 200


@pytest.mark.asyncio
async def test_get_terrorism_watch(client):
    resp = await client.get("/v1/bio-adn/persons/terrorism/some-id")
    assert resp.status_code == 200


@pytest.mark.asyncio
async def test_create_protection_order_success(client):
    payload = {
        "order_type": ProtectionOrderType.DOMESTIC_VIOLENCE.value,
        "issuing_court": "TPP Port-au-Prince",
        "issuing_judge": "Juge Pierre",
        "beneficiary_name": "Marie Duval",
        "protected_person": "Marie Duval",
        "restrained_person": "Pierre Duval",
        "restrictions": ["ne pas approcher à moins de 500m", "restitution arme"],
        "issue_date": "2026-06-01",
    }
    resp = await client.post("/v1/bio-adn/persons/protection-orders", json=payload)
    assert resp.status_code == 201
    data = resp.json()
    assert data["status"] == "ACTIVE"


@pytest.mark.asyncio
async def test_query_protection_orders(client):
    resp = await client.get("/v1/bio-adn/persons/protection-orders/query")
    assert resp.status_code == 200


@pytest.mark.asyncio
async def test_get_protection_order_urgent(client):
    resp = await client.get("/v1/bio-adn/persons/protection-orders/urgent/NIU-123")
    assert resp.status_code == 200


@pytest.mark.asyncio
async def test_create_supervised_release_success(client):
    payload = {
        "niu": "NIU-123456",
        "last_name": "Pierre",
        "first_name": "Jacques",
        "supervision_type": SupervisionType.PAROLE.value,
        "start_date": "2026-06-01",
        "conditions": ["pointage hebdomadaire", "interdiction de quitter le pays"],
        "supervising_officer": "OFF Michel",
        "supervising_agency": "SPA",
    }
    resp = await client.post("/v1/bio-adn/persons/supervised-release", json=payload)
    assert resp.status_code == 201
    data = resp.json()
    assert data["status"] == "ACTIVE"


@pytest.mark.asyncio
async def test_query_supervised_releases(client):
    resp = await client.get("/v1/bio-adn/persons/supervised-release/query")
    assert resp.status_code == 200


@pytest.mark.asyncio
async def test_get_supervised_release(client):
    resp = await client.get("/v1/bio-adn/persons/supervised-release/some-id")
    assert resp.status_code == 200


@pytest.mark.asyncio
async def test_update_sex_offender_risk(client):
    payload = {"risk_level": SexOffenderRiskLevel.HIGH.value}
    resp = await client.patch("/v1/bio-adn/persons/sex-offenders/some-id/risk", json=payload)
    assert resp.status_code == 200


@pytest.mark.asyncio
async def test_register_sex_offender_invalid_risk(client):
    resp = await client.post(
        "/v1/bio-adn/persons/sex-offenders",
        params={"niu": "TEST", "conviction_date": "2026-01-01", "conviction_court": "TPP", "offenses": ["viol"], "risk_level": "INVALID"},
    )
    assert resp.status_code == 422


@pytest.mark.asyncio
async def test_review_gang_member(client):
    payload = {"review_notes": "révision annuelle OK", "auto_removal_years": 5}
    resp = await client.post("/v1/bio-adn/persons/gang-members/some-id/review", json=payload)
    assert resp.status_code == 200


@pytest.mark.asyncio
async def test_sex_offenders_query(client):
    resp = await client.get("/v1/bio-adn/persons/sex-offenders/query")
    assert resp.status_code == 200


@pytest.mark.asyncio
async def test_gang_members_query(client):
    resp = await client.get("/v1/bio-adn/persons/gang-members/query")
    assert resp.status_code == 200


@pytest.mark.asyncio
async def test_created_foreign_fugitive_invalid_missing_notice(client):
    payload = {
        "last_name": "Test",
        "charges": ["test"],
        "issuing_country": "HT",
        "entering_agency": "DCPJ",
    }
    resp = await client.post("/v1/bio-adn/persons/foreign-fugitives", json=payload)
    assert resp.status_code == 422


@pytest.mark.asyncio
async def test_terrorism_watch_requires_charges(client):
    payload = {"last_name": "", "threat_type": "RADICALISATION", "entering_agency": "DCPJ", "approved_by_director": "DIR", "approved_by_pg": "PG"}
    resp = await client.post("/v1/bio-adn/persons/terrorism", json=payload)
    assert resp.status_code == 422


@pytest.mark.asyncio
async def test_missing_person_child_triggers_bpm(client):
    payload = {
        "last_name": "Enfant",
        "first_name": "Perdu",
        "category": MissingCategory.CHILD.value,
        "missing_date": "2026-06-01",
        "missing_location": "Jacmel",
        "entering_agency": "PNH",
        "citizen_portal_submission": True,
    }
    resp = await client.post("/v1/bio-adn/persons/missing", json=payload)
    assert resp.status_code == 201
    data = resp.json()
    assert data["status"] == "ACTIVE"


# ── PER-VIO: Known Violence ─────────────────────────────────────────────────


async def test_create_violence_record(client):
    payload = {
        "niu": "NIU-VIO-001",
        "incident_type": "DOMESTIC_VIOLENCE",
        "incident_date": "2026-06-01",
        "location": "Port-au-Prince, Delmas 33",
        "arresting_agency": "PNH-DELMAS",
        "risk_level": "HIGH",
    }
    resp = await client.post("/v1/bio-adn/persons/violence", json=payload)
    assert resp.status_code == 201
    data = resp.json()
    assert data["status"] == "ACTIVE"
    assert data["record_number"].startswith("VIO-")


async def test_query_violence_records(client):
    resp = await client.get("/v1/bio-adn/persons/violence/query", params={"niu": "NIU-VIO-001"})
    assert resp.status_code == 200


async def test_get_violence_record(client):
    resp = await client.get("/v1/bio-adn/persons/violence/some-id")
    assert resp.status_code == 200


# ── PER-IDV: Identity Theft ─────────────────────────────────────────────────


async def test_create_identity_theft(client):
    payload = {
        "victim_niu": "NIU-IDV-001",
        "fraud_type": "CIN_FRAUD",
        "document_type_used": "CIN",
        "report_date": "2026-06-01",
        "reporting_agency": "PNH-DELMAS",
    }
    resp = await client.post("/v1/bio-adn/persons/identity-theft", json=payload)
    assert resp.status_code == 201
    data = resp.json()
    assert data["status"] == "ACTIVE"
    assert data["record_number"].startswith("IDV-")


async def test_query_identity_thefts(client):
    resp = await client.get("/v1/bio-adn/persons/identity-theft/query", params={"victim_niu": "NIU-IDV-001"})
    assert resp.status_code == 200


async def test_get_identity_theft(client):
    resp = await client.get("/v1/bio-adn/persons/identity-theft/some-id")
    assert resp.status_code == 200
