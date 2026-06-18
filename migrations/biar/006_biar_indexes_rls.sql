BEGIN;

CREATE INDEX idx_biar_serial    ON biar_illicit_weapons(serial_number) WHERE serial_number IS NOT NULL;
CREATE INDEX idx_biar_gang      ON biar_illicit_weapons(gang_id) WHERE gang_id IS NOT NULL;
CREATE INDEX idx_biar_dept      ON biar_illicit_weapons(recovery_dept_code);
CREATE INDEX idx_biar_date      ON biar_illicit_weapons(recovery_date DESC);
CREATE INDEX idx_biar_iarms     ON biar_illicit_weapons(iarms_ref) WHERE iarms_ref IS NOT NULL;
CREATE INDEX idx_biar_origin    ON biar_illicit_weapons(origin_country);
CREATE INDEX idx_biar_batch_ref ON biar_batch_seizures(batch_reference);
CREATE INDEX idx_biar_batch_dept ON biar_batch_seizures(dept_code);
CREATE INDEX idx_biar_sync_status ON biar_iarms_sync_log(sync_status);
CREATE INDEX idx_biar_sync_weapon ON biar_iarms_sync_log(weapon_id);
CREATE INDEX idx_biar_route_origin ON biar_trafficking_routes(origin_country);
CREATE INDEX idx_biar_route_active ON biar_trafficking_routes(active) WHERE active = TRUE;

ALTER TABLE biar_illicit_weapons   ENABLE ROW LEVEL SECURITY;
ALTER TABLE biar_batch_seizures    ENABLE ROW LEVEL SECURITY;
ALTER TABLE biar_trafficking_routes ENABLE ROW LEVEL SECURITY;
ALTER TABLE biar_iarms_sync_log    ENABLE ROW LEVEL SECURITY;

CREATE POLICY biar_weapons_select ON biar_illicit_weapons FOR SELECT USING (
    current_setting('app.user_role') IN ('DCPJ_INTEL','DCPJ','PNH','DOUANES','BRI','MJSP')
);

CREATE POLICY biar_weapons_insert ON biar_illicit_weapons FOR INSERT WITH CHECK (
    current_setting('app.user_role') IN ('PNH','DOUANES','DCPJ')
);

CREATE POLICY biar_weapons_update ON biar_illicit_weapons FOR UPDATE USING (
    current_setting('app.user_role') IN ('DCPJ_INTEL','DCPJ')
);

CREATE POLICY biar_batches_select ON biar_batch_seizures FOR SELECT USING (
    current_setting('app.user_role') IN ('DCPJ_INTEL','DCPJ','PNH','DOUANES','BRI')
);

CREATE POLICY biar_batches_insert ON biar_batch_seizures FOR INSERT WITH CHECK (
    current_setting('app.user_role') = 'DCPJ_SUPERVISOR'
);

CREATE POLICY biar_routes_select ON biar_trafficking_routes FOR SELECT USING (
    current_setting('app.user_role') IN ('DCPJ_INTEL','DCPJ','ATF','BRI')
);

CREATE POLICY biar_sync_select ON biar_iarms_sync_log FOR SELECT USING (
    current_setting('app.user_role') = 'SUPERADMIN'
);

CREATE POLICY biar_sync_insert ON biar_iarms_sync_log FOR INSERT WITH CHECK (
    current_setting('app.user_role') = 'SUPERADMIN'
);

COMMIT;
