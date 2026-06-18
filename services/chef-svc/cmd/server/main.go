package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/snisid/platform/services/chef-svc/internal/api/rest"
	"github.com/snisid/platform/services/chef-svc/internal/repository/postgres"
	"github.com/snisid/platform/services/chef-svc/internal/service"
)

type noopPublisher struct{}

func (n *noopPublisher) PublishEvent(eventType string, payload interface{}) error {
	log.Printf("event: %s", eventType)
	return nil
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/chef_ht?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	memberRepo := postgres.NewMemberRepo(db)
	intelRepo := postgres.NewIntelNoteRepo(db)
	sightRepo := postgres.NewSightingRepo(db)
	publisher := &noopPublisher{}

	svc := service.NewMemberService(memberRepo, intelRepo, sightRepo, publisher)
	handler := rest.NewMemberHandler(svc)

	r := gin.Default()

	v1 := r.Group("/api/v1/chef")
	{
		v1.POST("/members", handler.CreateMember)
		v1.GET("/members/:id", handler.GetMember)
		v1.GET("/members/by-gang/:id", handler.GetByGang)
		v1.GET("/members/sanctioned", handler.GetSanctioned)
		v1.GET("/members/leaders", handler.GetLeaders)
		v1.POST("/members/:id/intel", handler.AddIntelligenceNote)
		v1.GET("/members/:id/intel", handler.GetIntelligenceNotes)
		v1.POST("/members/:id/sightings", handler.RecordSighting)
		v1.GET("/members/:id/sightings", handler.GetSightings)
		v1.PATCH("/members/:id/status", handler.UpdateStatus)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8097"
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	log.Printf("CHEF-HT service starting on port %s", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
