# 🟧 Runbook 08 — Cluster Temporal HS

**Severity :** Sev1
**Owner :** Astreinte Platform Engineering

## 1. Symptômes
- Workflows long-running ne progressent plus
- Workers ne polling pas (`temporal_worker_poll_succeed_per_rps == 0`)
- Frontend HTTP/gRPC 503
- Logs : `unable to acquire shard`, `Cassandra unavailable`, `history shard not found`

## 2. Diagnostic

```bash
# Pods Temporal
kubectl -n temporal get pods

# Health
kubectl -n temporal exec temporal-frontend-0 -- tctl --ad localhost:7233 cluster health

# Cassandra (backend par défaut)
kubectl -n cassandra exec cassandra-0 -- nodetool status
kubectl -n cassandra exec cassandra-0 -- nodetool tpstats | head -20

# Métriques
curl -s http://prometheus:9090/api/v1/query?query='temporal_persistence_latency_bucket'
```

Causes principales :
- **Cassandra down / surchargé** (storage backend par défaut)
- **Frontend overload** (limites RPS)
- **History service** : shard ownership stuck
- **Matching service** : task queue backlog

## 3. Remédiation

### Cas A — Cassandra surchargé / un node down
```bash
# Vérifier health Cassandra
kubectl -n cassandra exec cassandra-0 -- nodetool status | grep -E '^UN|^DN'

# Redémarrer le node down
kubectl -n cassandra delete pod cassandra-X

# Forcer compaction (si tombstones nombreux)
kubectl -n cassandra exec cassandra-0 -- nodetool compact temporal
```

### Cas B — Service Temporal HS
```bash
# Rolling restart par service (NE PAS tout en même temps)
kubectl -n temporal rollout restart deploy/temporal-frontend
sleep 60
kubectl -n temporal rollout restart deploy/temporal-history
sleep 60
kubectl -n temporal rollout restart deploy/temporal-matching
sleep 60
kubectl -n temporal rollout restart deploy/temporal-worker
```

### Cas C — Workflows bloqués
```bash
# Lister les workflows en cours
tctl --ns snisid-prod workflow list --query 'ExecutionStatus="Running"' --more

# Forcer reset d'un workflow bloqué
tctl --ns snisid-prod workflow reset \
  -w <workflowId> -r <runId> --event-id <eventId> --reason "stuck after outage"

# Terminer si impossible à réparer (sera reprise par déclencheur Kafka)
tctl --ns snisid-prod workflow terminate -w <workflowId> --reason "..."
```

### Cas D — Cluster perdu — bascule DC2
1. Mettre à jour `TEMPORAL_ADDRESS` dans tous les workers via External Secret.
2. Workers Temporal repollent automatiquement.
3. Cassandra DC2 doit être en `Active-Active` (configurée Phase 1).

## 4. Vérification
- `tctl cluster health` → `SERVING`
- `temporal_persistence_latency_bucket p99` < 50ms
- Workers actifs : `temporal_worker_active_count > 0`
- Workflows progressent : nouveau test
  ```bash
  tctl workflow start --tq snisid-default --wt birthSimpleWorkflow --input '{"test":true}'
  ```

## 5. Communication
- WGO immédiat
- Workflows Temporal = orchestrations long-running (mois) ; donc dégradation possiblement invisible immédiat mais critique sur 24-48h.

## 6. Post-mortem
- Capacity Cassandra (heap, compaction, IOPS)
- Limites RPS frontend
- Plan : passer en backend PostgreSQL si plus simple à opérer
