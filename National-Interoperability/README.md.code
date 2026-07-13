# SNISID — National Interoperability Platform (Phase 4)

## Objectif
La Plateforme d'Interopérabilité Nationale (X-Road Haïtien) est le backbone numérique de l'État. Elle connecte de manière souveraine, chiffrée et gouvernée l'ensemble des ministères et agences de la République d'Haïti. Elle garantit qu'une information (ex: un changement d'adresse ou un décès) saisie par une agence est instantanément et de manière sécurisée disponible pour les autres agences autorisées.

## Périmètre
1. **API Gateway Ecosystem (Kong)** : Point d'entrée unique, throttling, monétisation.
2. **Interoperability Bus (Kafka)** : Event-driven architecture asynchrone pour l'État.
3. **Data Exchange & Governance** : Modèle de données centralisé (Master Data) et règles d'échange.
4. **Identity Federation** : SSO inter-agences via Keycloak.
5. **Service Mesh (Istio)** : Chiffrement mTLS intra-cluster et gestion du trafic.
6. **BPMN Orchestration** : Processus cross-agences complexes.

## Standard d'Interopérabilité
- Tout échange doit transiter par le bus ou l'API Gateway (Aucune connexion Point-à-Point directe entre bases de données).
- Chaque consommateur doit être dûment authentifié (mTLS + OIDC).
- Tout échange est tracé (Immutable Audit Log).
