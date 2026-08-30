# SNISID National Master Plan 2026–2035
**Version** : 1.0 | **Date** : 2025-03-27 | **Statut** : 🟡 En validation

## 1. Gouvernance
- Créer le **Conseil National de Souveraineté Numérique** (CNSN) — 2026
- Créer l’**Autorité Nationale d’Identité** (ANI) — 2026
- Créer le **SOC National** — 2027

## 2. Architecture technique
- Modèle des données : NNU + attributs + événements
- Bus : Kafka multi-clusters, redondé géographiquement
- Stockage : PostgreSQL (Citus) + S3 compatible pour blobs

## 3. Workflow Factory (BPMN)
- Processus clés : enrollment, verification, recovery, revocation
- Moteur : Camunda Platform 8 (self-managed)

## 4. PKI
- Déploiement de la Root CA Nationale (2026)
- Distribution des certificats aux agences (2027)

## 5. SOC & Cybersécurité
- SIEM : Wazuh / Elastic Security
- SOAR : Shuffle / TheHive
- DFIR : outillage standardisé (Autopsy, KAPE)

## 6. Offline-first
- Edge nodes dans chaque département (10 sites)
- Kits mobiles pour les communes enclavées

## 7. Déploiement terrain
- Phase pilote : 3 départements (2027)
- Phase nationale : 10 départements (2029)

## 8. Plateforme données
- Data lake national, accès contrôlé par politique

## 9. Kubernetes (Infrastructure)
- Cluster principal : bare-metal, haute disponibilité, 3 sites
- GitOps : ArgoCD

## 10. Budget (estimé)
- Phase 0 (2025) : 0,5 M$ (études, architecture)
- Phase 1 (2026) : 2 M$ (gouvernance, PKI)
- Phase 2 (2027) : 8 M$ (infra, terrain)
- Total 2026-2035 : ≈ 45 M$

## 11. Formation
- Programme « SNISID Certified Engineer » — 2026
- Académie nationale du numérique — 2027

---

# National Domain Model

| Domaine          | Criticité | Owner proposé | Workflows clés                     |
|------------------|-----------|---------------|------------------------------------|
| Citizen          | CRITIQUE  | ANI           | Enrollment, Verification, Recovery |
| Identity         | CRITIQUE  | ANI           | NNU lifecycle                      |
| Civil Registry   | CRITIQUE  | Ministère Intérieur | Naissance, Mariage, Décès      |
| Justice          | ÉLEVÉ     | Ministère Justice | Casier judiciaire, mandats      |
| Elections        | ÉLEVÉ     | CEP           | Inscription, Vote                  |
| Security         | CRITIQUE  | PNH           | Flagging, Alerte                   |
| Tax              | ÉLEVÉ     | DGI           | Contribuable, Paiement             |
| Health           | MOYEN     | MSPP          | Dossier médical simplifié          |
| Education        | MOYEN     | MENFP         | Diplômes, identification scolaire  |
