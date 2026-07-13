#!/bin/bash
#===============================================================================
# SNISID - CENTRE DE SÉCURITÉ (SOC) 24/7
# Standard: NSA SOC Framework + Chine MLPS 2.0 + NIST CSF
# Classification: TOP SECRET / NOFORN - Gouvernement Haïtien
#===============================================================================

set -euo pipefail

# Configuration
SOC_DIR="/workspace/snisid-5-pillars/soc-247"
ORGANIZATION_DIR="${SOC_DIR}/organization"
PROCEDURES_DIR="${SOC_DIR}/procedures"
TOOLS_DIR="${SOC_DIR}/tools"
TRAINING_DIR="${SOC_DIR}/training"
WAR_ROOM_DIR="${SOC_DIR}/war-room"

# Couleurs
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

log() {
    echo -e "${CYAN}[$(date '+%Y-%m-%d %H:%M:%S')]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[$(date '+%Y-%m-%d %H:%M:%S')] ✓${NC} $1"
}

# Création de la structure
setup_directories() {
    log "Création de la structure SOC 24/7..."
    mkdir -p "${ORGANIZATION_DIR}" "${PROCEDURES_DIR}" "${TOOLS_DIR}" "${TRAINING_DIR}" "${WAR_ROOM_DIR}"
    log_success "Structure SOC créée"
}

# Organisation du SOC
create_organization_structure() {
    log "Génération de la structure organisationnelle..."
    
    cat > "${ORGANIZATION_DIR}/soc_organization.md.code" << 'EOF'
# ORGANISATION DU SOC SNISID - 24/7/365

## Mission du SOC

**Mission**: Détecter, analyser et répondre aux incidents de sécurité affectant le système SNISID en temps réel, 24 heures sur 24, 7 jours sur 7, 365 jours par an.

**Objectifs**:
- Détection des menaces en < 5 minutes
- Qualification des incidents en < 15 minutes
- Endiguement des attaques critiques en < 30 minutes
- Éradication complète en < 4 heures
- Retour d'expérience systématique sous 48 heures

## Structure Hiérarchique

### Niveau 1: Direction du SOC

```yaml
soc_director:
  title: "Directeur du Centre de Sécurité Opérationnelle"
  grade: "A-1 (Cadre Dirigeant)"
  reports_to: "Directeur Général SNISID"
  
  responsibilities:
    - "Définition de la stratégie de sécurité"
    - "Budget et ressources du SOC"
    - "Relations avec les autorités nationales"
    - "Communication de crise"
    - "Validation des procédures majeures"
    
  requirements:
    - "Master en cybersécurité ou équivalent"
    - "15+ ans expérience sécurité nationale"
    - "Certifications: CISSP, CISM, GIAC GSLC"
    - "Habilitation Top Secret"
    - "Expérience direction SOC gouvernemental"
    
  shifts: "Horaires bureau + astreinte 24/7"
```

### Niveau 2: Management Opérationnel

```yaml
soc_manager_team:
  soc_deputy_director:
    count: 2  # Rotation jour/nuit
    title: "Directeur Adjoint SOC"
    shift_pattern: "12h jour / 12h nuit / 48h repos"
    
  team_leads:
    count: 4  # Un par équipe (Alpha, Bravo, Charlie, Delta)
    title: "Chef d'Équipe SOC"
    shift_pattern: "Rotation hebdomadaire"
    
  responsibilities:
    - "Supervision opérations en temps réel"
    - "Escalade niveau 2 et 3"
    - "Coordination avec autres agences"
    - "Validation rapports d'incidents"
    - "Planning et gestion des effectifs"
```

### Niveau 3: Analystes SOC

```yaml
analyst_tiers:
  tier1_analysts:
    count: 16  # 4 par shift x 4 équipes
    title: "Analyste SOC Niveau 1"
    role: "Triaging et qualification initiale"
    
    responsibilities:
      - "Surveillance continue des alertes SIEM"
      - "Qualification initiale (faux positif vs réel)"
      - "Enrichissement contexte incident"
      - "Escalade vers Tier 2 si nécessaire"
      - "Documentation première réponse"
      
    requirements:
      - "Licence informatique ou cybersécurité"
      - "2+ ans expérience SOC"
      - "Certifications: Security+, CEH, GCIA"
      - "Habilitation Secret"
      
    shift_pattern: "4x12 (4 jours 12h / 4 jours repos)"
    
  tier2_analysts:
    count: 8  # 2 par shift
    title: "Analyste SOC Niveau 2"
    role: "Investigation approfondie"
    
    responsibilities:
      - "Investigation incidents escaladés"
      - "Analyse malware forensique"
      - "Chasse aux menaces (threat hunting)"
      - "Corrélation multi-sources"
      - "Rédaction rapports techniques"
      
    requirements:
      - "Master cybersécurité"
      - "5+ ans expérience SOC/forensique"
      - "Certifications: GCIH, GNFA, GCFA"
      - "Habilitation Top Secret"
      
  tier3_experts:
    count: 4  # 1 par shift + on-call
    title: "Expert SOC Niveau 3"
    role: "Réponse avancée et reverse engineering"
    
    responsibilities:
      - "Réponse aux incidents critiques"
      - "Reverse engineering malware"
      - "Analyse APT (Advanced Persistent Threat)"
      - "Développement signatures détection"
      - "Mentorat Tier 1 et Tier 2"
      
    requirements:
      - "Master + spécialisation offensive"
      - "10+ ans expérience cybersécurité"
      - "Certifications: OSEP, OSCE, GXPN"
      - "Habilitation Top Secret SCI"
```

### Niveau 4: Spécialistes

```yaml
specialists:
  threat_intelligence:
    count: 4
    title: "Analyste Renseignement Menaces"
    role: "Veille stratégique et tactique"
    
    responsibilities:
      - "Collecte renseignements sources ouvertes/fermées"
      - "Analyse groupes adverses (APT)"
      - "Production bulletins menace"
      - "Intégration feeds TI dans SIEM"
      - "Liaison avec services renseignement"
      
  forensics_incident_response:
    count: 4
    title: "Expert Forensique & Réponse Incident"
    role: "Investigation post-incident"
    
    responsibilities:
      - "Acquisition preuves numériques"
      - "Analyse forensique disque/mémoire"
      - "Reconstruction timeline attaque"
      - "Préservation chaîne de custode"
      - "Support investigations judiciaires"
      
  grc_specialists:
    count: 2
    title: "Spécialiste Gouvernance Risque Conformité"
    role: "Conformité réglementaire"
    
    responsibilities:
      - "Veille réglementaire"
      - "Audits internes"
      - "Rapports conformité ANSSI/MLPS"
      - "Gestion risques cyber"
      - "Politiques et procédures"
```

## Effectifs Totaux

| Catégorie | Effectif | Shifts Couverts |
|-----------|----------|-----------------|
| Directeur SOC | 1 | Jour + astreinte |
| Directeurs Adjoints | 2 | 24/7 rotation |
| Chefs d'Équipe | 4 | Rotation hebdo |
| Analystes Tier 1 | 16 | 24/7 (4x12) |
| Analystes Tier 2 | 8 | 24/7 |
| Experts Tier 3 | 4 | 24/7 + garde |
| Analystes TI | 4 | Jour + astreinte |
| Experts Forensiques | 4 | Jour + astreinte |
| Spécialistes GRC | 2 | Jour |
| **TOTAL** | **45** | **24/7/365** |

## Planning Type (Shift 12h)

### Équipe Alpha (Jour J1)
- 06:00 - 18:00: 1 Directeur Adjoint + 1 Chef d'Équipe + 4 Tier 1 + 2 Tier 2 + 1 Tier 3
- Relève avec Équipe Bravo

### Équipe Bravo (Nuit J1)
- 18:00 - 06:00: 1 Directeur Adjoint + 1 Chef d'Équipe + 4 Tier 1 + 2 Tier 2 + 1 Tier 3 (on-call)
- Relève avec Équipe Charlie

### Rotation Complète
- Semaine 1: Alpha (Jour) → Bravo (Nuit)
- Semaine 2: Charlie (Jour) → Delta (Nuit)
- Semaine 3-4: Repos et formation

## Salaire et Avantages

```yaml
compensation:
  currency: "USD"
  
  salaries_annual:
    soc_director: "$180,000 - $250,000"
    deputy_director: "$140,000 - $180,000"
    team_lead: "$110,000 - $140,000"
    tier3_expert: "$120,000 - $160,000"
    tier2_analyst: "$90,000 - $120,000"
    tier1_analyst: "$65,000 - $90,000"
    specialist: "$95,000 - $130,000"
    
  benefits:
    - "Prime shift nuit: +20%"
    - "Prime week-end: +50%"
    - "Prime jours fériés: +100%"
    - "Assurance santé complète famille"
    - "Mutuelle retraite avantageuse"
    - "Véhicule de fonction (management)"
    - "Logement de fonction ou allocation"
    - "Formation continue illimitée"
    - "Certifications payées + bonus"
    
  retention_programs:
    - "Bonus fidélité annuel (5-15% salaire)"
    - "Plan carrière clair avec promotions"
    - "Participation conférences internationales"
    - "Programme reconnaissance employé du mois"
```

## Localisation Physique

### War Room SOC

```yaml
war_room_specifications:
  location: "Centre de Commandement SNISID, Port-au-Prince"
  size: "200 m²"
  capacity: "30 postes de travail"
  
  zones:
    operations_floor:
      description: "Espace principal analystes"
      workstations: 24
      displays: "Mur vidéo 4x4 écrans 55 pouces"
      
    management_area:
      description: "Bureaux vitrés supervision"
      offices: 6
      view: "Vue directe sur operations floor"
      
    conference_room:
      description: "Salle crise sécurisée"
      capacity: 20
      features:
        - "Visioconférence加密"
        - "Tableau blanc numérique"
        - "Phone conférence"
        - "Isolation phonique"
        
    break_room:
      description: "Espace détente équipiers"
      features:
        - "Kitchenette complète"
        - "Canapés et lits repos"
        - "Jeux et divertissements"
        
    secure_storage:
      description: "Armoires fortes preuves"
      security: "Biométrie + combinaison"
      
  technical_infrastructure:
    network:
      - "Réseau rouge (données classifiées)"
      - "Réseau noir (internet)"
      - "Air-gapped physique entre réseaux"
    power:
      - "UPS redondant 2N"
      - "Générateur backup automatique"
    environmental:
      - "Climatisation précise 20°C ±1"
      - "Contrôle humidité 45% ±5%"
```

---
*Document d'organisation officiel - Niveau Confidentiel*
*Direction Générale SNISID - République d'Haïti*
EOF

    log_success "Structure organisationnelle générée"
}

# Procédures Opérationnelles
create_sop_procedures() {
    log "Génération des procédures opérationnelles..."
    
    cat > "${PROCEDURES_DIR}/soc_sop.md.code" << 'EOF'
# PROCÉDURES OPÉRATIONNELLES STANDARD (SOP) - SOC SNISID

## SOP-001: Surveillance Continue SIEM

### Objectif
Assurer une surveillance ininterrompue de tous les événements de sécurité du système SNISID.

### Périmètre
- Tous les logs systèmes (serveurs, réseau, applications)
- Flux réseau (NetFlow, PCAP)
- Alertes IDS/IPS
- Événements HSM et cryptographiques
- Accès physiques (badges, biométrie)

### Procédure

#### 1. Démarrage de Shift
```bash
# Checklist pré-shift
- [ ] Vérifier état de santé SIEM (Wazuh/Elasticsearch)
- [ ] Confirmer connectivité toutes sources de logs
- [ ] Review alertes non traitées du shift précédent
- [ ] Vérifier bulletins menace récents
- [ ] Tester canal communication sécurisé
```

#### 2. Surveillance Active
- **Dashboard Principal**: Vue globale santé sécurité
- **Dashboard Secondaire**: Alertes temps réel par criticité
- **Revue Périodique**: Toutes les 15 minutes minimum

#### 3. Tri des Alertes
```
Priorité P1 (Critique): Traitement immédiat (< 5 min)
  - Intrusion confirmée
  - Exfiltration données en cours
  - Compromission HSM/clés
  - Attaque ransomware active

Priorité P2 (Élevée): Traitement urgent (< 15 min)
  - Suspicion intrusion
  - Mouvement latéral détecté
  - Escalade privilèges anormale
  - Communication C2 suspecte

Priorité P3 (Moyenne): Traitement standard (< 1 heure)
  - Tentatives scanning/reconnaissance
  - Violations politiques mineures
  - Anomalies comportementales

Priorité P4 (Faible): Traitement différé (< 4 heures)
  - Faux positifs confirmés
  - Activités légitimes inhabituelles
  - Informations renseignement
```

#### 4. Documentation
Chaque alerte doit être documentée dans le ticket avec:
- Timestamp exact
- Source et destination
- Signature/règle déclenchée
- Actions entreprises
- Résolution et classification finale

### Métriques de Performance
- Temps moyen détection (MTTD): < 5 minutes
- Temps moyen réponse (MTTR): < 30 minutes
- Taux faux positifs: < 20%
- Couverture sources logs: 100%

---

## SOP-002: Qualification et Escalade d'Incident

### Objectif
Qualifier rapidement la nature et sévérité d'un incident pour déterminer la réponse appropriée.

### Matrice de Qualification

```yaml
severity_matrix:
  critical_sev1:
    criteria:
      - "Compromission confirmée données Tier 0 (biométrie, clés)"
      - "Attaque en cours sur infrastructure critique"
      - "Exfiltration massive données sensibles"
      - "Déni de service affectant service national"
    response_time: "< 5 minutes"
    escalation: "Directeur SOC + DG SNISID + Ministre"
    actions:
      - "Activation cellule crise immédiate"
      - "Isolement segments affectés"
      - "Notification autorités nationales"
      - "Préparation communication publique"
      
  high_sev2:
    criteria:
      - "Suspicion forte compromission"
      - "Mouvement latéral détecté"
      - "Tentative accès non autorisé données sensibles"
      - "Malware confirmé sur poste critique"
    response_time: "< 15 minutes"
    escalation: "Chef d'équipe + Directeur Adjoint"
    actions:
      - "Investigation approfondie Tier 2/3"
      - "Containment partiel"
      - "Prélèvement preuves forensiques"
      
  medium_sev3:
    criteria:
      - "Tentative intrusion bloquée"
      - "Violation politique sécurité"
      - "Anomalie nécessitant investigation"
    response_time: "< 1 heure"
    escalation: "Tier 2 analyst"
    actions:
      - "Investigation standard"
      - "Documentation complète"
      - "Recommandations prévention"
      
  low_sev4:
    criteria:
      - "Faux positif confirmé"
      - "Activité inhabituelle mais légitime"
      - "Information renseignement"
    response_time: "< 4 heures"
    escalation: "Tier 1 analyst"
    actions:
      - "Clôture avec documentation"
      - "Ajustement règles détection si nécessaire"
```

### Formulaire d'Escalade

```markdown
FORMULAIRE D'ESCALADE D'INCIDENT - SNISID-SOC-FORM-002

Numéro Ticket: _______________
Date/Heure: _______________
Analyste: _______________

DESCRIPTION INCIDENT:
_____________________
_____________________

SEVERITÉ ÉVALUÉE: □ SEV1 □ SEV2 □ SEV3 □ SEV4

CRITÈRES RENCONTRÉS:
_____________________
_____________________

ACTIONS DÉJÀ ENTREPRISES:
_____________________
_____________________

PERSONNES ESCALADÉES:
Nom: _______________ Heure: _______________
Nom: _______________ Heure: _______________

RECOMMANDATIONS:
_____________________
_____________________

SIGNATURE ANALYSTE: _______________
```

---

## SOP-003: Réponse à Incident Critique (SEV1)

### Objectif
Contenir et éradiquer une attaque critique dans les délais les plus courts.

### Phase 1: Détection et Alarme (0-5 minutes)

```bash
# Actions immédiates Tier 1
1. Confirmer l'alerte comme SEV1
2. Activer alarme sonore War Room
3. Notifier Chef d'Équipe et Directeur Adjoint
4. Ouvrir ticket incident majeur
5. Démarrer chronomètre réponse
```

### Phase 2: Mobilisation (5-15 minutes)

```bash
# Actions management
1. Activer cellule crise (physique ou virtuelle)
2. Convoquer experts Tier 3 et forensiques
3.Notifier DG SNISID et autorités
4. Établir canal communication sécurisé
5. Briefing situation initial
```

### Phase 3: Investigation Rapide (15-30 minutes)

```bash
# Actions Tier 2/3
1. Identifier vecteur attaque initial
2. Cartographier étendue compromission
3. Identifier données potentiellement exposées
4. Déterminer si attaque toujours active
5. Documenter findings en temps réel
```

### Phase 4: Endiguement (30-60 minutes)

```bash
# Décisions containment
Options selon scénario:

Option A: Isolation réseau
  - Déconnecter segments affectés
  - Basculer traffic vers backups sains
  - Bloquer IPs/domaines malveillants
  
Option B: Shutdown contrôlé
  - Arrêt progressif systèmes compromis
  - Préservation mémoire pour forensique
  - Activation mode dégradé
  
Option C: Contre-mesures actives
  - Déploiement règles IPS bloquantes
  - Invalidation sessions compromises
  - Rotation clés cryptographiques
```

### Phase 5: Éradication (1-4 heures)

```bash
# Actions éradication
1. Supprimer accès attaquants (comptes, backdoors)
2. Nettoyer/reconstruire systèmes compromis
3. Patch vulnérabilités exploitées
4. Renforcer monitoring zones affectées
5. Vérifier absence persistance
```

### Phase 6: Recovery (4-24 heures)

```bash
# Retour à normale
1. Validation sécurité par audit rapide
2. Restauration services depuis backups sains
3. Monitoring renforcé post-recovery
4. Communication statut aux parties prenantes
5. Reprise activité normale progressive
```

### Phase 7: Post-Incident (24-48 heures)

```bash
# Lessons learned
1. Réunion retour expérience complète
2. Rédaction rapport incident détaillé
3. Identification améliorations processus
4. Mise à jour playbooks et détections
5. Plan actions correctives
```

---

## SOP-004: Chasse aux Menaces (Threat Hunting)

### Objectif
Proactivement rechercher des indicateurs de compromission non détectés par les outils automatiques.

### Cycle de Hunting

#### 1. Hypothèse de Départ
Sources d'hypothèses:
- Renseignement menaces (TI) récent
- Nouvelles TTPs observées dans secteur public
- Anomalies subtiles dans les logs
- Questions "Et si...?" de l'équipe

#### 2. Plan de Chasse

```yaml
hunt_plan_template:
  hypothesis: "Un attaquant pourrait exfiltrer des données via DNS tunneling"
  
  data_sources:
    - "Logs DNS (queries, volumes)"
    - "NetFlow sortant"
    - "Proxy logs"
    
  tools:
    - "Elasticsearch queries"
    - "Python scripts analyse"
    - "Wireshark pour PCAP"
    
  indicators_research:
    - "Domaines longueur anormale"
    - "Volume queries élevé par host"
    - "Entropie élevée noms domaine"
    - "Connections vers IPs suspectes"
    
  timeline: "4 heures maximum"
  analysts_assigned: "1 Tier 2 + 1 Tier 3"
```

#### 3. Exécution

```bash
# Exemple query Elasticsearch - DNS Tunneling detection
GET /logs-dns-*/_search
{
  "query": {
    "bool": {
      "filter": [
        {"range": {"timestamp": {"gte": "now-24h"}}},
        {"term": {"query_type": "TXT"}},
        {"script": {
          "script": {
            "source": "params._source.query_name.length() > 50"
          }
        }}
      ]
    }
  },
  "aggs": {
    "top_clients": {
      "terms": {"field": "client_ip", "size": 20}
    }
  }
}
```

#### 4. Analyse et Validation

- Examiner résultats anormaux
- Corréler avec autres sources
- Éliminer faux positifs connus
- Confirmer ou infirmer hypothèse

#### 5. Documentation et Suivi

Si menace confirmée:
- Ouvrir ticket incident
- Escalader selon SOP-002
- Créer nouvelles règles détection

Si négatif:
- Documenter hypothèse infirmée
- Archiver recherches
- Partager learnings avec équipe

### Programme de Hunting Proactif

```yaml
weekly_hunts:
  semaine_1: "Recherche mouvement latéral"
  semaine_2: "Détection persistence mechanisms"
  semaine_3: "Analyse exfiltration potentielle"
  semaine_4: "Chasse basée TI récente"
  
monthly_focus:
  janvier: "Credential theft & misuse"
  février: "Privilege escalation"
  mars: "Data exfiltration"
  avril: "Malware & C2 communications"
  # ... cycle annuel complet
```

---

## SOP-005: Gestion des Preuves Numériques

### Objectif
Assurer la préservation légale des preuves numériques pour poursuites judiciaires.

### Chaîne de Custode

#### 1. Collecte

```bash
# Procédure collecte preuves
1. Documenter scène (photos, vidéos)
2. Noter timestamp exact découverte
3. Identifier personne ayant découvert
4. Utiliser équipement stérilisé
5. Créer image bit-à-bit support original
6. Calculer hash MD5/SHA256 image
7. Sceller preuve dans sac anti-preuve
```

#### 2. Stockage

```yaml
evidence_storage:
  location: "Armoire forte SOC, accès biométrique"
  access_log: "Automatisé avec vidéo surveillance"
  environmental_controls:
    temperature: "18-22°C"
    humidity: "40-50%"
    magnetic_shielding: "Oui"
  inventory_system: "Base de données dédiée avec QR codes"
```

#### 3. Transport

```markdown
PROCÈS-VERBAL DE TRANSPORT DE PREUVES

Numéro affaire: _______________
Date transport: _______________

PREUVE(S):
Description: _______________
Numéro série/scellé: _______________
Hash vérifié: _______________

TRANSPORTÉ PAR:
Nom: _______________
Signature: _______________

REÇU PAR:
Nom: _______________
Organisation: _______________
Signature: _______________

HEURE DÉPART: _______________
HEURE ARRIVÉE: _______________

INCIDENTS PENDANT TRANSPORT: _____________________
```

#### 4. Analyse

- Travailler uniquement sur copies, jamais originaux
- Documenter chaque étape analyse
- Utiliser outils validés scientifiquement
- Préserver reproductibilité

#### 5. Présentation Judiciaire

- Préparer rapport expert compréhensible
- Visualisations claires (timeline, graphes)
- Être prêt à témoigner sous serment
- Maintenir impartialité scientifique

---

*Manuel de procédures SOC - Niveau Confidentiel*
*SNISID - République d'Haïti*
EOF

    log_success "Procédures opérationnelles générées"
}

# Outils et Technologies SOC
create_soc_tools_config() {
    log "Configuration des outils SOC..."
    
    cat > "${TOOLS_DIR}/soc_toolstack.md.code" << 'EOF'
# STACK TECHNOLOGIQUE SOC SNISID

## SIEM (Security Information and Event Management)

### Solution Principale: Wazuh + Elasticsearch Stack

```yaml
siem_architecture:
  wazuh_manager:
    instances: 3  # Cluster haute disponibilité
    eps_capacity: "50,000 événements/seconde"
    
  elasticsearch_cluster:
    master_nodes: 3
    data_nodes: 12
    ingest_nodes: 3
    total_storage: "2 PB (hot+warm+cold architecture)"
    retention:
      hot_storage: "7 jours (NVMe SSD)"
      warm_storage: "30 jours (SSD)"
      cold_storage: "365 jours (HDD)"
      
  kibana_dashboards:
    custom_dashboards: 25
    real_time_views: 15
    compliance_reports: 10
    
  log_sources_integrated:
    operating_systems:
      - "Linux (auditd, syslog)"
      - "Windows (Event Log, Sysmon)"
    network_devices:
      - "Firewalls (Palo Alto, Fortinet)"
      - "Switches (Cisco, Arista)"
      - "Load balancers"
    applications:
      - "Web servers (Apache, Nginx)"
      - "Databases (PostgreSQL, MongoDB)"
      - "SNISID microservices"
    security_tools:
      - "IDS/IPS (Suricata)"
      - "EDR (Wazuh agents)"
      - "HSM audit logs"
      - "Physical access systems"
```

### Règles de Détection Customisées

```yaml
custom_rules_snisid:
  # Détection accès anormal aux données biométriques
  - rule_id: "SNISID-001"
    description: "Accès massif templates biométriques"
    severity: "Critical"
    query: >
      event.dataset:biometric AND 
      event.action:query AND 
      biometric.records_accessed:>1000
      
  # Détection tentative contournement authentification MFA
  - rule_id: "SNISID-002"
    description: "Échecs MFA multiples même utilisateur"
    severity: "High"
    query: >
      event.category:authentication AND 
      event.outcome:failure AND 
      mfa.challenge_failed:>5
      
  # Détection export non autorisé données citoyens
  - rule_id: "SNISID-003"
    description: "Export données hors périmètre autorisé"
    severity: "Critical"
    query: >
      event.action:data_export AND 
      NOT destination.ip:(10.0.0.0/8 OR 172.16.0.0/12)
      
  # Détection activité HSM anormale
  - rule_id: "SNISID-004"
    description: "Opération cryptographique hors horaires"
    severity: "High"
    query: >
      event.dataset:hsm_audit AND 
      NOT @timestamp:[06:00 TO 22:00]
```

## SOAR (Security Orchestration, Automation and Response)

### Solution: Shuffle + Playbooks Custom

```yaml
soar_configuration:
  platform: "Shuffle Open Source"
  integrations:
    - "Wazuh (SIEM)"
    - "TheHive (Case Management)"
    - "MISP (Threat Intelligence)"
    - "Slack/Teams (Communication)"
    - "Email (Notifications)"
    - "Firewall APIs (Containment)"
    
  automated_playbooks:
    playbook_phishing_response:
      trigger: "Alerte phishing confirmé"
      steps:
        - "Isoler URL malveillante dans sandbox"
        - "Extraire IOCs (URLs, hashes, IPs)"
        - "Bloquer URLs dans proxy/firewall"
        - "Identifier utilisateurs ayant cliqué"
        - "Envoyer notification aux utilisateurs"
        - "Créer ticket TheHive automatiquement"
        - "Notifier équipe SOC"
      execution_time: "< 2 minutes"
      
    playbook_brute_force:
      trigger: "Détection brute force SSH/RDP"
      steps:
        - "Identifier IP source attaque"
        - "Compter échecs authentication"
        - "Si seuil >100: bloquer IP firewall"
        - "Vérifier si compte compromis"
        - "Forcer reset mot de passe si besoin"
        - "Documenter incident"
      execution_time: "< 5 minutes"
      
    playbook_malware_detection:
      trigger: "EDR détecte malware"
      steps:
        - "Isoler host réseau automatiquement"
        - "Capturer mémoire et processus"
        - "Soumettre échantillon sandbox"
        - "Rechercher IOC sur autres hosts"
        - "Créer ticket forensique"
        - "Notifier Tier 2/3"
      execution_time: "< 3 minutes"
```

## Threat Intelligence Platform

### Solution: MISP + OpenCTI

```yaml
threat_intelligence_stack:
  misp_instance:
    purpose: "Partage IOCs et analyses"
    feeds_subscribed:
      - "AlienVault OTX"
      - "Abuse.ch"
      - "CISA Known Exploited Vulnerabilities"
      - "MITRE ATT&CK"
      - "Feeds gouvernements alliés (partenariats)"
      
  opencti_instance:
    purpose: "Plateforme connaissance menaces"
    features:
      - "Graph visualization APT campaigns"
      - "Correlation automatisée IOCs"
      - "Production rapports TI"
      - "Integration SIEM/SOAR"
      
  internal_intelligence:
    sources:
      - "Incidents SNISID historiques"
      - "Analyses malware internes"
      - "Honeypots dédiés"
      - "Darkweb monitoring (sous-traitance)"
      
  dissemination:
    daily_bulletin: "Résumé menaces 24h"
    weekly_report: "Analyse tendances semaine"
    monthly_assessment: "Évaluation menace mensuelle"
    flash_alerts: "Urgent (SEV1 uniquement)"
```

## Network Security Monitoring

### Solution: Suricata + Zeek + Arkime

```yaml
nsm_stack:
  suricata_ids_ips:
    deployment: "TAP ports sur switches coeur"
    throughput: "100 Gbps agrégés"
    rulesets:
      - "Emerging Threats Pro"
      - "Snort VRT"
      - "Règles custom SNISID"
    actions:
      - "Alert logging"
      - "IPS blocking (mode learning puis actif)"
      - "File extraction pour analyse"
      
  zeek_network_analysis:
    purpose: "Analyse comportementale réseau"
    logs_generated:
      - "conn.log (connections TCP/UDP)"
      - "http.log (transactions HTTP)"
      - "dns.log (requêtes DNS)"
      - "ssl.log (handshakes TLS)"
      - "files.log (fichiers transférés)"
      
  arkime_packet_capture:
    purpose: "Capture et recherche PCAP"
    retention: "30 jours full PCAP"
    storage_required: "~500 TB"
    features:
      - "Recherche full-text dans packets"
      - "Extraction fichiers automatique"
      - "Session reconstruction"
      - "Integration Wazuh/Kibana"
```

## Endpoint Detection and Response (EDR)

### Solution: Wazuh Agents + Osquery

```yaml
edr_deployment:
  wazuh_agents:
    coverage: "100% des endpoints SNISID"
    estimated_agents: 2000
    capabilities:
      - "File integrity monitoring"
      - "Rootkit detection"
      - "Vulnerability detection"
      - "Active response (block processes)"
      - "Command execution remote"
      
  osquery_fleet:
    purpose: "Query SQL sur endpoints"
    scheduled_queries:
      - "Processus démarrages récents"
      - "Connections réseau établies"
      - "Fichiers modifiés System32"
      - "Tâches planifiées suspectes"
      - "Extensions navigateur installées"
      
  agent_hardening:
    tamper_protection: "Activé"
    uninstall_password: "Obligatoire"
    communication_encryption: "TLS 1.3 mutual"
    heartbeat_interval: "60 secondes"
```

## Case Management

### Solution: TheHive + Cortex

```yaml
case_management:
  thehive_instance:
    purpose: "Gestion incidents et investigations"
    features:
      - "Création tickets structurés"
      - "Assignment analystes"
      - "Timeline investigation collaborative"
      - "Attachement preuves"
      - "Génération rapports"
      
  cortex_analyzers:
    automated_analyzers: 50
    examples:
      - "VirusTotal file/IP/hash lookup"
      - "Whois domain registration"
      - "Geolocation IP addresses"
      - "Malware sandbox submission"
      - "Password breach checking"
      
  workflows:
    incident_lifecycle:
      - "New → In Progress → Pending → Resolved → Closed"
    sla_tracking:
      - "Alerte temps réponse par sévérité"
      - "Escalation automatique si dépassement"
```

## Infrastructure SOC

```yaml
soc_workstations:
  analyst_tier1_tier2:
    model: "Dell Precision 5820 Tower"
    specs:
      cpu: "Intel Xeon W-2255 (10 cores)"
      ram: "64 GB DDR4 ECC"
      gpu: "NVIDIA RTX A4000"
      storage: "1 TB NVMe SSD + 4 TB HDD"
      monitors: "3x 27 pouces 4K"
      os: "SNISID-OS hardened"
      
  analyst_tier3_forensics:
    model: "Dell Precision 7920 Tower"
    specs:
      cpu: "Dual Intel Xeon Gold 6248R (48 cores)"
      ram: "256 GB DDR4 ECC"
      gpu: "NVIDIA RTX A6000"
      storage: "2 TB NVMe SSD + 8 TB HDD + 16 TB DAS"
      monitors: "4x 32 pouces 4K"
      os: "SNISID-OS + outils forensiques"
      
  network_requirements:
    bandwidth_per_workstation: "1 Gbps garanti"
    latency_to_siem: "< 10 ms"
    vpn_backup: "Connexion satellite Starlink backup"
```

---
*Documentation technique SOC - Niveau Confidentiel*
*SNISID - République d'Haïti*
EOF

    log_success "Configuration outils SOC générée"
}

# Programme de Formation
create_training_program() {
    log "Création du programme de formation..."
    
    cat > "${TRAINING_DIR}/soc_training_curriculum.md.code" << 'EOF'
# PROGRAMME DE FORMATION SOC SNISID

## Curriculum Global

### Durée Totale: 12 Semaines (3 Mois)

## Phase 1: Fondamentaux (Semaines 1-2)

### Module 1.1: Introduction SNISID
- **Durée**: 2 jours
- **Contenu**:
  - Mission et organisation SNISID
  - Architecture technique globale
  - Classification données (Tier 0-3)
  - Politiques sécurité nationales
  - Cadre juridique haïtien
  
- **Validation**: Quiz écrit (80% requis)

### Module 1.2: Cybersecurity Fundamentals
- **Durée**: 3 jours
- **Contenu**:
  - Concepts réseau (TCP/IP, DNS, HTTP/S)
  - Cryptographie appliquée
  - Authentification et autorisation
  - Types de menaces et attaques
  - Kill Chain et MITRE ATT&CK
  
- **Validation**: Laboratoire pratique

### Module 1.3: Linux Hardening
- **Durée**: 3 jours
- **Contenu**:
  - Administration Linux avancée
  - SNISID-OS spécificités
  - Auditd et logging système
  - AppArmor/SELinux
  - Détection rootkits
  
- **Validation**: TP hardening serveur

### Module 1.4: Legal and Compliance
- **Durée**: 2 jours
- **Contenu**:
  - Lois protection données Haïti
  - Standards NSA et chinois applicables
  - Chaîne de custode preuves
  - Témoignage judiciaire
  - Éthique et déontologie SOC
  
- **Validation**: Étude de cas juridique

---

## Phase 2: Outils SOC (Semaines 3-5)

### Module 2.1: SIEM Wazuh
- **Durée**: 5 jours
- **Contenu**:
  - Architecture Wazuh
  - Déploiement et configuration agents
  - Création règles custom
  - Dashboard Kibana
  - Investigation requêtes Elasticsearch
  
- **Laboratoire**: 
  - Déployer cluster Wazuh
  - Créer 10 règles détection
  - Investiguer scenario attaque
  
- **Validation**: Certification Wazuh Level 1

### Module 2.2: Network Security Monitoring
- **Durée**: 4 jours
- **Contenu**:
  - Analyse trafic avec Suricata
  - Logs Zeek interprétation
  - Recherche PCAP Arkime
  - Détection anomalies réseau
  - Reconstruction sessions
  
- **Laboratoire**:
  - Analyser capture attaque réelle
  - Extraire IOCs malware
  - Identifier exfiltration DNS
  
- **Validation**: Exercice pratique NSM

### Module 2.3: EDR et Endpoint Analysis
- **Durée**: 4 jours
- **Contenu**:
  - Wazuh agents advanced features
  - Osquery queries
  - Analyse processus suspects
  - Persistence mechanisms
  - Memory forensics basics
  
- **Laboratoire**:
  - Détecter malware sur endpoint
  - Identifier persistence
  - Nettoyer compromission
  
- **Validation**: Scenario EDR completion

### Module 2.4: SOAR et Automation
- **Durée**: 2 jours
- **Contenu**:
  - Principes SOAR
  - Shuffle playbooks
  - Integration APIs
  - Automated response
  
- **Laboratoire**:
  - Créer playbook phishing
  - Tester automation containment
  
- **Validation**: Playbook fonctionnel

---

## Phase 3: Réponse Incident (Semaines 6-8)

### Module 3.1: Incident Response Methodology
- **Durée**: 3 jours
- **Contenu**:
  - Framework NIST SP 800-61
  - Phases IR (Preparation → Lessons Learned)
  - Communication crise
  - Documentation incidents
  - Coordination équipes
  
- **Validation**: Tabletop exercise

### Module 3.2: Digital Forensics
- **Durée**: 5 jours
- **Contenu**:
  - Acquisition preuves (disque, mémoire)
  - Outils forensiques (Autopsy, Volatility)
  - Timeline analysis
  - Artefacts Windows/Linux
  - Anti-forensics detection
  
- **Laboratoire**:
  - Imager disque compromis
  - Analyser dump mémoire
  - Reconstruire timeline attaque
  
- **Validation**: Rapport forensique complet

### Module 3.3: Malware Analysis
- **Durée**: 4 jours
- **Contenu**:
  - Static analysis (strings, PE structure)
  - Dynamic analysis (sandbox, debugging)
  - Reverse engineering basics
  - YARA rules creation
  - IOC extraction
  
- **Laboratoire**:
  - Analyser échantillons malware réels
  - Créer règles YARA
  - Documenter TTPs
  
- **Validation**: Rapport analyse malware

### Module 3.4: Threat Hunting
- **Durée**: 3 jours
- **Contenu**:
  - Méthodologie hunting
  - Hypothesis-driven hunting
  - Data exploration techniques
  - Hunting with ElasticSearch
  - Reporting findings
  
- **Laboratoire**:
  - Chasser menace simulée
  - Développer hypotheses
  - Présenter résultats
  
- **Validation**: Hunt report

---

## Phase 4: Spécialisation (Semaines 9-10)

### Tracks au Choix (selon rôle futur)

#### Track A: Analyste Tier 1
- **Focus**: Triaging, qualification, escalation
- **Modules additionnels**:
  - Social engineering detection
  - Phishing analysis
  - Basic log analysis
  - Customer communication
  
#### Track B: Analyste Tier 2
- **Focus**: Investigation approfondie
- **Modules additionnels**:
  - Advanced log correlation
  - Lateral movement detection
  - Privilege escalation analysis
  - Incident documentation
  
#### Track C: Expert Tier 3
- **Focus**: Réponse avancée, reverse engineering
- **Modules additionnels**:
  - Advanced malware RE
  - APT TTPs deep dive
  - Custom tool development
  - Mentorship skills

---

## Phase 5: Simulation et Certification (Semaines 11-12)

### Module 5.1: Capstone Exercise
- **Durée**: 5 jours
- **Scenario**: Attaque APT complète simulée
- **Phases**:
  - Jour 1: Détection initiale
  - Jour 2-3: Investigation et containment
  - Jour 4: Éradication et recovery
  - Jour 5: Post-incident review
  
- **Évaluation**:
  - Temps de détection
  - Qualité investigation
  - Efficacité réponse
  - Documentation produite
  - Travail d'équipe

### Module 5.2: Certifications Officielles
- **Examens**:
  - SNISID SOC Analyst Level 1 (obligatoire)
  - Security+ ou équivalent
  - Certification outil primaire (Wazuh)
  
- **Conditions obtention**:
  - 80% minimum chaque examen
  - Participation 100% formations
  - Validation capstone exercise

### Module 5.3: Graduation Ceremony
- **Événement**: Remise certificats
- **Invités**: Direction SNISID, Ministère
- **Discours**: Meilleur étudiant de la promotion

---

## Formation Continue

### Requirements Annuels
- **Heures formation**: 40 heures/an minimum
- **Certifications renewals**: Tous les 2 ans
- **Exercices simulation**: Trimestriels
- **Veille technologique**: Hebdomadaire (4h)

### Programme Avancé (Annuel)

```yaml
advanced_training_yearly:
  q1:
    - "Nouvelles menaces et TTPs"
    - "Mise à jour outils SOC"
    - "Retour expériences incidents année"
    
  q2:
    - "Formation spécialisée (au choix)"
    - "Conférence sécurité (externe)"
    - "Exercice inter-agences"
    
  q3:
    - "Deep dive technique (malware, forensics)"
    - "Certification avancée"
    - "Mentorat nouveaux analystes"
    
  q4:
    - "Préparation audits"
    - "Amélioration procédures"
    - "Planning formation année suivante"
```

## Budget Formation

| Poste | Coût par Analyste | Total (45 analystes) |
|-------|-------------------|----------------------|
| Formateurs externes | $10,000 | $450,000 |
| Plateformes e-learning | $2,000 | $90,000 |
| Certifications exams | $3,000 | $135,000 |
| Conférences externes | $5,000 | $225,000 |
| Laboratoires pratiques | $8,000 | $360,000 |
| Documentation et livres | $1,000 | $45,000 |
| **TOTAL INITIAL** | **$29,000** | **$1,305,000** |
| **TOTAL ANNUEL (continue)** | **$15,000** | **$675,000** |

---
*Programme de formation officiel - Niveau Confidentiel*
*SNISID - Académie de Cybersécurité*
EOF

    log_success "Programme de formation créé"
}

# War Room Setup
create_war_room_setup() {
    log "Configuration War Room..."
    
    cat > "${WAR_ROOM_DIR}/war_room_setup.md.code" << 'EOF'
# CONFIGURATION WAR ROOM SOC SNISID

## Spécifications Physiques

### Local Principal

```yaml
war_room_main:
  location: "Centre de Commandement SNISID, Niveau -1"
  surface: "200 m²"
  hauteur_sous_plafond: "3.5 m"
  capacite_max: "30 personnes"
  
  zones_amenagement:
    operations_floor:
      surface: "120 m²"
      postes_analystes: 24
      disposition: "Îlots de 4 postes"
      ergonomie:
        - "Bureaux assis-debout électriques"
        - "Chaises Herman Miller Aeron"
        - "Repose-pieds ergonomiques"
        - "Éclairage individuel LED"
        
    management_overview:
      surface: "30 m²"
      bureaux_vitres: 6
      position: "Surélevé 50cm, vue panoramique"
      equipements:
        - "Écrans supplémentaires supervision"
        - "Intercom vers operations floor"
        - "Boutons appel urgence"
        
    crisis_conference_room:
      surface: "40 m²"
      capacite: "20 personnes"
      equipements:
        - "Table ovale 6 mètres"
        - "Écrans tactiles 86 pouces x2"
        - "Système visioconférence Poly Studio X50"
        - "Telephones конференц Yealink CP960"
        - "Tableau blanc numérique Microsoft Surface Hub"
        - "Isolation phonique STC 55+"
        
    support_areas:
      break_room:
        surface: "20 m²"
        equipements:
          - "Kitchenette complète"
          - "Machine café professionnelle"
          - "Réfrigérateur et micro-ondes"
          - "Canapés et fauteuils"
          - "TV 65 pouces détente"
          
      secure_locker_room:
        surface: "15 m²"
        casiers_biometriques: 50
       用途: "Stockage effets personnels (téléphones interdits)"
```

## Infrastructure Technique

### Système d'Affichage

```yaml
video_wall:
  type: "Mur vidéo LCD bezel-less"
  configuration: "4x4 écrans 55 pouces"
  resolution_totale: "15360 x 8640 (8K+)"
  luminosite: "700 nits"
  duree_vie: "100,000 heures"
  
  zones_affichage:
    zone_1_4: "Dashboard SIEM temps réel"
    zone_5_8: "Carte géographique attaques"
    zone_9_12: "Alertes prioritaires"
    zone_13_16: "Métriques SOC et SLA"
    
  controleurs:
    modele: "Datapath FX4-670"
    entrees: "12 sources HDMI/DisplayPort"
    sorties: "16 écrans synchronisés"
    features:
      - "Picture-by-picture"
      - "Overlapping windows"
      - "Preset configurations"
      - "API control integration"
```

### Postes de Travail Analystes

```yaml
analyst_workstation:
  quantite: 24
  
  hardware:
    workstation:
      modele: "Dell Precision 5820"
      cpu: "Intel Xeon W-2255 (10c/20t, 3.7GHz)"
      ram: "64 GB DDR4-2933 ECC"
      gpu: "NVIDIA RTX A4000 16GB"
      stockage:
        system: "1 TB NVMe PCIe Gen4"
        data: "4 TB HDD 7200rpm"
      alimentation: "950W Platinum"
      
    moniteurs:
      configuration: "3x Dell UltraSharp U2723QE 27 pouces"
      resolution: "3840x2160 (4K) chacun"
      ergonomie: "Pivot, tilt, height adjustable"
      calibration: "Delta E < 2 pour accuracy"
      
    peripheriques:
      clavier: "Logitech MX Keys (sans fil chiffré)"
      souris: "Logitech MX Master 3"
      casque: "Jabra Evolve2 65 (ANC, Teams certified)"
      webcam: "Logitech Brio 4K (pour visios)"
      
  software:
    os: "SNISID-OS hardened"
    browsers: "Firefox ESR durci + Chrome enterprise"
    productivity: "LibreOffice, OnlyOffice"
    soc_tools:
      - "Wazuh Dashboard"
      - "TheHive"
      - "Shuffle SOAR"
      - "MISP"
      - "Arkime"
      - "Terminal SSH sécurisé"
```

### Réseau et Connectivité

```yaml
network_infrastructure:
  segmentation:
    reseau_rouge:
     用途: "Données classifiées SNISID"
      isolation: "Air-gap physique"
      acces: "Postes analysts uniquement"
      
    reseau_noir:
     用途: "Internet et威胁 intelligence"
      securite: "Firewall nouvelle génération"
      filtrage: "Proxy avec inspection SSL"
      
    reseau_bleu:
     用途: "Administration infrastructure"
      acces: "Réservé admins SOC"
      
  bande_passante:
    internet_dedie: "10 Gbps symétrique"
    liaison_sites_snisid: "10 Gbps fibre noire"
    backup_satellite: "Starlink Business (150 Mbps)"
    
  commutation:
    switches_acces: "Cisco Catalyst 9300 48 ports"
    uplinks: "2x 10GbE SFP+ par poste"
    qos: "Priorisation traffic SOC critique"
    
  wifi:
    ssid_analystes: "SNISID-SOC-SECURE (WPA3-Enterprise)"
    ssid_invites: "SNISID-GUEST (isolé, captif portal)"
    couverture: "100% war room + zones adjacentes"
```

### Alimentation Électrique

```yaml
power_systems:
  ups:
    type: "Onduleurs double conversion online"
    capacite: "60 kVA"
    autonomie: "2 heures à pleine charge"
    redondance: "2N (deux systèmes indépendants)"
    
  pdus:
    type: "PDUs intelligentes APC APDU"
    outlets: "24 par rack, monitorées individuellement"
    features:
      - "Mesure consommation temps réel"
      - "Alarmes seuils configurables"
      - "Reset distant par outlet"
      
  generateur_backup:
    capacite: "200 kW"
    carburant_autonomie: "72 heures"
    demarrage: "Automatique < 10 secondes"
    contrat_refuel: "Priorité nationale"
```

### Environnement et Climatisation

```yaml
environmental_control:
  climatisation:
    type: "Precision cooling In-row"
    capacite: "80 kW"
    redondance: "N+1"
    consignes:
      temperature: "20°C ± 1°C"
      humidite: "45% ± 5%"
      
  qualite_air:
    filtration: "HEPA H13"
    renouvellement: "6 volumes/heure"
    pression_differentielle: "+15 Pa (surpression)"
    
  acoustique:
    traitement_phonique: "Panneaux absorbants plafond/murs"
    niveau_bruit_cible: "< 45 dB(A)"
    isolation_exterieure: "STC 50+"
```

## Sécurité Physique War Room

```yaml
physical_security:
  controles_acces:
    niveau_1_entree_batiment:
      - "Gardes armés 24/7"
      - "Barrières véhicules anti-bélier"
      - "Caméras LPR (License Plate Recognition)"
      
    niveau_2_entree_soc:
      - "Badge RFID + code PIN"
      - "Scanner biométrique empreintes"
      - "Détecteur métaux"
      
    niveau_3_entree_war_room:
      - "Reconnaissance faciale"
      - "Scanner iris"
      - "Sabot anti-poursuite (mantrap)"
      - "Fouille aléatoire sacs"
      
  surveillance_video:
    cameras_interieures: "8x Axis Q6075-E PTZ 4K"
    couverture: "100% sans angles morts"
    enregistrement: "90 jours minimum"
    analytics:
      - "Détection intrusion"
      - "Comptage personnes"
      - "Objets abandonnés"
      - "Mouvements suspects"
      
  alarmes:
    intrusion: "Capteurs ouverture portes/fenêtres"
    incendie: "VESDA + détecteurs thermiques"
    panic_buttons: "Sous chaque poste analyste"
    silent_alarm: "Vers centre police nationale"
```

## Procédures Opérationnelles War Room

### Ouverture de Shift

```markdown
CHECKLIST OUVERTURE SHIFT SOC

Heure: ___________
Équipe: ___________
Chef d'équipe: ___________

INFRASTRUCTURE:
[ ] Vérifier état vidéo wall (tous écrans OK)
[ ] Tester système audio/intercom
[ ] Confirmer connectivité réseau (rouge/noir)
[ ] Vérifier état UPS (charge 100%)
[ ] Tester alarmes panic buttons

OUTILS SOC:
[ ] Connexion SIEM Wazuh (dashboard vert)
[ ] Vérifier flux logs (EPS normal)
[ ] Tester SOAR Shuffle (playbooks ready)
[ ] Ouvrir TheHive (tickets en cours)
[ ] Confirmer feeds Threat Intel (MISP à jour)

COMMUNICATION:
[ ] Tester visioconférence (appel test)
[ ] Vérifier téléphones conférence
[ ] Confirmer canal Slack/Teams sécurisé
[ ] Distribuer badges temporaires (si visiteurs)

HANDOVER:
[ ] Briefing avec équipe sortante (30 min)
[ ] Review incidents en cours
[ ] Review alertes prioritaires
[ ] Notes particulières shift précédent

SIGNATURES:
Chef entrant: _________________
Chef sortant: _________________
```

### Fermeture de Shift

```markdown
CHECKLIST FERMETURE SHIFT SOC

Heure: ___________
Équipe: ___________

DOCUMENTATION:
[ ] Mettre à jour tous tickets en cours
[ ] Compléter logs d'activité shift
[ ] Documenter incidents résolus
[ ] Préparer handover notes équipe suivante

NETTOYAGE:
[ ] Effacer tableaux blancs (sauf infos persistantes)
[ ] ranger documents confidentiels (armoires fermées)
[ ] Vider corbeilles (destructeur documents)
[ ] Nettoyer postes de travail

SÉCURITÉ:
[ ] Verrouiller armoires preuves
[ ] Déconnecter sessions (lock screens)
[ ] Activer alarmes périmétriques
[ ] Briefing équipe entrante complété

SIGNATURE:
Chef d'équipe sortant: _________________
```

## Budget War Room

| Catégorie | Montant (USD) |
|-----------|---------------|
| Construction et aménagement | $800,000 |
| Mur vidéo 4x4 + contrôleurs | $250,000 |
| Postes de travail (24x) | $360,000 |
| Infrastructure réseau | $200,000 |
| Système UPS et électrique | $300,000 |
| Climatisation precision | $250,000 |
| Sécurité physique | $150,000 |
| Ameublement ergonomique | $120,000 |
| Systèmes communication | $80,000 |
| Installation et configuration | $150,000 |
| **TOTAL** | **$2,660,000** |

---
*Spécifications War Room SOC - Niveau Confidentiel*
*SNISID - République d'Haïti*
EOF

    log_success "Configuration War Room générée"
}

# Main execution
main() {
    echo ""
    echo "==============================================================================="
    echo "              SNISID - CENTRE DE SÉCURITÉ (SOC) 24/7/365                       "
    echo "                               RÉPUBLIQUE D'HAÏTI                              "
    echo "==============================================================================="
    echo ""
    
    setup_directories
    create_organization_structure
    create_sop_procedures
    create_soc_tools_config
    create_training_program
    create_war_room_setup
    
    echo ""
    echo "==============================================================================="
    echo "                    SOC 24/7 CONFIGURÉ AVEC SUCCÈS                             "
    echo "==============================================================================="
    echo ""
    echo "Répertoire principal: ${SOC_DIR}"
    echo ""
    echo "Documents générés:"
    echo "  ✓ Structure organisationnelle (45 personnels)"
    echo "  ✓ Procédures opérationnelles standard (SOP)"
    echo "  ✓ Stack technologique SOC complet"
    echo "  ✓ Programme de formation 12 semaines"
    echo "  ✓ Configuration War Room détaillée"
    echo ""
    echo "Budget total estimé: ~$4.6M USD (hors salaires)"
    echo "Budget salaires annuels: ~$4.5M USD"
    echo ""
    echo "PROCHAINES ÉTAPES:"
    echo "  1. Recrutement et habilitations sécurité"
    echo "  2. Construction et aménagement War Room"
    echo "  3. Achat et installation équipements"
    echo "  4. Formation initiale équipes"
    echo "  5. Exercice de validation avant mise en production"
    echo ""
    echo "CLASSIFICATION: CONFIDENTIEL DÉFENSE"
    echo "==============================================================================="
}

# Exécution du script
main "$@"
