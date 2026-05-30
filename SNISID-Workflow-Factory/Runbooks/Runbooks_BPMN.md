# Runbooks d'Opérations & de Secours BPMN
## Manuel de Résilience Opérationnelle — SNISID v4.0

---

## 1. RUNBOOK 01 : Récupération après Échec de Workflow (Workflow Failure Recovery)

### 1.1 Contexte & Symptômes
- Une instance de workflow est bloquée avec le statut `SYSTEM_ERROR` ou dans une boucle d'erreurs répétitives (Incident Camunda).
- Alertes d'échec de workflow reportées sur Prometheus (`snisid_workflow_failure_total`).
- Logs d'erreur dans Loki : `NullPointerException`, `ConnectorTimeoutException` ou `DataValidationException`.

### 1.2 Procédure de Résolution pas-à-pas
1. **Identifier l'Instance & l'Activité Bloquée** :
   Consulter le cockpit d'observabilité ou exécuter la commande d'extraction des incidents actifs :
   ```bash
   curl -X GET "https://bpmn.snisid.gouv/engine-rest/incident?failedActivityId=true" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"
   ```
2. **Analyser le Payload de Variables** :
   Déterminer si l'erreur provient de données corrompues transmises dans le dossier. S'il s'agit d'un problème de données (par exemple, un code de commune manquant dans le payload de naissance), corriger la variable de l'instance de dossier :
   ```bash
   curl -X PUT "https://bpmn.snisid.gouv/engine-rest/process-instance/<INSTANCE_ID>/variables/communeCode" \
     -H "Content-Type: application/json" \
     -d '{"value": "HT0111", "type": "String"}'
   ```
3. **Réessayer le Point de Défaillance (Retry Incident)** :
   Une fois la correction de donnée appliquée ou le microservice destinataire rétabli, déclencher le rejeu du token :
   ```bash
   curl -X POST "https://bpmn.snisid.gouv/engine-rest/incident/<INCIDENT_ID>/retry" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"
   ```
4. **Vérification** :
   Vérifier que le statut de l'instance passe de `SYSTEM_ERROR` à `RUNNING` ou `COMPLETED`.

---

## 2. RUNBOOK 02 : Escalade de Dépassement de SLA (SLA Breach Escalation)

### 2.1 Contexte & Symptômes
- Une tâche humaine dépasse sa limite légale d'examen.
- Alerte de priorité P1/P2 détectée par le `SLA Engine`.
- Alerte système générée dans le canal d'Urgence Nationale (`system.alarms`).

### 2.2 Procédure d'Escalade Automatique et Manuelle
1. **Détection par le Moteur** :
   Le `NationalSLAEngine` détecte que le timer d'une tâche `TASK-XYZ` a expiré.
2. **Vérification de la Priorité** :
   - **Si P1 (Critique)** : Déclencher immédiatement l'attribution au rôle supérieur `SUPERVISEUR_GENERAL_AGENCE` et alerter via la passerelle d'alerte gouvernementale.
   - **Si P2/P3 (Standard/Haute)** : Déclencher la délégation automatique vers la file d'attente d'escalade régionale.
3. **Réaffectation Manuelle Forcée (en cas d'obstruction systémique)** :
   Si aucun agent superviseur n'est disponible, l'administrateur système peut réattribuer la tâche directement à un pool d'agents de secours :
   ```bash
   curl -X POST "https://bpmn.snisid.gouv/engine-rest/task/<TASK_ID>/assignee" \
     -H "Content-Type: application/json" \
     -d '{"userId": "agent.backup.01@snisid.gouv"}'
   ```
4. **Journalisation d'Audit** :
   Chaque action d'escalade doit être scellée avec le tag `SLA_FORCE_ESCALATION` et transmise au journal immuable.

---

## 3. RUNBOOK 03 : Interruption et Restauration du Bus Kafka (Kafka Disruption Recovery)

### 3.1 Contexte & Symptômes
- Perte de liaison ou indisponibilité du cluster Kafka principal.
- Exceptions levées dans les services : `BrokerNotAvailableException` ou `TimeoutException` de production de messages.
- Files d'attente locales des microservices saturées.

### 3.2 Plan de Continuité d'Activité (PCA)
1. **Activation du Mode Dégradé local (Circuit Breaker)** :
   Dès que l'indicateur d'état du bus tombe à `DOWN`, les applications activent les disjoncteurs (Circuit Breakers) et basculent automatiquement l'écriture des événements vers une base de données de secours locale (Local Event Cache, ex: Redis/PostgreSQL local).
2. **Redémarrage / Redressement du Cluster Kafka** :
   Les administrateurs système déploient le script d'urgence de redressement de l'infrastructure de streaming :
   ```bash
   kubectl rollout restart statefulset/kafka -n snisid-infra
   ```
3. **Vérification de l'État du Cluster** :
   S'assurer que tous les courtiers (brokers) sont synchronisés et que le contrôleur est actif :
   ```bash
   kafka-topics.sh --bootstrap-server kafka-hl:9092 --list
   ```
4. **Récupération & Rejeu des Événements Stockés (Replay Phase)** :
   Une fois Kafka rétabli, déclencher le travailleur de synchronisation d'infrastructure pour consommer la base de sauvegarde Redis/PostgreSQL et renvoyer les messages dans les topics respectifs, dans l'ordre strict d'occurrence (First-In, First-Out).

---

## 4. RUNBOOK 04 : Résolution de Goulot d'Étranglement de Validation Humaine (Human Approval Bottleneck)

### 4.1 Contexte & Symptômes
- Accumulation excessive de dossiers en attente sur un rôle d'officier d'État Civil régional spécifique (`queue_depth > 100` dossiers).
- Les citoyens subissent des retards dans l'émission de certificats.

### 4.2 Mesures de Résolution des Goulots d'Étranglement
1. **Identifier la Région Affectée** :
   Identifier les files d'attente saturées :
   ```sql
   SELECT role_required, COUNT(*) as task_count 
   FROM human_tasks 
   WHERE status IN ('PENDING', 'ASSIGNED') 
   GROUP BY role_required 
   ORDER BY task_count DESC;
   ```
2. **Appliquer les Règles de Délégation Dynamique** :
   Si une région accuse un retard de plus de 50%, déléguer automatiquement une portion (ex: 30%) des dossiers les plus anciens vers une autre région moins sollicitée (Load Balancing de Validation Humaine).
3. **Appliquer une Réaffectation en Lot (Bulk Reassignment)** :
   Exécuter le script d'administration pour réattribuer 50 dossiers de la file d'attente d'un agent absent vers une file partagée :
   ```bash
   python3 -m snisid.bpm.admin_tools --bulk-reassign --from-agent agent.retard@snisid.gouv --to-pool pool.secours.national
   ```

---

## 5. RUNBOOK 05 : Résolution d'Anomalies de Synchronisation Hors-Ligne (Offline Sync Issue)

### 5.1 Contexte & Symptômes
- Un agent de terrain revient d'une mission en zone blanche, reconnecte son terminal mobile, mais la synchronisation des dossiers créés échoue.
- Logs du terminal : `SyncConflictError` ou `CryptographicSignatureMismatch`.

### 5.2 Résolution Manuelle des Conflits Récalcitrants
1. **Extraire le Journal d'Événements Hors-ligne du Terminal** :
   Télécharger le journal local chiffré depuis le terminal mobile de l'agent.
2. **Inspecter l'Événement Responsable du Blocage** :
   Vérifier si une signature de sécurité locale ne correspond plus en raison d'une clé d'agent expirée durant sa mission terrain.
3. **Forcer le Rejeu Sécurisé via l'Audit Tool** :
   Si la signature de l'agent est expirée mais valide historiquement, l'administrateur peut utiliser la clé d'habilitation système pour signer à nouveau l'enveloppe de données et forcer l'intégration :
   ```bash
   python3 -m snisid.offline.sync_tool --force-import --terminal-id TERM-9988 --override-signatures
   ```
4. **Vérifier l'alignement des états** :
   S'assurer que le dossier national créé hors-ligne est désormais visible et cohérent dans le registre centralisé avec l'état `SYNCHRONIZED`.
