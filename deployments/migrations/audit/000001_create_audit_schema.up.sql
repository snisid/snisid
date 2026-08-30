CREATE TABLE audit_events (
    event_id VARCHAR(50) PRIMARY KEY,
    correlation_id VARCHAR(100),
    event_type VARCHAR(100) NOT NULL,
    actor VARCHAR(100),
    action VARCHAR(100),
    resource VARCHAR(100),
    status VARCHAR(50),
    payload TEXT NOT NULL,
    previous_hash VARCHAR(100) NOT NULL,
    hash VARCHAR(100) NOT NULL,
    sequence_id BIGSERIAL UNIQUE NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_correlation_id ON audit_events(correlation_id);
CREATE INDEX idx_audit_actor ON audit_events(actor);
CREATE INDEX idx_audit_resource ON audit_events(resource);
CREATE INDEX idx_audit_sequence ON audit_events(sequence_id);
