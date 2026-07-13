# APPLICATION RUNBOOKS — SNISID
## Procédures d'Exploitation Applicative Nationale

| Attribut | Valeur |
|----------|--------|
| Document ID | SNISID-PH12-RUN-001 |
| Version | 1.0 |
| Statut | APPROUVÉ |

---

## 1. MOBILE SYNC FAILURE — Recovery

### 1.1 Détection

```
┌──────────────────────────────────────────────┐
│  SYMPTÔMES                                    │
├──────────────────────────────────────────────┤
│  🔴 Alert: "Sync failure rate > 5%"          │
│  🔴 Alert: "Sync queue backlog > 5000"       │
│  🔴 User reports: "Données non synchronisées"│
│  🟡 Dashboard: Sync health < 95%             │
└──────────────────────────────────────────────┘
```

### 1.2 Procédure

| Étape | Action | Responsable | Temps |
|-------|--------|-------------|-------|
| **1** | Vérifier l'alerte dans Grafana | Ops | 1 min |
| **2** | Identifier les appareils impactés | Ops | 2 min |
| **3** | Vérifier le serveur de sync | Ops | 2 min |
| **4a** | Si serveur HS → Redémarrer le service | Ops | 5 min |
| **4b** | Si réseau → Vérifier connectivité | Ops | 5 min |
| **5** | Forcer un sync de test | Ops | 2 min |
| **6** | Monitorer le taux de succès | Ops | 10 min |
| **7** | Si résolu → Fermer l'incident | Ops | 1 min |
| **8** | Si non résolu → Escalader niveau 2 | Ops | 1 min |

### 1.3 Commandes

```bash
# Vérifier le service sync
curl -s https://sync.snisid.gouv.ht/health | jq .

# Redémarrer le service sync
docker-compose -f /opt/snisid/docker-compose.yml restart sync-engine

# Voir les logs sync
journalctl -u snisid-sync -n 100 --no-pager

# Forcer un sync pour un device
curl -X POST https://sync.snisid.gouv.ht/api/v1/sync/force \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"device_id": "DEV-12345"}'
```

### 1.4 Vérification

```
Après résolution :
✅ Sync success rate > 98%
✅ Queue depth < 1000
✅ Aucune perte de données
✅ Appareils reconnectés
```

---

## 2. CITIZEN APP OUTAGE — Stabilisation

### 2.1 Détection

```
┌──────────────────────────────────────────────┐
│  SYMPTÔMES                                    │
├──────────────────────────────────────────────┤
│  🔴 Alert: "App crash rate > 0.5%"          │
│  🔴 Alert: "API latency p99 > 2s"           │
│  🔴 User reports: "App ne s'ouvre pas"       │
│  🟡 Support: Volume appels ×10              │
└──────────────────────────────────────────────┘
```

### 2.2 Procédure

| Étape | Action | Responsable | Temps |
|-------|--------|-------------|-------|
| **1** | Confirmer l'incident | Ops | 1 min |
| **2** | Identifier version et plateforme | Ops | 2 min |
| **3** | Vérifier dépendances (API, DB) | Ops | 3 min |
| **4** | Activer le mode dégradé si nécessaire | Ops | 2 min |
| **5a** | Si API HS → Failover vers backup | Ops | 5 min |
| **5b** | Si bug app → Rollback dernière version | Dev | 10 min |
| **6** | Communiquer statut aux utilisateurs | Support | 5 min |
| **7** | Monitorer recovery | Ops | 15 min |
| **8** | Root cause analysis | Dev | 24h |

### 2.3 Rollback Command

```bash
# Rollback Citizen App API
kubectl rollout undo deployment/citizen-app-api -n snisid

# Vérifier le rollback
kubectl rollout status deployment/citizen-app-api -n snisid

# Forcer cache refresh
curl -X POST https://api.snisid.gouv.ht/admin/cache/clear \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

### 2.4 Communication Template

```
SNISID — Statut Incident #INC-2026-XXXX

🟡 Dégradation détectée sur Citizen Super App
📱 Platformes: Android, iOS
🕐 Détecté: 2026-05-25T14:30:00Z

Actions en cours:
• Investigation en cours
• Mode dégradé activé

Prochaine mise à jour: dans 30 min
```

---

## 3. IDENTITY WALLET CORRUPTION — Recovery

### 3.1 Détection

```
┌──────────────────────────────────────────────┐
│  SYMPTÔMES                                    │
├──────────────────────────────────────────────┤
│  🔴 User: "Mon wallet ne s'ouvre pas"        │
│  🔴 App: "Wallet data integrity check failed"│
│  🟡 Log: "Secure storage corruption detected"│
└──────────────────────────────────────────────┘
```

### 3.2 Procédure

| Étape | Action | Responsable | Temps |
|-------|--------|-------------|-------|
| **1** | Isoler l'appareil | Support | 2 min |
| **2** | Vérifier l'intégrité du wallet | Support | 5 min |
| **3a** | Si backup disponible → Restaurer | Support | 10 min |
| **3b** | Si recovery phrase → Recovery | Support | 15 min |
| **4** | Vérifier les certificats restaurés | Support | 5 min |
| **5** | Forcer un sync de vérification | Support | 5 min |
| **6** | Clôturer le ticket avec confirmation | Support | 5 min |

### 3.3 Wallet Recovery

```
┌──────────────────────────────────────────────┐
│         WALLET RECOVERY FLOW                 │
├──────────────────────────────────────────────┤
│                                              │
│  Option A: Secure Cloud Backup               │
│  1. Authentifier l'utilisateur (MFA)        │
│  2. Vérifier identité officielle             │
│  3. Déclencher restauration                  │
│  4. Attendre confirmation                    │
│                                              │
│  Option B: Recovery Phrase                   │
│  1. Ouvrir wallet → "J'ai perdu mon wallet" │
│  2. Entrer la phrase de récupération (12 mots)│
│  3. Vérifier checksum                         │
│  4. Re-générer les clés                       │
│  5. Re-télécharger certificats               │
│                                              │
│  Option C: Bureau Gouvernemental             │
│  1. Se présenter avec pièce d'identité       │
│  2. Vérification biométrique                 │
│  3. Agent émet un nouveau wallet             │
│  4. Certificats ré-émis                      │
│                                              │
└──────────────────────────────────────────────┘
```

---

## 4. PUSH NOTIFICATION FAILURE — Recovery

### 4.1 Détection

```
┌──────────────────────────────────────────────┐
│  SYMPTÔMES                                    │
├──────────────────────────────────────────────┤
│  🔴 Alert: "Push delivery rate < 90%"        │
│  🔴 Alert: "FCM/APNs connection error"       │
│  🟡 User: "Je ne reçois pas les notifications"│
└──────────────────────────────────────────────┘
```

### 4.2 Procédure

| Étape | Action | Responsable | Temps |
|-------|--------|-------------|-------|
| **1** | Vérifier FCM/APNs status | Ops | 1 min |
| **2** | Vérifier credential expiry | Ops | 2 min |
| **3** | Vérifier rate limits | Ops | 2 min |
| **4a** | Si credentials expirés → Renouveler | Ops | 10 min |
| **4b** | Si rate limit → Throttle config | Ops | 5 min |
| **5** | Test push de vérification | Ops | 2 min |
| **6** | Monitorer delivery rate | Ops | 10 min |

### 4.3 Test Command

```bash
# Test notification send
curl -X POST https://notify.snisid.gouv.ht/api/v1/test \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "channels": ["push", "sms", "email"],
    "recipient": "admin@snisid.gouv.ht",
    "template": "test_notification"
  }'

# Voir les logs de notification
journalctl -u snisid-notification -n 50 --no-pager
```

---

## 5. OFFLINE NODE FAILURE — Synchronisation

### 5.1 Détection

```
┌──────────────────────────────────────────────┐
│  SYMPTÔMES                                    │
├──────────────────────────────────────────────┤
│  🔴 Device non synchronisé depuis > 7 jours  │
│  🔴 Queue depth > 10,000                     │
│  🟡 Alert: "Device storage > 90%"            │
└──────────────────────────────────────────────┘
```

### 5.2 Procédure

| Étape | Action | Responsable | Temps |
|-------|--------|-------------|-------|
| **1** | Identifier le device | Ops | 1 min |
| **2** | Vérifier dernier contact | Ops | 1 min |
| **3** | Essayer connexion distante | Ops | 5 min |
| **4a** | Si joignable → Forcer sync | Ops | 10 min |
| **4b** | Si pas joignable → Instructions terrain | Support | 5 min |
| **5** | Compresser et prioriser les données | Ops | 5 min |
| **6** | Monitorer sync completion | Ops | Variable |

### 5.3 Offline Node Recovery

```
┌──────────────────────────────────────────────┐
│     OFFLINE NODE RECOVERY                    │
├──────────────────────────────────────────────┤
│                                              │
│  Si le device est physique:                  │
│  1. Se connecter au réseau local             │
│  2. Lancer l'app de sync manuel             │
│  3. Vérifier la priorité de sync            │
│  4. Monitorer la progression                 │
│                                              │
│  Si le device est distant (>7 jours):        │
│  1. Tenter connexion VPN                    │
│  2. Forcer sync via SMS command              │
│  3. Si échec → Agent terrain requis          │
│  4. Sync via câble (mode air-gap)           │
│                                              │
└──────────────────────────────────────────────┘
```

---

## 6. INCIDENT MANAGEMENT SUMMARY

### 6.1 Runbook Reference

| Runbook | ID | Classification | SLA |
|---------|-----|---------------|-----|
| Mobile Sync Failure | RB-SYNC-001 | P1 | < 1h |
| Citizen App Outage | RB-APP-001 | P0 | < 30 min |
| Wallet Corruption | RB-WAL-001 | P1 | < 2h |
| Push Notification Failure | RB-NOT-001 | P2 | < 4h |
| Offline Node Failure | RB-OFF-001 | P2 | < 4h |
| Security Breach | RB-SEC-001 | P0 | < 15 min |
| Performance Degradation | RB-PERF-001 | P3 | < 8h |

### 6.2 Escalation Matrix

```
P0 ───▶ On-call Engineer (15 min)
 │
 ├──▶ Engineering Manager (30 min)
 │
 ├──▶ CTO (1h)
 │
 └──▶ CIO (2h)

P1 ───▶ On-call Engineer (1h)
 │
 ├──▶ Engineering Manager (4h)
 │
 └──▶ CTO (24h)

P2 ───▶ Team Lead (4h)
 │
 └──▶ Engineering Manager (24h)

P3 ───▶ Development Team (next business day)
```

---
*Fin du document — Application Runbooks v1.0*