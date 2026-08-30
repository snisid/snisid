# SNISID RUNBOOK — Kafka Recovery (National Messaging)
**Classification:** RESTREINT DEFENSE  
**Version:** 4.0.0  
**RTO cible:** < 30 minutes  
**Fréquence test:** Trimestriel

---

## 1. Architecture Kafka Nationale

- **Cluster:** 6 brokers (3 Core + 3 DR), replication factor 3 minimum
- **MirrorMaker2:** Core ↔ DR (active-active pour topics critiques)
- **Topics critiques:** `identites.nationales`, `enrollments.terrain`, `audit.securite`, `sync.edge`
- **Alerte déclencheur:** `SNISID_Kafka_OfflinePartition`, `SNISID_Kafka_BrokerDown`

## 2. Scénarios

| Code | Scénario | Symptômes |
|------|----------|-----------|
| KFK-001 | Broker unique down | 1/6 brokers NotReady, partitions under-replicated |
| KFK-002 | Perte quorum ZooKeeper/KRaft | Pas d'élection controller, topics inaccessibles |
| KFK-003 | Corruption log segments | Broker crash-loop, `CorruptRecordException` |
| KFK-004 | Consumer lag critique national | `sync.edge` lag > 100K, retard terrain |

## 3. Procédure KFK-001 : Remplacement broker

```bash
BROKER="kafka-core-02"
PARTITION="identites.nationales"

# 1. Vérifier under-replication
kafka-topics.sh --bootstrap-server kafka-core-01:9093 --describe --topic ${PARTITION} | grep "UnderReplicated"

# 2. Retirer le broker
kafka-server-stop.sh
kubectl delete pod ${BROKER} -n snisid-data

# 3. Vérifier PVC (Ceph RBD) — si corrompu, restaurer snapshot
kubectl get pvc -n snisid-data | grep ${BROKER}
# Snapshot Ceph : rbd snap ls snisid-rbd-tier1/kafka-core-02

# 4. Redémarrer via StatefulSet
kubectl rollout restart statefulset/kafka-core -n snisid-data

# 5. Validation
kafka-topics.sh --bootstrap-server kafka-core-01:9093 --describe --topic ${PARTITION}
kafka-consumer-groups.sh --bootstrap-server kafka-core-01:9093 --describe --group edge-sync-consumer
```

## 4. Procédure KFK-002 : Recovery controller quorum (KRaft)

> Kafka 3.x utilise KRaft (pas ZooKeeper). Le quorum est `node.id` ensemble.

```bash
# 1. Identifier le controller actif
kafka-metadata-quorum.sh --bootstrap-server kafka-core-01:9093 --describe

# 2. Si quorum perdu (moins de (N/2)+1 voters actifs)
# Forcer un nouveau quorum depuis le broker avec le log le plus récent
# Ceci est une opération d'urgence — notifier IGC préalablement

# Sur le broker survivant avec le plus haut offset:
kafka-storage.sh format -t $(cat /var/lib/kafka/cluster.id) -c /opt/kafka/config/kraft/server.properties --ignore-formatted

# 3. Reset les autres brokers comme followers
# Editer /var/lib/kafka/meta.properties pour retirer l'état quorum corrompu
# Démarrer les brokers un par un

# 4. Validation quorum
kafka-metadata-quorum.sh --bootstrap-server kafka-core-01:9093 --describe
```

## 5. Procédure KFK-004 : Lag critique national

```bash
# 1. Identifier le lag
kafka-consumer-groups.sh --bootstrap-server kafka-core-01:9093 --describe --group edge-sync-consumer

# 2. Si cause : partition skew
# Rééquilibrer les partitions (si clé mal distribuée)
kafka-reassign-partitions.sh --bootstrap-server kafka-core-01:9093 \
  --topics-to-move-json-file topics-to-move.json \
  --broker-list "1,2,3,4,5,6" \
  --generate

# 3. Si cause : consumer down en edge
# Déclencher alerte tier-4, vérifier connectivité edge
kubectl logs -n snisid-edge deployment/edge-kafka-consumer

# 4. Mesures d'urgence : augmenter consumers par partition
kubectl scale deployment edge-kafka-consumer --replicas=10 -n snisid-edge

# 5. Si edge offline > 72h : activer mobile sync via USB burst
```

---

*Test trimestriel sur cluster staging-dr. Logs obligatoires SIEM.*
