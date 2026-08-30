-- Migration 007 : Index de performance et optimisations

-- Index composites pour requêtes fréquentes
CREATE INDEX IF NOT EXISTS idx_citizens_dept_niu
    ON snisid_identity.citizens (departement_residence, niu);

CREATE INDEX IF NOT EXISTS idx_citizens_status_created
    ON snisid_identity.citizens (statut_identite, created_at);

CREATE INDEX IF NOT EXISTS idx_audit_entity_agent
    ON snisid_audit.audit_trail (entity_type, agent_niu, created_at);

CREATE INDEX IF NOT EXISTS idx_dna_active_profiles
    ON snisid_bio_adn.dna_profiles (profile_type, status)
    WHERE is_active = TRUE;

CREATE INDEX IF NOT EXISTS idx_person_active_records
    ON snisid_bio_adn.person_records (record_type, priority_level)
    WHERE is_active = TRUE;

CREATE INDEX IF NOT EXISTS idx_property_stolen
    ON snisid_bio_adn.property_records (record_type, theft_dept)
    WHERE status = 'STOLEN' AND is_active = TRUE;

-- Index fulltext pour recherche personne
CREATE INDEX IF NOT EXISTS idx_person_subject_name_fts
    ON snisid_bio_adn.person_records
    USING gin(to_tsvector('french', COALESCE(subject_name, '')));

-- Index pour le feature cache audit
CREATE INDEX IF NOT EXISTS idx_feature_audit_updated
    ON snisid_ml.feature_cache_audit (updated_at DESC);

-- Statistiques pour l'optimiseur de requêtes
ANALYZE snisid_identity.citizens;
ANALYZE snisid_audit.audit_trail;
ANALYZE snisid_bio_adn.dna_profiles;
ANALYZE snisid_bio_adn.person_records;
ANALYZE snisid_bio_adn.property_records;
