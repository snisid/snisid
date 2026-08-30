package fixtures

import (
	"time"

	"github.com/google/uuid"
	"github.com/snisid/vehicle-criminal-svc/internal/domain"
)

func CreateTestAlert() *domain.CriminalAlert {
	return &domain.CriminalAlert{
		AlertID:       uuid.New(),
		PlateNumber:   "PP-1234",
		PlateCategory: domain.PlateCategoryPrivate,
		VehicleType:   domain.VehicleTypeBerline,
		Make:          "Toyota",
		Model:         "Corolla",
		Year:          int16Ptr(2020),
		ColorPrimary:  "Blanc",
		CrimeCategory: domain.CrimeCategoryVehicleTheft,
		AlertLevel:    domain.AlertLevelWanted,
		Status:        domain.AlertStatusActive,
		ReportingUnit: "BLVV",
		IncidentDate:  time.Now().Add(-24 * time.Hour),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Version:       1,
	}
}

func CreateTestCriticalAlert() *domain.CriminalAlert {
	notes := "Occupants armés, faux policiers. Appeler renforts avant intervention."
	return &domain.CriminalAlert{
		AlertID:            uuid.New(),
		PlateNumber:        "SE-00871",
		PlateCategory:      domain.PlateCategoryState,
		VehicleType:        domain.VehicleTypeSUV,
		Make:               "Toyota",
		Model:              "Land Cruiser",
		Year:               int16Ptr(2023),
		ColorPrimary:       "Blanc",
		CrimeCategory:      domain.CrimeCategoryFakePolice,
		AlertLevel:         domain.AlertLevelCritical,
		Status:             domain.AlertStatusActive,
		ArmedAndDangerous:  true,
		DoNotStopAlone:     true,
		OfficerSafetyNotes: &notes,
		ReportingUnit:      "DCPJ",
		IncidentDate:       time.Now().Add(-2 * time.Hour),
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		Version:            1,
	}
}

func CreateTestStolenPlate() *domain.StolenPlate {
	return &domain.StolenPlate{
		PlateID:       uuid.New(),
		PlateNumber:   "PP-5678",
		PlateCategory: domain.PlateCategoryPrivate,
		TheftDate:     time.Now().Add(-48 * time.Hour),
		TheftDeptCode: "OU",
		TheftLocation: strPtr("Pétion-Ville, Delmas 33"),
		Status:        domain.StolenPlateStatusStolen,
		ReportingUnit: "BLVV",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func CreateTestKidnappingAlert() *domain.CriminalAlert {
	notes := "Véhicule utilisé pour enlèvement. Ne pas intercepter seul."
	return &domain.CriminalAlert{
		AlertID:            uuid.New(),
		PlateNumber:        "M-9999",
		PlateCategory:      domain.PlateCategoryMoto,
		VehicleType:        domain.VehicleTypeMoto,
		Make:               "Honda",
		Model:              "CRF250",
		ColorPrimary:       "Rouge",
		CrimeCategory:      domain.CrimeCategoryKidnapping,
		AlertLevel:         domain.AlertLevelCritical,
		Status:             domain.AlertStatusActive,
		ArmedAndDangerous:  true,
		DoNotStopAlone:     true,
		OfficerSafetyNotes: &notes,
		ReportingUnit:      "CAE",
		IncidentDate:       time.Now().Add(-1 * time.Hour),
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		Version:            1,
	}
}

func int16Ptr(v int16) *int16 {
	return &v
}

func strPtr(s string) *string {
	return &s
}
