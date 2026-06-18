BEGIN;

CREATE TABLE expl_legal_stocks (
    stock_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    holder_entity       VARCHAR(200) NOT NULL,
    holder_type         VARCHAR(30),
    explosive_type      expl_type NOT NULL,
    quantity_kg         DECIMAL(12,3) NOT NULL,
    storage_location    TEXT NOT NULL,
    dept_code           CHAR(2),
    license_ref         VARCHAR(50),
    last_audit_date     DATE,
    next_audit_date     DATE,
    is_secured          BOOLEAN DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
