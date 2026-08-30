package domain

import (
	"time"

	"github.com/google/uuid"
)

type ExplIncident struct {
	IncidentID          uuid.UUID  `json:"incident_id" db:"incident_id"`
	NationalExplID      string     `json:"national_expl_id" db:"national_expl_id"`
	IncidentType        string     `json:"incident_type" db:"incident_type"`
	ExplosiveType       ExplType   `json:"explosive_type" db:"explosive_type"`
	Status              ExplStatus `json:"status" db:"status"`
	Quantity            int        `json:"quantity" db:"quantity"`
	WeightKg            float64    `json:"weight_kg" db:"weight_kg"`
	Manufacturer        string     `json:"manufacturer" db:"manufacturer"`
	LotNumber           string     `json:"lot_number" db:"lot_number"`
	ManufactureCountry  string     `json:"manufacture_country" db:"manufacture_country"`
	EstimatedDate       *time.Time `json:"estimated_date" db:"estimated_date"`
	IncidentDate        time.Time  `json:"incident_date" db:"incident_date"`
	LocationDesc        string     `json:"location_desc" db:"location_desc"`
	DeptCode            string     `json:"dept_code" db:"dept_code"`
	Commune             string     `json:"commune" db:"commune"`
	Lat                 float64    `json:"lat" db:"lat"`
	Lng                 float64    `json:"lng" db:"lng"`
	RespondingUnit      string     `json:"responding_unit" db:"responding_unit"`
	EODOfficer          *uuid.UUID `json:"eod_officer" db:"eod_officer"`
	Casualties          int        `json:"casualties" db:"casualties"`
	GangID              *uuid.UUID `json:"gang_id" db:"gang_id"`
	FromPersonID        *uuid.UUID `json:"from_person_id" db:"from_person_id"`
	CaseReference       string     `json:"case_reference" db:"case_reference"`
	DNASampleTaken      bool       `json:"dna_sample_taken" db:"dna_sample_taken"`
	BioSampleRef        string     `json:"bio_sample_ref" db:"bio_sample_ref"`
	PhotoRefs           []string   `json:"photo_refs" db:"photo_refs"`
	InterpolExplointRef string     `json:"interpol_exploint_ref" db:"interpol_exploint_ref"`
	Notes               string     `json:"notes" db:"notes"`
	CreatedBy           uuid.UUID  `json:"created_by" db:"created_by"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
}

type LegalStock struct {
	StockID         uuid.UUID  `json:"stock_id" db:"stock_id"`
	HolderEntity    string     `json:"holder_entity" db:"holder_entity"`
	HolderType      string     `json:"holder_type" db:"holder_type"`
	ExplosiveType   ExplType   `json:"explosive_type" db:"explosive_type"`
	QuantityKg      float64    `json:"quantity_kg" db:"quantity_kg"`
	StorageLocation string     `json:"storage_location" db:"storage_location"`
	DeptCode        string     `json:"dept_code" db:"dept_code"`
	LicenseRef      string     `json:"license_ref" db:"license_ref"`
	LastAuditDate   *time.Time `json:"last_audit_date" db:"last_audit_date"`
	NextAuditDate   *time.Time `json:"next_audit_date" db:"next_audit_date"`
	IsSecured       bool       `json:"is_secured" db:"is_secured"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
}

type IncidentRepository interface {
	CreateIncident(incident *ExplIncident) error
	FindByID(id uuid.UUID) (*ExplIncident, error)
	FindByDept(deptCode string, limit, offset int) ([]ExplIncident, error)
	CreateLegalStock(stock *LegalStock) error
	GetLegalStocks(deptCode string, limit, offset int) ([]LegalStock, error)
	CountIncidents() (int, error)
}
