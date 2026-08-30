-- ============================================================
-- SNISID — NATIONAL IDENTITY REGISTRY DATABASE SCHEMA
-- Document ID: SNISID-IDN-DB-001
-- Version: 1.0.0
-- Date: Mai 2026
-- Classification: SOUVERAIN / INFRASTRUCTURE CRITIQUE
-- ============================================================

-- Extension requirements
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";  -- For fuzzy text search

-- ============================================================
-- ENUMERATIONS
-- ============================================================

CREATE TYPE identity_status AS ENUM (
    'PRE_REGISTERED',
    'PENDING_BIOMETRICS',
    'ACTIVE',
    'SUSPENDED',
    'REVOKED',
    'DECEASED',
    'ARCHIVED'
);

CREATE TYPE sexe_type AS ENUM ('M', 'F', 'INDETERMINE');

CREATE TYPE document_type AS ENUM (
    'ACTE_NAISSANCE',
    'JUGEMENT_SUPPLETIF',
    'PASSEPORT',
    'ANCIEN_CIN',
    'CERTIFICAT_BAPTEME',
    'DOSSIER_HOPITAL',
    'CARTE_ELECTORALE',
    'AFFIDAVIT',
    'ACTE_MARIAGE',
    'AUTRE'
);

CREATE TYPE document_verification_status AS ENUM (
    'PENDING',
    'VERIFIED_REGISTRY',
    'VERIFIED_COURT',
    'VERIFIED_MANUAL',
    'REJECTED',
    'SUSPECTED_FORGERY'
);

CREATE TYPE biometric_modality AS ENUM (
    'FINGERPRINT_10',
    'FINGERPRINT_PARTIAL',
    'IRIS_DUAL',
    'IRIS_SINGLE',
    'FACE_3D',
    'FACE_2D'
);

CREATE TYPE civil_act_type AS ENUM (
    'NAISSANCE_SIMPLE',
    'NAISSANCE_RECONNAISSANCE',
    'NAISSANCE_TARDIVE',
    'NAISSANCE_DECRET',
    'NAISSANCE_JUGEMENT',
    'MARIAGE_CIVIL',
    'MARIAGE_CONCORDATAIRE',
    'DIVORCE_CONTENTIEUX',
    'DIVORCE_MUTUEL',
    'DECES',
    'ADOPTION_SIMPLE',
    'ADOPTION_PLENIERE',
    'CORRECTION_ADMINISTRATIVE',
    'ANNOTATION_MARGINALE'
);

-- ============================================================
-- SCHEMA: snisid_identity
-- ============================================================
CREATE SCHEMA IF NOT EXISTS snisid_identity;
CREATE SCHEMA IF NOT EXISTS snisid_civil;
CREATE SCHEMA IF NOT EXISTS snisid_audit;
CREATE SCHEMA IF NOT EXISTS snisid_biometric;

SET search_path TO snisid_identity, snisid_civil, snisid_audit, snisid_biometric, public;

-- ============================================================
-- TABLE: citizens (Master Identity Record)
-- ============================================================
CREATE TABLE snisid_identity.citizens (
    -- Primary Key
    niu                         VARCHAR(10)         PRIMARY KEY,        -- Numéro d'Identification Unique (10 chiffres, cryptographiquement aléatoire)
    niu_version                 INTEGER             NOT NULL DEFAULT 1, -- Version du NUI (incrémenté à chaque renouvellement)

    -- État de l'identité
    statut_identite             identity_status     NOT NULL DEFAULT 'PRE_REGISTERED',
    statut_updated_at           TIMESTAMP WITH TIME ZONE,
    statut_updated_by           VARCHAR(50),        -- Agent ou système ayant changé le statut
    statut_reason               TEXT,               -- Raison du changement de statut

    -- Données démographiques nominatives
    nom                         VARCHAR(100)        NOT NULL,
    prenom                      VARCHAR(150)        NOT NULL,
    nom_jeune_fille             VARCHAR(100),       -- Nom de jeune fille (si applicable)
    autres_prenoms              VARCHAR(200),       -- Prénoms supplémentaires
    date_naissance              DATE                NOT NULL,
    date_naissance_incertaine   BOOLEAN             NOT NULL DEFAULT FALSE, -- Vrai si date approximative
    lieu_naissance_commune      VARCHAR(100)        NOT NULL,
    lieu_naissance_departement  VARCHAR(50)         NOT NULL,
    lieu_naissance_pays         VARCHAR(3)          NOT NULL DEFAULT 'HT',  -- ISO 3166-1 alpha-3
    sexe                        sexe_type           NOT NULL,
    nationalite_primaire        VARCHAR(3)          NOT NULL DEFAULT 'HT',  -- ISO 3166-1 alpha-3
    nationalite_secondaire      VARCHAR(3),

    -- Filiation
    pere_niu                    VARCHAR(10)         REFERENCES snisid_identity.citizens(niu) ON DELETE RESTRICT,
    pere_nom_complet            VARCHAR(200),       -- Si père non enregistré dans le système
    mere_niu                    VARCHAR(10)         REFERENCES snisid_identity.citizens(niu) ON DELETE RESTRICT,
    mere_nom_complet            VARCHAR(200),       -- Si mère non enregistrée dans le système

    -- Statut civil
    statut_matrimonial          VARCHAR(20)         DEFAULT 'CELIBATAIRE', -- CELIBATAIRE, MARIE, DIVORCE, VEUF
    conjoint_niu                VARCHAR(10)         REFERENCES snisid_identity.citizens(niu) ON DELETE SET NULL,
    date_mariage                DATE,
    date_divorce                DATE,
    date_deces                  DATE,               -- Renseigné par ANH lors du décès

    -- Adresse actuelle
    adresse_ligne1              VARCHAR(200),
    adresse_ligne2              VARCHAR(200),
    commune_residence           VARCHAR(100),
    departement_residence       VARCHAR(50),
    pays_residence              VARCHAR(3)          DEFAULT 'HT',
    adresse_gps_lat             DECIMAL(10, 7),
    adresse_gps_lon             DECIMAL(10, 7),

    -- Enregistrement
    date_enregistrement         TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    agent_enregistreur_id       VARCHAR(50)         NOT NULL,
    agent_enregistreur_nom      VARCHAR(150),
    centre_enregistrement_id    VARCHAR(20)         NOT NULL,
    centre_enregistrement_nom   VARCHAR(150),
    canal_enregistrement        VARCHAR(20)         NOT NULL, -- FIXED_SITE, MOBILE_KIT, OFFLINE_EMERGENCY
    session_enrollment_id       UUID                NOT NULL,

    -- Documents sources ayant servi à l'enregistrement
    document_source_primaire    document_type       NOT NULL,
    document_source_reference   VARCHAR(100),

    -- PKI / Certificats
    cert_auth_thumbprint        VARCHAR(64),        -- SHA-256 du certificat d'authentification
    cert_sign_thumbprint        VARCHAR(64),        -- SHA-256 du certificat de signature
    cert_issued_at              TIMESTAMP WITH TIME ZONE,
    cert_expires_at             TIMESTAMP WITH TIME ZONE,

    -- Fraude et scoring
    fraud_score_initial         DECIMAL(5,4),       -- Score de fraude à l'enregistrement (0.0000-1.0000)
    fraud_flags                 JSONB,              -- Flags de fraude actifs
    risk_level                  VARCHAR(10),        -- LOW, MEDIUM, HIGH

    -- Métadonnées système
    version                     BIGINT              NOT NULL DEFAULT 1, -- Optimistic locking
    created_at                  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at                  TIMESTAMP WITH TIME ZONE,               -- Soft delete (RGPD)
    data_hash                   VARCHAR(64),        -- SHA-256 de l'enregistrement pour tamper detection

    -- Conformité
    consentement_enregistre     BOOLEAN             NOT NULL DEFAULT FALSE,
    consentement_date           TIMESTAMP WITH TIME ZONE,
    gdpr_erasure_requested      BOOLEAN             NOT NULL DEFAULT FALSE,
    gdpr_erasure_date           TIMESTAMP WITH TIME ZONE

) PARTITION BY LIST (departement_residence);

-- Partitions par département d'Haïti
CREATE TABLE snisid_identity.citizens_ouest       PARTITION OF snisid_identity.citizens FOR VALUES IN ('Ouest');
CREATE TABLE snisid_identity.citizens_nord        PARTITION OF snisid_identity.citizens FOR VALUES IN ('Nord');
CREATE TABLE snisid_identity.citizens_nord_est    PARTITION OF snisid_identity.citizens FOR VALUES IN ('Nord-Est');
CREATE TABLE snisid_identity.citizens_nord_ouest  PARTITION OF snisid_identity.citizens FOR VALUES IN ('Nord-Ouest');
CREATE TABLE snisid_identity.citizens_artibonite  PARTITION OF snisid_identity.citizens FOR VALUES IN ('Artibonite');
CREATE TABLE snisid_identity.citizens_centre      PARTITION OF snisid_identity.citizens FOR VALUES IN ('Centre');
CREATE TABLE snisid_identity.citizens_sud         PARTITION OF snisid_identity.citizens FOR VALUES IN ('Sud');
CREATE TABLE snisid_identity.citizens_sud_est     PARTITION OF snisid_identity.citizens FOR VALUES IN ('Sud-Est');
CREATE TABLE snisid_identity.citizens_grande_anse PARTITION OF snisid_identity.citizens FOR VALUES IN ('Grande-Anse');
CREATE TABLE snisid_identity.citizens_nippes      PARTITION OF snisid_identity.citizens FOR VALUES IN ('Nippes');
CREATE TABLE snisid_identity.citizens_diaspora    PARTITION OF snisid_identity.citizens FOR VALUES IN ('Diaspora');
CREATE TABLE snisid_identity.citizens_unknown     PARTITION OF snisid_identity.citizens DEFAULT;

-- Indexes
CREATE INDEX idx_citizens_nom ON snisid_identity.citizens USING gin(to_tsvector('french', nom || ' ' || prenom));
CREATE INDEX idx_citizens_date_naissance ON snisid_identity.citizens(date_naissance);
CREATE INDEX idx_citizens_statut ON snisid_identity.citizens(statut_identite);
CREATE INDEX idx_citizens_commune ON snisid_identity.citizens(lieu_naissance_commune);
CREATE INDEX idx_citizens_pere ON snisid_identity.citizens(pere_niu) WHERE pere_niu IS NOT NULL;
CREATE INDEX idx_citizens_mere ON snisid_identity.citizens(mere_niu) WHERE mere_niu IS NOT NULL;
CREATE INDEX idx_citizens_conjoint ON snisid_identity.citizens(conjoint_niu) WHERE conjoint_niu IS NOT NULL;
CREATE INDEX idx_citizens_nom_trgm ON snisid_identity.citizens USING gin(nom gin_trgm_ops, prenom gin_trgm_ops);
CREATE INDEX idx_citizens_session ON snisid_identity.citizens(session_enrollment_id);

-- ============================================================
-- TABLE: identity_events (Event Store — IMMUTABLE)
-- ============================================================
CREATE TABLE snisid_identity.identity_events (
    event_id            UUID                PRIMARY KEY DEFAULT gen_random_uuid(),
    niu                 VARCHAR(10)         NOT NULL,
    event_type          VARCHAR(80)         NOT NULL,
    -- Ex: CitizenRegistered, IdentityActivated, IdentityUpdated,
    --     AddressUpdated, StatusChanged, CertificateIssued,
    --     BiometricEnrolled, IdentitySuspended, IdentityRevoked,
    --     IdentityDeceased, IdentityArchived, FraudFlagged, ConsentUpdated

    event_version       INTEGER             NOT NULL,   -- Séquence ordonnée par NIU
    event_timestamp     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    -- Payload de l'événement
    payload             JSONB               NOT NULL,   -- Données de l'événement (schéma par event_type)
    payload_schema_ver  VARCHAR(10)         NOT NULL DEFAULT 'v1', -- Version du schéma payload

    -- Metadata (immuable, piste d'audit)
    agent_id            VARCHAR(50)         NOT NULL,
    agent_nom           VARCHAR(150),
    session_id          UUID,
    trace_id            VARCHAR(64),        -- OpenTelemetry trace ID
    span_id             VARCHAR(32),
    source_system       VARCHAR(50)         NOT NULL,   -- ex: identity-service, workflow-engine
    source_channel      VARCHAR(20),        -- FIXED_SITE, MOBILE_KIT, API
    client_ip           INET,
    device_id           VARCHAR(100),

    -- Intégrité
    previous_hash       VARCHAR(64),        -- Hash de l'événement précédent (chaîne)
    event_hash          VARCHAR(64)         NOT NULL,   -- SHA-256 de ce event_id+payload+timestamp
    signature           TEXT,              -- Signature PKI de l'agent

    -- Kafka
    kafka_topic         VARCHAR(100),
    kafka_partition     INTEGER,
    kafka_offset        BIGINT,
    kafka_published_at  TIMESTAMP WITH TIME ZONE,

    UNIQUE (niu, event_version)
);

-- L'event store est APPEND-ONLY — aucun UPDATE/DELETE autorisé
CREATE RULE no_update_events AS ON UPDATE TO snisid_identity.identity_events DO INSTEAD NOTHING;
CREATE RULE no_delete_events AS ON DELETE TO snisid_identity.identity_events DO INSTEAD NOTHING;

CREATE INDEX idx_events_niu ON snisid_identity.identity_events(niu);
CREATE INDEX idx_events_type ON snisid_identity.identity_events(event_type);
CREATE INDEX idx_events_timestamp ON snisid_identity.identity_events(event_timestamp);
CREATE INDEX idx_events_agent ON snisid_identity.identity_events(agent_id);
CREATE INDEX idx_events_trace ON snisid_identity.identity_events(trace_id);

-- ============================================================
-- TABLE: identity_snapshots (CQRS Optimization)
-- ============================================================
CREATE TABLE snisid_identity.identity_snapshots (
    niu                 VARCHAR(10)         PRIMARY KEY REFERENCES snisid_identity.citizens(niu),
    latest_version      INTEGER             NOT NULL,
    state               JSONB               NOT NULL,   -- Snapshot complet de l'état courant
    state_hash          VARCHAR(64)         NOT NULL,
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- ============================================================
-- TABLE: biometric_references
-- ============================================================
CREATE TABLE snisid_biometric.biometric_references (
    bio_ref_id          UUID                PRIMARY KEY DEFAULT gen_random_uuid(),
    niu                 VARCHAR(10)         NOT NULL REFERENCES snisid_identity.citizens(niu),
    modality            biometric_modality  NOT NULL,
    abis_gallery_id     VARCHAR(100)        NOT NULL UNIQUE, -- ID dans le système ABIS
    template_hash       VARCHAR(64)         NOT NULL,        -- SHA-256 du template chiffré (pas le template lui-même)
    quality_score       DECIMAL(5,2),                        -- Score NFIQ2 ou équivalent
    enrolled_at         TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    enrolled_by         VARCHAR(50)         NOT NULL,
    capture_device_id   VARCHAR(100),
    is_active           BOOLEAN             NOT NULL DEFAULT TRUE,
    revoked_at          TIMESTAMP WITH TIME ZONE,
    revoked_reason      TEXT,

    -- Flags
    is_partial          BOOLEAN             NOT NULL DEFAULT FALSE, -- Capture partielle (doigts manquants)
    accommodation_type  VARCHAR(50),                               -- AMPUTATION, BLINDNESS, ELDERLY, INFANT
    accommodation_doc   VARCHAR(100),

    -- Audit
    last_verified_at    TIMESTAMP WITH TIME ZONE,
    verification_count  INTEGER             NOT NULL DEFAULT 0,

    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bio_niu ON snisid_biometric.biometric_references(niu);
CREATE INDEX idx_bio_modality ON snisid_biometric.biometric_references(modality);
CREATE INDEX idx_bio_active ON snisid_biometric.biometric_references(niu, is_active);

-- ============================================================
-- TABLE: identity_documents (Documents source)
-- ============================================================
CREATE TABLE snisid_identity.identity_documents (
    doc_id              UUID                PRIMARY KEY DEFAULT gen_random_uuid(),
    niu                 VARCHAR(10)         NOT NULL REFERENCES snisid_identity.citizens(niu),
    document_type       document_type       NOT NULL,
    reference_number    VARCHAR(100),       -- Numéro de l'acte, passeport, etc.
    issued_by           VARCHAR(200),       -- Organisme émetteur
    issued_at_commune   VARCHAR(100),
    issued_at_date      DATE,
    expiry_date         DATE,
    verification_status document_verification_status NOT NULL DEFAULT 'PENDING',
    verified_by         VARCHAR(50),
    verified_at         TIMESTAMP WITH TIME ZONE,
    verification_method VARCHAR(50),        -- REGISTRY_API, COURT_API, MANUAL, LEGACY_DB
    forgery_score       DECIMAL(5,4),       -- Score IA de falsification (0 = authentique)
    scan_reference      VARCHAR(200),       -- Référence au fichier scanné dans MinIO
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_docs_niu ON snisid_identity.identity_documents(niu);
CREATE INDEX idx_docs_type ON snisid_identity.identity_documents(document_type);
CREATE INDEX idx_docs_ref ON snisid_identity.identity_documents(reference_number);

-- ============================================================
-- TABLE: civil_acts (Actes d'état civil)
-- ============================================================
CREATE TABLE snisid_civil.civil_acts (
    act_id              UUID                PRIMARY KEY DEFAULT gen_random_uuid(),
    act_type            civil_act_type      NOT NULL,
    numero_acte         VARCHAR(50)         NOT NULL UNIQUE, -- Numéro officiel de l'acte
    date_acte           DATE                NOT NULL,
    lieu_acte_commune   VARCHAR(100)        NOT NULL,
    lieu_acte_departement VARCHAR(50)       NOT NULL,

    -- Parties concernées
    niu_principal       VARCHAR(10)         REFERENCES snisid_identity.citizens(niu),
    niu_secondaire      VARCHAR(10)         REFERENCES snisid_identity.citizens(niu),
    niu_pere            VARCHAR(10)         REFERENCES snisid_identity.citizens(niu),
    niu_mere            VARCHAR(10)         REFERENCES snisid_identity.citizens(niu),

    -- Officier d'état civil
    officier_etat_civil_id   VARCHAR(50)    NOT NULL,
    officier_etat_civil_nom  VARCHAR(150)   NOT NULL,
    signature_oec            TEXT,          -- Signature XAdES de l'OEC
    centre_etat_civil_id     VARCHAR(20)    NOT NULL,

    -- Témoins
    temoin1_niu         VARCHAR(10),
    temoin1_nom         VARCHAR(200),
    temoin2_niu         VARCHAR(10),
    temoin2_nom         VARCHAR(200),

    -- Document généré
    pdf_reference       VARCHAR(200),       -- Chemin MinIO du PDF/A-3 généré
    pdf_hash            VARCHAR(64),        -- SHA-256 du PDF
    qr_token            VARCHAR(200),       -- JWT QR code token
    xades_signature     TEXT,              -- Signature XAdES-LTA

    -- Référence judiciaire (si applicable)
    reference_jugement  VARCHAR(100),
    tribunal            VARCHAR(150),
    numero_dossier      VARCHAR(50),

    -- Workflow
    workflow_id         VARCHAR(100),       -- ID Temporal workflow
    statut_workflow     VARCHAR(30),        -- PENDING, IN_PROGRESS, COMPLETED, REJECTED
    validated_at        TIMESTAMP WITH TIME ZONE,

    -- Kafka
    kafka_event_id      VARCHAR(100),

    -- Audit
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by          VARCHAR(50)         NOT NULL
);

CREATE INDEX idx_civil_acts_niu ON snisid_civil.civil_acts(niu_principal);
CREATE INDEX idx_civil_acts_type ON snisid_civil.civil_acts(act_type);
CREATE INDEX idx_civil_acts_date ON snisid_civil.civil_acts(date_acte);
CREATE INDEX idx_civil_acts_commune ON snisid_civil.civil_acts(lieu_acte_commune);
CREATE INDEX idx_civil_acts_oec ON snisid_civil.civil_acts(officier_etat_civil_id);
CREATE INDEX idx_civil_acts_numero ON snisid_civil.civil_acts(numero_acte);

-- ============================================================
-- TABLE: officiers_etat_civil
-- ============================================================
CREATE TABLE snisid_civil.officiers_etat_civil (
    oec_id              VARCHAR(50)         PRIMARY KEY,
    niu                 VARCHAR(10)         REFERENCES snisid_identity.citizens(niu),
    nom_complet         VARCHAR(200)        NOT NULL,
    commune_assignee    VARCHAR(100)        NOT NULL,
    departement         VARCHAR(50)         NOT NULL,
    cert_thumbprint     VARCHAR(64),        -- Certificat PKI pour signature
    date_nomination     DATE                NOT NULL,
    date_expiration     DATE,
    statut              VARCHAR(20)         NOT NULL DEFAULT 'ACTIF',
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- ============================================================
-- TABLE: audit_trail (Piste d'audit complète)
-- ============================================================
CREATE TABLE snisid_audit.audit_trail (
    audit_id            UUID                PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type          VARCHAR(100)        NOT NULL,
    entity_type         VARCHAR(50)         NOT NULL,   -- CITIZEN, CIVIL_ACT, BIOMETRIC, etc.
    entity_id           VARCHAR(100)        NOT NULL,
    niu                 VARCHAR(10),

    -- Qui
    agent_id            VARCHAR(50)         NOT NULL,
    agency_id           VARCHAR(30)         NOT NULL,
    agency_name         VARCHAR(100),
    session_id          UUID,

    -- Quoi
    action              VARCHAR(50)         NOT NULL,   -- CREATE, READ, UPDATE, DELETE, VERIFY, EXPORT
    before_state        JSONB,
    after_state         JSONB,
    diff                JSONB,

    -- Contexte technique
    trace_id            VARCHAR(64),
    client_ip           INET,
    user_agent          VARCHAR(200),
    device_id           VARCHAR(100),
    geo_location        JSONB,              -- {commune, departement, lat, lon}

    -- Résultat
    status              VARCHAR(20)         NOT NULL,   -- SUCCESS, FAILURE, DENIED
    error_code          VARCHAR(50),
    error_message       TEXT,

    -- Intégrité
    audit_hash          VARCHAR(64)         NOT NULL,   -- SHA-256 de l'entrée complète
    chain_hash          VARCHAR(64),        -- Hash précédent (Merkle chain)

    timestamp           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (timestamp);

-- Partitions mensuelles pour la rétention (7 ans minimum)
CREATE TABLE snisid_audit.audit_trail_2026_01 PARTITION OF snisid_audit.audit_trail
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');
CREATE TABLE snisid_audit.audit_trail_2026_02 PARTITION OF snisid_audit.audit_trail
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');
-- ... (généré automatiquement par pg_partman)

CREATE INDEX idx_audit_entity ON snisid_audit.audit_trail(entity_type, entity_id);
CREATE INDEX idx_audit_niu ON snisid_audit.audit_trail(niu);
CREATE INDEX idx_audit_agent ON snisid_audit.audit_trail(agent_id);
CREATE INDEX idx_audit_timestamp ON snisid_audit.audit_trail(timestamp);
CREATE INDEX idx_audit_trace ON snisid_audit.audit_trail(trace_id);
CREATE INDEX idx_audit_agency ON snisid_audit.audit_trail(agency_id);

-- ============================================================
-- TABLE: verification_log (Log des vérifications)
-- ============================================================
CREATE TABLE snisid_audit.verification_log (
    request_id          UUID                PRIMARY KEY DEFAULT gen_random_uuid(),
    niu                 VARCHAR(10)         NOT NULL,
    requesting_agency   VARCHAR(30)         NOT NULL,
    requesting_agent    VARCHAR(50),
    verification_type   VARCHAR(30)         NOT NULL, -- IDENTITY, BIOMETRIC, QR, USSD, BATCH
    verification_level  INTEGER             NOT NULL, -- 1=basic, 2=photo, 3=biometric, 4=pki
    result              VARCHAR(20)         NOT NULL, -- MATCH, NO_MATCH, ERROR, SUSPENDED
    match_score         DECIMAL(5,4),
    disclosed_fields    VARCHAR[]           NOT NULL, -- Champs divulgués (minimisation)
    response_time_ms    INTEGER,
    timestamp           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    client_ip           INET,
    audit_hash          VARCHAR(64)         NOT NULL
);

CREATE INDEX idx_verif_niu ON snisid_audit.verification_log(niu);
CREATE INDEX idx_verif_agency ON snisid_audit.verification_log(requesting_agency);
CREATE INDEX idx_verif_timestamp ON snisid_audit.verification_log(timestamp);

-- ============================================================
-- ROW LEVEL SECURITY (RLS)
-- ============================================================

-- Activer RLS sur la table citizens
ALTER TABLE snisid_identity.citizens ENABLE ROW LEVEL SECURITY;

-- Politique: agents ONI peuvent voir leur département
CREATE POLICY oni_department_access ON snisid_identity.citizens
    FOR SELECT
    TO oni_agent_role
    USING (departement_residence = current_setting('app.agent_department', TRUE));

-- Politique: DGI peut voir statut fiscal uniquement
CREATE POLICY dgi_minimal_access ON snisid_identity.citizens
    FOR SELECT
    TO dgi_agent_role
    USING (statut_identite = 'ACTIVE');

-- Politique: PNH peut voir tout pour enquête
CREATE POLICY pnh_full_read ON snisid_identity.citizens
    FOR SELECT
    TO pnh_investigation_role
    USING (TRUE);

-- Politique: admin SNISID accès complet
CREATE POLICY snisid_admin_full ON snisid_identity.citizens
    FOR ALL
    TO snisid_admin_role
    USING (TRUE)
    WITH CHECK (TRUE);

-- ============================================================
-- TRIGGERS
-- ============================================================

-- Trigger: mise à jour automatique de updated_at
CREATE OR REPLACE FUNCTION snisid_identity.update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    NEW.version = OLD.version + 1;
    -- Calcul hash pour tamper detection
    NEW.data_hash = encode(sha256(
        (NEW.niu || NEW.nom || NEW.prenom || NEW.date_naissance::text || NEW.statut_identite::text || NEW.version::text)::bytea
    ), 'hex');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

CREATE TRIGGER citizens_update_timestamp
    BEFORE UPDATE ON snisid_identity.citizens
    FOR EACH ROW EXECUTE FUNCTION snisid_identity.update_timestamp();

-- Trigger: audit automatique sur citizens
CREATE OR REPLACE FUNCTION snisid_audit.audit_citizens_change()
RETURNS TRIGGER AS $$
DECLARE
    v_audit_hash VARCHAR(64);
BEGIN
    v_audit_hash := encode(sha256(
        (TG_OP || COALESCE(NEW.niu, OLD.niu) || NOW()::text || current_user)::bytea
    ), 'hex');

    INSERT INTO snisid_audit.audit_trail (
        event_type, entity_type, entity_id, niu,
        agent_id, agency_id,
        action, before_state, after_state,
        status, audit_hash, chain_hash, timestamp
    ) VALUES (
        'CITIZEN_' || TG_OP,
        'CITIZEN',
        COALESCE(NEW.niu, OLD.niu),
        COALESCE(NEW.niu, OLD.niu),
        COALESCE(current_setting('app.agent_id', TRUE), 'SYSTEM'),
        COALESCE(current_setting('app.agency_id', TRUE), 'SYSTEM'),
        TG_OP,
        CASE WHEN TG_OP != 'INSERT' THEN row_to_json(OLD) ELSE NULL END,
        CASE WHEN TG_OP != 'DELETE' THEN row_to_json(NEW) ELSE NULL END,
        'SUCCESS',
        v_audit_hash,
        NULL, -- Chaîne Merkle gérée séparément
        NOW()
    );

    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

CREATE TRIGGER citizens_audit_trigger
    AFTER INSERT OR UPDATE OR DELETE ON snisid_identity.citizens
    FOR EACH ROW EXECUTE FUNCTION snisid_audit.audit_citizens_change();

-- ============================================================
-- VIEWS
-- ============================================================

-- Vue: Profil citoyen minimal (pour vérification externe)
CREATE VIEW snisid_identity.v_citizen_minimal AS
SELECT
    c.niu,
    c.nom,
    c.prenom,
    c.date_naissance,
    c.lieu_naissance_commune,
    c.lieu_naissance_departement,
    c.sexe,
    c.statut_identite,
    c.cert_auth_thumbprint,
    c.cert_expires_at,
    EXISTS(SELECT 1 FROM snisid_biometric.biometric_references b WHERE b.niu = c.niu AND b.is_active) AS biometric_enrolled
FROM snisid_identity.citizens c
WHERE c.statut_identite != 'ARCHIVED';

-- Vue: Statistiques d'enregistrement par département
CREATE VIEW snisid_identity.v_enrollment_stats AS
SELECT
    departement_residence AS departement,
    statut_identite,
    canal_enregistrement,
    COUNT(*) AS total,
    DATE_TRUNC('day', created_at) AS enrollment_date
FROM snisid_identity.citizens
GROUP BY departement_residence, statut_identite, canal_enregistrement, DATE_TRUNC('day', created_at);

-- ============================================================
-- INITIAL DATA
-- ============================================================

-- Centres d'état civil (10 chefs-lieux de département)
INSERT INTO snisid_civil.officiers_etat_civil (oec_id, nom_complet, commune_assignee, departement, date_nomination, statut)
VALUES
    ('OEC-OUEST-001', 'Jean-Baptiste PIERRE', 'Port-au-Prince', 'Ouest', '2026-01-01', 'ACTIF'),
    ('OEC-NORD-001', 'Marie-Claire JEAN', 'Cap-Haïtien', 'Nord', '2026-01-01', 'ACTIF'),
    ('OEC-ARTIB-001', 'Joseph BLANC', 'Gonaïves', 'Artibonite', '2026-01-01', 'ACTIF'),
    ('OEC-CENTRE-001', 'Rosette PAUL', 'Hinche', 'Centre', '2026-01-01', 'ACTIF'),
    ('OEC-SUD-001', 'Frantz LAGUERRE', 'Les Cayes', 'Sud', '2026-01-01', 'ACTIF'),
    ('OEC-GRANDANSE-001', 'Carline DESROCHES', 'Jérémie', 'Grande-Anse', '2026-01-01', 'ACTIF'),
    ('OEC-NIPE-001', 'Edner CELESTIN', 'Miragoâne', 'Nippes', '2026-01-01', 'ACTIF'),
    ('OEC-SUDEST-001', 'Mireille SIMON', 'Jacmel', 'Sud-Est', '2026-01-01', 'ACTIF'),
    ('OEC-NORDEST-001', 'Pierre HENRY', 'Fort-Liberté', 'Nord-Est', '2026-01-01', 'ACTIF'),
    ('OEC-NORDOUEST-001', 'Lise MORENCY', 'Port-de-Paix', 'Nord-Ouest', '2026-01-01', 'ACTIF');

-- ============================================================
-- MAINTENANCE
-- ============================================================

-- Fonction de purge RGPD (anonymisation, pas suppression physique)
CREATE OR REPLACE FUNCTION snisid_identity.gdpr_anonymize(p_niu VARCHAR(10))
RETURNS VOID AS $$
BEGIN
    UPDATE snisid_identity.citizens SET
        nom = 'ANONYMISE',
        prenom = 'ANONYMISE',
        nom_jeune_fille = NULL,
        date_naissance = '1900-01-01',
        adresse_ligne1 = NULL,
        adresse_ligne2 = NULL,
        adresse_gps_lat = NULL,
        adresse_gps_lon = NULL,
        gdpr_erasure_requested = TRUE,
        gdpr_erasure_date = NOW(),
        updated_at = NOW()
    WHERE niu = p_niu;

    -- Log l'anonymisation
    INSERT INTO snisid_audit.audit_trail (
        event_type, entity_type, entity_id, niu,
        agent_id, agency_id, action, status, audit_hash, timestamp
    ) VALUES (
        'GDPR_ANONYMIZATION', 'CITIZEN', p_niu, p_niu,
        'SYSTEM', 'NDPA', 'ANONYMIZE', 'SUCCESS',
        encode(sha256((p_niu || NOW()::text)::bytea), 'hex'),
        NOW()
    );
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- ============================================================
-- GRANTS
-- ============================================================
GRANT USAGE ON SCHEMA snisid_identity TO identity_service_role, registry_service_role;
GRANT USAGE ON SCHEMA snisid_civil TO civil_registry_service_role;
GRANT USAGE ON SCHEMA snisid_audit TO audit_service_role;
GRANT USAGE ON SCHEMA snisid_biometric TO biometric_service_role;

GRANT SELECT, INSERT, UPDATE ON snisid_identity.citizens TO identity_service_role;
GRANT SELECT ON snisid_identity.citizens TO registry_service_role;
GRANT INSERT ON snisid_identity.identity_events TO identity_service_role;
GRANT SELECT ON snisid_identity.identity_events TO registry_service_role, audit_service_role;
GRANT SELECT, INSERT, UPDATE ON snisid_biometric.biometric_references TO biometric_service_role;
GRANT SELECT, INSERT, UPDATE ON snisid_civil.civil_acts TO civil_registry_service_role;
GRANT SELECT, INSERT ON snisid_audit.audit_trail TO audit_service_role;

-- ============================================================
-- END OF SCHEMA
-- SNISID-IDN-DB-001 v1.0.0
-- Approuvé par: DBA National SNISID / AND / CISO
-- ============================================================
