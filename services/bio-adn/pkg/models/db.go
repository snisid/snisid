package models

import "context"

type Database interface {
	CreateDNAProfile(ctx context.Context, p *DNAProfile) error
	GetDNAProfileByHash(ctx context.Context, hash string) (*DNAProfile, error)
	GetDNAProfileBySpecimen(ctx context.Context, specimen string) (*DNAProfile, error)
	SearchDNAProfiles(ctx context.Context, indexType string, limit, offset int) ([]DNAProfile, int, error)
	GetUnuploadedDNAProfiles(ctx context.Context, level string) ([]map[string]any, error)
	MarkUploaded(ctx context.Context, id, level string) error
	MarkExpunged(ctx context.Context, id string) error
	CreateWantedPerson(ctx context.Context, p *WantedPerson) error
	QueryWantedPersons(ctx context.Context, q *WantedQuery) ([]WantedPerson, int, error)
	GetWantedByID(ctx context.Context, id string) (*WantedPerson, error)
	UpdateWantedStatus(ctx context.Context, id, status string) error
	QueryPlateIndex(ctx context.Context, plate string) (*PlateHitResult, error)
	QueryVINIndex(ctx context.Context, vin string) (*PlateHitResult, error)
	QueryPlateClones(ctx context.Context, plate string) (int, error)
	CreateStolenVehicle(ctx context.Context, v *StolenVehicle) error
	UpdateVehicleStatus(ctx context.Context, id, status, location, agency string) error
	CreateStolenFirearm(ctx context.Context, f *StolenFirearm) error
	CreateStolenDocument(ctx context.Context, d *StolenDocument) error
	CreateStolenVessel(ctx context.Context, v *StolenVessel) error
	CreateStolenArticle(ctx context.Context, a *StolenArticle) error
	CreateStolenSecurity(ctx context.Context, s *StolenSecurity) error
	CreateForeignFugitive(ctx context.Context, f *ForeignFugitive) error
	QueryForeignFugitives(ctx context.Context, lastName, nationality, noticeType string, limit, offset int) ([]ForeignFugitive, int, error)
	GetForeignFugitiveByID(ctx context.Context, id string) (*ForeignFugitive, error)
	CreateUnidentifiedPerson(ctx context.Context, u *UnidentifiedPerson) error
	QueryUnidentifiedPersons(ctx context.Context, dept, gender string, ageMin, ageMax int, limit, offset int) ([]UnidentifiedPerson, int, error)
	GetUnidentifiedByID(ctx context.Context, id string) (*UnidentifiedPerson, error)
	CreateTerrorismWatch(ctx context.Context, t *TerrorismWatch) error
	QueryTerrorismWatches(ctx context.Context, riskLevel, threatType, nationality string, limit, offset int) ([]TerrorismWatch, int, error)
	GetTerrorismWatchByID(ctx context.Context, id string) (*TerrorismWatch, error)
	CreateProtectionOrder(ctx context.Context, p *ProtectionOrder) error
	QueryProtectionOrders(ctx context.Context, beneficiaryName, restrainedPerson, orderType string, limit, offset int) ([]ProtectionOrder, int, error)
	GetActiveProtectionOrdersByBeneficiary(ctx context.Context, beneficiaryNIU string) ([]ProtectionOrder, error)
	CreateSupervisedRelease(ctx context.Context, s *SupervisedRelease) error
	QuerySupervisedReleases(ctx context.Context, niu, supervisionType, status string, limit, offset int) ([]SupervisedRelease, int, error)
	GetSupervisedReleaseByID(ctx context.Context, id string) (*SupervisedRelease, error)
	UpdateSexOffenderRisk(ctx context.Context, id, riskLevel, address string) error
	RecordGangMemberReview(ctx context.Context, id string) error
	CreateLabEquipment(ctx context.Context, e *LabEquipment) error
	QueryLabEquipment(ctx context.Context, labCode string) ([]LabEquipment, error)
	GetLabEquipmentByID(ctx context.Context, id string) (*LabEquipment, error)
	UpdateEquipmentCalibration(ctx context.Context, id, calibrationDate, calibrationDue, status string) error
	CreateStaffTraining(ctx context.Context, t *StaffTraining) error
	QueryStaffTraining(ctx context.Context, staffNIU string) ([]StaffTraining, error)
	GetStaffTrainingByID(ctx context.Context, id string) (*StaffTraining, error)
	CheckDuplicateSpecimen(ctx context.Context, specimen string) (bool, error)
	MarkSpecimenSubmitted(ctx context.Context, specimen, sampleID string) error
	RecordCrossDeptHit(ctx context.Context, h *NdisCrossDeptHit) error
	QueryCrossDeptHits(ctx context.Context, sdis, matchType string, limit, offset int) ([]NdisCrossDeptHit, int, error)
	GetNdisStats(ctx context.Context) (*NdisStats, error)
	CreateNdisReport(ctx context.Context, r *NdisReport) error
	QueryNdisReports(ctx context.Context) ([]NdisReport, error)
	CreateInterpolSubmission(ctx context.Context, s *InterpolSubmission) error
	CreateIdentityLink(ctx context.Context, l *BioIdentityLink) error
	GetIdentityLinkBySampleID(ctx context.Context, sampleID string) (*BioIdentityLink, error)
	QueryIdentityLinksByNIU(ctx context.Context, niu string) ([]BioIdentityLink, error)
	CountInterpolSubmissionsThisWeek(ctx context.Context) (int, error)
	CreateViolenceRecord(ctx context.Context, v *ViolenceRecord) error
	QueryViolenceRecords(ctx context.Context, niu, incidentType, status string, limit, offset int) ([]ViolenceRecord, int, error)
	GetViolenceRecordByID(ctx context.Context, id string) (*ViolenceRecord, error)
	CreateIdentityTheft(ctx context.Context, i *IdentityTheft) error
	QueryIdentityThefts(ctx context.Context, victimNIU, fraudType, status string, limit, offset int) ([]IdentityTheft, int, error)
	GetIdentityTheftByID(ctx context.Context, id string) (*IdentityTheft, error)
}
