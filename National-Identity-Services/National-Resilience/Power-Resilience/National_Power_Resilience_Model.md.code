# National Power Resilience Model

## 1. Objectif
Assurer l'autonomie énergétique des infrastructures critiques SNISID face aux coupures prolongées, pénuries carburant et catastrophes affectant l'énergie.

## 2. Capacités
| Domaine | Support | Description |
|---|---:|---|
| Generator redundancy | Oui | générateurs N+1/N+2 pour sites critiques |
| Solar backup | Oui | solaire pour charges essentielles et edge |
| Battery systems | Oui | UPS/BESS pour transition et autonomie |
| Fuel reserve governance | Oui | réserves, contrats, sécurité, rotation carburant |

## 3. Autonomie cible
| Site | Autonomie cible | Configuration |
|---|---:|---|
| Primary National DC | 7-14 jours | UPS, générateurs N+1, carburant sécurisé |
| Secondary National DC | 7-14 jours | équivalent primary pour P0/P1 |
| Regional DR Site | 3-7 jours | générateur + batteries + solaire si possible |
| Offline Vault | 7 jours ponctuels | faible charge, générateur/batterie |
| Field/Edge Kits | 24-72 h | batteries, panneaux solaires, chargeurs |

## 4. Priorisation charges
| Priorité | Charges |
|---|---|
| E0 | IAM, registre identité, command center, sécurité physique |
| E1 | APIs inter-agences, stockage critique, monitoring, comms crise |
| E2 | portails citoyens, reporting, postes essentiels |
| E3 | analytics lourds, environnements non production |

## 5. Gouvernance carburant
Stock minimal par site, rotation carburant, contrats multi-fournisseurs, accès sécurisé, mesure consommation, plan réapprovisionnement, inventaire quotidien en L3/L4.

## 6. Mode économie énergie
Arrêter E3 → limiter E2 → réduire batch → activer régional/edge → déplacer charges → réserver énergie à P0/P1.

## 7. Tests et KPI
| Test/KPI | Fréquence/Objectif |
|---|---|
| bascule UPS | mensuel |
| démarrage générateurs charge réelle | mensuel |
| test autonomie site | trimestriel |
| inventaire carburant | hebdo, quotidien crise |
| autonomie restante | dashboard NRCC |
