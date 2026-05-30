package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/snisid/identity-service/internal/application/commands"
)

var enrollHandler *commands.EnrollCitizenHandler

func init() {
	// Initialize dependencies
	enrollHandler = commands.NewEnrollCitizenHandler()
}

func registerCitizenHttpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var cmd commands.EnrollCitizenCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := enrollHandler.Handle(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"citizenId":     result.CitizenID,
		"status":        result.Status,
		"correlationId": uuid.New().String(),
	})
}

func main() {
	http.HandleFunc("/v1/identity/citizens", registerCitizenHttpHandler)
	log.Println("Identity Service (CQRS Command Layer) listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
