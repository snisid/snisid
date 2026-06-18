CREATE TYPE mvsm_event_type AS ENUM (
    'POLITICAL_PROTEST', 'LABOR_STRIKE', 'COMMUNITY_ACTION',
    'RELIGIOUS_GATHERING', 'CULTURAL_EVENT', 'PEYI_LOK_BARRICADE',
    'GANG_MOBILIZATION', 'SPONTANEOUS_UNREST', 'OTHER'
);

CREATE TYPE mvsm_risk_level AS ENUM (
    'LOW', 'MODERATE', 'HIGH', 'CRITICAL'
);

CREATE TABLE mvsm_events (
    event_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_mvsm_id    VARCHAR(25) UNIQUE NOT NULL,
    event_type          mvsm_event_type NOT NULL,
    event_name          VARCHAR(200),
    risk_level          mvsm_risk_level NOT NULL DEFAULT 'LOW',
    status              VARCHAR(20) DEFAULT 'PLANNED',
    organizer_name      VARCHAR(200),
    organizer_snisid_id UUID,
    gang_id             UUID,
    scheduled_date      TIMESTAMPTZ NOT NULL,
    actual_start        TIMESTAMPTZ,
    actual_end          TIMESTAMPTZ,
    location_desc       VARCHAR(300),
    dept_code           CHAR(2),
    commune             VARCHAR(100),
    lat                 DECIMAL(10,7),
    lng                 DECIMAL(10,7),
    estimated_crowd     INTEGER,
    peak_crowd          INTEGER,
    deployed_units      TEXT[] DEFAULT '{}',
    incidents_during    INTEGER DEFAULT 0,
    casualties          INTEGER DEFAULT 0,
    arrests_made        INTEGER DEFAULT 0,
    weapons_found       INTEGER DEFAULT 0,
    vehicles_involved   INTEGER DEFAULT 0,
    sivc_alert_ids      UUID[] DEFAULT '{}',
    post_event_notes    TEXT,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE mvsm_real_time_updates (
    update_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id            UUID NOT NULL REFERENCES mvsm_events(event_id),
    update_time         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    current_crowd_est   INTEGER,
    situation           TEXT NOT NULL,
    risk_change         mvsm_risk_level,
    action_taken        TEXT,
    reported_by         UUID NOT NULL,
    lat                 DECIMAL(10,7),
    lng                 DECIMAL(10,7)
);

CREATE INDEX idx_mvsm_events_date ON mvsm_events(scheduled_date DESC);
CREATE INDEX idx_mvsm_events_dept ON mvsm_events(dept_code, scheduled_date DESC);
CREATE INDEX idx_mvsm_events_risk ON mvsm_events(risk_level) WHERE status IN ('PLANNED','ACTIVE');
CREATE INDEX idx_mvsm_updates_event ON mvsm_real_time_updates(event_id, update_time DESC);
