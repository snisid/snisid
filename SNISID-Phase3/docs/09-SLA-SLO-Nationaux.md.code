# 📏 SLA / SLO NATIONAUX

> **Phase 3 / Étape 9** — Administration mesurable.
> Version : 1.0.0

---

## 1. Définitions

- **SLA (Service Level Agreement)** : engagement contractuel envers le citoyen et les administrations.
- **SLO (Service Level Objective)** : objectif interne, marge avant le SLA.
- **SLI (Service Level Indicator)** : mesure réelle observée.
- **Error Budget** : marge d'erreur tolérée avant action corrective.

Règle d'or :
```
SLI ≤ SLO ≤ SLA
```

---

## 2. SLA / SLO par workflow critique

| Workflow | SLA citoyen | SLO interne | Mesure (SLI) |
|----------|-------------|-------------|--------------|
| `civil-registry.birth.simple` | **24 h** | 18 h p95 | `bpmn_duration{wf="civil.birth.simple"}` |
| `civil-registry.birth.recognition` | 72 h | 48 h p95 | idem |
| `civil-registry.birth.late-declaration` | 30 j | 20 j p95 | idem |
| `civil-registry.death.standard` | 24 h | 12 h p95 | idem |
| `civil-registry.death.disaster` | 7 j | 4 j p95 | idem |
| `civil-registry.marriage.civil` | 7 j | 5 j p95 | idem |
| `civil-registry.divorce.administrative` | 30 j | 20 j p95 | idem |
| `identity.enrollment.standard` | 24 h | 12 h p95 | idem |
| `identity.verification.online` | **5 min** | 30 s p99 | idem |
| `identity.verification.biometric` | 30 s | 5 s p99 | idem |
| `identity.revocation.administrative` | 24 h | 4 h p95 | idem |
| `judicial.validation.act` | 48 h | 24 h p95 | idem |
| `judicial.suspension.identity` | 4 h | 1 h p95 | idem |
| `elections.voter.validation` | 1 h | 5 min p95 | idem |
| `fraud.detection.automated` | 1 min | 10 s p99 | idem |
| `escalation.sla.breach` | 5 min | 1 min p99 | idem |

---

## 3. SLO Plateforme (transverses)

| Composant | SLO |
|-----------|-----|
| Zeebe gateway disponibilité | 99,95 % (mensuel) |
| Temporal disponibilité | 99,95 % |
| Kafka cluster disponibilité | 99,95 % |
| Latence p99 émission Kafka | < 100 ms |
| Latence p99 signature PKI | < 250 ms |
| Lag consumer max | < 10 000 messages |
| Réussite workflows critiques | ≥ 99,9 % |
| Disponibilité PKI/TSA | 99,9 % |

---

## 4. Politique Error Budget

| Service | Budget mensuel | Action si dépassement |
|---------|----------------|----------------------|
| Workflows CRITIQUES | 0,1 % | Freeze déploiements + audit |
| Workflows ÉLEVÉS | 0,5 % | Plan correctif obligatoire |
| Workflows MOYENS | 1 % | Revue trimestrielle |
| Kafka mesh | 0,05 % | Astreinte renforcée |

---

## 5. Escalation (cf. BPMN `escalation.sla.breach`)

| Niveau | Délai après breach | Acteurs notifiés |
|--------|-------------------|------------------|
| L1 | 0 min | Manager opérationnel |
| L2 | 15 min sans résolution | Astreinte + PagerDuty |
| L3 | 30 min sans résolution | Direction nationale + WGO |
| L4 | 60 min sans résolution | Présidence ONI + Cabinet |

---

## 6. Mesure & Reporting

- **Source de vérité** : événements Kafka `audit.workflow.transition.v1`
- **Calcul** : Prometheus + recording rules (`workflow_duration_seconds_bucket`)
- **Visualisation** : Grafana dashboard `SNISID-SLO-National`
- **Rapport public** : portail `https://transparence.snisid.ht` (extraits anonymisés)

---

## 7. Sanctions / Compensations (volet citoyen)

| Violation SLA | Conséquence |
|---------------|-------------|
| Naissance > 24 h sans escalade | Notification ombudsman |
| Identité non émise > 30 j | Procédure express + indemnité administrative |
| Verification > 5 min répétée | Audit obligatoire région |
| Suspension judiciaire > 4 h | Rapport au Procureur Général |

---

**Annexe :** Voir `sla/sla-catalog.yaml` pour la définition machine-readable.
