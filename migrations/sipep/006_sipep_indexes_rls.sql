BEGIN;

CREATE INDEX idx_sipep_inmates_facility ON sipep_inmates(current_facility) WHERE is_currently_detained = TRUE;
CREATE INDEX idx_sipep_inmates_person   ON sipep_inmates(snisid_person_id);
CREATE INDEX idx_sipep_detentions_case  ON sipep_detentions(case_reference);
CREATE INDEX idx_sipep_detentions_basis ON sipep_detentions(detention_basis, legal_status);
CREATE INDEX idx_sipep_movements_inmate ON sipep_movements(inmate_id);
CREATE INDEX idx_sipep_visits_inmate    ON sipep_visits(inmate_id);
CREATE INDEX idx_sipep_health_inmate    ON sipep_health_events(inmate_id);

CREATE MATERIALIZED VIEW sipep_facility_occupancy AS
SELECT
    current_facility,
    COUNT(*) AS current_count,
    current_dept_code
FROM sipep_inmates
WHERE is_currently_detained = TRUE
GROUP BY current_facility, current_dept_code;

COMMIT;
