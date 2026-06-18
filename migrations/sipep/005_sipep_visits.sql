BEGIN;

CREATE TABLE sipep_visits (
    visit_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inmate_id     UUID NOT NULL REFERENCES sipep_inmates(inmate_id),
    visitor_name  VARCHAR(200) NOT NULL,
    visitor_id    VARCHAR(50),
    relationship  VARCHAR(50),
    visit_date    TIMESTAMPTZ NOT NULL,
    check_in      TIMESTAMPTZ,
    check_out     TIMESTAMPTZ,
    authorized_by UUID,
    notes         TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sipep_health_events (
    event_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inmate_id         UUID NOT NULL REFERENCES sipep_inmates(inmate_id),
    event_type        VARCHAR(50) NOT NULL,
    event_date        TIMESTAMPTZ NOT NULL,
    description       TEXT,
    treating_facility VARCHAR(150),
    outcome           TEXT,
    reported_by       UUID,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
