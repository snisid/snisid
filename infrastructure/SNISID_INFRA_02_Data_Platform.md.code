# VOLUME 2 : Plateforme Nationale de Données (Data Platform)
## Infrastructure de Production Souveraine — SNISID

La plateforme de données de l'État doit combiner une consistance stricte (ACID) pour garantir qu'un NNI ne soit jamais dupliqué, avec une résilience aux coupures réseau massives (Théorème CAP).

---

## 🗄️ CHAPITRE 1 : ARCHITECTURE NEWSQL (COCKROACHDB)

L'utilisation de PostgreSQL classique monolithique est inadaptée pour résister à la chute d'un centre de données (DC). Le SNISID utilise **CockroachDB**, une base de données distribuée compatible PostgreSQL, conçue pour la géo-réplication.

### 1.1 Topologie de Survie (Géo-Partitioning)
*   Les nœuds (Nodes) sont répartis sur les 2 Datacenters principaux (Port-au-Prince et Cap-Haïtien) et les 10 clusters régionaux.
*   **Règle de Réplication (Quorum):** La donnée est répliquée 5 fois. Pour valider une transaction, la majorité (3 nœuds) doit valider l'écriture (Consensus Raft). Même si le DC de Port-au-Prince (2 nœuds) s'effondre, les 3 nœuds restants continuent d'opérer sans perte d'une seule milliseconde de données.
*   **Géo-Partitionnement :** Pour des raisons de performance (latence), les données d'identité des citoyens du Nord sont stockées physiquement (Leaseholder) sur les serveurs du Nord, tout en étant sauvegardées dans le Sud.

---

## 🚀 CHAPITRE 2 : EVENT SOURCING (APACHE KAFKA)

La base de données relationnelle représente "l'état actuel". Kafka représente "l'histoire de l'État".

### 2.1 Cluster Kafka National
*   **Broker Topology :** Architecture Strimzi opérant au-dessus de Kubernetes.
*   **Stockage Tiered (Tiered Storage) :** Les événements vieux de plus de 30 jours sont basculés automatiquement des disques NVMe coûteux des brokers Kafka vers le Data Lake S3 moins onéreux, permettant à Kafka de conserver un historique "infini" sans saturer les disques primaires.
*   **MirrorMaker 2 :** Synchronisation asynchrone des topics Kafka entre le Cluster Core et les Clusters Edge régionaux.

---

## 🌊 CHAPITRE 3 : DATA LAKE SOUVERAIN ET ANALYTIQUE

Le Gouvernement haïtien a besoin de capacités de Big Data (Statistiques de naissances, surveillance sanitaire, élections) sans envoyer les données chez Google Cloud (BigQuery) ou AWS (Redshift).

### 3.1 MinIO Enterprise (S3 Compatible)
*   **Data Lake National :** Déploiement d'une architecture MinIO massivement parallèle sur des serveurs JBOD (Just a Bunch of Disks). 
*   **Erasure Coding :** Algorithme mathématique permettant de perdre la moitié des disques durs d'un rack serveur sans perdre la moindre donnée, supprimant le besoin de RAID matériel.

### 3.2 OpenSearch (Moteur de Recherche & Analytique)
*   Utilisé par l'Institut Haïtien de Statistique et d'Informatique (IHSI).
*   Se branche sur le Data Lake S3 pour permettre des requêtes complexes en temps réel sur des milliards d'événements sans impacter la base de données de production CockroachDB.

---

## 🔒 CHAPITRE 4 : STOCKAGE D'AUDIT IMMUABLE (WORM)

La conformité légale (Legal Compliance) du SNISID exige que certains événements (mariage, naturalisation) ne soient jamais modifiables, même par le super-administrateur Root de la base de données.

### 4.1 Implémentation S3 Object Lock (Compliance Mode)
*   Lorsque le moteur de Workflow SNISID publie un certificat finalisé, il l'écrit sur un Bucket MinIO configuré en "Compliance Mode" (Write Once, Read Many).
*   **Protection Ransomware/Insider Threat :** Le protocole S3 refuse au niveau du firmware logiciel toute requête `DELETE` ou `PUT` (écrasement) sur cet objet avant la date de péremption définie (ex: 100 ans).
*   Si le cluster de base de données est détruit par un ransomware ou un acteur malveillant interne, l'État peut reconstruire la vérité absolue depuis ce stockage WORM cryptographiquement sûr.
