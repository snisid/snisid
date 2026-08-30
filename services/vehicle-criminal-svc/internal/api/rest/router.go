package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/snisid/vehicle-criminal-svc/internal/service"
)

func NewRouter(
	alertSvc *service.CriminalAlertService,
	plateSvc *service.StolenPlateService,
	intelSvc *service.VehicleIntelService,
	interpolSvc *service.InterpolSyncService,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(AuditMiddleware)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1/sivc", func(r chi.Router) {
		r.Post("/alerts", CreateAlertHandler(alertSvc))
		r.Get("/alerts", ListAlertsHandler(alertSvc))
		r.Get("/alerts/{id}", GetAlertHandler(alertSvc))
		r.Patch("/alerts/{id}/status", UpdateAlertStatusHandler(alertSvc))
		r.Post("/alerts/{id}/sighting", CreateSightingHandler(alertSvc))

		r.Get("/check/plate/{plate}", CheckPlateHandler(alertSvc))
		r.Get("/check/vin/{vin}", CheckVINHandler(alertSvc))

		r.Post("/stolen-plates", DeclareStolenPlateHandler(plateSvc))
		r.Get("/stolen-plates/{plate}", CheckStolenPlateHandler(plateSvc))
		r.Patch("/stolen-plates/{id}/recovered", MarkPlateRecoveredHandler(plateSvc))

		r.Get("/search", SearchHandler(alertSvc))

		r.Post("/intel-reports", CreateIntelReportHandler(intelSvc))
		r.Get("/intel-reports/{id}", GetIntelReportHandler(intelSvc))

		r.Get("/interpol/sync-status", InterpolSyncStatusHandler(interpolSvc))
	})

	return r
}
