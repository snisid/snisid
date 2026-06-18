BEGIN;

CREATE INDEX idx_trafar_routes_type     ON trafar_routes(route_type, activity_level);
CREATE INDEX idx_trafar_routes_origin   ON trafar_routes(origin_country);
CREATE INDEX idx_trafar_routes_entry    ON trafar_routes(entry_dept_code);
CREATE INDEX idx_trafar_shipments_route ON trafar_shipments(route_id);
CREATE INDEX idx_trafar_shipments_date  ON trafar_shipments(shipment_date DESC);

ALTER TABLE trafar_routes ENABLE ROW LEVEL SECURITY;
ALTER TABLE trafar_shipments ENABLE ROW LEVEL SECURITY;
ALTER TABLE trafar_suppliers ENABLE ROW LEVEL SECURITY;

CREATE POLICY trafar_routes_select ON trafar_routes
    FOR SELECT USING (
        current_setting('snisid.user_role') IN ('DCPJ','DCPJ_INTEL','PNH_ADMIN','ATF')
    );

CREATE POLICY trafar_routes_insert ON trafar_routes
    FOR INSERT WITH CHECK (
        current_setting('snisid.user_role') = 'DCPJ_INTEL'
    );

CREATE POLICY trafar_shipments_select ON trafar_shipments
    FOR SELECT USING (
        current_setting('snisid.user_role') IN ('DCPJ','DCPJ_INTEL','DOUANES','PNH_ADMIN')
    );

CREATE POLICY trafar_shipments_insert ON trafar_shipments
    FOR INSERT WITH CHECK (
        current_setting('snisid.user_role') IN ('DCPJ','DOUANES')
    );

COMMIT;
