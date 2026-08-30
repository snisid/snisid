BEGIN;

CREATE TABLE gang_territories (
    territory_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    gang_id      UUID NOT NULL REFERENCES gang_organizations(gang_id) ON DELETE CASCADE,
    dept_code    CHAR(2) NOT NULL,
    commune      VARCHAR(100) NOT NULL,
    locality     VARCHAR(200),
    geojson      JSONB,
    is_claimed   BOOLEAN DEFAULT TRUE,
    is_contested BOOLEAN DEFAULT FALSE,
    contested_with UUID[] DEFAULT '{}',
    controlled_since DATE,
    notes        TEXT,
    created_by   UUID NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(gang_id, dept_code, commune)
);

COMMIT;
