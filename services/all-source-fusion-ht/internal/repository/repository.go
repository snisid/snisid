package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/snisid/all-source-fusion-ht/internal/domain"
)

type FusionRepository interface {
	CreateProduct(ctx context.Context, p *domain.IntelProduct) error
	GetRecentProducts(ctx context.Context, limit int) ([]domain.IntelProduct, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*domain.IntelProduct, error)
	CreateThreatActor(ctx context.Context, a *domain.ThreatActor) error
	GetHighRiskActors(ctx context.Context) ([]domain.ThreatActor, error)
	CreateCorrelation(ctx context.Context, c *domain.CrossDisciplineCorrelation) error
	GetSourceMap(ctx context.Context, productID uuid.UUID) (*domain.IntelProduct, error)
	GetNationalEstimates(ctx context.Context) ([]domain.IntelProduct, error)
}

type fusionRepo struct {
	db *sql.DB
}

func NewFusionRepository(db *sql.DB) FusionRepository {
	return &fusionRepo{db: db}
}

func (r *fusionRepo) CreateProduct(ctx context.Context, p *domain.IntelProduct) error {
	query := `INSERT INTO fusion_intel_products
		(product_id, title, classification, source_disciplines, sigint_refs, humint_refs, geoint_refs,
		 osint_refs, analyst_assessment, confidence_level, related_threat_actors, related_regions,
		 nie_ref, created_by, approved_by, valid_from, valid_until)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err := r.db.ExecContext(ctx, query,
		p.ProductID, p.Title, p.Classification,
		pq.Array(p.SourceDisciplines), pq.Array(p.SigintRefs), pq.Array(p.HumintRefs),
		pq.Array(p.GeointRefs), pq.Array(p.OsintRefs),
		p.AnalystAssessment, p.ConfidenceLevel,
		pq.Array(p.RelatedThreatActors), pq.Array(p.RelatedRegions),
		p.NIERef, p.CreatedBy, p.ApprovedBy, p.ValidFrom, p.ValidUntil)
	return err
}

func (r *fusionRepo) GetRecentProducts(ctx context.Context, limit int) ([]domain.IntelProduct, error) {
	query := `SELECT product_id, title, classification, source_disciplines, sigint_refs, humint_refs,
		geoint_refs, osint_refs, analyst_assessment, confidence_level, related_threat_actors,
		related_regions, nie_ref, created_by, approved_by, valid_from, valid_until
		FROM fusion_intel_products ORDER BY valid_from DESC LIMIT $1`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProducts(rows)
}

func (r *fusionRepo) GetProductByID(ctx context.Context, id uuid.UUID) (*domain.IntelProduct, error) {
	query := `SELECT product_id, title, classification, source_disciplines, sigint_refs, humint_refs,
		geoint_refs, osint_refs, analyst_assessment, confidence_level, related_threat_actors,
		related_regions, nie_ref, created_by, approved_by, valid_from, valid_until
		FROM fusion_intel_products WHERE product_id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	return scanProduct(row)
}

func (r *fusionRepo) CreateThreatActor(ctx context.Context, a *domain.ThreatActor) error {
	query := `INSERT INTO fusion_threat_actors
		(actor_id, name, aliases, type, cap_level, intent_level, opportunity_level, overall_risk,
		 last_activity_at, primary_region, associated_groups, ofac_designated, notes)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`
	_, err := r.db.ExecContext(ctx, query,
		a.ActorID, a.Name, pq.Array(a.Aliases), a.Type, a.CapLevel, a.IntentLevel,
		a.OpportunityLevel, a.OverallRisk, a.LastActivityAt, a.PrimaryRegion,
		pq.Array(a.AssociatedGroups), a.OFACDesignated, a.Notes)
	return err
}

func (r *fusionRepo) GetHighRiskActors(ctx context.Context) ([]domain.ThreatActor, error) {
	query := `SELECT actor_id, name, aliases, type, cap_level, intent_level, opportunity_level,
		overall_risk, last_activity_at, primary_region, associated_groups, ofac_designated, notes
		FROM fusion_threat_actors WHERE overall_risk >= 4 ORDER BY overall_risk DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []domain.ThreatActor
	for rows.Next() {
		var a domain.ThreatActor
		if err := rows.Scan(&a.ActorID, &a.Name, pq.Array(&a.Aliases), &a.Type,
			&a.CapLevel, &a.IntentLevel, &a.OpportunityLevel, &a.OverallRisk,
			&a.LastActivityAt, &a.PrimaryRegion, pq.Array(&a.AssociatedGroups),
			&a.OFACDesignated, &a.Notes); err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	if result == nil {
		result = []domain.ThreatActor{}
	}
	return result, rows.Err()
}

func (r *fusionRepo) CreateCorrelation(ctx context.Context, c *domain.CrossDisciplineCorrelation) error {
	query := `INSERT INTO fusion_correlations
		(correlation_id, discipline_a, reference_a, discipline_b, reference_b,
		 correlation_type, analyst_notes, score, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := r.db.ExecContext(ctx, query,
		c.CorrelationID, c.DisciplineA, c.ReferenceA, c.DisciplineB, c.ReferenceB,
		c.CorrelationType, c.AnalystNotes, c.Score, c.CreatedAt)
	return err
}

func (r *fusionRepo) GetSourceMap(ctx context.Context, productID uuid.UUID) (*domain.IntelProduct, error) {
	return r.GetProductByID(ctx, productID)
}

func (r *fusionRepo) GetNationalEstimates(ctx context.Context) ([]domain.IntelProduct, error) {
	query := `SELECT product_id, title, classification, source_disciplines, sigint_refs, humint_refs,
		geoint_refs, osint_refs, analyst_assessment, confidence_level, related_threat_actors,
		related_regions, nie_ref, created_by, approved_by, valid_from, valid_until
		FROM fusion_intel_products WHERE nie_ref IS NOT NULL ORDER BY valid_from DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProducts(rows)
}

func scanProducts(rows *sql.Rows) ([]domain.IntelProduct, error) {
	var result []domain.IntelProduct
	for rows.Next() {
		var p domain.IntelProduct
		if err := rows.Scan(&p.ProductID, &p.Title, &p.Classification,
			pq.Array(&p.SourceDisciplines), pq.Array(&p.SigintRefs), pq.Array(&p.HumintRefs),
			pq.Array(&p.GeointRefs), pq.Array(&p.OsintRefs),
			&p.AnalystAssessment, &p.ConfidenceLevel,
			pq.Array(&p.RelatedThreatActors), pq.Array(&p.RelatedRegions),
			&p.NIERef, &p.CreatedBy, &p.ApprovedBy, &p.ValidFrom, &p.ValidUntil); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	if result == nil {
		result = []domain.IntelProduct{}
	}
	return result, rows.Err()
}

func scanProduct(row *sql.Row) (*domain.IntelProduct, error) {
	var p domain.IntelProduct
	if err := row.Scan(&p.ProductID, &p.Title, &p.Classification,
		pq.Array(&p.SourceDisciplines), pq.Array(&p.SigintRefs), pq.Array(&p.HumintRefs),
		pq.Array(&p.GeointRefs), pq.Array(&p.OsintRefs),
		&p.AnalystAssessment, &p.ConfidenceLevel,
		pq.Array(&p.RelatedThreatActors), pq.Array(&p.RelatedRegions),
		&p.NIERef, &p.CreatedBy, &p.ApprovedBy, &p.ValidFrom, &p.ValidUntil); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

var _ FusionRepository = (*fusionRepo)(nil)
