package routes

import (
	"clinicplus/internal/employee"
	"clinicplus/internal/iam"
	"clinicplus/internal/shared/config"
	"log"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func RegisterRoutes(r *mux.Router, db *gorm.DB) {

	db.AutoMigrate(&iam.User{}, &employee.Employee{}, &employee.Attendance{}, &employee.Shift{}, &employee.EmployeeShift{})

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
	shiftRouter := r.PathPrefix("/shifts").Subrouter()
	employeeRouter.HandleFunc("", employeeHandler.GetEmployees).Methods("GET")
	employeeRouter.HandleFunc("/search", employeeHandler.SearchEmployees).Methods("GET")
	employeeRouter.HandleFunc("", employeeHandler.CreateEmployee).Methods("POST")
	employeeRouter.HandleFunc("/{id}", employeeHandler.GetEmployee).Methods("GET")
	employeeRouter.HandleFunc("/{id}", employeeHandler.UpdateEmployee).Methods("PUT")
	employeeRouter.HandleFunc("/{id}", employeeHandler.DeleteEmployee).Methods("DELETE")
	employeeRouter.HandleFunc("/{id}/clockin", employeeHandler.ClockInEmployee).Methods("POST")
	employeeRouter.HandleFunc("/{id}/clockout", employeeHandler.ClockOutEmployee).Methods("POST")
	employeeRouter.HandleFunc("/{id}/assign_shift", employeeHandler.AssignShift).Methods("POST")

	// Shift Management Routes
	shiftRouter.HandleFunc("", employeeHandler.CreateShift).Methods("POST")
	shiftRouter.HandleFunc("/{id}", employeeHandler.GetShift).Methods("GET")
	shiftRouter.HandleFunc("/{id}", employeeHandler.UpdateShift).Methods("PUT")
	shiftRouter.HandleFunc("/{id}", employeeHandler.DeleteShift).Methods("DELETE")

	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, err1 := route.GetPathTemplate()
		met, err2 := route.GetMethods()
		log.Println(tpl, err1, met, err2)
		return nil
	})

}
