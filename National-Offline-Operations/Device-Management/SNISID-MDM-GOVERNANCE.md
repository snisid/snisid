---
# ============================================================
# SNISID-Edge — National Field Device Management (MDM)
# Sécurité des Terminaux Mobiles et Remote Wipe
# Document ID: SNISID-MDM-001
# Version: 1.0.0
# ============================================================

## 1. MOBILE DEVICE MANAGEMENT (MDM)

Le déploiement de milliers de tablettes gouvernementales dans la nature représente un risque majeur d'exfiltration de données s'il n'est pas strictement encadré.

## 2. GOUVERNANCE DES TERMINAUX

Toutes les tablettes sont enrolées (Android Enterprise Dedicated Device / Kiosk Mode).
- L'utilisateur (Policier/Agent ONI) ne peut pas installer d'applications, ni sortir de l'application métier SNISID.
- Le port USB est verrouillé (désactivé au niveau de l'OS) pour empêcher l'extraction de données via câble.

## 3. PROTOCOLE "TIME BOMB" & REMOTE WIPE

### 3.1 Remote Wipe (Destruction à Distance)
Si une tablette est déclarée volée ou perdue par un officier, l'administrateur MDM envoie une commande de `Remote Wipe`. Dès que la tablette attrape un signal (Wi-Fi/4G), elle efface cryptographiquement tout son contenu et se "brique" (Factory Reset verrouillé).

### 3.2 Time-Bomb (Bombe à Retardement Offline)
Que se passe-t-il si le voleur place la tablette dans une cage de Faraday ou coupe les antennes pour empêcher la réception de l'ordre de destruction ?
- **Politique Time-Bomb :** Si la tablette ne réussit pas un "Handshake" chiffré (mTLS) avec le serveur central (ou un Edge Node régional) pendant **7 jours consécutifs**, le logiciel MDM interne déclenche automatiquement l'effacement local des clés de chiffrement. Toutes les données biométriques et les caches locaux deviennent illisibles.

---
*Document ID: SNISID-MDM-001 | Approuvé par: CISO National*
