# SNISID Offline-First Operational Strategy
**Version** : 1.0 | **Date** : 2025-03-27 | **Statut** : ✅ Approuvé

## Principes
- Le système **fonctionne sans internet** dans les zones isolées.
- La synchronisation est **asynchrone** et tolère des pannes réseau longues (> 7 jours).

## Composants offline

### Edge computing régional
- Serveur local (mini-PC, Raspberry Pi 4 ou NUC) dans chaque commune
- Stockage local de la base des citoyens (uniquement les données nécessaires)
- Capacité : ~50 000 enregistrements locaux

### Offline enrollment
- Agent mobile (application Flutter) : capture biométrique, scan documents
- Données stockées en local, cryptées (AES-256)
- Transfert vers edge node via WiFi ou clé USB quotidienne

### Delayed sync
- File d’attente FIFO, horodatage UTC
- Résolution de conflits : règle « dernier événement gagnant » (avec audit)

### Satellite backup (lien d’urgence)
- Iridium / Starlink pour sync critique (flag sécurité)
- Bande passante limitée, priorisation des événements

### Solar operations
- Tous les edge nodes alimentés par panneaux solaires + batterie LiFePO4
- Autonomie : 3 jours sans soleil

### Mobile units
- Véhicule 4x4 avec équipement complet d’enrôlement
- Serveur embarqué, satellite, panneaux solaires

### Offline validation
- Génération de QR code signé avec vérification cryptographique locale
- Attestation PDF signée (clé privée locale) – peut être vérifiée hors ligne
