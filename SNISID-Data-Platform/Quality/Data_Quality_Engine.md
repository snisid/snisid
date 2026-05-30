# Data Quality Engine

## Objectif
Garantir la qualité nationale des données par contrôles automatisés et validation inter-agences.

## Contrôles obligatoires

| Contrôle | Description | Action en échec |
|---|---|---|
| Duplicate detection | Détection doublons | Quarantaine + MDM review |
| Missing values | Données manquantes | Blocage si champ critique |
| Format validation | Standards formats | Rejet ou correction contrôlée |
| Integrity checks | Cohérence référentielle | Blocage Silver/Gold |
| Cross-agency validation | Corrélation inter-agences | Alerte steward |
| Freshness | Fraîcheur données | Alerte pipeline |
| Range checks | Valeurs plausibles | Quarantaine ligne |
| Consent/purpose check | Usage autorisé | Blocage accès |

## Seuils de qualité

| Classification | Score minimal | Publication Gold |
|---|---:|---|
| NATIONAL_SECRET | 99.9% | Approbation manuelle + automatique |
| RESTRICTED | 99.0% | Approbation owner/steward |
| CONFIDENTIAL | 97.0% | Automatique si conforme |
| INTERNAL | 95.0% | Automatique |
| PUBLIC | 98.0% | Validation publication |

## Workflow qualité

1. Profilage automatique.
2. Application règles standard + règles domaine.
3. Calcul score qualité.
4. Quarantaine anomalies.
5. Notification steward/owner.
6. Correction ou dérogation approuvée.
7. Publication uniquement si seuil atteint.

## KPI

- DQ score par domaine.
- Nombre d'anomalies critiques.
- Temps moyen de résolution.
- Doublons détectés/fusionnés.
- Taux de publication Gold conforme.
