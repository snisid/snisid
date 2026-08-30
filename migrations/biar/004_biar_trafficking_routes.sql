BEGIN;

CREATE TABLE biar_trafficking_routes (
    route_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    origin_country      CHAR(3) NOT NULL,
    transit_countries   CHAR(3)[] DEFAULT '{}',
    destination_country CHAR(3) NOT NULL DEFAULT 'HTI',
    import_method       VARCHAR(50),
    weapon_types        TEXT[] DEFAULT '{}',
    total_weapons       INTEGER DEFAULT 0,
    first_detected      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_detected       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    active              BOOLEAN DEFAULT TRUE,
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
