package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/fips-cert-svc/internal/domain"
)

type Repository interface {
	CreateModule(ctx context.Context, m *domain.CryptoModule) error
	FindModuleByID(ctx context.Context, id uuid.UUID) (*domain.CryptoModule, error)
	ListModules(ctx context.Context) ([]domain.CryptoModule, error)
	UpdateValidation(ctx context.Context, id uuid.UUID, status domain.ValidationStatus, certNumber string, validationDate time.Time) error
	CreateCVEResult(ctx context.Context, r *domain.CVEScanResult) error
	ListCVEsByModule(ctx context.Context, moduleID uuid.UUID) ([]domain.CVEScanResult, error)
	GetComplianceByService(ctx context.Context, service string) (*domain.ComplianceReport, error)
	GetDashboard(ctx context.Context) ([]domain.ComplianceReport, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateModule(ctx context.Context, m *domain.CryptoModule) error {
	q := `INSERT INTO fips_crypto_modules (module_id, name, version, vendor, fips_level, algorithms, cert_number, validation_date, expiry_date, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.db.ExecContext(ctx, q,
		m.ModuleID, m.Name, m.Version, m.Vendor, m.FIPSLevel,
		algoSlice(m.Algorithms), m.CertNumber, m.ValidationDate, m.ExpiryDate,
		m.Status, time.Now().UTC(), time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert module: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindModuleByID(ctx context.Context, id uuid.UUID) (*domain.CryptoModule, error) {
	q := `SELECT module_id, name, version, vendor, fips_level, algorithms, cert_number, validation_date, expiry_date, status, created_at, updated_at
		FROM fips_crypto_modules WHERE module_id = $1`
	m := &domain.CryptoModule{}
	var algos []string
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&m.ModuleID, &m.Name, &m.Version, &m.Vendor, &m.FIPSLevel,
		&algos, &m.CertNumber, &m.ValidationDate, &m.ExpiryDate,
		&m.Status, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("module not found: %s", id)
		}
		return nil, fmt.Errorf("query module: %w", err)
	}
	m.Algorithms = parseAlgos(algos)
	return m, nil
}

func (r *postgresRepo) ListModules(ctx context.Context) ([]domain.CryptoModule, error) {
	q := `SELECT module_id, name, version, vendor, fips_level, algorithms, cert_number, validation_date, expiry_date, status, created_at, updated_at
		FROM fips_crypto_modules ORDER BY name`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list modules: %w", err)
	}
	defer rows.Close()

	var modules []domain.CryptoModule
	for rows.Next() {
		var m domain.CryptoModule
		var algos []string
		if err := rows.Scan(
			&m.ModuleID, &m.Name, &m.Version, &m.Vendor, &m.FIPSLevel,
			&algos, &m.CertNumber, &m.ValidationDate, &m.ExpiryDate,
			&m.Status, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		m.Algorithms = parseAlgos(algos)
		modules = append(modules, m)
	}
	return modules, rows.Err()
}

func (r *postgresRepo) UpdateValidation(ctx context.Context, id uuid.UUID, status domain.ValidationStatus, certNumber string, validationDate time.Time) error {
	q := `UPDATE fips_crypto_modules SET status = $1, cert_number = $2, validation_date = $3, updated_at = $4 WHERE module_id = $5`
	_, err := r.db.ExecContext(ctx, q, status, certNumber, validationDate, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("update validation: %w", err)
	}
	return nil
}

func (r *postgresRepo) CreateCVEResult(ctx context.Context, cve *domain.CVEScanResult) error {
	q := `INSERT INTO fips_cve_results (scan_id, module_id, cve_id, severity, discovered, patched, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, q,
		cve.ScanID, cve.ModuleID, cve.CVEID, cve.Severity, cve.Discovered, cve.Patched, cve.Notes,
	)
	if err != nil {
		return fmt.Errorf("insert cve: %w", err)
	}
	return nil
}

func (r *postgresRepo) ListCVEsByModule(ctx context.Context, moduleID uuid.UUID) ([]domain.CVEScanResult, error) {
	q := `SELECT scan_id, module_id, cve_id, severity, discovered, patched, notes FROM fips_cve_results WHERE module_id = $1 ORDER BY discovered DESC`
	rows, err := r.db.QueryContext(ctx, q, moduleID)
	if err != nil {
		return nil, fmt.Errorf("list cves: %w", err)
	}
	defer rows.Close()

	var results []domain.CVEScanResult
	for rows.Next() {
		var r domain.CVEScanResult
		if err := rows.Scan(&r.ScanID, &r.ModuleID, &r.CVEID, &r.Severity, &r.Discovered, &r.Patched, &r.Notes); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

func (r *postgresRepo) GetComplianceByService(ctx context.Context, service string) (*domain.ComplianceReport, error) {
	q := `SELECT
		$1::text AS service_name,
		COUNT(*) AS module_count,
		COUNT(*) FILTER (WHERE status = 'VALIDATED') AS validated_count,
		COUNT(*) FILTER (WHERE status = 'PENDING') AS pending_count,
		COUNT(*) FILTER (WHERE status = 'EXPIRED') AS expired_count
		FROM fips_crypto_modules WHERE vendor = $1 OR name ILIKE '%' || $1 || '%'`
	rep := &domain.ComplianceReport{}
	err := r.db.QueryRowContext(ctx, q, service).Scan(
		&rep.ServiceName, &rep.ModuleCount, &rep.ValidatedCount, &rep.PendingCount, &rep.ExpiredCount,
	)
	if err != nil {
		return nil, fmt.Errorf("compliance query: %w", err)
	}
	rep.LastChecked = time.Now().UTC()

	cveQ := `SELECT COUNT(*) FROM fips_cve_results cv JOIN fips_crypto_modules m ON cv.module_id = m.module_id WHERE (m.vendor = $1 OR m.name ILIKE '%' || $1 || '%') AND (cv.patched IS NULL OR cv.patched = false)`
	r.db.QueryRowContext(ctx, cveQ, service).Scan(&rep.OpenCVEs)

	if rep.ModuleCount == 0 {
		rep.OverallStatus = "NO_MODULES"
	} else if rep.ValidatedCount == rep.ModuleCount && rep.OpenCVEs == 0 {
		rep.OverallStatus = "COMPLIANT"
	} else if rep.ExpiredCount > 0 || rep.OpenCVEs > 0 {
		rep.OverallStatus = "NON_COMPLIANT"
	} else {
		rep.OverallStatus = "PARTIALLY_COMPLIANT"
	}
	return rep, nil
}

func (r *postgresRepo) GetDashboard(ctx context.Context) ([]domain.ComplianceReport, error) {
	q := `SELECT COALESCE(vendor, 'unknown') AS service_name,
		COUNT(*) AS module_count,
		COUNT(*) FILTER (WHERE status = 'VALIDATED') AS validated_count,
		COUNT(*) FILTER (WHERE status = 'PENDING') AS pending_count,
		COUNT(*) FILTER (WHERE status = 'EXPIRED') AS expired_count
		FROM fips_crypto_modules GROUP BY vendor ORDER BY vendor`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("dashboard query: %w", err)
	}
	defer rows.Close()

	var reports []domain.ComplianceReport
	for rows.Next() {
		var rep domain.ComplianceReport
		if err := rows.Scan(&rep.ServiceName, &rep.ModuleCount, &rep.ValidatedCount, &rep.PendingCount, &rep.ExpiredCount); err != nil {
			return nil, err
		}
		rep.LastChecked = time.Now().UTC()
		reports = append(reports, rep)
	}
	return reports, rows.Err()
}

func algoSlice(algos []domain.CertAlgo) []string {
	s := make([]string, len(algos))
	for i, a := range algos {
		s[i] = string(a)
	}
	return s
}

func parseAlgos(s []string) []domain.CertAlgo {
	a := make([]domain.CertAlgo, len(s))
	for i, v := range s {
		a[i] = domain.CertAlgo(v)
	}
	return a
}
