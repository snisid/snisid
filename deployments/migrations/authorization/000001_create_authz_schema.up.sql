CREATE TABLE policies (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    module TEXT NOT NULL,
    enabled BOOLEAN DEFAULT TRUE,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE role_grants (
    id VARCHAR(50) PRIMARY KEY,
    role VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL,
    resource VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_policies_enabled ON policies(enabled);
CREATE INDEX idx_role_grants_role ON role_grants(role);

-- Insert base RBAC policy snippet as an initial seed
INSERT INTO policies (id, name, module) VALUES 
('pol-1', 'snisid.abac', 'package snisid.abac\n\ndefault allow = false\nallow {\n    input.action == "read"\n    input.user.roles[_] == "admin"\n}');
