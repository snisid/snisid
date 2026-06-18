package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/internal/platform/events"
	"github.com/snisid/platform/internal/platform/logger"
	"github.com/snisid/platform/internal/platform/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type EnrollmentStatus string

const (
	EnrollmentPending    EnrollmentStatus = "PENDING"
	EnrollmentBiometrics EnrollmentStatus = "BIOMETRICS_CAPTURED"
	EnrollmentVerified   EnrollmentStatus = "VERIFIED"
	EnrollmentCompleted  EnrollmentStatus = "COMPLETED"
	EnrollmentRejected   EnrollmentStatus = "REJECTED"
	EnrollmentExpired    EnrollmentStatus = "EXPIRED"
)

type Enrollment struct {
	ID             string           `json:"id" gorm:"primaryKey"`
	NNU            string           `json:"nnu" gorm:"uniqueIndex"`
	FirstName      string           `json:"firstName"`
	LastName       string           `json:"lastName"`
	MiddleName     string           `json:"middleName"`
	DateOfBirth    string           `json:"dateOfBirth"`
	PlaceOfBirth   string           `json:"placeOfBirth"`
	Gender         string           `json:"gender"`
	Nationality    string           `json:"nationality"`
	Email          string           `json:"email"`
	Phone          string           `json:"phone"`
	AddressJSON    string           `json:"addressJson"`
	PhotoURL       string           `json:"photoUrl"`
	AgencyID       string           `json:"agencyId"`
	AgentID        string           `json:"agentId"`
	LocationID     string           `json:"locationId"`
	DeviceID       string           `json:"deviceId"`
	Status         EnrollmentStatus `json:"status"`
	StatusReason   string           `json:"statusReason"`
	BiometricHash  string           `json:"biometricHash"`
	DocumentHashes string           `json:"documentHashes"`
	QualityScore   float64          `json:"qualityScore"`
	FraudScore     int              `json:"fraudScore"`
	Version        int              `json:"version"`
	OfflineSync    bool             `json:"offlineSync"`
	OtpVerified    bool             `json:"otpVerified"`
	ConsentGiven   bool             `json:"consentGiven"`
	ConsentRecord  string           `json:"consentRecord"`
	CreatedAt      time.Time        `json:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt"`
	CompletedAt    *time.Time       `json:"completedAt,omitempty"`
	ExpiresAt      time.Time        `json:"expiresAt"`
}

type EnrollmentEvent struct {
	EventType    string      `json:"eventType"`
	EnrollmentID string      `json:"enrollmentId"`
	NNU          string      `json:"nnu"`
	Status       string      `json:"status"`
	Timestamp    time.Time   `json:"timestamp"`
	Data         interface{} `json:"data,omitempty"`
}

type CreateEnrollmentRequest struct {
	FirstName    string                 `json:"firstName" binding:"required"`
	LastName     string                 `json:"lastName" binding:"required"`
	MiddleName   string                 `json:"middleName"`
	DateOfBirth  string                 `json:"dateOfBirth" binding:"required"`
	PlaceOfBirth string                 `json:"placeOfBirth" binding:"required"`
	Gender       string                 `json:"gender" binding:"required"`
	Nationality  string                 `json:"nationality" binding:"required"`
	Email        string                 `json:"email"`
	Phone        string                 `json:"phone"`
	Address      map[string]interface{} `json:"address"`
	AgencyID     string                 `json:"agencyId" binding:"required"`
	AgentID      string                 `json:"agentId" binding:"required"`
	DeviceID     string                 `json:"deviceId"`
	LocationID   string                 `json:"locationId"`
	ConsentGiven bool                   `json:"consentGiven"`
	OfflineSync  bool                   `json:"offlineSync"`
}

type SearchEnrollmentRequest struct {
	Status     string `form:"status"`
	AgencyID   string `form:"agencyId"`
	AgentID    string `form:"agentId"`
	SearchTerm string `form:"search"`
	FromDate   string `form:"fromDate"`
	ToDate     string `form:"toDate"`
	Page       int    `form:"page"`
	PageSize   int    `form:"pageSize"`
}

type CaptureBiometricRequest struct {
	BiometricHash string  `json:"biometricHash" binding:"required"`
	QualityScore  float64 `json:"qualityScore" binding:"required"`
}

type VerifyEnrollmentRequest struct {
	OtpVerified bool   `json:"otpVerified"`
	OtpCode     string `json:"otpCode,omitempty"`
}

type RejectEnrollmentRequest struct {
	Reason string `json:"reason" binding:"required"`
}

type EnrollmentService struct {
	db       *gorm.DB
	producer *events.Producer
}

func NewEnrollmentService(db *gorm.DB, producer *events.Producer) *EnrollmentService {
	return &EnrollmentService{db: db, producer: producer}
}

func (s *EnrollmentService) CreateEnrollment(req CreateEnrollmentRequest) (*Enrollment, error) {
	enrollment := &Enrollment{
		ID:           newUUID(),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		MiddleName:   req.MiddleName,
		DateOfBirth:  req.DateOfBirth,
		PlaceOfBirth: req.PlaceOfBirth,
		Gender:       req.Gender,
		Nationality:  req.Nationality,
		Email:        req.Email,
		Phone:        req.Phone,
		AgencyID:     req.AgencyID,
		AgentID:      req.AgentID,
		DeviceID:     req.DeviceID,
		LocationID:   req.LocationID,
		Status:       EnrollmentPending,
		ConsentGiven: req.ConsentGiven,
		Version:      1,
		OfflineSync:  req.OfflineSync,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		ExpiresAt:    time.Now().UTC().Add(72 * time.Hour),
	}

	if req.Address != nil {
		b, _ := json.Marshal(req.Address)
		enrollment.AddressJSON = string(b)
	}

	if err := s.db.Create(enrollment).Error; err != nil {
		return nil, fmt.Errorf("failed to create enrollment: %w", err)
	}

	s.publishEvent("enrollment.created", enrollment.ID, "", "PENDING", map[string]interface{}{
		"firstName": req.FirstName,
		"lastName":  req.LastName,
		"agencyId":  req.AgencyID,
	})

	return enrollment, nil
}

func (s *EnrollmentService) CaptureBiometrics(id string, biometricHash string, qualityScore float64) (*Enrollment, error) {
	enrollment, err := s.getEnrollment(id)
	if err != nil {
		return nil, err
	}
	if enrollment.Status != EnrollmentPending {
		return nil, fmt.Errorf("cannot capture biometrics for enrollment in status: %s", enrollment.Status)
	}
	if qualityScore < 0.5 {
		return nil, fmt.Errorf("biometric quality too low: %.2f", qualityScore)
	}

	enrollment.BiometricHash = biometricHash
	enrollment.QualityScore = qualityScore
	enrollment.Status = EnrollmentBiometrics
	enrollment.Version++
	enrollment.UpdatedAt = time.Now().UTC()

	if err := s.db.Save(enrollment).Error; err != nil {
		return nil, err
	}

	s.publishEvent("enrollment.biometrics_captured", enrollment.ID, enrollment.NNU, string(EnrollmentBiometrics),
		map[string]interface{}{"qualityScore": qualityScore})
	return enrollment, nil
}

func (s *EnrollmentService) VerifyEnrollment(id string, otpVerified bool) (*Enrollment, error) {
	enrollment, err := s.getEnrollment(id)
	if err != nil {
		return nil, err
	}
	if enrollment.Status != EnrollmentBiometrics {
		return nil, fmt.Errorf("cannot verify enrollment in status: %s", enrollment.Status)
	}

	nnu, err := s.generateNNU(enrollment)
	if err != nil {
		return nil, fmt.Errorf("failed to generate NNU: %w", err)
	}

	enrollment.NNU = nnu
	enrollment.Status = EnrollmentVerified
	enrollment.OtpVerified = otpVerified
	enrollment.Version++
	enrollment.UpdatedAt = time.Now().UTC()

	if err := s.db.Save(enrollment).Error; err != nil {
		return nil, err
	}

	s.publishEvent("enrollment.verified", enrollment.ID, nnu, string(EnrollmentVerified),
		map[string]interface{}{"nnu": nnu})
	return enrollment, nil
}

func (s *EnrollmentService) CompleteEnrollment(id string) (*Enrollment, error) {
	enrollment, err := s.getEnrollment(id)
	if err != nil {
		return nil, err
	}
	if enrollment.Status != EnrollmentVerified {
		return nil, fmt.Errorf("cannot complete enrollment in status: %s", enrollment.Status)
	}

	now := time.Now().UTC()
	enrollment.Status = EnrollmentCompleted
	enrollment.CompletedAt = &now
	enrollment.Version++
	enrollment.UpdatedAt = now

	if err := s.db.Save(enrollment).Error; err != nil {
		return nil, err
	}
	s.publishEvent("enrollment.completed", enrollment.ID, enrollment.NNU, string(EnrollmentCompleted), nil)
	return enrollment, nil
}

func (s *EnrollmentService) RejectEnrollment(id, reason string) (*Enrollment, error) {
	enrollment, err := s.getEnrollment(id)
	if err != nil {
		return nil, err
	}
	enrollment.Status = EnrollmentRejected
	enrollment.StatusReason = reason
	enrollment.Version++
	enrollment.UpdatedAt = time.Now().UTC()

	if err := s.db.Save(enrollment).Error; err != nil {
		return nil, err
	}
	s.publishEvent("enrollment.rejected", enrollment.ID, enrollment.NNU, string(EnrollmentRejected),
		map[string]interface{}{"reason": reason})
	return enrollment, nil
}

func (s *EnrollmentService) GetEnrollment(id string) (*Enrollment, error) {
	return s.getEnrollment(id)
}

func (s *EnrollmentService) SearchEnrollments(req SearchEnrollmentRequest) ([]Enrollment, int64, error) {
	query := s.db.Model(&Enrollment{})
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.AgencyID != "" {
		query = query.Where("agency_id = ?", req.AgencyID)
	}
	if req.AgentID != "" {
		query = query.Where("agent_id = ?", req.AgentID)
	}
	if req.SearchTerm != "" {
		search := "%" + req.SearchTerm + "%"
		query = query.Where("(first_name ILIKE ? OR last_name ILIKE ? OR nnu ILIKE ?)", search, search, search)
	}
	if req.FromDate != "" {
		query = query.Where("created_at >= ?", req.FromDate)
	}
	if req.ToDate != "" {
		query = query.Where("created_at <= ?", req.ToDate)
	}

	var total int64
	query.Count(&total)

	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var enrollments []Enrollment
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&enrollments)

	return enrollments, total, nil
}

func (s *EnrollmentService) getEnrollment(id string) (*Enrollment, error) {
	var enrollment Enrollment
	if err := s.db.First(&enrollment, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("enrollment not found: %s", id)
	}
	return &enrollment, nil
}

func (s *EnrollmentService) generateNNU(enrollment *Enrollment) (string, error) {
	for attempts := 0; attempts < 10; attempts++ {
		department := enrollment.LocationID
		if len(department) >= 3 {
			department = department[:3]
		} else {
			department = "001"
		}
		birthYear := "0000"
		if len(enrollment.DateOfBirth) >= 4 {
			birthYear = enrollment.DateOfBirth[:4]
		}
		randomBytes := make([]byte, 4)
		if _, err := rand.Read(randomBytes); err != nil {
			return "", err
		}
		randomPart := hex.EncodeToString(randomBytes)[:7]
		nnu := strings.ToUpper(fmt.Sprintf("%s%s%s", department, birthYear[2:], randomPart))
		if len(nnu) > 12 {
			nnu = nnu[:12]
		} else if len(nnu) < 12 {
			nnu = nnu + strings.Repeat("0", 12-len(nnu))
		}

		var count int64
		s.db.Model(&Enrollment{}).Where("nnu = ?", nnu).Count(&count)
		if count == 0 {
			return nnu, nil
		}
	}
	return "", fmt.Errorf("failed to generate unique NNU after 10 attempts")
}

func (s *EnrollmentService) publishEvent(eventType, enrollmentID, nnu, status string, data interface{}) {
	if s.producer == nil {
		return
	}
	evt := EnrollmentEvent{
		EventType:    eventType,
		EnrollmentID: enrollmentID,
		NNU:          nnu,
		Status:       status,
		Timestamp:    time.Now().UTC(),
		Data:         data,
	}
	_ = s.producer.Publish(context.Background(), enrollmentID, evt)
}

func main() {
	jwtSecret := getEnv("JWT_SECRET", "dev-secret")
	broker := getEnv("KAFKA_BROKER", "localhost:9092")
	dbURL := getEnv("DATABASE_URL", "host=localhost user=snisid password=snisid dbname=snisid port=5432 sslmode=disable")
	port := getEnv("PORT", "8083")

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	if getEnv("ENV", "dev") == "dev" {
		db.AutoMigrate(&Enrollment{})
	}

	producer := events.NewProducer([]string{broker}, "enrollment.events")
	defer producer.Close()

	svc := NewEnrollmentService(db, producer)

	r := gin.Default()
	r.Use(middleware.RateLimit(30, 60))
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	api := r.Group("/v1/enrollments", middleware.Auth(jwtSecret))
	{
		api.POST("", func(c *gin.Context) {
			var req CreateEnrollmentRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			enrollment, err := svc.CreateEnrollment(req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, enrollment)
		})

		api.GET("", func(c *gin.Context) {
			var req SearchEnrollmentRequest
			c.BindQuery(&req)
			enrollments, total, err := svc.SearchEnrollments(req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"items": enrollments, "total": total, "page": req.Page, "pageSize": req.PageSize})
		})

		api.GET("/:id", func(c *gin.Context) {
			enrollment, err := svc.GetEnrollment(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, enrollment)
		})

		api.POST("/:id/biometrics", func(c *gin.Context) {
			var req CaptureBiometricRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			enrollment, err := svc.CaptureBiometrics(c.Param("id"), req.BiometricHash, req.QualityScore)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, enrollment)
		})

		api.POST("/:id/verify", func(c *gin.Context) {
			var req VerifyEnrollmentRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			enrollment, err := svc.VerifyEnrollment(c.Param("id"), req.OtpVerified)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, enrollment)
		})

		api.POST("/:id/complete", func(c *gin.Context) {
			enrollment, err := svc.CompleteEnrollment(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, enrollment)
		})

		api.POST("/:id/reject", func(c *gin.Context) {
			var req RejectEnrollmentRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			enrollment, err := svc.RejectEnrollment(c.Param("id"), req.Reason)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, enrollment)
		})
	}

	srv := &http.Server{Addr: ":" + port, Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("enrollment service failed: %v", err)
		}
	}()
	logger.Info(context.Background(), "Enrollment service started on port "+port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	logger.Info(context.Background(), "shutting down enrollment-service...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func newUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
