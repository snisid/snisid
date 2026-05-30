# Duplicate Detection Engine

> **Moteur national de détection de doublons d'identité SNISID**  
> **Version :** 1.0.0  
> **Classification :** SOUVERAIN — RESTREINT  
> **Dernière mise à jour :** 2026-05-25

---

## 1. OBJECTIF

Empêcher la création d'identités multiples pour un même individu. Le moteur de dé-duplication est la **barrière antifraude primaire** du SNISID.

> **Les doublons doivent déclencher investigation.**

---

## 2. FONCTIONS DU MOTEUR

### 2.1 Biometric Deduplication

**Pipeline de Vérification :**

```
[Capture Biométrique] → [Qualité Check] → [Liveness Check] → [Extraction Features] → [Comparaison N:N] → [Scoring] → [Décision]
```

**Algorithmes de Comparaison :**

| Biométrie | Algorithme | Seuil match | Seuil suspect |
|-----------|-----------|-------------|---------------|
| Empreinte digitale | Minutiae + Deep Learning | ≥ 95% | 80-94% |
| Reconnaissance faciale | ArcFace / FaceNet | ≥ 90% | 75-89% |
| Iris | Daugman / Gabor | ≥ 92% | 78-91% |
| Voix | x-vector + PLDA | ≥ 85% | 70-84% |

### 2.2 Identity Correlation

**Vérification multi-dimensionnelle :**

| Dimension | Critères |
|-----------|----------|
| Données démographiques | Nom + date de naissance + lieu de naissance + noms parents |
| Documents | Même numéro acte de naissance / carte d'identité / passeport |
| Géographie | Même adresse d'enrôlement / commune / département |
| Liens familiaux | Frères/sœurs/parents enregistrés |
| Comportement | Mêmes accès historiques / services utilisés |

### 2.3 Fraud Scoring

```yaml
fraud_scoring_engine:
  scoring_factors:
    biometric_match:        {weight: 40, threshold: 80}
    demographic_similarity: {weight: 25, threshold: 70}
    document_correlation:   {weight: 20, threshold: 60}
    geographic_proximity:   {weight:  5, threshold: 50}
    behavioral_patterns:    {weight:  5, threshold: 40}
    temporal_anomalies:     {weight:  5, threshold: 30}

  risk_levels:
    LOW:      {range: "0-30",  action: "Auto-approve"}
    MEDIUM:   {range: "31-60", action: "Enhanced verification"}
    HIGH:     {range: "61-80", action: "Manual review required"}
    CRITICAL: {range: "81-100", action: "Block + investigation"}

  scoring_formula: |
    total_score = Σ(factor_weight × factor_score) / Σ(factor_weights)
```

### 2.4 Human Review

**Déclencheurs :**
- Score MEDIUM → vérification renforcée
- Score HIGH → revue manuelle obligatoire
- Score CRITICAL → blocage immédiat

**Équipe de revue :**
- Senior Enrollment Officer
- Identity Analyst
- Fraud Investigator (si CRITICAL)

**Processus :**
1. Notification automatique
2. Ouverture ticket investigation
3. Analyse dossier complet
4. Entretien citoyen (si nécessaire)
5. Décision : approuver / rejeter / escalader
6. Documentation décision
7. Mise à jour registre

### 2.5 Judicial Escalation

**Déclencheurs :**
- Fraude confirmée
- Tentative d'usurpation d'identité
- Document falsifié détecté
- Réseau de fraude identifié

**Délais :**
- Détection → Suspension : < 1 heure
- Suspension → Rapport : < 24 heures
- Rapport → Décision judiciaire : < 72 heures

---

## 3. DÉCLENCHEURS DE DÉTECTION

### 3.1 Automatiques

| Événement | Action immédiate |
|-----------|-----------------|
| Nouvel enrôlement | Vérification N:N complète |
| Modification identité | Re-vérification ciblée |
| Vérification biométrique | Comparaison avec registre |
| Connexion suspecte | Analyse comportementale |

### 3.2 Temporels

| Tâche | Fréquence |
|-------|-----------|
| Scan complet registre | Quotidien (02:00 UTC) |
| Analyse corrélations | Hebdomadaire |
| Audit doublons potentiels | Mensuel |

### 3.3 Manuels

| Action | Rôle requis |
|--------|-------------|
| Signalement fraude | Tout agent |
| Demande investigation | Senior Officer |
| Escalade judiciaire | Fraud Investigator |

---

## 4. RÉSOLUTION DE DOUBLONS

### 4.1 Algorithme

1. **Identification** — Trouver toutes les identités candidates
2. **Analyse** — Déterminer l'identité "source" (plus ancienne)
3. **Fusion** — Conserver source, migrer données, préserver historique
4. **Nettoyage** — Marquer doublons comme merged, révoquer NNU
5. **Audit** — Journaliser, signer, notifier

### 4.2 Matrice de Décision

| Scénario | Action | Rôle décisionnel |
|----------|--------|-----------------|
| 2 identités, 1 biométrie match | Fusion | Système automatique |
| 2 identités, données similaires | Investigation | Senior Officer |
| 2 identités, fraude confirmée | Révocation + Escalade | Fraud Investigator |
| 3+ identités, même individu | Fusion multiple | Comité IAM |
| Identité usurpée | Révocation + Judiciaire | Parquet |

---

## 5. MÉTRIQUES

| Métrique | Cible | Mesure |
|----------|-------|--------|
| Taux de faux positifs | < 0.1% | Mensuel |
| Taux de faux négatifs | < 0.01% | Mensuel |
| Temps de détection | < 5 secondes | Par transaction |
| Taux de résolution | > 99% | Mensuel |

---

## 6. CONFIGURATION TECHNIQUE

```yaml
biometric_search_engine:
  technology: Elasticsearch + FAISS
  fingerprints:
    algorithm: Minutiae Cylinder-Code
    index_type: IVF-PQ
  face:
    algorithm: ArcFace (ResNet-100)
    index_type: HNSW
  iris:
    algorithm: Daugman
    index_type: IVF-Flat
  voice:
    algorithm: x-vector + PLDA
    index_type: HNSW
  search_params:
    timeout_ms: 5000
    max_results: 10
    similarity_threshold: 0.80
```

---

> **Le Duplicate Detection Engine est la première ligne de défense contre la fraude identitaire.**
