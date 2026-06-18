package middleware

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuth_MissingBearer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Auth("secret"))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuth_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Auth("secret"))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuth_RoleAllowed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Auth("secret", "admin"))
	r.GET("/test", func(c *gin.Context) {
		claims, _ := c.Get("claims")
		if claims == nil {
			t.Error("claims should be set in context")
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// Generate a valid token for role "admin"
	token := generateTestToken("usr-001", "admin", "secret")
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestAuth_RoleDenied(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Auth("secret", "admin"))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	token := generateTestToken("usr-001", "viewer", "secret")
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestAuth_NoRoleCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Auth("secret")) // No roles specified, any valid token passes
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	token := generateTestToken("usr-001", "any-role", "secret")
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}
}

// generateTestToken creates a simple JWT for testing (not using real crypto)
func generateTestToken(subject, role, secret string) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(
		`{"sub":"` + subject + `","role":"` + role + `","exp":9999999999}`))
	sig := base64.RawURLEncoding.EncodeToString([]byte(secret + "." + header + "." + payload))
	return strings.Join([]string{header, payload, sig}, ".")
}

func TestAuth_EmptySecret(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Auth(""))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	token := generateTestToken("usr-001", "admin", "")
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code == http.StatusUnauthorized {
		t.Log("Auth with empty secret correctly rejects tokens")
	}
}
