# SNISID Regional Activation Strategy
## Stratégie d'Activation Territoriale par Département — République d'Haïti

---

## 1. Cadre de Déploiement Géographique

Le déploiement territorial du SNISID ne peut pas s'effectuer d'un seul bloc sans risquer de surcharger l'infrastructure centrale et de générer du désordre opérationnel sur le terrain. La **Stratégie d'Activation Régionale** impose un modèle d'activation structuré, basé sur la vérification systématique de la maturité locale de chaque département.

```
       [Évaluation de Maturité Régionale (Readiness Scorecard)]
                                |
                                v (Seuil minimum : 85%)
             [Ordre d'Activation Territoriale (Green-Light)]
                                |
       +------------------------+------------------------+
       |                                                 |
       v                                                 v
[Déploiement Physique & Logistique]             [Onboarding & Formation]
 - Acheminement Edge Nodes                       - Certifications Opérateurs
 - Liaison Satellite Starlink                    - Simulations Incidentologie
 - Raccordement Énergie Solaire                  - Campagne de Com Locale
```

---

## 2. Indicateurs de Maturité Régionale (Regional Readiness Scorecard)

Avant d'autoriser le déploiement physique dans un département, la cellule nationale de pilotage évalue la préparation locale selon 5 axes critiques. Le score global de préparation doit atteindre **au moins 85/100** pour déclencher le feu vert (Green-Light) de déploiement.

### 2.1 Grille d'Évaluation de la Maturité

| Catégorie | Indicateur de Validation | Poids | Condition de Validation |
| :--- | :--- | :--- | :--- |
| **Sécurité (Sec)** | Niveau de menace locale évalué par la DCPJ / PNH inférieur au seuil critique. | 25% | Sécurisation du transport de matériel sensible et des locaux des bureaux de liaison. |
| **Énergie & Site (Eng)** | Disponibilité d'une source d'énergie primaire stabilisée et d'un kit solaire autonome complet. | 20% | Autonomie minimale de 48 heures vérifiée par test de décharge batterie. |
| **Connectivité (Net)** | Liaison Starlink ou liaison terrestre validée avec un temps de latence inférieur à 100ms vers le DC Central. | 20% | Double antenne (Active/Backup) installée et fonctionnelle. |
| **Ressources Humaines (HR)** | Au moins 4 agents formés et certifiés par poste d'enrôlement biométrique. | 20% | Réussite au test théorique et pratique de gestion d'identité (Score minimum 90%). |
| **Gouvernance Locale (Gov)** | Signature du protocole d'entente avec les maires, préfets et responsables ONI locaux. | 15% | Établissement du comité de coordination départemental. |

---

## 3. Matrice de Préparation de Connectivité (Connectivity Readiness Matrix)

La connectivité réseau est adaptée au relief géographique difficile de chaque région d'Haïti :

```
                                CONNECTIVITY PROFILES
+----------------------+--------------------------+-------------------------------------+
| Département          | Profil Topologique       | Solution de Connectivité Retenue    |
+----------------------+--------------------------+-------------------------------------+
| Ouest                | Urbain / Semi-montagneux | Fibre Optique + Backup 4G LTE/Multi |
| Nord / Nord-Est      | Côtier / Plaine          | Liaison Micro-onde + Satellite LEO  |
| Sud / Grand'Anse     | Montagneux / Côtier      | Satellite LEO (Starlink) + Solaire  |
| Centre / Artibonite  | Plateaux / Plaines       | Satellite LEO + Backup 4G/3G Dual   |
| Nord-Ouest / Nippes  | Enclavé / Montagneux     | Satellite LEO Exclusif + Solaire    |
+----------------------+--------------------------+-------------------------------------+
```

*   **Profil A (Urbain Connecté) :** Liaison fibre optique primaire de 100 Mbps symétrique + Liaison de sauvegarde 4G LTE multi-opérateurs (Digicel & Natcom).
*   **Profil B (Régional Standard) :** Liaison satellite LEO (Starlink Business) de 150 Mbps avec antenne motorisée auto-chauffante + Liaison micro-onde locale.
*   **Profil C (Zone Enclavée / Isolée) :** Liaison satellite Starlink à basse consommation d'énergie + Point d'accès Wi-Fi local restreint pour les agents + Mode Offline synchrone toutes les 24h.

---

## 4. Programme de Formation Locale (Local Training & Certification)

Le succès sur le terrain repose entièrement sur la compétence des opérateurs locaux du SNISID. Un programme d'habilitation de 5 jours est obligatoire :

### 4.1 Programme de Formation des Opérateurs de Liaison (POL)

*   **Jour 1 : Fondations et Déontologie de l'Identité**
    *   Cadre légal de la protection des données personnelles en Haïti.
    *   Lutte contre la fraude documentaire et la corruption d'identité.
*   **Jour 2 : Maîtrise des Matériels d'Enrôlement**
    *   Utilisation des caméras d'acquisition faciale (Norme ICAO/OACI).
    *   Capture des empreintes digitales (Capteurs FAP-45 certifiés FBI) et capture de l'iris.
*   **Jour 3 : Opérations Offlines & Edge Computing**
    *   Gestion de la base de données locale du Edge Node.
    *   Résolution de conflits de données en local.
    *   Procédure de sauvegarde physique cryptée.
*   **Jour 4 : Simulations et Gestion d'Incidents**
    *   Simulation de pannes d'énergie, de perte de connectivité ou d'agression sur le site.
    *   Exercices de communication de crise avec la population.
*   **Jour 5 : Examen Final & Certification**
    *   Évaluation pratique de rapidité et d'exactitude (enrôlement complet de test en moins de 3 minutes).
    *   Remise de la carte d'accès cryptographique personnelle de l'opérateur (YubiKey sécurisée).

---

## 5. Processus d'Approbation Régional (Regional Sign-Off Gateway)

Le déploiement effectif d'une région suit une série de portes de contrôle obligatoires :

```
[Porte 1 : Audit de Readiness] -> [Porte 2 : Déploiement Physique] -> [Porte 3 : Activation Logique] -> [Go-Live]
```

L'activation logique d'un département nécessite la signature conjointe du Directeur de Déploiement SNISID et du Préfet du département concerné, attestant que toutes les étapes ont été validées conformément aux standards de sécurité nationaux.
