# 🔧 RUNBOOK — Pipeline Failure Recovery

**ID** : RB-ANL-001
**Sévérité max** : HIGH
**Propriétaire** : Data Engineering
**SLA résolution** : 2h (HIGH) / 30min (CRITICAL)

---

## 1. SYMPTÔMES

- Alerte Prometheus `PipelineFailureHigh`
- DAG Airflow en état `failed`
- Job Spark/Flink en boucle de redémarrage
- Lag Kafka > seuil
- Datasets Gold non rafraîchis (alerte fraîcheur)

## 2. IMPACT

- Données décisionnelles obsolètes
- Cockpits affichant valeurs périmées
- Risque sur prédictions / scoring
- Possible cascade vers BI / ML serving

## 3. DIAGNOSTIC

```bash
# 1. Identifier le pipeline en échec
kubectl get pods -n airflow | grep -v Running
airflow tasks list <dag_id> --tree

# 2. Examiner les logs
kubectl logs -n airflow <task_pod> --tail=500
# ou Loki :
logcli query '{service="airflow", dag_id="<dag>"} |= "ERROR"'

# 3. Vérifier dépendances upstream
#   - Kafka (lag, brokers)
kafka-consumer-groups --bootstrap-server kafka:9092 \
  --describe --group <consumer_group>
#   - MinIO (santé bucket)
mc admin info snisid-minio
#   - Postgres source (CDC slot)
psql -c "SELECT * FROM pg_replication_slots;"

# 4. Vérifier ressources
kubectl top pods -n analytics
```

## 4. PROCÉDURE DE REMÉDIATION

### 4.1 Cas A — Échec transitoire (réseau, ressource)
```bash
airflow dags trigger <dag_id> --conf '{"backfill": true, "start": "<ts>"}'
```

### 4.2 Cas B — Erreur code / schéma
1. Identifier commit fautif (Git blame DAG)
2. Revert via PR d'urgence
3. Redéployer image Airflow / job
4. Backfill période impactée

### 4.3 Cas C — Source upstream KO
1. Notifier équipe source (CDC, Kafka, API)
2. Activer mode **degraded** : marquer dataset stale dans catalogue
3. Bannière sur dashboards : "Données figées au {ts}"
4. Reprise dès rétablissement source

### 4.4 Cas D — Saturation ressources
```bash
# Augmenter executors Spark
spark-submit --conf spark.executor.instances=20 ...

# Scaler Flink
flink modify <job_id> -p 16

# Scale Airflow workers
kubectl scale deployment airflow-worker --replicas=10 -n airflow
```

## 5. VÉRIFICATION

- [ ] DAG/job en état `success`
- [ ] Lag Kafka redescendu < seuil
- [ ] Dataset cible mis à jour (fraîcheur OK)
- [ ] DQ score ≥ 90
- [ ] Dashboards downstream actualisés
- [ ] Alertes éteintes dans Alertmanager

## 6. POST-MORTEM

Si pipeline critique (NRIC, présidentiel) :
- Post-mortem sous 72h
- Document dans `runbooks/postmortems/`
- Actions correctives → backlog
