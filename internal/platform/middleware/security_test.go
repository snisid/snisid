package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSecurityHeaders_Present(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(SecurityHeaders())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	headers := w.Header()

	if headers.Get("X-Frame-Options") != "DENY" {
		t.Errorf("X-Frame-Options = %s, want DENY", headers.Get("X-Frame-Options"))
	}
	if headers.Get("X-Content-Type-Options") != "nosniff" {
		t.Errorf("X-Content-Type-Options = %s, want nosniff", headers.Get("X-Content-Type-Options"))
	}
	if headers.Get("X-XSS-Protection") != "1; mode=block" {
		t.Errorf("X-XSS-Protection = %s, want 1; mode=block", headers.Get("X-XSS-Protection"))
	}
	if headers.Get("Strict-Transport-Security") != "max-age=31536000; includeSubDomains; preload" {
		t.Errorf("HSTS = %s, want max-age=31536000...", headers.Get("Strict-Transport-Security"))
	}
	if headers.Get("Referrer-Policy") != "no-referrer" {
		t.Errorf("Referrer-Policy = %s, want no-referrer", headers.Get("Referrer-Policy"))
	}
}

func TestSecurityHeaders_CSP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(SecurityHeaders())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	csp := w.Header().Get("Content-Security-Policy")
	if csp == "" {
		t.Error("Content-Security-Policy header should be set")
	}
	if csp != "default-src 'self'; script-src 'self'; object-src 'none';" {
		t.Errorf("CSP = %s, want proper CSP", csp)
	}
}

func TestSecurityHeaders_AllHeadersSet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(SecurityHeaders())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	expectedHeaders := []string{
		"X-Frame-Options",
		"X-Content-Type-Options",
		"X-XSS-Protection",
		"Content-Security-Policy",
		"Strict-Transport-Security",
		"Referrer-Policy",
	}

	for _, h := range expectedHeaders {
		if w.Header().Get(h) == "" {
			t.Errorf("Header %s should not be empty", h)
		}
	}
}

func TestSecurityHeaders_DoesNotBreakResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(SecurityHeaders())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}
	if w.Body.String() == "" {
		t.Error("Response body should not be empty")
	}
}
