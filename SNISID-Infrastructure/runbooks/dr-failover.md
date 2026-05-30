# SNISID RUNBOOK — DR Failover (Core DC → DR DC)
**Classification:** TOP-SECRET  
**Version:** 4.0.0  
**RTO cible:** < 30 minutes  
**Fréquence test:** Mensuel (dernier lundi)

---

## 1. Topologie active-active

| DC | Rôle | Clusters | Poids DNS |
|----|------|----------|-----------|
| Core | Primaire actif | Core, Identity, Data, BPMN, Cyber, Obs | 100 |
| DR   | Secondaire actif | Core-DR, Identity-DR, Data-DR, BPMN-DR, Cyber-DR, Obs-DR | 0 (standby) |

> **Active-active** = données répliquées en continu, mais DNS dirige 100% vers Core. DR est chaud (workloads déployés, prêts à servir).

## 2. Triggers failover

| Condition | Durée | Action |
|-----------|-------|--------|
| API Core indisponible | > 3 min | Failover automatique DNS + Istio |
| Vault scellé Core | > 1 min | Promotion DR Vault |
| Ceph HEALTH_ERR Core | > 5 min | Promotion Ceph DR + failover |
| Blackout national Core | Immédiat | DR total manuel (hors auto) |

## 3. Procédure failover automatique

```bash
# Exécuté par CronJob/Operator dans namespace snisid-cyber
/opt/snisid/scripts/failover-operator.sh
```

### 3.1 DNS Failover (CoreDNS interne)
```yaml
# Patch ConfigMap coredns
apiVersion: v1
kind: ConfigMap
metadata:
  name: snisid-dns-failover
  namespace: kube-system
data:
  Corefile: |
    snisid.gouv.local {
      loadbalance round_robin
      # Failover : Core (down) → DR (promoted)
      template IN A {
        match "^api\\.core\\.snisid\\.gouv\\.local$"
        answer "{{ .Name }} 30 IN A 10.2.10.11"  # DR VIP
        fallthrough
      }
      forward . /etc/resolv.conf
    }
```

### 3.2 Istio Locality Failover
```yaml
# DestinationRule — priorité DR quand Core unhealthy
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: snisid-api-locality
  namespace: istio-system
spec:
  host: "*.snisid.gouv.local"
  trafficPolicy:
    connectionPool:
      tcp:
        maxConnections: 100
    loadBalancer:
      simple: LEAST_CONN
      localityLbSetting:
        enabled: true
        failover:
          - from: core
            to: dr
    outlierDetection:
      consecutive5xxErrors: 5
      interval: 30s
      baseEjectionTime: 30s
```

### 3.3 Ceph Promotion DR
```bash
# Sur cluster Ceph DR
ceph osd pool ls  # vérifier pools snisid-rbd-tier0, snisid-rbd-tier1
rbd mirror pool promote snisid-rbd-tier0  # promotion écriture DR
rbd mirror pool promote snisid-rbd-tier1
# Le mirroring inverse (DR → Core) commence quand Core revient
```

### 3.4 Kafka MirrorMaker2 switch
```bash
# Pause mirroring Core → DR
kubectl exec mm2-cluster -n snisid-data -- mm2-admin --stop source->target
# Activer production directe sur DR topics
kubectl patch kafkatopic identites.nationales-dr -p '{"spec":{"config":{"min.insync.replicas":2}}}'
```

### 3.5 ArgoCD promotion
```bash
# Patch des Applications pour utiliser values-dr-prod.yaml
argocd app set snisid-core-api --helm-set-file values=values-dr-prod.yaml
argocd app sync snisid-core-api
# Sync tous les apps du projet national
argocd app sync -l snisid.gov.region=dr --prune
```

## 4. Procédure failback (retour Core)

> **Règle :** Le failback n'est JAMAIS automatique. Validation manuelle obligatoire.

1. Vérifier Core 100% sain pendant 30 minutes
2. Rétrograder Ceph DR (demote → resync Core)
3. Réactiver MirrorMaker2 Core → DR
4. Modifier DNS Core poids 100, DR poids 0
5. Modifier Istio locality : priorité Core
6. Resync ArgoCD vers values-core-prod.yaml
7. Test bout-en-bout : enrollment biométrique → API → DB → Ceph
8. Validation IGC avant fermeture incident

## 5. Drill mensuel (dernier lundi 02h00)

```bash
# 1. Notifier SOC (drill programmé)
# 2. Exécuter CronJob snisid-dr-drill
kubectl create job --from=cronjob/snisid-dr-drill drill-$(date +%s) -n snisid-cyber

# 3. Observer le failover (doit être < 30 min)
# 4. Valider services sur DR
./scripts/validate-dr-health.sh

# 5. Exécuter failback manuel
# 6. Rapport drill → wiki.interne.snisid.gouv.local/dr-reports/YYYY-MM
```

---

*Toute erreur de failover est un incident critique de niveau national.*
