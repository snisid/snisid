package rest

import (
	"encoding/json"
	"net/http"

	"github.com/snisid/vehicle-criminal-svc/internal/service"
)

func InterpolSyncStatusHandler(svc *service.InterpolSyncService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		syncs, err := svc.GetPendingSyncs(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"pending_syncs": syncs,
			"count":         len(syncs),
		})
	}
}
