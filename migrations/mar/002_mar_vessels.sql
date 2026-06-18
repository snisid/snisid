BEGIN;

CREATE TABLE mar_vessels (
    vessel_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_mar_id     VARCHAR(25) UNIQUE NOT NULL,
    vessel_name         VARCHAR(150),
    imo_number          VARCHAR(20),
    mmsi                VARCHAR(15),
    call_sign           VARCHAR(15),
    vessel_type         mar_vessel_type NOT NULL,
    flag_country        CHAR(3),
    hull_color          VARCHAR(50),
    length_m            DECIMAL(8,2),
    tonnage_gt          INTEGER,
    engine_count        SMALLINT,
    horsepower          INTEGER,
    owner_name          VARCHAR(200),
    owner_snisid_id     UUID,
    registration_number VARCHAR(50),
    registration_port   VARCHAR(100),
    status              mar_vessel_status NOT NULL DEFAULT 'REGISTERED',
    gang_id             UUID,
    interpol_svd_ref    VARCHAR(50),
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
