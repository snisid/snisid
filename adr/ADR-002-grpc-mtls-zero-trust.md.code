# ADR-002: mTLS Obligatoire pour toutes les Communications gRPC

**Statut:** Accepté
**Date:** 2026-06-18
**Décideurs:** Security Owner, Lead Architecte

## Contexte
Les communications gRPC utilisent actuellement `insecure.NewCredentials()` sans TLS, exposant les données biométriques et les décisions de scoring à des attaques Man-in-the-Middle.

## Décision
Imposer mTLS (Mutual TLS) pour toutes les communications gRPC via:
- Certificats SPIFFE pour l'identité des workloads (SPIRE)
- Istio mTLS (PeerAuthentication STRICT) pour la couche service mesh
- Certificats signés par Vault PKI pour les connexions directes
- Rejet de toute connexion sans certificat valide

## Conséquences
Positives:
- Chiffrement de bout en bout de toutes les communications
- Identité forte de chaque workload (SPIFFE ID)
- Protection contre les attaques MITM et replay

Négatives:
- Overhead de performance (négociation TLS)
- Complexité de gestion des certificats
- Délai accru pour les connexions initiales

## Alternatives considérées
1. TLS unidirectionnel seul: Rejeté (pas d'authentification mutuelle)
2. JWT-only pour l'authentification: Rejeté (pas de chiffrement du transport)
