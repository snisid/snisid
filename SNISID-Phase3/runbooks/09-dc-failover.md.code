# 🔴 Runbook 09 — Failover Datacenter (DC1 → DC2 / DC3)

**Severity :** Sev0 (crise nationale technique)
**Owner :** WGO + Direction ONI + Platform Engineering
**Délai cible :** RTO < 15 min, RPO < 60 s

## 1. Quand déclencher
- DC1 (Port-au-Prince) totalement inaccessible :
  - Panne électrique > 30 min sans diesel
  - Catastrophe naturelle (séisme, ouragan)
  - Cyberattaque massive
  - Coupure fibre principale + secours
- Tests de bascule planifiés (semestriels obligatoires WGO)

## 2. Pré-requis (vérifier mensuellement)
- ✅ Réplication Kafka MirrorMaker2 DC1↔DC2 active, lag < 5 s
- ✅ PostgreSQL Camunda en streaming replication vers DC2
- ✅ Cassandra Temporal en Active-Active DC1/DC2
- ✅ DNS GSLB configuré (TTL ≤ 30 s)
- ✅ Certificats PKI/TSA répliqués DC2
- ✅ Last DR drill < 6 mois

## 3. Procédure de bascule DC1 → DC2

### Étape 1 — Activation du plan crise national
```bash
# Démarrer le workflow national
zbctl create instance escalation.crisis.national \
  --variables '{"trigger":"DC1_DOWN","approver":"DG_ONI","timestamp":"'$(date -Iseconds)'"}'
```

### Étape 2 — Bascule DNS (GSLB)
```bash
# Pointer api.snisid.ht, portal.snisid.ht, etc. vers DC2
./scripts/gslb-failover.sh --from dc1-pap --to dc2-cap
# Vérifier propagation
dig +short api.snisid.ht  # doit retourner les IPs DC2
```

### Étape 3 — Promotion PostgreSQL DC2
```bash
# Promouvoir la replica
kubectl -n postgres-dc2 exec postgres-replica-0 -- \
  patronictl failover --master postgres-dc2-replica-0
```

### Étape 4 — Kafka : DC2 devient leader
```bash
# Élection forcée des nouveaux leaders côté DC2
kubectl -n kafka-dc2 exec kafka-13-0 -- \
  kafka-leader-election.sh --bootstrap-server kafka-13:9093 \
  --election-type PREFERRED --all-topic-partitions

# Stopper MirrorMaker2 DC1→DC2 (devient DC2→DC1 quand DC1 revient)
kubectl -n kafka scale deploy/mirrormaker-dc1-to-dc2 --replicas=0
```

### Étape 5 — Zeebe DC2
```bash
# Les brokers DC2 prennent le leadership des partitions
kubectl -n zeebe-dc2 scale sts/zeebe-broker --replicas=3
zbctl status --address zeebe-gateway-dc2:26500
```

### Étape 6 — Temporal DC2
```bash
# Cassandra DC2 a déjà tous les rings (Active-Active)
# Frontend Temporal DC2 démarre automatiquement
kubectl -n temporal-dc2 scale deploy/temporal-frontend --replicas=3
```

### Étape 7 — PKI / TSA DC2
```bash
# Activation des HSM secondaires
./scripts/pki-activate-dc2.sh
```

### Étape 8 — Mode offline-first amplifié
```bash
# Activer côté tous les kits terrain via push
./scripts/edge-broadcast.sh --command "offline.mode.enable"
```

### Étape 9 — Vérification e2e
```bash
# Workflow test bout-en-bout
zbctl create instance identity.verification.online \
  --variables '{"nin":"HT-NIN-TEST-0001"}' \
  --address zeebe-gateway-dc2:26500
# Doit retourner un résultat sous 30s
```

## 4. Procédure DC2 → DC3 (DR ultime)

DC3 (Les Cayes) est **standby asynchrone** :
- RPO ≤ 5 min
- RTO ≤ 1 h

À déclencher si DC1 ET DC2 inaccessibles (catastrophe nationale).
Procédure quasi-identique, mais :
- Cassandra DC3 doit rejoindre l'anneau (`nodetool bootstrap`)
- Kafka : promotion des brokers 19-24 + reconstruction des replicas
- **Décision exclusive Direction ONI + Présidence**

## 5. Retour à la normale (DC1 réparé)

```bash
# 1. Re-synchroniser DC1 depuis DC2
./scripts/dc-resync.sh --source dc2 --target dc1

# 2. Vérifier intégrité (Merkle audit chain)
./scripts/audit-chain-verify.sh --range "ALL"

# 3. Bascule planifiée DC2 → DC1 (en heures creuses)
./scripts/gslb-failover.sh --from dc2-cap --to dc1-pap --planned

# 4. Reprise MirrorMaker2 DC1↔DC2 bidirectionnel
kubectl -n kafka scale deploy/mirrormaker-dc1-to-dc2 --replicas=1
kubectl -n kafka scale deploy/mirrormaker-dc2-to-dc1 --replicas=1
```

## 6. Vérification finale
- 4 dashboards Grafana → tous verts
- Aucune perte de transaction confirmée par `audit.chain.verify`
- Lag réplication < 5 s
- Test e2e workflow critique réussi

## 7. Communication

| Audience | Délai | Canal |
|----------|-------|-------|
| Présidence ONI + Primature | < 5 min | Téléphone direct |
| Tous Ministères | < 15 min | Email + SMS direction |
| Citoyens | < 30 min | Communiqué officiel portail + radio nationale |
| Diaspora | < 1 h | Réseaux sociaux ONI |

## 8. Post-mortem
- **OBLIGATOIRE** dans les 7 jours, présenté au Conseil des Ministres
- Audit indépendant si > 4 h indisponibilité
- Mise à jour du plan DR + drill spécifique
