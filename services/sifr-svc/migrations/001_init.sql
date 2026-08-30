BEGIN;
CREATE TYPE sifr_crossing_direction AS ENUM ('ENTRY', 'EXIT');
CREATE TYPE sifr_doc_type AS ENUM ('PASSPORT','NATIONAL_ID','LAISSEZ_PASSER','BIRTH_CERTIFICATE','TRAVEL_DOCUMENT','NONE');
CREATE TYPE sifr_alert_type AS ENUM ('WANTED_PERSON','STOLEN_DOCUMENT','BLACKLIST','ACTIVE_WARRANT','SANCTIONS','CUSTOMS_ALERT');
CREATE TABLE sifr_border_posts (
    post_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), post_code VARCHAR(10) UNIQUE NOT NULL,
    name VARCHAR(150) NOT NULL, dept_code CHAR(2) NOT NULL, border_country CHAR(3) NOT NULL DEFAULT 'DOM',
    post_lat DECIMAL(10,7), post_lng DECIMAL(10,7), is_official BOOLEAN DEFAULT TRUE,
    is_active BOOLEAN DEFAULT TRUE, lanes_count SMALLINT DEFAULT 2,
    has_biometric_scanner BOOLEAN DEFAULT FALSE, has_vehicle_scanner BOOLEAN DEFAULT FALSE,
    operating_hours VARCHAR(50), commanding_officer UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE sifr_crossings (
    crossing_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), post_id UUID NOT NULL REFERENCES sifr_border_posts(post_id),
    direction sifr_crossing_direction NOT NULL, crossing_datetime TIMESTAMPTZ NOT NULL,
    snisid_person_id UUID, document_type sifr_doc_type NOT NULL DEFAULT 'PASSPORT',
    document_number VARCHAR(100), document_country CHAR(3), document_expiry DATE,
    traveler_name VARCHAR(200) NOT NULL, traveler_dob DATE, traveler_nationality CHAR(3),
    vehicle_plate VARCHAR(20), lane_number SMALLINT, processing_officer UUID NOT NULL,
    alert_triggered BOOLEAN DEFAULT FALSE, alert_type sifr_alert_type,
    alert_action_taken TEXT, processing_time_sec INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE sifr_alerts_log (
    alert_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), crossing_id UUID REFERENCES sifr_crossings(crossing_id),
    post_id UUID NOT NULL, alert_type sifr_alert_type NOT NULL, snisid_person_id UUID,
    document_number VARCHAR(100), vehicle_plate VARCHAR(20), alert_source VARCHAR(50),
    source_record_id UUID, notified_units TEXT[] DEFAULT '{}', action_taken TEXT,
    resolved BOOLEAN DEFAULT FALSE, resolved_by UUID, resolved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE sifr_clandestine_crossings (
    report_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), location_desc VARCHAR(300),
    dept_code CHAR(2), lat DECIMAL(10,7), lng DECIMAL(10,7), reported_date TIMESTAMPTZ NOT NULL,
    crossing_type VARCHAR(50), estimated_persons INTEGER, gang_related BOOLEAN DEFAULT FALSE,
    gang_id UUID, trafficking_type TEXT, reported_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_sifr_crossings_datetime ON sifr_crossings(crossing_datetime DESC);
CREATE INDEX idx_sifr_crossings_person ON sifr_crossings(snisid_person_id) WHERE snisid_person_id IS NOT NULL;
CREATE INDEX idx_sifr_crossings_doc ON sifr_crossings(document_number);
CREATE INDEX idx_sifr_crossings_alert ON sifr_crossings(alert_triggered) WHERE alert_triggered = TRUE;
CREATE INDEX idx_sifr_crossings_post ON sifr_crossings(post_id, crossing_datetime DESC);
COMMIT;
