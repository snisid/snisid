#!/bin/bash
###############################################################################
# SNISID ENTITY RESOLUTION ENGINE
# Inspiré du système "Palantir Gotham" (FBI/NSA) et "SkyNet" (Chine)
# Objectif: Fusionner les identités, détecter les doublons, mapper les réseaux criminels
###############################################################################

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LOG_FILE="${SCRIPT_DIR}/logs/entity_resolution.log"
DATA_DIR="${SCRIPT_DIR}/../data/graph"
NEO4J_HOST="localhost"
NEO4J_PORT="7687"
NEO4J_USER="neo4j"
NEO4J_PASS="${SNISID_GRAPH_PASS:-ChangeMe123!}"

# Couleurs pour les logs
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log() {
    local level=$1
    shift
    local message="$@"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo -e "${timestamp} [${level}] ${message}" | tee -a "$LOG_FILE"
}

log_info() { log "${BLUE}INFO${NC}" "$@"; }
log_success() { log "${GREEN}SUCCESS${NC}" "$@"; }
log_warning() { log "${YELLOW}WARNING${NC}" "$@"; }
log_error() { log "${RED}ERROR${NC}" "$@"; }

# Vérifier les dépendances
check_dependencies() {
    log_info "Vérification des dépendances..."
    
    if ! command -v cypher-shell &> /dev/null; then
        log_warning "cypher-shell non trouvé. Installation requise pour Neo4j."
        log_info "Téléchargement de cypher-shell..."
        wget -q https://neo4j.com/artifact.php?name=neo4j-admin-5.14.0-unix.tar.gz -O /tmp/neo4j-admin.tar.gz
        tar -xzf /tmp/neo4j-admin.tar.gz -C /opt/
        ln -sf /opt/neo4j-admin-5.14.0/bin/cypher-shell /usr/local/bin/cypher-shell 2>/dev/null || true
    fi
    
    log_success "Dépendances vérifiées."
}

# Initialiser le schéma de graphe (Style FBI Link Analysis)
init_graph_schema() {
    log_info "Initialisation du schéma de graphe SNISID..."
    
    cypher-shell -u "$NEO4J_USER" -p "$NEO4J_PASS" -a "bolt://${NEO4J_HOST}:${NEO4J_PORT}" <<EOF
// ============================================================================
// SCHÉMA DE BASE DE DONNÉES GRAPHES - STYLE FBI/DEA/NSA
// ============================================================================

// 1. CRÉATION DES CONTRAINTES D'UNICITÉ (Data Quality)
CREATE CONSTRAINT person_id_unique IF NOT EXISTS FOR (p:Person) REQUIRE p.id IS UNIQUE;
CREATE CONSTRAINT person_nin_unique IF NOT EXISTS FOR (p:Person) REQUIRE p.nin IS UNIQUE;
CREATE CONSTRAINT organization_name_unique IF NOT EXISTS FOR (o:Organization) REQUIRE o.name IS UNIQUE;
CREATE CONSTRAINT phone_number_unique IF NOT EXISTS FOR (ph:PhoneNumber) REQUIRE ph.number IS UNIQUE;
CREATE CONSTRAINT email_address_unique IF NOT EXISTS FOR (em:EmailAddress) REQUIRE em.address IS UNIQUE;
CREATE CONSTRAINT device_imei_unique IF NOT EXISTS FOR (d:Device) REQUIRE d.imei IS UNIQUE;
CREATE CONSTRAINT vehicle_vin_unique IF NOT EXISTS FOR (v:Vehicle) REQUIRE d.vin IS UNIQUE;
CREATE CONSTRAINT address_full_unique IF NOT EXISTS FOR (a:Address) REQUIRE a.full_address IS UNIQUE;

// 2. CRÉATION DES INDEX POUR LA PERFORMANCE (Recherche rapide)
CREATE INDEX person_name_idx IF NOT EXISTS FOR (p:Person) ON (p.last_name, p.first_name);
CREATE INDEX person_dob_idx IF NOT EXISTS FOR (p:Person) ON (p.date_of_birth);
CREATE INDEX organization_type_idx IF NOT EXISTS FOR (o:Organization) ON (o.type);
CREATE INDEX event_date_idx IF NOT EXISTS FOR (e:Event) ON (e.timestamp);
CREATE INDEX location_coords_idx IF NOT EXISTS FOR (l:Location) ON (l.latitude, l.longitude);

// 3. CRÉATION DES NODE LABELS (Typologie FBI/DEA)
// Personnes: Suspect, Informant, Victim, Witness, Associate, PEP (Politically Exposed Person)
// Organisations: Gang, Cartel, NGO, Government, Business, ShellCompany
// Événements: Meeting, Call, Transaction, Travel, Arrest, Surveillance

log_success "Schéma de graphe initialisé avec succès.";
EOF
}

# Algorithme de Entity Resolution (Fusion d'identités)
run_entity_resolution() {
    log_info "Lancement de l'algorithme de Entity Resolution..."
    
    cypher-shell -u "$NEO4J_USER" -p "$NEO4J_PASS" -a "bolt://${NEO4J_HOST}:${NEO4J_PORT}" <<EOF
// ============================================================================
// ALGORITHME DE FUSION D'IDENTITÉS (Entity Resolution)
// Détecte les doublons et fusionne les profils (Style Palantir/Matraix)
// ============================================================================

// ÉTAPE 1: DÉTECTION DES DOUBLONS POTENTIELS
// Critères: Même nom + même date de naissance OU même NIN
MATCH (p1:Person), (p2:Person)
WHERE p1.id < p2.id
AND (
    (p1.nin IS NOT NULL AND p2.nin IS NOT NULL AND p1.nin = p2.nin)
    OR 
    (p1.last_name = p2.last_name AND p1.first_name = p2.first_name AND p1.date_of_birth = p2.date_of_birth)
)
MERGE (p1)-[:POTENTIAL_DUPLICATE {score: 0.95, algorithm: 'exact_match'}]-(p2);

// ÉTAPE 2: DÉTECTION DES LIENS FAIBLES (Fuzzy Matching)
// Critères: Nom similaire (Levenshtein) + même numéro de téléphone ou adresse
MATCH (p1:Person)-[:HAS_PHONE]->(ph:PhoneNumber)<-[:HAS_PHONE]-(p2:Person)
WHERE p1.id < p2.id
AND p1.last_name <> p2.last_name
MERGE (p1)-[:POTENTIAL_DUPLICATE {score: 0.75, algorithm: 'shared_phone', reason: 'Partage même numéro de téléphone'}]-(p2);

// ÉTAPE 3: CRÉATION DE SUPER-NŒUDS (Identity Clusters)
// Regroupe toutes les entités liées en un cluster unique
CALL gds.graph.project('identityGraph', ['Person'], ['POTENTIAL_DUPLICATE'])
YIELD graphName, nodeCount, relationshipCount;

CALL gds.wcc.stream('identityGraph')
YIELD nodeId, componentId
WITH gds.util.asNode(nodeId) AS person, componentId
MERGE (c:IdentityCluster {id: componentId})
MERGE (person)-[:BELONGS_TO_CLUSTER]->(c);

log_success "Entity Resolution complétée. Clusters d'identité créés.";
EOF
}

# Analyse de réseaux criminels (Link Analysis)
run_link_analysis() {
    log_info "Lancement de l'analyse de liens criminels..."
    
    cypher-shell -u "$NEO4J_USER" -p "$NEO4J_PASS" -a "bolt://${NEO4J_HOST}:${NEO4J_PORT}" <<EOF
// ============================================================================
// LINK ANALYSIS - CARTOGRAPHIE DES RÉSEAUX CRIMINELS
// Algorithmes: PageRank, Betweenness Centrality, Community Detection
// ============================================================================

// 1. PAGE RANK: Identifier les leaders/influenceurs du réseau
CALL gds.pageRank.stream('identityGraph', {maxIterations: 20, dampingFactor: 0.85})
YIELD nodeId, score
WITH gds.util.asNode(nodeId) AS person, score
WHERE score > 0.05
SET person.pageRankScore = score
RETURN person.nin AS NIN, person.last_name AS Nom, person.pageRankScore AS Score
ORDER BY person.pageRankScore DESC
LIMIT 100;

// 2. BETWEENNESS CENTRALITY: Identifier les intermédiaires critiques (brokers)
CALL gds.betweenness.stream('identityGraph')
YIELD nodeId, score
WITH gds.util.asNode(nodeId) AS person, score
WHERE score > 0.01
SET person.betweennessScore = score
RETURN person.nin AS NIN, person.last_name AS Nom, person.betweennessScore AS Score
ORDER BY person.betweennessScore DESC
LIMIT 50;

// 3. LOUVAIN COMMUNITY DETECTION: Identifier les gangs/cellules criminelles
CALL gds.louvain.stream('identityGraph')
YIELD nodeId, communityId
WITH gds.util.asNode(nodeId) AS person, communityId
MERGE (g:CriminalGang {id: communityId})
MERGE (person)-[:MEMBER_OF_GANG]->(g)
RETURN g.id AS GangID, count(person) AS Membres
ORDER BY Membres DESC;

log_success "Analyse de liens criminels terminée.";
EOF
}

# Génération de rapports pour les agences
generate_intel_report() {
    log_info "Génération du rapport de renseignement..."
    
    local report_file="${SCRIPT_DIR}/../reports/intel_report_$(date +%Y%m%d_%H%M%S).json"
    
    cypher-shell -u "$NEO4J_USER" -p "$NEO4J_PASS" -a "bolt://${NEO4J_HOST}:${NEO4J_PORT}" --format json <<EOF
MATCH (g:CriminalGang)<-[:MEMBER_OF_GANG]-(m:Person)
WITH g.id AS gangId, collect({
    nin: m.nin,
    name: m.first_name + ' ' + m.last_name,
    dob: m.date_of_birth,
    pageRank: m.pageRankScore,
    role: CASE 
        WHEN m.pageRankScore > 0.1 THEN 'LEADER'
        WHEN m.betweennessScore > 0.05 THEN 'BROKER'
        ELSE 'MEMBER'
    END
}) AS members
RETURN {
    gang_id: gangId,
    member_count: size(members),
    leaders: [m IN members WHERE m.role = 'LEADER'],
    brokers: [m IN members WHERE m.role = 'BROKER'],
    full_roster: members
} AS gang_profile
ORDER BY size(members) DESC;
EOF
    log_success "Rapport généré: ${report_file}"
}

# Main execution
main() {
    log_info "=========================================="
    log_info "SNISID Entity Resolution Engine v2.0"
    log_info "Inspiré: FBI Palantir / NSA DataWave / Chine SkyNet"
    log_info "=========================================="
    
    mkdir -p "${SCRIPT_DIR}/logs" "${SCRIPT_DIR}/../reports"
    
    check_dependencies
    
    # Démarrer la chaîne de traitement
    init_graph_schema
    run_entity_resolution
    run_link_analysis
    generate_intel_report
    
    log_success "=========================================="
    log_success "Traitement terminé avec succès!"
    log_success "Prochaine étape: Validation par les analystes du CNS (Centre National de Sécurité)"
    log_success "=========================================="
}

# Exécution
main "$@"
