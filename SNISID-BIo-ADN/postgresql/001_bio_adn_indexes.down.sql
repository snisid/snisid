-- Rollback BIO-ADN indexes
DROP TABLE IF EXISTS snisid_bio_adn.property_records;
DROP TYPE IF EXISTS snisid_bio_adn.property_status;
DROP TYPE IF EXISTS snisid_bio_adn.property_record_type;
DROP TABLE IF EXISTS snisid_bio_adn.person_records;
DROP TYPE IF EXISTS snisid_bio_adn.person_status;
DROP TYPE IF EXISTS snisid_bio_adn.person_record_type;
DROP TABLE IF EXISTS snisid_bio_adn.dna_matches;
DROP TABLE IF EXISTS snisid_bio_adn.dna_profiles;
DROP TYPE IF EXISTS snisid_bio_adn.dna_match_type;
DROP TYPE IF EXISTS snisid_bio_adn.dna_profile_status;
DROP SCHEMA IF EXISTS snisid_bio_adn CASCADE;
