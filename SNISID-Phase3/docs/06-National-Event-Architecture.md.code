# ⚡ NATIONAL EVENT ARCHITECTURE

> **Phase 3 / Étape 6** — Système nerveux numérique de l'État.
> Version : 1.0.0

---

## 1. Vision

SNISID est un **Event-Driven Government**.
Chaque acte administratif génère un **événement immuable**, signé, versionné, distribué via Kafka, consommé par tous les acteurs autorisés.

> Pas d'événement = pas d'acte. Pas d'acte = pas d'événement.

---

## 2. Principes Event-Driven

| Principe | Description |
|----------|-------------|
| **Event = source of truth** | L'état est dérivé des événements (event sourcing) |
| **Append-only** | Aucun événement n'est modifié ou supprimé |
| **Immuable + signé** | Signature PKI + hash chaîné Merkle |
| **Versionné** | Schéma Avro versionné (compatibilité BACKWARD_TRANSITIVE) |
| **Replay-able** | Reconstruction d'état possible via replay |
| **Multi-consumers** | Un événement → N consommateurs (registres, audit, BI, SIEM) |
| **Idempotent** | Identifiant unique `eventId` (UUIDv7) |

---

## 3. Anatomie d'un Événement National

```json
{
  "eventId": "01HXY8K3M2W5N9P7R6T4QZ1V0E",
  "eventType": "civil.birth.created.v1",
  "occurredAt": "2026-05-24T14:35:22.123Z",
  "producer": {
    "service": "civil-registry-service",
    "version": "2.3.1",
    "node": "dc1-pap-node-04",
    "spiffeId": "spiffe://snisid.ht/civil-registry"
  },
  "correlation": {
    "workflowId": "civil-registry.birth.simple",
    "workflowInstanceId": "...",
    "traceId": "00-...",
    "causationId": "..."
  },
  "subject": {
    "personId": "HT-NIN-..."
  },
  "payload": { /* spécifique au type */ },
  "context": {
    "agentId": "...",
    "deviceId": "...",
    "location": { "commune": "Port-au-Prince", "department": "Ouest" }
  },
  "integrity": {
    "hash": "sha384:...",
    "prevHash": "sha384:...",
    "signature": "pki:...",
    "tsa": "rfc3161:..."
  },
  "schemaVersion": "1.0.0"
}
```

---

## 4. Catalogue des Domaines Événementiels

### 4.1 Civil Registry (`civil.*`)

| Événement | Producer | Consommateurs |
|-----------|----------|----------------|
| `civil.birth.declared.v1` | Mairie | Civil-Registry-Svc |
| `civil.birth.created.v1` | Civil-Registry-Svc | Identity-Svc, Health-Svc, Stats, BI |
| `civil.birth.recognized.v1` | Civil-Registry-Svc | Identity-Svc, Audit, BI |
| `civil.death.registered.v1` | Civil-Registry-Svc | Identity-Svc (révoke), Tax, Pension, Banking |
| `civil.marriage.civil.registered.v1` | Civil-Registry-Svc | Identity-Svc, Tax, Notaries |
| `civil.divorce.judicial.finalized.v1` | Court-Svc | Identity-Svc, Notaries |
| `civil.adoption.national.completed.v1` | Court-Svc | Identity-Svc, IBESR |

### 4.2 Identity (`identity.*`)

| Événement | Description |
|-----------|-------------|
| `identity.enrollment.started.v1` | Démarrage enrôlement |
| `identity.enrolled.v1` | NIN attribué |
| `identity.verified.v1` | Vérification 1:1 réussie |
| `identity.duplicate.detected.v1` | Doublon détecté |
| `identity.revoked.v1` | Révocation |
| `identity.suspended.v1` | Suspension |
| `identity.corrected.v1` | Correction appliquée |
| `identity.appeal.opened.v1` | Contestation citoyen |

### 4.3 Judicial (`judicial.*`)

| Événement | Description |
|-----------|-------------|
| `judicial.flagged.v1` | Drapeau judiciaire posé |
| `judicial.validated.v1` | Acte validé légalement |
| `judicial.case.opened.v1` | Dossier ouvert |
| `judicial.case.closed.v1` | Dossier clôturé |
| `judicial.order.suspension.v1` | Ordre suspension |
| `judicial.appeal.filed.v1` | Appel déposé |

### 4.4 Elections (`elections.*`)

| Événement | Description |
|-----------|-------------|
| `elections.voter.registered.v1` | Inscription |
| `elections.voter.validated.v1` | Validation électeur |
| `elections.candidate.approved.v1` | Candidat agréé |
| `elections.results.published.v1` | Résultats publiés |

### 4.5 Fraud (`fraud.*`)

| Événement | Description |
|-----------|-------------|
| `fraud.detected.v1` | Score > seuil |
| `fraud.case.opened.v1` | Enquête ouverte |
| `fraud.case.confirmed.v1` | Fraude confirmée |
| `fraud.case.dismissed.v1` | Classée sans suite |

### 4.6 Security (`security.*`)

| Événement | Description |
|-----------|-------------|
| `security.alert.raised.v1` | Alerte SIEM |
| `security.access.denied.v1` | Tentative refusée |
| `security.policy.violation.v1` | Violation OPA |

### 4.7 Audit (`audit.*`)

| Événement | Description |
|-----------|-------------|
| `audit.workflow.transition.v1` | Transition workflow |
| `audit.user.action.v1` | Action utilisateur |
| `audit.data.access.v1` | Accès donnée sensible |

### 4.8 Offline (`offline.*`)

| Événement | Description |
|-----------|-------------|
| `offline.batch.uploaded.v1` | Batch reçu |
| `offline.conflict.detected.v1` | Conflit CRDT |
| `offline.sync.completed.v1` | Sync terminée |

---

## 5. Règles d'Émission

> Tous les workflows **DOIVENT** :
> 1. Émettre `*.started` au démarrage
> 2. Émettre `*.<state>` à chaque transition majeure
> 3. Émettre `*.completed` ou `*.failed` à la fin
> 4. Toujours via service task `kafka.emit`
> 5. Signature PKI du payload obligatoire
> 6. Schéma Avro enregistré au Schema Registry

---

## 6. Schémas Avro (extraits)

Voir `kafka/schemas/` pour les définitions Avro complètes :
- `civil.birth.created.v1.avsc`
- `identity.enrolled.v1.avsc`
- `judicial.flagged.v1.avsc`
- ...

---

## 7. Patterns Recommandés

| Pattern | Usage |
|---------|-------|
| **Outbox Pattern** | Garantir l'émission après commit DB |
| **Saga (Choreography)** | Coordination cross-domain |
| **CQRS** | Lecture optimisée (read models) |
| **Event Sourcing** | Source of truth = log événements |
| **Dead Letter Queue** | `*.dlq` par topic |
| **Compensating Events** | `*.compensated.v1` pour annuler |

---

## 8. Anti-patterns à Bannir

| ❌ Anti-pattern | Pourquoi |
|----------------|----------|
| Modifier un événement émis | Casse l'immuabilité |
| `delete` d'événement | Casse l'audit |
| Schéma sans version | Casse les consommateurs |
| Payload non signé | Inadmissible légalement |
| Émettre depuis un endpoint REST direct (pas de workflow) | Casse la gouvernance |

---

**Maintenu par :** Workflow Governance Office + Event Architecture Board
