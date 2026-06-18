import pytest
import asyncio
import json
import uuid
from datetime import datetime, timezone


@pytest.fixture
def event_timestamp():
    return int(datetime.now(timezone.utc).timestamp() * 1000)


@pytest.fixture
def sample_wanted_event(event_timestamp):
    return {
        "event_id": str(uuid.uuid4()),
        "event_type": "CREATED",
        "niu": "HT12345678",
        "full_name": "JEAN DUPONT",
        "risk_level": "HIGH",
        "status": "ACTIVE",
        "warrants_active": 3,
        "aliases": ["JEANO", "TONTON"],
        "gang_affiliations": ["Gang 400 Mawozo"],
        "timestamp": event_timestamp,
        "source": "test",
        "agency": "PNH",
        "correlation_id": str(uuid.uuid4()),
    }


@pytest.fixture
def sample_alpr_event(event_timestamp):
    return {
        "reading_id": str(uuid.uuid4()),
        "plate": "AB-123-CD",
        "make": "TOYOTA",
        "model": "HILUX",
        "year": 2023,
        "color": "BLACK",
        "location": "18.5333,-72.3333",
        "speed_kmh": 85,
        "timestamp": event_timestamp,
        "camera_id": "CAM-PAP-001",
        "owner_niu": "HT12345678",
    }


@pytest.fixture
def sample_transaction_event(event_timestamp):
    return {
        "transaction_id": str(uuid.uuid4()),
        "sender_niu": "HT98765432",
        "beneficiary_niu": "HT12345678",
        "amount": 50000.00,
        "currency": "USD",
        "bank": "Banque Nationale de Credit",
        "transaction_type": "WIRE_TRANSFER",
        "country": "HT",
        "timestamp": event_timestamp,
    }


@pytest.fixture
def sample_border_event(event_timestamp):
    return {
        "crossing_id": str(uuid.uuid4()),
        "niu": "HT12345678",
        "passport_number": "HT123456",
        "port_of_entry": "Aéroport Toussaint Louverture",
        "direction": "EXIT",
        "origin_country": "HT",
        "destination_country": "US",
        "timestamp": event_timestamp,
    }


@pytest.fixture
def sample_cyber_event(event_timestamp):
    return {
        "ioc_id": str(uuid.uuid4()),
        "ioc_type": "IP",
        "ioc_value": "185.220.101.42",
        "threat_type": "C2",
        "confidence": 0.95,
        "tags": ["botnet", "ransomware"],
        "timestamp": event_timestamp,
    }


@pytest.fixture
def sample_missing_event(event_timestamp):
    return {
        "niu": "HT87654321",
        "full_name": "MARIE JEAN",
        "missing_since": event_timestamp - 86400000,
        "last_seen_location": "Pétion-Ville",
        "age_at_disappearance": 14,
        "alert_type": "AMBER",
        "status": "ACTIVE",
        "reported_by": "Famille Jean",
        "agency_code": "PNH",
    }
