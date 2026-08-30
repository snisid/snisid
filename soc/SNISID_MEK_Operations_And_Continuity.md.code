# SNISID : Opérations de Terrain et Continuité d'Activité des MEK (v2.0)

**Classification :** SOUVERAIN / CRITIQUE  
**Recommandation de Référence :** SNISID v2.0 — MP-008

Ce document définit les protocoles opérationnels de télémétrie, de synchronisation et de sécurité d'urgence pour les **Kits Mobiles d'Enrôlement (MEK)** déployés sur le territoire de la République d'Haïti.

---

## 1. Cartographie en Temps Réel des MEK Actifs (Dashboard NOC)

Chaque MEK transmet à intervalle régulier (toutes les 15 minutes si connecté, ou empilé localement si déconnecté) une charge utile de télémétrie vers le Network Operations Center (NOC) du SNISID.

### Sujet NATS JetStream : `snisid.mek.telemetry.<device_id>`

### Format JSON du Payload de Télémétrie :

```json
{
  "device_id": "MEK-HT-042",
  "timestamp": "2026-05-24T21:24:00Z",
  "status": "ONLINE",
  "gps": {
    "latitude": 19.7589,
    "longitude": -72.2036,
    "altitude_meters": 120.5,
    "geofence_compliant": true
  },
  "power": {
    "battery_level_percent": 88,
    "battery_voltage": 13.2,
    "power_source": "SOLAR",
    "solar_input_watts": 145.2,
    "remaining_runtime_minutes": 2160
  },
  "queue": {
    "pending_sync_enrollments": 42,
    "oldest_unsynced_timestamp": "2026-05-24T12:00:00Z",
    "storage_used_bytes": 85899345
  },
  "integrity": {
    "tpm_attestation_passed": true,
    "chassis_secured": true,
    "agent_session_active": true,
    "active_agent_id": "ONI-AGENT-8809"
  }
}
```

### Règles d'alerting du NOC (SOC Trigger) :
* **Alerte de Déconnexion (Priorité Haute) :** Si un MEK n'a pas publié de télémétrie depuis **> 48 heures**, une alerte de Niveau 2 est créée au SOC, notifiant le responsable de zone terrain pour enquête physique.
* **Alerte Geofence (Priorité Critique) :** Si un MEK sort de sa commune assignée de plus de 5 km (détection par calcul de distance haversine sur le GPS), le SOC lance la procédure d'enquête et verrouille temporairement la session de l'agent.

---

## 2. Zéroisation d'Urgence par Signal Distant (Remote Wipe)

En cas de vol confirmé ou de capture d'un MEK par des groupes hostiles, le CISO ou le Directeur Général de l'ONI peut émettre un ordre de destruction cryptographique (zéroisation) via le réseau NATS.

### 2.1. Format du Message d'Ordre de Zéroisation
Le message de commande est signé à l'aide de la clé ECDSA privée du CISO (secp256r1).

### Sujet NATS de Commande : `snisid.mek.commands.MEK-HT-042`

```json
{
  "command": "ZEROIZE",
  "target_device_id": "MEK-HT-042",
  "command_id": "CMD-20260524-9981",
  "timestamp": "2026-05-24T21:24:15Z",
  "reason": "CONFIRMED_THEFT_HAITI_SUD",
  "issued_by": "CISO-SNISID-001",
  "signature": "3045022100a89d7d24ab8c89528fa2d46e9df07eb82e185c90d6e87f8976bcf65103a890db022005e839e99c8bfd9a8c7df9a0f0d2c6e61f21a8d052bc76bf3389ac3544d6e902"
}
```

### 2.2. Cinématique Locale du Démon de Sécurité
Lors de la réception de ce message valide :
1. **Validation de Signature :** Le démon de sécurité local du MEK (qui tourne sur le NUC Edge) extrait et valide la signature cryptographique par rapport à la clé publique publique du CISO stockée dans le TPM 2.0.
2. **Purge des Clés LUKS (Cryptsetup) :** Il détruit immédiatement les slots de clés LUKS en mémoire vive et écrase les slots de métadonnées LUKS du NVMe à l'aide de données aléatoires :
   ```bash
   # Éliminer la clé de déchiffrement de la RAM
   cryptsetup luksSuspend /dev/nvme0n1p3
   # Écraser les entêtes de clés physiques de l'enveloppe LUKS
   dd if=/dev/urandom of=/dev/nvme0n1p3 bs=512 count=4096 conv=fdatasync
   ```
3. **Purge RAM :** Déclenche un crash noyau forcé par SysRq pour purger instantanément la RAM avant extraction d'éventuels résidus thermiques :
   ```bash
   echo c > /proc/sysrq-trigger
   ```

### 2.3. Protection Physique (Tamper Auto-Wipe via TPM 2.0)
Si le boîtier physique du MEK est ouvert sans clé d'autorisation ou si le NVMe est extrait :
* Les commutateurs d'intrusion du châssis modifient les registres PCR (Platform Configuration Registers) du TPM 2.0.
* Au démarrage suivant, le TPM constate le changement des PCR 4, 7 et 14 et refuse d'unsealer la clé principale de déchiffrement LUKS.
* Sans intervention manuelle combinée du SRE Lead et du CISO via une clé d'administration USB signée, les données restent chiffrées en AES-256-XTS et totalement inaccessibles.

---

## 3. Cadre Légal Haïtien pour la Zéroisation d'Urgence

L'effacement à distance de données citoyennes et étatiques est une mesure de souveraineté nationale de dernier recours. Elle s'inscrit dans le cadre juridique suivant en République d'Haïti :

### 3.1. Fondements Juridiques
* **Loi sur la Cybercriminalité (2017) :** L'article 42 autorise l'État haïtien à prendre toutes les mesures technologiques nécessaires pour empêcher l'accès non autorisé à des systèmes de données classés "Infrastructures Critiques de l'État".
* **Décret sur la Protection des Données Personnelles :** Conforme au principe de souveraineté numérique, stipulant que la perte de contrôle physique d'un support de stockage contenant des données biométriques impose son invalidation immédiate pour protéger l'identité des citoyens.
* **Décret d'Urgence Nationale :** En cas de catastrophe naturelle (Séisme, Cyclone Cat 5) ou d'invasion de zone, le Directeur Général de l'ONI dispose des pleins pouvoirs pour ordonner la mise hors service à distance des terminaux compromis.

### 3.2. Protocole de Consignation Légal
Avant toute émission de commande de zéroisation :
1. **Rapport d'Incident SOC :** Le SOC documente la preuve de perte/vol (perte de télémétrie > 48h, hors-zone GPS, ou rapport de l'agent de terrain).
2. **Signature du Formulaire de Zéroisation :** Le Directeur de l'ONI et le Directeur Central de la Police Judiciaire (DCPJ) doivent contresigner le formulaire d'autorisation légale.
3. **Consignation dans le Registre WORM :** La signature CISO et l'autorisation administrative sont horodatées et inscrites de façon immuable dans les logs WORM centraux.

---

*Ce protocole garantit l'intégrité absolue de l'identité haïtienne en interdisant toute compromission de données sur le terrain.*
