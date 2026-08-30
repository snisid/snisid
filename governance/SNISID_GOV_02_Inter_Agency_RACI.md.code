# VOLUME 2 : Gouvernance Inter-Agences (RACI)
## Architecture de Gouvernance — SNISID

L'efficacité du SNISID réside dans sa capacité à briser les silos administratifs via l'interopérabilité (X-Road). Cela nécessite une matrice stricte des responsabilités (RACI) pour éviter les conflits de juridiction entre les Ministères.

---

## 🤝 CHAPITRE 1 : PROPRIÉTÉ DES DONNÉES ET DES APIs

Dans l'écosystème SNISID, l'entité qui produit la donnée en est la "Propriétaire Autorisée" (Owner). Les autres agences sont de simples "Consommatrices" soumises à approbation.

| Domaine de Donnée | Propriétaire (Data Owner) | Consommateurs (Exemples) | API Externe (X-Road) |
| :--- | :--- | :--- | :--- |
| **Identité Biométrique (AFIS)** | ONI (Office National d'Identification) | DCPJ, Immigration | `api.oni.gouv.ht/bio-match` |
| **État Civil (Naissances/Décès)** | ANH (Archives Nationales d'Haïti) | ONI, Santé, Éducation | `api.anh.gouv.ht/civil-registry` |
| **Filiation Pénale (Casier)** | DCPJ (Police Judiciaire) / MJSP | ONI, Immigration | `api.dcpj.gouv.ht/background` |
| **Statut Fiscal (NIF)** | DGI (Direction Générale des Impôts) | Douane, Banques | `api.dgi.gouv.ht/tax-status` |
| **Passeports & Visas** | Direction de l'Immigration (DIE) | Douane, ONI | `api.die.gouv.ht/travel` |
| **Identité Électorale** | CEP (Conseil Électoral Provisoire) | (Usage restreint) | `api.cep.gouv.ht/voter-roll` |

### 1.1 Gouvernance de l'Approbation (Approval Chains)
Si le Ministère de la Santé (MSPP) souhaite accéder à l'API d'État Civil de l'ANH pour tracer les statistiques de mortalité :
1.  **Requête Légale :** Le MSPP soumet une demande formelle justifiant la base légale.
2.  **Validation Gouvernance :** Le Bureau de l'Interopérabilité du SNISID valide la sécurité technique.
3.  **Approbation Data Owner :** Le Directeur de l'ANH signe cryptographiquement l'approbation.
4.  **Application SOAR :** La politique d'accès X-Road est automatiquement mise à jour (Policy-as-Code).

---

## 📊 CHAPITRE 2 : MATRICE RACI OPÉRATIONNELLE

*(Responsible, Accountable, Consulted, Informed)*

Pour la gestion d'un incident critique (ex: Suspicion de faux certificats de naissance massivement injectés depuis un hôpital).

| Phase de l'Incident | ONI | ANH | DCPJ (Police) | SNISID SOC |
| :--- | :--- | :--- | :--- | :--- |
| **Détection de l'Anomalie (UEBA)** | I | C | I | **R / A** |
| **Confinement Technique (Isolation IP)** | I | I | I | **R / A** |
| **Enquête Métier (Faux en écriture)** | C | **R / A** | C | I |
| **Enquête Judiciaire / Arrestations** | I | C | **R / A** | C (Preuves) |
| **Révocation des Faux NNI** | **R** | **A** | I | C |

---

## ⏱️ CHAPITRE 3 : GOUVERNANCE DES SLA (Service Level Agreements)

Les agences interconnectées sont liées par des contrats de performance légaux.
*   **Temps de Réponse API (P99) :** La DGI exige que l'API de validation d'identité de l'ONI réponde en moins de 2 secondes.
*   **Disponibilité (Uptime) :** Si l'API des Archives Nationales tombe en dessous de 99.9%, le tableau de bord public de transparence signale une défaillance de service, déclenchant un audit par le Conseil National.
*   **Escalade de Performance :** Les pannes inter-agences sont escaladées via un canal d'urgence dédié du SOC, permettant un diagnostic "Cross-Agency" (OpenTelemetry) pour identifier l'origine exacte du ralentissement (Réseau, BDD, ou Code).
