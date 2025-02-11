package main

import (
    "clinicplus/models"
    "clinicplus/routes"
    "github.com/gorilla/mux"
    "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
    "github.com/joho/godotenv"
    "log"
    "net/http"
    "os"
)

var db *gorm.DB

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Connect to the database
    dsn := os.Getenv("DATABASE_URL")
    db, err = gorm.Open("postgres", dsn)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    // Migrate the models
    db.AutoMigrate(&models.Employee{}, &models.Attendance{}, &models.Leave{}, &models.Payroll{}, &models.Expense{}, &models.User{})

	// Set the database connection for routes
    routes.SetDB(db)

    r := mux.NewRouter()
    routes.RegisterRoutes(r)

    log.Println("Server started at :8080")
    http.ListenAndServe(":8080", r)
}