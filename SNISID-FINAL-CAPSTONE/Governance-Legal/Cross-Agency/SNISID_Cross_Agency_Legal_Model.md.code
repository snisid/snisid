# SNISID — MODÈLE JURIDIQUE INTER-AGENCES

**Classification :** CADRE JURIDIQUE — INTEROPÉRABILITÉ
**Référence :** SNISID-XAGN-001
**Version :** 1.0
**Date :** 25 mai 2026

---

## 1. OBJECTIF

Encadrer juridiquement les échanges de données entre les agences connectées au SNISID, garantissant légalité, traçabilité, responsabilité et protection des droits des citoyens dans tout partage de données.

---

## 2. PRINCIPES FONDAMENTAUX

| Principe | Description |
|----------|-------------|
| Légalité | Tout échange requiert une base légale explicite |
| Nécessité | Seules les données nécessaires sont échangées |
| Proportionnalité | L'échange est proportionné à la finalité |
| Réciprocité contrôlée | Droits et obligations mutuels |
| Traçabilité | Chaque échange est journalisé |
| Responsabilité | Chaque agence est responsable de son usage |
| Réversibilité | L'accès peut être révoqué |

---

## 3. FONCTIONS SUPPORTÉES

### 3.1 Data Sharing Agreements (Accords de Partage de Données)

**Structure d'un accord type :**

| Section | Contenu |
|---------|---------|
| 1. Parties | Agences signataires, responsables |
| 2. Objet | Finalité précise de l'échange |
| 3. Base légale | Texte de loi autorisant l'échange |
| 4. Données concernées | Liste exhaustive des données échangées |
| 5. Personnes concernées | Catégories de personnes dont les données sont partagées |
| 6. Flux de données | Direction (unidirectionnel/bidirectionnel), fréquence, volume |
| 7. Modalités techniques | API, format, protocole, chiffrement |
| 8. Sécurité | Mesures de protection obligatoires |
| 9. Durée et renouvellement | Durée déterminée, conditions de renouvellement |
| 10. Droits des personnes | Mécanismes de respect des droits |
| 11. Incidents | Procédure de notification mutuelle |
| 12. Audit | Droits d'audit réciproques |
| 13. Résiliation | Conditions et effets de la résiliation |
| 14. Responsabilité | Répartition des responsabilités |
| 15. Signatures | Signatures qualifiées des responsables |

**Accords-cadres SNISID :**

| Accord | Parties | Données | Finalité |
|--------|---------|---------|----------|
| AC-001 | ONI ↔ PNH | Identité + photo | Vérification d'identité policière |
| AC-002 | ONI ↔ DGI | NNI + état civil | Identification fiscale |
| AC-003 | ONI ↔ Justice | Identité complète | Procédures judiciaires |
| AC-004 | PNH ↔ Justice | Données d'enquête | Poursuites pénales |
| AC-005 | ONI ↔ Santé | NNI + état civil | Identification patients |
| AC-006 | ONI ↔ Éducation | NNI + état civil | Inscription scolaire |
| AC-007 | PNH ↔ Immigration | Identité + biométrie | Contrôle frontalier |
| AC-008 | Toutes agences ↔ CERT-HT | Logs de sécurité | Cybersécurité |

### 3.2 Legal Interoperability (Interopérabilité Juridique)

**Cadre d'interopérabilité juridique :**

| Dimension | Standard |
|-----------|---------|
| Sémantique | Vocabulaire commun des données (ontologie SNISID) |
| Syntaxique | Formats d'échange standardisés (JSON-LD, XML SNISID) |
| Organisationnelle | Processus alignés entre agences |
| Juridique | Base légale commune, accords harmonisés |
| Technique | API standardisées, protocoles communs |

**Registre d'interopérabilité :**
- Catalogue des services d'échange disponibles
- Spécifications techniques de chaque service
- Conditions juridiques d'utilisation
- Niveaux de service garantis
- Processus d'habilitation

### 3.3 Jurisdiction Rules (Règles de Juridiction)

| Règle | Description |
|-------|-------------|
| Compétence territoriale | Les données SNISID restent sous juridiction haïtienne |
| Compétence fonctionnelle | Chaque agence est compétente dans son domaine |
| Résolution de conflits | Le CNGN arbitre les conflits inter-agences |
| Priorité de juridiction | En cas de conflit : sécurité nationale > justice > administratif |
| Souveraineté des données | L'agence source reste propriétaire de ses données |
| Droit applicable | Droit haïtien exclusivement |

**Matrice de compétences :**

| Donnée | Agence propriétaire | Agences autorisées | Conditions |
|--------|--------------------|--------------------|-----------|
| Identité civile | ONI | Toutes (selon accord) | Finalité justifiée |
| Données biométriques | ONI | PNH, Immigration | Enquête / contrôle |
| Casier judiciaire | Justice | PNH | Enquête autorisée |
| Données fiscales | DGI | Justice | Réquisition judiciaire |
| Données de santé | Santé | ONI (décès) | Base légale spécifique |
| Données d'enquête | PNH | Justice | Transmission légale |

### 3.4 Access Delegation (Délégation d'Accès)

**Niveaux de délégation :**

| Niveau | Description | Approbation |
|--------|-------------|------------|
| Standard | Accès prédéfini par l'accord-cadre | Automatique (système) |
| Étendu | Accès au-delà de l'accord standard | Responsable SNISID agence |
| Exceptionnel | Accès d'urgence hors cadre normal | Directeur agence + BNC |
| Temporaire | Accès limité dans le temps pour un cas précis | Responsable SNISID + justificatif |

**Processus de délégation :**
```
1. Demande d'accès par l'agence requérante
2. Vérification de la base légale
3. Vérification du niveau d'habilitation de l'agent
4. Approbation selon le niveau de délégation
5. Configuration technique de l'accès
6. Journalisation de la délégation
7. Notification à l'agence propriétaire
8. Monitoring de l'utilisation
9. Revue périodique / expiration
```

### 3.5 Accountability (Responsabilité)

**Matrice de responsabilité :**

| Rôle | Responsabilité |
|------|---------------|
| Agence propriétaire | Qualité et exactitude des données source |
| Agence destinataire | Usage conforme, sécurité locale, non-rediffusion |
| BNC | Vérification de conformité, coordination |
| ANPD | Protection des droits des personnes |
| CNGN | Gouvernance globale, arbitrage |
| Agent individuel | Usage conforme, non-divulgation, traçabilité |

**Mécanismes de responsabilité :**

| Mécanisme | Description |
|-----------|-------------|
| Journalisation complète | Qui a accédé à quoi, quand, pourquoi |
| Audit trail croisé | Vérification côté émetteur et destinataire |
| Engagement écrit | Chaque agent signe une charte de responsabilité |
| Sanctions | Sanctions individuelles en cas de violation |
| Rapport de transparence | Publication périodique des statistiques d'échange |
| Droit d'audit | Chaque agence peut auditer l'usage de ses données |

---

## 4. PROCÉDURES INTER-AGENCES

### 4.1 Procédure de Demande d'Accès

```
Agence A soumet demande à Agence B via plateforme SNISID
    → Vérification automatique : accord-cadre applicable ?
        → Oui → Vérification habilitation agent
            → Habilitée → Accès accordé → Journalisation
            → Non habilitée → Refus → Notification
        → Non → Demande exceptionnelle
            → Justification + base légale
            → Approbation Directeur B + BNC
            → Si approuvée → Accès temporaire → Journalisation
            → Si refusée → Notification motivée → Recours CNGN
```

### 4.2 Procédure d'Incident Inter-Agences

```
Incident détecté (violation, accès non autorisé, fuite)
    → Notification immédiate aux agences concernées
    → Notification CERT-HT (sécurité) et ANPD (données)
    → Gel des accès si nécessaire
    → Investigation conjointe
    → Rapport d'incident
    → Mesures correctives
    → Revue de l'accord-cadre si nécessaire
    → Sanctions le cas échéant
```

### 4.3 Procédure de Résiliation d'Accord

```
Demande de résiliation (par une partie ou par le CNGN)
    → Préavis de 90 jours (sauf urgence sécuritaire)
    → Plan de désengagement
    → Destruction des données partagées chez le destinataire
    → Vérification de la destruction (audit)
    → Clôture de l'accord
    → Archivage de l'accord et de son historique
```

---

## 5. GOUVERNANCE DES ÉCHANGES

### 5.1 Comité des Échanges Inter-Agences
- Composition : 1 représentant par agence + BNC + ANPD (observateur)
- Fréquence : Trimestrielle
- Missions : Revue des accords, résolution de problèmes, évolution des besoins

### 5.2 Registre National des Échanges
Registre centralisé et public contenant :
- Liste des accords en vigueur
- Statistiques d'échange (volume, fréquence)
- Incidents liés aux échanges
- Résultats des audits

---

*Document cadre préparé dans le cadre de la Phase 14 — SNISID National Legal Framework*
