---
# ============================================================
# SNISID-Interop — Data Governance Model
# Master Data, Lignage et Propriété
# Document ID: SNISID-DATA-GOV-001
# Version: 1.0.0
# ============================================================

## 1. GOUVERNANCE NATIONALE DES DONNÉES

Le chaos des données gouvernementales (ex: une adresse différente à la DGI, au Passeport et à l'ONI) est résolu par une gouvernance stricte (Master Data Management).

## 2. MODÈLE DE PROPRIÉTÉ (DATA OWNERSHIP)

1. **System of Record (SoR) :** L'unique base de données faisant autorité pour un champ donné.
2. **System of Reference :** Les autres agences qui gardent une copie (cache) de la donnée en lecture seule.

| Attribut | System of Record | Règle de Modification |
|----------|------------------|-----------------------|
| `Nom, Prénom, Sexe` | SNISID Identity | Juge Civil uniquement |
| `Adresse Résidence` | SNISID Identity | Mise à jour par citoyen (Portail) validée par Mairie |
| `NIF (Numéro Fiscal)` | DGI | Généré par la DGI, lié au NIU |
| `Passeport Actif` | Immigration | Généré par l'Immigration, lié au NIU |

## 3. DATA LINEAGE (Lignage de la donnée)

Pour chaque donnée critique (ex: un changement de nom suite à un mariage), le système SNISID conserve :
- `source_agency`: Quelle agence a fait la modification (ex: OEC Port-au-Prince).
- `agent_id`: Qui a cliqué.
- `timestamp`: Quand.
- `legal_basis`: Document justifiant (ex: Certificat de Mariage MAR-2026-001).
- `consumers_notified`: Liste des agences ayant accusé réception du changement (ex: DGI, CEP).

---
*Document ID: SNISID-DATA-GOV-001 | Approuvé par: Comité de Gouvernance Numérique*
