---
# ============================================================
# SNISID-Interop — Interoperability Bus
# Backbone Event-Driven (Kafka)
# Document ID: SNISID-INTEROP-BUS-001
# Version: 1.0.0
# ============================================================

## 1. BACKBONE ASYNCHRONE DE L'ÉTAT

L'architecture gouvernementale repose sur le paradigme **Publish/Subscribe**. Les agences ne doivent pas faire des requêtes "GET" constantes pour savoir si une donnée a changé (Polling).
Elles "s'abonnent" aux événements pertinents.

## 2. EVENT TOPOLOGY

### 2.1 Topics Gouvernementaux Standardisés

| Topic Kafka | Producteur Autorisé | Consommateurs Types | Qualité de Service (QoS) |
|-------------|---------------------|---------------------|--------------------------|
| `snisid.identity.lifecycle` | SNISID Identity | DGI, MSPP, Justice | Exactly-Once, Infinite Retention |
| `snisid.civil.vital` | Civil Registry | DGI (Héritage), CEP (Électeurs) | Exactly-Once, Infinite |
| `snisid.justice.warrants` | Justice | PNH, Frontières, DCPJ | At-Least-Once, Compacted |
| `snisid.health.births` | MSPP (Hôpitaux) | Civil Registry | Exactly-Once, 30 days |

### 2.2 Avro Schema Registry

Un topic Kafka ne contient pas de JSON libre, mais du binaire validé par un **Schema Registry** (Apicurio).
Si le Ministère de la Santé (MSPP) tente de publier un événement de Naissance sans le champ obligatoire `mother_niu`, le message est rejeté avant même d'entrer dans le bus.

```json
// Exemple Schema Avro: DeathRegisteredEvent
{
  "type": "record",
  "name": "DeathRegistered",
  "namespace": "ht.gov.snisid.events",
  "fields": [
    {"name": "niu", "type": "string"},
    {"name": "date_of_death", "type": "string", "logicalType": "timestamp-millis"},
    {"name": "certifying_doctor_id", "type": "string"},
    {"name": "cause_code", "type": ["null", "string"]}
  ]
}
```

### 2.3 Dead-Letter Queues (DLQ)
Si le système de la DGI n'arrive pas à traiter un événement `DeathRegistered` après 5 tentatives, le message est mis dans une DLQ (`snisid.dgi.dlq`). Des opérateurs humains peuvent analyser pourquoi l'intégration a échoué et rejouer le message manuellement (Replayable Events).

---
*Document ID: SNISID-INTEROP-BUS-001 | Approuvé par: Architecte Souverain*
