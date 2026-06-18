package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type SADHandler struct {
	db         *sqlx.DB
	baseURL    string
	apiKey     string
	httpClient *http.Client
	logger     *zap.Logger
}

func NewSADHandler(db *sqlx.DB, baseURL, apiKey string, logger *zap.Logger) *SADHandler {
	return &SADHandler{
		db:      db,
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

type PendingSADSync struct {
	SyncID      string `db:"sync_id"`
	PlateID     string `db:"stolen_plate_id"`
	PlateNumber string `db:"plate_number"`
}

func (h *SADHandler) SyncPending(ctx context.Context) error {
	var pending []PendingSADSync
	query := `
		SELECT s.sync_id, s.stolen_plate_id, p.plate_number
		FROM sivc_interpol_sync_log s
		JOIN sivc_stolen_plates p ON s.stolen_plate_id = p.plate_id
		WHERE s.sync_status = 'PENDING' AND s.sync_direction = 'OUTBOUND'
		ORDER BY s.sync_timestamp ASC
		LIMIT 50
	`

	if err := h.db.SelectContext(ctx, &pending, query); err != nil {
		return fmt.Errorf("failed to fetch pending SAD syncs: %w", err)
	}

	for _, p := range pending {
		if err := h.syncRecord(ctx, p); err != nil {
			h.logger.Error("Failed to sync SAD record",
				zap.String("sync_id", p.SyncID),
				zap.Error(err),
			)
			h.updateStatus(ctx, p.SyncID, "FAILED", err.Error())
			continue
		}
		h.updateStatus(ctx, p.SyncID, "SUCCESS", "")
	}

	return nil
}

func (h *SADHandler) syncRecord(ctx context.Context, record PendingSADSync) error {
	payload, _ := json.Marshal(map[string]string{
		"stolen_plate_id": record.PlateID,
		"plate_number":    record.PlateNumber,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", h.baseURL+"/sad/documents", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.apiKey)

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("INTERPOL SAD returned status %d", resp.StatusCode)
	}

	return nil
}

func (h *SADHandler) updateStatus(ctx context.Context, syncID, status, errorMsg string) {
	query := `UPDATE sivc_interpol_sync_log SET sync_status = $1, error_message = $2, processed_at = NOW() WHERE sync_id = $3`
	h.db.ExecContext(ctx, query, status, errorMsg, syncID)
}
