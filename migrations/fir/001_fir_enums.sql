BEGIN;

CREATE TYPE fir_offense_class AS ENUM (
    'CONTRAVENTION',
    'DELIT',
    'CRIME',
    'FELONY_FOREIGN'
);

CREATE TYPE fir_case_status AS ENUM (
    'OPEN',
    'PENDING_TRIAL',
    'CONVICTED',
    'ACQUITTED',
    'DISMISSED',
    'APPEAL_PENDING',
    'EXPUNGED'
);

CREATE TYPE fir_sentence_type AS ENUM (
    'PRISON',
    'SUSPENDED',
    'FINE',
    'COMMUNITY_SERVICE',
    'DEATH_PENALTY',
    'ACQUITTAL',
    'PROBATION'
);

CREATE TYPE fir_movement_type AS ENUM (
    'RECORD_CREATED',
    'ARREST_ADDED',
    'CONVICTION_ADDED',
    'ALIAS_ADDED',
    'ALIAS_REMOVED',
    'RECORD_EXPUNGED',
    'RECORD_REACTIVATED',
    'CERTIFICATE_ISSUED'
);

CREATE TYPE fir_certificate_result AS ENUM (
    'CLEAN',
    'HAS_RECORD'
);

COMMIT;
