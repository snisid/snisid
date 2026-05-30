# VOLUME 4 : Gouvernance Opérationnelle et Crise (Crisis Governance)
## Architecture de Gouvernance — SNISID

En Haïti, la crise n'est pas une anomalie, c'est un paramètre de conception. La gouvernance ne doit pas s'effondrer lorsque les communications sont coupées ou lorsque l'État est attaqué.

---

## 🌪️ CHAPITRE 1 : GOUVERNANCE EN CAS DE CATASTROPHE (DISASTER GOVERNANCE)

**Scénario :** Tremblement de terre majeur (ex: Magnitude 7+ à Port-au-Prince). Destruction du Palais National et du Datacenter Primaire.

### 1.1 Délégation d'Autorité Régionale (Devolution of Command)
*   Par défaut, l'autorité de validation réside à Port-au-Prince.
*   Si le centre est détruit, le protocole **"Ligne de Succession Numérique"** s'active. L'autorité de traitement de la base de données bascule automatiquement sur Cap-Haïtien (Failover).
*   L'autorité *administrative* (capacité à ordonner des validations d'urgence) est déléguée aux Préfets régionaux, qui peuvent utiliser leurs terminaux mobiles déconnectés (Kits MEK) pour valider l'identité des survivants et des cadavres via cryptographie PKI hors-ligne.

### 1.2 Continuité de l'État Civil (Continuity Governance)
*   Les décès massifs sont enregistrés via une "Procédure d'Exception Sismique", allégeant le circuit d'approbation (BPMN) pour privilégier la rapidité d'identification des victimes, tout en maintenant l'auditabilité stricte (WORM).

---

## ⚔️ CHAPITRE 2 : GOUVERNANCE EN CYBER-GUERRE (CYBERWAR GOVERNANCE)

**Scénario :** Attaque étatique ciblée (Wiper ou Ransomware) sur les institutions d'Haïti.

### 2.1 Loi Martiale Numérique
L'Incident Commander du SOC National a le pouvoir unilatéral de déclarer la "Loi Martiale Numérique".
1.  **Isolation (Island Mode) :** Coupure des liens internationaux.
2.  **Révocation Massive :** Annulation de tous les jetons d'accès API (X-Road) des Ministères jusqu'à preuve de leur intégrité.
3.  **No-Ransom Policy :** La loi interdit formellement le paiement de toute rançon. La gouvernance impose la restauration immédiate depuis les bunkers WORM (Air-Gapped), peu importe la perte de données des 2 dernières heures.

---

## 🗳️ CHAPITRE 3 : GOUVERNANCE DES ÉLECTIONS NATIONALES

Le Conseil Électoral Provisoire (CEP) dépend du SNISID pour l'établissement des listes électorales.

### 3.1 Séparation des Pouvoirs (Data Segregation)
*   Le CEP n'a pas accès à la base de données centrale. Il interroge une **"Vue Matérialisée Cryptographique"** (Snapshot) générée 90 jours avant l'élection.
*   **Audit d'Intégrité :** La liste électorale est hashée et publiée sur une blockchain d'État (ou registre public auditable). Les partis politiques peuvent télécharger l'empreinte cryptographique pour garantir que la liste n'a subi aucune altération entre son émission par le SNISID et le jour du vote.
*   **Biométrie Électorale :** Les terminaux de vote ne contiennent aucune base de données. Ils lisent la carte eID du citoyen et vérifient la signature électronique de l'AN-PKI (zéro connexion internet requise).
