package routes

import (
    "github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {

    // Health Check Routes
    r.HandleFunc("/health", HealthCheck).Methods("GET")

    // Authentication Routes
    r.HandleFunc("/login", Login).Methods("POST")
    r.HandleFunc("/logout", Logout).Methods("POST")

    // Employee Management Routes
    employeeRouter := r.PathPrefix("/employees").Subrouter()
    employeeRouter.HandleFunc("", GetEmployees).Methods("GET")
    employeeRouter.HandleFunc("", CreateEmployee).Methods("POST")
    employeeRouter.HandleFunc("/{id}", GetEmployee).Methods("GET")
    employeeRouter.HandleFunc("/{id}", UpdateEmployee).Methods("PUT")
    employeeRouter.HandleFunc("/{id}", DeleteEmployee).Methods("DELETE")
}
