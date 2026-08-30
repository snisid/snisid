package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/reso-svc/internal/domain"
)

type networkRepo struct {
	pool *pgxpool.Pool
}

func NewNetworkRepo(pool *pgxpool.Pool) *networkRepo {
	return &networkRepo{pool: pool}
}

func (r *networkRepo) GetAllPersons() ([]domain.Person, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT snisid_id, name, aliases, nationality, dob, risk_score,
		        is_gang_member, is_sanctioned, created_at, updated_at
		 FROM reso_persons ORDER BY risk_score DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var persons []domain.Person
	for rows.Next() {
		var p domain.Person
		if err := rows.Scan(
			&p.SnisidID, &p.Name, &p.Aliases, &p.Nationality, &p.DOB,
			&p.RiskScore, &p.IsGangMember, &p.IsSanctioned, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		persons = append(persons, p)
	}
	return persons, nil
}

func (r *networkRepo) GetAllGangs() ([]domain.Gang, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT gang_id, name, primary_activity, territory_dept, activity_level,
		        member_count, created_at
		 FROM reso_gangs ORDER BY member_count DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gangs []domain.Gang
	for rows.Next() {
		var g domain.Gang
		if err := rows.Scan(
			&g.GangID, &g.Name, &g.PrimaryActivity, &g.TerritoryDept,
			&g.ActivityLevel, &g.MemberCount, &g.CreatedAt,
		); err != nil {
			return nil, err
		}
		gangs = append(gangs, g)
	}
	return gangs, nil
}

func (r *networkRepo) GetPersonGangMemberships() ([]domain.PersonGang, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT person_id, gang_id, role, since, confidence
		 FROM reso_person_gang`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memberships []domain.PersonGang
	for rows.Next() {
		var m domain.PersonGang
		if err := rows.Scan(&m.PersonID, &m.GangID, &m.Role, &m.Since, &m.Confidence); err != nil {
			return nil, err
		}
		memberships = append(memberships, m)
	}
	return memberships, nil
}

func (r *networkRepo) GetGangRelations() ([]domain.GangRelation, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT gang_id_1, gang_id_2, relation_type, since, confidence
		 FROM reso_gang_relations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var relations []domain.GangRelation
	for rows.Next() {
		var rel domain.GangRelation
		if err := rows.Scan(&rel.GangID1, &rel.GangID2, &rel.RelationType, &rel.Since, &rel.Confidence); err != nil {
			return nil, err
		}
		relations = append(relations, rel)
	}
	return relations, nil
}

func (r *networkRepo) GetPersonAssociations() ([]domain.PersonAssociation, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT person_id_1, person_id_2, association_type, confidence, source, created_at
		 FROM reso_person_associations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var associations []domain.PersonAssociation
	for rows.Next() {
		var a domain.PersonAssociation
		if err := rows.Scan(&a.PersonID1, &a.PersonID2, &a.AssociationType, &a.Confidence, &a.Source, &a.CreatedAt); err != nil {
			return nil, err
		}
		associations = append(associations, a)
	}
	return associations, nil
}

func (r *networkRepo) GetNetworks() ([]domain.CriminalNetwork, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT network_id, network_name, member_ids, community_size,
		        modularity_score, detected_at, analysis_version
		 FROM reso_criminal_networks ORDER BY detected_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var networks []domain.CriminalNetwork
	for rows.Next() {
		var n domain.CriminalNetwork
		if err := rows.Scan(
			&n.NetworkID, &n.NetworkName, &n.MemberIDs, &n.CommunitySize,
			&n.ModularityScore, &n.DetectedAt, &n.AnalysisVersion,
		); err != nil {
			return nil, err
		}
		networks = append(networks, n)
	}
	return networks, nil
}

func (r *networkRepo) SaveNetwork(network *domain.CriminalNetwork) error {
	ctx := context.Background()
	_, err := r.pool.Exec(ctx,
		`INSERT INTO reso_criminal_networks
		 (network_id, network_name, member_ids, community_size, modularity_score, detected_at, analysis_version)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (network_id) DO UPDATE SET
		   network_name = EXCLUDED.network_name,
		   member_ids = EXCLUDED.member_ids,
		   community_size = EXCLUDED.community_size,
		   modularity_score = EXCLUDED.modularity_score,
		   analysis_version = EXCLUDED.analysis_version`,
		network.NetworkID, network.NetworkName, network.MemberIDs, network.CommunitySize,
		network.ModularityScore, network.DetectedAt, network.AnalysisVersion,
	)
	return err
}

func (r *networkRepo) GetPersonNetwork(personID uuid.UUID) ([]domain.PersonAssociation, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT person_id_1, person_id_2, association_type, confidence, source, created_at
		 FROM reso_person_associations
		 WHERE person_id_1 = $1 OR person_id_2 = $1`,
		personID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var associations []domain.PersonAssociation
	for rows.Next() {
		var a domain.PersonAssociation
		if err := rows.Scan(&a.PersonID1, &a.PersonID2, &a.AssociationType, &a.Confidence, &a.Source, &a.CreatedAt); err != nil {
			return nil, err
		}
		associations = append(associations, a)
	}
	return associations, nil
}

func (r *networkRepo) GetGangOverlap(gangID1, gangID2 uuid.UUID) ([]uuid.UUID, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT a.person_id
		 FROM reso_person_gang a
		 JOIN reso_person_gang b ON a.person_id = b.person_id
		 WHERE a.gang_id = $1 AND b.gang_id = $2`,
		gangID1, gangID2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commonMembers []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		commonMembers = append(commonMembers, id)
	}
	return commonMembers, nil
}

func scanPersons(rows pgx.Rows) ([]domain.Person, error) {
	var persons []domain.Person
	for rows.Next() {
		var p domain.Person
		if err := rows.Scan(
			&p.SnisidID, &p.Name, &p.Aliases, &p.Nationality, &p.DOB,
			&p.RiskScore, &p.IsGangMember, &p.IsSanctioned, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		persons = append(persons, p)
	}
	return persons, nil
}
