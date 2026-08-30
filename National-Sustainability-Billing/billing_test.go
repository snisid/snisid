package main

import (
	"testing"
)

func TestProcessTransaction(t *testing.T) {
	// Initialize an in-memory SQLite DB for testing
	err := InitDB("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to initialize test DB: %v", err)
	}
	defer DB.Close()

	// Clear DB (since cache=shared might retain data across subtests)
	DB.Exec("UPDATE ledger SET treasury_account=0, hardware_refresh_fund=0, total_transactions=0 WHERE id=1")

	err = ProcessTransaction(100.0)
	if err != nil {
		t.Fatalf("Transaction failed: %v", err)
	}

	treasury, refresh, count := GetLedgerBalance()

	if count != 1 {
		t.Errorf("Expected 1 transaction, got %d", count)
	}

	// 100 * 0.30 = 30
	if refresh != 30.0 {
		t.Errorf("Expected Refresh Fund to be 30.0, got %f", refresh)
	}

	// 100 - 30 = 70
	if treasury != 70.0 {
		t.Errorf("Expected Treasury to be 70.0, got %f", treasury)
	}
}
