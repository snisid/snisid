package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/snisid/vehicle-criminal-svc/internal/repository"
	"github.com/snisid/vehicle-criminal-svc/internal/service"
)

func SearchHandler(svc *service.CriminalAlertService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			http.Error(w, "search query required", http.StatusBadRequest)
			return
		}

		filter := repository.AlertFilter{
			DeptCode: r.URL.Query().Get("dept"),
			Category: r.URL.Query().Get("category"),
		}

		if p := r.URL.Query().Get("page"); p != "" {
			filter.Page, _ = strconv.Atoi(p)
		}
		if l := r.URL.Query().Get("limit"); l != "" {
			filter.Limit, _ = strconv.Atoi(l)
		}

		alerts, total, err := svc.SearchAlerts(r.Context(), query, filter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"alerts": alerts,
			"total":  total,
			"query":  query,
		})
	}
}
