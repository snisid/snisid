BEGIN;

CREATE TABLE mar_incidents (
    incident_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vessel_id           UUID REFERENCES mar_vessels(vessel_id),
    incident_type       mar_incident_type NOT NULL,
    incident_date       TIMESTAMPTZ NOT NULL,
    lat                 DECIMAL(10,7),
    lng                 DECIMAL(10,7),
    zone_desc           VARCHAR(100),
    responding_unit     VARCHAR(50),
    outcome             TEXT,
    persons_involved    INTEGER DEFAULT 0,
    snisid_person_ids   UUID[] DEFAULT '{}',
    drug_types          TEXT[] DEFAULT '{}',
    drug_weight_kg      DECIMAL(12,3),
    weapons_found       BOOLEAN DEFAULT FALSE,
    weapons_count       INTEGER DEFAULT 0,
    migrants_count      INTEGER DEFAULT 0,
    biar_refs           UUID[] DEFAULT '{}',
    case_reference      VARCHAR(100),
    photo_refs          TEXT[] DEFAULT '{}',
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
