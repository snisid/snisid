BEGIN;

CREATE TABLE reso_persons (
    snisid_id           UUID PRIMARY KEY,
    name                VARCHAR(200) NOT NULL,
    aliases             TEXT[] DEFAULT '{}',
    nationality         VARCHAR(3),
    dob                 DATE,
    risk_score          DECIMAL(3,2) DEFAULT 0.0,
    is_gang_member      BOOLEAN DEFAULT FALSE,
    is_sanctioned       BOOLEAN DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE reso_gangs (
    gang_id             UUID PRIMARY KEY,
    name                VARCHAR(200) NOT NULL,
    primary_activity    VARCHAR(100),
    territory_dept      CHAR(2),
    activity_level      VARCHAR(20) DEFAULT 'ACTIVE',
    member_count        INTEGER DEFAULT 0,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE reso_person_gang (
    person_id           UUID NOT NULL REFERENCES reso_persons(snisid_id),
    gang_id             UUID NOT NULL REFERENCES reso_gangs(gang_id),
    role                VARCHAR(50) NOT NULL,
    since               DATE,
    confidence          DECIMAL(3,2) DEFAULT 1.0,
    PRIMARY KEY (person_id, gang_id)
);

CREATE TABLE reso_gang_relations (
    gang_id_1           UUID NOT NULL REFERENCES reso_gangs(gang_id),
    gang_id_2           UUID NOT NULL REFERENCES reso_gangs(gang_id),
    relation_type       VARCHAR(30) NOT NULL,
    since               DATE,
    confidence          DECIMAL(3,2) DEFAULT 1.0,
    PRIMARY KEY (gang_id_1, gang_id_2)
);

CREATE TABLE reso_person_associations (
    person_id_1         UUID NOT NULL REFERENCES reso_persons(snisid_id),
    person_id_2         UUID NOT NULL REFERENCES reso_persons(snisid_id),
    association_type    VARCHAR(30) NOT NULL,
    confidence          DECIMAL(3,2) DEFAULT 1.0,
    source              VARCHAR(50),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (person_id_1, person_id_2)
);

CREATE TABLE reso_criminal_networks (
    network_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    network_name        VARCHAR(200),
    member_ids          UUID[] DEFAULT '{}',
    community_size      INTEGER DEFAULT 0,
    modularity_score    DECIMAL(5,4),
    detected_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    analysis_version    VARCHAR(20) DEFAULT '1.0'
);

CREATE INDEX idx_reso_pers_gang ON reso_person_gang(gang_id);
CREATE INDEX idx_reso_assoc_1   ON reso_person_associations(person_id_1);
CREATE INDEX idx_reso_assoc_2   ON reso_person_associations(person_id_2);
CREATE INDEX idx_reso_gang_rel  ON reso_gang_relations(gang_id_1);

COMMIT;
