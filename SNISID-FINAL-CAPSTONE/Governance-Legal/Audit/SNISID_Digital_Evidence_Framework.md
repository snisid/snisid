# SNISID — CADRE DE PREUVE NUMÉRIQUE

**Classification :** CADRE OPÉRATIONNEL — PREUVE NUMÉRIQUE
**Référence :** SNISID-EVID-001
**Version :** 1.0
**Date :** 25 mai 2026

---

## 1. OBJECTIF

Garantir que toutes les preuves numériques générées, collectées et conservées par le SNISID sont juridiquement exploitables, admissibles devant les tribunaux haïtiens et conformes aux standards forensiques internationaux.

---

## 2. FONCTIONS SUPPORTÉES

### 2.1 Evidence Preservation (Préservation des Preuves)

| Mécanisme | Description | Standard |
|-----------|-------------|---------|
| Journalisation immutable | Tous les événements SNISID enregistrés en append-only | WORM storage |
| Hachage cryptographique | SHA-384 sur chaque entrée de journal | RFC 6234 |
| Horodatage qualifié | TSA nationale sur chaque transaction | RFC 3161 |
| Réplication | Copie temps réel sur site secondaire | Synchronisation < 1s |
| Archivage long terme | Migration de format périodique | 10 ans minimum |
| Chiffrement au repos | AES-256 pour les preuves stockées | FIPS 197 |
| Signature de lot | Merkle tree toutes les 60 secondes | Intégrité de bloc |

**Architecture de préservation :**
```
Événement système
    → Journalisation immédiate (buffer < 100ms)
    → Hash SHA-384 de l'entrée
    → Ajout au Merkle tree courant
    → Horodatage qualifié du bloc (toutes les 60s)
    → Signature du bloc par le service
    → Réplication synchrone sur site DR
    → Archivage quotidien vers stockage long terme
    → Vérification d'intégrité hebdomadaire
```

### 2.2 Chain of Custody (Chaîne de Traçabilité)

| Étape | Documentation Requise | Responsable |
|-------|----------------------|-------------|
| Génération | Système source, timestamp, contexte | Système automatique |
| Collection | Méthode, agent, hash, PV | Agent forensique |
| Transport | Mode, protection, vérification hash | Agent transporteur |
| Stockage | Localisation, sécurité, accès | Gardien des preuves |
| Analyse | Méthodes, outils, hash avant/après | Analyste forensique |
| Présentation | Contexte, intégrité vérifiée, PV | Agent habilité |
| Restitution/Destruction | Décision judiciaire, méthode | Gardien des preuves |

**Registre de traçabilité (Chain of Custody Log) :**
```json
{
  "evidence_id": "EVD-2026-XXXX-XXXX",
  "type": "audit_log_extract",
  "classification": "CONFIDENTIEL",
  "hash_sha384": "...",
  "chain": [
    {
      "action": "GENERATED",
      "timestamp": "2026-05-25T10:00:00Z",
      "actor": "SYSTEM:SNISID-CORE",
      "location": "DC-PRIMARY",
      "hash_verified": true
    },
    {
      "action": "COLLECTED",
      "timestamp": "2026-05-25T14:30:00Z",
      "actor": "AGENT:FOR-1234",
      "authorization": "MANDAT-2026-5678",
      "method": "forensic_copy_bit_for_bit",
      "hash_verified": true
    },
    {
      "action": "TRANSFERRED",
      "timestamp": "2026-05-25T15:00:00Z",
      "from": "AGENT:FOR-1234",
      "to": "STORAGE:EVIDENCE-VAULT-01",
      "transport": "encrypted_usb_sealed",
      "hash_verified": true
    }
  ]
}
```

### 2.3 Forensic Integrity (Intégrité Forensique)

| Principe | Implementation |
|----------|---------------|
| Non-altération | Copie forensique (bit-for-bit), original non modifié |
| Vérifiabilité | Hash avant, pendant et après chaque opération |
| Reproductibilité | Méthodes documentées, résultats reproductibles |
| Indépendance | Outils forensiques certifiés et calibrés |
| Documentation | Chaque étape documentée dans le PV |

**Standards forensiques appliqués :**
- ISO 27037 — Lignes directrices pour l'identification, la collecte, l'acquisition et la préservation de preuves numériques
- ISO 27041 — Orientation sur la garantie d'aptitude des méthodes d'investigation
- ISO 27042 — Lignes directrices pour l'analyse et l'interprétation des preuves numériques
- ISO 27043 — Principes et processus d'investigation

**Outils forensiques autorisés :**

| Catégorie | Exigence |
|-----------|----------|
| Copie forensique | Outil certifié, write-blocker obligatoire |
| Analyse disque | Outil validé, hash vérifié |
| Analyse réseau | Capture complète, horodatée |
| Analyse mémoire | Dump authentifié, hash vérifié |
| Analyse mobile | Outil certifié, PV détaillé |

### 2.4 Court Admissibility (Admissibilité Judiciaire)

**Critères d'admissibilité :**

| Critère | Test | Preuve |
|---------|------|--------|
| Authenticité | L'origine est-elle vérifiable ? | Signature + certificat |
| Intégrité | Les données sont-elles intactes ? | Hash vérifié à chaque étape |
| Fiabilité | La méthode est-elle fiable ? | Conformité ISO 27037 |
| Pertinence | Les données sont-elles pertinentes ? | Lien avec les faits |
| Licéité | La collecte est-elle légale ? | Mandat / base légale |
| Complétude | La chaîne est-elle continue ? | Chain of custody log |

**Modèle de rapport pour le tribunal :**
```
RAPPORT DE PREUVE NUMÉRIQUE
═══════════════════════════
Affaire : [Référence]
Juridiction : [Tribunal]
Expert : [Nom, qualification, numéro]

1. IDENTIFICATION DE LA PREUVE
   - Identifiant : EVD-XXXX
   - Type : [journal d'audit / document / biométrie / etc.]
   - Source : [système / dispositif]
   - Date de collecte : [date + heure UTC]

2. CHAÎNE DE TRAÇABILITÉ
   [Détail complet avec signatures]

3. INTÉGRITÉ
   - Hash à la collecte : [SHA-384]
   - Hash actuel : [SHA-384]
   - Statut : [IDENTIQUE / DIFFÉRENT]

4. ANALYSE
   [Méthodes, outils, constatations]

5. CONCLUSIONS
   [Avis technique de l'expert]

6. ANNEXES
   [Données techniques, captures, logs]

Signature de l'expert : [Signature qualifiée]
Date : [Horodatage qualifié]
```

### 2.5 Immutable Audit Logs (Journaux d'Audit Immutables)

**Architecture des journaux :**

| Composant | Description |
|-----------|-------------|
| Write-Once Storage | Stockage WORM (Write Once Read Many) |
| Append-Only Database | Aucune modification ni suppression possible |
| Hash Chain | Chaque entrée liée à la précédente par hash |
| Merkle Tree | Arbre de hash pour vérification de bloc |
| Signature de bloc | Signature serveur toutes les 60 secondes |
| Horodatage TSA | Horodatage qualifié par bloc |
| Réplication | Copie synchrone sur 2 sites minimum |

**Format d'entrée de journal :**
```json
{
  "log_id": "LOG-20260525-100000-000001",
  "timestamp": "2026-05-25T10:00:00.000Z",
  "tsa_timestamp": "2026-05-25T10:00:00.123Z",
  "previous_hash": "SHA384:...",
  "event_type": "IDENTITY_VERIFICATION",
  "actor": {
    "type": "AGENT",
    "id": "AGT-12345",
    "agency": "PNH",
    "authentication_level": "LOA3"
  },
  "subject": {
    "type": "CITIZEN",
    "nni_hash": "SHA384:...",
    "data_accessed": ["identity_basic", "photo"]
  },
  "action": {
    "type": "READ",
    "result": "SUCCESS",
    "justification": "ENQUETE-2026-7890"
  },
  "context": {
    "source_ip_hash": "SHA384:...",
    "workstation": "WS-PNH-PAP-042",
    "session_id": "SES-...",
    "geo_location": "PAP-CENTRAL"
  },
  "entry_hash": "SHA384:..."
}
```

**Vérification d'intégrité :**
- Vérification automatique horaire de la chaîne de hash
- Comparaison quotidienne entre sites
- Audit mensuel d'intégrité complet
- Alerte immédiate en cas de divergence

---

## 3. PROCÉDURES FORENSIQUES

### 3.1 Procédure de Collecte de Preuve SNISID

```
1. DEMANDE
   - Requête officielle (mandat judiciaire ou demande habilitée)
   - Vérification de l'autorité du demandeur
   - Enregistrement de la demande

2. PRÉPARATION
   - Identification des données à collecter
   - Préparation de l'environnement forensique
   - Vérification des outils

3. COLLECTE
   - Extraction par agent forensique qualifié
   - Copie forensique (bit-for-bit)
   - Calcul et enregistrement du hash
   - Procès-verbal de collecte

4. VÉRIFICATION
   - Vérification du hash
   - Vérification de la complétude
   - Validation par un second agent

5. MISE SOUS SCELLÉ
   - Scellé numérique (hash + signature + horodatage)
   - Stockage sécurisé
   - Enregistrement dans le registre des preuves

6. TRANSMISSION
   - Remise au demandeur contre décharge
   - Mise à jour de la chaîne de traçabilité
```

---

## 4. RÉTENTION ET DESTRUCTION

| Type de preuve | Durée de rétention | Destruction |
|---------------|-------------------|------------|
| Journaux d'audit opérationnel | 10 ans | Automatique + PV |
| Preuves judiciaires | Jusqu'à décision définitive + 5 ans | Sur décision du magistrat |
| Preuves de crimes graves | 30 ans minimum | Sur décision du Procureur |
| Données d'incident de sécurité | 7 ans | Automatique + PV |
| Archives historiques | Permanent | Interdite |

---

*Document cadre préparé dans le cadre de la Phase 14 — SNISID National Legal Framework*
