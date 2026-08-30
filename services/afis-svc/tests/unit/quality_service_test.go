package unit

import (
	"testing"

	"github.com/snisid/platform/services/afis-svc/internal/domain"
	"github.com/snisid/platform/services/afis-svc/internal/service"
)

func TestNFIQ2_Score_Threshold(t *testing.T) {
	svc := service.NewQualityService(60)

	tests := []struct {
		score   int16
		wantErr bool
	}{
		{75, false},
		{60, false},
		{59, true},
		{0, true},
		{100, false},
	}

	for _, tt := range tests {
		err := svc.ValidateScore(tt.score)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateScore(%d): got err=%v, wantErr=%v", tt.score, err, tt.wantErr)
		}
	}
}

func TestNFIQ2_HighQuality(t *testing.T) {
	svc := service.NewQualityService(60)

	if !svc.IsHighQuality(85) {
		t.Error("expected score 85 to be high quality")
	}
	if svc.IsHighQuality(70) {
		t.Error("expected score 70 to NOT be high quality")
	}
}

func TestNFIQ2_QualityTooLow_Error(t *testing.T) {
	err := domain.ErrQualityTooLow
	if err == nil {
		t.Fatal("expected ErrQualityTooLow")
	}
}
