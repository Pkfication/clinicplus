package server

import (
	"clinicplus/internal/shared/db"
	"clinicplus/internal/shared/observability"
	"clinicplus/internal/shared/routes"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

func NewServer() *mux.Router {
	db.InitDatabase()

	r := mux.NewRouter()
	
	// Add observability middleware
	r.Use(otelmux.Middleware("clinicplus-api"))
	r.Use(observability.MetricsMiddleware())
	
	// Register observability routes
	r.Handle("/metrics", observability.MetricsHandler(observability.InitMetrics())).Methods("GET")
	r.PathPrefix("/debug/pprof/").Handler(observability.PprofHandler())
	
	// Register application routes
	routes.RegisterRoutes(r, db.DB)

	return r
}
