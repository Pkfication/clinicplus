package utils

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

// LoadEnv loads environment variables from .env file
func LoadEnv() {
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found. Using system environment.")
    }
}

// GetEnvString retrieves a string environment variable with a default value
func GetEnvString(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}

// GetJWTSecret retrieves the JWT secret from environment
func GetJWTSecret() []byte {
    secret := GetEnvString("JWT_SECRET", "")
    if secret == "" {
        log.Fatal("JWT_SECRET not set in environment")
    }
    return []byte(secret)
}

// GetDatabaseURL retrieves the database connection string
func GetDatabaseURL() string {
    dbURL := GetEnvString("DATABASE_URL", "")
    if dbURL == "" {
        log.Fatal("DATABASE_URL not set in environment")
    }
    return dbURL
}

// GetServerPort retrieves the server port
func GetServerPort() string {
    return GetEnvString("PORT", "8080")
}