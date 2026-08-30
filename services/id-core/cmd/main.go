package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/snisid/idcore-svc/internal/handler"
	"github.com/snisid/idcore-svc/internal/kafka"
	"github.com/snisid/idcore-svc/internal/milvus"
	"github.com/snisid/idcore-svc/internal/nin"
	"github.com/snisid/idcore-svc/internal/repository"
	"github.com/snisid/idcore-svc/internal/service"
)

func main() {
	dbHost := getEnv("IDCORE_DB_HOST", "localhost")
	dbPort := getEnv("IDCORE_DB_PORT", "26257")
	dbName := getEnv("IDCORE_DB_NAME", "snisid_idcore")
	dbUser := getEnv("IDCORE_DB_USER", "root")
	dbSSLMode := getEnv("IDCORE_DB_SSLMODE", "disable")

	milvusAddr := getEnv("IDCORE_MILVUS_ADDR", "localhost:19530")
	kafkaBrokers := getEnv("IDCORE_KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("IDCORE_KAFKA_TOPIC", "snisid.idcore.events")
	port := getEnv("IDCORE_SERVICE_PORT", "8081")
	bioThreshold := getEnv("IDCORE_DEDUP_BIO_THRESHOLD", "0.95")
	demoThreshold := getEnv("IDCORE_DEDUP_DEMO_THRESHOLD", "0.85")

	dbURL := fmt.Sprintf("postgresql://%s@%s:%s/%s?sslmode=%s", dbUser, dbHost, dbPort, dbName, dbSSLMode)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to CockroachDB: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping CockroachDB: %v", err)
	}
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := runMigrations(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	milvusClient, err := milvus.NewClient(milvusAddr)
	if err != nil {
		log.Fatalf("failed to connect to Milvus: %v", err)
	}
	defer milvusClient.Close()

	producer := kafka.NewProducer([]string{kafkaBrokers}, kafkaTopic)
	defer producer.Close()

	ninGen := nin.NewGenerator(db)

	repo := repository.NewCockroachRepo(db)

	svc := service.NewIdentityService(repo, milvusClient, producer, ninGen, bioThreshold, demoThreshold)

	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h := handler.NewHandler(svc)
	api := r.Group("/api/v1/idcore")
	h.RegisterRoutes(api)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("id-core service started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run id-core service: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down id-core service...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server forced to shutdown: %v", err)
	}
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TYPE IF NOT EXISTS id_status AS ENUM (
			'ACTIVE', 'SUSPENDED', 'DECEASED', 'CANCELLED', 'MERGED_DUPLICATE'
		)`,
		`CREATE TYPE IF NOT EXISTS id_enrollment_type AS ENUM (
			'BIRTH', 'ADULT_FIRST_ENROLLMENT', 'NATURALIZATION',
			'REGULARIZATION', 'REFUGEE_STATUS', 'RECONSTRUCTION_LOST_RECORDS'
		)`,
		`CREATE TABLE IF NOT EXISTS citizens (
			citizen_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			nin                  VARCHAR(13) UNIQUE NOT NULL,
			status               id_status NOT NULL DEFAULT 'ACTIVE',
			enrollment_type      id_enrollment_type NOT NULL,
			full_name_legal      VARCHAR(200) NOT NULL,
			first_name           VARCHAR(100) NOT NULL,
			middle_names         VARCHAR(100),
			last_name            VARCHAR(100) NOT NULL,
			maiden_name          VARCHAR(100),
			dob                  DATE NOT NULL,
			pob_commune          VARCHAR(100),
			pob_dept_code        CHAR(2),
			gender               VARCHAR(10),
			nationality          CHAR(3) NOT NULL DEFAULT 'HTI',
			dept_code            CHAR(2) NOT NULL,
			current_address      TEXT,
			current_commune      VARCHAR(100),
			biometric_template_id UUID,
			photo_ref             VARCHAR(500),
			mother_nin            VARCHAR(13),
			father_nin            VARCHAR(13),
			date_of_death         DATE,
			death_certificate_ref VARCHAR(100),
			is_merged             BOOLEAN DEFAULT FALSE,
			merged_into_citizen_id UUID,
			created_by            UUID NOT NULL,
			created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS id_change_history (
			history_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			citizen_id            UUID NOT NULL,
			field_changed         VARCHAR(100) NOT NULL,
			old_value             TEXT,
			new_value             TEXT,
			change_reason         TEXT,
			authorized_by         UUID NOT NULL,
			changed_at             TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS id_dedup_candidates (
			candidate_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			citizen_id_a           UUID NOT NULL,
			citizen_id_b           UUID NOT NULL,
			biometric_score         DECIMAL(5,4),
			demographic_score       DECIMAL(5,4),
			composite_score         DECIMAL(5,4),
			status                  VARCHAR(20) DEFAULT 'PENDING_REVIEW',
			reviewed_by              UUID,
			resolution                VARCHAR(30),
			created_at                TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_citizens_nin ON citizens(nin)`,
		`CREATE INDEX IF NOT EXISTS idx_citizens_dob ON citizens(dob)`,
		`CREATE INDEX IF NOT EXISTS idx_citizens_status ON citizens(status)`,
		`CREATE INDEX IF NOT EXISTS idx_history_citizen ON id_change_history(citizen_id, changed_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_citizens_name_fts ON citizens USING gin(to_tsvector('simple', full_name_legal))`,
	}

	for _, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			return fmt.Errorf("migration failed: %s: %w", m[:60], err)
		}
	}

	partitionStmt := `CREATE TABLE IF NOT EXISTS citizens_ouest PARTITION OF citizens FOR VALUES IN ('OU')`
	if _, err := db.Exec(partitionStmt); err != nil {
		if !isAlreadyPartitioned(err) {
			return fmt.Errorf("partition migration: %w", err)
		}
	}

	return nil
}

func isAlreadyPartitioned(err error) bool {
	return err != nil && (contains(err.Error(), "already exists") || contains(err.Error(), "already a partition"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func getEnv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
