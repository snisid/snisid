---
# ============================================================
# SNISID-Cyber — National SOAR Platform
# Automatisation de la Réponse aux Incidents (Playbooks)
# Document ID: SNISID-SOAR-001
# Version: 1.0.0
# ============================================================

## 1. AUTOMATISATION (Machine Speed Defense)

Le SOAR (Security Orchestration, Automation, and Response) exécute des actions à la vitesse de la machine pour contrer des attaques automatisées (ex: Ransomware, Brute-Force). Un analyste humain mettrait 10 minutes à bloquer une IP ; le SOAR le fait en 2 secondes.

## 2. PLAYBOOKS DE RÉPONSE

Les Playbooks sont des workflows automatisés (scripts Python / n8n / TheHive Cortex).

### 2.1 Playbook: Ransomware Containment (Endpoint)
**Déclencheur :** L'EDR (Wazuh) détecte un processus chiffrant massivement des fichiers locaux sur l'ordinateur d'un greffier de Justice.
**Actions SOAR :**
1. **Isolation (T+1s) :** Envoie une commande API à l'agent EDR pour isoler l'ordinateur du réseau (Coupure de la carte réseau virtuelle). Seul le port de l'EDR reste ouvert.
2. **Snapshot (T+5s) :** Demande à l'hyperviseur de prendre un snapshot de la RAM de la VM (Forensic Evidence).
3. **IAM Revoke (T+10s) :** Désactive le compte Active Directory/Keycloak de l'utilisateur compromis.
4. **Notification (T+15s) :** Ouvre un ticket P1 dans le JIRA du SOC et envoie un message sur le canal sécurisé (Mattermost) de l'équipe de crise.

### 2.2 Playbook: Phishing Escalation
**Déclencheur :** Un employé signale un email suspect.
**Actions SOAR :**
1. Extraction des URLs et des pièces jointes de l'email.
2. Soumission de la pièce jointe à une "Sandbox" Malware (ex: Cuckoo) pour analyse comportementale.
3. Requête aux bases de Threat Intelligence (VirusTotal / MISP) pour la réputation de l'IP/URL.
4. Si malveillant : Purge (Hard Delete) de cet email dans TOUTES les boîtes aux lettres gouvernementales.

---
*Document ID: SNISID-SOAR-001 | Approuvé par: SOC Manager*
