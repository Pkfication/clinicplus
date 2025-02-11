package main

import (
    "clinicplus/models"
    "clinicplus/routes"
    "clinicplus/utils"
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    "golang.org/x/crypto/bcrypt"
)

var db *gorm.DB

func setupCORS(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        handler.ServeHTTP(w, r)
    })
}

func initDatabase() *gorm.DB {
    // Load environment variables
    utils.LoadEnv()

    // Get database connection string from environment
    dbURL := utils.GetDatabaseURL()

    // Open database connection
    db, err := gorm.Open("postgres", dbURL)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Auto migrate models
    db.AutoMigrate(&models.Employee{}, &models.User{})

    // Optional: Create a default admin user if not exists
    createDefaultAdminUser(db)

    return db
}

func createDefaultAdminUser(db *gorm.DB) {
    var user models.User
    // Check if admin user already exists
    if err := db.Where("username = ?", "admin").First(&user).Error; err != nil {
        // Admin user doesn't exist, create one
        hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("adminpassword"), bcrypt.DefaultCost)
        
        adminUser := models.User{
            Username:     "admin",
            PasswordHash: string(hashedPassword),
            Role:         "Admin",
        }

        if err := db.Create(&adminUser).Error; err != nil {
            log.Printf("Failed to create default admin user: %v", err)
        } else {
            log.Println("Default admin user created")
        }
    }
}

func main() {
    // Initialize database
    db := initDatabase()
    defer db.Close()

    // Get JWT secret from environment
    jwtSecret := utils.GetJWTSecret()

    // Create router
    router := mux.NewRouter()

    // Register routes and pass JWT secret
    routes.SetDB(db)
    routes.SetJWTKey(jwtSecret)
    routes.RegisterRoutes(router)

	   // Log registered routes
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err == nil {
			methods, _ := route.GetMethods()
			log.Printf("Registered Route: %v %v", methods, path)
		}
		return nil
    })

    if err != nil {
        log.Printf("Error walking routes: %v", err)
    }

    // Get port from environment
    port := utils.GetServerPort()

    // Setup CORS middleware
    handler := setupCORS(router)

    // Start server
    serverAddr := fmt.Sprintf(":%s", port)
    log.Printf("Server starting on port %s", port)
    
    // Start the server
    log.Fatal(http.ListenAndServe(serverAddr, handler))
}