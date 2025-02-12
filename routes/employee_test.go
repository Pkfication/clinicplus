package routes

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
	"strconv"  
    "testing"

    "clinicplus/models"
    "clinicplus/utils"
    "github.com/gorilla/mux"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "github.com/stretchr/testify/assert"
)

var testDB *gorm.DB

func setupTestDB() {
    // Initialize test database connection
    var err error
    testDB, err = gorm.Open("sqlite3", ":memory:")
    if err != nil {
        panic("failed to connect database")
    }

    // Migrate the schema
    testDB.AutoMigrate(&models.Employee{}, &models.User{})

    // Set the database for routes
    SetDB(testDB)
}

func TestCreateEmployee(t *testing.T) {
    // Setup test database
    setupTestDB()
    defer testDB.Close()

    // Prepare test employee creation request
    createRequest := struct {
        models.Employee
        Username string `json:"username"`
        Password string `json:"password"`
    }{
        Employee: models.Employee{
            Name:  "John Doe",
            Email: "john.doe@company.com",
            Role:  "Employee",
        },
        Username: "johndoe",
        Password: "securepassword123",
    }

    // Convert request to JSON
    jsonRequest, err := json.Marshal(createRequest)
    assert.NoError(t, err)

    // Create request
    req, err := http.NewRequest("POST", "/employees", bytes.NewBuffer(jsonRequest))
    assert.NoError(t, err)

    // Create response recorder
    rr := httptest.NewRecorder()

    // Call the handler
    CreateEmployee(rr, req)

    // Check the status code
    assert.Equal(t, http.StatusCreated, rr.Code)

    // Parse the response
    var response utils.StandardResponse
    err = json.Unmarshal(rr.Body.Bytes(), &response)
    assert.NoError(t, err)

    // Convert response data to employee
    var createdEmployee models.Employee
    data, _ := json.Marshal(response.Data)
    err = json.Unmarshal(data, &createdEmployee)
    assert.NoError(t, err)

    // Verify employee details
    assert.Equal(t, "John Doe", createdEmployee.Name)
    assert.Equal(t, "john.doe@company.com", createdEmployee.Email)

    // Verify the employee was actually saved in the database
    var dbEmployee models.Employee
    err = testDB.Where("email = ?", "john.doe@company.com").First(&dbEmployee).Error
    assert.NoError(t, err)
    assert.Equal(t, "John Doe", dbEmployee.Name)

    // Verify user was created
    var dbUser models.User
    err = testDB.Where("username = ?", "johndoe").First(&dbUser).Error
    assert.NoError(t, err)
    assert.Equal(t, "johndoe", dbUser.Username)
    assert.Equal(t, dbEmployee.ID, dbUser.EmployeeID)
}

// Add additional test cases for validation errors
func TestCreateEmployeeValidation(t *testing.T) {
    // Setup test database
    setupTestDB()
    defer testDB.Close()

    testCases := []struct {
        name           string
        requestBody    interface{}
        expectedStatus int
        expectedError  string
    }{
        {
            name: "Missing Name",
            requestBody: struct {
                models.Employee
                Username string `json:"username"`
                Password string `json:"password"`
            }{
                Employee: models.Employee{
                    Email: "john.doe@company.com",
                    Role:  "Employee",
                },
                Username: "johndoe",
                Password: "securepassword123",
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Name is required",
        },
        {
            name: "Missing Email",
            requestBody: struct {
                models.Employee
                Username string `json:"username"`
                Password string `json:"password"`
            }{
                Employee: models.Employee{
                    Name: "John Doe",
                    Role: "Employee",
                },
                Username: "johndoe",
                Password: "securepassword123",
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Email must be explicitly provided",
        },
        {
            name: "Missing Username",
            requestBody: struct {
                models.Employee
                Password string `json:"password"`
            }{
                Employee: models.Employee{
                    Name:  "John Doe",
                    Email: "john.doe@company.com",
                    Role:  "Employee",
                },
                Password: "securepassword123",
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Username must be explicitly provided",
        },
        {
            name: "Missing Password",
            requestBody: struct {
                models.Employee
                Username string `json:"username"`
            }{
                Employee: models.Employee{
                    Name:  "John Doe",
                    Email: "john.doe@company.com",
                    Role:  "Employee",
                },
                Username: "johndoe",
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Password must be explicitly provided",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Convert request to JSON
            jsonRequest, err := json.Marshal(tc.requestBody)
            assert.NoError(t, err)

            // Create request
            req, err := http.NewRequest("POST", "/employees", bytes.NewBuffer(jsonRequest))
            assert.NoError(t, err)

            // Create response recorder
            rr := httptest.NewRecorder()

            // Call the handler
            CreateEmployee(rr, req)

            // Check the status code
            assert.Equal(t, tc.expectedStatus, rr.Code)

            // Parse the response
            var response utils.StandardResponse
            err = json.Unmarshal(rr.Body.Bytes(), &response)
            assert.NoError(t, err)

            // Verify error message
            assert.Contains(t, rr.Body.String(), tc.expectedError)
        })
    }
}

func TestGetEmployee(t *testing.T) {
    // Setup test database
    setupTestDB()
    defer testDB.Close()

    // Create a test employee
    employee := models.Employee{
        Name:  "John Doe", 
        Email: "john.doe@example.com",
        Role:  "Employee",
    }
    testDB.Create(&employee)

    // Create request
    // Convert ID to string correctly
    req, err := http.NewRequest("GET", "/employees/"+strconv.Itoa(int(employee.ID)), nil)
    assert.NoError(t, err)

    // Create response recorder
    rr := httptest.NewRecorder()

    // Set URL vars for mux
    req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(employee.ID))})

    // Call the handler
    GetEmployee(rr, req)

    // Check the status code
    assert.Equal(t, http.StatusOK, rr.Code)

    // Print raw response body for debugging
    rawBody := rr.Body.Bytes()
    t.Logf("Raw Response Body: %s", string(rawBody))

    // Parse the response
    var response utils.StandardResponse
    err = json.Unmarshal(rawBody, &response)
    assert.NoError(t, err, "Should be able to unmarshal StandardResponse")

    // Convert response data to employee
    var retrievedEmployee models.Employee
    
    // Convert Data to JSON first, then unmarshal
    dataBytes, err := json.Marshal(response.Data)
    assert.NoError(t, err)
    
    err = json.Unmarshal(dataBytes, &retrievedEmployee)
    assert.NoError(t, err, "Should be able to unmarshal employee")

    // Verify employee details
    assert.Equal(t, "John Doe", retrievedEmployee.Name)
    assert.Equal(t, "john.doe@example.com", retrievedEmployee.Email)
}

func TestUpdateEmployee(t *testing.T) {
    // Setup test database
    setupTestDB()
    defer testDB.Close()

    // Create a test employee
    employee := models.Employee{
        Name:  "John Doe", 
        Email: "john.doe@example.com",
        Role:  "Employee",
    }
    testDB.Create(&employee)

    // Prepare updated employee data
    updatedEmployee := models.Employee{
        Name:  "John Smith",
        Email: "john.smith@example.com",
        Role:  "Manager",
    }

    // Convert updated employee to JSON
    jsonEmployee, err := json.Marshal(updatedEmployee)
    assert.NoError(t, err)

    // Create request
    // Convert ID to string correctly
    req, err := http.NewRequest("PUT", "/employees/"+strconv.Itoa(int(employee.ID)), bytes.NewBuffer(jsonEmployee))
    assert.NoError(t, err)

    // Set URL vars for mux
    req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(employee.ID))})

    // Create response recorder
    rr := httptest.NewRecorder()

    // Call the handler
    UpdateEmployee(rr, req)

    // Check the status code
    assert.Equal(t, http.StatusOK, rr.Code)

    // Print raw response body for debugging
    rawBody := rr.Body.Bytes()

    // Parse the response
    var response utils.StandardResponse
    err = json.Unmarshal(rawBody, &response)
    assert.NoError(t, err, "Should be able to unmarshal StandardResponse")

    // Convert response data to employee
    var returnedEmployee models.Employee
    
    // Convert Data to JSON first, then unmarshal
    dataBytes, err := json.Marshal(response.Data)
    assert.NoError(t, err)
    
    err = json.Unmarshal(dataBytes, &returnedEmployee)
    assert.NoError(t, err, "Should be able to unmarshal employee")

    // Verify updated employee details
    assert.Equal(t, "John Smith", returnedEmployee.Name)
    assert.Equal(t, "john.smith@example.com", returnedEmployee.Email)

    // Verify the employee was actually updated in the database
    var dbEmployee models.Employee
    err = testDB.First(&dbEmployee, employee.ID).Error
    assert.NoError(t, err)
    assert.Equal(t, "John Smith", dbEmployee.Name)
}

func TestDeleteEmployee(t *testing.T) {
    // Setup test database
    setupTestDB()
    defer testDB.Close()

    // Create a test employee
    employee := models.Employee{
        Name:  "John Doe", 
        Email: "john.doe@example.com",
        Role:  "Employee",
    }
    testDB.Create(&employee)

    // Create request
    // Convert ID to string correctly
    req, err := http.NewRequest("DELETE", "/employees/"+strconv.Itoa(int(employee.ID)), nil)
    assert.NoError(t, err)

    // Set URL vars for mux
    req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(employee.ID))})

    // Create response recorder
    rr := httptest.NewRecorder()

    // Call the handler
    DeleteEmployee(rr, req)

    // Check the status code
    assert.Equal(t, http.StatusOK, rr.Code)

    // Print raw response body for debugging
    rawBody := rr.Body.Bytes()

    // Parse the response
    var response utils.StandardResponse
    err = json.Unmarshal(rawBody, &response)
    assert.NoError(t, err, "Should be able to unmarshal StandardResponse")

    // Verify the employee was actually deleted from the database
    var dbEmployee models.Employee
    err = testDB.First(&dbEmployee, employee.ID).Error
    assert.Error(t, err, "Expect a 'record not found' error after deletion")
}

func TestGetEmployees(t *testing.T) {
    // Setup test database
    setupTestDB()
    defer testDB.Close()

    // Prepare test data: create multiple employees
    testEmployees := []models.Employee{
        {
            Name:  "John Doe",
            Email: "john.doe@company.com",
            Role:  "Employee",
        },
        {
            Name:  "Jane Smith",
            Email: "jane.smith@company.com",
            Role:  "Manager",
        },
        {
            Name:  "Bob Johnson",
            Email: "bob.johnson@company.com",
            Role:  "Admin",
        },
    }

    // Insert test employees into the database
    for _, emp := range testEmployees {
        err := testDB.Create(&emp).Error
        assert.NoError(t, err)
    }

    testCases := []struct {
        name           string
        page           int
        limit          int
        expectedCount  int
        expectedStatus int
        expectedTotal  int
    }{
        {
            name:           "Default Pagination",
            page:           1,
            limit:          10,
            expectedCount:  3,
            expectedStatus: http.StatusOK,
            expectedTotal:  3,
        },
        {
            name:           "First Page with Limited Results",
            page:           1,
            limit:          2,
            expectedCount:  2,
            expectedStatus: http.StatusOK,
            expectedTotal:  3,
        },
        {
            name:           "Second Page with Limited Results",
            page:           2,
            limit:          2,
            expectedCount:  1,
            expectedStatus: http.StatusOK,
            expectedTotal:  3,
        },
        {
            name:           "Page Beyond Total Results",
            page:           3,
            limit:          2,
            expectedCount:  0,
            expectedStatus: http.StatusOK,
            expectedTotal:  3,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Create request with pagination parameters
            req, err := http.NewRequest("GET", fmt.Sprintf("/employees?page=%d&limit=%d", tc.page, tc.limit), nil)
            assert.NoError(t, err)

            // Create response recorder
            rr := httptest.NewRecorder()

            // Call the handler
            GetEmployees(rr, req)

            // Check the status code
            assert.Equal(t, tc.expectedStatus, rr.Code)

            // Parse the response
            var response utils.StandardResponse
            err = json.Unmarshal(rr.Body.Bytes(), &response)
            assert.NoError(t, err)

            // Check metadata
            assert.NotNil(t, response.Meta)
            meta, ok := response.Meta.(map[string]interface{})
            assert.True(t, ok)

            // Verify metadata fields
            assert.Equal(t, float64(tc.page), meta["page"])
            assert.Equal(t, float64(tc.limit), meta["limit"])
            assert.Equal(t, float64(tc.expectedTotal), meta["total"])
            
            // Verify employees
            employees, ok := response.Data.([]interface{})
            assert.True(t, ok)
            assert.Equal(t, tc.expectedCount, len(employees))

            // If there are employees, verify their structure
            if tc.expectedCount > 0 {
                for _, emp := range employees {
                    empMap, ok := emp.(map[string]interface{})
                    assert.True(t, ok)
                    assert.Contains(t, []string{"John Doe", "Jane Smith", "Bob Johnson"}, empMap["name"])
                }
            }
        })
    }
}