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
	"github.com/snisid/platform/services/afis/internal/api/rest"
	"github.com/snisid/platform/services/afis/internal/service"
)

func main() {
	quality := service.NewQualityService(60)
	enrollment := service.NewEnrollmentService(quality)
	search := service.NewSearchService()
	latentSvc := service.NewLatentService(search, quality)

	enrollHandler := rest.NewEnrollHandler(enrollment, search)
	searchHandler := rest.NewSearchHandler(search)
	latentHandler := rest.NewLatentHandler(latentSvc, search)
	qualityHandler := rest.NewQualityHandler(quality)

	router := gin.Default()

	v1 := router.Group("/api/v1/afis")
	{
		v1.POST("/enroll", enrollHandler.Enroll)
		v1.POST("/search/tenprint", searchHandler.SearchTenprint)
		v1.POST("/search/latent", searchHandler.SearchLatent)
		v1.GET("/subjects/:id", enrollHandler.GetSubject)
		v1.GET("/subjects/:id/history", enrollHandler.GetSubjectHistory)
		v1.POST("/latents", latentHandler.SubmitLatent)
		v1.GET("/latents/:id", latentHandler.GetLatent)
		v1.GET("/latents", latentHandler.ListLatents)
		v1.PATCH("/latents/:id/match", latentHandler.ConfirmMatch)
		v1.GET("/quality/check", qualityHandler.CheckQuality)
		v1.GET("/stats", enrollHandler.GetStats)
	}

	srv := &http.Server{
		Addr:    ":8091",
		Handler: router,
	}

	go func() {
		log.Printf("AFIS-HT listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down AFIS-HT server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %s", err)
	}
	log.Println("AFIS-HT server exited")
}
