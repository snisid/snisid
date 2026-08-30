BEGIN;

CREATE TABLE sivc_intelligence_reports (
    report_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_number       VARCHAR(50) UNIQUE NOT NULL,
    title               VARCHAR(250) NOT NULL,
    report_type         VARCHAR(30) NOT NULL,
    classification      VARCHAR(20) NOT NULL DEFAULT 'RESTRICTED',
    summary             TEXT NOT NULL,
    full_report         JSONB,
    alert_ids           UUID[] DEFAULT '{}',
    plate_ids           UUID[] DEFAULT '{}',
    person_ids          UUID[] DEFAULT '{}',
    originating_unit    VARCHAR(50) NOT NULL,
    author_id           UUID NOT NULL,
    recipient_units     TEXT[] DEFAULT '{}',
    published_at        TIMESTAMPTZ,
    expiry_date         TIMESTAMPTZ,
    attachments         TEXT[] DEFAULT '{}',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sivc_intel_type   ON sivc_intelligence_reports(report_type, classification);
CREATE INDEX idx_sivc_intel_unit   ON sivc_intelligence_reports(originating_unit);
CREATE INDEX idx_sivc_intel_alerts ON sivc_intelligence_reports USING gin(alert_ids);

COMMIT;
