package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/biar-svc/internal/domain"
)

type weaponRepo struct {
	pool *pgxpool.Pool
}

func NewWeaponRepo(pool *pgxpool.Pool) *weaponRepo {
	return &weaponRepo{pool: pool}
}

const weaponColumns = `weapon_id, national_biar_id, serial_number, serial_obliterated, make, model, caliber,
	weapon_type, manufacture_country, estimated_manufacture_year, recovery_date,
	recovery_context, recovery_location, recovery_dept_code, recovery_commune,
	recovery_lat, recovery_lng, seizing_unit, seizing_officer, case_reference,
	from_person_id, gang_id, crime_category, associated_cases, origin_country,
	transit_countries, trafficking_route, import_method, iarms_ref, atf_etrace_ref,
	reported_to_interpol, interpol_reported_at, disposition, disposal_date, disposal_auth,
	quantity_ammunition, ammunition_type, photos_refs, notes, created_by, created_at, updated_at`

func (r *weaponRepo) CreateWeapon(w *domain.IllicitWeapon) (*domain.IllicitWeapon, error) {
	ctx := context.Background()
	w.WeaponID = uuid.New()
	w.CreatedAt = time.Now().UTC()
	w.UpdatedAt = time.Now().UTC()

	if w.AssociatedCases == nil {
		w.AssociatedCases = []string{}
	}
	if w.TransitCountries == nil {
		w.TransitCountries = []string{}
	}
	if w.PhotosRefs == nil {
		w.PhotosRefs = []string{}
	}

	err := r.pool.QueryRow(ctx,
		`INSERT INTO biar_illicit_weapons
		 (weapon_id, national_biar_id, serial_number, serial_obliterated, make, model, caliber,
		  weapon_type, manufacture_country, estimated_manufacture_year, recovery_date,
		  recovery_context, recovery_location, recovery_dept_code, recovery_commune,
		  recovery_lat, recovery_lng, seizing_unit, seizing_officer, case_reference,
		  from_person_id, gang_id, crime_category, associated_cases, origin_country,
		  transit_countries, trafficking_route, import_method, iarms_ref, atf_etrace_ref,
		  reported_to_interpol, interpol_reported_at, disposition, disposal_date, disposal_auth,
		  quantity_ammunition, ammunition_type, photos_refs, notes, created_by, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36,$37,$38,$39,$40,$41,$42)
		 RETURNING created_at, updated_at`,
		w.WeaponID, w.NationalBIARID, w.SerialNumber, w.SerialObliterated, w.Make, w.Model,
		w.Caliber, w.WeaponType, w.ManufactureCountry, w.EstimatedManufactureYear,
		w.RecoveryDate, w.RecoveryContext, w.RecoveryLocation, w.RecoveryDeptCode,
		w.RecoveryCommune, w.RecoveryLat, w.RecoveryLng, w.SeizingUnit, w.SeizingOfficer,
		w.CaseReference, w.FromPersonID, w.GangID, w.CrimeCategory, w.AssociatedCases,
		w.OriginCountry, w.TransitCountries, w.TraffickingRoute, w.ImportMethod,
		w.IARMSRef, w.ATFEtraceRef, w.ReportedToInterpol, w.InterpolReportedAt,
		w.Disposition, w.DisposalDate, w.DisposalAuth,
		w.QuantityAmmunition, w.AmmunitionType, w.PhotosRefs, w.Notes, w.CreatedBy,
		w.CreatedAt, w.UpdatedAt,
	).Scan(&w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (r *weaponRepo) FindByID(id uuid.UUID) (*domain.IllicitWeapon, error) {
	ctx := context.Background()
	w := &domain.IllicitWeapon{}
	err := r.pool.QueryRow(ctx,
		`SELECT `+weaponColumns+` FROM biar_illicit_weapons WHERE weapon_id = $1`, id).Scan(
		&w.WeaponID, &w.NationalBIARID, &w.SerialNumber, &w.SerialObliterated, &w.Make,
		&w.Model, &w.Caliber, &w.WeaponType, &w.ManufactureCountry,
		&w.EstimatedManufactureYear, &w.RecoveryDate, &w.RecoveryContext,
		&w.RecoveryLocation, &w.RecoveryDeptCode, &w.RecoveryCommune, &w.RecoveryLat,
		&w.RecoveryLng, &w.SeizingUnit, &w.SeizingOfficer, &w.CaseReference,
		&w.FromPersonID, &w.GangID, &w.CrimeCategory, &w.AssociatedCases,
		&w.OriginCountry, &w.TransitCountries, &w.TraffickingRoute, &w.ImportMethod,
		&w.IARMSRef, &w.ATFEtraceRef, &w.ReportedToInterpol, &w.InterpolReportedAt,
		&w.Disposition, &w.DisposalDate, &w.DisposalAuth,
		&w.QuantityAmmunition, &w.AmmunitionType, &w.PhotosRefs, &w.Notes,
		&w.CreatedBy, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (r *weaponRepo) FindBySerial(sn string) ([]domain.IllicitWeapon, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT `+weaponColumns+` FROM biar_illicit_weapons WHERE serial_number = $1`, sn)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanWeapons(rows)
}

func (r *weaponRepo) UpdateWeapon(w *domain.IllicitWeapon) (*domain.IllicitWeapon, error) {
	ctx := context.Background()
	w.UpdatedAt = time.Now().UTC()

	_, err := r.pool.Exec(ctx,
		`UPDATE biar_illicit_weapons SET
		 serial_number=$1, serial_obliterated=$2, make=$3, model=$4, caliber=$5,
		 weapon_type=$6, manufacture_country=$7, estimated_manufacture_year=$8,
		 recovery_date=$9, recovery_context=$10, recovery_location=$11,
		 recovery_dept_code=$12, recovery_commune=$13, recovery_lat=$14, recovery_lng=$15,
		 seizing_unit=$16, seizing_officer=$17, case_reference=$18,
		 from_person_id=$19, gang_id=$20, crime_category=$21, associated_cases=$22,
		 origin_country=$23, transit_countries=$24, trafficking_route=$25,
		 import_method=$26, iarms_ref=$27, atf_etrace_ref=$28,
		 reported_to_interpol=$29, interpol_reported_at=$30, disposition=$31,
		 disposal_date=$32, disposal_auth=$33, quantity_ammunition=$34,
		 ammunition_type=$35, photos_refs=$36, notes=$37, updated_at=$38
		 WHERE weapon_id = $39`,
		w.SerialNumber, w.SerialObliterated, w.Make, w.Model, w.Caliber,
		w.WeaponType, w.ManufactureCountry, w.EstimatedManufactureYear,
		w.RecoveryDate, w.RecoveryContext, w.RecoveryLocation,
		w.RecoveryDeptCode, w.RecoveryCommune, w.RecoveryLat, w.RecoveryLng,
		w.SeizingUnit, w.SeizingOfficer, w.CaseReference,
		w.FromPersonID, w.GangID, w.CrimeCategory, w.AssociatedCases,
		w.OriginCountry, w.TransitCountries, w.TraffickingRoute,
		w.ImportMethod, w.IARMSRef, w.ATFEtraceRef,
		w.ReportedToInterpol, w.InterpolReportedAt, w.Disposition,
		w.DisposalDate, w.DisposalAuth, w.QuantityAmmunition,
		w.AmmunitionType, w.PhotosRefs, w.Notes, w.UpdatedAt, w.WeaponID)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (r *weaponRepo) CreateBatch(b *domain.BatchSeizure) (*domain.BatchSeizure, error) {
	ctx := context.Background()
	b.BatchID = uuid.New()
	b.CreatedAt = time.Now().UTC()

	if b.WeaponIDs == nil {
		b.WeaponIDs = []string{}
	}
	if b.PartneringAgencies == nil {
		b.PartneringAgencies = []string{}
	}

	err := r.pool.QueryRow(ctx,
		`INSERT INTO biar_batch_seizures
		 (batch_id, batch_reference, operation_name, seizure_date, location_desc, dept_code,
		  total_weapons, weapon_ids, seizing_unit, lead_officer, partnering_agencies, notes, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		 RETURNING created_at`,
		b.BatchID, b.BatchReference, b.OperationName, b.SeizureDate, b.LocationDesc,
		b.DeptCode, b.TotalWeapons, b.WeaponIDs, b.SeizingUnit, b.LeadOfficer,
		b.PartneringAgencies, b.Notes, b.CreatedAt,
	).Scan(&b.CreatedAt)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r *weaponRepo) GetWeaponsByGang(gangID uuid.UUID) ([]domain.IllicitWeapon, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT `+weaponColumns+` FROM biar_illicit_weapons WHERE gang_id = $1 ORDER BY recovery_date DESC`, gangID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanWeapons(rows)
}

func (r *weaponRepo) GetWeaponsByOrigin(origin string) ([]domain.IllicitWeapon, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT `+weaponColumns+` FROM biar_illicit_weapons WHERE origin_country = $1 ORDER BY recovery_date DESC`, origin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanWeapons(rows)
}

func (r *weaponRepo) GetStatsByGang() ([]map[string]interface{}, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT gang_id::text, COUNT(*) as count FROM biar_illicit_weapons
		 WHERE gang_id IS NOT NULL GROUP BY gang_id ORDER BY count DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanStats(rows)
}

func (r *weaponRepo) GetStatsByOrigin() ([]map[string]interface{}, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT origin_country, COUNT(*) as count FROM biar_illicit_weapons
		 WHERE origin_country IS NOT NULL GROUP BY origin_country ORDER BY count DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanStats(rows)
}

func (r *weaponRepo) GetRoutes() ([]map[string]interface{}, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT trafficking_route, COUNT(*) as count FROM biar_illicit_weapons
		 WHERE trafficking_route IS NOT NULL AND trafficking_route != ''
		 GROUP BY trafficking_route ORDER BY count DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanStats(rows)
}

func (r *weaponRepo) UpsertFromIARMS(w *domain.IllicitWeapon) (*domain.IllicitWeapon, error) {
	ctx := context.Background()

	if w.AssociatedCases == nil {
		w.AssociatedCases = []string{}
	}
	if w.TransitCountries == nil {
		w.TransitCountries = []string{}
	}
	if w.PhotosRefs == nil {
		w.PhotosRefs = []string{}
	}

	err := r.pool.QueryRow(ctx,
		`INSERT INTO biar_illicit_weapons
		 (weapon_id, national_biar_id, serial_number, serial_obliterated, make, model, caliber,
		  weapon_type, manufacture_country, estimated_manufacture_year, recovery_date,
		  recovery_context, recovery_location, recovery_dept_code, recovery_commune,
		  recovery_lat, recovery_lng, seizing_unit, seizing_officer, case_reference,
		  from_person_id, gang_id, crime_category, associated_cases, origin_country,
		  transit_countries, trafficking_route, import_method, iarms_ref, atf_etrace_ref,
		  reported_to_interpol, interpol_reported_at, disposition, disposal_date, disposal_auth,
		  quantity_ammunition, ammunition_type, photos_refs, notes, created_by, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36,$37,$38,$39,$40,$41,$42)
		 ON CONFLICT (national_biar_id) DO UPDATE SET
		  serial_number=EXCLUDED.serial_number, serial_obliterated=EXCLUDED.serial_obliterated,
		  make=EXCLUDED.make, model=EXCLUDED.model, caliber=EXCLUDED.caliber,
		  weapon_type=EXCLUDED.weapon_type, manufacture_country=EXCLUDED.manufacture_country,
		  estimated_manufacture_year=EXCLUDED.estimated_manufacture_year,
		  recovery_date=EXCLUDED.recovery_date, recovery_context=EXCLUDED.recovery_context,
		  recovery_location=EXCLUDED.recovery_location, recovery_dept_code=EXCLUDED.recovery_dept_code,
		  recovery_commune=EXCLUDED.recovery_commune, recovery_lat=EXCLUDED.recovery_lat,
		  recovery_lng=EXCLUDED.recovery_lng, seizing_unit=EXCLUDED.seizing_unit,
		  seizing_officer=EXCLUDED.seizing_officer, case_reference=EXCLUDED.case_reference,
		  origin_country=EXCLUDED.origin_country, transit_countries=EXCLUDED.transit_countries,
		  trafficking_route=EXCLUDED.trafficking_route, import_method=EXCLUDED.import_method,
		  iarms_ref=EXCLUDED.iarms_ref, atf_etrace_ref=EXCLUDED.atf_etrace_ref,
		  disposition=EXCLUDED.disposition, updated_at=NOW()
		 RETURNING created_at, updated_at`,
		w.WeaponID, w.NationalBIARID, w.SerialNumber, w.SerialObliterated, w.Make, w.Model,
		w.Caliber, w.WeaponType, w.ManufactureCountry, w.EstimatedManufactureYear,
		w.RecoveryDate, w.RecoveryContext, w.RecoveryLocation, w.RecoveryDeptCode,
		w.RecoveryCommune, w.RecoveryLat, w.RecoveryLng, w.SeizingUnit, w.SeizingOfficer,
		w.CaseReference, w.FromPersonID, w.GangID, w.CrimeCategory, w.AssociatedCases,
		w.OriginCountry, w.TransitCountries, w.TraffickingRoute, w.ImportMethod,
		w.IARMSRef, w.ATFEtraceRef, w.ReportedToInterpol, w.InterpolReportedAt,
		w.Disposition, w.DisposalDate, w.DisposalAuth,
		w.QuantityAmmunition, w.AmmunitionType, w.PhotosRefs, w.Notes, w.CreatedBy,
		w.CreatedAt, w.UpdatedAt,
	).Scan(&w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (r *weaponRepo) CreateSyncLog(log *domain.IARMSyncLog) error {
	ctx := context.Background()
	log.SyncID = uuid.New()
	log.CreatedAt = time.Now().UTC()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO biar_iarms_sync_log
		 (sync_id, weapon_id, direction, iarms_ref, sync_status, synced_at, error_message, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		log.SyncID, log.WeaponID, log.Direction, log.IARMSRef,
		log.SyncStatus, log.SyncedAt, log.ErrorMessage, log.CreatedAt)
	return err
}

func scanWeapons(rows pgx.Rows) ([]domain.IllicitWeapon, error) {
	var weapons []domain.IllicitWeapon
	for rows.Next() {
		var w domain.IllicitWeapon
		if err := rows.Scan(
			&w.WeaponID, &w.NationalBIARID, &w.SerialNumber, &w.SerialObliterated, &w.Make,
			&w.Model, &w.Caliber, &w.WeaponType, &w.ManufactureCountry,
			&w.EstimatedManufactureYear, &w.RecoveryDate, &w.RecoveryContext,
			&w.RecoveryLocation, &w.RecoveryDeptCode, &w.RecoveryCommune, &w.RecoveryLat,
			&w.RecoveryLng, &w.SeizingUnit, &w.SeizingOfficer, &w.CaseReference,
			&w.FromPersonID, &w.GangID, &w.CrimeCategory, &w.AssociatedCases,
			&w.OriginCountry, &w.TransitCountries, &w.TraffickingRoute, &w.ImportMethod,
			&w.IARMSRef, &w.ATFEtraceRef, &w.ReportedToInterpol, &w.InterpolReportedAt,
			&w.Disposition, &w.DisposalDate, &w.DisposalAuth,
			&w.QuantityAmmunition, &w.AmmunitionType, &w.PhotosRefs, &w.Notes,
			&w.CreatedBy, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		weapons = append(weapons, w)
	}
	return weapons, nil
}

func scanStats(rows pgx.Rows) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	for rows.Next() {
		var key string
		var count int
		if err := rows.Scan(&key, &count); err != nil {
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"label": key,
			"count": count,
		})
	}
	return results, nil
}
