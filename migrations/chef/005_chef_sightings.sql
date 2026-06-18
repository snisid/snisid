BEGIN;

CREATE TABLE chef_sightings (
    sighting_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    member_id           UUID NOT NULL REFERENCES chef_criminal_members(member_id),
    sighted_at          TIMESTAMPTZ NOT NULL,
    location_desc       VARCHAR(300),
    dept_code           CHAR(2),
    commune             VARCHAR(100),
    lat                 DECIMAL(10,7),
    lng                 DECIMAL(10,7),
    source_type         VARCHAR(30),
    confidence          SMALLINT,
    photo_ref           VARCHAR(500),
    reported_by         UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
