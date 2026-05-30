# 🕵 Runbook 04 — Escalation Fraude Critique

**Severity :** Sev1 (Sev0 si réseau organisé)
**Owner :** Cellule Anti-Fraude + Police judiciaire + WGO

## 1. Symptômes
- Alerte Prometheus : `fraud_score > 0.9` plus de N fois en 5 min
- Spike `fraud.detected.v1` sur une commune / un agent / une fenêtre temporelle
- Signalement humain (agent ONI, citoyen)
- Anomalie graphe (cluster de NIN liés)

## 2. Procédure

1. **Geler les workflows** concernés :
   ```bash
   ./scripts/freeze-workflow.sh identity.enrollment.standard \
     --filter "context.commune=='Cap-Haïtien' AND context.agentId=='AGT-1234'"
   ```
2. **Suspendre les identités** suspectes :
   - Workflow `identity.suspension.judicial` (besoin ordre judiciaire ; à défaut, ordre administratif temporaire 72 h).
3. **Saisir** la cellule anti-fraude (`fraud-investigators`).
4. **Saisir** la police judiciaire si organisé.
5. **Préserver les preuves** :
   ```bash
   ./scripts/forensic-snapshot.sh --case-id <CASE_ID> \
     --include audit.workflow.transition.v1,identity.*,fraud.*,security.* \
     --output s3://snisid-forensic/<CASE_ID>/
   ```
6. **Communiquer** au WGO + Direction + Procureur de la République.

## 3. Vérification
- Identités suspendues : `kubectl exec -it identity-svc -- snisid-ctl identity status <NIN>` → `SUSPENDED`
- Workflows gelés : confirmé via Camunda Operate UI
- Snapshot forensic : checksum SHA-384 enregistré dans le case

## 4. Suivi
- Ouvrir un **dossier judiciaire** via `judicial.investigation.fraud`
- Rapport quotidien à la cellule jusqu'à clôture
- Post-mortem WGO + LVB + Cyber
