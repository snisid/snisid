# 🟧 Runbook 07 — Cluster Zeebe HS

**Severity :** Sev1 (Sev0 si total + > 15 min)
**Owner :** Astreinte Platform Engineering + WGO

## 1. Symptômes
- `up{job="zeebe-gateway"} == 0` ou pods en `CrashLoopBackOff`
- API gateway répond `UNAVAILABLE` / `DEADLINE_EXCEEDED`
- Aucun nouveau process instance ne démarre
- Jobs en attente explosent : `zeebe_jobs_activated_total` plat
- Logs broker : `partition not healthy`, `Raft leader election failed`

## 2. Diagnostic

```bash
# État du cluster
kubectl -n zeebe get pods -o wide
zbctl status --address zeebe-gateway:26500 --insecure=false

# Health des brokers
kubectl -n zeebe exec zeebe-broker-0 -- curl -s localhost:9600/actuator/health | jq

# Partition leaders
kubectl -n zeebe exec zeebe-broker-0 -- curl -s localhost:9600/actuator/partitions | jq

# Disque (zeebe = TRÈS sensible au disque plein)
kubectl -n zeebe exec zeebe-broker-0 -- df -h /usr/local/zeebe/data
```

Causes courantes :
- **Disque plein** (cause #1 chez Zeebe → blocages quorum Raft)
- **Quorum Raft perdu** (> 1 broker indisponible sur 3)
- **OOM** sur broker (allouer 4-8 Go heap min en prod)
- **Network partition** DC1↔DC2

## 3. Remédiation

### Cas A — Disque plein
```bash
# 1. Augmenter le PVC à chaud (StorageClass avec allowVolumeExpansion=true)
kubectl -n zeebe patch pvc data-zeebe-broker-0 -p '{"spec":{"resources":{"requests":{"storage":"500Gi"}}}}'
kubectl -n zeebe delete pod zeebe-broker-0  # remount

# 2. Forcer la compaction des journaux
kubectl -n zeebe exec zeebe-broker-0 -- curl -X POST localhost:9600/actuator/compact

# 3. Vérifier la rétention configurée
# zeebe.broker.data.snapshotPeriod (par défaut 15min) — réduire si besoin
```

### Cas B — Perte de quorum
1. Identifier les brokers manquants : `kubectl -n zeebe get pods`
2. Si réseau OK mais pods KO → redémarrer un par un avec **rolling restart** :
   ```bash
   kubectl -n zeebe rollout restart sts/zeebe-broker
   ```
3. Surveiller : `zbctl status` doit retrouver `LEADER` pour chaque partition.
4. **NE PAS** restart tous les brokers en même temps → perte définitive de données.

### Cas C — Cluster totalement KO (perte données critique)
1. **Bascule DR DC3** (cf. runbook 09).
2. Restaurer depuis le dernier snapshot S3 (`s3://snisid-zeebe-snapshots/`) :
   ```bash
   kubectl -n zeebe-dr apply -f deploy/zeebe-restore-from-snapshot.yaml
   ```
3. Lecture seule pendant la restauration (mode dégradé : workflows en outbox PG).

### Cas D — Workflows actifs perdus
1. Replay des événements `audit.workflow.transition.v1` depuis Kafka :
   ```bash
   ./scripts/replay-workflows-from-audit.sh --from "2026-05-24T14:00:00Z"
   ```
2. Reprise sur Temporal des workflows long-running (Temporal a sa propre persistance).

## 4. Vérification
- `zbctl status` → tous les brokers `HEALTHY`, leaders distribués
- Nouveaux process créables : test rapide
  ```bash
  zbctl create instance civil-registry.birth.simple --variables '{"test":true}'
  ```
- Backlog jobs résorbé : `kafka-consumer-groups --describe --group zeebe-jobs` lag → 0
- Pas d'incidents en attente : `zbctl list incidents | jq length` = 0

## 5. Communication
| Audience | Délai |
|----------|-------|
| WGO + Direction | immédiat |
| Métier (toutes administrations) | 15 min |
| Citoyens (si > 1 h) | message portail "service maintenance" |

## 6. Post-mortem (obligatoire si > 30 min ou impact citoyen)
- Cause racine technique
- Pourquoi le monitoring n'a pas alerté plus tôt
- Plan d'action : capacité disque, alertes early, runbook ajusté
