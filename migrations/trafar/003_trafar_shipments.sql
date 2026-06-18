BEGIN;

CREATE TABLE trafar_shipments (
    shipment_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route_id            UUID REFERENCES trafar_routes(route_id),
    shipment_date       TIMESTAMPTZ NOT NULL,
    intercepted         BOOLEAN DEFAULT FALSE,
    interception_date   TIMESTAMPTZ,
    interception_location VARCHAR(300),
    interception_unit   VARCHAR(50),
    weapons_count       INTEGER,
    weapons_types       TEXT[] DEFAULT '{}',
    estimated_value_usd DECIMAL(12,2),
    linked_persons      UUID[] DEFAULT '{}',
    port_ht_ref         UUID,
    mar_ht_ref          UUID,
    case_reference      VARCHAR(100),
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
