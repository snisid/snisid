---
# ============================================================
# SNISID-Cyber — National Endpoint Security Platform
# EDR, Device Compliance et Isolation
# Document ID: SNISID-ENDPOINT-001
# Version: 1.0.0
# ============================================================

## 1. PROTECTION DES POSTES DE TRAVAIL (Endpoints)

L'humain (le fonctionnaire) est le maillon faible (Phishing, clé USB infectée). L'Endpoint Detection and Response (EDR) est la dernière ligne de défense sur l'ordinateur lui-même.

## 2. EDR (Endpoint Detection and Response)

L'agent EDR (ex: Wazuh Agent / CrowdStrike) est installé sur 100% des postes de l'État (Windows/Linux/macOS) et des serveurs.
Contrairement à un antivirus classique basé sur des signatures connues, l'EDR analyse le comportement (Behavioral Analysis).
- **Détection :** Si Word (`winword.exe`) essaie soudainement d'exécuter PowerShell et de télécharger un fichier `.exe` chiffré depuis une adresse IP russe, l'EDR bloque l'action instantanément (Runtime Protection).

## 3. DEVICE COMPLIANCE (Gouvernance des Terminaux)

L'accès au réseau SNISID (via le NAC de la Phase 5 ou Zero Trust) est conditionné à la conformité (Compliance) du poste.
Règles de conformité obligatoires :
1. Disque dur entièrement chiffré (BitLocker / LUKS).
2. Agent EDR actif et communicant avec le SOC.
3. Système d'exploitation mis à jour (Patching < 15 jours).
4. Aucun logiciel non autorisé (AppLocker) installé.

Si un de ces critères n'est pas respecté (ex: l'utilisateur a désactivé l'EDR), le poste est placé en quarantaine réseau dynamique (VLAN Quarantaine) et ne peut plus accéder aux applications métier.

---
*Document ID: SNISID-ENDPOINT-001 | Approuvé par: SOC Manager*
