BEGIN;

CREATE TABLE chef_cross_gang_links (
    link_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    member_a_id         UUID NOT NULL REFERENCES chef_criminal_members(member_id),
    member_b_id         UUID NOT NULL REFERENCES chef_criminal_members(member_id),
    link_type           VARCHAR(50),
    confidence          SMALLINT,
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT no_self_link CHECK (member_a_id <> member_b_id)
);

COMMIT;
