---
# ============================================================
# SNISID-Data — National Data Governance Model
# Rétention, Accès et DGO (Data Governance Office)
# Document ID: SNISID-DATA-GOV-001
# Version: 1.0.0
# ============================================================

## 1. LE DATA GOVERNANCE OFFICE (DGO)

Le DGO est l'entité légale et technique qui possède (Data Ownership) les tables du Lakehouse. Aucun ingénieur (même root) ne peut accéder aux données sans l'approbation du DGO, transcrite en règles "Policy-as-Code" (Zero Trust, Phase 6).

## 2. DATA RETENTION & ARCHIVAL

- **Données Transitoires (Logs applicatifs) :** Conservées 90 jours dans OpenSearch (Phase 6), puis archivées dans le Lakehouse (Iceberg/MinIO) sur des disques magnétiques lents (Cold Storage) pour 10 ans.
- **Données Biométriques / État Civil :** Rétention Perpétuelle (Immuable). Si une erreur est faite, on écrit un événement de correction (CQRS, Phase 2) mais on n'efface *jamais* l'événement d'origine.
- **Droit à l'oubli (Privacy) :** Strictement encadré. Réservé aux cas légaux (Ex: Changement d'identité autorisé par un Juge, Protection des témoins).

## 3. DATA ACCESS MODEL (ABAC)

L'accès est régi par l'ABAC (Attribute-Based Access Control).
- Un Analyste du Ministère de la Santé (MSPP) cherchant à croiser le registre d'État Civil avec les naissances ne verra que les "Communes" et les "Dates". 
- La colonne `nom` et `niu` (Numéro d'Identification) sera masquée à la volée (Dynamic Data Masking).

---
*Document ID: SNISID-DATA-GOV-001 | Approuvé par: Délégué à la Protection des Données (DPO)*
