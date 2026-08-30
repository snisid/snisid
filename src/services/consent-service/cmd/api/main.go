package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/snisid/consent-service/internal/application/commands"
)

var grantConsentHandler *commands.GrantConsentHandler

func init() {
	grantConsentHandler = commands.NewGrantConsentHandler()
}

func grantConsentHttpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var cmd commands.GrantConsentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := grantConsentHandler.Handle(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func main() {
	http.HandleFunc("/v1/consent/grant", grantConsentHttpHandler)
	log.Println("Consent Service listening on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
