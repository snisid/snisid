package service

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/fpr-ht/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockFprRepo struct {
	saveWarrantFn           func(w *domain.Warrant) error
	findWarrantsByNameFn    func(name string) ([]domain.Warrant, error)
	saveCheckLogFn          func(cl *domain.CheckLog) error
	saveSightingFn          func(s *domain.Sighting) error
	updateWarrantExecutedFn func(id uuid.UUID, executedAt time.Time) error
	getArmedDangerousFn     func() ([]domain.Warrant, error)
	getDashboardStatsFn     func() (*domain.DashboardStats, error)
}

func (m *mockFprRepo) SaveWarrant(w *domain.Warrant) error {
	return m.saveWarrantFn(w)
}
func (m *mockFprRepo) FindWarrantsByName(name string) ([]domain.Warrant, error) {
	return m.findWarrantsByNameFn(name)
}
func (m *mockFprRepo) SaveCheckLog(cl *domain.CheckLog) error {
	return m.saveCheckLogFn(cl)
}
func (m *mockFprRepo) SaveSighting(s *domain.Sighting) error {
	return m.saveSightingFn(s)
}
func (m *mockFprRepo) UpdateWarrantExecuted(id uuid.UUID, executedAt time.Time) error {
	return m.updateWarrantExecutedFn(id, executedAt)
}
func (m *mockFprRepo) GetArmedDangerousWarrants() ([]domain.Warrant, error) {
	return m.getArmedDangerousFn()
}
func (m *mockFprRepo) GetDashboardStats() (*domain.DashboardStats, error) {
	return m.getDashboardStatsFn()
}

func TestCreateWarrant(t *testing.T) {
	tests := []struct {
		name    string
		warrant *domain.Warrant
		repoErr error
		wantErr bool
	}{
		{
			name: "success",
			warrant: &domain.Warrant{
				FullName:     "John Doe",
				WarrantType:  domain.WarrantTypeArrest,
				Charges:      []string{"Theft", "Burglary"},
				IssuingCourt: "Port-au-Prince Tribunal",
				IssuedAt:     time.Now(),
			},
		},
		{
			name: "repo error",
			warrant: &domain.Warrant{
				FullName:     "Jane Doe",
				WarrantType:  domain.WarrantTypeBench,
				Charges:      []string{"Contempt"},
				IssuingCourt: "Cap-Haitien Court",
				IssuedAt:     time.Now(),
			},
			repoErr: errors.New("insert error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockFprRepo{
				saveWarrantFn: func(w *domain.Warrant) error {
					return tt.repoErr
				},
			}
			svc := NewFprService(repo, nil)
			err := svc.CreateWarrant(tt.warrant)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, tt.warrant.ID)
			assert.False(t, tt.warrant.CreatedAt.IsZero())
		})
	}
}

func TestCheckCitizen(t *testing.T) {
	citizenID := "CIT-001"
	tests := []struct {
		name       string
		citizenID  string
		warrants   []domain.Warrant
		repoErr    error
		wantErr    bool
		wantWanted bool
	}{
		{
			name:       "clear",
			citizenID:  citizenID,
			warrants:   []domain.Warrant{},
			wantWanted: false,
		},
		{
			name:      "wanted",
			citizenID: citizenID,
			warrants: []domain.Warrant{
				{ID: uuid.New(), FullName: "John Doe", WarrantType: domain.WarrantTypeArrest},
			},
			wantWanted: true,
		},
		{
			name:      "repo error",
			citizenID: citizenID,
			repoErr:   errors.New("query error"),
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockFprRepo{
				findWarrantsByNameFn: func(name string) ([]domain.Warrant, error) {
					return tt.warrants, tt.repoErr
				},
				saveCheckLogFn: func(cl *domain.CheckLog) error {
					return nil
				},
			}
			svc := NewFprService(repo, nil)
			result, err := svc.CheckCitizen(tt.citizenID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantWanted, result.WarrantFound)
			if tt.wantWanted {
				require.NotNil(t, result.Warrant)
				assert.Equal(t, "WANTED", result.CheckLog.Result)
			} else {
				assert.Equal(t, "CLEAR", result.CheckLog.Result)
			}
		})
	}
}

func TestCheckByName(t *testing.T) {
	tests := []struct {
		name       string
		searchName string
		warrants   []domain.Warrant
		repoErr    error
		wantErr    bool
		wantWanted bool
	}{
		{
			name:       "clear",
			searchName: "Unknown Person",
			warrants:   []domain.Warrant{},
			wantWanted: false,
		},
		{
			name:       "wanted",
			searchName: "John Doe",
			warrants: []domain.Warrant{
				{ID: uuid.New(), FullName: "John Doe", WarrantType: domain.WarrantTypeArrest},
			},
			wantWanted: true,
		},
		{
			name:       "repo error",
			searchName: "Error Test",
			repoErr:    errors.New("query error"),
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockFprRepo{
				findWarrantsByNameFn: func(name string) ([]domain.Warrant, error) {
					return tt.warrants, tt.repoErr
				},
				saveCheckLogFn: func(cl *domain.CheckLog) error {
					return nil
				},
			}
			svc := NewFprService(repo, nil)
			result, err := svc.CheckByName(tt.searchName)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantWanted, result.WarrantFound)
		})
	}
}

func TestReportSighting(t *testing.T) {
	warrantID := uuid.New()
	tests := []struct {
		name      string
		warrantID uuid.UUID
		sighting  *domain.Sighting
		repoErr   error
		wantErr   bool
	}{
		{
			name:      "success",
			warrantID: warrantID,
			sighting: &domain.Sighting{
				CitizenID:   "CIT-001",
				Description: "Suspicious activity at market",
				ReportedBy:  "OFFICER-01",
				SightedAt:   time.Now(),
			},
		},
		{
			name:      "repo error",
			warrantID: warrantID,
			sighting: &domain.Sighting{
				Description: "Test sighting",
				ReportedBy:  "OFFICER-99",
				SightedAt:   time.Now(),
			},
			repoErr: errors.New("insert error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockFprRepo{
				saveSightingFn: func(s *domain.Sighting) error {
					return tt.repoErr
				},
			}
			svc := NewFprService(repo, nil)
			err := svc.ReportSighting(tt.warrantID, tt.sighting)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.warrantID, tt.sighting.WarrantID)
			assert.NotEqual(t, uuid.Nil, tt.sighting.ID)
		})
	}
}

func TestExecuteWarrant(t *testing.T) {
	warrantID := uuid.New()
	tests := []struct {
		name    string
		id      uuid.UUID
		repoErr error
		wantErr bool
	}{
		{
			name: "success",
			id:   warrantID,
		},
		{
			name:    "repo error",
			id:      warrantID,
			repoErr: errors.New("update error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockFprRepo{
				updateWarrantExecutedFn: func(id uuid.UUID, executedAt time.Time) error {
					return tt.repoErr
				},
			}
			svc := NewFprService(repo, nil)
			err := svc.ExecuteWarrant(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestGetArmedDangerous(t *testing.T) {
	tests := []struct {
		name    string
		repoRes []domain.Warrant
		repoErr error
		wantErr bool
	}{
		{
			name: "success",
			repoRes: []domain.Warrant{
				{
					ID:          uuid.New(),
					FullName:    "Armed Suspect",
					WarrantType: domain.WarrantTypeArrest,
					DangerLevel: dangerPtr(domain.DangerLevelArmedAndDangerous),
				},
			},
		},
		{
			name:    "empty",
			repoRes: []domain.Warrant{},
		},
		{
			name:    "repo error",
			repoErr: errors.New("query error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockFprRepo{
				getArmedDangerousFn: func() ([]domain.Warrant, error) {
					return tt.repoRes, tt.repoErr
				},
			}
			svc := NewFprService(repo, nil)
			got, err := svc.GetArmedDangerous()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Len(t, got, len(tt.repoRes))
		})
	}
}

func TestGetDashboardStats(t *testing.T) {
	tests := []struct {
		name    string
		repoRes *domain.DashboardStats
		repoErr error
		wantErr bool
	}{
		{
			name: "success",
			repoRes: &domain.DashboardStats{
				TotalWarrants:    100,
				ActiveWarrants:   75,
				ExecutedWarrants: 25,
				ArmedDangerous:   5,
				TotalSightings:   30,
				TotalChecks:      500,
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
			repo := &mockFprRepo{
				getDashboardStatsFn: func() (*domain.DashboardStats, error) {
					return tt.repoRes, tt.repoErr
				},
			}
			svc := NewFprService(repo, nil)
			got, err := svc.GetDashboardStats()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.repoRes.TotalWarrants, got.TotalWarrants)
		})
	}
}

func dangerPtr(d domain.DangerLevel) *domain.DangerLevel {
	return &d
}
