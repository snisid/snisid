CREATE TABLE identities (
    id VARCHAR(50) PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    dob DATE,
    gender VARCHAR(20),
    agency VARCHAR(100),
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE identity_histories (
    history_id VARCHAR(50) PRIMARY KEY,
    identity_id VARCHAR(50) NOT NULL REFERENCES identities(id),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    dob DATE,
    gender VARCHAR(20),
    agency VARCHAR(100),
    status VARCHAR(50),
    version INT NOT NULL,
    changed_at TIMESTAMP WITH TIME ZONE NOT NULL,
    changed_by VARCHAR(100),
    reason TEXT
);

CREATE TABLE biometric_references (
    id VARCHAR(50) PRIMARY KEY,
    identity_id VARCHAR(50) NOT NULL REFERENCES identities(id),
    type VARCHAR(50) NOT NULL,
    reference_uri TEXT NOT NULL,
    quality_score FLOAT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE document_associations (
    id VARCHAR(50) PRIMARY KEY,
    identity_id VARCHAR(50) NOT NULL REFERENCES identities(id),
    document_type VARCHAR(50) NOT NULL,
    document_uri TEXT NOT NULL,
    verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_identities_agency ON identities(agency);
CREATE INDEX idx_identities_status ON identities(status);
CREATE INDEX idx_identity_histories_identity_id ON identity_histories(identity_id);
