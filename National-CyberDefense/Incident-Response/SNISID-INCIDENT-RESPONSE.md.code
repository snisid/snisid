---
# ============================================================
# SNISID-Cyber — National Incident Response Framework
# Playbooks de Crise et Coordination
# Document ID: SNISID-IR-001
# Version: 1.0.0
# ============================================================

## 1. INCIDENT RESPONSE (Réagir sous le feu)

Un plan de réponse aux incidents n'a de valeur que s'il est documenté et testé AVANT l'attaque. En cas de crise majeure, la panique est le pire ennemi.

## 2. PHASES DE RÉPONSE (Standard NIST 800-61)
1. **Préparation :** Entraînement des équipes (Phase 6, Etape 17).
2. **Détection & Analyse :** Identification via le SIEM, qualification de la gravité (P1 à P4).
3. **Confinement, Éradication & Récupération :**
   - *Confinement (Containment) :* Isoler le sous-réseau infecté, révoquer les accès VPN.
   - *Éradication :* Supprimer le malware, patcher la vulnérabilité (Zero-Day).
   - *Récupération :* Restaurer les données depuis les backups immuables (MinIO WORM).
4. **Post-Incident (Lessons Learned) :** Réunion obligatoire pour ajuster l'architecture (Comment sont-ils rentrés ? Comment éviter que cela ne se reproduise ?).

## 3. PLAYBOOK CRISIS : RANSOMWARE (Nation-State)

Si un groupe étatique chiffre l'infrastructure centrale :
- **Déclenchement "Code Rouge" :** Basculement du pouvoir décisionnel au "Crisis Commander" (Gouvernemental).
- **Communication Out-of-Band :** Utilisation de canaux de communication de secours (ex: Signal/Threema sur téléphones dédiés) car l'infrastructure interne (emails) est présumée compromise.
- **Politique de non-paiement :** L'État haïtien ne paie jamais de rançon. Restauration obligatoire depuis les Backups Offline et Tape (Phase 5).

---
*Document ID: SNISID-IR-001 | Approuvé par: Conseil National de Sécurité*
