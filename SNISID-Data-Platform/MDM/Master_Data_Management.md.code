# Master Data Management — Source de vérité nationale

## Objectif
Établir les **golden records** nationaux et éliminer les vérités concurrentes entre agences.

## Domaines MDM prioritaires

| Domaine | Criticité | Identifiant maître | Propriétaire recommandé |
|---|---|---|---|
| Citizens | CRITIQUE | National Citizen ID | Autorité nationale identité/registre |
| Identity | CRITIQUE | Identity Master ID | Agence identité numérique |
| Civil Registry | CRITIQUE | Civil Registry Record ID | État civil |
| Agencies | ÉLEVÉ | Agency ID | Administration centrale |
| Devices | ÉLEVÉ | Device ID | Autorité technique SNISID |

## Fonctions MDM

- Résolution d'identité.
- Détection et fusion de doublons.
- Survivorship rules.
- Gestion d'historique et changements.
- Référentiels officiels et listes contrôlées.
- Publication API des golden records.
- Workflow d'arbitrage par data steward.

## Règles de golden record

| Règle | Description |
|---|---|
| Unicité | Une entité nationale = une vérité officielle |
| Traçabilité | Toute fusion/séparation est justifiée et auditée |
| Survivorship | Les champs officiels suivent règles de priorité par source |
| Validation croisée | Les données critiques sont vérifiées avec agences sources |
| Non-destruction | Les anciennes valeurs sont historisées |

## Exemple de survivorship Citizens

| Attribut | Source prioritaire | Source secondaire | Validation |
|---|---|---|---|
| Nom légal | Civil Registry | Identity System | Document officiel |
| Date de naissance | Civil Registry | Passport/ID | Contrôle format + cohérence |
| Adresse | Citizen Services | Tax/municipality | Fraîcheur et preuve |
| Statut vital | Civil Registry | Health/Municipality | Validation inter-agence |

## KPI MDM

- Taux de doublons citoyens.
- Nombre de conflits ouverts.
- Délai moyen de résolution steward.
- Taux de golden records certifiés.
- Nombre de modifications critiques auditées.
