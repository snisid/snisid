package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/snisid/vehicle-criminal-svc/internal/domain"
	"github.com/snisid/vehicle-criminal-svc/internal/service"
)

func DeclareStolenPlateHandler(svc *service.StolenPlateService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req domain.DeclareStolenPlateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		userID := r.Context().Value(userIDKey).(uuid.UUID)

		plate, err := svc.DeclareStolen(r.Context(), req, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(plate)
	}
}

func CheckStolenPlateHandler(svc *service.StolenPlateService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		plate := chi.URLParam(r, "plate")

		result, err := svc.CheckPlate(r.Context(), plate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if result == nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"plate_number": plate,
				"is_stolen":    false,
			})
			return
		}
		json.NewEncoder(w).Encode(result)
	}
}

func MarkPlateRecoveredHandler(svc *service.StolenPlateService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "invalid plate ID", http.StatusBadRequest)
			return
		}

		var req struct {
			Location string `json:"location"`
			DeptCode string `json:"dept_code"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if err := svc.MarkRecovered(r.Context(), id, req.Location, req.DeptCode, uuid.Nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
