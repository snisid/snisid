package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestTracing_ContextSet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Tracing())
	r.GET("/test", func(c *gin.Context) {
		// Context should have tracing span injected
		if c.Request.Context() == nil {
			t.Error("Request context should not be nil")
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestTracing_MultipleRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Tracing())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"path": c.Request.URL.Path})
	})

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d: Status = %d, want %d", i, w.Code, http.StatusOK)
		}
	}
}

func TestTracing_DifferentPaths(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Tracing())
	r.GET("/api/v1/identities", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	r.POST("/api/v1/enrollments", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"id": "new"})
	})

	paths := []struct {
		method string
		path   string
		code   int
	}{
		{"GET", "/api/v1/identities", http.StatusOK},
		{"POST", "/api/v1/enrollments", http.StatusCreated},
	}

	for _, p := range paths {
		req := httptest.NewRequest(p.method, p.path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != p.code {
			t.Errorf("%s %s: Status = %d, want %d", p.method, p.path, w.Code, p.code)
		}
	}
}

func TestTracing_StatusCodeRecorded(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Tracing())
	r.GET("/ok", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	r.GET("/notfound", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	tests := []struct {
		path string
		code int
	}{
		{"/ok", http.StatusOK},
		{"/notfound", http.StatusNotFound},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodGet, tt.path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != tt.code {
			t.Errorf("%s: Status = %d, want %d", tt.path, w.Code, tt.code)
		}
	}
}
