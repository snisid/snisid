BEGIN;

CREATE TABLE sifr_border_posts (
    post_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_code           VARCHAR(10) UNIQUE NOT NULL,
    name                VARCHAR(150) NOT NULL,
    dept_code           CHAR(2) NOT NULL,
    border_country      CHAR(3) NOT NULL DEFAULT 'DOM',
    post_lat            DECIMAL(10,7),
    post_lng            DECIMAL(10,7),
    is_official         BOOLEAN DEFAULT TRUE,
    is_active           BOOLEAN DEFAULT TRUE,
    lanes_count         SMALLINT DEFAULT 2,
    has_biometric_scanner BOOLEAN DEFAULT FALSE,
    has_vehicle_scanner BOOLEAN DEFAULT FALSE,
    operating_hours     VARCHAR(50),
    commanding_officer  UUID,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
