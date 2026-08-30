# 🔧 RUNBOOK — Analytics Overload Scaling

**ID** : RB-ANL-005
**Sévérité max** : HIGH
**Propriétaire** : Plateforme Analytics
**SLA résolution** : 30 min

---

## 1. SYMPTÔMES

- Latence Trino/Superset > seuil
- File d'attente requêtes > N
- CPU/RAM workers > 85 %
- Lag Kafka croissant
- Saturation MinIO IOPS

## 2. IMPACT

- Cockpits lents
- Pipelines retardés
- Risque de chute de service en cascade

## 3. DIAGNOSTIC

```bash
# Trino : top requêtes coûteuses
curl https://trino.snisid.ht/v1/query | jq '.[] | select(.state=="RUNNING")' \
  | jq -r '[.queryId, .query[0:80], .state] | @tsv'

# Spark History
curl https://spark-history.snisid.ht/api/v1/applications

# Kubernetes ressources
kubectl top nodes
kubectl top pods -A | sort -k3 -n -r | head -20

# Kafka lag
kafka-consumer-groups --bootstrap-server kafka:9092 --describe --all-groups
```

## 4. PROCÉDURE DE REMÉDIATION

### 4.1 Scaling horizontal

```bash
# Trino workers
helm upgrade trino bitnami/trino -n trino --set worker.replicas=20

# Spark dynamic allocation
# (déjà actif via KEDA, mais boost manuel possible)
kubectl scale deploy spark-thrift-server --replicas=6

# Superset workers
kubectl scale deploy superset-worker -n bi --replicas=12

# Flink parallelism
flink modify <job_id> -p 32
```

### 4.2 Throttling
- Activer `query_max_memory_per_node` Trino
- Limiter concurrent queries par user
- Bloquer requêtes > 10 min en dehors heures off-peak

### 4.3 Killing requêtes runaway
```sql
CALL system.runtime.kill_query('<id>', 'overload mitigation');
```

### 4.4 Cache & matérialisation
- Activer cache résultats Trino
- Précalculer dashboards lourds via dbt incremental
- Pousser KPIs critiques vers Druid (sub-seconde)

### 4.5 Délestage non critique
- Pauser DAGs non critiques (rapports mensuels, analyses ad-hoc)
- Reporter retrainings ML
- Garder uniquement : ingestion, fraude temps réel, cockpit présidentiel

## 5. VÉRIFICATION

- [ ] Latence P95 < SLA
- [ ] CPU < 70 %
- [ ] Lag Kafka décroissant
- [ ] File requêtes vide < 60 s d'attente
- [ ] Cockpits fluides

## 6. POST-MORTEM

Si overload structurel récurrent :
- Revoir capacity planning
- Optimiser requêtes top consommatrices
- Considérer scaling permanent (CAPEX)
