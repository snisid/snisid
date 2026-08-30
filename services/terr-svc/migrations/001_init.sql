BEGIN;

CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS postgis_topology;

CREATE TYPE terr_control_level AS ENUM (
    'FULL_CONTROL', 'STRONG_INFLUENCE', 'CONTESTED',
    'WEAK_INFLUENCE', 'STATE_CONTROLLED', 'NO_MAN_LAND'
);

CREATE TYPE terr_source AS ENUM (
    'PNH_FIELD_REPORT', 'SATELLITE_ANALYSIS', 'INFORMANT',
    'NGO_REPORT', 'ACLED', 'MEDIA_CROSS_CHECK', 'LAPI_ANALYSIS'
);

CREATE TABLE terr_zones (
    zone_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    gang_id             UUID NOT NULL,
    zone_name           VARCHAR(150),
    dept_code           CHAR(2) NOT NULL,
    commune             VARCHAR(100),
    section_communale   VARCHAR(100),
    geom                GEOMETRY(MultiPolygon, 4326) NOT NULL,
    area_km2            DECIMAL(10,3),
    centroid_lat        DECIMAL(10,7),
    centroid_lng        DECIMAL(10,7),
    control_level       terr_control_level NOT NULL,
    estimated_population INTEGER,
    strategic_importance SMALLINT CHECK (strategic_importance BETWEEN 1 AND 10),
    controls_national_road BOOLEAN DEFAULT FALSE,
    road_numbers        TEXT[] DEFAULT '{}',
    controls_port       BOOLEAN DEFAULT FALSE,
    controls_airport    BOOLEAN DEFAULT FALSE,
    controls_market     BOOLEAN DEFAULT FALSE,
    valid_from          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    valid_to            TIMESTAMPTZ,
    is_current          BOOLEAN DEFAULT TRUE,
    intelligence_source terr_source NOT NULL,
    confidence_level    SMALLINT CHECK (confidence_level BETWEEN 1 AND 10),
    analyst_notes       TEXT,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION terr_compute_geometry_props()
RETURNS TRIGGER AS $$
BEGIN
    NEW.area_km2 := ST_Area(NEW.geom::geography) / 1000000;
    NEW.centroid_lat := ST_Y(ST_Centroid(NEW.geom));
    NEW.centroid_lng := ST_X(ST_Centroid(NEW.geom));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_terr_geometry
    BEFORE INSERT OR UPDATE ON terr_zones
    FOR EACH ROW EXECUTE FUNCTION terr_compute_geometry_props();

CREATE TABLE terr_zone_history (
    history_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    zone_id             UUID NOT NULL REFERENCES terr_zones(zone_id),
    change_type         VARCHAR(30) NOT NULL,
    previous_control    terr_control_level,
    new_control         terr_control_level,
    change_date         TIMESTAMPTZ NOT NULL,
    trigger_event       TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE terr_checkpoints (
    checkpoint_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    gang_id             UUID NOT NULL,
    location            GEOMETRY(Point, 4326) NOT NULL,
    location_desc       VARCHAR(300),
    dept_code           CHAR(2),
    road_number         VARCHAR(20),
    is_armed            BOOLEAN DEFAULT TRUE,
    extortion_type      VARCHAR(100),
    reported_at         TIMESTAMPTZ NOT NULL,
    is_active           BOOLEAN DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_terr_zones_geom        ON terr_zones USING GIST(geom) WHERE is_current = TRUE;
CREATE INDEX idx_terr_zones_dept        ON terr_zones(dept_code) WHERE is_current = TRUE;
CREATE INDEX idx_terr_zones_gang        ON terr_zones(gang_id) WHERE is_current = TRUE;
CREATE INDEX idx_terr_zones_control     ON terr_zones(control_level) WHERE is_current = TRUE;
CREATE INDEX idx_terr_checkpoints_geom  ON terr_checkpoints USING GIST(location) WHERE is_active = TRUE;

COMMIT;
