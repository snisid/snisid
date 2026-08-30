# SNI-SIDE Python SDK

Client Python officiel pour l'API SNI-SIDE.

## Installation

```bash
pip install sniside-sdk
```

## Usage

```python
from sniside import SNISIDClient

client = SNISIDClient(api_key="sk-...", base_url="https://api.sniside.ht")

# === RECHERCHE UNIFIÉE ===
results = client.search("HT12345678")
# => {"detected_type": "NIU", "results": [...], "execution_time_ms": 45}

# === NCID ===
person = client.ncid.get_wanted("HT12345678")
persons = client.ncid.search_wanted(risk_level="HIGH", status="ACTIVE")

# === BIOMETRICS ===
match = client.biometrics.verify(face_image="base64...", niu="HT12345678")
matches = client.biometrics.search(fingerprint="base64...")

# === ALPR ===
reads = client.alpr.search(plate="AB-123-CD", since="2026-01-01")
client.alpr.ingest_bulk(reads=[...])

# === GRAPH INTELLIGENCE ===
report = client.graphrag.generate_report(
    entity_id="HT12345678",
    report_type="ENTITY_PROFILE",
)

# === ALERTS ===
alerts = client.alerts.list(severity="CRITICAL", limit=20)

# === FINANCIAL ===
txns = client.financial.search_suspicious(amount_min=10000, currency="USD")

# === CYBER ===
iocs = client.cyber.search_iocs(ioc_type="IP", value="185.220.101.42")

# === WATCHLIST ===
entries = client.watchlist.search(category="TERRORISM")

# === AI FUSION ===
risk = client.ai.score_entity("HT12345678")
```

## Authentication

```python
# API Key
client = SNISIDClient(api_key="sk-...")

# JWT Token (Keycloak)
client = SNISIDClient(token="eyJ...")

# mTLS certificate
client = SNISIDClient(cert=("client.pem", "client-key.pem"))
```

## Async Support

```python
from sniside import AsyncSNISIDClient

async with AsyncSNISIDClient(api_key="sk-...") as client:
    results = await client.search("AB-123-CD")
    report = await client.graphrag.generate_report("HT12345678", "ENTITY_PROFILE")
```
