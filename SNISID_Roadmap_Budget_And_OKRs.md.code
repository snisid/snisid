# SNISID : Roadmap, Budget et Gouvernance par OKRs (v2.0)

**Classification :** SOUVERAIN / DIRECTION GÉNÉRALE  
**Recommandation de Référence :** SNISID v2.0 — MP-010

Ce document définit la planification temporelle ajustée au risque d'Haïti, le cadre budgétaire consolidé (~54M USD) assorti de télémétrie financière, et le modèle de validation des phases par OKRs automatisés.

---

## 📂 Fichiers de Référence dans le Workspace

* **Gouvernance & Suivi** :
  * [ROADMAP.md](file:///c:/Users/sopil/Desktop/snisid%20system/ROADMAP.md) — Plan global consolidé et vision finale.
  * [verify_okrs.py](file:///c:/Users/sopil/Desktop/snisid%20system/pki/scripts/verify_okrs.py) — Moteur automatique d'évaluation des OKRs de Phase 2 (ABIS).

---

## 1. Planification Temporelle Ajustée aux Risques (Buffer 20%)

Haïti présente des risques logistiques et politiques chroniques (pénuries de carburant, cyclones, blocages routiers). Pour garantir la faisabilité, un **buffer systématique de 20%** est appliqué sur chaque phase.

### Tableau Comparatif des Délais :

| Phase | Livrable Principal | Durée Théorique Nominale | Durée Ajustée aux Risques (+20%) | Échéance Calendaire Estimée |
|---|---|---|---|---|
| **Phase 1** | Fondation, DC1, PKI, IAM base, Root CA | Mois 1 - Mois 12 | **Mois 1 - Mois 14.4** (14.4 mois) | T1 Année 2 |
| **Phase 2** | Enrôlement, ABIS, API Gateway, eID | Mois 6 - Mois 18 | **Mois 7.2 - Mois 21.6** (14.4 mois) | T3 Année 2 |
| **Phase 3** | ONI, DGI, DCPJ, ANH via X-Road | Mois 12 - Mois 24 | **Mois 14.4 - Mois 28.8** (14.4 mois) | T1 Année 3 |
| **Phase 4** | SIEM opérationnel, SOC 24/7, DC2 | Mois 18 - Mois 28 | **Mois 21.6 - Mois 33.6** (12.0 mois) | T3 Année 3 |
| **Phase 5** | 15M citoyens, ISO 27001, PRA validé | Mois 24 - Mois 36 | **Mois 28.8 - Mois 43.2** (14.4 mois) | T2 Année 4 |
| **TOTAL** | **Vision complète nationale** | **36 mois** | **43.2 mois** | **43.2 mois au total** |

### Jalons Critiques GO/NO-GO Binaires :
Le passage d'une phase à la suivante est conditionné par un jalon binaire (0 ou 1) validé par le système. Un échec bloque le déploiement de la phase suivante via ArgoCD.

* **Jalon GO/NO-GO Phase 1 -> 2 :** 
  * *Critère 1 :* Cérémonie de la Root CA exécutée et signée cryptographiquement par le quorum 5-of-9.
  * *Critère 2 :* Disponibilité du cluster de base DC1 $\ge 99.99\%$.
* **Jalon GO/NO-GO Phase 2 -> 3 :**
  * *Critère 1 :* Enrôlement de 100K citoyens tests effectué sans aucune régression.
  * *Critère 2 :* ABIS validé en production avec un taux de fausse acceptation (FAR) $< 0.001\%$.
* **Jalon GO/NO-GO Phase 3 -> 4 :**
  * *Critère 1 :* 4 agences clés interconnectées en X-Road avec signature électronique mutuelle fonctionnelle.
* **Jalon GO/NO-GO Phase 4 -> 5 :**
  * *Critère 1 :* Failover DC1 ↔ DC2 testé automatiquement avec un RTO $< 15$ min et RPO $< 1$ min.

---

## 2. Dashboard Financier en Temps Réel

Le NOC intègre une télémétrie financière qui interroge les ERPs gouvernementaux pour calculer le taux de déviation budgétaire par phase.

### Charge Utile JSON de Télémétrie Financière :

```json
{
  "project_id": "SNISID-HAITI",
  "timestamp": "2026-05-24T21:34:00Z",
  "consolidated_budget_usd": 54000000.0,
  "actual_spend_to_date_usd": 11500000.0,
  "phases": [
    {
      "phase_id": "PHASE_1",
      "planned_budget_usd": 12000000.0,
      "actual_spend_usd": 11500000.0,
      "committed_outstanding_usd": 400000.0,
      "deviation_percent": -0.83,
      "alert_status": "NOMINAL"
    },
    {
      "phase_id": "PHASE_2",
      "planned_budget_usd": 18000000.0,
      "actual_spend_usd": 1200000.0,
      "committed_outstanding_usd": 15000000.0,
      "deviation_percent": 10.0,
      "alert_status": "WARNING"
    }
  ],
  "burn_rate_monthly_usd": 950000.0,
  "projection_at_completion_usd": 54850000.0,
  "overall_deviation_percent": 1.57,
  "trigger_soc_escalation": false
}
```

### Règle d'alerte financière :
Si la déviation d'une phase (`deviation_percent`) dépasse **10%**, le système :
1. Verrouille les budgets d'extra-engagement non approuvés.
2. Émet automatiquement un rapport d'alerte mensuel envoyé par email cryptographiquement signé au Comité de Pilotage National.

---

## 3. OKRs par Phase et Mesure Automatique

Chaque phase est évaluée selon des OKRs quantitatifs mesurés par des requêtes de monitoring automatiques.

### OKRs Phase 2 — Validation ABIS en Production
* **Objectif (O) :** Valider la précision et la résilience de l'ABIS en conditions réelles.
  * **KR 2.1 :** Atteindre un taux de fausse acceptation (FAR) $< 0.001\%$ sur un panel de 100K citoyens tests.
  * **KR 2.2 :** Compléter 100K enrôlements biométriques sans régression système.
  * **KR 2.3 :** Détecter et empêcher 100% des tentatives de doublons biométriques (0 doublon non détecté).

### Requête SQL de Vérification Automatique (KR 2.3) :
Cette requête est exécutée quotidiennement pour vérifier qu'aucun doublon biométrique n'a été inséré sans être bloqué ou flaggé pour examen :

```sql
-- Vérifie s'il existe des enregistrements validés partageant un template similaire non flaggé
SELECT COUNT(*) as undetected_duplicates
FROM citizen_biometrics cb1
JOIN citizen_biometrics cb2 ON cb1.citizen_id <> cb2.citizen_id
WHERE cb1.iris_hash = cb2.iris_hash 
   OR cb1.fingerprint_slap_hash = cb2.fingerprint_slap_hash
   AND cb1.status = 'APPROVED' 
   AND cb2.status = 'APPROVED';
```

---

*Ce cadre garantit la transparence budgétaire et technique absolue du SNISID.*
