package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/snisid/dr-svc/internal/domain"
)

type Repository interface {
	InsertRegion(ctx context.Context, region *domain.DRRegion) error
	FindAllRegions(ctx context.Context) ([]domain.DRRegion, error)
	FindRegionByName(ctx context.Context, name string) (*domain.DRRegion, error)
	UpdateRegionHealth(ctx context.Context, name string, health domain.HealthStatus) error
	UpdateRegionActive(ctx context.Context, name string, active bool) error
	InsertFailoverPlan(ctx context.Context, plan *domain.FailoverPlan) error
	FindAllFailoverPlans(ctx context.Context) ([]domain.FailoverPlan, error)
	FindFailoverPlanByID(ctx context.Context, planID uuid.UUID) (*domain.FailoverPlan, error)
	UpdateFailoverPlanExecuted(ctx context.Context, planID uuid.UUID) error
	InsertFailoverExecution(ctx context.Context, exec *domain.FailoverExecution) error
	FindExecutionsByPlan(ctx context.Context, planID uuid.UUID) ([]domain.FailoverExecution, error)
	InsertBackupManifest(ctx context.Context, manifest *domain.BackupManifest) error
	FindAllBackupManifests(ctx context.Context) ([]domain.BackupManifest, error)
	FindBackupManifestByID(ctx context.Context, manifestID uuid.UUID) (*domain.BackupManifest, error)
	InsertRecoveryPoint(ctx context.Context, point *domain.RecoveryPoint) error
	InsertDRTestResult(ctx context.Context, result *domain.DRTestResult) error
	FindAllDRTestResults(ctx context.Context) ([]domain.DRTestResult, error)
	InsertReplicationStatus(ctx context.Context, status *domain.ReplicationStatus) error
	FindReplicationStatus(ctx context.Context) ([]domain.ReplicationStatus, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) InsertRegion(ctx context.Context, region *domain.DRRegion) error {
	query := `INSERT INTO dr_regions (region_id, name, endpoint, is_active, health, last_checked, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query,
		region.RegionID, region.Name, region.Endpoint, region.IsActive, region.Health, region.LastChecked, region.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert region: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindAllRegions(ctx context.Context) ([]domain.DRRegion, error) {
	query := `SELECT region_id, name, endpoint, is_active, health, last_checked, created_at FROM dr_regions ORDER BY name`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query regions: %w", err)
	}
	defer rows.Close()
	var regions []domain.DRRegion
	for rows.Next() {
		var reg domain.DRRegion
		if err := rows.Scan(&reg.RegionID, &reg.Name, &reg.Endpoint, &reg.IsActive, &reg.Health, &reg.LastChecked, &reg.CreatedAt); err != nil {
			return nil, err
		}
		regions = append(regions, reg)
	}
	return regions, rows.Err()
}

func (r *postgresRepo) FindRegionByName(ctx context.Context, name string) (*domain.DRRegion, error) {
	query := `SELECT region_id, name, endpoint, is_active, health, last_checked, created_at FROM dr_regions WHERE name = $1`
	reg := &domain.DRRegion{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&reg.RegionID, &reg.Name, &reg.Endpoint, &reg.IsActive, &reg.Health, &reg.LastChecked, &reg.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("region not found: %s", name)
		}
		return nil, fmt.Errorf("query region: %w", err)
	}
	return reg, nil
}

func (r *postgresRepo) UpdateRegionHealth(ctx context.Context, name string, health domain.HealthStatus) error {
	query := `UPDATE dr_regions SET health = $1, last_checked = NOW() WHERE name = $2`
	_, err := r.db.ExecContext(ctx, query, health, name)
	if err != nil {
		return fmt.Errorf("update region health: %w", err)
	}
	return nil
}

func (r *postgresRepo) UpdateRegionActive(ctx context.Context, name string, active bool) error {
	query := `UPDATE dr_regions SET is_active = $1 WHERE name = $2`
	_, err := r.db.ExecContext(ctx, query, active, name)
	if err != nil {
		return fmt.Errorf("update region active: %w", err)
	}
	return nil
}

func (r *postgresRepo) InsertFailoverPlan(ctx context.Context, plan *domain.FailoverPlan) error {
	query := `INSERT INTO failover_plans (plan_id, name, source_region, target_region, is_automated, created_at, is_executed)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query,
		plan.PlanID, plan.Name, plan.SourceRegion, plan.TargetRegion, plan.IsAutomated, plan.CreatedAt, plan.IsExecuted,
	)
	if err != nil {
		return fmt.Errorf("insert failover plan: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindAllFailoverPlans(ctx context.Context) ([]domain.FailoverPlan, error) {
	query := `SELECT plan_id, name, source_region, target_region, is_automated, created_at, is_executed FROM failover_plans ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failover plans: %w", err)
	}
	defer rows.Close()
	var plans []domain.FailoverPlan
	for rows.Next() {
		var p domain.FailoverPlan
		if err := rows.Scan(&p.PlanID, &p.Name, &p.SourceRegion, &p.TargetRegion, &p.IsAutomated, &p.CreatedAt, &p.IsExecuted); err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	return plans, rows.Err()
}

func (r *postgresRepo) FindFailoverPlanByID(ctx context.Context, planID uuid.UUID) (*domain.FailoverPlan, error) {
	query := `SELECT plan_id, name, source_region, target_region, is_automated, created_at, is_executed FROM failover_plans WHERE plan_id = $1`
	p := &domain.FailoverPlan{}
	err := r.db.QueryRowContext(ctx, query, planID).Scan(
		&p.PlanID, &p.Name, &p.SourceRegion, &p.TargetRegion, &p.IsAutomated, &p.CreatedAt, &p.IsExecuted,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("plan not found: %s", planID)
		}
		return nil, fmt.Errorf("query failover plan: %w", err)
	}
	return p, nil
}

func (r *postgresRepo) UpdateFailoverPlanExecuted(ctx context.Context, planID uuid.UUID) error {
	query := `UPDATE failover_plans SET is_executed = TRUE WHERE plan_id = $1`
	_, err := r.db.ExecContext(ctx, query, planID)
	if err != nil {
		return fmt.Errorf("update failover plan executed: %w", err)
	}
	return nil
}

func (r *postgresRepo) InsertFailoverExecution(ctx context.Context, exec *domain.FailoverExecution) error {
	query := `INSERT INTO failover_executions (execution_id, plan_id, started_at, completed_at, is_successful, error_message)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query,
		exec.ExecutionID, exec.PlanID, exec.StartedAt, exec.CompletedAt, exec.IsSuccessful, exec.ErrorMessage,
	)
	if err != nil {
		return fmt.Errorf("insert failover execution: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindExecutionsByPlan(ctx context.Context, planID uuid.UUID) ([]domain.FailoverExecution, error) {
	query := `SELECT execution_id, plan_id, started_at, completed_at, is_successful, error_message
		FROM failover_executions WHERE plan_id = $1 ORDER BY started_at DESC`
	rows, err := r.db.QueryContext(ctx, query, planID)
	if err != nil {
		return nil, fmt.Errorf("query executions: %w", err)
	}
	defer rows.Close()
	var execs []domain.FailoverExecution
	for rows.Next() {
		var e domain.FailoverExecution
		if err := rows.Scan(&e.ExecutionID, &e.PlanID, &e.StartedAt, &e.CompletedAt, &e.IsSuccessful, &e.ErrorMessage); err != nil {
			return nil, err
		}
		execs = append(execs, e)
	}
	return execs, rows.Err()
}

func (r *postgresRepo) InsertBackupManifest(ctx context.Context, manifest *domain.BackupManifest) error {
	query := `INSERT INTO backup_manifests (manifest_id, region, backup_path, backup_size_mb, started_at, completed_at, is_valid)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query,
		manifest.ManifestID, manifest.Region, manifest.BackupPath, manifest.BackupSizeMB, manifest.StartedAt, manifest.CompletedAt, manifest.IsValid,
	)
	if err != nil {
		return fmt.Errorf("insert backup manifest: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindAllBackupManifests(ctx context.Context) ([]domain.BackupManifest, error) {
	query := `SELECT manifest_id, region, backup_path, backup_size_mb, started_at, completed_at, is_valid FROM backup_manifests ORDER BY completed_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query backup manifests: %w", err)
	}
	defer rows.Close()
	var manifests []domain.BackupManifest
	for rows.Next() {
		var m domain.BackupManifest
		if err := rows.Scan(&m.ManifestID, &m.Region, &m.BackupPath, &m.BackupSizeMB, &m.StartedAt, &m.CompletedAt, &m.IsValid); err != nil {
			return nil, err
		}
		manifests = append(manifests, m)
	}
	return manifests, rows.Err()
}

func (r *postgresRepo) FindBackupManifestByID(ctx context.Context, manifestID uuid.UUID) (*domain.BackupManifest, error) {
	query := `SELECT manifest_id, region, backup_path, backup_size_mb, started_at, completed_at, is_valid FROM backup_manifests WHERE manifest_id = $1`
	m := &domain.BackupManifest{}
	err := r.db.QueryRowContext(ctx, query, manifestID).Scan(
		&m.ManifestID, &m.Region, &m.BackupPath, &m.BackupSizeMB, &m.StartedAt, &m.CompletedAt, &m.IsValid,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("manifest not found: %s", manifestID)
		}
		return nil, fmt.Errorf("query backup manifest: %w", err)
	}
	return m, nil
}

func (r *postgresRepo) InsertRecoveryPoint(ctx context.Context, point *domain.RecoveryPoint) error {
	query := `INSERT INTO recovery_points (point_id, manifest_id, recovery_time, is_restored, restored_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query,
		point.PointID, point.ManifestID, point.RecoveryTime, point.IsRestored, point.RestoredAt,
	)
	if err != nil {
		return fmt.Errorf("insert recovery point: %w", err)
	}
	return nil
}

func (r *postgresRepo) InsertDRTestResult(ctx context.Context, result *domain.DRTestResult) error {
	query := `INSERT INTO dr_test_results (test_id, plan_id, test_name, started_at, completed_at, is_successful, details)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query,
		result.TestID, result.PlanID, result.TestName, result.StartedAt, result.CompletedAt, result.IsSuccessful, result.Details,
	)
	if err != nil {
		return fmt.Errorf("insert dr test result: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindAllDRTestResults(ctx context.Context) ([]domain.DRTestResult, error) {
	query := `SELECT test_id, plan_id, test_name, started_at, completed_at, is_successful, details FROM dr_test_results ORDER BY completed_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query dr test results: %w", err)
	}
	defer rows.Close()
	var results []domain.DRTestResult
	for rows.Next() {
		var res domain.DRTestResult
		if err := rows.Scan(&res.TestID, &res.PlanID, &res.TestName, &res.StartedAt, &res.CompletedAt, &res.IsSuccessful, &res.Details); err != nil {
			return nil, err
		}
		results = append(results, res)
	}
	return results, rows.Err()
}

func (r *postgresRepo) InsertReplicationStatus(ctx context.Context, status *domain.ReplicationStatus) error {
	query := `INSERT INTO replication_status (replication_id, source_region, target_region, lag_seconds, is_healthy, last_checked_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query,
		status.ReplicationID, status.SourceRegion, status.TargetRegion, status.LagSeconds, status.IsHealthy, status.LastCheckedAt,
	)
	if err != nil {
		return fmt.Errorf("insert replication status: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindReplicationStatus(ctx context.Context) ([]domain.ReplicationStatus, error) {
	query := `SELECT replication_id, source_region, target_region, lag_seconds, is_healthy, last_checked_at
		FROM replication_status ORDER BY last_checked_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query replication status: %w", err)
	}
	defer rows.Close()
	var statuses []domain.ReplicationStatus
	for rows.Next() {
		var s domain.ReplicationStatus
		if err := rows.Scan(&s.ReplicationID, &s.SourceRegion, &s.TargetRegion, &s.LagSeconds, &s.IsHealthy, &s.LastCheckedAt); err != nil {
			return nil, err
		}
		statuses = append(statuses, s)
	}
	return statuses, rows.Err()
}
