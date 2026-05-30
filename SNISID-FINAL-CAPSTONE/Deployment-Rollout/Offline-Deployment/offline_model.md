# SNISID Offline-First Deployment Model
## Architecture et Fonctionnement dans les Zones de Faible Connectivité d'Haïti

---

## 1. Philosophie du Système Offline-First

En Haïti, l'infrastructure de télécommunication est sujette à des pannes régulières en raison de la topographie montagneuse, des intempéries cycloniques et de l'instabilité du réseau électrique. Une plateforme d'identification nationale qui exigerait une connexion Internet permanente pour fonctionner serait inutilisable sur plus de 60% du territoire en dehors de l'aire métropolitaine de Port-au-Prince.

Le SNISID est conçu selon le principe de **l'Offline-First**. Toutes les fonctionnalités critiques d'enrôlement et de vérification d'identité peuvent être effectuées localement sur un **Local Edge Node (LEN)**, de manière totalement déconnectée du Datacenter Central, et synchronisées ultérieurement dès qu'un canal réseau devient disponible.

```
       +-------------------------------------------------------+
       |                  LOCAL EDGE NODE (LEN)                |
       |  (Bénéficie d'un cache local crypté LUKS AES-256)     |
       +-------------------------------------------------------+
            |                                       |
            v (Offline Enrollment)                  v (Offline Verification)
[Nouveau Citoyen Enrôlé localement]        [Lecture et validation biométrique]
- ID temporaire généré                     - Recherche dans le cache local
- Données chiffrées en attente             - Signature numérique locale
            |                                       |
            +-------------------+-------------------+
                                |
                                v (Delayed Sync Protocol)
             +-------------------------------------+
             |  SNISID SECURE CENTRAL DATACENTER   |
             |       (Conflict Resolution)         |
             +-------------------------------------+
```

---

## 2. Le Local Edge Node (LEN) - Spécifications Matérielles et Logicielles

Chaque bureau de liaison communal (BLC) ou unité mobile du SNISID est équipé d'un serveur ultra-robuste et sécurisé appelé **Local Edge Node**.

### 2.1 Spécifications Critiques du LEN
*   **Sécurité Physique :** Boîtier blindé inviolable (tamper-evident) équipé de détecteurs d'ouverture physique (Intrusion Detection System). En cas de violation détectée, le serveur détruit instantanément ses clés de déchiffrement en mémoire RAM (Zeroisation).
*   **Sécurité Logique :** Système d'exploitation durci Linux (Alpine Linux) avec noyau de sécurité renforcé. Disques SSD configurés en RAID 1 chiffrés par mot de passe dynamique lié au module TPM 2.0 (Trusted Platform Module) matériel.
*   **Base de Données Locale :** Base SQLite ultra-rapide contenant un sous-ensemble chiffré des identités de la commune pour permettre la vérification d'identité locale instantanée (Cache Régional).

---

## 3. Protocole d'Enrôlement Hors-ligne (Offline Enrollment)

Lorsqu'un citoyen se présente dans un bureau de liaison hors-ligne :
1. **Saisie et Capture :** L'opérateur saisit les données démographiques et capture les données biométriques (empreintes et photo) sur sa station de travail reliée au LEN par réseau local sécurisé (Ethernet chiffré TLS 1.3).
2. **Attribution d'un Identifiant Temporaire (IUI-T) :** Le LEN génère un identifiant d'identité temporaire unique encodé localement, signé par la clé privée du LEN (ex: `HT-TEMP-LEN04-00291-2026`).
3. **Chiffrement du Dossier :** Le dossier complet est chiffré en AES-256-GCM avec une clé unique générée par transaction, puis placé dans la file d'attente d'exportation sécurisée (*Outbound Sync Queue*).

---

## 4. Protocole de Synchronisation Différée (Delayed Sync Protocol)

Dès qu'une connexion réseau (Starlink, 4G ou liaison filaire) est établie, ou lors de la transmission physique de la clé de synchronisation sécurisée par un agent de liaison certifié, le protocole de synchronisation se déclenche.

```
                  DELAYED SYNCHRONIZATION SEQUENCE
                  
[LEN (Local Node)]                                 [Central DC (Souverain)]
        |                                                    |
        | ------ 1. TLS 1.3 Handshake & Auth (mTLS) -------> |
        | <----- 2. Node Status Verified & Token ----------- |
        |                                                    |
        | ------ 3. Upload Signed Transaction Queue --------> |
        |                                                    |
        |                                           [Reconciliation Factory]
        |                                           - Biometric deduplication
        |                                           - Demographic verification
        |                                                    |
        | <----- 4. Sync Confirmation & Sync Log ----------- |
        |                                                    |
   [Purge Queue]
```

### 4.1 Logique de Résolution de Conflits de Synchronisation (Conflict Resolution Heuristics)

La base centrale utilise une logique déterministe pour résoudre les conflits entre les données locales modifiées hors-ligne et les données centrales :
*   **Conflit de type Écriture/Écriture :** Si un profil a été mis à jour simultanément en local et en central, le principe de la **Dernière Écriture Gagnante (Last-Write-Wins - LWW)** s'applique, basée sur l'horloge réseau synchronisée par NTP sécurisé du LEN, à condition que le score de signature cryptographique soit intègre.
*   **Conflit d'Identité Unique :** Si deux personnes distinctes ont été enrôlées hors-ligne avec la même empreinte biométrique, la base centrale rejette la seconde inscription, marque l'IUI-T comme conflictuel et place la transaction dans la file de quarantaine pour audit immédiat par la cellule de crise.
