# 💰 SNISID — National Financial Model

**Document N° :** SNISID-FIN-013
**Étape Phase 0 :** 13/16
**Principe :** *SNISID doit survivre politiquement et économiquement.*

---

## 1. Doctrine de Financement

SNISID est une **infrastructure publique critique**. Son financement doit être :
- **Pluriannuel** (sanctuarisé sur 5 ans minimum)
- **Mixte** (national + bailleurs, sans aliénation souveraine)
- **Soutenable** (transition vers autofinancement opérationnel d'ici 2032)
- **Transparent** (publication budget annuel + audit Cour des Comptes)

---

## 2. CAPEX — Investissements (2026–2030)

| Poste | Montant (USD) |
|-------|---------------|
| Datacenter primaire PaP (Tier III) | 20 M |
| Datacenter DR Cap-Haïtien | 12 M |
| 10 Edge nodes départementaux | 6 M |
| 100 Offline kits terrain (montée en charge) | 7 M |
| Plateforme logicielle (licences, intégration, dev custom) | 15 M |
| Équipements biométriques nationaux (AFIS, capteurs) | 8 M |
| PKI nationale (HSM, infra, audit) | 3 M |
| Réseau (fibre, satellite, équipements) | 5 M |
| Programme enrôlement de masse (matériel + déploiement) | 10 M |
| Cybersécurité (SOC outils, SIEM, EDR) | 4 M |
| Études, conseil, conception | 5 M |
| Formation et transfert de compétences | 5 M |
| **TOTAL CAPEX 5 ans** | **~100 M USD** |

---

## 3. OPEX — Exploitation Annuelle (en régime à partir de 2029)

| Poste | Annuel (USD) |
|-------|--------------|
| Masse salariale (180 ETP en régime) | 8-10 M |
| Énergie + connectivité datacenters | 2-3 M |
| Maintenance matérielle | 2 M |
| Licences logicielles (essentiellement OSS — support entreprise) | 1,5 M |
| Cybersécurité (audits, CTI, bug bounty, pentests) | 1 M |
| Communication & sensibilisation | 1 M |
| Formation continue | 0,5 M |
| Coût des cartes CIN (consommables) | 3-5 M |
| Astreintes, support 24/7 | 1 M |
| Imprévus (~5 %) | 1 M |
| **TOTAL OPEX annuel** | **~22-26 M USD** |

---

## 4. Sources de Financement

### 4.1 Budget National Haïtien
- Cible : ≥ 0,5 % du budget national sanctuarisé
- Justification : infrastructure critique, retour sur investissement (lutte fraude fiscale, identification populations vulnérables, amélioration recouvrement)

### 4.2 Bailleurs Internationaux

| Bailleur | Vocation | Montant potentiel |
|----------|----------|-------------------|
| **BID (Banque Interaméricaine)** | Identification + état civil + interop | 30-50 M USD |
| **Banque Mondiale (ID4D)** | Identification universelle | 25-40 M USD |
| **Union Européenne** | Gouvernance, état de droit | 10-20 M EUR |
| **USAID** | Cybersécurité, formation | 5-15 M USD |
| **ONU (PNUD)** | État civil, ODD 16.9 | 5-10 M USD |
| **GIZ (Allemagne)** | GovStack | 3-5 M EUR |
| **AFD (France)** | Modernisation État | 5-10 M EUR |
| **Banque Mondiale CAT-DDO** | Résilience | 10 M USD |

### 4.3 Recettes Propres (à partir de 2028)
- Frais de délivrance CIN (modique, exemption pour vulnérables)
- Frais de KYC pour banques/MFI/télécoms (B2B)
- Frais signature électronique qualifiée
- Frais services à valeur ajoutée (extraits multilingues, certifications)

**Cible recettes propres 2030 :** 8-12 M USD/an → couverture partielle OPEX.

---

## 5. Donor Governance

Pour éviter les dérives :
- **Conseil unique de coordination** des bailleurs sous présidence AND
- **Une seule architecture, un seul programme** (pas de silos par bailleur)
- Tous les SI financés par bailleurs respectent **standards SNISID**
- **Sortie progressive** des bailleurs planifiée (autonomie opérationnelle 2032)
- **No backdoor** : aucune obligation contractuelle imposant un éditeur étranger spécifique
- **Open-source priorisé** pour souveraineté pérenne

---

## 6. Modèle d'Amortissement

| Catégorie | Durée |
|-----------|-------|
| Génie civil datacenter | 25 ans |
| Équipements IT serveurs | 5 ans |
| Équipements réseau | 7 ans |
| Postes opérateurs / kits | 4 ans |
| Logiciels développés | 5 ans |
| HSM / sécurité physique | 7 ans |

---

## 7. Soutenabilité — Trajectoire

```
USD M/an
   ▲
30 │                              ┌─── OPEX total ───────
   │                              │
20 │       ┌──── OPEX ──────┐     │
   │       │                │     │
10 │  ┌────┘                └─────┘  Recettes propres
   │  │              ┌──────────────
 0 └──┴──────────────┴────────────────────────────────▶
    2026   2027   2028   2029   2030   2031   2032
```

**Hypothèse de soutenabilité :**
- 2026-2028 : 100 % financement national + bailleurs (CAPEX dominant)
- 2029-2030 : OPEX stabilisé, recettes propres commencent
- 2031+ : couverture progressive jusqu'à 50 % OPEX par recettes propres
- 2032+ : pleine soutenabilité avec budget national couvrant le reste

---

## 8. Risques Financiers & Mitigations

| Risque | Mitigation |
|--------|------------|
| Coupes budgétaires politiques | Loi-cadre verrouillant budget pluriannuel |
| Dépréciation gourde | Contrats USD/EUR pour CAPEX |
| Dépendance à un bailleur | Diversification (min. 4 bailleurs majeurs) |
| Sur-coûts projets | Provision 10 % imprévus + revues trimestrielles |
| Sanctions / gel d'aide | Plan de continuité dégradé (services prioritaires) |
| Corruption | Audit indépendant + Cour des Comptes + e-passation marchés |

---

## 9. Gouvernance Budgétaire

- Vote annuel au Parlement (ligne dédiée)
- Reporting trimestriel publié
- Audit annuel Cour des Comptes
- Plateforme citoyenne de suivi budgétaire (transparence)

---

## 10. KPI Financiers

| KPI | Cible |
|-----|-------|
| Taux d'exécution budgétaire | ≥ 90 % |
| Coût par citoyen enrôlé | < 8 USD |
| Coût par acte numérique généré | < 0,30 USD |
| Recettes propres / OPEX (2032) | ≥ 50 % |
| Audit annuel sans réserve majeure | 100 % |

---
*Fin du document — Étape 13/16*
