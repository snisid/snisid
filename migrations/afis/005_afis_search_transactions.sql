BEGIN;

CREATE TABLE afis_search_transactions (
    transaction_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_type    VARCHAR(20) NOT NULL,
    query_subject_id    UUID,
    query_latent_id     UUID,
    hits_count          SMALLINT DEFAULT 0,
    top_score           DECIMAL(5,2),
    top_match_id        UUID,
    search_duration_ms  INTEGER,
    requested_by        UUID NOT NULL,
    requesting_unit     VARCHAR(50),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
