package routes

import (
	"clinicplus/internal/employee"
	"clinicplus/internal/iam"
	"clinicplus/internal/shared/config"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func RegisterRoutes(r *mux.Router, db *gorm.DB) {

	// Health Check Routes
	r.HandleFunc("/health", HealthCheck).Methods("GET")

	// IAM Routes
	iamService := iam.NewAuthService(db, config.GetJWTSecret())
	iamHandler := iam.NewAuthHandler(iamService)
	r.HandleFunc("/login", iamHandler.Login).Methods("POST")
	r.HandleFunc("/logout", iamHandler.Logout).Methods("POST")

	// Employee Management Routes
	employeeService := employee.NewEmployeeService(db)
	employeeHandler := employee.NewEmployeeHandler(employeeService)
	employeeRouter := r.PathPrefix("/employees").Subrouter()
	employeeRouter.HandleFunc("", employeeHandler.GetEmployees).Methods("GET")
	employeeRouter.HandleFunc("", employeeHandler.CreateEmployee).Methods("POST")
	employeeRouter.HandleFunc("/{id}", employeeHandler.GetEmployee).Methods("GET")
	employeeRouter.HandleFunc("/{id}", employeeHandler.UpdateEmployee).Methods("PUT")
	employeeRouter.HandleFunc("/{id}", employeeHandler.DeleteEmployee).Methods("DELETE")
}
