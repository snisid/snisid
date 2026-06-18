package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/snisid/platform/services/identity-provider/internal/client"
	"github.com/snisid/platform/services/identity-provider/internal/consent"
	"github.com/snisid/platform/services/identity-provider/internal/oidc"
	"github.com/snisid/platform/services/identity-provider/internal/session"
)

func main() {
	jwtSecret := getEnv("JWT_SECRET", "dev-secret-identity-provider")
	port := getEnv("PORT", "8090")
	issuer := getEnv("ISSUER", "http://localhost:8090")

	clientManager := client.NewManager()
	consentEngine := consent.NewEngine()
	sessionStore := session.NewStore()

	cfg := oidc.Config{
		Issuer:        issuer,
		JWTSecret:     jwtSecret,
		ClientManager: clientManager,
		ConsentEngine: consentEngine,
		SessionStore:  sessionStore,
	}
	handler := oidc.NewHandler(cfg)

	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/.well-known/openid-configuration", handler.Discovery)
	r.GET("/.well-known/jwks", handler.JWKS)

	oidcGroup := r.Group("/oidc")
	{
		oidcGroup.GET("/authorize", handler.Authorize)
		oidcGroup.POST("/token", handler.Token)
		oidcGroup.GET("/userinfo", handler.UserInfo)
		oidcGroup.POST("/introspect", handler.Introspect)
		oidcGroup.POST("/revoke", handler.Revoke)
		oidcGroup.GET("/session", handler.SessionInfo)
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("identity-provider started on port %s (issuer: %s)", port, issuer)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run identity-provider: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down identity-provider...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server forced to shutdown: %v", err)
	}
}

func getEnv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
