---
# ============================================================
# SNISID-Security — National Offline Security Operations Model
# Opérations PNH Déconnectées & Synchronisation Edge
# Document ID: SNISID-OFFLINE-SEC-001
# Version: 1.0.0
# ============================================================

## 1. NÉCESSITÉ DU MODE OFFLINE

En Haïti, la connectivité internet (fibre, 4G) n'est pas garantie à 100%, particulièrement dans les villes de province ou lors de catastrophes naturelles. Les opérations de police et de sécurité ne peuvent pas s'arrêter.

## 2. TOPOLOGIE OFFLINE-FIRST (PNH & FRONTIÈRES)

Le modèle s'appuie sur la technologie Edge définie dans la Phase 1 (K3s + NATS).

### 2.1 Unités Mobiles (Smartphone/Tablette Agent)
- Possèdent une base de données locale (SQLite/SQLCipher).
- Contiennent un "Hot Cache" (Les 10 000 individus les plus recherchés, téléchargé le matin).
- Permettent de créer des `ArrestReport` hors ligne.

### 2.2 K3s Edge Node (Commissariat / Poste Frontière)
- Maintient un "Regional Cache" (Les données de la juridiction).
- Reçoit les données des unités mobiles via réseau local (Wi-Fi intra-commissariat).
- Synchronise avec le SNISID-Core Central dès que la liaison (VSAT/4G/Fibre) revient.

## 3. RÉSOLUTION DES CONFLITS (POLICE)

Si deux agents arrêtent le même individu offline et synchronisent plus tard :
- L'Identity Registry acceptera le premier `ArrestEvent` (par horodatage).
- Le second déclenchera un `DuplicateArrestWarning`, nécessitant la fusion manuelle par un superviseur (DDO).

---
*Document ID: SNISID-OFFLINE-SEC-001 | Approuvé par: Architecte Souverain*
