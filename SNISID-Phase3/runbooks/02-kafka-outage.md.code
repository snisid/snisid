# 🚨 Runbook 02 — Panne Kafka / Event Mesh

**Severity :** Sev1 (Sev0 si > 30 min)
**Owner :** Astreinte Platform Engineering

## 1. Symptômes
- `kafka_cluster_active_controllers != 1` (controller perdu)
- `kafka_under_replicated_partitions > 0`
- Producers : exceptions `NotLeaderForPartitionException`, `OutOfOrderSequenceException`
- Consumers : lag explose (`kafka_consumergroup_lag > 100000`)
- Workflows en attente d'émission Kafka → cascade de timeouts

## 2. Diagnostic

```bash
# État du cluster
kubectl -n kafka exec kafka-1-0 -- kafka-topics.sh --bootstrap-server localhost:9093 \
  --command-config /etc/kafka/client.properties --describe --under-replicated-partitions

# Controller
kubectl -n kafka exec kafka-1-0 -- kafka-metadata-shell.sh --snapshot \
  /var/kafka-logs/__cluster_metadata-0/00000000000000000000.log

# Broker logs
kubectl -n kafka logs kafka-X-0 --tail=200
```

## 3. Remédiation

### Cas A — Un broker HS
1. Vérifier le pod (`kubectl describe pod`) : disque ? OOM ? réseau ?
2. Relancer si nécessaire : `kubectl rollout restart sts/kafka -n kafka`
3. Attendre re-sync des partitions (peut prendre 10-30 min selon volume).

### Cas B — Plus d'un broker HS
1. Activer le **mode dégradé** côté producers :
   ```bash
   kubectl set env deploy/snisid-workflow-engine PRODUCER_DEGRADED=true
   ```
   → activations switchent vers **outbox PostgreSQL local** (workflows non bloqués).
2. Restaurer le quorum brokers (priorité absolue).
3. Une fois cluster sain, **flusher l'outbox** :
   ```bash
   ./scripts/flush-outbox-to-kafka.sh
   ```

### Cas C — Cluster KO total
1. Déclencher le **bascule DR DC3** :
   ```bash
   kubectl apply -f deploy/dr/kafka-failover-dc3.yaml
   ```
2. Mettre à jour `KAFKA_BROKERS` dans tous les services (via External Secret).
3. Lancer le workflow `escalation.crisis.national`.

## 4. Vérification
- `kafka_under_replicated_partitions == 0`
- Lag consumer < 1 000
- Workflows critiques reprennent
- Outbox vide

## 5. Communication
- WGO + Direction ONI : immédiat
- Citoyens : message "service ralenti, fonctionnement nominal en cours de rétablissement"

## 6. Post-mortem
Obligatoire si > 15 min OU > 1 broker affecté.
