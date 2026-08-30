CREATE TABLE user_credentials (
    user_id VARCHAR(50) PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    mfa_enabled BOOLEAN DEFAULT FALSE,
    mfa_secret TEXT,
    roles TEXT NOT NULL,
    locked_until TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE webauthn_credentials (
    id BYTEA PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL REFERENCES user_credentials(user_id),
    public_key BYTEA NOT NULL,
    attestation_type VARCHAR(50),
    sign_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_credentials_username ON user_credentials(username);
CREATE INDEX idx_webauthn_credentials_userid ON webauthn_credentials(user_id);
