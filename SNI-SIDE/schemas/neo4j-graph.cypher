// ============================================================
// SNI-SIDE: National Sovereign Intelligence Graph
// Neo4j Enterprise 5.x — Graph Schema
// ============================================================

// ============ CONSTRAINTS ============
CREATE CONSTRAINT citizen_niu_unique IF NOT EXISTS FOR (c:Citizen) REQUIRE c.niu IS UNIQUE;
CREATE CONSTRAINT vehicle_vin_unique IF NOT EXISTS FOR (v:Vehicle) REQUIRE v.vin IS UNIQUE;
CREATE CONSTRAINT vehicle_plate_unique IF NOT EXISTS FOR (v:Vehicle) REQUIRE v.plate IS UNIQUE;
CREATE CONSTRAINT phone_number_unique IF NOT EXISTS FOR (p:Phone) REQUIRE p.number IS UNIQUE;
CREATE CONSTRAINT address_id_unique IF NOT EXISTS FOR (a:Address) REQUIRE a.id IS UNIQUE;
CREATE CONSTRAINT org_name_unique IF NOT EXISTS FOR (o:Organization) REQUIRE o.name IS UNIQUE;
CREATE CONSTRAINT bank_account_unique IF NOT EXISTS FOR (b:BankAccount) REQUIRE b.account_number IS UNIQUE;
CREATE CONSTRAINT passport_number_unique IF NOT EXISTS FOR (p:Passport) REQUIRE p.passport_number IS UNIQUE;
CREATE CONSTRAINT weapon_serial_unique IF NOT EXISTS FOR (w:Weapon) REQUIRE w.serial_number IS UNIQUE;
CREATE CONSTRAINT dna_profile_unique IF NOT EXISTS FOR (d:DNA) REQUIRE d.profile_id IS UNIQUE;
CREATE CONSTRAINT biometric_id_unique IF NOT EXISTS FOR (b:Biometric) REQUIRE b.biometric_id IS UNIQUE;
CREATE CONSTRAINT case_number_unique IF NOT EXISTS FOR (c:Case) REQUIRE c.case_number IS UNIQUE;
CREATE CONSTRAINT crossing_id_unique IF NOT EXISTS FOR (b:BorderCrossing) REQUIRE b.crossing_id IS UNIQUE;
CREATE CONSTRAINT document_number_unique IF NOT EXISTS FOR (d:Document) REQUIRE d.document_number IS UNIQUE;
CREATE CONSTRAINT device_id_unique IF NOT EXISTS FOR (d:Device) REQUIRE d.device_id IS UNIQUE;
CREATE CONSTRAINT ip_address_unique IF NOT EXISTS FOR (i:IP) REQUIRE i.address IS UNIQUE;
CREATE CONSTRAINT domain_name_unique IF NOT EXISTS FOR (d:Domain) REQUIRE d.name IS UNIQUE;
CREATE CONSTRAINT wallet_address_unique IF NOT EXISTS FOR (w:Wallet) REQUIRE w.address IS UNIQUE;
CREATE CONSTRAINT gang_name_unique IF NOT EXISTS FOR (g:Gang) REQUIRE g.name IS UNIQUE;
CREATE CONSTRAINT criminal_org_name_unique IF NOT EXISTS FOR (co:CriminalOrganization) REQUIRE co.name IS UNIQUE;
CREATE CONSTRAINT watchlist_entry_unique IF NOT EXISTS FOR (w:WatchlistEntry) REQUIRE w.entry_id IS UNIQUE;

// ============ INDEXES ============
CREATE INDEX citizen_name IF NOT EXISTS FOR (c:Citizen) ON (c.full_name);
CREATE INDEX citizen_nationality IF NOT EXISTS FOR (c:Citizen) ON (c.nationality);
CREATE INDEX citizen_risk IF NOT EXISTS FOR (c:Citizen) ON (c.risk_score);
CREATE INDEX vehicle_make_model IF NOT EXISTS FOR (v:Vehicle) ON (v.make, v.model);
CREATE INDEX phone_owner IF NOT EXISTS FOR (p:Phone) ON (p.owner_niu);
CREATE INDEX case_status IF NOT EXISTS FOR (c:Case) ON (c.status);
CREATE INDEX case_type IF NOT EXISTS FOR (c:Case) ON (c.case_type);
CREATE INDEX crossing_date IF NOT EXISTS FOR (b:BorderCrossing) ON (b.crossing_date);
CREATE INDEX wallet_risk IF NOT EXISTS FOR (w:Wallet) ON (w.risk_score);
CREATE INDEX ip_risk IF NOT EXISTS FOR (i:IP) ON (i.risk_score);

// ============ NODE LABELS ============
// 
// (:Citizen) — National Identity (NIU)
//   Properties: niu, full_name, alias[], date_of_birth, nationality, 
//               gender, risk_score, status, is_wanted, is_missing, 
//               is_pep, is_terrorist_watch
//
// (:Vehicle) — Vehicle Intelligence
//   Properties: vin, plate, make, model, year, color, status,
//               risk_score, last_seen
//
// (:Phone) — Phone Intelligence
//   Properties: number, imei, operator, owner_niu, risk_score,
//               is_prepaid, status, last_seen
//
// (:Address) — Location Intelligence
//   Properties: id, full_address, city, department, country,
//               latitude, longitude, risk_score
//
// (:Organization) — Entity Intelligence
//   Properties: name, type, registration_number, country,
//               status, risk_score
//
// (:BankAccount) — Financial Intelligence
//   Properties: account_number, bank_name, bank_country,
//               iban, swift, owner_niu, risk_score
//
// (:Passport) — Document Intelligence
//   Properties: passport_number, country, issue_date, expiry,
//               holder_niu, status, is_fraudulent
//
// (:Weapon) — Firearms Intelligence
//   Properties: serial_number, make, model, caliber, type,
//               owner_niu, status, in_evidence
//
// (:DNA) — DNA Intelligence
//   Properties: profile_id, profile_type, niu, sample_id,
//               lab_id, collection_date
//
// (:Biometric) — Biometric Intelligence
//   Properties: biometric_id, niu, type (face/fp/iris/voice),
//               quality_score, enrollment_date
//
// (:Case) — Criminal Case
//   Properties: case_number, case_type, status, jurisdiction,
//               lead_agency, incident_date, risk_level
//
// (:BorderCrossing) — Border Movement
//   Properties: crossing_id, niu, passport, direction, border_point,
//               crossing_date, method, risk_score
//
// (:Document) — Identity Document
//   Properties: document_number, document_type, country, holder_niu,
//               issue_date, expiry_date, is_fraudulent
//
// (:Device) — Digital Device
//   Properties: device_id, type, imei, mac_address, os,
//               owner_niu, risk_score
//
// (:IP) — IP Address
//   Properties: address, type (v4/v6), isp, country, risk_score,
//               first_seen, last_seen
//
// (:Domain) — Domain Intelligence
//   Properties: name, registrar, registration_date, is_malicious,
//               risk_score, threat_actor
//
// (:Wallet) — Cryptocurrency Wallet
//   Properties: address, blockchain, type, risk_score,
//               total_received, total_sent, transaction_count
//
// (:Gang) — Criminal Gang
//   Properties: name, territory, criminal_activities[], status,
//               member_count, risk_level
//
// (:CriminalOrganization) — Organized Crime
//   Properties: name, type, country, status, risk_level
//
// (:WatchlistEntry) — National Watchlist
//   Properties: entry_id, type, category, risk_level, status,
//               listing_authority, confidence
//
// (:Incident) — Cyber/Physical Incident
//   Properties: incident_id, type, severity, date, location,
//               status, description
//
// (:Evidence) — Digital Evidence
//   Properties: evidence_id, type, file_hash, case_number,
//               captured_date, status
//
// (:SocialMedia) — Social Media Profile
//   Properties: platform, profile_id, username, display_name,
//               url, risk_score

// ============ RELATIONSHIPS ============
//
// Ownership and possession
// (:Citizen)-[:OWNS]->(:Vehicle)          — Vehicle ownership
// (:Citizen)-[:OWNS]->(:Phone)            — Phone ownership
// (:Citizen)-[:OWNS]->(:Weapon)           — Weapon ownership
// (:Citizen)-[:OWNS]->(:BankAccount)      — Bank account ownership
// (:Citizen)-[:OWNS]->(:Passport)         — Passport ownership
// (:Citizen)-[:OWNS]->(:Document)         — Document ownership
// (:Citizen)-[:OWNS]->(:Device)           — Device ownership
// (:Citizen)-[:OWNS]->(:Wallet)           — Crypto wallet ownership
// (:Organization)-[:OWNS]->(:Vehicle)      — Corporate vehicle
// (:Organization)-[:OWNS]->(:BankAccount)  — Corporate account
//
// Usage and association
// (:Citizen)-[:USES]->(:Phone)            — Phone usage
// (:Citizen)-[:USES]->(:IP)              — IP usage
// (:Citizen)-[:USES]->(:Device)           — Device usage
// (:Citizen)-[:USES]->(:SocialMedia)      — Social media usage
// (:Citizen)-[:USES]->(:Wallet)           — Wallet usage
// (:Vehicle)-[:USES]->(:Phone)            — Vehicle phone usage
// (:Citizen)-[:USES]->(:Domain)           — Domain usage (email)
//
// Visitation and location
// (:Citizen)-[:VISITED]->(:Address)       — Residential address
// (:Citizen)-[:VISITED]->(:BorderCrossing) — Border crossing
// (:Vehicle)-[:VISITED]->(:Address)        — Vehicle location
// (:Phone)-[:VISITED]->(:Address)          — Phone location history
//
// Association
// (:Citizen)-[:ASSOCIATED_WITH]->(:Citizen) — Person-to-person link
// (:Citizen)-[:ASSOCIATED_WITH]->(:Organization) — Affiliation
// (:Citizen)-[:ASSOCIATED_WITH]->(:Gang) — Gang membership
// (:Citizen)-[:ASSOCIATED_WITH]->(:CriminalOrganization) — Org membership
// (:Gang)-[:ASSOCIATED_WITH]->(:Gang) — Gang rivalry/alliance
// (:Organization)-[:ASSOCIATED_WITH]->(:Organization) — Corporate links
//
// Communication
// (:Phone)-[:CONNECTED_TO]->(:Phone)       — Call records
// (:Device)-[:CONNECTED_TO]->(:Device)     — Device-to-device
// (:IP)-[:CONNECTED_TO]->(:IP)             — Network connections
// (:IP)-[:CONNECTED_TO]->(:Domain)         — DNS resolution
// (:Phone)-[:CONNECTED_TO]->(:Citizen)     — Contact list
//
// Financial
// (:Citizen)-[:FINANCED_BY]->(:Organization) — Funding source
// (:Organization)-[:FINANCED_BY]->(:Organization) — Funding chain
// (:Citizen)-[:FINANCED_BY]->(:BankAccount) — Income source
// (:Citizen)-[:FINANCED_BY]->(:Citizen)    — Person-to-person transfers
// (:BankAccount)-[:FINANCED_BY]->(:BankAccount) — Money flow
// (:Wallet)-[:FINANCED_BY]->(:Wallet)      — Crypto flow
//
// Employment
// (:Citizen)-[:WORKS_FOR]->(:Organization) — Employment
// (:Citizen)-[:WORKS_FOR]->(:Citizen)      — Employer-employee
//
// Travel
// (:Citizen)-[:TRAVELLED_WITH]->(:Citizen)  — Travel companions
// (:Citizen)-[:TRAVELLED_WITH]->(:Vehicle)   — Travel vehicle
// (:Citizen)-[:TRAVELLED_ON]->(:Passport)   — Passport used
//
// Case links
// (:Citizen)-[:LINKED_TO]->(:Case)          — Person involved in case
// (:Vehicle)-[:LINKED_TO]->(:Case)          — Vehicle involved in case
// (:Phone)-[:LINKED_TO]->(:Case)            — Phone involved in case
// (:Weapon)-[:LINKED_TO]->(:Case)           — Weapon involved in case
// (:DNA)-[:LINKED_TO]->(:Case)              — DNA linked to case
// (:Evidence)-[:LINKED_TO]->(:Case)         — Evidence in case
// (:WatchlistEntry)-[:LINKED_TO]->(:Case)   — Watchlist linked to case
//
// Identification
// (:Citizen)-[:HAS_BIOMETRIC]->(:Biometric) — Biometric reference
// (:Citizen)-[:HAS_DNA]->(:DNA)             — DNA profile
// (:Citizen)-[:HAS_ALIAS]->(:Citizen)       — Alias identity
//
// Incident correlation
// (:Case)-[:CORRELATED_WITH]->(:Case)       — Case-to-case link
// (:Incident)-[:CORRELATED_WITH]->(:Case)   — Incident-to-case link
// (:IP)-[:INVOLVED_IN]->(:Incident)         — IP involved in incident

// ============ EXAMPLE GRAPH QUERIES ============

// --- 1. Find criminal network of a person (degrees of separation) ---
// MATCH (c:Citizen {niu: '0000000001'})
// CALL apoc.path.expand(c, 'ASSOCIATED_WITH|WORKS_FOR|FINANCED_BY|OWNS', '>Citizen|Organization|Gang|BankAccount', 1, 4)
// YIELD path
// RETURN path

// --- 2. Detect money laundering rings ---
// MATCH path = (c:Citizen)-[:FINANCED_BY]->(ba:BankAccount)<-[:FINANCED_BY]-(c2:Citizen)-[:FINANCED_BY]->(ba2:BankAccount)<-[:FINANCED_BY]-(c3:Citizen)
// WHERE c.risk_score > 0.7 AND c2.risk_score > 0.7
// RETURN path

// --- 3. Find common associates between two persons ---
// MATCH (c1:Citizen {niu: '0000000001'})-[:ASSOCIATED_WITH]-(common:Citizen)-[:ASSOCIATED_WITH]-(c2:Citizen {niu: '0000000002'})
// RETURN common

// --- 4. Vehicle network analysis ---
// MATCH (v:Vehicle {plate: 'AA-1234-BB'})<-[:OWNS]-(c:Citizen)-[:ASSOCIATED_WITH]->(c2:Citizen)-[:OWNS]->(v2:Vehicle)
// RETURN v, c, c2, v2

// --- 5. Cross-border criminal movement ---
// MATCH (c:Citizen)-[:VISITED]->(bc:BorderCrossing)-[:VISITED]-(c2:Citizen)
// WHERE bc.crossing_date > datetime('2026-01-01') AND bc.risk_score > 0.6
// RETURN c, c2, bc

// --- 6. Phone network intelligence ---
// MATCH (p:Phone {number: '+509XXXXXXXX'})<-[:USES]-(c:Citizen)-[:ASSOCIATED_WITH]->(c2:Citizen)-[:OWNS]->(p2:Phone)
// RETURN p, c, c2, p2

// --- 7. Combined narco-terrorism link ---
// MATCH (c:Citizen)-[:OWNS]->(v:Vehicle)
// MATCH (c)-[:ASSOCIATED_WITH]->(org:CriminalOrganization {type: 'CARTEL'})
// MATCH (c)-[:LINKED_TO]->(case:Case {case_type: 'NARCOTICS'})
// MATCH (v)-[:LINKED_TO]->(case)
// RETURN c, v, org, case

// --- 8. Watchlist hit with full context ---
// MATCH (w:WatchlistEntry {status: 'ACTIVE', risk_level: 'CRITICAL'})-[:LINKED_TO]->(case:Case)
// MATCH (w)-[:MATCHES]->(c:Citizen)
// OPTIONAL MATCH (c)-[:OWNS]->(v:Vehicle)
// OPTIONAL MATCH (c)-[:USES]->(p:Phone)
// OPTIONAL MATCH (c)-[:VISITED]->(bc:BorderCrossing)
// RETURN w, c, case, v, p, bc

// --- 9. Financial flow analysis (PEP) ---
// MATCH (pep:Citizen {is_pep: true})-[:FINANCED_BY]->(ba:BankAccount)
// MATCH (ba)<-[:FINANCED_BY]-(c2:Citizen)
// MATCH (c2)-[:FINANCED_BY]->(ba2:BankAccount)
// WHERE ba2.country IN ['PAN', 'CAY', 'UAE', 'CHE']
// RETURN pep, ba, c2, ba2

// --- 10. Digital forensics correlation ---
// MATCH (e:Evidence)-[:LINKED_TO]->(case:Case)
// MATCH (e)<-[:LINKED_TO]-(device:Device)
// MATCH (device)-[:USED_BY]->(c:Citizen)
// MATCH (ip:IP)-[:CONNECTED_TO]->(device)
// RETURN e, case, device, c, ip
