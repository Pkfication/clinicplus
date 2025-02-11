package routes

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "clinicplus/models"
    "github.com/gorilla/mux"
    "github.com/jinzhu/gorm"
    "github.com/stretchr/testify/assert"
    "golang.org/x/crypto/bcrypt"
)

func setupAuthTestDB() *gorm.DB {
    // Setup test database (similar to employee test setup)
    testDB, err := gorm.Open("sqlite3", ":memory:")
    if err != nil {
        panic("failed to connect to database")
    }
    testDB.AutoMigrate(&models.User{})
    
    // Create a test user
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
    testUser := models.User{
        Username:     "testuser",
        PasswordHash: string(hashedPassword),
    }
    testDB.Create(&testUser)

    SetDB(testDB)
    return testDB
}

func TestLogin(t *testing.T) {
    // Setup test database
    testDB := setupAuthTestDB()
    defer testDB.Close()

    // Setup router
    router := mux.NewRouter()
    RegisterRoutes(router)

    // Test cases
    testCases := []struct {
        name           string
        requestBody    LoginRequest
        expectedStatus int
        expectedError  string
    }{
        {
            name: "Successful Login",
            requestBody: LoginRequest{
                Username: "testuser",
                Password: "testpassword",
            },
            expectedStatus: http.StatusOK,
            expectedError:  "",
        },
        {
            name: "Invalid Username",
            requestBody: LoginRequest{
                Username: "nonexistentuser",
                Password: "testpassword",
            },
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Invalid username or password",
        },
        {
            name: "Invalid Password",
            requestBody: LoginRequest{
                Username: "testuser",
                Password: "wrongpassword",
            },
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Invalid username or password",
        },
        {
            name: "Empty Credentials",
            requestBody: LoginRequest{
                Username: "",
                Password: "",
            },
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Invalid username or password",
        },
    }

    // Run test cases
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Prepare request body
            jsonData, _ := json.Marshal(tc.requestBody)
            req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
            req.Header.Set("Content-Type", "application/json")

            // Create a response recorder
            rr := httptest.NewRecorder()
            router.ServeHTTP(rr, req)

            // Check the status code
            assert.Equal(t, tc.expectedStatus, rr.Code)

            // Parse the response
            var response StandardResponse
            err := json.Unmarshal(rr.Body.Bytes(), &response)
            assert.NoError(t, err)

            // Check error message for unsuccessful logins
            if tc.expectedStatus != http.StatusOK {
                assert.Equal(t, tc.expectedError, response.Error)
                assert.Nil(t, response.Data)
            } else {
                // For successful login, check data structure
                assert.Nil(t, response.Error)
                
                // Convert data to map for easier checking
                data, ok := response.Data.(map[string]interface{})
                assert.True(t, ok)
                
                assert.Contains(t, data, "token")
                assert.Contains(t, data, "username")
                assert.Equal(t, "testuser", data["username"])
            }

            // Check meta for successful login
            if tc.expectedStatus == http.StatusOK {
                assert.NotNil(t, response.Meta)
                meta, ok := response.Meta.(map[string]interface{})
                assert.True(t, ok)
                assert.Contains(t, meta, "expires_at")
            }
        })
    }
}