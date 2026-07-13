# 📶 Runbook 05 — Récupération Sync Terrain (Offline)

**Severity :** Sev2 (Sev1 en crise nationale)
**Owner :** Astreinte Sync Hub

## 1. Symptômes
- Backlog `offline.batch.uploaded.v1` croît
- Kits terrain signalent erreurs sync
- Conflits CRDT > seuil
- Citoyens : enrôlements terrain n'apparaissent pas en central

## 2. Diagnostic
```bash
# Backlog par région
curl -s https://sync-hub.snisid.ht/admin/backlog | jq

# Taux de conflits
promql 'rate(offline_conflict_detected_total[1h])'

# État des kits
kubectl exec -it sync-hub-0 -- snisid-sync-ctl kits status
```

## 3. Remédiation

### Cas A — Connectivité dégradée
1. Vérifier liens fibre/4G/Starlink des sites concernés
2. Augmenter le **batch interval** côté kits :
   ```
   snisid-edge-cli set sync.interval=15m  (au lieu de 5m)
   ```

### Cas B — Conflits CRDT massifs
1. Lancer le job de réconciliation :
   ```bash
   kubectl create job --from=cronjob/crdt-reconciler crdt-fix-$(date +%s)
   ```
2. Si conflits non auto-résolvables : escalader vers WGO + LVB pour décision manuelle.

### Cas C — Hub central HS
1. Bascule vers réplica DC2 :
   ```bash
   kubectl -n snisid patch svc sync-hub -p '{"spec":{"selector":{"dc":"dc2"}}}'
   ```
2. Reprendre la résorption.

## 4. Vérification
- Backlog redescend sous le seuil normal (< 10 000)
- Conflits résorbés
- Nouveaux uploads acquittés (`offline.sync.completed.v1`)

## 5. Communication
- Sites terrain concernés : SMS Astreinte
- WGO + Direction : si > 4 h
