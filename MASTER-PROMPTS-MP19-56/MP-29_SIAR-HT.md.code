# MP-29 — SIAR-HT
## Système National d'Information sur les Armes à Feu d'Haïti
### Master Prompt d'Implémentation — SNISID

```
Classification   : SECRET / USAGE OFFICIEL EXCLUSIF
Module           : MP-29 | Code : SIAR-HT
Dépendances      : BIAR-HT (MP-30), TRAF-AR (MP-32), FIR-HT (MP-20), GANG-HT (MP-24)
Normes           : ATF eTRACE, INTERPOL iARMS, Traité sur le Commerce des Armes (TCA),
                   Loi haïtienne sur les armes (Décret 1/2/1973 révisé)
Acteurs          : DCPJ, PNH Unité Armes, Douanes, MJSP, ATF liaison USA
```

---

## 1. CONTEXTE

Environ 500,000 armes illicites circulent en Haïti selon l'ONU. Plus de 80% proviennent
des États-Unis via des trafiquants (straw purchases en Floride, Géorgie). Il n'existe
aucun registre national des armes légales ou des licences. Ce module crée le premier
registre complet : armes légales enregistrées, armes saisies, licences, et traçage.

### Catégories d'armes en Haïti

| Catégorie               | Exemple                    | Régulation          |
|-------------------------|----------------------------|---------------------|
| Armes légères (SALW)    | Pistolets, fusils           | Permis MJSP requis  |
| Armes automatiques      | AR-15, AK-47               | Interdites civils   |
| Armes lourdes           | RPG, MG                    | Militaires/police   |
| Munitions               | Cartouches, chargeurs       | Contrôle quota      |
| Armes artisanales       | Tremblé (pistolet local)   | Fabrication illicite|

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE siar_weapon_type AS ENUM (
    'HANDGUN','RIFLE','SHOTGUN','SUBMACHINE_GUN','ASSAULT_RIFLE',
    'MACHINE_GUN','SNIPER','RPG','GRENADE','HOMEMADE','OTHER'
);

CREATE TYPE siar_status AS ENUM (
    'REGISTERED','REPORTED_STOLEN','SEIZED','DESTROYED',
    'REPORTED_LOST','TRANSFERRED','DEACTIVATED'
);

CREATE TYPE siar_registration_type AS ENUM (
    'CIVILIAN','POLICE','MILITARY','SECURITY_COMPANY',
    'EMBASSY','ILLEGAL_FOUND','HISTORICAL'
);

CREATE TABLE siar_firearms (
    firearm_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_siar_id    VARCHAR(25) UNIQUE NOT NULL,  -- SIAR-HT-NNNNNN
    serial_number       VARCHAR(100),
    make                VARCHAR(100) NOT NULL,
    model               VARCHAR(100) NOT NULL,
    caliber             VARCHAR(30) NOT NULL,
    weapon_type         siar_weapon_type NOT NULL,
    manufacture_year    SMALLINT,
    manufacture_country CHAR(3),
    status              siar_status NOT NULL DEFAULT 'REGISTERED',
    reg_type            siar_registration_type NOT NULL,

    -- Proprietaire (si legale)
    owner_snisid_id     UUID,
    owner_entity_name   VARCHAR(200),     -- Si organisation
    license_number      VARCHAR(50),
    license_expiry      DATE,

    -- Traçage importation
    import_date         DATE,
    import_country      CHAR(3),
    import_permit_ref   VARCHAR(100),
    importer_name       VARCHAR(200),
    customs_entry_ref   VARCHAR(100),

    -- Localisation actuelle
    current_dept_code   CHAR(2),
    storage_location    TEXT,

    -- Liens criminels
    fir_record_id       UUID,
    gang_id             UUID,
    case_references     TEXT[] DEFAULT '{}',

    -- INTERPOL iARMS
    iarms_ref           VARCHAR(50),
    atf_etrace_ref      VARCHAR(50),

    notes               TEXT,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE siar_licenses (
    license_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    license_number      VARCHAR(50) UNIQUE NOT NULL,
    holder_snisid_id    UUID NOT NULL,
    license_type        VARCHAR(50) NOT NULL,   -- CARRY, POSSESS, DEALER, COLLECTOR
    firearms_authorized INTEGER DEFAULT 1,
    issue_date          DATE NOT NULL,
    expiry_date         DATE NOT NULL,
    issuing_authority   VARCHAR(100) NOT NULL,
    is_active           BOOLEAN DEFAULT TRUE,
    revocation_reason   TEXT,
    revoked_at          TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE siar_transfers (
    transfer_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    firearm_id          UUID NOT NULL REFERENCES siar_firearms(firearm_id),
    from_owner_id       UUID,
    to_owner_id         UUID,
    transfer_type       VARCHAR(50),  -- SALE, GIFT, INHERITANCE, CONFISCATION
    transfer_date       DATE NOT NULL,
    permit_ref          VARCHAR(100),
    authorized_by       UUID,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE siar_seizures (
    seizure_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    firearm_id          UUID REFERENCES siar_firearms(firearm_id),
    seizure_date        TIMESTAMPTZ NOT NULL,
    seizing_unit        VARCHAR(50) NOT NULL,
    seizing_officer     UUID,
    location_desc       VARCHAR(300),
    dept_code           CHAR(2),
    context             TEXT,          -- Circonstances de saisie
    from_person_id      UUID,          -- Personne chez qui saisie
    gang_id             UUID,
    case_reference      VARCHAR(100),
    disposed_of         BOOLEAN DEFAULT FALSE,
    disposal_method     VARCHAR(50),   -- DESTROYED, KEPT_AS_EVIDENCE, RETURNED
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_siar_serial    ON siar_firearms(serial_number) WHERE serial_number IS NOT NULL;
CREATE INDEX idx_siar_status    ON siar_firearms(status);
CREATE INDEX idx_siar_gang      ON siar_firearms(gang_id) WHERE gang_id IS NOT NULL;
CREATE INDEX idx_siar_owner     ON siar_firearms(owner_snisid_id) WHERE owner_snisid_id IS NOT NULL;
CREATE INDEX idx_siar_iarms     ON siar_firearms(iarms_ref) WHERE iarms_ref IS NOT NULL;
CREATE INDEX idx_siar_licenses  ON siar_licenses(holder_snisid_id, is_active);

COMMIT;
```

---

## 3. SERVICE GO CLÉ

```go
package service

import (
    "context"
    "github.com/google/uuid"
    "github.com/snisid/siar-svc/internal/domain"
)

func (s *FirearmService) ReportSeizure(
    ctx context.Context,
    req domain.SeizureRequest,
    officerID uuid.UUID,
) (*domain.Seizure, error) {
    seizure := domain.NewSeizure(req, officerID)

    // Chercher si l arme est deja dans le registre (par numero de serie)
    if req.SerialNumber != "" {
        existing, _ := s.repo.FindBySerial(ctx, req.SerialNumber)
        if existing != nil {
            seizure.FirearmID = &existing.FirearmID
            existing.Status = domain.StatusSeized
            _ = s.repo.UpdateFirearm(ctx, existing)
        }
    }

    if err := s.repo.CreateSeizure(ctx, seizure); err != nil {
        return nil, err
    }

    // Soumettre a INTERPOL iARMS si arme pas encore enregistree
    if seizure.FirearmID == nil {
        _ = s.iarmsClient.SubmitStolenFirearm(ctx, seizure)
    }

    _ = s.kafka.Publish(ctx, "siar.seizure.reported", seizure)
    return seizure, nil
}
```

---

## 4. API REST

| Méthode | Endpoint                          | Rôle             | Description                     |
|---------|-----------------------------------|------------------|---------------------------------|
| `POST`  | `/api/v1/siar/firearms`           | PNH, MJSP        | Enregistrer arme                |
| `GET`   | `/api/v1/siar/firearms/:id`       | PNH, DCPJ        | Détail arme                     |
| `GET`   | `/api/v1/siar/check/serial/:sn`   | PNH_OFFICER      | Vérifier numéro de série        |
| `POST`  | `/api/v1/siar/seizures`           | PNH_OFFICER      | Signaler saisie d'arme          |
| `POST`  | `/api/v1/siar/stolen`             | PNH_OFFICER      | Déclarer vol d'arme             |
| `GET`   | `/api/v1/siar/licenses/:person`   | PNH, MJSP        | Licences d'une personne         |
| `POST`  | `/api/v1/siar/licenses`           | MJSP_ADMIN       | Créer licence                   |
| `GET`   | `/api/v1/siar/stats/by-type`      | DCPJ, MJSP       | Statistiques par type d'arme    |

---

## 5. VARIABLES D'ENVIRONNEMENT

```dotenv
SIAR_DB_HOST=localhost
SIAR_DB_NAME=snisid_siar
SIAR_IARMS_GATEWAY=https://i247-gateway.pnh.gov.ht/iarms
SIAR_ATF_ETRACE_URL=https://etrace.atf.gov/api
SIAR_GANG_SERVICE_URL=http://gang-svc:8095
SIAR_LICENSE_VALIDITY_YEARS=2
SIAR_SERVICE_PORT=8102
```

---
*MP-29 — SIAR-HT — Armes à Feu Légales — SNISID — République d'Haïti*
