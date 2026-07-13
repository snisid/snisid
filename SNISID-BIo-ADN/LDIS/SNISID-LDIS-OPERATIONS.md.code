# SNISID-BIO-ADN — LDIS-HT Opérations Locales
**Document ID :** SNISID-LDIS-OPS-001 | **Version :** 1.0.0

---

## 1. LISTE DES LABORATOIRES LDIS-HT

| Code | Nom | Institution | Département | Statut |
|------|-----|-------------|-------------|--------|
| `LDIS-PAP-001` | Labo Médico-Légal PAP | DCPJ / MSPP | Ouest | PRIORITAIRE |
| `LDIS-CAP-001` | Labo Médico-Légal Cap-Haïtien | PNH / MSPP | Nord | PRIORITAIRE |
| `LDIS-LES-001` | Labo Médico-Légal Les Cayes | MSPP | Sud | À CRÉER |
| `LDIS-GON-001` | Labo Médico-Légal Gonaïves | MSPP | Artibonite | À CRÉER |
| `LDIS-JAC-001` | Labo Médico-Légal Jacmel | MSPP | Sud-Est | À CRÉER |
| `LDIS-HIN-001` | Labo Médico-Légal Hinche | MSPP | Centre | À CRÉER |

> Phase 1 : Démarrer avec LDIS-PAP-001 et LDIS-CAP-001.
> Phase 2 (an 2) : Ouvrir les 4 laboratoires départementaux restants.

---

## 2. PROCÉDURES OPÉRATIONNELLES LDIS

### 2.1 Collecte d'un prélèvement ADN — SOP-LDIS-001

```
ÉTAPE 1 — Autorisation
    └── Vérifier ordonnance judiciaire (BIO-CON, BIO-ARR) ou
        consentement écrit (BIO-DIS volontaire)

ÉTAPE 2 — Prélèvement
    └── Écouvillon buccal (2 écouvillons stériles minimum)
    └── Étiquetage : NOM_DOSSIER_DATE_INITIALES_TECHNICIEN
    └── Scellé tamper-evident numéroté

ÉTAPE 3 — Chain of custody (CoC)
    └── Formulaire CoC signé par : agent arrêtant + technicien labo
    └── Scan du CoC dans le système SNISID-BIO-ADN
    └── Attribution du specimen_number automatique

ÉTAPE 4 — Analyse STR
    └── Extraction ADN (Chelex 100 ou kit commercial accrédité)
    └── PCR + électrophorèse capillaire (ABI 3500, 3130, ou équivalent)
    └── Interprétation 20 loci CODIS Core

ÉTAPE 5 — Contrôle qualité
    └── Quality score calculé automatiquement
    └── Si quality_score < seuil → répéter analyse (max 3 tentatives)
    └── Si toujours < seuil → rapport d'échec + conservation échantillon

ÉTAPE 6 — Soumission SNISID-BIO-ADN
    └── API POST /dna/profiles
    └── Confirmation sample_id reçue
    └── Archivage physique de l'échantillon (conservation 25 ans minimum)
```

### 2.2 Upload LDIS → SDIS (SOP-LDIS-002)

- **Fréquence :** Quotidien à 02h00 (heure Haïti, UTC-4)
- **Trigger :** Automatique (scheduler) ou manuel (superviseur)
- **Condition :** Quality score ≥ seuil minimal par index
- **Réseau :** Tunnel WireGuard chiffré SNISID-PKI vers SDIS

```bash
# Commande manuelle upload (en cas d'urgence)
snisid-bio-cli upload \
  --level ldis-to-sdis \
  --lab-code LDIS-PAP-001 \
  --date-from 2026-06-01 \
  --date-to 2026-06-09 \
  --operator-niu HTI-XXXXXXXXXX
```

### 2.3 Gestion des conflits de scellés

Si un specimen_number est soumis deux fois :
1. Le système rejette le doublon (contrainte UNIQUE)
2. Un événement `BIO-DUPLICATE-SPECIMEN` est publié sur Kafka
3. Le superviseur labo est notifié pour investigation

---

## 3. ÉQUIPEMENTS REQUIS PAR LABORATOIRE LDIS

| Équipement | Modèle recommandé | Rôle |
|------------|-------------------|------|
| Séquenceur capillaire | ABI 3500 / 3500xL | Analyse STR |
| Extracteur ADN | QIAsymphony SP | Extraction automatisée |
| Thermocycleur | Applied Biosystems Veriti | PCR |
| Lecteur de codes-barres | Honeywell Xenon | Traçabilité scellés |
| Workstation SNISID | Dell OptiPlex 7010 (validé SNISID) | Interface système |
| Congélateur -80°C | Thermo Scientific | Conservation échantillons |
| Sorbonne | Thermo Scientific | Manipulation sécurisée |

---

## 4. ACCRÉDITATION DES LABORATOIRES

Les laboratoires LDIS-HT doivent satisfaire :

1. **ISO/IEC 17025** — Accréditation laboratoires d'essais
2. **Normes SWGDAM** (Scientific Working Group for DNA Analysis Methods) adaptées
3. **Audit SNISID** semestriel par la Direction Nationale
4. **Contrôle qualité externe** : échange d'échantillons témoins entre labs

---

## 5. FORMATION DU PERSONNEL LDIS

| Formation | Durée | Fréquence | Dispensé par |
|-----------|-------|-----------|--------------|
| STR Analysis | 40h | Initiale | Direction SNISID + partenaire international |
| SNISID-BIO-ADN System | 16h | Initiale | Direction SNISID |
| Chain of Custody | 8h | Annuelle | Directeur labo |
| Sécurité des données | 4h | Annuelle | SNISID IT Security |
| Courtroom testimony | 8h | Initiale | MJSP / Parquet |
