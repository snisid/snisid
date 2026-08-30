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
	"github.com/snisid/platform/services/fir/internal/api/rest"
	"github.com/snisid/platform/services/fir/internal/service"
)

func main() {
	recordSvc := service.NewRecordService()
	chargeSvc := service.NewChargeService()
	aliasSvc := service.NewAliasService()

	recordHandler := rest.NewRecordHandler(recordSvc, chargeSvc)
	searchHandler := rest.NewSearchHandler(recordSvc, chargeSvc)
	aliasHandler := rest.NewAliasHandler(aliasSvc, recordSvc)

	router := gin.Default()

	v1 := router.Group("/api/v1/fir")
	{
		v1.POST("/records", recordHandler.CreateRecord)
		v1.GET("/records", recordHandler.ListRecords)
		v1.GET("/records/:id", recordHandler.GetRecord)
		v1.GET("/records/person/:person_id", recordHandler.GetRecordByPerson)
		v1.POST("/records/:id/arrests", recordHandler.AddArrest)
		v1.POST("/records/:id/convictions", recordHandler.AddConviction)
		v1.POST("/records/:id/expunge", recordHandler.ExpungeRecord)

		v1.GET("/search", searchHandler.Search)
		v1.GET("/search/person", searchHandler.SearchByPerson)

		v1.POST("/records/:id/aliases", aliasHandler.AddAlias)
		v1.GET("/records/:id/aliases", aliasHandler.ListAliases)
		v1.DELETE("/records/:id/aliases/:alias_id", aliasHandler.RemoveAlias)
	}

	srv := &http.Server{
		Addr:    ":8093",
		Handler: router,
	}

	go func() {
		log.Printf("FIR-HT listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down FIR-HT server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %s", err)
	}
	log.Println("FIR-HT server exited")
}
