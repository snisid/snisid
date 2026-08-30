---
# ============================================================
# SNISID-Cyber — Privileged Access Management (PAM)
# Gouvernance des Super Admins (Teleport / Vault)
# Document ID: SNISID-PAM-001
# Version: 1.0.0
# ============================================================

## 1. LE PROBLÈME DES ACCÈS PRIVILÉGIÉS

Les cyberattaques les plus dévastatrices (ou la corruption interne) impliquent des comptes "Root" ou "Administrator". Un administrateur système corrompu pourrait potentiellement effacer ou modifier des identités (Insiders Threats).

## 2. ARCHITECTURE PAM (TELEPORT & VAULT)

Pour contrecarrer cela, **aucun humain ne connaît le mot de passe Root d'aucun serveur ni base de données.**

### 2.1 Just-In-Time Access (JIT)
Si un ingénieur SRE doit se connecter en SSH sur un serveur de production :
1. Il se connecte au portail Teleport via SSO + YubiKey (MFA Matériel).
2. Il demande un accès temporaire (Access Request) justifié par un ticket JIRA.
3. Un autre ingénieur ou le SOC approuve la demande (Dual-Control).
4. Teleport génère un certificat éphémère valable **15 minutes**. L'ingénieur est connecté.

### 2.2 Session Recording (Audit Visuel)
Toute frappe clavier (Keystrokes) et l'enregistrement vidéo du terminal (TTY) de la session SSH de l'ingénieur sont enregistrés de manière immuable et streamés en direct vers le SOC.

### 2.3 HashiCorp Vault (Gestion des Secrets)
Les applications ne stockent pas de mots de passe en dur.
L'API SNISID contacte Vault au démarrage pour obtenir un mot de passe temporaire pour CockroachDB (Database Secrets Engine), qui change toutes les heures. Si une base de données fuit, le mot de passe est déjà obsolète.

---
*Document ID: SNISID-PAM-001 | Approuvé par: CISO National*
