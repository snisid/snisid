package indexes

import (
	"context"
	"fmt"
	"time"

	"github.com/snisid/platform/services/bio-adn/pkg/models"
)

type PersonsIndex struct {
	db models.Database
}

func NewPersonsIndex(db models.Database) *PersonsIndex {
	return &PersonsIndex{db: db}
}

func (idx *PersonsIndex) CreateWanted(ctx context.Context, p *models.WantedPerson) error {
	if p.WarrantType == "" {
		return fmt.Errorf("warrant_type is required")
	}
	if len(p.Charges) == 0 {
		return fmt.Errorf("at least one charge is required")
	}
	if (p.WarrantType == "MAN-ARR" || p.WarrantType == "MAN-EXT") && p.WarrantNumber == "" {
		return fmt.Errorf("warrant_number is required for MAN-ARR and MAN-EXT")
	}
	return idx.db.CreateWantedPerson(ctx, p)
}

func (idx *PersonsIndex) QueryWanted(ctx context.Context, q *models.WantedQuery) ([]models.WantedPerson, int, error) {
	if q.Limit <= 0 {
		q.Limit = 20
	}
	return idx.db.QueryWantedPersons(ctx, q)
}

func (idx *PersonsIndex) GetWanted(ctx context.Context, id string) (*models.WantedPerson, error) {
	return idx.db.GetWantedByID(ctx, id)
}

func (idx *PersonsIndex) UpdateStatus(ctx context.Context, id, status string) error {
	valid := map[string]bool{"ACTIVE": true, "CLEARED": true, "EXPIRED": true, "SUSPENDED": true}
	if !valid[status] {
		return fmt.Errorf("invalid status: %s", status)
	}
	return idx.db.UpdateWantedStatus(ctx, id, status)
}

func (idx *PersonsIndex) CreateForeignFugitive(ctx context.Context, f *models.ForeignFugitive) error {
	if f.InterpolNoticeNumber == "" {
		return fmt.Errorf("interpol_notice_number is required")
	}
	return idx.db.CreateForeignFugitive(ctx, f)
}

func (idx *PersonsIndex) QueryForeignFugitives(ctx context.Context, lastName, nationality, noticeType string, limit, offset int) ([]models.ForeignFugitive, int, error) {
	return idx.db.QueryForeignFugitives(ctx, lastName, nationality, noticeType, limit, offset)
}

func (idx *PersonsIndex) GetForeignFugitive(ctx context.Context, id string) (*models.ForeignFugitive, error) {
	return idx.db.GetForeignFugitiveByID(ctx, id)
}

func (idx *PersonsIndex) CreateUnidentifiedPerson(ctx context.Context, u *models.UnidentifiedPerson) error {
	return idx.db.CreateUnidentifiedPerson(ctx, u)
}

func (idx *PersonsIndex) QueryUnidentifiedPersons(ctx context.Context, dept, gender string, ageMin, ageMax, limit, offset int) ([]models.UnidentifiedPerson, int, error) {
	return idx.db.QueryUnidentifiedPersons(ctx, dept, gender, ageMin, ageMax, limit, offset)
}

func (idx *PersonsIndex) GetUnidentifiedPerson(ctx context.Context, id string) (*models.UnidentifiedPerson, error) {
	return idx.db.GetUnidentifiedByID(ctx, id)
}

func (idx *PersonsIndex) CreateTerrorismWatch(ctx context.Context, t *models.TerrorismWatch) error {
	if t.ApprovedByDirector == "" || t.ApprovedByPG == "" {
		return fmt.Errorf("dual approval from DCPJ director and AG required")
	}
	return idx.db.CreateTerrorismWatch(ctx, t)
}

func (idx *PersonsIndex) QueryTerrorismWatches(ctx context.Context, riskLevel, threatType, nationality string, limit, offset int) ([]models.TerrorismWatch, int, error) {
	return idx.db.QueryTerrorismWatches(ctx, riskLevel, threatType, nationality, limit, offset)
}

func (idx *PersonsIndex) GetTerrorismWatch(ctx context.Context, id string) (*models.TerrorismWatch, error) {
	return idx.db.GetTerrorismWatchByID(ctx, id)
}

func (idx *PersonsIndex) CreateProtectionOrder(ctx context.Context, po *models.ProtectionOrder) error {
	return idx.db.CreateProtectionOrder(ctx, po)
}

func (idx *PersonsIndex) QueryProtectionOrders(ctx context.Context, beneficiaryName, restrainedPerson, orderType string, limit, offset int) ([]models.ProtectionOrder, int, error) {
	return idx.db.QueryProtectionOrders(ctx, beneficiaryName, restrainedPerson, orderType, limit, offset)
}

func (idx *PersonsIndex) GetActiveProtectionOrders(ctx context.Context, beneficiaryNIU string) ([]models.ProtectionOrder, error) {
	return idx.db.GetActiveProtectionOrdersByBeneficiary(ctx, beneficiaryNIU)
}

func (idx *PersonsIndex) CreateSupervisedRelease(ctx context.Context, s *models.SupervisedRelease) error {
	return idx.db.CreateSupervisedRelease(ctx, s)
}

func (idx *PersonsIndex) QuerySupervisedReleases(ctx context.Context, niu, supervisionType, status string, limit, offset int) ([]models.SupervisedRelease, int, error) {
	return idx.db.QuerySupervisedReleases(ctx, niu, supervisionType, status, limit, offset)
}

func (idx *PersonsIndex) GetSupervisedRelease(ctx context.Context, id string) (*models.SupervisedRelease, error) {
	return idx.db.GetSupervisedReleaseByID(ctx, id)
}

func (idx *PersonsIndex) CreateLabEquipment(ctx context.Context, e *models.LabEquipment) error {
	return idx.db.CreateLabEquipment(ctx, e)
}

func (idx *PersonsIndex) QueryLabEquipment(ctx context.Context, labCode string) ([]models.LabEquipment, error) {
	return idx.db.QueryLabEquipment(ctx, labCode)
}

func (idx *PersonsIndex) GetLabEquipment(ctx context.Context, id string) (*models.LabEquipment, error) {
	return idx.db.GetLabEquipmentByID(ctx, id)
}

func (idx *PersonsIndex) UpdateEquipmentCalibration(ctx context.Context, id, calibrationDate, calibrationDue, status string) error {
	return idx.db.UpdateEquipmentCalibration(ctx, id, calibrationDate, calibrationDue, status)
}

func (idx *PersonsIndex) CreateStaffTraining(ctx context.Context, t *models.StaffTraining) error {
	return idx.db.CreateStaffTraining(ctx, t)
}

func (idx *PersonsIndex) QueryStaffTraining(ctx context.Context, staffNIU string) ([]models.StaffTraining, error) {
	return idx.db.QueryStaffTraining(ctx, staffNIU)
}

func (idx *PersonsIndex) GetStaffTraining(ctx context.Context, id string) (*models.StaffTraining, error) {
	return idx.db.GetStaffTrainingByID(ctx, id)
}

// ── PER-VIO: Known Violence ──────────────────────────────────────────────────

func (idx *PersonsIndex) CreateViolenceRecord(ctx context.Context, v *models.ViolenceRecord) error {
	return idx.db.CreateViolenceRecord(ctx, v)
}

func (idx *PersonsIndex) QueryViolenceRecords(ctx context.Context, niu, incidentType, status string, limit, offset int) ([]models.ViolenceRecord, int, error) {
	return idx.db.QueryViolenceRecords(ctx, niu, incidentType, status, limit, offset)
}

func (idx *PersonsIndex) GetViolenceRecord(ctx context.Context, id string) (*models.ViolenceRecord, error) {
	return idx.db.GetViolenceRecordByID(ctx, id)
}

// ── PER-IDV: Identity Theft ─────────────────────────────────────────────────

func (idx *PersonsIndex) CreateIdentityTheft(ctx context.Context, i *models.IdentityTheft) error {
	return idx.db.CreateIdentityTheft(ctx, i)
}

func (idx *PersonsIndex) QueryIdentityThefts(ctx context.Context, victimNIU, fraudType, status string, limit, offset int) ([]models.IdentityTheft, int, error) {
	return idx.db.QueryIdentityThefts(ctx, victimNIU, fraudType, status, limit, offset)
}

func (idx *PersonsIndex) GetIdentityTheft(ctx context.Context, id string) (*models.IdentityTheft, error) {
	return idx.db.GetIdentityTheftByID(ctx, id)
}

func timePtr(t time.Time) *string {
	s := t.Format(time.RFC3339)
	return &s
}
