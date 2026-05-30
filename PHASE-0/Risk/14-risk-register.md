# ⚠️ SNISID — National Risk Register

**Document N° :** SNISID-RSK-014
**Étape Phase 0 :** 14/16
**Principe :** *Les risques doivent être connus AVANT construction.*

---

## 1. Méthodologie

Approche **EBIOS Risk Manager** + **ISO 31000**.

Évaluation : **Probabilité (1-5)** × **Impact (1-5)** = **Criticité (1-25)**.

Classification :
- 🟢 Faible (1-6)
- 🟡 Modéré (7-12)
- 🟠 Élevé (13-19)
- 🔴 Critique (20-25)

---

## 2. Registre Maître

| ID | Risque | Catégorie | P | I | C | Niv. | Stratégie | Owner |
|----|--------|-----------|---|---|---|------|-----------|-------|
| R01 | Cyberattaque étatique (APT) ciblant données identité | Cyber | 4 | 5 | 20 | 🔴 | Mitiger | CISO |
| R02 | Ransomware sur infrastructure | Cyber | 4 | 5 | 20 | 🔴 | Mitiger | CISO |
| R03 | Fuite massive de données biométriques | Cyber/Privacy | 3 | 5 | 15 | 🟠 | Mitiger | CISO + NDPA |
| R04 | DDoS sur services publics | Cyber | 4 | 3 | 12 | 🟡 | Mitiger | NOC + SOC |
| R05 | Fraude interne (administrateur malveillant) | Cyber/RH | 3 | 5 | 15 | 🟠 | Mitiger | CISO + RH |
| R06 | Corruption au moment de l'enrôlement | Intégrité | 4 | 4 | 16 | 🟠 | Mitiger | AND + IG |
| R07 | Détournement marchés publics | Financier | 3 | 4 | 12 | 🟡 | Mitiger | AND + CSC |
| R08 | Pannes datacenter (énergie, réseau) | Opérationnel | 4 | 4 | 16 | 🟠 | Mitiger | NOC |
| R09 | Catastrophe naturelle (séisme, cyclone) | Naturel | 4 | 5 | 20 | 🔴 | Transférer + Mitiger | AND |
| R10 | Inondation datacenter PaP | Naturel | 3 | 5 | 15 | 🟠 | Mitiger | Infra |
| R11 | Instabilité politique / changement régime | Politique | 4 | 5 | 20 | 🔴 | Mitiger | Direction |
| R12 | Coupe budgétaire majeure | Financier | 3 | 5 | 15 | 🟠 | Mitiger | AND + Finances |
| R13 | Dépendance bailleur unique | Financier/Souv. | 3 | 4 | 12 | 🟡 | Mitiger | AND |
| R14 | Lock-in éditeur propriétaire | Souveraineté | 3 | 4 | 12 | 🟡 | Éviter | Architectes |
| R15 | Insécurité limitant déploiement terrain | Sécurité | 4 | 4 | 16 | 🟠 | Mitiger | Logistique |
| R16 | Rejet citoyen / défiance | Adoption | 3 | 5 | 15 | 🟠 | Mitiger | Communication |
| R17 | Doublons / erreurs biométriques massives | Qualité | 3 | 4 | 12 | 🟡 | Mitiger | Data Steward |
| R18 | Discrimination algorithmique (IA biaisée) | Éthique | 3 | 4 | 12 | 🟡 | Mitiger | Comité Éthique |
| R19 | Vol matériel terrain (kits, biométrie) | Sécurité | 4 | 3 | 12 | 🟡 | Mitiger | Logistique |
| R20 | Compromission de la PKI nationale | Cyber | 2 | 5 | 10 | 🟡 | Mitiger | CISO |
| R21 | Indisponibilité prolongée (>24h) | Opérationnel | 2 | 5 | 10 | 🟡 | Mitiger | NOC |
| R22 | Litiges juridiques (CN, recours) | Légal | 3 | 3 | 9 | 🟡 | Mitiger | Juridique |
| R23 | Sanctions internationales | Géopolitique | 2 | 4 | 8 | 🟡 | Surveiller | Direction |
| R24 | Pénurie compétences techniques en Haïti | RH | 4 | 4 | 16 | 🟠 | Mitiger | RH + Formation |
| R25 | Fuite des talents (brain drain) | RH | 4 | 4 | 16 | 🟠 | Mitiger | RH |
| R26 | Coupure Internet nationale prolongée | Opérationnel | 3 | 4 | 12 | 🟡 | Mitiger | Offline-first |
| R27 | Crise sanitaire (épidémie) | Sanitaire | 2 | 3 | 6 | 🟢 | Surveiller | Direction |
| R28 | Vulnérabilité critique stack OSS (0-day) | Cyber | 4 | 3 | 12 | 🟡 | Mitiger | SOC |
| R29 | Sur-coûts projet (>20 %) | Financier | 3 | 3 | 9 | 🟡 | Mitiger | PMO |
| R30 | Manque adoption par agences | Organisationnel | 3 | 4 | 12 | 🟡 | Mitiger | AND |

---

## 3. Top 5 Risques Critiques — Plans de Traitement Détaillés

### R01 / R02 — Cyberattaques avancées
**Mesures :**
- SOC 24/7 + SIEM/SOAR/EDR
- Zero Trust + microsegmentation
- Backups immuables (WORM) + air-gap
- Tests d'intrusion semestriels
- Threat intelligence + chasse aux menaces
- Plan de réponse ransomware testé
- Cyber-assurance souveraine (à étudier)

### R09 — Catastrophe naturelle (séisme, cyclone)
**Mesures :**
- 2 datacenters distants > 250 km (PaP / Cap-Haïtien)
- Bâtiments anti-sismiques zone 4 + cyclones cat 5
- Backups offline multi-sites
- DRP testé semestriellement
- Convention d'hébergement de secours régional (République Dominicaine, CARICOM)

### R11 — Instabilité politique
**Mesures :**
- Loi-cadre adoptée à majorité qualifiée (verrouillage)
- Mandat AND déphasé du cycle électoral
- Engagement international (conventions BID/BM verrouillant)
- Reddition de comptes publique forte
- Personnel statutaire indépendant
- Continuité écrite (manuel de gouvernance)

### R06 — Corruption à l'enrôlement
**Mesures :**
- Biométrie obligatoire (élimine doublons frauduleux)
- Audit logs immuables
- Rotation des opérateurs
- Caméra obligatoire sur poste d'enrôlement
- Hotline anti-corruption (ULCC + AND)
- Sanctions exemplaires publiées

### R24 / R25 — RH (compétences + brain drain)
**Mesures :**
- Bourses + programme universitaire dédié (UEH, ESIH, UNDH...)
- Grilles salariales compétitives
- Plan carrière + formation continue
- Diaspora program (rapatriement temporaire experts haïtiens)
- Partenariats internationaux avec transfert obligatoire

---

## 4. Risques par Catégorie

```
Cyber:        ████████████████████  6 risques majeurs
Naturel:      ████████              2 critiques
Politique:    █████                 1 critique + autres
Financier:    ████████              4 risques
RH:           ████████              2 majeurs
Opérationnel: ███████               3 risques
Adoption:     ████                  2 risques
Éthique:      ██                    1 risque
```

---

## 5. Gouvernance des Risques

- **Comité Risques mensuel** (présidé par DG AND)
- Revue trimestrielle du registre
- Test annuel **Crisis Tabletop** (simulation crise)
- Test annuel **Disaster Recovery Drill**
- Reporting au CNN sur risques 🔴 / 🟠
- Risk appetite défini et signé par CNN

---

## 6. KPI Risk Management

| KPI | Cible |
|-----|-------|
| Risques 🔴 résiduels après traitement | 0 |
| Couverture plans de traitement | 100 % |
| Tests réalisés (DR + crisis) | ≥ 2/an |
| Délai notification incident grave aux autorités | < 24 h |
| Mise à jour registre | Trimestriel min. |

---
*Fin du document — Étape 14/16*
