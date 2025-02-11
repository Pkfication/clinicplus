package routes

import (
    "bytes"
    "encoding/json"
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

func TestGetEmployees(t *testing.T) {
    // Setup test database
    setupTestDB()
    defer testDB.Close()

    // Create some test employees
    employees := []models.Employee{
        {Name: "John Doe", Email: "john.doe@example.com", Role: "Employee"},
        {Name: "Jane Smith", Email: "jane.smith@example.com", Role: "Manager"},
    }
    for _, emp := range employees {
        testDB.Create(&emp)
    }

    // Create request
    req, err := http.NewRequest("GET", "/employees?page=1&limit=10", nil)
    assert.NoError(t, err)

    // Create response recorder
    rr := httptest.NewRecorder()

    // Call the handler
    GetEmployees(rr, req)

    // Check the status code
    assert.Equal(t, http.StatusOK, rr.Code)

    // Print raw response body for debugging
    rawBody := rr.Body.Bytes()
    t.Logf("Raw Response Body: %s", string(rawBody))

    // Parse the response
    var response utils.StandardResponse
    err = json.Unmarshal(rawBody, &response)
    if err != nil {
        t.Fatalf("Failed to unmarshal response: %v", err)
    }

    // Convert response data to employees
    var retrievedEmployees []models.Employee
    
    // Convert Data to JSON first, then unmarshal
    dataBytes, err := json.Marshal(response.Data)
    assert.NoError(t, err)
    
    err = json.Unmarshal(dataBytes, &retrievedEmployees)
    if err != nil {
        t.Fatalf("Failed to unmarshal employees: %v", err)
    }

    // Verify employees
    assert.GreaterOrEqual(t, len(retrievedEmployees), 2)
    assert.NotNil(t, response.Meta)
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

func TestCreateEmployee(t *testing.T) {
    // Setup test database
    setupTestDB()
    defer testDB.Close()

    // Prepare test employee
    employee := models.Employee{
        Name:  "John Doe",
        Email: "john.doe@company.com",
        Role:  "Employee",
    }

    // Convert employee to JSON
    jsonEmployee, err := json.Marshal(employee)
    assert.NoError(t, err)

    // Create request
    req, err := http.NewRequest("POST", "/employees", bytes.NewBuffer(jsonEmployee))
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