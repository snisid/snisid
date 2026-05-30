---
# ============================================================
# SNISID-Edge — National Offline Identity Verification
# Vérification Cryptographique par QR Code
# Document ID: SNISID-OFFLINE-ID-001
# Version: 1.0.0
# ============================================================

## 1. LE PROBLÈME DE LA VÉRIFICATION OFFLINE

Comment un policier peut-il vérifier qu'une carte d'identité n'est pas un faux s'il n'a pas accès à la base de données centrale pour interroger le numéro d'identité (NIU) ?

## 2. LE QR CODE SÉCURISÉ (Visible Digital Seal - VDS)

Chaque document SNISID (Carte d'identité, Permis, Passeport) inclut un QR code haute densité (CEV - Cachet Électronique Visible).

### 2.1 Contenu du QR Code
- Nom, Prénom, Date de naissance, NIU.
- Un Hash (empreinte courte) de la photo du citoyen.
- **La Signature Numérique (ECDSA)** émise par la PKI Souveraine du SNISID.

### 2.2 Processus de vérification (Offline)
1. L'application mobile (PNH) lit le QR code.
2. Elle compare la signature numérique contenue dans le QR code avec le certificat public de l'État Haïtien préchargé dans l'application.
3. Si la signature est valide, l'application certifie que les données du QR code ont bien été émises par l'État et n'ont pas été altérées.
4. Le policier vérifie visuellement que le visage de la personne correspond à la photo imprimée (ou recalcule le hash de la photo lue sur la puce NFC).

**Prévention Fraude :** Un faussaire peut imprimer une fausse carte, mais il ne peut pas générer un faux QR code valide sans voler la clé privée (HSM) de la PKI nationale (Zero Trust).

---
*Document ID: SNISID-OFFLINE-ID-001 | Approuvé par: Direction de l'Identité Nationale*
