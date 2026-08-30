CREATE TYPE sisal_hazard_type AS ENUM (
    'EARTHQUAKE', 'HURRICANE', 'FLOOD', 'TSUNAMI', 'LANDSLIDE',
    'SECURITY_GANG', 'SECURITY_MASS_CASUALTY', 'EPIDEMIC',
    'INDUSTRIAL', 'COMPOSITE'
);

CREATE TYPE sisal_severity AS ENUM (
    'ADVISORY', 'WATCH', 'WARNING', 'EMERGENCY', 'CATASTROPHE'
);

CREATE TYPE sisal_channel AS ENUM (
    'SMS_MASS', 'PUSH_NOTIFICATION', 'RADIO_BROADCAST',
    'SIRENE', 'SOCIAL_MEDIA', 'OFFICIAL_AGENCIES', 'LOUD_SPEAKER'
);

CREATE TABLE sisal_alerts (
    alert_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_sisal_id   VARCHAR(25) UNIQUE NOT NULL,
    hazard_type         sisal_hazard_type NOT NULL,
    severity            sisal_severity NOT NULL,
    title               VARCHAR(200) NOT NULL,
    message_fr          TEXT NOT NULL,
    message_ht          TEXT NOT NULL,
    affected_depts      CHAR(2)[] DEFAULT '{}',
    affected_communes   TEXT[] DEFAULT '{}',
    affected_pop_est    INTEGER,
    issued_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    valid_until         TIMESTAMPTZ,
    source_agency       VARCHAR(100) NOT NULL,
    source_event_id     UUID,
    recommended_actions TEXT[],
    channels_used       sisal_channel[] DEFAULT '{}',
    sms_count_sent      INTEGER DEFAULT 0,
    push_count_sent     INTEGER DEFAULT 0,
    is_cancelled        BOOLEAN DEFAULT FALSE,
    cancelled_at        TIMESTAMPTZ,
    cancel_reason       TEXT,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sisal_data_feeds (
    feed_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    feed_name           VARCHAR(100) NOT NULL,
    feed_url            VARCHAR(500),
    feed_type           VARCHAR(30),
    hazard_types        sisal_hazard_type[] DEFAULT '{}',
    polling_interval_sec INTEGER DEFAULT 60,
    is_active           BOOLEAN DEFAULT TRUE,
    last_poll           TIMESTAMPTZ,
    last_alert_generated TIMESTAMPTZ,
    alert_thresholds    JSONB,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sisal_subscriptions (
    sub_id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    snisid_person_id    UUID,
    phone_number        VARCHAR(30),
    email               VARCHAR(200),
    dept_code           CHAR(2),
    commune             VARCHAR(100),
    hazard_types        sisal_hazard_type[] DEFAULT '{}',
    min_severity        sisal_severity DEFAULT 'WARNING',
    channels            sisal_channel[] DEFAULT ARRAY['SMS_MASS'::sisal_channel],
    is_active           BOOLEAN DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sisal_alerts_dept ON sisal_alerts USING gin(affected_depts);
CREATE INDEX idx_sisal_alerts_severity ON sisal_alerts(severity, issued_at DESC);
CREATE INDEX idx_sisal_alerts_hazard ON sisal_alerts(hazard_type, issued_at DESC);
CREATE INDEX idx_sisal_subs_dept ON sisal_subscriptions(dept_code) WHERE is_active = TRUE;
