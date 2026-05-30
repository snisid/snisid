# OFFLINE MOBILE SYNC ENGINE — SNISID
## Moteur de Synchronisation Mobile Résilient

| Attribut | Valeur |
|----------|--------|
| Document ID | SNISID-PH12-OFF-001 |
| Version | 1.0 |
| Statut | APPROUVÉ |

---

## 1. PRÉSENTATION

Moteur de synchronisation mobile conçu pour survivre aux longues coupures internet et garantir l'intégrité des données terrain dans les conditions les plus difficiles d'Haïti.

### 1.1 Défis Haïtiens

```
┌──────────────────────────────────────────────┐
│        CONDITIONS RÉELLES                     │
├──────────────────────────────────────────────┤
│  🌐 Coupures internet: jusqu'à 7 jours       │
│  📶 Réseau: 2G/3G intermittent               │
│  ⚡ Pas d'électricité: recharge solaire      │
│  🌧️ Intempéries: cyclones, inondations       │
│  🏔️ Zones montagneuses: sans signal          │
│  📱 Appareils: entrée/milieu de gamme        │
└──────────────────────────────────────────────┘
```

---

## 2. ARCHITECTURE

### 2.1 Sync Engine Components

```
┌─────────────────────────────────────────────┐
│           OFFLINE SYNC ENGINE                │
├─────────────────────────────────────────────┤
│  ┌──────────────────────────────────────┐   │
│  │      SYNC COORDINATOR                │   │
│  │  • Priority Queue                    │   │
│  │  • Retry Logic                       │   │
│  │  • Network Monitor                   │   │
│  └──────────────────────────────────────┘   │
│                                             │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐    │
│  │  Cache   │ │ Conflict │ │  Sync    │    │
│  │ Manager  │ │Resolver  │ │ Protocol │    │
│  └──────────┘ └──────────┘ └──────────┘    │
│                                             │
│  ┌──────────────────────────────────────┐   │
│  │      STORAGE LAYER                   │   │
│  │  ┌──────────┐ ┌──────────┐          │   │
│  │  │  Local   │ │  Event   │          │   │
│  │  │  Cache   │ │  Buffer  │          │   │
│  │  └──────────┘ └──────────┘          │   │
│  └──────────────────────────────────────┘   │
│                                             │
│  ┌──────────────────────────────────────┐   │
│  │      NETWORK LAYER                    │   │
│  │  • Connection Monitor                │   │
│  │  • Bandwidth Estimator               │   │
│  │  • Resume Support                    │   │
│  └──────────────────────────────────────┘   │
└─────────────────────────────────────────────┘
```

### 2.2 Data Flow

```
┌──────────┐    ┌──────────┐    ┌──────────┐
│  App     │───▶│  Write   │───▶│  Is      │
│  Action  │    │  Local   │    │  Online? │
└──────────┘    └──────────┘    └────┬─────┘
                                     │
                          ┌──────────┴──────────┐
                          │                     │
                    ┌─────▼─────┐         ┌─────▼─────┐
                    │   Yes     │         │    No     │
                    └─────┬─────┘         └─────┬─────┘
                          │                     │
                    ┌─────▼─────┐         ┌─────▼─────┐
                    │  Sync     │         │  Queue    │
                    │  Now      │         │  For Sync │
                    └─────┬─────┘         └─────┬─────┘
                          │                     │
                    ┌─────▼─────┐         ┌─────▼─────┐
                    │  Success? │         │  Notify   │
                    └─────┬─────┘         │  Pending  │
                          │               └───────────┘
               ┌──────────┴──────────┐
               │                     │
          ┌────▼────┐          ┌────▼────┐
          │   Yes   │          │   No    │
          └─────────┘          │  Retry  │
                               │  Later  │
                               └─────────┘
```

---

## 3. OFFLINE CACHING

### 3.1 Cache Strategy

| Stratégie | Description | Usage |
|-----------|-------------|-------|
| **Cache First** | Lire depuis le cache, refresh async | Profils, ID Cards |
| **Write Behind** | Écrire dans le cache, sync plus tard | Enrollments |
| **Cache Aside** | Lire depuis cache, miss → API → cache | Références |
| **Stale While Revalidate** | Servir cache, refresh en arrière-plan | Listes |

### 3.2 Cache Invalidation

| Événement | Action | Impact |
|-----------|--------|--------|
| Sync success | Invalider les items syncés | Minimal |
| Time-to-live expiré | Rafraîchir au prochain online | Cache stale |
| Force refresh | Tout invalider | Re-download |
| Conflict detected | Invalider version locale | Résolution |

### 3.3 Cache Tiers

```
┌──────────────────────────────────────────────┐
│              CACHE TIERS                     │
├──────────────────────────────────────────────┤
│  Tier 1 - Critical (50 MB max)               │
│  • User profile, ID Card                     │
│  • Current session data                      │
│  • App configuration                         │
│  Priority: Never evict                       │
│                                              │
│  Tier 2 - Important (200 MB max)             │
│  • Recent certificates                       │
│  • Recent notifications                      │
│  • Reference data (cities, depts)            │
│  Priority: LRU eviction                      │
│                                              │
│  Tier 3 - Normal (500 MB max)                │
│  • Historical data                           │
│  • Old notifications                         │
│  • Media files                               │
│  Priority: LRU + size eviction              │
└──────────────────────────────────────────────┘
```

---

## 4. DELAYED SYNC

### 4.1 Sync Queue

| Priorité | Max Delay | Retry | Data Types |
|----------|-----------|-------|------------|
| **P0 - Critical** | 1 min (si réseau) | ∞ | Emergency ops, Alerts |
| **P1 - High** | 1 hour | 10 attempts | Enrollments, Bio data |
| **P2 - Normal** | 24 hours | 5 attempts | ID Verifications |
| **P3 - Low** | 7 days | 3 attempts | Audit logs |
| **P4 - Bulk** | 30 days | 1 attempt | Reports, statistics |

### 4.2 Retry Logic

```
┌──────────────────────────────────────────────┐
│           EXPONENTIAL BACKOFF                │
├──────────────────────────────────────────────┤
│  Attempt 1: 1 min                            │
│  Attempt 2: 5 min                            │
│  Attempt 3: 15 min                           │
│  Attempt 4: 30 min                           │
│  Attempt 5: 1 hour                           │
│  Attempt 6: 2 hours                          │
│  Attempt 7: 4 hours                          │
│  Attempt 8: 8 hours                          │
│  Attempt 9: 12 hours                         │
│  Attempt 10+: 24 hours                       │
│                                              │
│  Jitter: ±20% random                         │
│  Max retries: ∞ for P0, configurable         │
└──────────────────────────────────────────────┘
```

---

## 5. CONFLICT RESOLUTION

### 5.1 Conflict Types

| Type | Example | Strategy |
|------|---------|----------|
| **Create** | Même ID créé offline/serveur | UUID conflict → Merge |
| **Update** | Même champ modifié 2x | Last Write Wins (LWW) |
| **Delete** | Supprimé localement mais modifié serveur | Delete wins |
| **Schema** | Version différente du modèle | Schema migration |

### 5.2 Resolution Algorithm

```
┌──────────────────────────────────────────────┐
│         CONFLICT RESOLUTION                  │
├──────────────────────────────────────────────┤
│  For each conflicting item:                  │
│                                              │
│  1. Compare timestamps                       │
│     - If 1 version newer: take newest        │
│     - If equal timestamps: go to step 2      │
│                                              │
│  2. Compare field-level changes              │
│     - If non-overlapping: merge fields       │
│     - If overlapping fields: go to step 3    │
│                                              │
│  3. Apply business rules:                    │
│     - Identity: official record wins         │
│     - Enrollment: latest biometric wins      │
│     - Audit: no conflict, append both        │
│                                              │
│  4. If still unresolved:                     │
│     - Flag for manual review                 │
│     - Notify admin                           │
│     - Keep both versions                     │
└──────────────────────────────────────────────┘
```

### 5.3 Version Vector

```json
{
  "document_id": "ID-12345",
  "version": {
    "local": { "counter": 5, "timestamp": "2026-05-25T10:00:00Z" },
    "server": { "counter": 7, "timestamp": "2026-05-25T12:00:00Z" }
  },
  "conflict": false,
  "last_sync": "2026-05-24T08:00:00Z",
  "fields": {
    "name": { "local": "Jean", "server": "Jean-Marie", "resolution": "server" },
    "address": { "local": "P-au-P", "server": "P-au-P", "resolution": "no_conflict" }
  }
}
```

---

## 6. SECURE SYNC

### 6.1 Sync Security

```
┌──────────────────────────────────────────────┐
│           SECURE SYNC PROTOCOL               │
├──────────────────────────────────────────────┤
│  Step 1: Device Authentication               │
│  • Device attestation token                  │
│  • User JWT                                  │
│  • Sync session key exchange                 │
│                                              │
│  Step 2: Data Encryption                     │
│  • Compress payload (gzip)                   │
│  • Encrypt with session key (AES-256-GCM)   │
│  • Add HMAC for integrity                    │
│                                              │
│  Step 3: Transfer                            │
│  • HTTPS with certificate pinning            │
│  • Chunked transfer (1MB chunks)            │
│  • Resume support (byte ranges)             │
│                                              │
│  Step 4: Server Verification                 │
│  • Verify HMAC                              │
│  • Decrypt payload                          │
│  • Process and acknowledge                  │
│  • Return checksum                          │
│                                              │
│  Step 5: Client Acknowledgment               │
│  • Verify server checksum                    │
│  • Mark items as synced                      │
│  • Remove from buffer                        │
└──────────────────────────────────────────────┘
```

### 6.2 Sync Payload Format

```json
{
  "sync_id": "sync-20260525-abc123",
  "device_id": "DEV-98765",
  "user_id": "USR-12345",
  "timestamp": "2026-05-25T14:30:00Z",
  "items": [
    {
      "type": "enrollment",
      "id": "ENR-001",
      "action": "create",
      "data_hash": "sha256:abc...",
      "payload": "<encrypted_base64>"
    }
  ],
  "checksum": "sha256:def...",
  "signature": "sig:ghi..."
}
```

---

## 7. MONITORING

### 7.1 Sync Metrics

| Métrique | Seuil | Alerte |
|----------|-------|--------|
| Sync Success Rate | > 98% | < 95% |
| Queue Depth | < 1000 | > 5000 |
| Sync Latency (P0) | < 1 min | > 5 min |
| Conflict Rate | < 1% | > 5% |
| Storage Usage | < 80% | > 90% |
| Battery Impact | < 5%/h | > 10%/h |

### 7.2 Sync Dashboard

```
┌──────────────────────────────────────────────┐
│         SYNC ENGINE DASHBOARD                │
├──────────────────────────────────────────────┤
│  📊 Overall Sync Health: 🟢 99.2%           │
│                                              │
│  Active Devices: 1,234                       │
│  Queue Depth:  56  (threshold: 5000)        │
│  Today's Syncs: 12,345                      │
│  Conflicts: 12 (0.1%)                       │
│                                              │
│  ┌──────────────────────────────────────┐   │
│  │  Sync Success Rate (24h)             │   │
│  │  ██████████████████████████ 99.2%   │   │
│  └──────────────────────────────────────┘   │
│                                              │
│  Recent Failures: 3                         │
│  • DEV-1234: Network timeout (2m ago)       │
│  • DEV-5678: Conflict (5m ago)              │
│  • DEV-9012: Auth failed (10m ago)          │
└──────────────────────────────────────────────┘
```

---

## 8. PERFORMANCE

| Métrique | Cible | Pire Cas |
|----------|-------|----------|
| Sync 1 record (online) | < 100ms | < 1s |
| Sync 1000 records | < 10s | < 30s |
| Conflict resolution | < 50ms | < 200ms |
| Cache read | < 10ms | < 50ms |
| Offline queue capacity | 100,000 records | 500,000 |
| Storage overhead | < 10% of data | < 20% |
| Sync resume efficiency | 95% | 80% |

---

## 9. RÉSILIENCE

| Failure Mode | Recovery | Impact |
|-------------|----------|--------|
| Network loss mid-sync | Resume on reconnect | Minimal |
| App crash during sync | Atomic transactions | None |
| Storage full | Auto-cleanup LRU | Data loss (non-critical) |
| Data corruption | Checksum + recovery | Minimal |
| Power loss | Transaction rollback | None |
| Sync server down | Local queue | Extended offline |

---
*Fin du document — Offline Mobile Sync Engine v1.0*