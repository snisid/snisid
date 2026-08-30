# Module AEGIS/JUDAS - Système Anti-Corruption Interne

## 🎯 Objectif

Détecter et neutraliser toute tentative de contournement du système SNISID par un technicien, administrateur ou utilisateur interne.

## 🔍 Fonctionnalités

### 1. Surveillance Comportementale
- Analyse en temps réel des actions utilisateurs
- Détection des accès hors profil (RBAC + ABAC)
- Surveillance des tentatives d'accès aux fichiers sensibles
- Tracking des commandes exécutées avec privilèges élevés

### 2. Alerte Automatique "Traitor Watch"
En cas de comportement suspect :
```json
{
  "alerte_id": "JUDAS-2024-0001",
  "timestamp": "2024-01-15T14:32:45Z",
  "utilisateur": {
    "nom": "DUPONT",
    "prenom": "Jean",
    "photo_base64": "iVBORw0KGgoAAAANSUhEUgAA...",
    "nif": "123456789",
    "cin": "A1B2C3D4E",
    "niveau_acces": "Niveau 3 - Opérateur"
  },
  "action_tentee": {
    "type": "ACCES_NON_AUTORISE",
    "fichier_cible": "/data/biometrie/templates_iris.db",
    "commande_executee": "sudo cat /etc/shadow",
    "localisation": "Hub Regional Nord - Cap-Haïtien",
    "adresse_ip": "192.168.10.45"
  },
  "niveau_menace": "CRITIQUE",
  "action_requise": "VERROUILLAGE_IMMEDIAT"
}
```

### 3. Protocole d'Alerte
- **Fréquence** : Toutes les 15 minutes tant que la menace persiste
- **Destination** : Bureau Central SNISID (Port-au-Prince)
- **Canal** : WebSocket sécurisé + SMS + Email crypté
- **Affichage** : Dashboard central avec photo, nom, NIF en rouge clignotant

### 4. Rapport Automatique
Génération d'un rapport PDF signé cryptographiquement contenant :
- Chronologie complète des actions
- Captures d'écran des sessions suspectes
- Logs système corrélés
- Recommandations disciplinaires

### 5. Contre-Mesures Automatiques
- Verrouillage immédiat du compte utilisateur
- Révocation des tokens de session
- Notification à la hiérarchie directe
- Enregistrement dans la blockchain d'audit immuable

## 🔧 Installation

```bash
#!/bin/bash
# Installation du module AEGIS/JUDAS

echo "[*] Installation AEGIS/JUDAS..."

# Installation des dépendances
apt-get install -y auditd ossec-hids fail2ban python3-pip

# Configuration auditd
cat > /etc/audit/rules.d/snisid-judas.rules << EOF
# Surveillance des accès root
-a always,exit -F arch=b64 -S setuid -k priv_esc
-a always,exit -F arch=b64 -S setgid -k priv_esc

# Surveillance des fichiers sensibles
-w /etc/shadow -p wa -k shadow_access
-w /etc/passwd -p wa -k passwd_access
-w /data/biometrie/ -p rwa -k bio_access
-w /data/clefs_privees/ -p rwa -k key_access

# Surveillance des commandes critiques
-a always,exit -F arch=b64 -S execve -k all_exec
EOF

# Activation du service
systemctl enable auditd
systemctl restart auditd

echo "[+] AEGIS/JUDAS installé avec succès"
```

## 📊 Dashboard Central

L'administrateur central accède aux alertes via :
```bash
https://central.snisid.ht/judas/dashboard
Code d'activation requis : [CODE_6_CHIFFRES]
```

## 🔐 Sécurité du Module

- Le module JUDAS est lui-même protégé contre la désactivation
- Toute tentative de désactivation déclenche une alerte maximale
- Redondance sur 3 serveurs distincts
- Journalisation dans blockchain privée (4 nœuds)

## ⚖️ Cadre Légal

Conformément à l'Arrêté Présidentiel du 15 Janvier 2024 :
- Article 1 : Tout agent de l'État acceptant un poste SNISID consent à cette surveillance
- Article 2 : La tentative de contournement est punie de 10 ans de réclusion
- Article 3 : Les preuves générées par JUDAS sont recevables devant tous les tribunaux

---

**Classification :** TOP SECRET  
**© 2024 Gouvernement Haïtien**
