BEGIN;

-- Deportees indexes
CREATE INDEX idx_rdep_deportees_person      ON rdep_deportees(snisid_person_id);
CREATE INDEX idx_rdep_deportees_country     ON rdep_deportees(deportation_country);
CREATE INDEX idx_rdep_deportees_risk        ON rdep_deportees(criminal_risk_level) WHERE criminal_risk_level IN ('HIGH','VERY_HIGH');
CREATE INDEX idx_rdep_deportees_gang        ON rdep_deportees(gang_affiliated) WHERE gang_affiliated = TRUE;
CREATE INDEX idx_rdep_deportees_monitoring  ON rdep_deportees(monitoring_status) WHERE monitoring_required = TRUE;
CREATE INDEX idx_rdep_deportees_flight      ON rdep_deportees(flight_id);
CREATE INDEX idx_rdep_deportees_date        ON rdep_deportees(deportation_date DESC);

-- Extraditions indexes
CREATE INDEX idx_rdep_extraditions_person   ON rdep_extraditions(snisid_person_id);
CREATE INDEX idx_rdep_extraditions_status   ON rdep_extraditions(extradition_status);
CREATE INDEX idx_rdep_extraditions_country  ON rdep_extraditions(requesting_country);

-- Flights indexes
CREATE INDEX idx_rdep_flights_number        ON rdep_flights(flight_number);
CREATE INDEX idx_rdep_flights_date          ON rdep_flights(arrival_time DESC);
CREATE INDEX idx_rdep_flights_origin        ON rdep_flights(origin_country);

-- Foreign records indexes
CREATE INDEX idx_rdep_foreign_records_deportee ON rdep_foreign_records(deportee_id);
CREATE INDEX idx_rdep_foreign_records_country  ON rdep_foreign_records(country);

-- Monitoring events indexes
CREATE INDEX idx_rdep_monitoring_deportee   ON rdep_monitoring_events(deportee_id);
CREATE INDEX idx_rdep_monitoring_date       ON rdep_monitoring_events(event_date DESC);

-- Row-Level Security for multi-tenant access
ALTER TABLE rdep_deportees ENABLE ROW LEVEL SECURITY;
ALTER TABLE rdep_foreign_records ENABLE ROW LEVEL SECURITY;
ALTER TABLE rdep_monitoring_events ENABLE ROW LEVEL SECURITY;
ALTER TABLE rdep_extraditions ENABLE ROW LEVEL SECURITY;
ALTER TABLE rdep_flights ENABLE ROW LEVEL SECURITY;

-- Policy: DGI agents see all deportees
CREATE POLICY rdep_deportees_dgi_all ON rdep_deportees
    FOR ALL
    USING (current_setting('app.role') IN ('DGI_ADMIN', 'DGI_AGENT', 'SIFR_AGENT'));

-- Policy: PNH officers see monitoring-relevant data
CREATE POLICY rdep_deportees_pnh_monitoring ON rdep_deportees
    FOR SELECT
    USING (monitoring_required = TRUE AND current_setting('app.role') IN ('PNH_OFFICER', 'DCPJ_AGENT', 'PNH_ADMIN'));

COMMIT;
