package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/snisid/vehicle-criminal-svc/internal/service"
)

func CheckPlateHandler(svc *service.CriminalAlertService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		plate := chi.URLParam(r, "plate")
		if plate == "" {
			http.Error(w, "plate number required", http.StatusBadRequest)
			return
		}

		result, err := svc.CheckPlate(r.Context(), plate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

func CheckVINHandler(svc *service.CriminalAlertService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vin := chi.URLParam(r, "vin")
		if vin == "" {
			http.Error(w, "VIN required", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"vin":    vin,
			"status": "not_implemented",
		})
	}
}
