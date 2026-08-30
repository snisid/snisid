---
# ============================================================
# SNISID-Security — National Security Search & Correlation Engine
# Moteur de Recherche Inter-Agences et Alertes
# Document ID: SNISID-SEC-CORRELATION-001
# Version: 1.0.0
# ============================================================

## 1. CONCEPT : LE MOTEUR DE CORRÉLATION NATIONAL

Alors que le système DCPJ est orienté "Graphes et Réseaux", le **Correlation Engine** est le système nerveux central en **temps réel**. Il repose sur OpenSearch (Elasticsearch souverain) et Kafka Streams pour croiser les données structurées et non structurées de toutes les agences.

Il permet une recherche globale type "Google-like" sur des milliards de documents, filtrée par les habilitations de l'agent effectuant la recherche (ABAC).

## 2. SOURCES DE DONNÉES INDEXÉES (OpenSearch)

| Index | Source | Contenu | Sensibilité |
|-------|--------|---------|-------------|
| `idx_identities` | Identity Registry | NIU, Nom, Photo Hash, Statut (Vivant/Décédé) | Medium |
| `idx_civil_acts` | Civil Registry | Actes de naissance, mariage, divorce | High |
| `idx_criminal_cases`| Justice System | Résumé des affaires, décisions | High |
| `idx_warrants` | Justice System | Mandats actifs, annulés, exécutés | Medium (Actifs) |
| `idx_police_incidents`| PNH | Mains courantes, rapports terrain | High |
| `idx_border_events`| Immigration | Entrées/Sorties territoire, Refus | Medium |
| `idx_prison_stays` | Pénitentiaire | Séjours carcéraux, incidents prison | High |

## 3. STREAMING CORRELATION (Kafka Streams)

Des topologies Kafka Streams évaluent chaque événement transitant sur le bus en temps réel pour détecter des corrélations critiques.

```java
// Exemple conceptuel de topologie Kafka Streams (Correlation Engine)
KStream<String, BorderEvent> borderEvents = builder.stream("snisid.border.events");
KTable<String, Warrant> activeWarrants = builder.table("snisid.justice.warrants.active");

// Jointure Temps Réel: Détection d'un fugitif à la frontière
KStream<String, Alert> fugitiveAlerts = borderEvents
    .join(activeWarrants,
        (borderEvent, warrant) -> {
            if (borderEvent.getNiu().equals(warrant.getTargetNiu())) {
                return new Alert(
                    "FUGITIVE_AT_BORDER",
                    "Mandat " + warrant.getId() + " détecté à " + borderEvent.getPortId(),
                    Severity.CRITICAL
                );
            }
            return null;
        }
    )
    .filter((key, alert) -> alert != null);

fugitiveAlerts.to("snisid.security.alerts.critical");
```

## 4. ARCHITECTURE DE L'ALERTING (ESCALATION WORKFLOW)

Les alertes critiques générées par le Correlation Engine sont routées vers un système de gestion d'incidents (similaire à PagerDuty, mais souverain).

```mermaid
flowchart TD
    A[Kafka Streams Correlation] -->|Alerte Générée| B(Alert Router)
    
    B -->|Fugitif Frontière| C[DGIE / POLIFRONT / DCPJ]
    B -->|Évasion Prison| D[PNH Locale / Unités Spéciales]
    B -->|Fraude Identité| E[ONI / Parquet]
    
    C --> F{Acquittement?}
    F -->|Non (5 min)| G[Escalade Niveau 2 - Directeur]
    F -->|Oui| H[Intervention en cours]
```

## 5. RECHERCHE MULTI-AGENCE & SÉCURITÉ

La recherche est unifiée via un point d'entrée API `/v1/search/global`.

**Le filtre OPA (Open Policy Agent) injecte silencieusement des clauses dans la requête OpenSearch :**
- Un agent PNH cherche "Jean M". La requête est modifiée pour exclure `idx_civil_acts` et ne retourner que les `idx_identities` (infos publiques) et `idx_police_incidents` (de sa juridiction).
- Un juge d'instruction cherche "Jean M". Il voit `idx_civil_acts` et `idx_criminal_cases` liés à sa juridiction.

---
*Document ID: SNISID-SEC-CORRELATION-001 | Approuvé par: Architecte Souverain*
