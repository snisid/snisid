# 🔧 Runbook 01 — Échec d'un Workflow

**Severity :** Sev2 (par défaut), Sev1 si workflow CRITIQUE
**Owner :** Astreinte Workflow Engine

## 1. Symptômes
- Alerte Prometheus : `workflow_failures_total{wf="..."} rate(5m) > 0.05`
- Grafana : panneau "Workflow failure rate" passe au rouge
- Citoyen reçoit erreur "Demande en échec"
- Logs Zeebe : `INCIDENT` ou jobs en état `FAILED`

## 2. Diagnostic

```bash
# 1. Identifier le workflow et l'incident
zbctl status
zbctl list incidents | jq '.[] | select(.errorType=="JOB_NO_RETRIES")'

# 2. Examiner les logs du worker concerné
kubectl logs -n snisid -l app=zeebe-worker --tail=200 | grep ERROR

# 3. Tracing OTel
# Aller sur Tempo : trace_id depuis l'event Kafka audit.workflow.transition.v1
```

Causes courantes :
- **Service downstream HS** (biométrie, PKI...) → cf. runbooks dédiés
- **Données invalides** → `ValidationError` non-retryable
- **Timeout activité** → augmenter `startToCloseTimeout` ou réparer service
- **Fraude critique** → `FraudCriticalError` (cf. runbook 04)

## 3. Remédiation

### Cas A — Erreur retryable (service tiers HS)
```bash
# Une fois le service réparé : relancer manuellement les incidents
zbctl list incidents | jq -r '.[].key' | xargs -I{} zbctl resolve incident {}
```

### Cas B — Erreur non-retryable
1. Identifier la cause exacte (logs, traces).
2. Ouvrir un ticket WGO si correction de données nécessaire.
3. Annuler proprement via :
   ```bash
   zbctl cancel instance <processInstanceKey> --reason "ERR_DATA_INVALID"
   ```
4. Émettre l'événement de compensation manuelle :
   ```bash
   ./scripts/emit-event.sh civil-registry.birth.cancelled.v1 <payload.json>
   ```
5. Notifier le citoyen (workflow `notification.send`).

### Cas C — Workflow corrompu (bug code)
- Voir Runbook **03 — BPMN Rollback**.

## 4. Vérification

```bash
# Vérifier que le taux d'échec redescend
curl -s http://prometheus:9090/api/v1/query?query='rate(workflow_failures_total[5m])'
```
- Tableau de bord Grafana **"Workflow Health"** → vert
- Aucun incident en attente : `zbctl list incidents | jq length` = 0

## 5. Communication

| Audience | Canal | Délai |
|----------|-------|-------|
| Astreinte | PagerDuty | immédiat |
| Manager WGO | Email/Slack `#snisid-incidents` | 15 min |
| Direction | Email | 1 h |
| Citoyens | Notification SMS/Email | dès résolution |

## 6. Post-mortem (template)
- Date / durée :
- Workflows impactés :
- Citoyens impactés :
- Cause racine :
- Mesures correctives (court / long terme) :
- Action items dans `gitlab.snisid.ht/wgo/post-mortems`
