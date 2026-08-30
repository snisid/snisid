-- SNISID Sovereign Identity Master Ledger
-- Powered by CockroachDB Distributed SQL

CREATE TABLE IF NOT EXISTS identity_events (
    event_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregate_id UUID NOT NULL,
    aggregate_type VARCHAR(50) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    event_data JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    actor_id VARCHAR(100) NOT NULL,
    signature VARCHAR(256) NOT NULL,
    INDEX idx_aggregate (aggregate_id, created_at ASC)
);

CREATE TABLE IF NOT EXISTS citizens_current_state (
    niu UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL,
    place_of_birth VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    biometric_hash VARCHAR(256) UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_event_id UUID REFERENCES identity_events(event_id),
    INDEX idx_status (status)
);

-- Audit Trigger to maintain updated_at
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_citizen_state
BEFORE UPDATE ON citizens_current_state
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
