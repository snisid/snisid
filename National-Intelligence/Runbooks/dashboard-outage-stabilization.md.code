# 🔧 RUNBOOK — Dashboard Outage Stabilization

**ID** : RB-ANL-002
**Sévérité max** : CRITICAL
**Propriétaire** : BI Platform Team
**SLA résolution** : 15 min (présidentiel) / 1h (autres)

---

## 1. SYMPTÔMES

- Probe `DashboardDown` rouge
- Utilisateurs signalent 5xx ou écrans blancs
- Latence Superset/Grafana > 10s
- Pods Superset/Grafana en `CrashLoopBackOff`

## 2. IMPACT

- **CRITIQUE** si cockpit présidentiel ou crise
- Aveuglement décisionnel
- Perte confiance dans plateforme

## 3. DIAGNOSTIC

```bash
# Statut services
kubectl get pods -n bi
kubectl describe pod <pod> -n bi

# Health endpoints
curl -k https://superset.snisid.ht/health
curl -k https://grafana.snisid.ht/api/health

# Backend BD Superset
kubectl exec -n bi deploy/superset -- \
  python -c "from superset import db; print(db.engine.execute('SELECT 1').scalar())"

# Trino health (datasource)
curl https://trino.snisid.ht/v1/info

# Charge requête
# Loki:
{service="superset"} |= "slow query" | json
```

## 4. PROCÉDURE DE REMÉDIATION

### 4.1 Pod crash
```bash
kubectl rollout restart deploy/superset -n bi
kubectl rollout status deploy/superset -n bi
```

### 4.2 Surcharge backend
```bash
# Scaler Superset workers
kubectl scale deploy/superset-worker --replicas=8 -n bi
# Scaler Trino coordinators / workers
helm upgrade trino bitnami/trino --set worker.replicas=12 -n trino
```

### 4.3 Requête lourde bloquante
```sql
-- Trino: kill query
CALL system.runtime.kill_query(query_id => '<id>', message => 'ops kill');
```

### 4.4 Mode dégradé cockpit présidentiel
1. Activer **fallback statique** (cache HTML signé pré-généré)
2. Redirection DNS vers `cockpit-fallback.snisid.ht`
3. Notification décideurs : "Mode dégradé, valeurs T-1h"
4. Investiguer en parallèle

### 4.5 Cache & CDN interne
```bash
# Vider cache front
kubectl exec deploy/superset-redis -- redis-cli FLUSHDB
```

## 5. VÉRIFICATION

- [ ] `probe_success{job="superset_probe"} == 1` depuis 5 min
- [ ] Latence P95 < 2 s
- [ ] Aucun pod en erreur
- [ ] Cockpit présidentiel charge < 2 s
- [ ] Échantillon utilisateurs confirme

## 6. POST-MORTEM

Obligatoire si cockpit présidentiel down > 5 min.
Communication : note interne direction + ministères impactés.
