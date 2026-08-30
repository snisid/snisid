# SNISID National Data Architecture

## 1. But
Définir l'architecture officielle des données nationales SNISID afin de fournir une plateforme unifiée, souveraine, gouvernée, sécurisée, observable et exploitable en temps réel.

## 2. Vue d'ensemble

```text
[Operational Databases] ---> [CDC/Event Streams] ---> [Sovereign Data Lakehouse]
        |                           |                         |
        |                           v                         v
        |                    [Audit Data Fabric]       [Analytics Platform]
        |                           |                         |
        v                           v                         v
[MDM Golden Records] <------ [Metadata/Lineage] ------ [AI/ML Platform]
        |                                                   |
        v                                                   v
[Data Access Governance] --------------------------> [APIs/Dashboards]
```

## 3. Domaines architecturaux

| Domaine | Fonction | Technologies recommandées | Gouvernance obligatoire |
|---|---|---|---|
| Operational Databases | Transactions officielles | PostgreSQL, MariaDB, Cassandra selon besoin | Classification, ownership, audit |
| Data Lakehouse | Historique, analytique, conservation | MinIO/S3, Iceberg, Trino, Spark | Catalogage, lineage, qualité |
| Event Streams | Temps réel et intégration | Kafka, Schema Registry, Kafka Connect | Contrats d'événements, rétention, audit |
| Metadata Platform | Gouvernance et catalogue | OpenMetadata, DataHub ou Apache Atlas | Propriétaire, sensibilité, contrats |
| Analytics Platform | Intelligence décisionnelle | Trino, Superset/Metabase/Grafana | Accès contrôlé, masquage, audit |
| AI Platform | IA souveraine et contrôlée | MLflow, Feature Store, Model Registry | Explicabilité, validation, audit |
| Audit Data Fabric | Traçabilité nationale | Kafka, OpenSearch, Iceberg, immutability | WORM, horodatage, intégrité |

## 4. Zones du Lakehouse

| Zone | Description | Accès | Règles |
|---|---|---|---|
| Landing | Données brutes reçues | Très restreint | Immutable, chiffrée, auditée |
| Bronze | Données brutes normalisées | Ingénierie data | Schéma minimal, metadata obligatoire |
| Silver | Données nettoyées et validées | Domain teams | Qualité, lineage, contrôles |
| Gold | Produits analytiques officiels | Décideurs/API/BI | Certifié, propriétaire officiel |
| Restricted | Données sensibles/secret national | Accès autorisé uniquement | ABAC, masquage, justification |
| Archive | Conservation longue durée | Archivistes autorisés | Legal hold, WORM, rétention |

## 5. Classification nationale des données

| Niveau | Description | Exemples | Exigences |
|---|---|---|---|
| PUBLIC | Données publiables | Statistiques agrégées | Validation publication |
| INTERNAL | Usage gouvernemental | Rapports internes | RBAC, journalisation |
| CONFIDENTIAL | Données sensibles | Dossiers administratifs | Chiffrement, ABAC, masquage |
| RESTRICTED | Données critiques personnelles | Identité, biométrie, registre civil | MFA, consentement/purpose, audit renforcé |
| NATIONAL_SECRET | Sécurité nationale | Investigations, menaces | Segmentation forte, WORM, accès exceptionnel |

## 6. Règles de gouvernance universelles

- Aucun dataset sans `data_owner`.
- Aucun dataset sans `classification`.
- Aucun dataset sans politique de rétention.
- Aucun pipeline sans lineage technique.
- Aucun accès sans justification et journalisation.
- Aucun modèle IA sans fiche modèle, jeux de données approuvés et audit.
- Aucun partage inter-agence sans contrat de données.

## 7. Flux de données standard

1. Source opérationnelle produit transaction ou événement.
2. CDC/Kafka capture l'événement avec schéma versionné.
3. Donnée arrive en zone Landing/Bronze chiffrée.
4. Contrôles qualité automatiques.
5. Enrichissement MDM et validation inter-agences.
6. Publication Silver/Gold avec certification.
7. Catalogue metadata mis à jour.
8. Lineage technique et métier enregistré.
9. Accès via politiques RBAC/ABAC et purpose limitation.
10. Tous les accès, transformations et exports sont audités.

## 8. RACI synthétique

| Activité | Data Owner | Data Steward | Platform Team | Security | Audit |
|---|---|---|---|---|---|
| Classification | A | R | C | C | I |
| Qualité | A | R | C | I | I |
| Accès | A | C | R | A | I |
| Rétention | A | R | C | C | C |
| Lineage | C | R | R | I | I |
| Audit | I | I | C | C | A |

R=Responsible, A=Accountable, C=Consulted, I=Informed.
