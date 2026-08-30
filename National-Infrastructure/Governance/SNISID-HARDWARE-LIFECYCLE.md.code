---
# ============================================================
# SNISID-Infra — Hardware Governance Framework
# Sécurité de la chaîne d'approvisionnement (Supply Chain)
# Document ID: SNISID-HW-GOV-001
# Version: 1.0.0
# ============================================================

## 1. SOVEREIGN SUPPLY CHAIN SECURITY

La souveraineté numérique commence par le matériel. Un serveur backdooré en usine compromet tout le reste de l'édifice (Zero Trust échoue si le hardware est compromis).

## 2. RÈGLES D'APPROVISIONNEMENT (Procurement)

1. **Fournisseurs Agréés :** Le matériel réseau et compute doit provenir de fournisseurs certifiés (ex: interdiction des marques sous embargo international ou sans audit de sécurité transparent).
2. **Firmware Validation :** Avant qu'un serveur ne soit raccordé au Datacenter, son BIOS/UEFI, le firmware de ses cartes réseaux (NIC), et de la carte mère (BMC/iDRAC) doivent être flashés avec une version validée par l'équipe de sécurité.
3. **Hardware Trust Anchors :** Tous les serveurs doivent être équipés d'un TPM 2.0 (Trusted Platform Module) pour garantir un "Secure Boot". Si le système d'exploitation a été altéré, le serveur refusera de démarrer.

## 3. LIFECYCLE MANAGEMENT (Cycle de vie)

| Étape | Processus | Outil |
|-------|-----------|-------|
| **1. Réception** | Scan MAC/Série, flashage initial, placement en Quarantaine Réseau. | Netbox / Snipe-IT |
| **2. Provisioning**| Installation de l'OS via PXE (MaaS/Ironic) avec Secure Boot. | OpenStack Ironic |
| **3. Production** | Le serveur rejoint le cluster. Surveillance de l'intégrité via eBPF. | Prometheus / Cilium |
| **4. Retrait** | Destruction physique (Degaussing/Shredding) des disques NVMe/SSD contenant des données régaliennes avant mise au rebut. | Procédure Manuelle |

---
*Document ID: SNISID-HW-GOV-001 | Approuvé par: Directeur des Acquisitions (AND)*
