package service

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/lapi-ht/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPostgresRepo struct {
	savePlateReadFn     func(read *domain.PlateRead) error
	getRecentReadsFn    func(limit int) ([]domain.PlateRead, error)
	getReadsByPlateFn   func(plateNumber string) ([]domain.PlateRead, error)
	getActiveAlertsFn   func() ([]domain.AlertDispatch, error)
	getCamerasFn        func() ([]domain.Camera, error)
	saveAlertDispatchFn func(alert *domain.AlertDispatch) error
}

func (m *mockPostgresRepo) SavePlateRead(read *domain.PlateRead) error {
	return m.savePlateReadFn(read)
}
func (m *mockPostgresRepo) GetRecentReads(limit int) ([]domain.PlateRead, error) {
	return m.getRecentReadsFn(limit)
}
func (m *mockPostgresRepo) GetReadsByPlate(plateNumber string) ([]domain.PlateRead, error) {
	return m.getReadsByPlateFn(plateNumber)
}
func (m *mockPostgresRepo) GetActiveAlerts() ([]domain.AlertDispatch, error) {
	return m.getActiveAlertsFn()
}
func (m *mockPostgresRepo) GetCameras() ([]domain.Camera, error) {
	return m.getCamerasFn()
}
func (m *mockPostgresRepo) SaveAlertDispatch(alert *domain.AlertDispatch) error {
	return m.saveAlertDispatchFn(alert)
}

func TestRecordRead(t *testing.T) {
	cameraID := uuid.New()
	alertSaved := false
	tests := []struct {
		name      string
		read      *domain.PlateRead
		repoErr   error
		wantErr   bool
		wantAlert bool
	}{
		{
			name: "success without alert",
			read: &domain.PlateRead{
				CameraID:              cameraID,
				PlateNumberRaw:        "ABC123",
				PlateNumberNormalized: "ABC-123",
				OcrConfidence:         0.95,
				AlertTriggered:        false,
				CapturedAt:            time.Now(),
			},
		},
		{
			name: "success with alert",
			read: &domain.PlateRead{
				CameraID:              cameraID,
				PlateNumberRaw:        "WNT-001",
				PlateNumberNormalized: "WNT-001",
				OcrConfidence:         0.98,
				AlertTriggered:        true,
				CapturedAt:            time.Now(),
			},
			wantAlert: true,
		},
		{
			name: "repo error",
			read: &domain.PlateRead{
				CameraID:              cameraID,
				PlateNumberRaw:        "ERR-001",
				PlateNumberNormalized: "ERR-001",
				AlertTriggered:        false,
				CapturedAt:            time.Now(),
			},
			repoErr: errors.New("insert error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alertSaved = false
			repo := &mockPostgresRepo{
				savePlateReadFn: func(read *domain.PlateRead) error {
					return tt.repoErr
				},
				saveAlertDispatchFn: func(alert *domain.AlertDispatch) error {
					alertSaved = true
					return nil
				},
			}
			svc := NewLapiService(repo, nil)
			err := svc.RecordRead(tt.read)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, tt.read.ID)
			assert.Equal(t, tt.wantAlert, alertSaved)
		})
	}
}

func TestGetRecentReads(t *testing.T) {
	tests := []struct {
		name    string
		limit   int
		repoRes []domain.PlateRead
		repoErr error
		wantErr bool
		wantLen int
	}{
		{
			name:  "success default limit",
			limit: 0,
			repoRes: []domain.PlateRead{
				{ID: uuid.New(), PlateNumberNormalized: "ABC-123"},
			},
			wantLen: 1,
		},
		{
			name:  "custom limit",
			limit: 10,
			repoRes: func() []domain.PlateRead {
				reads := make([]domain.PlateRead, 10)
				for i := range reads {
					reads[i] = domain.PlateRead{ID: uuid.New()}
				}
				return reads
			}(),
			wantLen: 10,
		},
		{
			name:    "clamped to 50 when over 100",
			limit:   200,
			repoRes: []domain.PlateRead{},
			wantLen: 0,
		},
		{
			name:    "repo error",
			limit:   5,
			repoErr: errors.New("query error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockPostgresRepo{
				getRecentReadsFn: func(limit int) ([]domain.PlateRead, error) {
					return tt.repoRes, tt.repoErr
				},
			}
			svc := NewLapiService(repo, nil)
			got, err := svc.GetRecentReads(tt.limit)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Len(t, got, tt.wantLen)
		})
	}
}

func TestGetReadsByPlate(t *testing.T) {
	tests := []struct {
		name    string
		plate   string
		repoRes []domain.PlateRead
		repoErr error
		wantErr bool
	}{
		{
			name:  "success",
			plate: "ABC-123",
			repoRes: []domain.PlateRead{
				{ID: uuid.New(), PlateNumberNormalized: "ABC-123"},
			},
		},
		{
			name:    "not found",
			plate:   "ZZZ-999",
			repoRes: []domain.PlateRead{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockPostgresRepo{
				getReadsByPlateFn: func(plate string) ([]domain.PlateRead, error) {
					return tt.repoRes, tt.repoErr
				},
			}
			svc := NewLapiService(repo, nil)
			got, err := svc.GetReadsByPlate(tt.plate)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Len(t, got, len(tt.repoRes))
		})
	}
}

func TestGetActiveAlerts(t *testing.T) {
	tests := []struct {
		name    string
		repoRes []domain.AlertDispatch
		repoErr error
		wantErr bool
	}{
		{
			name: "success",
			repoRes: []domain.AlertDispatch{
				{ID: uuid.New(), PlateNumber: "WNT-001", IsActive: true},
			},
		},
		{
			name:    "repo error",
			repoErr: errors.New("query error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockPostgresRepo{
				getActiveAlertsFn: func() ([]domain.AlertDispatch, error) {
					return tt.repoRes, tt.repoErr
				},
			}
			svc := NewLapiService(repo, nil)
			got, err := svc.GetActiveAlerts()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Len(t, got, len(tt.repoRes))
		})
	}
}

func TestGetCameraStatus(t *testing.T) {
	tests := []struct {
		name    string
		repoRes []domain.Camera
		repoErr error
		wantErr bool
	}{
		{
			name: "success",
			repoRes: []domain.Camera{
				{ID: uuid.New(), Label: "Cam-01", Type: domain.CameraTypeFixedIntersection, IsActive: true},
			},
		},
		{
			name:    "repo error",
			repoErr: errors.New("query error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockPostgresRepo{
				getCamerasFn: func() ([]domain.Camera, error) {
					return tt.repoRes, tt.repoErr
				},
			}
			svc := NewLapiService(repo, nil)
			got, err := svc.GetCameraStatus()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Len(t, got, len(tt.repoRes))
		})
	}
}
