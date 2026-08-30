import itertools
from collections.abc import AsyncGenerator
from unittest.mock import AsyncMock, MagicMock

import pytest
from uuid import uuid4

_specimen_counter = itertools.count(1)
from fastapi import FastAPI
from httpx import AsyncClient, ASGITransport

from services.bio_adn.api import router as bio_adn_router
from shared.events import KafkaProducer


@pytest.fixture
def app() -> FastAPI:
    application = FastAPI()
    application.include_router(bio_adn_router)
    return application


@pytest.fixture
async def client(app: FastAPI) -> AsyncGenerator[AsyncClient, None]:
    async with AsyncClient(
        transport=ASGITransport(app=app), base_url="http://test"
    ) as ac:
        yield ac


@pytest.fixture
def mock_kafka() -> MagicMock:
    mock = MagicMock(spec=KafkaProducer)
    mock.publish = AsyncMock(return_value=None)
    return mock


@pytest.fixture(autouse=True)
def patch_kafka(monkeypatch, mock_kafka):
    monkeypatch.setattr("services.bio_adn.api._producer", mock_kafka)


@pytest.fixture
def valid_str_profile_submit() -> dict:
    n = next(_specimen_counter)
    return {
        "specimen_number": f"FSC-2026-{n:04d}",
        "correlation_id": str(uuid4()),
        "index_type": "BIO-FSC",
        "loci_data": {
            "CSF1PO": {"value1": "10", "value2": "12"},
            "D3S1358": {"value1": "15", "value2": "17"},
            "D5S818": {"value1": "11", "value2": "13"},
            "D7S820": {"value1": "8", "value2": "10"},
            "D8S1179": {"value1": "13", "value2": "14"},
            "D13S317": {"value1": "9", "value2": "11"},
            "D16S539": {"value1": "11", "value2": "12"},
            "D18S51": {"value1": "14", "value2": "16"},
            "D21S11": {"value1": "29", "value2": "30"},
            "FGA": {"value1": "21", "value2": "23"},
            "TH01": {"value1": "7", "value2": "9"},
            "TPOX": {"value1": "8", "value2": "11"},
            "vWA": {"value1": "16", "value2": "17"},
            "D1S1656": {"value1": "13", "value2": "15"},
            "D2S441": {"value1": "10", "value2": "11"},
            "D2S1338": {"value1": "19", "value2": "23"},
            "D10S1248": {"value1": "13", "value2": "15"},
            "D12S391": {"value1": "18", "value2": "20"},
            "D19S433": {"value1": "13", "value2": "14"},
            "D22S1045": {"value1": "15", "value2": "16"},
        },
        "quality_score": 0.92,
        "case_number": "DCPJ-2026-001",
        "collected_date": "2026-06-01",
    }


@pytest.fixture
def valid_wanted_request() -> dict:
    return {
        "last_name": "Pierre",
        "first_name": "Jean",
        "charges": ["Vol à main armée", "Recel"],
        "warrant_type": "MAN-ARR",
        "warrant_number": "W-2026-001",
        "issuing_date": "2026-06-01",
        "danger_level": "HIGH",
        "mco_contact": "PNH-DELMAS",
        "armed_dangerous": True,
    }


@pytest.fixture
def valid_stolen_vehicle_request() -> dict:
    return {
        "plate_number": "AA-1234",
        "vehicle_make": "Toyota",
        "vehicle_model": "Hilux",
        "vehicle_year": 2020,
        "vehicle_color": "Blanc",
        "theft_date": "2026-06-10",
        "theft_location": "Delmas 33",
        "theft_department": "OUEST",
    }
