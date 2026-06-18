BEGIN;

CREATE TABLE trafar_suppliers (
    supplier_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    supplier_name       VARCHAR(200),
    supplier_type       VARCHAR(50),
    country             CHAR(3) NOT NULL,
    city                VARCHAR(100),
    snisid_person_id    UUID,
    linked_routes       UUID[] DEFAULT '{}',
    atf_subject_ref     VARCHAR(50),
    interpol_notice_ref VARCHAR(50),
    is_active           BOOLEAN DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
