package routes

import (
	"bytes"
    "clinicplus/models"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strconv"
    "github.com/gorilla/mux"
    "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
    "github.com/stretchr/testify/assert"
    "testing"
)

var testDB *gorm.DB

// Setup function to initialize the database before tests
func setup() {
    var err error
    testDB, err = gorm.Open("sqlite3", ":memory:")
    if err != nil {
        panic("failed to connect to database")
    }
    testDB.AutoMigrate(&models.Employee{}, &models.User{}) // Ensure the Employee and User models are migrated
}

// Test for creating an employee and user account
func TestCreateEmployee(t *testing.T) {
    setup()
    router := mux.NewRouter()
    SetDB(testDB)
    RegisterRoutes(router)

    // Prepare the request body with employee and user details
    requestBody := map[string]interface{}{
        "name":         "John Doe",
        "role":         "Doctor",
        "salary":       60000,
        "email":        "john@example.com",
        "phone_number": "1234567890",
        "hire_date":    "2025-02-11",
        "leave_balance": 10,
        "username":     "johndoe",
        "password":     "securepassword",
    }

    jsonData, _ := json.Marshal(requestBody)
    req, _ := http.NewRequest("POST", "/employees", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")

    // Create a response recorder
    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusCreated, rr.Code)

    // Check if the employee was created
    var createdEmployee models.Employee
    json.Unmarshal(rr.Body.Bytes(), &createdEmployee)
    assert.Equal(t, "John Doe", createdEmployee.Name)

    // Check if the user was created in the database
    var user models.User
    err := testDB.Where("username = ?", "johndoe").First(&user).Error
    assert.NoError(t, err) // Ensure no error occurred
    assert.Equal(t, "johndoe", user.Username) // Check the username
}

// Test for getting all employees
func TestGetEmployees(t *testing.T) {
    setup()
    router := mux.NewRouter()
    SetDB(testDB)
    RegisterRoutes(router)

    req, _ := http.NewRequest("GET", "/employees", nil)
    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)
}

// Test for getting a single employee
func TestGetEmployee(t *testing.T) {
    setup()
    router := mux.NewRouter()
    SetDB(testDB)
    RegisterRoutes(router)

    // First create an employee to retrieve
    employee := models.Employee{
        Name:         "Jane Doe",
        Role:         "Nurse",
        Salary:       50000,
        Email:        "jane@example.com",
        PhoneNumber:  "0987654321",
        HireDate:     "2025-02-11",
        LeaveBalance: 15,
    }
    testDB.Create(&employee)

    req, _ := http.NewRequest("GET", "/employees/"+strconv.Itoa(int(employee.ID)), nil)
    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)
}

// Additional tests for UpdateEmployee and DeleteEmployee can be added similarly.
// Test for updating an employee
func TestUpdateEmployee(t *testing.T) {
    setup()
    router := mux.NewRouter()
    SetDB(testDB)
    RegisterRoutes(router)

    // Create an employee to update
    employee := models.Employee{
        Name:         "Alice Smith",
        Role:         "Nurse",
        Salary:       55000,
        Email:        "alice@example.com",
        PhoneNumber:  "1234567890",
        HireDate:     "2025-02-11",
        LeaveBalance: 12,
    }
    testDB.Create(&employee)

    // Update the employee's salary
    employee.Salary = 60000
    jsonData, _ := json.Marshal(employee)
    req, _ := http.NewRequest("PUT", "/employees/"+strconv.Itoa(int(employee.ID)), bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")

    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)

    var updatedEmployee models.Employee
    json.Unmarshal(rr.Body.Bytes(), &updatedEmployee)
    assert.Equal(t, employee.Salary, updatedEmployee.Salary)
}

// Test for deleting an employee
func TestDeleteEmployee(t *testing.T) {
    setup()
    router := mux.NewRouter()
    SetDB(testDB)
    RegisterRoutes(router)

    // Create an employee to delete
    employee := models.Employee{
        Name:         "Bob Johnson",
        Role:         "Admin",
        Salary:       65000,
        Email:        "bob@example.com",
        PhoneNumber:  "0987654321",
        HireDate:     "2025-02-11",
        LeaveBalance: 8,
    }
    testDB.Create(&employee)

    // Delete the employee
    req, _ := http.NewRequest("DELETE", "/employees/"+strconv.Itoa(int(employee.ID)), nil)
    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusNoContent, rr.Code)

    // Check if the employee is deleted
    var deletedEmployee models.Employee
    err := testDB.First(&deletedEmployee, employee.ID).Error
    assert.Equal(t, gorm.ErrRecordNotFound, err)
}