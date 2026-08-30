CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE sigeo_admin_boundaries (
    boundary_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    level               SMALLINT NOT NULL,
    code                VARCHAR(20) UNIQUE NOT NULL,
    name                VARCHAR(150) NOT NULL,
    parent_code         VARCHAR(20),
    geom                GEOMETRY(MultiPolygon, 4326) NOT NULL,
    population_est      INTEGER,
    area_km2            DECIMAL(10,3),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sigeo_incidents_unified (
    event_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_module       VARCHAR(20) NOT NULL,
    source_record_id    UUID NOT NULL,
    event_type          VARCHAR(50) NOT NULL,
    event_date          TIMESTAMPTZ NOT NULL,
    lat                 DECIMAL(10,7) NOT NULL,
    lng                 DECIMAL(10,7) NOT NULL,
    geom                GEOMETRY(Point, 4326) GENERATED ALWAYS AS (
                            ST_SetSRID(ST_MakePoint(lng, lat), 4326)
                        ) STORED,
    dept_code           CHAR(2),
    commune             VARCHAR(100),
    h3_index_8          VARCHAR(20),
    h3_index_10         VARCHAR(20),
    severity            SMALLINT CHECK (severity BETWEEN 1 AND 10),
    gang_id             UUID,
    description         TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sigeo_checkpoints (
    cp_id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cp_type             VARCHAR(30) NOT NULL,
    location            GEOMETRY(Point, 4326) NOT NULL,
    dept_code           CHAR(2),
    road_number         VARCHAR(10),
    description         VARCHAR(300),
    controlling_gang_id UUID,
    is_active           BOOLEAN DEFAULT TRUE,
    source_module       VARCHAR(20),
    source_record_id    UUID,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sigeo_incidents_geom ON sigeo_incidents_unified USING GIST(geom);
CREATE INDEX idx_sigeo_incidents_date ON sigeo_incidents_unified(event_date DESC);
CREATE INDEX idx_sigeo_incidents_dept ON sigeo_incidents_unified(dept_code, event_date DESC);
CREATE INDEX idx_sigeo_cp_geom ON sigeo_checkpoints USING GIST(location) WHERE is_active = TRUE;
