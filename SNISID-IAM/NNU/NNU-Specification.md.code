# NNU — Numéro National Unique

> **Spécification de l'identifiant national maître SNISID**  
> **Version :** 1.0.0  
> **Classification :** SOUVERAIN — RESTREINT  
> **Dernière mise à jour :** 2026-05-25

---

## 1. OBJECTIF

Le **NNU (Numéro National Unique)** est la clé nationale souveraine. Il identifie de manière **unique, persistante et non réutilisable** chaque citoyen haïtien dans le système SNISID.

> **Le NNU est la référence absolue d'identité nationale.**

---

## 2. CARACTÉRISTIQUES OBLIGATOIRES

| Critère | Exigence | Vérification |
|---------|----------|-------------|
| Unique | ✅ Un seul NNU par citoyen | Index unique en base |
| Non réutilisable | ✅ Jamais recyclé | Politique nationale |
| Persistant | ✅ À vie du citoyen | Stockage immuable |
| Audit trail | ✅ Chaque usage tracé | Journal d'audit |
| Versionné | ✅ Historique des changements | Séquentiel |
| Sécurisé | ✅ Chiffré et signé | PKI nationale |

---

## 3. STRUCTURE DU NNU

### 3.1 Format

```
NNU-HA-YYYYMMDD-XXXXXXXX-CC-V

Segments :
  NNU  = Préfixe fixe (3 caractères)
  HA   = Code pays ISO 3166-1 alpha-2 (Haïti)
  YYYYMMDD = Date d'enrôlement (8 caractères)
  XXXXXXXX = Séquence cryptographique unique (8 caractères alphanumériques A-Z, 2-9)
  CC   = Checksum (2 caractères de contrôle)
  V    = Version du NNU (toujours 1 pour la création)
```

**Exemple complet :** `NNU-HA-20260525-A7K9M2X4-73-1`

### 3.2 Génération de la Séquence

```
Séquence = HMAC-SHA256(
  clé_secrète_nationale,
  timestamp_enrôlement + random_128bit
)[:8]
```

- **Clé secrète nationale** : stockée dans le HSM souverain
- **Random** : CSPRNG certifié (NIST SP 800-90A)
- **Résultat** : 8 caractères alphanumériques (A-Z, 2-9)

### 3.3 Checksum

Formule ISO 7064 MOD 97-10, adaptée au format NNU.

---

## 4. ENREGISTREMENT NNU

### 4.1 Processus de Création

```
[Enrôlement Citoyen] → [Validation Documents] → [Génération NNU] → [Signature PKI] → [Archivage] → [Notification]
```

### 4.2 Étapes Détaillées

| Étape | Action | Acteur | Système |
|-------|--------|--------|---------|
| 1 | Collecte données citoyennes | Enrollment Officer | Terminal certifié |
| 2 | Vérification documents légaux | Officer + Système | OCR + Validation |
| 3 | Capture biométrique | Officer | Capteur certifié |
| 4 | Dé-duplication | Système | Moteur de détection |
| 5 | Génération NNU | Système | HSM souverain |
| 6 | Signature PKI | Système | PKI nationale |
| 7 | Archivage | Système | Archive immuable |
| 8 | Notification citoyen | Système | SMS / Email / Wallet |

---

## 5. BASE DE DONNÉES NNU

```sql
CREATE TABLE nnu_registry (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nnu             VARCHAR(32) UNIQUE NOT NULL,
    citizen_id      UUID NOT NULL REFERENCES citizens(id),
    country_code    VARCHAR(2) DEFAULT 'HA' NOT NULL,
    enrollment_date DATE NOT NULL,
    sequence_code   VARCHAR(8) NOT NULL,
    checksum        VARCHAR(2) NOT NULL,
    version         INTEGER DEFAULT 1 NOT NULL,
    hmac_signature  BYTEA NOT NULL,
    encrypted_payload BYTEA NOT NULL,
    pkcs_signature  BYTEA NOT NULL,
    status          VARCHAR(20) DEFAULT 'active' CHECK (
        status IN ('active', 'suspended', 'revoked', 'deceased')
    ),
    created_by      UUID NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at      TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    deleted_at      TIMESTAMPTZ,
    previous_version UUID REFERENCES nnu_registry(id),
    audit_trail_id  UUID NOT NULL REFERENCES audit_trail(id)
);

CREATE TABLE nnu_audit_log (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nnu_ref         VARCHAR(32) NOT NULL,
    operation       VARCHAR(20) NOT NULL CHECK (
        operation IN ('create', 'verify', 'suspend', 'revoke', 'restore', 'access')
    ),
    performed_by    UUID NOT NULL,
    performed_at    TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    ip_address      INET NOT NULL,
    device_id       VARCHAR(64),
    location        POINT,
    reason          TEXT,
    previous_status VARCHAR(20),
    new_status      VARCHAR(20),
    signature       BYTEA NOT NULL
);
```

---

## 6. OPÉRATIONS NNU

### 6.1 API Nationale

| Endpoint | Méthode | Description | Rôle requis |
|----------|---------|-------------|-------------|
| `/api/v1/nnu` | POST | Générer un nouveau NNU | enrollment_officer |
| `/api/v1/nnu/{nnu}/verify` | POST | Vérifier la validité d'un NNU | Tout rôle authentifié |
| `/api/v1/nnu/{nnu}/suspend` | POST | Suspendre un NNU | admin, judicial_officer |
| `/api/v1/nnu/{nnu}/revoke` | POST | Révoquer un NNU | admin, judicial_officer |
| `/api/v1/nnu/{nnu}/history` | GET | Historique complet d'un NNU | admin, citizen (son NNU) |

### 6.2 Règles Opérationnelles

| Opération | Autorisation | MFA | Audit | Signature |
|-----------|-------------|-----|-------|-----------|
| Création | Enrollment Officer | ✅ | ✅ | ✅ |
| Vérification | Tout rôle authentifié | ✅ | ✅ | ❌ |
| Suspension | Admin / Judicial | ✅ | ✅ | ✅ |
| Révocation | Admin + Judicial | ✅ | ✅ | ✅ |
| Restauration | Admin + Comité | ✅ | ✅ | ✅ |
| Consultation | Citoyen (son NNU) | ✅ | ✅ | ❌ |

---

## 7. SÉCURITÉ NNU

| Couche | Mécanisme |
|--------|-----------|
| Stockage | Chiffrement AES-256-GCM |
| Transmission | TLS 1.3 + mTLS |
| Accès | MFA + RBAC + ABAC |
| Signature | PKI nationale (ECDSA P-384) |
| Audit | Journal immuable |
| HSM | Module matériel certifié FIPS 140-3 |

### 7.1 Règle d'Intégrité

> **NNU est IMMUABLE une fois créé.** Aucune modification de la séquence, du checksum ou du format. Seul le statut peut évoluer (active → suspended → revoked).

### 7.2 Détection de Fraude

| Signal | Action |
|--------|--------|
| NNU utilisé depuis 2+ géolocalisations simultanées | Alerte fraude |
| Tentative de vérification avec biométrie non correspondante | Alerte tentative usurpation |
| NNU demandé avec documents suspects | Investigation |
| Multiples suspensions/révocations sur même NNU | Escalade judiciaire |

---

## 8. POLITIQUES NATIONALES

### 8.1 Non-Réutilisabilité

> **Un NNU révoqué ne sera JAMAIS réattribué.**

Même en cas de décès, le NNU reste dans le registre avec le statut `deceased`.

### 8.2 Persistance

> **Le NNU est attribué à vie.** Il ne change jamais, même en cas de :
- Changement de nom (mariage, correction)
- Changement d'adresse
- Changement de statut civil

### 8.3 Confidentialité

> **Le NNU n'est jamais affiché en entier dans les interfaces publiques.**

Masquage : `NNU-HA-****-****-**-1`

Affichage complet uniquement :
- Dans le wallet citoyen (après MFA)
- Pour les agents autorisés (après vérification de rôle)

---

## 9. VALIDATION DU NNU

### 9.1 Algorithme de Validation

```python
def validate_nnu(nnu: str) -> dict:
    """Valide un NNU selon les règles nationales."""
    errors = []
    
    # Format regex
    pattern = r'^NNU-HA-\d{8}-[A-Z2-9]{8}-\d{2}-\d$'
    if not re.match(pattern, nnu):
        errors.append("Format invalide")
        return {"valid": False, "errors": errors}
    
    # Vérifier date d'enrôlement
    date_str = nnu[7:15]
    enrollment_date = datetime.strptime(date_str, '%Y%m%d')
    if enrollment_date > datetime.now():
        errors.append("Date d'enrôlement future")
    
    # Vérifier checksum
    body = nnu.replace('-', '').replace('NNU', '').replace('HA', '')
    expected_checksum = compute_checksum(body[:-2])
    actual_checksum = body[-2:]
    if expected_checksum != actual_checksum:
        errors.append("Checksum invalide")
    
    # Vérifier dans le registre
    registry_check = lookup_in_registry(nnu)
    if not registry_check:
        errors.append("NNU non trouvé dans le registre national")
    
    return {"valid": len(errors) == 0, "errors": errors}
```

### 9.2 Tests de Validation

| Test | NNU | Résultat attendu |
|------|-----|-----------------|
| Format valide | `NNU-HA-20260525-A7K9M2X4-73-1` | ✅ |
| Format invalide | `NNU-HA-20260525-a7k9m2x4-73-1` | ❌ |
| Checksum invalide | `NNU-HA-20260525-A7K9M2X4-00-1` | ❌ |
| Date future | `NNU-HA-20271201-A7K9M2X4-73-1` | ❌ |
| NNU inconnu | `NNU-HA-20260525-XXXXXXXX-73-1` | ❌ |

---

## 10. RÉFÉRENCES

| Document | Lien |
|----------|------|
| National Identity Domain Model | `../Identity-Model/National-Identity-Domain-Model.md` |
| Duplicate Detection Engine | `../Identity-Model/Duplicate-Detection-Engine.md` |
| PKI Specification | `../Governance/Identity-Governance-Specification.md` |

---

> **Le NNU est la clé nationale souveraine. Son intégrité est la priorité absolue du SNISID.**
