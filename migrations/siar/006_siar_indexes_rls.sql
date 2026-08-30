BEGIN;

-- Performance indexes
CREATE INDEX idx_siar_firearms_serial ON siar_firearms(serial_number) WHERE serial_number IS NOT NULL;
CREATE INDEX idx_siar_firearms_status ON siar_firearms(status);
CREATE INDEX idx_siar_firearms_gang ON siar_firearms(gang_id) WHERE gang_id IS NOT NULL;
CREATE INDEX idx_siar_firearms_owner ON siar_firearms(owner_snisid_id) WHERE owner_snisid_id IS NOT NULL;
CREATE INDEX idx_siar_firearms_iarms ON siar_firearms(iarms_ref) WHERE iarms_ref IS NOT NULL;
CREATE INDEX idx_siar_firearms_type ON siar_firearms(weapon_type);
CREATE INDEX idx_siar_firearms_national_id ON siar_firearms(national_siar_id);
CREATE INDEX idx_siar_firearms_dept ON siar_firearms(current_dept_code) WHERE current_dept_code IS NOT NULL;

CREATE INDEX idx_siar_licenses_holder ON siar_licenses(holder_snisid_id, is_active);
CREATE INDEX idx_siar_licenses_number ON siar_licenses(license_number);
CREATE INDEX idx_siar_licenses_expiry ON siar_licenses(expiry_date) WHERE is_active = TRUE;

CREATE INDEX idx_siar_transfers_firearm ON siar_transfers(firearm_id);
CREATE INDEX idx_siar_transfers_date ON siar_transfers(transfer_date);

CREATE INDEX idx_siar_seizures_firearm ON siar_seizures(firearm_id) WHERE firearm_id IS NOT NULL;
CREATE INDEX idx_siar_seizures_date ON siar_seizures(seizure_date);
CREATE INDEX idx_siar_seizures_unit ON siar_seizures(seizing_unit);
CREATE INDEX idx_siar_seizures_person ON siar_seizures(from_person_id) WHERE from_person_id IS NOT NULL;
CREATE INDEX idx_siar_seizures_gang ON siar_seizures(gang_id) WHERE gang_id IS NOT NULL;

CREATE INDEX idx_siar_dealers_status ON siar_dealers(status);
CREATE INDEX idx_siar_dealers_owner ON siar_dealers(owner_snisid_id);
CREATE INDEX idx_siar_dealers_dept ON siar_dealers(dept_code) WHERE dept_code IS NOT NULL;

-- Row-Level Security (RLS)
ALTER TABLE siar_firearms ENABLE ROW LEVEL SECURITY;
ALTER TABLE siar_licenses ENABLE ROW LEVEL SECURITY;
ALTER TABLE siar_transfers ENABLE ROW LEVEL SECURITY;
ALTER TABLE siar_seizures ENABLE ROW LEVEL SECURITY;
ALTER TABLE siar_dealers ENABLE ROW LEVEL SECURITY;

-- PNH officers can read all firearms
CREATE POLICY siar_firearms_pnh_select ON siar_firearms
    FOR SELECT USING (
        current_setting('snisid.user_role') IN ('PNH_OFFICER','PNH_ADMIN','DCPJ','MJSP_ADMIN')
    );

-- MJSP_ADMIN can insert/update
CREATE POLICY siar_firearms_mjsp_insert ON siar_firearms
    FOR INSERT WITH CHECK (
        current_setting('snisid.user_role') IN ('MJSP_ADMIN','PNH_ADMIN','DCPJ')
    );

CREATE POLICY siar_firearms_mjsp_update ON siar_firearms
    FOR UPDATE USING (
        current_setting('snisid.user_role') IN ('MJSP_ADMIN','PNH_ADMIN','DCPJ')
    );

-- License RLS
CREATE POLICY siar_licenses_select ON siar_licenses
    FOR SELECT USING (
        current_setting('snisid.user_role') IN ('MJSP_ADMIN','PNH_OFFICER','PNH_ADMIN','DCPJ')
    );

CREATE POLICY siar_licenses_insert ON siar_licenses
    FOR INSERT WITH CHECK (
        current_setting('snisid.user_role') = 'MJSP_ADMIN'
    );

-- Seizure RLS
CREATE POLICY siar_seizures_select ON siar_seizures
    FOR SELECT USING (
        current_setting('snisid.user_role') IN ('PNH_OFFICER','PNH_ADMIN','DCPJ','MJSP_ADMIN')
    );

CREATE POLICY siar_seizures_insert ON siar_seizures
    FOR INSERT WITH CHECK (
        current_setting('snisid.user_role') IN ('PNH_OFFICER','PNH_ADMIN')
    );

-- Dealer RLS
CREATE POLICY siar_dealers_select ON siar_dealers
    FOR SELECT USING (
        current_setting('snisid.user_role') IN ('MJSP_ADMIN','PNH_ADMIN','DCPJ')
    );

CREATE POLICY siar_dealers_insert ON siar_dealers
    FOR INSERT WITH CHECK (
        current_setting('snisid.user_role') = 'MJSP_ADMIN'
    );

COMMIT;
