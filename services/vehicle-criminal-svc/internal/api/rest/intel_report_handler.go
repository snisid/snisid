package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/snisid/vehicle-criminal-svc/internal/domain"
	"github.com/snisid/vehicle-criminal-svc/internal/service"
)

func CreateIntelReportHandler(svc *service.VehicleIntelService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var report domain.IntelligenceReport
		if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		userID := r.Context().Value(userIDKey).(uuid.UUID)

		if err := svc.CreateReport(r.Context(), &report, userID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(report)
	}
}

func GetIntelReportHandler(svc *service.VehicleIntelService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "invalid report ID", http.StatusBadRequest)
			return
		}

		report, err := svc.GetReport(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if report == nil {
			http.Error(w, "report not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(report)
	}
}
