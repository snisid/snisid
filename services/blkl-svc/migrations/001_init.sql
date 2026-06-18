CREATE TABLE IF NOT EXISTS blkl_blacklist (
    id UUID PRIMARY KEY,
    entry_id VARCHAR(50) UNIQUE NOT NULL,
    national_blkl_id VARCHAR(100),
    snisid_person_id UUID NOT NULL,
    restriction_type VARCHAR(30) NOT NULL,
    source VARCHAR(40) NOT NULL,
    source_record_id VARCHAR(100),
    reason TEXT NOT NULL,
    court_order_ref VARCHAR(100),
    ordered_by VARCHAR(200),
    effective_date TIMESTAMP NOT NULL,
    expiry_date TIMESTAMP,
    is_permanent BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    alert_level VARCHAR(20),
    armed_dangerous BOOLEAN DEFAULT false,
    created_by VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_blkl_person_id ON blkl_blacklist(snisid_person_id);
CREATE INDEX idx_blkl_active ON blkl_blacklist(is_active);
CREATE INDEX idx_blkl_expiry ON blkl_blacklist(expiry_date) WHERE is_active = true;

CREATE TABLE IF NOT EXISTS blkl_alerts_log (
    id UUID PRIMARY KEY,
    blacklist_id UUID NOT NULL REFERENCES blkl_blacklist(id),
    person_id UUID NOT NULL,
    alert_type VARCHAR(50) NOT NULL,
    message TEXT NOT NULL,
    acknowledged BOOLEAN DEFAULT false,
    acknowledged_by VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_blkl_alerts_person ON blkl_alerts_log(person_id);
