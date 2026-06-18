package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/snisid/vehicle-criminal-svc/internal/domain"
	"github.com/snisid/vehicle-criminal-svc/internal/repository"
	"github.com/snisid/vehicle-criminal-svc/internal/service"
)

func CreateAlertHandler(svc *service.CriminalAlertService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req domain.CreateAlertRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		userID := r.Context().Value(userIDKey).(uuid.UUID)

		alert, err := svc.CreateAlert(r.Context(), req, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(alert)
	}
}

func ListAlertsHandler(svc *service.CriminalAlertService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := repository.AlertFilter{
			DeptCode:      r.URL.Query().Get("dept"),
			Category:      r.URL.Query().Get("category"),
			Level:         r.URL.Query().Get("level"),
			Status:        r.URL.Query().Get("status"),
			ReportingUnit: r.URL.Query().Get("unit"),
		}

		if p := r.URL.Query().Get("page"); p != "" {
			filter.Page, _ = strconv.Atoi(p)
		}
		if l := r.URL.Query().Get("limit"); l != "" {
			filter.Limit, _ = strconv.Atoi(l)
		}

		alerts, total, err := svc.ListAlerts(r.Context(), filter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"alerts": alerts,
			"total":  total,
			"page":   filter.Page,
			"limit":  filter.Limit,
		})
	}
}

func GetAlertHandler(svc *service.CriminalAlertService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "invalid alert ID", http.StatusBadRequest)
			return
		}

		alert, err := svc.GetAlert(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if alert == nil {
			http.Error(w, "alert not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(alert)
	}
}

func UpdateAlertStatusHandler(svc *service.CriminalAlertService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "invalid alert ID", http.StatusBadRequest)
			return
		}

		var req struct {
			Status domain.AlertStatus `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		userID := r.Context().Value(userIDKey).(uuid.UUID)

		if err := svc.UpdateAlertStatus(r.Context(), id, req.Status, userID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
