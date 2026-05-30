---
# ============================================================
# SNISID-Data — Sovereign AI & ML Platform
# Modèles d'Intelligence Artificielle de l'État
# Document ID: SNISID-DATA-AI-001
# Version: 1.0.0
# ============================================================

## 1. INTELLIGENCE ARTIFICIELLE SOUVERAINE

L'État Haïtien s'interdit d'utiliser des API Cloud externes (OpenAI, Anthropic, Google Cloud AI) pour traiter des données gouvernementales (Risque d'espionnage ou de fuite). Tous les modèles IA sont open-source et hébergés "On-Premise" sur le cluster GPU du Datacenter (Phase 5).

## 2. CAS D'USAGE DE L'IA GOUVERNEMENTALE

- **Reconnaissance Faciale (Computer Vision) :** Modèles convolutionnels (CNN) optimisés pour la détection anti-fraude (Liveness Detection, Phase 8) et l'ABIS central.
- **Analyse Prédictive du Crime (Machine Learning) :** Le SOC de la Police Nationale utilise des modèles de régression pour identifier des corrélations (Ex: Augmentation des vols de véhicules liée à des événements spécifiques).
- **LLM Gouvernemental (Generative AI) :** Déploiement de modèles type "Llama 3" ou "Mistral" via vLLM en interne. Permet aux agents de la Justice de résumer des centaines de pages de casiers judiciaires de manière instantanée, sans que le texte ne quitte jamais les serveurs sécurisés de l'État.

## 3. MLOps (Machine Learning Operations)

Le cycle de vie de l'IA (Entraînement, Déploiement, Inférence) est géré via MLflow et Kubeflow sur RKE2. Les données d'entraînement proviennent exclusivement des couches "Silver/Gold" du Lakehouse (Iceberg), anonymisées selon les règles du DGO.

---
*Document ID: SNISID-DATA-AI-001 | Approuvé par: Directeur de l'IA Gouvernementale*
