# 🔧 RUNBOOK — Model Rollback

**ID** : RB-ANL-004
**Sévérité max** : HIGH
**Propriétaire** : Data Science + Plateforme
**SLA résolution** : 1h

---

## 1. SYMPTÔMES

- `ModelDriftDetected` alerte
- Augmentation taux faux positifs (NRIC analyste)
- Métriques business dégradées (recall fraude ↓)
- Plaintes utilisateurs sur recommandations
- Distribution prédictions anormale

## 2. IMPACT

- Décisions IA potentiellement biaisées
- Risque atteinte droits citoyens
- Perte crédibilité IA gouvernementale

## 3. DIAGNOSTIC

```bash
# Lister versions modèle
mlflow models list-versions --name fraud_scoring

# Métriques actuelles
curl -s https://prometheus.snisid.ht/api/v1/query \
  --data-urlencode 'query=model_drift_score{model="fraud_scoring"}'

# Comparer distribution prédictions vs baseline
python scripts/drift_report.py --model fraud_scoring --window 7d
```

## 4. PROCÉDURE DE REMÉDIATION

### 4.1 Rollback KServe (production)
```bash
# Identifier version stable précédente (ex: v3.1.7)
PREV=v3.1.7

# Patcher l'InferenceService
kubectl patch isvc fraud-scoring -n ml-serving --type=merge \
  -p "{\"spec\":{\"predictor\":{\"model\":{\"storageUri\":\"s3://mlflow/artifacts/fraud_scoring/${PREV}/\"}}}}"

kubectl rollout status isvc fraud-scoring -n ml-serving
```

### 4.2 Mode shadow inversé
Optionnel : router 100 % trafic sur ancien modèle, 0 % nouveau, garder logs comparatifs.

```yaml
# istio VirtualService
http:
- match: [{headers: {x-canary: {exact: "true"}}}]
  route: [{destination: {host: fraud-scoring-v3-2-1}, weight: 0}]
- route: [{destination: {host: fraud-scoring-v3-1-7}, weight: 100}]
```

### 4.3 Notification & blocage
- Alerter NRIC
- Bloquer dépréciation des scores du nouveau modèle dans le case management
- Informer comité d'éthique IA

### 4.4 Investigation root cause
- Drift données amont (changement distribution features) ?
- Drift conceptuel (changement règles fraude) ?
- Bug pipeline d'entraînement ?
- Données labellisées biaisées ?

### 4.5 Plan de retour
1. Corriger root cause
2. Réentraîner sur données récentes propres
3. Tests biais + explainability
4. Shadow mode 30 jours
5. Validation comité éthique
6. Promotion progressive (canary 5% → 25% → 100%)

## 5. VÉRIFICATION

- [ ] Latence inférence OK
- [ ] Distribution prédictions retour baseline
- [ ] Recall / précision business stables
- [ ] Pas de pic de faux positifs (24h)
- [ ] NRIC analystes confirment qualité

## 6. POST-MORTEM

Obligatoire. Documenter dans `postmortems/`.
Mise à jour Model Card avec leçons apprises.
