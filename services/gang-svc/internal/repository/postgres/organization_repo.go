package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/gang-svc/internal/domain"
)

type OrganizationRepo struct {
	pool *pgxpool.Pool
}

func NewOrganizationRepo(pool *pgxpool.Pool) *OrganizationRepo {
	return &OrganizationRepo{pool: pool}
}

func (r *OrganizationRepo) Create(ctx context.Context, org *domain.Organization) error {
	query := `
		INSERT INTO gang_organizations 
			(gang_id, national_gang_id, name, aliases, structure_type, primary_activity,
			 activity_level, estimated_members, armed_members_pct, heavy_weapons,
			 primary_dept_code, territory_communes, territory_geojson,
			 estimated_revenue_usd_monthly, primary_income_sources, un_designation_date,
			 ofac_designation, ofac_sdn_ref, allied_gang_ids, rival_gang_ids,
			 established_date, current_leader_id, intel_confidence, last_intel_update,
			 is_active, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
				$17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28)
	`
	_, err := r.pool.Exec(ctx, query,
		org.GangID, org.NationalGangID, org.Name, org.Aliases, org.StructureType,
		org.PrimaryActivity, org.ActivityLevel, org.EstimatedMembers,
		org.ArmedMembersPct, org.HeavyWeapons, org.PrimaryDeptCode,
		org.TerritoryCommunes, org.TerritoryGeoJSON, org.EstimatedRevenueUSDMonthly,
		org.PrimaryIncomeSources, org.UNDesignationDate, org.OFACDesignation,
		org.OFACSDNRef, org.AlliedGangIDs, org.RivalGangIDs, org.EstablishedDate,
		org.CurrentLeaderID, org.IntelConfidence, org.LastIntelUpdate,
		org.IsActive, org.CreatedBy, org.CreatedAt, org.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create organization: %w", err)
	}
	return nil
}

func (r *OrganizationRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	query := `
		SELECT gang_id, national_gang_id, name, aliases, structure_type, primary_activity,
			   activity_level, estimated_members, armed_members_pct, heavy_weapons,
			   primary_dept_code, territory_communes, territory_geojson,
			   estimated_revenue_usd_monthly, primary_income_sources, un_designation_date,
			   ofac_designation, ofac_sdn_ref, allied_gang_ids, rival_gang_ids,
			   established_date, current_leader_id, intel_confidence, last_intel_update,
			   is_active, created_by, created_at, updated_at
		FROM gang_organizations
		WHERE gang_id = $1
	`
	org := &domain.Organization{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&org.GangID, &org.NationalGangID, &org.Name, &org.Aliases, &org.StructureType,
		&org.PrimaryActivity, &org.ActivityLevel, &org.EstimatedMembers,
		&org.ArmedMembersPct, &org.HeavyWeapons, &org.PrimaryDeptCode,
		&org.TerritoryCommunes, &org.TerritoryGeoJSON, &org.EstimatedRevenueUSDMonthly,
		&org.PrimaryIncomeSources, &org.UNDesignationDate, &org.OFACDesignation,
		&org.OFACSDNRef, &org.AlliedGangIDs, &org.RivalGangIDs, &org.EstablishedDate,
		&org.CurrentLeaderID, &org.IntelConfidence, &org.LastIntelUpdate,
		&org.IsActive, &org.CreatedBy, &org.CreatedAt, &org.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find organization: %w", err)
	}
	return org, nil
}

func (r *OrganizationRepo) FindAll(ctx context.Context) ([]*domain.Organization, error) {
	query := `
		SELECT gang_id, national_gang_id, name, aliases, structure_type, primary_activity,
			   activity_level, estimated_members, armed_members_pct, heavy_weapons,
			   primary_dept_code, territory_communes, territory_geojson,
			   estimated_revenue_usd_monthly, primary_income_sources, un_designation_date,
			   ofac_designation, ofac_sdn_ref, allied_gang_ids, rival_gang_ids,
			   established_date, current_leader_id, intel_confidence, last_intel_update,
			   is_active, created_by, created_at, updated_at
		FROM gang_organizations
		WHERE is_active = TRUE
		ORDER BY activity_level DESC, estimated_members DESC
	`
	return r.queryOrganizations(ctx, query)
}

func (r *OrganizationRepo) FindByDeptCode(ctx context.Context, deptCode string) ([]*domain.Organization, error) {
	query := `
		SELECT gang_id, national_gang_id, name, aliases, structure_type, primary_activity,
			   activity_level, estimated_members, armed_members_pct, heavy_weapons,
			   primary_dept_code, territory_communes, territory_geojson,
			   estimated_revenue_usd_monthly, primary_income_sources, un_designation_date,
			   ofac_designation, ofac_sdn_ref, allied_gang_ids, rival_gang_ids,
			   established_date, current_leader_id, intel_confidence, last_intel_update,
			   is_active, created_by, created_at, updated_at
		FROM gang_organizations
		WHERE primary_dept_code = $1 AND is_active = TRUE
		ORDER BY activity_level DESC
	`
	return r.queryOrganizations(ctx, query, deptCode)
}

func (r *OrganizationRepo) FindSanctioned(ctx context.Context) ([]*domain.Organization, error) {
	query := `
		SELECT gang_id, national_gang_id, name, aliases, structure_type, primary_activity,
			   activity_level, estimated_members, armed_members_pct, heavy_weapons,
			   primary_dept_code, territory_communes, territory_geojson,
			   estimated_revenue_usd_monthly, primary_income_sources, un_designation_date,
			   ofac_designation, ofac_sdn_ref, allied_gang_ids, rival_gang_ids,
			   established_date, current_leader_id, intel_confidence, last_intel_update,
			   is_active, created_by, created_at, updated_at
		FROM gang_organizations
		WHERE ofac_designation = TRUE OR un_designation_date IS NOT NULL
		ORDER BY activity_level DESC
	`
	return r.queryOrganizations(ctx, query)
}

func (r *OrganizationRepo) Update(ctx context.Context, org *domain.Organization) error {
	query := `
		UPDATE gang_organizations
		SET name = $3, aliases = $4, structure_type = $5, primary_activity = $6,
			activity_level = $7, estimated_members = $8, armed_members_pct = $9,
			heavy_weapons = $10, territory_communes = $11, territory_geojson = $12,
			estimated_revenue_usd_monthly = $13, primary_income_sources = $14,
			ofac_designation = $15, ofac_sdn_ref = $16, allied_gang_ids = $17,
			rival_gang_ids = $18, current_leader_id = $19, intel_confidence = $20,
			last_intel_update = $21, updated_at = $22
		WHERE gang_id = $1 AND national_gang_id = $2
	`
	_, err := r.pool.Exec(ctx, query,
		org.GangID, org.NationalGangID, org.Name, org.Aliases, org.StructureType,
		org.PrimaryActivity, org.ActivityLevel, org.EstimatedMembers,
		org.ArmedMembersPct, org.HeavyWeapons, org.TerritoryCommunes,
		org.TerritoryGeoJSON, org.EstimatedRevenueUSDMonthly, org.PrimaryIncomeSources,
		org.OFACDesignation, org.OFACSDNRef, org.AlliedGangIDs, org.RivalGangIDs,
		org.CurrentLeaderID, org.IntelConfidence, org.LastIntelUpdate, org.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update organization: %w", err)
	}
	return nil
}

func (r *OrganizationRepo) queryOrganizations(ctx context.Context, query string, args ...interface{}) ([]*domain.Organization, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query organizations: %w", err)
	}
	defer rows.Close()

	var orgs []*domain.Organization
	for rows.Next() {
		org := &domain.Organization{}
		err := rows.Scan(
			&org.GangID, &org.NationalGangID, &org.Name, &org.Aliases, &org.StructureType,
			&org.PrimaryActivity, &org.ActivityLevel, &org.EstimatedMembers,
			&org.ArmedMembersPct, &org.HeavyWeapons, &org.PrimaryDeptCode,
			&org.TerritoryCommunes, &org.TerritoryGeoJSON, &org.EstimatedRevenueUSDMonthly,
			&org.PrimaryIncomeSources, &org.UNDesignationDate, &org.OFACDesignation,
			&org.OFACSDNRef, &org.AlliedGangIDs, &org.RivalGangIDs, &org.EstablishedDate,
			&org.CurrentLeaderID, &org.IntelConfidence, &org.LastIntelUpdate,
			&org.IsActive, &org.CreatedBy, &org.CreatedAt, &org.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan organization: %w", err)
		}
		orgs = append(orgs, org)
	}
	return orgs, nil
}
