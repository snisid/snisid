BEGIN;
CREATE TYPE trafar_route_type AS ENUM ('MARITIME_DIRECT','MARITIME_VIA_BAHAMAS','AIR_CARGO','AIR_PASSENGER','LAND_BORDER_DOM','LAND_BORDER_OTHER','POSTAL','MIXED');
CREATE TYPE trafar_method AS ENUM ('STRAW_PURCHASE','STOLEN_DIVERTED','CORRUPT_OFFICIAL','FALSE_END_USER','DARK_WEB','DIPLOMATIC_POUCH','CONCEALED_CARGO','DRUGS_FOR_GUNS_SWAP','UNKNOWN');
CREATE TABLE trafar_routes (
    route_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), route_name VARCHAR(150) NOT NULL,
    route_type trafar_route_type NOT NULL, trafficking_method trafar_method NOT NULL,
    origin_country CHAR(3) NOT NULL, origin_city VARCHAR(100), transit_points JSONB,
    entry_point_haiti VARCHAR(100), entry_dept_code CHAR(2), associated_gang_ids UUID[] DEFAULT '{}',
    known_suppliers TEXT[] DEFAULT '{}', activity_level VARCHAR(20) DEFAULT 'ACTIVE',
    estimated_volume_monthly INTEGER, weapon_types TEXT[] DEFAULT '{}',
    intel_confidence SMALLINT CHECK (intel_confidence BETWEEN 1 AND 10),
    first_detected DATE, last_confirmed DATE, linked_case_refs TEXT[] DEFAULT '{}',
    biar_weapon_ids UUID[] DEFAULT '{}', atf_case_refs TEXT[] DEFAULT '{}',
    unodc_ref VARCHAR(50), analyst_notes TEXT, created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE trafar_shipments (
    shipment_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), route_id UUID REFERENCES trafar_routes(route_id),
    shipment_date TIMESTAMPTZ NOT NULL, intercepted BOOLEAN DEFAULT FALSE,
    interception_date TIMESTAMPTZ, interception_location VARCHAR(300),
    interception_unit VARCHAR(50), weapons_count INTEGER, weapons_types TEXT[] DEFAULT '{}',
    estimated_value_usd DECIMAL(12,2), linked_persons UUID[] DEFAULT '{}',
    port_ht_ref UUID, mar_ht_ref UUID, case_reference VARCHAR(100),
    notes TEXT, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE trafar_suppliers (
    supplier_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), supplier_name VARCHAR(200),
    supplier_type VARCHAR(50), country CHAR(3) NOT NULL, city VARCHAR(100),
    snisid_person_id UUID, linked_routes UUID[] DEFAULT '{}',
    atf_subject_ref VARCHAR(50), interpol_notice_ref VARCHAR(50),
    is_active BOOLEAN DEFAULT TRUE, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_trafar_routes_type ON trafar_routes(route_type, activity_level);
CREATE INDEX idx_trafar_routes_origin ON trafar_routes(origin_country);
CREATE INDEX idx_trafar_routes_entry ON trafar_routes(entry_dept_code);
CREATE INDEX idx_trafar_shipments_route ON trafar_shipments(route_id);
CREATE INDEX idx_trafar_shipments_date ON trafar_shipments(shipment_date DESC);
COMMIT;
