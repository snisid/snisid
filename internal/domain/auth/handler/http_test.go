package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/internal/domain/auth/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockAuthService struct {
	registerFn func(ctx context.Context, username, password, roles string) error
	loginFn    func(ctx context.Context, username, password, clientIP, device string) (*usecase.TokenPair, error)
	refreshFn  func(ctx context.Context, sessionID, oldRefreshToken, clientIP string) (*usecase.TokenPair, error)
	logoutFn   func(ctx context.Context, sessionID string) error
}

func (m *mockAuthService) Register(ctx context.Context, username, password, roles string) error {
	if m.registerFn != nil {
		return m.registerFn(ctx, username, password, roles)
	}
	return nil
}

func (m *mockAuthService) Login(ctx context.Context, username, password, clientIP, device string) (*usecase.TokenPair, error) {
	if m.loginFn != nil {
		return m.loginFn(ctx, username, password, clientIP, device)
	}
	return &usecase.TokenPair{AccessToken: "mock-token", RefreshToken: "mock-refresh", SessionID: "mock-session"}, nil
}

func (m *mockAuthService) Refresh(ctx context.Context, sessionID, oldRefreshToken, clientIP string) (*usecase.TokenPair, error) {
	if m.refreshFn != nil {
		return m.refreshFn(ctx, sessionID, oldRefreshToken, clientIP)
	}
	return &usecase.TokenPair{AccessToken: "new-token", RefreshToken: "new-refresh", SessionID: "new-session"}, nil
}

func (m *mockAuthService) Logout(ctx context.Context, sessionID string) error {
	if m.logoutFn != nil {
		return m.logoutFn(ctx, sessionID)
	}
	return nil
}

func setupRouter(svc usecase.AuthService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewHttpHandler(svc)
	h.RegisterRoutes(r.Group("/v1/auth"))
	return r
}

func TestRegister_Success(t *testing.T) {
	svc := &mockAuthService{}
	router := setupRouter(svc)

	body := `{"username":"jdoe","password":"StrongP@ss1","roles":"user"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "registered", resp["status"])
}

func TestRegister_InvalidBody(t *testing.T) {
	svc := &mockAuthService{}
	router := setupRouter(svc)

	body := `{invalid json}`
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegister_ServiceError(t *testing.T) {
	svc := &mockAuthService{
		registerFn: func(ctx context.Context, username, password, roles string) error {
			return errors.New("db unavailable")
		},
	}
	router := setupRouter(svc)

	body := `{"username":"jdoe","password":"StrongP@ss1","roles":"user"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestLogin_Success(t *testing.T) {
	svc := &mockAuthService{
		loginFn: func(ctx context.Context, username, password, clientIP, device string) (*usecase.TokenPair, error) {
			return &usecase.TokenPair{
				AccessToken:  "access-token-123",
				RefreshToken: "refresh-token-456",
				SessionID:    "session-789",
			}, nil
		},
	}
	router := setupRouter(svc)

	body := `{"username":"jdoe","password":"StrongP@ss1","device":"chrome-120"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp usecase.TokenPair
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "access-token-123", resp.AccessToken)
	assert.Equal(t, "refresh-token-456", resp.RefreshToken)
}

func TestLogin_AccountLocked(t *testing.T) {
	svc := &mockAuthService{
		loginFn: func(ctx context.Context, username, password, clientIP, device string) (*usecase.TokenPair, error) {
			return nil, usecase.ErrAccountLocked
		},
	}
	router := setupRouter(svc)

	body := `{"username":"locked","password":"AnyP@ss1","device":"chrome"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}

func TestLogin_WrongCredentials(t *testing.T) {
	svc := &mockAuthService{
		loginFn: func(ctx context.Context, username, password, clientIP, device string) (*usecase.TokenPair, error) {
			return nil, usecase.ErrInvalidCredentials
		},
	}
	router := setupRouter(svc)

	body := `{"username":"jdoe","password":"WrongP@ss1","device":"firefox"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRefresh_Success(t *testing.T) {
	svc := &mockAuthService{}
	router := setupRouter(svc)

	body := `{"sessionId":"sess-1","refreshToken":"refresh-token-abc"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/refresh", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp usecase.TokenPair
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
}

func TestRefresh_Invalid(t *testing.T) {
	svc := &mockAuthService{
		refreshFn: func(ctx context.Context, sessionID, oldRefreshToken, clientIP string) (*usecase.TokenPair, error) {
			return nil, usecase.ErrInvalidSession
		},
	}
	router := setupRouter(svc)

	body := `{"sessionId":"sess-expired","refreshToken":"bad-token"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/refresh", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogout_Success(t *testing.T) {
	svc := &mockAuthService{}
	router := setupRouter(svc)

	body := `{"sessionId":"sess-1"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/logout", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLogout_InvalidBody(t *testing.T) {
	svc := &mockAuthService{}
	router := setupRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/v1/auth/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
