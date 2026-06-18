BEGIN;

CREATE INDEX idx_afis_subjects_snisid    ON afis_subjects(snisid_person_id) WHERE snisid_person_id IS NOT NULL;
CREATE INDEX idx_afis_prints_subject     ON afis_fingerprints(subject_id, finger_position);
CREATE INDEX idx_afis_prints_quality     ON afis_fingerprints(quality_accepted) WHERE quality_accepted = TRUE;
CREATE INDEX idx_afis_latent_case        ON afis_latent_prints(case_reference);
CREATE INDEX idx_afis_latent_identified  ON afis_latent_prints(is_identified);

COMMIT;
