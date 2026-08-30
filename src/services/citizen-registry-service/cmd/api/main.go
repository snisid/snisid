package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/snisid/citizen-registry-service/internal/application/queries"
	"github.com/snisid/citizen-registry-service/internal/domain"
)

var getCitizenHandler *queries.GetCitizenHandler

func init() {
	// Initialize OpenSearch client and dependencies
	getCitizenHandler = queries.NewGetCitizenHandler()
}

func getCitizenHttpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract NIU from path (e.g. /v1/registry/citizens/1234567890)
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	niu := parts[4]

	query := queries.GetCitizenQuery{NIU: niu}

	result, err := getCitizenHandler.Handle(r.Context(), query)
	if err != nil {
		if err == domain.ErrCitizenNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func main() {
	http.HandleFunc("/v1/registry/citizens/", getCitizenHttpHandler)
	log.Println("Citizen Registry Service (CQRS Read Layer) listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
