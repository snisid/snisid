# 📒 SNISID — Runbooks Opérationnels

> **Phase 3 / Étape 10**
> Toute opération doit être : **répétable et documentée**.

| Runbook | Sujet | Severity |
|---------|-------|----------|
| [01](01-workflow-failure.md) | Échec workflow | Sev2-Sev1 |
| [02](02-kafka-outage.md) | Panne Kafka / Event Mesh | Sev1 |
| [03](03-bpmn-rollback.md) | Rollback BPMN | Sev1 |
| [04](04-fraud-escalation.md) | Escalation fraude critique | Sev1 |
| [05](05-offline-sync-recovery.md) | Récupération sync terrain | Sev2 |
| [06](06-pki-failure.md) | Panne PKI / TSA | Sev1 |
| [07](07-zeebe-cluster-failure.md) | Cluster Zeebe HS | Sev1 |
| [08](08-temporal-cluster-failure.md) | Cluster Temporal HS | Sev1 |
| [09](09-dc-failover.md) | Failover DC1 → DC2 / DC3 | Sev0 |
| [10](10-mass-events.md) | Catastrophe nationale | Sev0 |

## Niveaux de Sévérité

| Sev | Description | Réponse |
|-----|-------------|---------|
| 0 | Crise nationale | < 5 min, direction prévenue |
| 1 | Incident majeur (service critique HS) | < 15 min |
| 2 | Dégradation service | < 1 h |
| 3 | Anomalie sans impact citoyen | < 1 j |

## Convention

Chaque runbook contient :
1. Symptômes
2. Alertes Prometheus / dashboards Grafana
3. Diagnostic
4. Procédure de remédiation (étape par étape)
5. Vérification post-remédiation
6. Communication
7. Post-mortem (template)
