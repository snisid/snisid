package rest

import (
	"encoding/json"
	"net/http"

	"github.com/snisid/vehicle-criminal-svc/internal/service"
)

func CreateSightingHandler(svc *service.CriminalAlertService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "sighting_endpoint_active",
		})
	}
}
