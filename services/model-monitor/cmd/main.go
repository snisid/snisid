package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/snisid/platform/services/model-monitor/internal/service"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	metrics := &service.Metrics{}

	mux := http.NewServeMux()
	mux.HandleFunc("/predict", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			Pred    float64 `json:"pred"`
			Label   int     `json:"label"`
			Version string  `json:"version"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}
		if req.Version == "" {
			req.Version = "v1"
		}
		metrics.Update(req.Pred, req.Label, req.Version)
		drift := metrics.CalculateDrift(0.5, req.Version)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"version": req.Version,
			"total":   metrics.Total,
			"drift":   drift,
		})
	})
	mux.HandleFunc("/drift", func(w http.ResponseWriter, r *http.Request) {
		baseline := 0.5
		if b := r.URL.Query().Get("baseline"); b != "" {
			if v, err := strconv.ParseFloat(b, 64); err == nil {
				baseline = v
			}
		}
		version := r.URL.Query().Get("version")
		if version == "" {
			version = "v1"
		}
		drift := metrics.CalculateDrift(baseline, version)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"drift":   drift,
			"version": version,
		})
	})
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	srv := &http.Server{Addr: ":8101", Handler: mux}
	go func() {
		slog.Info("model-monitor service starting", "addr", ":8101")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to serve", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down model-monitor service")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
