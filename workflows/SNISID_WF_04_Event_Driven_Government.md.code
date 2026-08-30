# VOLUME 4 : Gouvernement Piloté par les Événements (Event-Driven Government)
## Usine Nationale des Workflows — SNISID

L'architecture du SNISID est découplée. Le moteur de workflow orchestre les sagas et publie les résultats finaux de manière immuable sur un Cluster Apache Kafka Hautement Disponible, servant de source unique de vérité pour tout l'État haïtien.

---

## 🚀 CHAPITRE 1 : ARCHITECTURE DE DISTRIBUTION KAFKA

L'architecture de distribution garantit qu'aucune donnée n'est perdue lors de coupures des faisceaux hertziens ou fibres optiques provinciales.

```mermaid
graph TD
    WF[Workflow Engine (Temporal)] -->|Publish Event| SR[Schema Registry (Protobuf)]
    SR -->|Valid Schema| K[Kafka Central (Port-au-Prince)]
    K -->|MirrorMaker 2| DR[Kafka DR (Cap-Haïtien)]
    
    K --> DGI[Consumer: DGI]
    K --> CEP[Consumer: Conseil Electoral]
    K --> MSPP[Consumer: Santé Publique]
    K --> WORM[Consumer: Audit Archiver]
```

---

## 📂 CHAPITRE 2 : REGISTRE NATIONAL DES TOPICS (TOPIC REGISTRY)

Le nommage des topics respecte une convention stricte : `tg.<domaine>.<sous-domaine>.<type-événement>` (tg = Transition Gouvernementale).

### Domaines État Civil (Civil Registration)
*   `tg.civil.birth.events` : Naissances enregistrées. (Partition key : NNI Enfant).
*   `tg.civil.death.events` : Décès enregistrés. Déclenche la suspension des droits. (Partition key : NNI).
*   `tg.civil.marriage.events` : Mariages.
*   `tg.civil.divorce.events` : Divorces.

### Domaines Identité (Identity Security)
*   `tg.identity.enrollment.success` : Nouvelle identité certifiée après hit ABIS négatif.
*   `tg.identity.status.updated` : Suspension, révocation ou gel d'une identité.
*   `tg.security.fraud.alerts` : Alertes ABIS, détection de fraude, vélocité anormale.

---

## 📝 CHAPITRE 3 : SCHÉMA D'ÉVÉNEMENT GOUVERNEMENTAL (PROTOBUF)

Tous les événements publiés sont fortement typés avec Protobuf, garantissant la rétrocompatibilité via le Schema Registry.

**Exemple de Schéma de Décès (Death Event) :**
```protobuf
syntax = "proto3";
package tg.civil.death;

message DeathEvent {
  string event_id = 1;               // UUID WORM
  string timestamp_iso8601 = 2;      // Heure exacte d'enregistrement
  string nni = 3;                    // NNI du défunt
  string death_certificate_id = 4;   // Référence MSPP
  string medical_officer_nni = 5;    // Identité du médecin constatant
  string signature_hash = 6;         // Signature cryptographique PKI du médecin
  
  enum DeathType {
    STANDARD = 0;
    JUDICIAL = 1;
    DISASTER = 2;
  }
  DeathType type = 7;
}
```

---

## 🔄 CHAPITRE 4 : RÉSILIENCE ET STRATÉGIES DE REJEU (REPLAY)

### 4.1 Idempotence Absolue
Étant donné la nature du réseau en Haïti, les consommateurs d'événements (agences gouvernementales) peuvent recevoir le même événement plusieurs fois. Tous les consommateurs doivent utiliser l'`event_id` comme clé de déduplication (Idempotency Key) dans une base de données ACID (ex: PostgreSQL / CockroachDB).

### 4.2 Pattern de Rejeu (Replay Strategy)
En cas de désastre majeur (ex: serveur de la DGI détruit par un incendie), l'État peut rejouer tous les événements Kafka à partir de l'offset "Zéro" ou depuis une date donnée, reconstruisant ainsi l'intégralité de la base de données citoyenne de l'agence sinistrée en quelques heures, à partir de la dorsale souveraine SNISID.

### 4.3 Outbox Pattern
Le moteur de workflow n'écrit jamais directement dans Kafka au milieu d'une transaction de base de données. Il utilise le **Transactional Outbox Pattern** : 
1. Il écrit le changement de statut du citoyen ET l'événement à publier dans la MÊME transaction de base de données.
2. Un processus CDC (Change Data Capture) comme Debezium lit la table *Outbox* et publie de manière asynchrone et fiable l'événement vers Kafka. Cela garantit une cohérence absolue (Zero Data Loss).
