package oidc

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler serves OIDC endpoints.
type Handler struct {
	cfg Config
}

// Config holds dependencies for the OIDC handler.
type Config struct {
	Issuer        string
	JWTSecret     string
	ClientManager ClientManager
	ConsentEngine ConsentEngine
	SessionStore  SessionStore
}

// ClientManager defines the interface for OAuth2 client operations.
type ClientManager interface {
	ValidateClient(clientID, clientSecret string) bool
	GetRedirectURI(clientID string) string
	GetClientScopes(clientID string) []string
}

// ConsentEngine defines the interface for user consent operations.
type ConsentEngine interface {
	HasConsent(userID, clientID, scope string) bool
	RecordConsent(userID, clientID, scope string)
}

// SessionStore defines the interface for SSO session operations.
type SessionStore interface {
	CreateSession(userID, clientID string) string
	GetSession(sessionID string) *SessionData
	DeleteSession(sessionID string)
}

// SessionData represents an authenticated user session.
type SessionData struct {
	UserID   string
	ClientID string
	Subject  string
	Scopes   []string
}

// NewHandler creates a new OIDC handler.
func NewHandler(cfg Config) *Handler {
	return &Handler{cfg: cfg}
}

// Discovery serves the OpenID Connect discovery document.
// GET /.well-known/openid-configuration
func (h *Handler) Discovery(c *gin.Context) {
	iss := h.cfg.Issuer
	c.JSON(http.StatusOK, gin.H{
		"issuer":                                iss,
		"authorization_endpoint":                fmt.Sprintf("%s/oidc/authorize", iss),
		"token_endpoint":                        fmt.Sprintf("%s/oidc/token", iss),
		"userinfo_endpoint":                     fmt.Sprintf("%s/oidc/userinfo", iss),
		"introspection_endpoint":                fmt.Sprintf("%s/oidc/introspect", iss),
		"revocation_endpoint":                   fmt.Sprintf("%s/oidc/revoke", iss),
		"jwks_uri":                              fmt.Sprintf("%s/.well-known/jwks", iss),
		"scopes_supported":                      []string{"openid", "profile", "email", "address", "phone", "snisid:identity", "snisid:biometric", "snisid:document", "offline_access"},
		"response_types_supported":              []string{"code", "id_token", "token"},
		"grant_types_supported":                 []string{"authorization_code", "refresh_token", "client_credentials"},
		"acr_values_supported":                  []string{"1", "2"},
		"subject_types_supported":               []string{"public", "pairwise"},
		"id_token_signing_alg_values_supported": []string{"RS256", "ES256"},
		"token_endpoint_auth_methods_supported": []string{"client_secret_basic", "client_secret_post"},
		"claims_supported":                      []string{"sub", "iss", "aud", "exp", "iat", "name", "given_name", "family_name", "email", "phone_number", "address", "snisid_national_id", "snisid_status"},
		"code_challenge_methods_supported":      []string{"S256"},
	})
}

func respondError(c *gin.Context, status int, errType, desc string) {
	c.JSON(status, gin.H{
		"error":             errType,
		"error_description": desc,
	})
}
