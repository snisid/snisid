# SNISID National API & Interoperability Architecture

## 1. Introduction
Cette architecture définit le socle de l'interopérabilité souveraine pour l'État Haïtien (SNISID). Elle permet l'échange fluide, sécurisé et standardisé de données entre toutes les agences gouvernementales.

## 2. Vue d'Ensemble des Composants
| Domaine | Fonction | Technologie Cible |
| :--- | :--- | :--- |
| **API Gateway** | Point d'entrée unique, routage, rate limiting | Kong / APISIX |
| **Service Mesh** | Communication Est-Ouest, mTLS, Observabilité | Istio / Linkerd |
| **API Registry** | Catalogue central des APIs gouvernementales | Custom Registry / Backstage |
| **Event Bus** | Messagerie asynchrone pour l'intégration temps réel | Apache Kafka |
| **Federation Layer** | Agrégation des identités et services inter-agences | Keycloak / OIDC |
| **Security Layer** | Modèle Zero Trust, WAF, mTLS, JWT | OPA / Vault |

## 3. Flux de Données
1. **Externe vers Interne**: Trafic passant par le National API Gateway.
2. **Inter-Services (Est-Ouest)**: Trafic géré par le Service Mesh avec chiffrement mTLS systématique.
3. **Asynchrone**: Événements publiés sur le Kafka National pour notification aux agences.

## 4. Principes Directeurs
- **Souveraineté**: Hébergement local et contrôle total des clés de chiffrement.
- **Standardisation**: OpenAPI 3.0, JSON, gRPC.
- **Sécurité par Design**: Authentification forte (OIDC) et autorisation granulaire (RBAC/ABAC).
- **Résilience**: Support du mode déconnecté (queues persistantes).

---
*Document Version: 1.0.0*
*Status: Approved*
