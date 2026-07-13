# SNISID-BIO-ADN — Sécurité, Chiffrement & Contrôle d'Accès
**Document ID :** SNISID-BIO-SEC-001 | **Version :** 1.0.0

---

## 1. MODÈLE DE MENACES

| Menace | Impact | Contrôle |
|--------|--------|----------|
| Exfiltration profils ADN | CRITIQUE | Chiffrement HSM + RLS PostgreSQL |
| Accès non autorisé à BIO-identity_links | CRITIQUE | Politique RLS + audit immuable |
| Injection SQL / API abuse | ÉLEVÉ | Parameterized queries + OPA Rego |
| Corruption de données | ÉLEVÉ | Signatures ECDSA par entrée |
| Synchronisation LDIS→NDIS compromise | ÉLEVÉ | mTLS mutuel + certificat SNISID PKI |
| Insider threat (agent corrompu) | ÉLEVÉ | Audit log immuable + RLS par rôle |
| Déni de service moteur matching | MOYEN | Rate limiting + K8s HPA |

---

## 2. CHIFFREMENT AU REPOS

### 2.1 Profils STR ADN (priorité maximale)

```
Profil STR (JSON 20 loci)
        │
        ▼
AES-256-GCM (clé data-key 256 bits)
        │
        ▼
data-key chiffrée par HSM Luna (PKCS#11)
        │
        ▼
Seul le slot HSM bio-adn peut déchiffrer la data-key
```

```go
// bio-adn-service/internal/crypto/hsm.go
package crypto

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/json"
    "fmt"
    "io"

    "github.com/miekg/pkcs11"
)

type HSMCrypto struct {
    p       *pkcs11.Ctx
    session pkcs11.SessionHandle
    slotID  uint
}

// EncryptSTRProfile chiffre un profil STR avec AES-256-GCM via HSM
func (h *HSMCrypto) EncryptSTRProfile(profile STRLociData) ([]byte, error) {
    // 1. Sérialiser le profil
    plaintext, err := json.Marshal(profile)
    if err != nil {
        return nil, fmt.Errorf("marshal: %w", err)
    }

    // 2. Générer data-key aléatoire 256 bits
    dataKey := make([]byte, 32)
    if _, err := io.ReadFull(rand.Reader, dataKey); err != nil {
        return nil, fmt.Errorf("key gen: %w", err)
    }

    // 3. Chiffrer le profil avec la data-key
    block, _ := aes.NewCipher(dataKey)
    gcm, _ := cipher.NewGCM(block)
    nonce := make([]byte, gcm.NonceSize())
    io.ReadFull(rand.Reader, nonce)
    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

    // 4. Chiffrer la data-key via HSM (wrapping)
    wrappedKey, err := h.wrapKeyWithHSM(dataKey)
    if err != nil {
        return nil, fmt.Errorf("hsm wrap: %w", err)
    }

    // 5. Stocker : [4 bytes longueur wrappedKey][wrappedKey][ciphertext]
    result := make([]byte, 4+len(wrappedKey)+len(ciphertext))
    // ... assemblage binaire ...
    return result, nil
}
```

### 2.2 Données nominales sensibles

Les champs `intelligence_notes`, `current_address`, `employer` dans les tables
`per_sex_offenders` et `per_gang_members` sont stockés chiffrés avec
`pgcrypto.pgp_sym_encrypt()` (clé Vault SNISID, rotation 90 jours).

---

## 3. CONTRÔLE D'ACCÈS (RBAC + ABAC)

### 3.1 Rôles définis

| Rôle Keycloak | Accès autorisé | Restrictions |
|---------------|----------------|-------------|
| `bio.lab.technician` | Créer/modifier profils STR propre lab | Pas d'accès cross-lab |
| `bio.lab.supervisor` | Valider profils, uploader vers SDIS | Pas d'accès NDIS |
| `bio.sdis.operator` | Lire tous profils département, matcher SDIS | Pas d'accès identity_links |
| `bio.ndis.analyst` | Matcher NDIS, lire hits, pas d'identity_links | Audit obligatoire |
| `bio.dcpj.investigator` | Lire wanted/missing/gang, créer mandats | Journalisation systématique |
| `bio.dcpj.director` | Accès identity_links, déclassification | MFA obligatoire + audit |
| `bio.admin` | Administration technique | Pas d'accès données opérationnelles |
| `bio.auditor` | Lecture seule audit_log | Isolation totale |

### 3.2 Politique OPA Rego

```rego
# bio-adn-service/policies/access.rego
package snisid.bio_adn.access

import future.keywords.if
import future.keywords.in

default allow := false

# Accès aux profils STR : lab technician uniquement pour son lab
allow if {
    input.action == "read"
    input.resource == "bio_str_profiles"
    "bio.lab.technician" in input.user.roles
    input.user.lab_id == input.record.lab_id
}

# Accès identity_links : DCPJ Director uniquement + ordonnance judiciaire requise
allow if {
    input.action in {"read", "create"}
    input.resource == "bio_identity_links"
    "bio.dcpj.director" in input.user.roles
    input.context.court_order_ref != ""
    input.context.mfa_verified == true
}

# Matching NDIS : analystes NDIS uniquement
allow if {
    input.action == "match"
    input.resource == "bio_str_profiles"
    "bio.ndis.analyst" in input.user.roles
    input.context.case_number != ""
    input.context.purpose in {"criminal_investigation", "missing_person", "identification"}
}

# Lecture mandats : agents PNH / DCPJ
allow if {
    input.action == "read"
    input.resource == "per_wanted_persons"
    some role in ["bio.dcpj.investigator", "bio.ndis.analyst"]
    role in input.user.roles
}

# Refus explicite : personne ne peut supprimer un log d'audit
deny if {
    input.action == "delete"
    input.resource == "bio_audit_log"
}
```

---

## 4. AUTHENTIFICATION DES SERVICES

Tous les appels inter-services SNISID-BIO-ADN utilisent **mTLS** :

```yaml
# Certificats SPIFFE/SPIRE — identité des workloads
spiffe://snisid.gov.ht/bio-adn/ldis-pap        # Labo LDIS Port-au-Prince
spiffe://snisid.gov.ht/bio-adn/sdis-ouest       # SDIS Département Ouest
spiffe://snisid.gov.ht/bio-adn/ndis-central     # NDIS National
spiffe://snisid.gov.ht/bio-adn/matching-engine  # Moteur de matching
spiffe://snisid.gov.ht/bio-adn/api-gateway      # Gateway REST
```

---

## 5. AUDIT IMMUABLE

Chaque opération sur les tables forensiques génère une entrée dans
`bio_audit_log`, signée numériquement (ECDSA-P256, clé privée HSM).
La signature empêche la modification a posteriori des logs.

```go
// Trigger PostgreSQL pour audit automatique
// (défini dans la migration SQL)
CREATE OR REPLACE FUNCTION bio_audit_trigger() RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO bio_audit_log (
        event_type, table_name, record_id,
        officer_niu, agency_code, purpose,
        action, details, created_at
    ) VALUES (
        'data.' || TG_OP,
        TG_TABLE_NAME,
        COALESCE(NEW.record_id, OLD.record_id),
        current_setting('snisid.officer_niu', true),
        current_setting('snisid.agency_code', true),
        current_setting('snisid.purpose', true),
        TG_OP,
        row_to_json(COALESCE(NEW, OLD)),
        NOW()
    );
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;
```

---

## 6. RÉTENTION ET EXPUNGEMENT

| Index | Durée de rétention | Condition d'expungement |
|-------|--------------------|------------------------|
| BIO-CON (condamnés) | Durée condamnation + 10 ans | Décision judiciaire finale |
| BIO-ARR (arrestés) | 3 ans max si pas de condamnation | Acquittement / classement sans suite |
| BIO-FSC (scènes crime) | Durée prescription + 5 ans | Aucun (preuve judiciaire) |
| BIO-DIS (disparus) | Jusqu'à identification | Identification ou décès confirmé |
| PER-REC (recherchés) | Durée mandat + 2 ans | Arrestation / annulation mandat |
| PER-DIS (disparus) | Jusqu'à localisation | Retour / décès confirmé |
| BIE-VEH (véhicules) | 5 ans après récupération | Récupération + rapport |

L'expungement est **irréversible** et génère une entrée `bio_audit_log`
de type `EXPUNGE` signée par le DCPJ Director + juge compétent.
