---
# ============================================================
# SNISID-Data — Data Catalog & Lineage
# Traçabilité et Dictionnaire de Données
# Document ID: SNISID-DATA-META-001
# Version: 1.0.0
# ============================================================

## 1. DATA CATALOG (Le Dictionnaire de l'État)

Avec des milliers de tables dans le Lakehouse, un data scientist ne peut pas "deviner" où se trouve l'information. Le SNISID utilise un Data Catalog (ex: DataHub ou Amundsen).
Ce catalogue permet de rechercher des termes métier (Ex: "Trouver la table contenant le numéro de passeport").

## 2. DATA LINEAGE (Traçabilité)

Le système conserve la généalogie de chaque donnée.
Si un tableau de bord (Dashboard) montre qu'il y a 10 millions d'électeurs :
- Le Lineage montre que ce chiffre provient de la vue `Gold_Electeurs`...
- ...qui elle-même a été calculée à partir de la table `Silver_Citoyens`...
- ...qui elle-même a été ingérée depuis le topic Kafka `IdentityCreated` (Phase 4).

## 3. IMPACT ANALYSIS

Le Lineage permet aussi d'évaluer l'impact d'un changement.
Si un développeur veut supprimer la colonne `adresse_ancienne` de la base de données source, le catalogue lui indique immédiatement que "Ce champ est utilisé par 3 modèles d'IA et 2 Dashboards de la Police. La suppression casserait ces systèmes en production."

---
*Document ID: SNISID-DATA-META-001 | Approuvé par: Chief Data Officer*
