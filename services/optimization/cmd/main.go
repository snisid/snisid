package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	sandbox "github.com/snisid/platform/services/optimization/rl-sandbox"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	monitor := &sandbox.DriftMonitor{Baseline: 0.5}

	mux := http.NewServeMux()
	mux.HandleFunc("/reward", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var state sandbox.EnvState
		if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}
		reward := sandbox.Reward(state)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]float64{"reward": reward})
	})
	mux.HandleFunc("/drift", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			Current float64 `json:"current"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}
		risk := monitor.DetectCollapseRisk(req.Current)
		if risk {
			sandbox.SafetyFreeze()
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"collapse_risk": risk,
			"current":       req.Current,
			"baseline":      monitor.Baseline,
		})
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	srv := &http.Server{Addr: ":8103", Handler: mux}
	go func() {
		slog.Info("optimization service starting", "addr", ":8103")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to serve", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down optimization service")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
