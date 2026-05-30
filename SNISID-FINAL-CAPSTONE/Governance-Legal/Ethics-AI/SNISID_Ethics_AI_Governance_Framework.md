# SNISID — CADRE D'ÉTHIQUE ET DE GOUVERNANCE DE L'INTELLIGENCE ARTIFICIELLE

**Classification :** CADRE ÉTHIQUE — GOUVERNANCE IA
**Référence :** SNISID-ETHI-001
**Version :** 1.0
**Date :** 25 mai 2026

---

## 1. OBJECTIF

Encadrer l'utilisation de l'intelligence artificielle au sein du SNISID pour garantir que tous les systèmes d'IA restent gouvernés humainement, éthiquement responsables, transparents et exempts de biais discriminatoires.

---

## 2. CHAMP D'APPLICATION

Ce cadre s'applique à tous les systèmes d'IA utilisés dans le SNISID, notamment :

| Système IA | Fonction | Classification de Risque |
|-----------|---------|------------------------|
| Matching biométrique | Identification / vérification | Élevé |
| Détection de doublons | Déduplication du registre | Élevé |
| Reconnaissance faciale | Identification policière | Très élevé |
| Analyse de documents | Vérification d'authenticité | Moyen |
| Scoring de risque | Évaluation sécuritaire | Très élevé |
| Détection d'anomalies | Fraude, comportements suspects | Élevé |
| Traitement automatisé | Décisions administratives | Élevé |
| Chatbot / assistant | Service citoyen | Faible |
| Analytics prédictif | Planification de ressources | Moyen |

---

## 3. PRINCIPES ÉTHIQUES FONDAMENTAUX

| Principe | Description |
|----------|-------------|
| **Dignité humaine** | L'IA au service de l'humain, jamais l'inverse |
| **Non-discrimination** | Aucun biais basé sur l'ethnie, le genre, la religion, l'origine |
| **Transparence** | Fonctionnement explicable et documenté |
| **Responsabilité** | Un humain toujours responsable des décisions |
| **Proportionnalité** | Usage proportionné au besoin |
| **Contrôle humain** | Supervision humaine permanente |
| **Sécurité** | Robustesse contre les attaques adversariales |
| **Vie privée** | Minimisation des données, protection intégrée |

---

## 4. FONCTIONS DE GOUVERNANCE

### 4.1 Explainable AI (IA Explicable)

**Exigences d'explicabilité par niveau de risque :**

| Niveau de Risque | Exigences d'Explicabilité |
|-----------------|--------------------------|
| Très élevé | Explication individuelle pour chaque décision, traçabilité complète du raisonnement, documentation du modèle |
| Élevé | Explication disponible sur demande, facteurs principaux identifiés, documentation du modèle |
| Moyen | Documentation du fonctionnement général, métriques de performance publiées |
| Faible | Documentation basique du modèle |

**Standards d'explicabilité :**

| Composant | Exigence |
|-----------|----------|
| Documentation du modèle | Carte modèle (Model Card) pour chaque système IA |
| Facteurs de décision | Identification des facteurs ayant influencé chaque décision |
| Niveau de confiance | Score de confiance accompagnant chaque résultat |
| Données d'entraînement | Description de la nature et de la source des données |
| Limitations connues | Documentation des limites et cas d'échec |
| Explication citoyenne | Explication en langage clair pour les personnes affectées |

**Modèle de Model Card SNISID :**
```
╔════════════════════════════════════════╗
║         MODEL CARD — [NOM DU MODÈLE]  ║
╠════════════════════════════════════════╣
║ Version: X.Y                          ║
║ Date de déploiement: JJ/MM/AAAA       ║
║ Responsable: [Nom + Agence]           ║
║ Classification risque: [Niveau]       ║
║                                       ║
║ DESCRIPTION                           ║
║ [Ce que fait le modèle]               ║
║                                       ║
║ DONNÉES D'ENTRAÎNEMENT                ║
║ Source: [Description]                 ║
║ Volume: [Nombre d'exemples]           ║
║ Période: [Dates]                      ║
║ Biais connus: [Description]           ║
║                                       ║
║ PERFORMANCE                           ║
║ Précision globale: XX%                ║
║ Taux de faux positifs: X.X%           ║
║ Taux de faux négatifs: X.X%           ║
║ Performance par sous-groupe: [Détail] ║
║                                       ║
║ LIMITATIONS                           ║
║ [Cas où le modèle est moins fiable]   ║
║                                       ║
║ USAGE AUTORISÉ                        ║
║ [Finalités approuvées]                ║
║                                       ║
║ USAGE INTERDIT                        ║
║ [Finalités explicitement interdites]  ║
║                                       ║
║ SUPERVISION HUMAINE                   ║
║ [Niveau de supervision requis]        ║
║                                       ║
║ DERNIÈRE REVUE ÉTHIQUE: JJ/MM/AAAA   ║
║ PROCHAINE REVUE: JJ/MM/AAAA          ║
╚════════════════════════════════════════╝
```

### 4.2 Bias Controls (Contrôle des Biais)

**Types de biais surveillés :**

| Type de Biais | Description | Mitigation |
|-------------|-------------|-----------|
| Biais de données | Données d'entraînement non représentatives | Audit de diversité des données |
| Biais algorithmique | Algorithme favorisant certains groupes | Tests de parité |
| Biais de sélection | Sélection biaisée des cas d'entraînement | Échantillonnage stratifié |
| Biais de mesure | Qualité de mesure inégale entre groupes | Calibration par groupe |
| Biais historique | Reproduction de discriminations historiques | Analyse critique des données |
| Biais de confirmation | Renforcement de résultats attendus | Diversité d'évaluation |

**Processus de détection de biais :**
```
1. AVANT DÉPLOIEMENT
   - Audit de diversité des données d'entraînement
   - Tests de performance par sous-groupe démographique
   - Vérification de parité statistique
   - Rapport de biais pré-déploiement

2. PENDANT L'EXPLOITATION
   - Monitoring continu des métriques d'équité
   - Alertes automatiques en cas de déviation
   - Analyse trimestrielle des résultats par sous-groupe
   - Feedback loop citoyen

3. REMÉDIATION
   - Ré-entraînement avec données rééquilibrées
   - Ajustement des seuils
   - Modification algorithmique
   - Retrait du modèle si biais irréductible
```

**Métriques d'équité obligatoires :**

| Métrique | Définition | Seuil |
|----------|-----------|-------|
| Parité démographique | Taux de résultat positif égal entre groupes | Ratio ≥ 0.8 |
| Égalité des chances | Taux de vrais positifs égal entre groupes | Ratio ≥ 0.8 |
| Parité prédictive | Précision égale entre groupes | Ratio ≥ 0.8 |
| Calibration | Probabilités prédites calibrées par groupe | Déviation < 5% |

### 4.3 Human Oversight (Supervision Humaine)

**Niveaux de supervision humaine :**

| Niveau | Description | Systèmes concernés |
|--------|-------------|-------------------|
| Human-in-the-loop | Humain approuve chaque décision IA | Reconnaissance faciale policière, scoring de risque |
| Human-on-the-loop | Humain supervise et peut intervenir en temps réel | Matching biométrique, détection anomalies |
| Human-in-command | Humain définit les règles, supervise globalement | Analytics prédictif, chatbot |

**Exigences de supervision :**

| Exigence | Description |
|----------|-------------|
| Droit de veto | L'opérateur humain peut toujours annuler une décision IA |
| Bouton d'arrêt | Capacité de désactiver tout système IA immédiatement |
| Escalade | Cas complexes ou incertains escaladés automatiquement |
| Formation | Opérateurs formés à comprendre et évaluer les résultats IA |
| Fatigue | Rotation pour éviter la complaisance de supervision |
| Indépendance | Le superviseur n'est pas évalué sur la productivité du système IA |

### 4.4 AI Auditability (Auditabilité de l'IA)

| Composant Auditable | Détail |
|---------------------|--------|
| Données d'entraînement | Source, volume, diversité, étiquetage |
| Algorithme | Code source, architecture, hyperparamètres |
| Processus d'entraînement | Logs d'entraînement, métriques d'évaluation |
| Décisions individuelles | Entrées, facteurs, résultat, confiance |
| Performance en production | Métriques continues, dérive du modèle |
| Modifications | Historique des modifications du modèle |
| Incidents | Erreurs, faux positifs/négatifs significatifs |

**Exigences de journalisation IA :**
```json
{
  "prediction_id": "PRED-2026-XXXXX",
  "model_id": "MOD-BIOMETRIC-MATCH-V3",
  "model_version": "3.2.1",
  "timestamp": "2026-05-25T10:00:00Z",
  "input_hash": "SHA384:...",
  "output": {
    "decision": "MATCH",
    "confidence": 0.97,
    "threshold_applied": 0.95,
    "top_factors": ["minutiae_match: 0.98", "face_match: 0.96"]
  },
  "human_review": {
    "required": true,
    "reviewer": "AGT-12345",
    "decision": "CONFIRMED",
    "timestamp": "2026-05-25T10:02:30Z"
  }
}
```

### 4.5 Ethical Reviews (Revues Éthiques)

**Comité National d'Éthique Numérique :**

| Membre | Profil |
|--------|--------|
| Président | Éthicien / philosophe |
| Membre 1 | Juriste spécialisé droits fondamentaux |
| Membre 2 | Expert en IA |
| Membre 3 | Sociologue |
| Membre 4 | Représentant société civile |
| Membre 5 | Représentant des droits des femmes |
| Membre 6 | Représentant des minorités |

**Processus de revue éthique :**

| Étape | Description | Délai |
|-------|-------------|-------|
| Soumission | Fiche de demande de déploiement IA | - |
| Triage | Classification du niveau de risque | 5 jours |
| Analyse | Étude d'impact éthique | 30 jours |
| Audition | Présentation par l'équipe technique | 1 session |
| Délibération | Vote du comité | 15 jours |
| Décision | Approuvé / Approuvé avec conditions / Refusé | - |
| Suivi | Revue annuelle des systèmes approuvés | Annuel |

**Cas de saisine obligatoire du comité :**
- Tout système IA à risque élevé ou très élevé avant déploiement
- Toute modification majeure d'un système existant
- Tout incident éthique (discrimination détectée, erreur à impact)
- Sur demande de tout citoyen ou organisation (auto-saisine)

---

## 5. USAGES INTERDITS DE L'IA

| Usage Interdit | Justification |
|---------------|---------------|
| Scoring social généralisé | Atteinte à la dignité humaine |
| Surveillance de masse indiscriminée | Atteinte aux libertés fondamentales |
| Manipulation comportementale | Atteinte à l'autonomie |
| Décision automatisée sans recours | Atteinte au droit de défense |
| Profilage ethnique ou religieux | Discrimination |
| Prédiction criminelle individuelle | Présomption d'innocence |
| Reconnaissance d'émotions pour décisions administratives | Fiabilité insuffisante, dignité |

---

## 6. INDICATEURS ÉTHIQUES

| KPI | Cible | Mesure |
|-----|-------|--------|
| % systèmes IA avec Model Card | 100% | Continu |
| % systèmes IA avec revue éthique | 100% (risque élevé+) | Continu |
| Ratio de parité démographique | ≥ 0.8 pour tous les groupes | Trimestriel |
| % décisions IA supervisées humainement | 100% (très élevé), ≥ 95% (élevé) | Continu |
| Délai de correction de biais détecté | ≤ 30 jours | Sur événement |
| Taux de satisfaction des citoyens impactés | ≥ 80% | Semestriel |

---

*Document cadre préparé dans le cadre de la Phase 14 — SNISID National Legal Framework*
