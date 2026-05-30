# 🕸 KAFKA EVENT MESH NATIONAL

> **Phase 3 / Étape 7** — Interconnexion gouvernementale temps réel
> Version : 1.0.0

---

## 1. Topologie Cible

```
                      ┌──────────────────────────────────────────┐
                      │           NATIONAL KAFKA EVENT MESH      │
                      │            (KRaft mode, mTLS, ACL)       │
                      └──────────┬──────────────────────┬────────┘
                                 │                      │
       ┌─────────────────────────┴───┐   ┌──────────────┴────────────────┐
       │   DC1 — Port-au-Prince      │   │   DC2 — Cap-Haïtien           │
       │   3× controllers + 6× brokers│  │   3× controllers + 6× brokers │
       │   Leader (writes)            │  │   Active stretch (writes)     │
       └─────────────────────────────┘   └───────────────────────────────┘
                                 │
                      ┌──────────┴─────────┐
                      │ DC3 — Les Cayes    │
                      │ Cluster MirrorMaker2│
                      │ DR + replay         │
                      └────────────────────┘

Domain Topic Groups
├── civil-registry.*    (Birth, Death, Marriage, Divorce, Adoption)
├── identity.*          (Enrollment, Verification, Revocation...)
├── judicial.*          (Validation, Suspension, Investigation...)
├── elections.*         (Voter, Candidate, Results)
├── immigration.*       (Entry, Visa, Deportation)
├── tax.*               (Registration, Declaration, Audit)
├── health.*            (Vaccination, Epidemic)
├── fraud.*             (Detection, Case)
├── audit.*             (Workflow, User, Data)
├── security.*          (Alerts, Access, Policy)
├── offline.*           (Batch, Conflict, Sync)
└── _internal.*         (DLQ, Compaction, Replay)
```

---

## 2. Standards de Nommage des Topics

```
<domain>.<entity>.<event-type>.v<MAJOR>
```

Exemples valides :
- `civil-registry.birth.created.v1`
- `identity.enrollment.completed.v1`
- `judicial.case.opened.v2`
- `fraud.detection.alert.v1`

Variantes :
- `*.dlq` — Dead letter
- `*.replay` — Replay manuel
- `*.compact` — Topic compaction (états)

---

## 3. Configuration des Topics (Standards)

| Profil | Cas d'usage | Partitions | RF | min.insync.replicas | Retention | Compaction |
|--------|-------------|------------|----|--------------------:|-----------|------------|
| **CRITICAL** | civil.*, identity.*, judicial.* | 24 | 3 | 2 | 10 ans (legal) | non |
| **HIGH** | fraud.*, elections.*, security.* | 12 | 3 | 2 | 5 ans | non |
| **MEDIUM** | tax.*, immigration.* | 6 | 3 | 2 | 2 ans | non |
| **LOW** | health.* (non sensible) | 3 | 2 | 2 | 1 an | non |
| **STATE** | `*.state` (snapshot) | 6 | 3 | 2 | infini | oui |
| **AUDIT** | audit.* | 24 | 3 | 3 | 30 ans + WORM | non |
| **DLQ** | *.dlq | 3 | 3 | 2 | 90 jours | non |

---

## 4. Sécurité

| Contrôle | Mise en œuvre |
|----------|---------------|
| Transport | TLS 1.3 obligatoire (mTLS) |
| Authentification clients | SPIFFE/SPIRE X.509 SVID |
| Authentification inter-broker | mTLS |
| Autorisation | ACL Kafka + OPA gateway |
| Chiffrement payload | AES-256-GCM (clé par domaine via KMS) |
| Signature | PKI sur header `signature` |
| Audit accès | Topic `_internal.kafka.access.audit` |

### ACL exemples

```bash
# Civil-Registry-Service peut écrire sur civil-registry.*
kafka-acls --add --producer \
  --topic 'civil-registry.*' \
  --principal "User:CN=civil-registry-service,O=SNISID,C=HT" \
  --resource-pattern-type PREFIXED

# Identity-Service peut lire civil-registry.birth.*
kafka-acls --add --consumer --group identity-service \
  --topic 'civil-registry.birth.' \
  --resource-pattern-type PREFIXED \
  --principal "User:CN=identity-service,O=SNISID,C=HT"
```

---

## 5. Versioning et Compatibilité

- **Schema Registry** (Confluent / Apicurio) en HA
- **Stratégie** : `BACKWARD_TRANSITIVE` par défaut
- **MAJOR** : nouveau topic suffixé `.vN`
- **MINOR/PATCH** : compatible sur le même topic
- **Deprecation** : durée minimale **18 mois** avant retrait

---

## 6. Replay & Résilience

| Capacité | Mise en œuvre |
|----------|---------------|
| Replay par offset | Outil `snisid-kafka-replay` |
| Replay par timestamp | `--from 2026-05-24T00:00:00Z` |
| DLQ retry exponentiel | 30 s → 5 min → 1 h → 24 h |
| MirrorMaker 2 DC1 ↔ DC2 | Active-Active |
| MirrorMaker 2 → DC3 | DR asynchrone |
| RPO | ≤ 60 s |

---

## 7. Quotas & Capacity

| Producteur | Quota / sec |
|------------|-------------|
| civil-registry | 5 000 msg/s |
| identity (read-heavy) | 10 000 msg/s |
| judicial | 1 000 msg/s |
| fraud (alerts) | 500 msg/s |
| audit | 20 000 msg/s |

Capacité initiale par broker : **200 MB/s in / 600 MB/s out**, **NVMe RAID-10**, **64 GB RAM**, **16 vCPU**.

---

## 8. Conventions de Headers Kafka

| Header | Description |
|--------|-------------|
| `event-id` | UUIDv7 unique |
| `event-type` | `civil.birth.created.v1` |
| `schema-id` | ID Schema Registry |
| `trace-id` | W3C Trace Context |
| `producer-spiffe` | SPIFFE ID producteur |
| `signature` | Signature PKI (base64) |
| `tsa-timestamp` | RFC 3161 |
| `correlation-id` | Workflow instance ID |
| `causation-id` | Event ayant déclenché celui-ci |

---

## 9. Observabilité Kafka

- **Cluster** : JMX → Prometheus → Grafana ("Kafka National Mesh")
- **Lag** : Burrow / Kafka Lag Exporter
- **Traces** : OpenTelemetry sur producers/consumers
- **Alertes** : `consumer_lag > 10000` → PagerDuty
- **SLO** : 99,95 % de disponibilité, latence p99 < 100 ms

---

## 10. Inventaire des Brokers (Production)

| ID | DC | Rôle |
|----|----|------|
| 1-3 | DC1 | Controller |
| 4-9 | DC1 | Broker |
| 10-12 | DC2 | Controller |
| 13-18 | DC2 | Broker |
| 19-24 | DC3 | DR brokers (MirrorMaker2 target) |

---

**Maintenu par :** Workflow Governance Office + Platform Engineering
