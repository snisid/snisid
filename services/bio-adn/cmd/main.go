package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/services/bio-adn/internal/db"
	"github.com/snisid/platform/services/bio-adn/internal/engine"
	"github.com/snisid/platform/services/bio-adn/internal/indexes"
	"github.com/snisid/platform/services/bio-adn/internal/interpol"
	"github.com/snisid/platform/services/bio-adn/internal/kafka"
	"github.com/snisid/platform/services/bio-adn/internal/sync"
	"github.com/snisid/platform/services/bio-adn/pkg/models"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	dbURL := getEnv("DATABASE_URL", "postgres://snisid:snisid@localhost:5432/snisid_bio?sslmode=disable")
	kafkaBrokers := getEnv("KAFKA_BROKERS", "localhost:9092")

	pg, err := db.NewPostgres(ctx, dbURL)
	if err != nil {
		logger.Fatal("failed to connect to postgres", zap.Error(err))
	}
	defer pg.Close()

	matcher := engine.NewMatcher()
	lapiSvc := indexes.NewLAPIQueryService(redisAddr, pg)
	ndisMatcher := engine.NewNDISMatcher(pg)
	reportScheduler := sync.NewNDISReportScheduler(pg, logger)
	interpolGateway := interpol.NewGateway(
		getEnv("INTERPOL_ENDPOINT", "gateway.interpol.int:443"),
		getEnv("INTERPOL_CERT", "SNISID-PKI/certs/bcn-dcpj.crt"),
		getEnv("INTERPOL_KEY", "SNISID-PKI/keys/bcn-dcpj.key"),
		getEnv("INTERPOL_CA", "SNISID-PKI/ca/interpol-ca.crt"),
	)
	_ = interpolGateway

	producer := kafka.NewProducer(kafkaBrokers, logger)
	lapiCache := indexes.NewLAPICache(redisAddr)

	handlers := map[string]kafka.HandlerFunc{
		"snisid.bio.profile.created": func(ctx context.Context, msg map[string]any) error {
			sampleID, ok := msg["sample_id"].(string)
			if !ok {
				logger.Warn("profile.created missing required field: sample_id")
				return nil
			}
			specimenNumber, ok := msg["specimen_number"].(string)
			if !ok {
				logger.Warn("profile.created missing required field: specimen_number")
				return nil
			}
			indexType, ok := msg["index_type"].(string)
			if !ok {
				logger.Warn("profile.created missing required field: index_type")
				return nil
			}
			qualityScore, _ := msg["quality_score"].(float64)
			correlationID, _ := msg["correlation_id"].(string)

			lociRaw, ok := msg["loci_data"].(map[string]any)
			if !ok {
				logger.Warn("profile.created missing loci_data", zap.String("sample_id", sampleID))
				return nil
			}

			loci := make(engine.STRLoci)
			for locus, valsRaw := range lociRaw {
				vals, ok := valsRaw.(map[string]any)
				if !ok {
					continue
				}
				val1, _ := vals["value1"].(string)
				val2, _ := vals["value2"].(string)
				loci[locus] = engine.Locus{Value1: val1, Value2: val2}
			}

			hash := matcher.HashProfile(loci)
			profile := &models.DNAProfile{
				SampleID:       sampleID,
				SpecimenNumber: specimenNumber,
				IndexType:      indexType,
				LociHash:       hash,
				QualityScore:   qualityScore,
				LociCount:      len(loci),
				LabID:          "LDIS-PAP-001",
			}

			if err := pg.CreateDNAProfile(ctx, profile); err != nil {
				logger.Error("failed to save DNA profile", zap.String("sample_id", sampleID), zap.Error(err))
				return err
			}

			existing, err := pg.GetDNAProfileByHash(ctx, hash)
			if err != nil {
				logger.Error("failed to check duplicate hash", zap.String("sample_id", sampleID), zap.Error(err))
				return err
			}
			if existing != nil && existing.SampleID != sampleID {
				result := matcher.Compare(loci, loci)
				logger.Warn("duplicate DNA profile detected",
					zap.String("new_sample_id", sampleID),
					zap.String("existing_sample_id", existing.SampleID),
					zap.Float64("score", result.Score),
				)
				hitType := kafka.MatchTypeFull
				_ = hitType
				producer.PublishHitDetected(ctx, &kafka.DNAHitDetected{
					EventEnvelope: kafka.EventEnvelope{
						EventID:       fmt.Sprintf("HIT-%d", time.Now().UnixNano()),
						EventType:     "HIT_FOUND",
						CorrelationID: correlationID,
						Timestamp:     time.Now().UnixMilli(),
					},
					HitID:          fmt.Sprintf("HIT-%d", time.Now().UnixNano()),
					QuerySampleID:  sampleID,
					MatchSampleID:  existing.SampleID,
					MatchType:      kafka.MatchTypeFull,
					Confidence:     float32(result.Score),
					MatchedLoci:    int32(result.MatchedLoci),
					TotalLoci:      int32(result.TotalLoci),
					HitLevel:       result.AlertLevel,
					AlertLevel:     kafka.AlertLevelCritical,
					QueryIndexType: indexType,
					MatchIndexType: existing.IndexType,
				})
			}

			logger.Info("profile created processed", zap.String("sample_id", sampleID), zap.Int("loci_count", len(loci)))
			return nil
		},
		"snisid.bio.profile.uploaded": func(ctx context.Context, msg map[string]any) error {
			sampleID, ok := msg["sample_id"].(string)
			if !ok {
				logger.Warn("profile.uploaded missing required field: sample_id")
				return nil
			}
			indexType, ok := msg["index_type"].(string)
			if !ok {
				logger.Warn("profile.uploaded missing required field: index_type")
				return nil
			}
			sourceSDIS, _ := msg["source_sdis"].(string)

			profile, err := pg.GetDNAProfileBySpecimen(ctx, sampleID)
			if err != nil || profile == nil {
				logger.Warn("profile not found for NDIS matching", zap.String("sample_id", sampleID), zap.Error(err))
				return nil
			}

			result, err := ndisMatcher.MatchCrossDept(ctx, sampleID, profile.LociHash, indexType, sourceSDIS)
			if err != nil {
				logger.Error("NDIS cross-dept matching failed", zap.String("sample_id", sampleID), zap.Error(err))
				return nil
			}
			if result != nil {
				logger.Warn("NDIS cross-dept hit detected",
					zap.String("hit_id", result.HitID),
					zap.String("query_sample_id", result.QuerySampleID),
					zap.String("match_sample_id", result.MatchSampleID),
					zap.String("match_type", result.MatchType),
					zap.String("alert_level", result.AlertLevel),
				)
				producer.PublishCrossDeptHit(ctx, &kafka.CrossDeptHitDetected{
					EventEnvelope: kafka.EventEnvelope{
						EventID:   result.HitID,
						EventType: "CROSS_DEPT_HIT_DETECTED",
						Timestamp: time.Now().UnixMilli(),
					},
					HitID:         result.HitID,
					QuerySampleID: result.QuerySampleID,
					MatchSampleID: result.MatchSampleID,
					MatchType:     result.MatchType,
					Confidence:    result.Confidence,
					QuerySDIS:     result.QuerySDIS,
					MatchSDIS:     result.MatchSDIS,
				})
			}

			logger.Info("profile uploaded for NDIS matching", zap.String("sample_id", sampleID))
			return nil
		},
		"snisid.bio.hits": func(ctx context.Context, msg map[string]any) error {
			hitID, ok := msg["hit_id"].(string)
			if !ok {
				logger.Warn("hits event missing required field: hit_id")
				return nil
			}
			logger.Info("hit detected", zap.String("hit_id", hitID))
			producer.PublishAuditEvent(ctx, map[string]any{
				"table_name": "bio_hits",
				"record_id":  hitID,
				"action":     "HIT",
				"details":    msg,
			})
			return nil
		},
		"snisid.bio.lapi.query": func(ctx context.Context, msg map[string]any) error {
			plate, ok := msg["plate_number"].(string)
			if !ok {
				logger.Warn("lapi.query missing required field: plate_number")
				return nil
			}
			cameraID, _ := msg["camera_id"].(string)
			hit, err := lapiSvc.QueryPlate(ctx, plate)
			if err != nil {
				return err
			}

			cloneWarning := false
			cloneCount, err := pg.QueryPlateClones(ctx, plate)
			if err == nil && cloneCount > 1 {
				cloneWarning = true
				logger.Warn("cloned plate detected", zap.String("plate", plate), zap.Int("records", cloneCount))
			}

			if cameraID != "" {
				_ = lapiCache.RecordPlateSighting(ctx, plate, cameraID)
			}

			if hit.HitFound {
				hitType := "STOLEN_VEHICLE"
				recordNumber := hit.RecordNumber
				alertLevel := hit.AlertLevel
				mcoContact := hit.MCOContact
				if cloneWarning {
					alertLevel = "CRITICAL"
				}
				producer.PublishLAPIResponse(ctx, &kafka.LAPIPlateResponse{
					PlateNumber:  plate,
					HitFound:     true,
					HitType:      &hitType,
					RecordNumber: &recordNumber,
					AlertLevel:   &alertLevel,
					MCOContact:   &mcoContact,
					ResponseMs:   int32(hit.ResponseMs),
				})
			}
			return nil
		},
		"snisid.bio.vehicle.recovered": func(ctx context.Context, msg map[string]any) error {
			id, ok := msg["record_id"].(string)
			if !ok {
				logger.Warn("vehicle.recovered missing required field: record_id")
				return nil
			}
			location, _ := msg["recovered_location"].(string)
			agency, _ := msg["recovering_agency"].(string)
			logger.Info("vehicle recovered", zap.String("record_id", id))
			return pg.UpdateVehicleStatus(ctx, id, "RECOVERED", location, agency)
		},
		"snisid.bio.arm.hit": func(ctx context.Context, msg map[string]any) error {
			logger.Warn("firearm crime scene hit", zap.Any("hit_id", msg["hit_id"]))
			return nil
		},
		"snisid.bio.expunge.events": func(ctx context.Context, msg map[string]any) error {
			sampleID, ok := msg["sample_id"].(string)
			if !ok {
				logger.Warn("expunge event missing required field: sample_id")
				return nil
			}
			logger.Warn("expunge event", zap.String("sample_id", sampleID))
			return pg.MarkExpunged(ctx, sampleID)
		},
		"snisid.bio.audit.events": func(ctx context.Context, msg map[string]any) error {
			return pg.WriteAuditLog(ctx, msg)
		},
		"snisid.oni.document.revoked": func(ctx context.Context, msg map[string]any) error {
			logger.Warn("ONI document revoked", zap.String("doc_number", toString(msg["document_number"])))
			return nil
		},
		"snisid.bio.fugitive.events": func(ctx context.Context, msg map[string]any) error {
			logger.Info("foreign fugitive", zap.String("record_id", toString(msg["record_id"])), zap.String("country", toString(msg["issuing_country"])))
			return nil
		},
		"snisid.bio.unidentified.events": func(ctx context.Context, msg map[string]any) error {
			logger.Info("unidentified person", zap.String("record_id", toString(msg["record_id"])), zap.String("location", toString(msg["discovery_location"])))
			return nil
		},
		"snisid.bio.terrorism.events": func(ctx context.Context, msg map[string]any) error {
			logger.Warn("terrorism watch", zap.String("record_id", toString(msg["record_id"])), zap.String("threat_type", toString(msg["threat_type"])))
			return nil
		},
		"snisid.bio.protection.events": func(ctx context.Context, msg map[string]any) error {
			logger.Info("protection order", zap.String("record_id", toString(msg["record_id"])), zap.String("order_type", toString(msg["order_type"])))
			return nil
		},
		"snisid.bio.supervised.events": func(ctx context.Context, msg map[string]any) error {
			logger.Info("supervised release", zap.String("record_id", toString(msg["record_id"])), zap.String("type", toString(msg["supervision_type"])))
			return nil
		},
		"snisid.bio.lab.duplicate": func(ctx context.Context, msg map[string]any) error {
			logger.Warn("duplicate specimen detected", zap.String("specimen", toString(msg["specimen_number"])))
			return nil
		},
		"snisid.bio.lab.equipment": func(ctx context.Context, msg map[string]any) error {
			logger.Info("equipment registered", zap.String("equipment_id", toString(msg["equipment_id"])))
			return nil
		},
		"snisid.bio.lab.training": func(ctx context.Context, msg map[string]any) error {
			logger.Info("training recorded", zap.String("training_id", toString(msg["training_id"])))
			return nil
		},
		"snisid.bio.lab.upload": func(ctx context.Context, msg map[string]any) error {
			labCode := toString(msg["lab_code"])
			count, ok := msg["uploaded_count"].(float64)
			if !ok {
				logger.Warn("lab.upload missing required field: uploaded_count")
				return nil
			}
			logger.Info("LDIS upload completed", zap.String("lab_code", labCode), zap.Int("count", int(count)))
			return nil
		},
		"snisid.bio.ndis.crossdept.hit": func(ctx context.Context, msg map[string]any) error {
			logger.Warn("cross-dept NDIS hit",
				zap.String("hit_id", toString(msg["hit_id"])),
				zap.String("query_sdis", toString(msg["query_sdis"])),
				zap.String("match_sdis", toString(msg["match_sdis"])),
			)
			return nil
		},
		"snisid.bio.ndis.interpol": func(ctx context.Context, msg map[string]any) error {
			logger.Info("INTERPOL submission", zap.String("submission_id", toString(msg["submission_id"])))
			return nil
		},
		"snisid.bio.ndis.reports": func(ctx context.Context, msg map[string]any) error {
			logger.Info("NDIS report", zap.String("report_id", toString(msg["report_id"])), zap.String("type", toString(msg["report_type"])))
			return nil
		},
		"snisid.bio.violence.events": func(ctx context.Context, msg map[string]any) error {
			logger.Info("violence record", zap.String("record_id", toString(msg["record_id"])), zap.String("type", toString(msg["incident_type"])))
			return nil
		},
		"snisid.bio.identitytheft.events": func(ctx context.Context, msg map[string]any) error {
			logger.Warn("identity theft", zap.String("record_id", toString(msg["record_id"])), zap.String("fraud_type", toString(msg["fraud_type"])))
			return nil
		},
		"snisid.bio.identity.linked": func(ctx context.Context, msg map[string]any) error {
			logger.Warn("identity linked",
				zap.String("sample_id", toString(msg["sample_id"])),
				zap.String("niu", toString(msg["niu"])),
				zap.String("linked_by", toString(msg["linked_by"])),
			)
			return nil
		},
	}

	consumers := kafka.StartAllConsumers(ctx, kafkaBrokers, logger, handlers)
	defer func() {
		for _, c := range consumers {
			c.Close()
		}
	}()

	syncScheduler := sync.NewSyncScheduler("LDIS", pg, producer)
	syncScheduler.Start(ctx)

	reportScheduler.Start(ctx)

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/ready", func(c *gin.Context) {
		if err := pg.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})
	r.GET("/v1/bio-adn/lapi/plate/:plate", func(c *gin.Context) {
		plate := c.Param("plate")
		cameraID := c.Query("camera_id")
		hit, err := lapiSvc.QueryPlate(ctx, plate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		cloneWarning := false
		cloneCount, err := pg.QueryPlateClones(ctx, plate)
		if err == nil && cloneCount > 1 {
			cloneWarning = true
		}

		if cameraID != "" {
			go func() { _ = lapiCache.RecordPlateSighting(context.Background(), plate, cameraID) }()
		}

		result := map[string]any{
			"hit_found":     hit.HitFound,
			"hit_type":      hit.HitType,
			"record_number": hit.RecordNumber,
			"alert_level":   hit.AlertLevel,
			"mco_contact":   hit.MCOContact,
			"response_ms":   hit.ResponseMs,
			"clone_warning": cloneWarning,
		}
		if cloneWarning {
			result["alert_level"] = "CRITICAL"
		}
		c.JSON(http.StatusOK, result)
	})
	r.GET("/v1/bio-adn/lapi/vin/:vin", func(c *gin.Context) {
		vin := c.Param("vin")
		hit, err := lapiSvc.QueryVIN(ctx, vin)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, hit)
	})

	srv := &http.Server{Addr: ":" + getEnv("PORT", "8092"), Handler: r}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}

func toString(v any) string {
	switch s := v.(type) {
	case string:
		return s
	case int:
		return strconv.Itoa(s)
	default:
		return ""
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
