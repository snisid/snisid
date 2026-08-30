# MP-25 — CHEF-HT
## Fichier des Leaders et Membres Identifiés d'Organisations Criminelles
### Master Prompt d'Implémentation — SNISID

```
Classification   : SECRET / USAGE OFFICIEL EXCLUSIF
Module           : MP-25 | Code : CHEF-HT
Dépendances      : GANG-HT (MP-24), FIR-HT (MP-20), AFIS-HT (MP-19), RDEP-HT (MP-22)
Normes           : INTERPOL Organized Crime Notices, Résolution CSNU 2653, OFAC, DHS
Acteurs          : DCPJ BAC, Cellule Intelligence, DEA liaison, Panel experts ONU Haïti
```

---

## 1. CONTEXTE

Les leaders de gangs haïtiens sous sanctions internationales identifiées incluent :
- **Jimmy Chérizier alias "Barbecue"** — Chef de Viv Ansanm / ex-G9, désigné OFAC
- **Gabriel Jean-Pierre alias "Ti Gabriel"** — 400 Mawozo, désigné OFAC
- **Lanmo San Jou** — G9 an Fanm (Cité Soleil)

Ce module indexe tous les membres identifiés de toutes organisations criminelles avec
rôle hiérarchique, biographie criminelle, contacts, statut sanctions et intelligence terrain.

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE chef_role_type AS ENUM (
    'SUPREME_LEADER','ZONE_COMMANDER','LIEUTENANT',
    'SOLDIER','ASSOCIATE','FINANCIER','ENABLER','INFORMANT'
);

CREATE TYPE chef_status AS ENUM (
    'ACTIVE','ARRESTED','DETAINED','DECEASED','FLED_COUNTRY','UNKNOWN'
);

CREATE TABLE chef_criminal_members (
    member_id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_chef_id        VARCHAR(25) UNIQUE NOT NULL,  -- CHEF-HT-NNNNNN
    snisid_person_id        UUID NOT NULL,
    fir_record_id           UUID,
    afis_subject_id         UUID,
    rdep_deportee_id        UUID,

    primary_gang_id         UUID NOT NULL,
    role_in_gang            chef_role_type NOT NULL,
    role_description        TEXT,
    joined_date             DATE,
    rank_level              SMALLINT,

    aliases                 TEXT[] DEFAULT '{}',
    known_languages         TEXT[] DEFAULT '{}',
    tattoo_description      TEXT,
    physical_description    TEXT,
    photo_refs              TEXT[] DEFAULT '{}',

    territory_dept          CHAR(2),
    territory_communes      TEXT[] DEFAULT '{}',

    known_armed             BOOLEAN DEFAULT FALSE,
    weapon_types            TEXT[] DEFAULT '{}',
    trained_combatant       BOOLEAN DEFAULT FALSE,

    status                  chef_status NOT NULL DEFAULT 'ACTIVE',
    un_designated           BOOLEAN DEFAULT FALSE,
    un_designation_date     TIMESTAMPTZ,
    ofac_designated         BOOLEAN DEFAULT FALSE,
    ofac_sdn_ref            VARCHAR(50),
    interpol_notice_ref     VARCHAR(50),

    last_known_address      TEXT,
    last_known_dept         CHAR(2),
    last_seen_at            TIMESTAMPTZ,
    last_seen_location      VARCHAR(300),

    estimated_wealth_usd    DECIMAL(15,2),
    known_assets            TEXT[] DEFAULT '{}',

    intel_classification    VARCHAR(20) DEFAULT 'SECRET',
    intel_confidence        SMALLINT CHECK (intel_confidence BETWEEN 1 AND 10),
    created_by              UUID NOT NULL,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE chef_intelligence_notes (
    note_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    member_id           UUID NOT NULL REFERENCES chef_criminal_members(member_id),
    note_date           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    intel_type          VARCHAR(50),  -- SIGHTING, COMM_INTERCEPT, INFORMANT, ARREST
    content             TEXT NOT NULL,
    source_classif      VARCHAR(20),
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE chef_cross_gang_links (
    link_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    member_a_id         UUID NOT NULL REFERENCES chef_criminal_members(member_id),
    member_b_id         UUID NOT NULL REFERENCES chef_criminal_members(member_id),
    link_type           VARCHAR(50),  -- FAMILY, ASSOCIATE, SUPPLIER, RIVAL
    confidence          SMALLINT,
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT no_self_link CHECK (member_a_id <> member_b_id)
);

CREATE TABLE chef_sightings (
    sighting_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    member_id           UUID NOT NULL REFERENCES chef_criminal_members(member_id),
    sighted_at          TIMESTAMPTZ NOT NULL,
    location_desc       VARCHAR(300),
    dept_code           CHAR(2),
    commune             VARCHAR(100),
    lat                 DECIMAL(10,7),
    lng                 DECIMAL(10,7),
    source_type         VARCHAR(30),  -- LAPI, INFORMANT, FIELD_REPORT, CHECKPOINT
    confidence          SMALLINT,
    photo_ref           VARCHAR(500),
    reported_by         UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_chef_gang        ON chef_criminal_members(primary_gang_id, role_in_gang);
CREATE INDEX idx_chef_status      ON chef_criminal_members(status) WHERE status = 'ACTIVE';
CREATE INDEX idx_chef_ofac        ON chef_criminal_members(ofac_designated) WHERE ofac_designated = TRUE;
CREATE INDEX idx_chef_un          ON chef_criminal_members(un_designated) WHERE un_designated = TRUE;
CREATE INDEX idx_chef_dept        ON chef_criminal_members(territory_dept) WHERE status = 'ACTIVE';
CREATE INDEX idx_chef_sightings   ON chef_sightings(member_id, sighted_at DESC);

COMMIT;
```

---

## 3. SERVICE GO CLÉ

```go
package service

import (
    "context"
    "github.com/google/uuid"
    "github.com/snisid/chef-svc/internal/domain"
)

func (s *MemberService) UpdateStatus(
    ctx context.Context,
    memberID uuid.UUID,
    newStatus domain.ChefStatus,
    updatedBy uuid.UUID,
    notes string,
) error {
    member, err := s.repo.FindByID(ctx, memberID)
    if err != nil {
        return err
    }
    old := member.Status
    member.Status = newStatus
    if err := s.repo.Update(ctx, member); err != nil {
        return err
    }

    // Publier changement de statut
    _ = s.kafka.Publish(ctx, "chef.status.changed", domain.StatusChangedEvent{
        MemberID:   memberID,
        GangID:     member.PrimaryGangID,
        OldStatus:  old,
        NewStatus:  newStatus,
        ChangedBy:  updatedBy,
        Notes:      notes,
    })

    // Si ARRESTED -> notifier SIPEP pour enrôlement
    if newStatus == domain.StatusArrested {
        _ = s.kafka.Publish(ctx, "chef.arrested", domain.MemberArrestedEvent{
            MemberID:  memberID,
            PersonID:  member.SNISIDPersonID,
            GangID:    member.PrimaryGangID,
            RoleInGang: member.RoleInGang,
        })
    }
    return nil
}
```

---

## 4. API REST

| Méthode | Endpoint                              | Rôle            | Description                     |
|---------|---------------------------------------|-----------------|---------------------------------|
| `POST`  | `/api/v1/chef/members`                | DCPJ_INTEL      | Créer fiche membre              |
| `GET`   | `/api/v1/chef/members/:id`            | DCPJ, BAC       | Profil complet membre           |
| `GET`   | `/api/v1/chef/members/by-gang/:id`    | DCPJ, BAC       | Membres d'un gang               |
| `GET`   | `/api/v1/chef/members/sanctioned`     | DCPJ, MJSP      | Membres sous sanctions ONU/OFAC |
| `POST`  | `/api/v1/chef/members/:id/intel`      | DCPJ_INTEL      | Ajouter note intelligence       |
| `POST`  | `/api/v1/chef/members/:id/sightings`  | PNH_OFFICER     | Enregistrer observation         |
| `PATCH` | `/api/v1/chef/members/:id/status`     | DCPJ_SUPERVISOR | Changer statut                  |
| `GET`   | `/api/v1/chef/network/:id`            | DCPJ_INTEL      | Réseau connexions du membre     |
| `GET`   | `/api/v1/chef/members/leaders`        | DCPJ            | Tous les chefs actifs           |

---

## 5. INTÉGRATIONS

- **GANG-HT** : `primary_gang_id` → structure organisationnelle complète
- **AFIS-HT** : empreintes obligatoires à la création de chaque fiche
- **FPR-HT** : membres ACTIVE → alerte automatique si LAPI capte plaque liée
- **SIVC-HT** : véhicules connus du membre → surveillance LAPI
- **SANC-HT** : `ofac_sdn_ref` / `un_designated` → vérification croisée quotidienne
- **RDEP-HT** : `rdep_deportee_id` → lien antécédents étrangers

---

## 6. VARIABLES D'ENVIRONNEMENT

```dotenv
CHEF_DB_HOST=localhost
CHEF_DB_NAME=snisid_chef
CHEF_NEO4J_URI=bolt://neo4j:7687
CHEF_GANG_SERVICE_URL=http://gang-svc:8095
CHEF_AFIS_SERVICE_URL=http://afis-svc:8091
CHEF_FIR_SERVICE_URL=http://fir-svc:8093
CHEF_KAFKA_BROKERS=kafka:9092
CHEF_SERVICE_PORT=8097
```

---
*MP-25 — CHEF-HT — Leaders Criminels — SNISID — République d'Haïti*
