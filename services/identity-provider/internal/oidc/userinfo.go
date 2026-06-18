package oidc

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// UserInfo serves the OIDC UserInfo endpoint.
// GET /oidc/userinfo (requires Bearer token)
func (h *Handler) UserInfo(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		respondError(c, http.StatusUnauthorized, "invalid_token", "Missing or invalid access token")
		return
	}

	tokenStr := strings.TrimPrefix(auth, "Bearer ")
	claims, err := h.validateToken(tokenStr)
	if err != nil {
		respondError(c, http.StatusUnauthorized, "invalid_token", "Token validation failed: "+err.Error())
		return
	}

	sub, _ := claims.GetSubject()
	c.JSON(http.StatusOK, gin.H{
		"sub":                sub,
		"snisid_national_id": sub,
		"name":               "Citizen " + sub,
		"updated_at":         0,
	})
}
