# Metadata Platform — Catalogue national des données

## Objectif
Créer un catalogue national permettant découverte, gouvernance, ownership, classification, sensibilité et contrats de données.

## Capacités obligatoires

| Fonction | Support |
|---|---:|
| Data discovery | Oui |
| Classification | Oui |
| Ownership | Oui |
| Sensitivity levels | Oui |
| Data contracts | Oui |
| Glossaire métier | Oui |
| Lineage intégré | Oui |
| Certification dataset | Oui |

## Métadonnées minimales obligatoires

| Champ | Obligatoire | Exemple |
|---|---:|---|
| dataset_id | Oui | snisid.citizens.gold_citizen_profile |
| business_name | Oui | Profil citoyen officiel |
| data_owner | Oui | Direction identité nationale |
| data_steward | Oui | Steward registre citoyen |
| classification | Oui | RESTRICTED |
| sensitivity | Oui | PII, national_id |
| retention_policy | Oui | LIFE_PLUS_10 |
| legal_basis | Oui | Loi/mandat administratif |
| allowed_purposes | Oui | service_delivery, fraud_control |
| source_systems | Oui | civil_registry, identity |
| dq_score | Oui | 98.5 |
| lineage_url | Oui | lien OpenLineage |
| access_policy | Oui | policy://citizens/restricted |

## États de certification

| Statut | Description | Utilisation |
|---|---|---|
| Draft | En cours de création | Non exploitable officiellement |
| Registered | Metadata minimale complète | Usage interne contrôlé |
| Certified | Validé owner/steward/DQ/security | Usage officiel |
| Deprecated | Remplacé | Lecture limitée |
| Retired | Supprimé/archivé | Accès exceptionnels |

## Règle anti-orphelin
Tout objet data sans owner, classification ou rétention est automatiquement placé en quarantaine et bloqué pour usage analytique.
