# National BPMN Standards & Conventions
## Document de Standardisation Technologique — SNISID v4.0

---

## 1. CONVENTIONS DE NOMMAGE (BPMN NAMING)

Afin d'éviter le chaos organisationnel et de garantir une lisibilité absolue par les outils d'observabilité automatique, les règles de nommage suivantes s'appliquent de manière stricte :

### 1.1 Identifiants de Processus (Process Definition Key)
Format : `snisid-[domaine]-[sous_domaine]-[nom_workflow]` (tout en minuscules, séparé par des tirets).
- *Exemple* : `snisid-civil-birth-standard` (Naissance Simple)
- *Exemple* : `snisid-justice-warrant-arrest` (Mandat d'arrêt)

### 1.2 Identifiants d'Activités (Task IDs)
Format : `[type]_[action]_[entite]` (camelCase ou snake_case homogène).
- *Tâche Utilisateur* : `user_validate_birth_payload`
- *Tâche Service (Script/Worker)* : `service_sign_identity_certificate`
- *Événement Message de Fin* : `end_event_identity_revoked`

---

## 2. STRATÉGIE DE VERSIONNEMENT DES WORKFLOWS (VERSIONING)

Tous les processus déployés dans la Workflow Factory doivent suivre le schéma **Semantic Versioning (SemVer 2.0.0)** : `MAJOR.MINOR.PATCH`

1. **MAJOR (Majeur)** : Changement structurel modifiant le modèle de données requis, supprimant des étapes de validation humaine critiques, ou altérant la rétrocompatibilité des jetons (tokens) en cours d'exécution. Les instances actives ne peuvent pas être migrées automatiquement ; elles doivent se terminer sur l'ancienne version ou être migrées via un script de migration de jetons dédié.
2. **MINOR (Mineur)** : Ajout d'étapes non-bloquantes, de notifications, d'indicateurs de performance (KPI) ou de règles de routage alternatives sans rupture de compatibilité. Migration automatique possible des instances à l'état "en attente".
3. **PATCH (Correctif)** : Correction d'une erreur de syntaxe, correction de fautes dans les notifications, ajustements mineurs des timers de SLA ou de labels. Migration transparente.

---

## 3. SCHÉMAS STANDARDISÉS POUR LES ÉVÉNEMENTS (EVENT SCHEMAS)

Tous les messages échangés sur le bus d'événements Kafka doivent respecter la spécification **CloudEvents (v1.0)**. 

### Schéma Universel d'Événement SNISID (JSON Schema)
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "SNISIDEvent",
  "type": "object",
  "properties": {
    "specversion": { "type": "string", "const": "1.0" },
    "id": { "type": "string", "format": "uuid" },
    "source": { "type": "string", "format": "uri" },
    "type": { "type": "string" },
    "subject": { "type": "string" },
    "time": { "type": "string", "format": "date-time" },
    "datacontenttype": { "type": "string", "const": "application/json" },
    "datacryptosignature": { "type": "string" },
    "data": {
      "type": "object",
      "properties": {
        "correlationId": { "type": "string", "format": "uuid" },
        "operatorId": { "type": "string" },
        "agency": { "type": "string" },
        "payload": { "type": "object" }
      },
      "required": ["correlationId", "operatorId", "agency", "payload"]
    }
  },
  "required": ["specversion", "id", "source", "type", "subject", "time", "datacontenttype", "datacryptosignature", "data"]
}
```

---

## 4. DÉFINITIONS ET NIVEAUX DE SLA (SLA DEFINITIONS)

Les accords de niveau de service (SLA) sont catégorisés par niveau de priorité administrative :

| Niveau | Désignation | Temps Limite de Validation (T1) | Temps Limite d'Escalade (T2) | Action d'Urgence (SLA Breach Action) |
| :--- | :--- | :--- | :--- | :--- |
| **P1** | **Urgence Nationale / Critique** | 30 Minutes | 1 Heure | Transfert hiérarchique immédiat au Directeur National de l'Agence + Notification SMS/Email de crise. |
| **P2** | **Haute Priorité** | 4 Heures | 8 Heures | Réaffectation automatique à la file d'attente d'escalade d'équipe (Niveau 2). |
| **P3** | **Standard** | 48 Heures | 72 Heures | Alerte visuelle orange dans le tableau de bord de l'agent + rappel email quotidien. |
| **P4** | **Faible Priorité / De fond** | 10 Jours | 15 Jours | Rappel hebdomadaire par email. |

---

## 5. RÈGLES DE CLASSIFICATION DE SÉCURITÉ (SECURITY CLASSIFICATION)

Chaque workflow et chaque champ de données manipulé au sein du système de Case Management possède un niveau de classification définissant ses contrôles de sécurité :

1. **PUBLIC (Niveau 1)** : Données ouvertes, statistiques anonymisées de naissances/décès. Aucune restriction d'accès.
2. **RESTREINT (Niveau 2)** : Accès restreint aux agents habilités d'une administration donnée. Exemple : Changement d'adresse simple.
3. **CONFIDENTIEL (Niveau 3)** : Données personnelles identifiables (PII) - Identité civile complète, informations biométriques. Requiert authentification multi-facteurs (MFA) et journalisation de chaque consultation.
4. **SECRET-DÉFENSE / JUDICIAIRE (Niveau 4)** : Informations d'investigation de police, casier judiciaire complet, mandats actifs DCPJ. Chiffrement asymétrique de bout en bout. Les journaux d'accès sont transmis immédiatement à l'auditeur national et cryptés par clé de sécurité étatique.
