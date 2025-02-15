package server

import (
	"clinicplus/internal/shared/db"
	"clinicplus/internal/shared/routes"

	"github.com/gorilla/mux"
)

func NewServer() *mux.Router {
	db.InitDatabase()

	r := mux.NewRouter()
	routes.RegisterRoutes(r, db.DB)

	return r
}
