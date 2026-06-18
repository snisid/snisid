package middleware_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/internal/platform/middleware"
)

// CreateIdentityDTO represents a request with strict validation rules.
type CreateIdentityDTO struct {
	FirstName string `json:"firstName" validate:"required,min=2,max=100"`
	Email     string `json:"email" validate:"required,email"`
	Age       int    `json:"age" validate:"required,gte=18"`
}

func TestValidateBody_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	var extractedDTO CreateIdentityDTO

	router.POST("/test", middleware.ValidateBody[CreateIdentityDTO](), func(c *gin.Context) {
		dto, ok := middleware.GetValidated[CreateIdentityDTO](c)
		if ok {
			extractedDTO = dto
		}
		c.Status(http.StatusOK)
	})

	validPayload := `{"firstName": "John", "email": "john@example.com", "age": 30}`
	req, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(validPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	if extractedDTO.FirstName != "John" {
		t.Errorf("Expected extracted FirstName to be 'John', got '%s'", extractedDTO.FirstName)
	}
}

func TestValidateBody_Failure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/test", middleware.ValidateBody[CreateIdentityDTO](), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Invalid payload: missing email, age < 18, firstName too short
	invalidPayload := `{"firstName": "J", "email": "invalid", "age": 16}`
	req, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(invalidPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	errObj := resp["error"].(map[string]interface{})
	if errObj["code"] != "invalid_argument" {
		t.Errorf("Expected code 'invalid_argument', got '%v'", errObj["code"])
	}
}
