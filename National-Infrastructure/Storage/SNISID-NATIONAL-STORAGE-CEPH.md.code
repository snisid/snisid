---
# ============================================================
# SNISID-Infra — Sovereign Storage & Immutable Backups
# Ceph, MinIO et Rétention Légale
# Document ID: SNISID-STORAGE-001
# Version: 1.0.0
# ============================================================

## 1. ARCHITECTURE DE STOCKAGE DISTRIBUÉ (CEPH)

Pour garantir qu'aucun disque dur défectueux ne cause de perte de données, tout le stockage est géré par **Ceph**. Ceph regroupe les disques durs physiques de tous les serveurs en un seul "pool" massif et réplique la donnée.

### 1.1 Réplication Triplée (3x Replication)
Chaque bloc de données (ex: une ligne de la base CockroachDB) est écrit sur 3 disques durs physiques situés dans 3 baies (Racks) différentes.
- Tolérance de panne : 2 disques complets ou 2 serveurs complets peuvent brûler simultanément sans AUCUNE perte de donnée.

## 2. IMMUTABLE BACKUPS (Protection Ransomware)

Face à la menace croissante des ransomwares (Cryptolockers), les sauvegardes (Backups) utilisent la fonctionnalité "Object Lock" (WORM - Write Once Read Many) de **MinIO**.

### 2.1 Cycle de Vie de la Sauvegarde (Backup Lifecycle)
1. Velero (sur Kubernetes) prend un snapshot des bases de données toutes les heures.
2. Le snapshot est chiffré (AES-256 GCM) puis envoyé vers le cluster MinIO.
3. **Immutabilité :** MinIO verrouille l'objet pour 30 jours (Compliance Mode).
4. **Scénario Ransomware :** Si un hacker (ou un admin corrompu) accède au système et lance la commande `aws s3 rm --recursive`, le cluster MinIO *refusera* l'effacement. Les données sont garanties inaltérables pendant 30 jours.

## 3. ARCHIVES À LONG TERME (Tape Backup / Cold Storage)

Pour les exigences légales (ex: conservation de l'état civil à perpétuité), des sauvegardes "Cold" sont écrites sur bandes magnétiques (LTO-9).
- Une copie est gardée sur le site Primaire.
- Une copie est envoyée sur le site DR.
- Une copie "Air-Gapped" (hors ligne, dans un coffre-fort physique ignifugé de la Banque de la République d'Haïti - BRH).

---
*Document ID: SNISID-STORAGE-001 | Approuvé par: CISO National*
