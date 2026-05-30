---
# ============================================================
# SNISID API Ecosystem — API Audit Fabric
# Traçabilité Légale et Non-Répudiation B2B
# Document ID: SNISID-API-AUDIT-001
# Version: 1.0.0
# ============================================================

## 1. L'IMPÉRATIF DE LA PREUVE LÉGALE

Lorsqu'une entité privée interroge une base de l'État, l'État doit pouvoir prouver devant un tribunal *qui* a demandé l'information, *quand*, et *quelle* réponse a été donnée. C'est la Non-Répudiation.

## 2. ARCHITECTURE DE LA TRAÇABILITÉ (Audit Fabric)

- **Interception Asynchrone :** Kong Gateway utilise un plugin asynchrone (ex: `http-log` ou `kafka-log`) pour copier les métadonnées de *chaque* requête B2B vers un topic Kafka dédié `b2b-audit-logs`.
- **Ce qui est loggé :** 
  - ID de la transaction.
  - Certificat client utilisé.
  - Endpoint appelé (ex: `/b2b/v1/identity/verify`).
  - Statut de la réponse (HTTP 200, Match: True).
  - Horodatage cryptographique.
- **Ce qui n'est JAMAIS loggé (Data Privacy) :** Les empreintes biométriques (WSQ/IrisCode) ou les mots de passe.

## 3. IMMUABILITÉ ET CONSERVATION (WORM)

Les messages Kafka sont ingérés par OpenSearch (Phase 6) pour analyse en temps réel (SOC), mais ils sont également archivés dans le système WORM (Write Once, Read Many) du Datacenter (Phase 5). 
Même un administrateur système "root" de l'État ne peut pas altérer les logs B2B sur les disques WORM, garantissant ainsi l'intégrité de la preuve légale.

---
*Document ID: SNISID-API-AUDIT-001 | Approuvé par: Inspection Générale de l'État*
