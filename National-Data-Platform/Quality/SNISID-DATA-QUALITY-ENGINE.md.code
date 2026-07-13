---
# ============================================================
# SNISID-Data — National Data Quality Engine
# Nettoyage, Déduplication et Data Contracts
# Document ID: SNISID-DATA-QUAL-001
# Version: 1.0.0
# ============================================================

## 1. LA PROBLÉMATIQUE DE LA QUALITÉ DES DONNÉES

En Haïti, les variations orthographiques des noms (Ex: "Jean-Baptiste" vs "Janbatis") et l'absence d'adresses standardisées créent des doublons. Le moteur de qualité (Quality Engine) nettoie le "Data Swamp" pour en faire un "Data Lakehouse" propre.

## 2. DATA CONTRACTS (Contrats de Données)

Plutôt que de nettoyer a posteriori, le SNISID impose des contrats stricts aux producteurs de données (Ex: La Police, Les Mairies).
Un contrat est un fichier YAML définissant exactement le schéma attendu.
Si un système tente de publier un événement Kafka où le champ `date_naissance` est "12 Mars 1990" au lieu de "1990-03-12" (ISO 8601), le message est rejeté dans une "Dead Letter Queue" (File d'attente d'erreur).

## 3. NETTOYAGE ET DÉDUPLICATION (Entity Resolution)

Le pipeline de qualité (Apache Spark) exécute des algorithmes heuristiques :
- **Phonétique (Soundex/Metaphone adapté au Créole) :** Rapprochement des orthographes.
- **Levenshtein Distance :** Détection de fautes de frappe.
- **Biometric Deduplication :** Le dernier recours infaillible (ABIS GPU, Phase 2) pour lier deux dossiers sous le même NIU si les empreintes matchent.

---
*Document ID: SNISID-DATA-QUAL-001 | Approuvé par: Chief Data Officer (CDO)*
