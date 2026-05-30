# Standards Data Nationaux SNISID

## Objectif
Éviter le chaos des données par des conventions nationales obligatoires.

## Standards

| Domaine | Standard |
|---|---:|
| Naming conventions | Oui |
| Schema evolution | Oui |
| Event schemas | Oui |
| Retention policies | Oui |
| Data quality rules | Oui |
| Metadata minimal | Oui |
| Classification | Oui |

## Naming conventions

### Datasets
`<domain>.<zone>.<entity>_<purpose>_v<major>`

Exemple : `identity.gold.citizen_profile_v1`

### Kafka topics
`snisid.<domain>.<entity>.<event>.v<major>`

Exemple : `snisid.identity.citizen.created.v1`

### Colonnes
- snake_case.
- suffixes standards : `_id`, `_ts`, `_date`, `_code`, `_hash`.
- pas d'acronymes non documentés.

## Évolution de schéma

| Changement | Compatibilité | Approbation |
|---|---|---|
| Ajouter champ optionnel | Compatible | Steward |
| Ajouter champ obligatoire | Breaking | Owner + consumers |
| Renommer champ | Breaking | Migration formelle |
| Supprimer champ | Breaking | Dépréciation préalable |
| Changer type | Breaking | Comité architecture |

## Standards événements

Chaque événement doit inclure :
- event_id,
- event_type,
- event_version,
- occurred_at,
- producer,
- correlation_id,
- subject_id,
- classification,
- payload,
- hash/signature pour critique.

## Politique qualité minimale

- Contrôle unicité sur identifiants maîtres.
- Contrôle format sur identifiants, dates, codes.
- Contrôle valeurs autorisées via référentiels.
- Contrôle référentiel sur foreign keys logiques.
- Contrôle fraîcheur pour données opérationnelles.
