CREATE TYPE sipci_asset_category AS ENUM (
    'ENERGY', 'TRANSPORT', 'WATER', 'TELECOMS', 'HEALTH',
    'FINANCE', 'GOVERNMENT', 'EDUCATION', 'FOOD_SUPPLY'
);

CREATE TYPE sipci_threat_level AS ENUM (
    'NORMAL', 'ELEVATED', 'HIGH', 'SEVERE', 'CRITICAL'
);

CREATE TABLE sipci_assets (
    asset_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_sipci_id   VARCHAR(25) UNIQUE NOT NULL,
    asset_name          VARCHAR(200) NOT NULL,
    asset_category      sipci_asset_category NOT NULL,
    owner_entity        VARCHAR(200),
    operating_org       VARCHAR(200),
    location_desc       VARCHAR(300),
    dept_code           CHAR(2) NOT NULL,
    commune             VARCHAR(100),
    lat                 DECIMAL(10,7) NOT NULL,
    lng                 DECIMAL(10,7) NOT NULL,
    criticality_score   SMALLINT CHECK (criticality_score BETWEEN 1 AND 10),
    population_served   INTEGER,
    dependency_assets   UUID[] DEFAULT '{}',
    single_point_failure BOOLEAN DEFAULT FALSE,
    current_threat_level sipci_threat_level NOT NULL DEFAULT 'NORMAL',
    is_in_gang_zone     BOOLEAN DEFAULT FALSE,
    controlling_gang_id UUID,
    under_extortion     BOOLEAN DEFAULT FALSE,
    extors_case_id      UUID,
    incident_count_12m  INTEGER DEFAULT 0,
    last_incident_date  TIMESTAMPTZ,
    protection_unit     VARCHAR(50),
    security_guards     INTEGER DEFAULT 0,
    has_cctv            BOOLEAN DEFAULT FALSE,
    cctv_count          SMALLINT DEFAULT 0,
    has_perimeter       BOOLEAN DEFAULT FALSE,
    has_backup_power    BOOLEAN DEFAULT FALSE,
    site_manager_name   VARCHAR(200),
    site_manager_phone  VARCHAR(30),
    emergency_contact   VARCHAR(200),
    notes               TEXT,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sipci_incidents (
    incident_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id            UUID NOT NULL REFERENCES sipci_assets(asset_id),
    incident_type       VARCHAR(50) NOT NULL,
    incident_date       TIMESTAMPTZ NOT NULL,
    perpetrator_type    VARCHAR(30),
    gang_id             UUID,
    description         TEXT NOT NULL,
    impact_severity     SMALLINT CHECK (impact_severity BETWEEN 1 AND 10),
    population_affected INTEGER,
    service_disruption_hours DECIMAL(8,2),
    economic_loss_usd   DECIMAL(15,2),
    responding_units    TEXT[] DEFAULT '{}',
    sivc_alert_ids      UUID[] DEFAULT '{}',
    case_reference      VARCHAR(100),
    resolution_date     TIMESTAMPTZ,
    resolution_notes    TEXT,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sipci_assets_category ON sipci_assets(asset_category, current_threat_level);
CREATE INDEX idx_sipci_assets_dept ON sipci_assets(dept_code);
CREATE INDEX idx_sipci_assets_gang ON sipci_assets(is_in_gang_zone) WHERE is_in_gang_zone = TRUE;
CREATE INDEX idx_sipci_assets_critical ON sipci_assets(criticality_score DESC) WHERE single_point_failure = TRUE;
CREATE INDEX idx_sipci_incidents_date ON sipci_incidents(incident_date DESC, asset_id);
