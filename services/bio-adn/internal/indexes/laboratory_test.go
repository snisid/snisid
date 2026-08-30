package indexes

import (
	"context"
	"testing"

	"github.com/snisid/platform/services/bio-adn/pkg/models"
)

type mockLabOpsDB struct {
	models.Database
	equipment []models.LabEquipment
	training  []models.StaffTraining
}

func newMockLabOpsDB() *mockLabOpsDB {
	return &mockLabOpsDB{
		equipment: make([]models.LabEquipment, 0),
		training:  make([]models.StaffTraining, 0),
	}
}

func (m *mockLabOpsDB) CreateLabEquipment(ctx context.Context, e *models.LabEquipment) error {
	m.equipment = append(m.equipment, *e)
	return nil
}

func (m *mockLabOpsDB) CreateStaffTraining(ctx context.Context, t *models.StaffTraining) error {
	m.training = append(m.training, *t)
	return nil
}

func TestLabOps_RegisterEquipment(t *testing.T) {
	db := newMockLabOpsDB()
	idx := NewPersonsIndex(db)
	err := idx.CreateLabEquipment(context.Background(), &models.LabEquipment{
		ID:            "EQ-001",
		LabCode:       "LDIS-PAP-001",
		EquipmentName: "ABI 3500",
		SerialNumber:  "ABI-2026-001",
		Role:          "Analyse STR",
		Status:        "ACTIVE",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLabOps_RecordTraining(t *testing.T) {
	db := newMockLabOpsDB()
	idx := NewPersonsIndex(db)
	err := idx.CreateStaffTraining(context.Background(), &models.StaffTraining{
		ID:            "TR-001",
		StaffNIU:      "HTI-12345",
		TrainingName:  "STR Analysis",
		TrainingCode:  "STR-40H",
		DurationHours: 40,
		CompletedDate: "2026-06-01",
		IssuedBy:      "Direction SNISID",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
