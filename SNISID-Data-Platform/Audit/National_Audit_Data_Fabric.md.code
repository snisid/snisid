# National Audit Data Fabric

## Objectif
Centraliser la traçabilité nationale complète de SNISID.

## Données centralisées

| Domaine | Données | Centralisé |
|---|---|---:|
| IAM events | login, MFA, role change, token use | Oui |
| BPMN events | workflow start/end/task/decision | Oui |
| API calls | endpoint, subject, purpose, response class | Oui |
| PKI operations | cert issue/revoke/sign/verify | Oui |
| Administrative actions | création, modification, approbation | Oui |
| Data access | query, export, dashboard view | Oui |
| Data changes | insert/update/delete/merge | Oui |
| AI/ML decisions | model, version, score, explanation | Oui |

## Exigences

- Horodatage synchronisé NTP sécurisé.
- Identité acteur/service obligatoire.
- Correlation ID bout en bout.
- Stockage immutable/WORM pour événements critiques.
- Hachage et signature des lots d'audit.
- Conservation selon politique légale.
- Recherche et export probatoire contrôlé.

## Schéma audit minimal

```json
{
  "event_id": "uuid",
  "timestamp": "2026-05-25T12:00:00Z",
  "actor_id": "user-or-service",
  "actor_type": "USER|SERVICE|SYSTEM",
  "action": "DATA_READ|DATA_WRITE|API_CALL|LOGIN|MODEL_SCORE",
  "resource": "dataset/api/workflow",
  "purpose": "service_delivery|audit|fraud_control",
  "classification": "RESTRICTED",
  "decision": "ALLOW|DENY",
  "correlation_id": "uuid",
  "source_ip": "x.x.x.x",
  "hash": "sha256"
}
```

## Détection

- Accès massif inhabituel.
- Export hors heures normales.
- Accès sans finalité valide.
- Changement de rôle administrateur.
- Requêtes sur données critiques par compte dormant.
- Divergence entre consentement et usage.
