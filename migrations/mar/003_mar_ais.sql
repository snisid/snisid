BEGIN;

CREATE TABLE mar_ais_sightings (
    sighting_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vessel_id           UUID REFERENCES mar_vessels(vessel_id),
    mmsi                VARCHAR(15),
    vessel_name         VARCHAR(150),
    sighting_timestamp  TIMESTAMPTZ NOT NULL,
    lat                 DECIMAL(10,7) NOT NULL,
    lng                 DECIMAL(10,7) NOT NULL,
    speed_knots         DECIMAL(5,2),
    heading_degrees     SMALLINT,
    destination         VARCHAR(100),
    source_type         VARCHAR(30),
    zone_code           VARCHAR(20),
    alert_triggered     BOOLEAN DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
