-- Migration 003 : Références biométriques

CREATE TYPE snisid_biometric.biometric_modality AS ENUM (
    'FINGERPRINT_10',
    'FINGERPRINT_4_4_2',
    'FACE_2D',
    'FACE_3D',
    'IRIS_BOTH',
    'IRIS_LEFT',
    'IRIS_RIGHT',
    'SIGNATURE'
);

CREATE TABLE snisid_biometric.biometric_references (
    bio_ref_id          UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    niu                 CHAR(10) NOT NULL,
    modality            snisid_biometric.biometric_modality NOT NULL,
    abis_gallery_id     VARCHAR(100),
    vault_template_ref  VARCHAR(200),
    quality_score       DECIMAL(5,4) CHECK (quality_score BETWEEN 0 AND 1),
    capture_device      VARCHAR(100),
    capture_date        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    verification_count  INTEGER NOT NULL DEFAULT 0,
    last_verified_at    TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (niu, modality, is_active)
);

CREATE INDEX idx_bio_niu ON snisid_biometric.biometric_references (niu, modality);
CREATE INDEX idx_bio_gallery ON snisid_biometric.biometric_references (abis_gallery_id) WHERE abis_gallery_id IS NOT NULL;

CREATE UNIQUE INDEX idx_bio_active_unique ON snisid_biometric.biometric_references (niu, modality) WHERE is_active = TRUE;

CREATE TYPE snisid_identity.document_type AS ENUM (
    'ACTE_NAISSANCE', 'PASSEPORT', 'CIN', 'PERMIS_CONDUIRE',
    'TITRE_SEJOUR', 'JUGEMENT_SUPPLETTIF', 'AUTRE'
);

CREATE TABLE snisid_identity.identity_documents (
    doc_id              UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    niu                 CHAR(10) NOT NULL,
    document_type       snisid_identity.document_type NOT NULL,
    reference_number    VARCHAR(100),
    verification_status VARCHAR(20) DEFAULT 'PENDING' CHECK (verification_status IN ('PENDING','VERIFIED','REJECTED','EXPIRED')),
    forgery_score       DECIMAL(5,4),
    scan_reference      VARCHAR(500),
    issued_date         DATE,
    expiry_date         DATE,
    issuing_authority   VARCHAR(100),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_docs_niu ON snisid_identity.identity_documents (niu, document_type);
CREATE INDEX idx_docs_ref ON snisid_identity.identity_documents (reference_number) WHERE reference_number IS NOT NULL;
