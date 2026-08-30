package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func chargeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := ProcessTransaction(req.Amount); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success", "message":"Fonds répartis selon la politique de pérennité (Phase 17)"}`))
}

func fundsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	treasury, refresh, count := GetLedgerBalance()

	resp := map[string]interface{}{
		"treasury_account":      treasury,
		"hardware_refresh_fund": refresh,
		"total_transactions":    count,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	err := InitDB("file:snisid_billing.db?cache=shared&mode=rwc")
	if err != nil {
		log.Fatalf("Impossible d'initialiser la base SQLite : %v", err)
	}
	defer DB.Close()

	http.HandleFunc("/api/billing/charge", chargeHandler)
	http.HandleFunc("/api/billing/funds", fundsHandler)

	fmt.Println("SNISID Sustainability Billing API (Phase 17) démarré sur le port 8082...")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
