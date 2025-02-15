// internal/iam/service.go
package iam

import (
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(username, password string) (string, string, time.Time, error)
	Logout(token string) error
}

type authService struct {
	db     *gorm.DB
	jwtKey []byte
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func NewAuthService(db *gorm.DB, jwtKey []byte) AuthService {
	return &authService{db: db, jwtKey: jwtKey}
}

func (s *authService) Login(username, password string) (string, string, time.Time, error) {
	var user User

	// Find user by username
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		log.Printf("Invalid login attempt for user: %s\n", username)
		return "", "", time.Time{}, errors.New("invalid username or password")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		log.Printf("Invalid password for user: %s\n", username)
		return "", "", time.Time{}, errors.New("invalid username or password")
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

	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		log.Println("Error signing token:", err)
		return "", "", time.Time{}, errors.New("error generating authentication token")
	}

	log.Printf("User %s logged in successfully\n", user.Username)
	return tokenString, user.Username, expirationTime, nil
}

func (s *authService) Logout(token string) error {
	// In a stateless JWT system, true logout happens on the client-side
	// by discarding the token. This function is a placeholder for future enhancements.
	log.Printf("Logout attempt with token: %s\n", token)
	return nil
}
