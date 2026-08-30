BEGIN;

CREATE TABLE chef_intelligence_notes (
    note_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    member_id           UUID NOT NULL REFERENCES chef_criminal_members(member_id),
    note_date           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    intel_type          VARCHAR(50),
    content             TEXT NOT NULL,
    source_classif      VARCHAR(20),
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
