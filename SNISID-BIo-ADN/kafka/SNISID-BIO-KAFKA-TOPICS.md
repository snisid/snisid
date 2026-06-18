# SNISID-BIO-ADN — Kafka Topics & Schémas Avro
**Document ID :** SNISID-BIO-KFK-001 | **Version :** 1.0.0

---

## 1. TOPICS KAFKA — Vue d'ensemble

| Topic | Partitions | Rétention | Producteur | Consommateurs |
|-------|-----------|-----------|------------|---------------|
| `snisid.bio.profile.created` | 6 | 30 jours | LDIS labs | SDIS sync, audit |
| `snisid.bio.profile.uploaded` | 6 | 30 jours | SDIS sync | NDIS matching |
| `snisid.bio.hits` | 12 | 90 jours | Matching Engine | Alertes, DCPJ, audit |
| `snisid.bio.wanted.events` | 12 | 90 jours | PNH/DCPJ | PER-REC index, FPR |
| `snisid.bio.missing.events` | 6 | 90 jours | PNH/citoyens | PER-DIS index |
| `snisid.bio.vehicle.stolen` | 6 | 30 jours | PNH/FOVeS | BIE-VEH, LAPI |
| `snisid.bio.document.stolen` | 6 | 30 jours | ONI/PNH | BIE-DOC, frontières |
| `snisid.bio.vessel.stolen` | 3 | 30 jours | CGFADH/PNH | BIE-EMB |
| `snisid.bio.arm.stolen` | 3 | 90 jours | DCPJ | BIE-ARM |
| `snisid.bio.lapi.query` | 24 | 1 jour | LAPI (MP-16) | BIE-VEH, BIE-PLQ |
| `snisid.bio.expunge.events` | 3 | 365 jours | dcpj.director | Tous les index |
| `snisid.bio.audit.events` | 6 | 365 jours | Triggers PG | Audit service |

---

## 2. SCHÉMAS AVRO

### 2.1 DNAProfileCreated

```json
{
  "type": "record",
  "name": "DNAProfileCreated",
  "namespace": "ht.gov.snisid.bio",
  "fields": [
    {"name": "event_id",       "type": "string"},
    {"name": "event_type",     "type": "string", "default": "DNAProfileCreated"},
    {"name": "sample_id",      "type": "string"},
    {"name": "specimen_number","type": "string"},
    {"name": "index_type",     "type": {
      "type": "enum", "name": "IndexType",
      "symbols": ["BIO_CON","BIO_ARR","BIO_FSC","BIO_DIS","BIO_RNI"]
    }},
    {"name": "lab_id",         "type": "string"},
    {"name": "lab_level",      "type": "string"},
    {"name": "quality_score",  "type": "float"},
    {"name": "loci_count",     "type": "int"},
    {"name": "amelogenin",     "type": ["null","string"], "default": null},
    {"name": "case_number",    "type": ["null","string"], "default": null},
    {"name": "collected_date", "type": "string"},
    {"name": "correlation_id", "type": "string"},
    {"name": "timestamp",      "type": "long",   "logicalType": "timestamp-millis"}
  ]
}
```

### 2.2 DNAHitDetected

```json
{
  "type": "record",
  "name": "DNAHitDetected",
  "namespace": "ht.gov.snisid.bio",
  "fields": [
    {"name": "event_id",         "type": "string"},
    {"name": "hit_id",           "type": "string"},
    {"name": "query_sample_id",  "type": "string"},
    {"name": "match_sample_id",  "type": "string"},
    {"name": "match_type",       "type": {
      "type": "enum", "name": "MatchType",
      "symbols": ["FULL_MATCH","PARTIAL","FAMILIAL"]
    }},
    {"name": "confidence",       "type": "float"},
    {"name": "matched_loci",     "type": "int"},
    {"name": "total_loci",       "type": "int"},
    {"name": "hit_level",        "type": "string"},
    {"name": "alert_level",      "type": {
      "type": "enum", "name": "AlertLevel",
      "symbols": ["LOW","MEDIUM","HIGH","CRITICAL"]
    }},
    {"name": "query_index_type", "type": "string"},
    {"name": "match_index_type", "type": "string"},
    {"name": "case_number",      "type": ["null","string"], "default": null},
    {"name": "timestamp",        "type": "long", "logicalType": "timestamp-millis"}
  ]
}
```

### 2.3 WantedPersonCreated

```json
{
  "type": "record",
  "name": "WantedPersonCreated",
  "namespace": "ht.gov.snisid.bio",
  "fields": [
    {"name": "event_id",          "type": "string"},
    {"name": "record_id",         "type": "string"},
    {"name": "record_number",     "type": "string"},
    {"name": "niu",               "type": ["null","string"], "default": null},
    {"name": "warrant_type",      "type": "string"},
    {"name": "danger_level",      "type": "string"},
    {"name": "armed_dangerous",   "type": "boolean"},
    {"name": "charges",           "type": {"type": "array", "items": "string"}},
    {"name": "entering_agency",   "type": "string"},
    {"name": "interpol_notice",   "type": ["null","string"], "default": null},
    {"name": "expiry_date",       "type": ["null","string"], "default": null},
    {"name": "correlation_id",    "type": "string"},
    {"name": "timestamp",         "type": "long", "logicalType": "timestamp-millis"}
  ]
}
```

### 2.4 LAPIQuery (temps réel)

```json
{
  "type": "record",
  "name": "LAPIPlateQuery",
  "namespace": "ht.gov.snisid.lapi",
  "fields": [
    {"name": "query_id",      "type": "string"},
    {"name": "plate_number",  "type": "string"},
    {"name": "camera_id",     "type": "string"},
    {"name": "location",      "type": "string"},
    {"name": "image_ref",     "type": ["null","string"], "default": null},
    {"name": "timestamp",     "type": "long", "logicalType": "timestamp-millis"}
  ]
}
```

### 2.5 LAPIQueryResponse

```json
{
  "type": "record",
  "name": "LAPIPlateResponse",
  "namespace": "ht.gov.snisid.lapi",
  "fields": [
    {"name": "query_id",      "type": "string"},
    {"name": "plate_number",  "type": "string"},
    {"name": "hit_found",     "type": "boolean"},
    {"name": "hit_type",      "type": ["null","string"], "default": null},
    {"name": "record_number", "type": ["null","string"], "default": null},
    {"name": "alert_level",   "type": ["null","string"], "default": null},
    {"name": "mco_contact",   "type": ["null","string"], "default": null},
    {"name": "response_ms",   "type": "int"},
    {"name": "timestamp",     "type": "long", "logicalType": "timestamp-millis"}
  ]
}
```

---

## 3. CONSUMER GROUPS

```yaml
# consumers/bio-adn-consumers.yaml

consumers:
  - group_id: snisid-bio-sdis-sync
    topic: snisid.bio.profile.created
    description: Synchronisation LDIS→SDIS des nouveaux profils ADN
    sla_max_lag: 1000   # messages max en retard

  - group_id: snisid-bio-ndis-matcher
    topic: snisid.bio.profile.uploaded
    description: Déclenchement du matching NDIS après upload SDIS
    sla_max_lag: 500

  - group_id: snisid-bio-alert-dispatcher
    topic: snisid.bio.hits
    description: Dispatch alertes hits ADN vers agences concernées
    sla_max_lag: 100    # critique — alertes temps réel

  - group_id: snisid-lapi-responder
    topic: snisid.bio.lapi.query
    description: Réponse requêtes LAPI < 200ms
    sla_max_lag: 50     # ultra-critique — temps réel terrain

  - group_id: snisid-bio-audit-writer
    topics:
      - snisid.bio.hits
      - snisid.bio.wanted.events
      - snisid.bio.expunge.events
    description: Écriture audit log immuable
    sla_max_lag: 5000   # moins critique
```

---

## 4. TOPOLOGIE KAFKA POUR BIO-ADN

```
LDIS Lab (Port-au-Prince)
    │  bio.profile.created
    ▼
Kafka Cluster (SNISID infra)
    │
    ├──► SDIS Sync Consumer ──► SDIS PostgreSQL ──► bio.profile.uploaded
    │
    ├──► NDIS Matcher Consumer ──► Matching Engine
    │         │
    │         └──► bio.hits ──► Alert Dispatcher ──► DIDComm / Email / PNH App
    │
    ├──► PNH/DCPJ producers ──► bio.wanted.events ──► PER-REC index
    │
    ├──► LAPI (MP-16) ──► bio.lapi.query ──► BIE-VEH responder ──► bio.lapi.response
    │
    └──► Audit Writer ──► bio_audit_log (PostgreSQL immuable)
```
