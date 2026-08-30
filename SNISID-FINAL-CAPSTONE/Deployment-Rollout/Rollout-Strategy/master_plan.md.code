# SNISID National Rollout & Migration Master Plan
## Système National d'Identité Sociale et d'Identification (SNISID) — République d'Haïti
**Version:** 4.2.0  
**Statut:** Approuvé pour Déploiement National  
**Classification:** Confidentiel - Gouvernemental  

---

## 1. Vision et Cadre Stratégique

Le présent document définit le **Master Plan de Déploiement National et de Migration** du Système National d'Identité Sociale et d'Identification (SNISID). L'objectif est de transformer le SNISID d'un système pilote industrialisé en une **infrastructure critique nationale**, déployée de manière opérationnelle et sécurisée sur l'ensemble du territoire de la République d'Haïti (les 10 départements géographiques).

```
+---------------------------------------------------------------------------------+
|                                 SNISID CORE                                     |
|              (Identité Unique, Registre Social, Sécurité Souveraine)            |
+---------------------------------------------------------------------------------+
                                       |
       +-------------------------------+-------------------------------+
       |                               |                               |
       v                               v                               v
[Déploiement Territorial]     [Onboarding Institutionnel]     [Migration de l'Historique]
  - 10 Départements             - ONI, MAST, DIE, PNH           - Facteur de Migration
  - 140 Communes                - Archives Nationales           - Déduplication Biométrique
  - Edge Nodes Hors-ligne       - Réconciliation d'Identité     - Nettoyage des Données
```

---

## 2. Déploiement Régional (Regional Deployment)

Le déploiement est structuré en **4 vagues progressives (Vagues A, B, C, D)** pour minimiser les risques opérationnels, tester la logistique territoriale et assurer le transfert de compétences.

### 2.1 Découpage Territorial et Vagues de Rollout

| Vague | Département | Communes Pilotes | Population Estimée | Priorité & Logistique |
| :--- | :--- | :--- | :--- | :--- |
| **Vague A** | **Ouest** (Département Pilote) | Port-au-Prince, Delmas, Pétion-Ville, Carrefour | ~3 100 000 | **Très Haute** (Hub Central, Proximité des Datacenters) |
| **Vague B** | **Nord, Nord-Est, Centre** | Cap-Haïtien, Fort-Liberté, Hinche | ~2 200 000 | **Haute** (Zone frontalière, Échanges commerciaux) |
| **Vague C** | **Artibonite, Sud, Sud-Est** | Gonaïves, Les Cayes, Jacmel | ~2 800 000 | **Moyenne** (Hubs agricoles et côtiers) |
| **Vague D** | **Nord-Ouest, Grand'Anse, Nippes** | Port-de-Paix, Jérémie, Miragoâne | ~1 500 000 | **Standard** (Zones enclavées, focus Connectivité Critique) |

### 2.2 Jalons par Département (Durée standard : 6 semaines par département)
1. **Semaines 1-2 : Audit de Readiness Régionale** (Énergie, Bâtiment physique, Réseau, Personnel).
2. **Semaine 3 : Installation de l'infrastructure locale** (Edge Nodes, Terminaux Biométriques).
3. **Semaine 4 : Formation intensive des opérateurs locaux** (Certifications SNISID).
4. **Semaine 5 : Phase de Coexistence Parallèle** (Saisie double et tests d'intégration réels).
5. **Semaine 6 : Go-Live Officiel et Bascule** (Activation du mode de production exclusif).

---

## 3. Onboarding des Agences Partenaires (Agency Onboarding)

L'intégration des institutions de l'État se fait via une architecture d'API unifiée (SNISID Secure Gateway) et des protocoles d'habilitation stricts.

```
+------------------+     +------------------+     +------------------+
|   ONI Legacy     |     |   Civil Registry |     |  Police Records  |
+------------------+     +------------------+     +------------------+
         |                        |                        |
         +------------------------+------------------------+
                                  |
                                  v (API Gateway / Migration Factory)
                       +----------------------+
                       |  SNISID RECONCILED   |
                       |   NATIONAL ENGINE    |
                       +----------------------+
```

### 3.1 Agences Cibles et Protocoles d'Onboarding

*   **Office National d'Identification (Collectif ONI) :**
    *   *Rôle :* Fournisseur de l'identité civile de base (Carte d'Identification Nationale - CIN).
    *   *Protocole :* Synchronisation biométrique temps réel (FAP-45 / Ten-Print) avec le système ABIS national du SNISID.
*   **Archives Nationales d'Haïti (ANH) :**
    *   *Rôle :* Référentiel des actes de naissance, mariage et décès.
    *   *Protocole :* Extraction par lots (Batch XML/JSON via SFTP sécurisé), rapprochement des registres d'état civil papier numérisés.
*   **Direction de l'Immigration et de l'Émigration (DIE) :**
    *   *Rôle :* Contrôle des passeports et des flux transfrontaliers.
    *   *Protocole :* Requêtes synchrones gRPC pour la validation d'identité aux postes frontières (Malpasse, Ouanaminthe, Belladère, Anse-à-Pitres) et aéroports.
*   **Police Nationale d'Haïti (PNH / DCPJ) :**
    *   *Rôle :* Vérification judiciaire et gestion des antécédents.
    *   *Protocole :* Interrogation cryptée unidirectionnelle. Pas d'accès direct au registre d'identité générale, mais croisement de jetons d'anonymisation (Tokens cryptographiques).
*   **Ministère des Affaires Sociales et du Travail (MAST) / SIMAST :**
    *   *Rôle :* Registre social des bénéficiaires de programmes nationaux.
    *   *Protocole :* Fusion du SIMAST historique dans la base de données de registre social du SNISID, garantissant l'unicité des aides publiques.

---

## 4. Déploiement de l'Infrastructure & Connectivité (Infrastructure & Connectivity)

L'architecture réseau d'Haïti présentant des vulnérabilités géographiques et de disponibilité électrique, le SNISID repose sur une architecture d'infrastructure hautement résiliente.

```
                                  +-----------------------+
                                  |   MAIN DATACENTER     |
                                  |   (Port-au-Prince)    |
                                  +-----------------------+
                                              ^
                                              | (Fibre Optique / Micro-onde)
                                              v
                                  +-----------------------+
                                  |  BACKUP DATACENTER    |
                                  |     (Cap-Haïtien)     |
                                  +-----------------------+
                                              ^
                      +-----------------------+-----------------------+
                      |                       |                       |
                      v                       v                       v
               [Edge Node - Ouest]     [Edge Node - Nord]     [Edge Node - Sud]
               - Cache local           - Cache local          - Cache local
               - Sync Hybride          - Sync Hybride         - Sync Hybride
```

### 4.1 Datacenters Nationaux
*   **Datacenter Principal (DC-1 - Port-au-Prince) :** Héberge le cluster Kubernetes maître (Bare-metal sécurisé), l'ABIS central, et la base de données relationnelle principale (PostgreSQL multi-région).
*   **Datacenter Secondaire de Secours (DC-2 - Cap-Haïtien) :** Réplication synchrone pour les données critiques, asynchrone pour les archives volumineuses. Capable de prendre le relais en moins de 5 minutes (RTO = 5 min, RPO = 10 sec) en cas de sinistre dans la zone métropolitaine.

### 4.2 Stratégie de Connectivité Multi-Transport
Pour assurer la continuité de service des bureaux de liaison communaux (BLC) d'identification :
1. **Liaison Primaire :** Fibre Optique terrestre ou Liaison micro-onde dédiée là où disponible (principaux centres urbains).
2. **Liaison Secondaire :** Liaison par Satellite Orbite Basse (Starlink Business) déployée sur chaque site distant avec double alimentation solaire autonome (panels photovoltaïques + batteries LiFePO4).
3. **Liaison de Secours :** Routeurs multi-SIM 4G/LTE (Digicel & Natcom) avec VPN IPSec chiffré AES-256-GCM.
4. **Mode Dégradé (Offline-First) :** Si toutes les liaisons échouent, le système bascule de manière transparente sur les bases locales des Edge Nodes (voir Étape 8).

---

## 5. Gestion de Crise et Continuité d'Activité (Crisis Fallback)

Une cellule de crise nationale (Deployment War Room) est immédiatement activée en cas d'événement majeur.

### 5.1 Matrice de Scénarios de Crise et Fallback

| Incident Identifié | Niveau de Gravité | Action Immédiate (Action de Contournement) | Procédure de Fallback (Reprise) |
| :--- | :--- | :--- | :--- |
| **Panne totale de réseau national** | Critique | Activation automatique du mode Offline-First sur les Edge Nodes locaux. | Capture locale des inscriptions et synchronisation physique par clé USB scellée ou attente de rétablissement. |
| **Instabilité politique / Troubles sécuritaires** | Majeure | Fermeture temporaire physique des bureaux de liaison affectés. | Bascule des opérations sur les unités mobiles d'enrôlement ou redirection des citoyens vers les communes voisines sécurisées. |
| **Corruption ou vol d'un Edge Node physique** | Critique | Révocation immédiate de la clé de chiffrement du nœud par le DC-1 (Remote Wipe instantané). | Les données locales sur le disque dur SSD chiffré par LUKS/AES-256 deviennent illisibles. Déploiement d'un nouveau nœud pré-configuré. |
| **Panne d'électricité prolongée (Grid Failure)** | Majeure | Commutation automatique sur le circuit d'alimentation solaire (Silo Solaire + Batterie). | Autonomie garantie de 72 heures sans soleil. Suspension des services accessoires pour préserver l'enrôlement de base. |

---

## 6. Approbation et Signatures

Ce document a été examiné et validé par les autorités compétentes et sert de feuille de route unique et obligatoire pour l'ensemble des acteurs impliqués dans le déploiement de la phase 15.

*   **Pour le Ministère de l'Intérieur et des Collectivités Territoriales (MICT) :** `[Signé - 25 Mai 2026]`
*   **Pour l'Office National d'Identification (ONI) :** `[Signé - 25 Mai 2026]`
*   **Pour la Direction Technique du SNISID :** `[Signé - 25 Mai 2026]`
