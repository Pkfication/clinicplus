package routes

import (
    "clinicplus/models"
    "encoding/json"
    "net/http"
    "github.com/dgrijalva/jwt-go"
    "golang.org/x/crypto/bcrypt"
    "time"
	"log"
)

type Claims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}

var jwtKey = []byte("your_secret_key") // Change this to a secure key in production

// LoginRequest represents the structure of the login request
// StandardResponse represents the standard API response structure
type StandardResponse struct {
    Data  interface{} `json:"data"`
    Error interface{} `json:"error"`
    Meta  interface{} `json:"meta"`
}

// LoginResponse represents the data returned on successful login
type LoginResponse struct {
    Token     string `json:"token"`
    Username  string `json:"username"`
}

// LoginRequest represents the structure of the login request
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

// sendJSONResponse is a helper function to send standardized JSON responses
func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}, err interface{}, meta interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    response := StandardResponse{
        Data:  data,
        Error: err,
        Meta:  meta,
    }
    
    json.NewEncoder(w).Encode(response)
}

// Login function
func Login(w http.ResponseWriter, r *http.Request) {
    var user models.User
    var loginRequest LoginRequest

    // Decode request body
    if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
        log.Println("Error decoding request body:", err)
        sendJSONResponse(w, http.StatusBadRequest, nil, "Invalid request body", nil)
        return
    }

    // Find user by username
    if err := db.Where("username = ?", loginRequest.Username).First(&user).Error; err != nil {
        log.Printf("Invalid login attempt for user: %s\n", loginRequest.Username)
        sendJSONResponse(w, http.StatusUnauthorized, nil, "Invalid username or password", nil)
        return
    }

    // Compare password
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginRequest.Password)); err != nil {
        log.Printf("Invalid password for user: %s\n", loginRequest.Username)
        sendJSONResponse(w, http.StatusUnauthorized, nil, "Invalid username or password", nil)
        return
    }

    // Create JWT token
    expirationTime := time.Now().Add(5 * time.Minute)
    claims := &Claims{
        Username: user.Username,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        log.Println("Error signing token:", err)
        sendJSONResponse(w, http.StatusInternalServerError, nil, "Error generating authentication token", nil)
        return
    }

    // Prepare login response
    loginResponse := LoginResponse{
        Token:    tokenString,
        Username: user.Username,
    }

    log.Printf("User %s logged in successfully\n", user.Username)

    // Send standardized response
    sendJSONResponse(w, http.StatusOK, loginResponse, nil, map[string]interface{}{
        "expires_at": expirationTime.Format(time.RFC3339),
    })
}

// Logout function (optional)
func Logout(w http.ResponseWriter, r *http.Request) {
    // In a stateless JWT system, true logout happens on the client-side
    // by discarding the token. However, we'll provide a standardized response.

    // Extract the token from the Authorization header (optional, for logging)
    authHeader := r.Header.Get("Authorization")
    
    log.Printf("Logout attempt with token: %s\n", authHeader)

    // Prepare logout response
    logoutResponse := map[string]string{
        "message": "Successfully logged out",
    }

    // Send standardized response
    sendJSONResponse(w, http.StatusOK, logoutResponse, nil, map[string]interface{}{
        "timestamp": time.Now().Format(time.RFC3339),
    })
}