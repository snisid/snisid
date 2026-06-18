BEGIN;

CREATE TABLE rdep_flights (
    flight_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flight_number         VARCHAR(20) NOT NULL,
    flight_type           rdep_flight_type NOT NULL,
    origin_country        rdep_deportation_country NOT NULL,
    departure_airport     VARCHAR(100) NOT NULL,
    arrival_airport       VARCHAR(100) NOT NULL,           -- PAP, CAP, etc.
    departure_time        TIMESTAMPTZ NOT NULL,
    arrival_time          TIMESTAMPTZ NOT NULL,
    deporting_agency      VARCHAR(100),
    total_passengers      INTEGER NOT NULL DEFAULT 0,
    manifest_ref          VARCHAR(200),
    notes                 TEXT,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
