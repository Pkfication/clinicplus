package routes

import (
    "clinicplus/models"
    "clinicplus/utils"
    "encoding/json"
    "net/http"
    "strconv"
    "log"
    "fmt"
    "strings"
    "github.com/gorilla/mux"
    "github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

var db *gorm.DB

// SetDB sets the database connection
func SetDB(database *gorm.DB) {
    db = database
}

// GetEmployees retrieves all employees with pagination
func GetEmployees(w http.ResponseWriter, r *http.Request) {
    var employees []models.Employee
    
    // Pagination support
    page := 1
    limit := 10
    
    pageStr := r.URL.Query().Get("page")
    limitStr := r.URL.Query().Get("limit")
    
    if pageStr != "" {
        page, _ = strconv.Atoi(pageStr)
    }
    
    if limitStr != "" {
        limit, _ = strconv.Atoi(limitStr)
    }
    
    offset := (page - 1) * limit
    
    // Fetch employees with pagination
    if err := db.Offset(offset).Limit(limit).Find(&employees).Error; err != nil {
        log.Printf("Error fetching employees: %v", err)
        utils.SendJSONResponse(w, http.StatusInternalServerError, nil, "Failed to retrieve employees", nil)
        return
    }
    
    // Count total employees for metadata
    var total int
    db.Model(&models.Employee{}).Count(&total)
    
    // Prepare metadata
    meta := map[string]interface{}{
        "page":       page,
        "limit":      limit,
        "total":      total,
        "total_pages": (total + limit - 1) / limit,
    }
    
    // Wrap employees in the data field
    utils.SendJSONResponse(w, http.StatusOK, employees, nil, meta)
}

func GetEmployee(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]
    
    // Convert ID to integer
    id, err := strconv.Atoi(idStr)
    if err != nil {
        log.Printf("Invalid employee ID: %v", err)
        utils.SendJSONResponse(w, http.StatusBadRequest, nil, "Invalid employee ID", nil)
        return
    }
    
    var employee models.Employee
    if err := db.First(&employee, id).Error; err != nil {
        if gorm.IsRecordNotFoundError(err) {
            utils.SendJSONResponse(w, http.StatusNotFound, nil, "Employee not found", nil)
        } else {
            log.Printf("Error fetching employee: %v", err)
            utils.SendJSONResponse(w, http.StatusInternalServerError, nil, "Failed to retrieve employee", nil)
        }
        return
    }
    
    utils.SendJSONResponse(w, http.StatusOK, employee, nil, nil)
}

// CreateEmployee creates a new employee
func CreateEmployee(w http.ResponseWriter, r *http.Request) {
    // Start a database transaction
    tx := db.Begin()
    if tx.Error != nil {
        log.Printf("Error starting transaction: %v", tx.Error)
        utils.SendJSONResponse(w, http.StatusInternalServerError, nil, "Failed to start transaction", nil)
        return
    }

    // Create a struct to capture both employee and user details
    var createRequest struct {
        models.Employee
        Username string `json:"username"`
        Password string `json:"password"`
    }

    // Decode request body
    if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
        tx.Rollback()
        log.Printf("Error decoding request body: %v", err)
        utils.SendJSONResponse(w, http.StatusBadRequest, nil, "Invalid request body", nil)
        return
    }

    // Validate required fields
    var validationErrors []string
    
    if createRequest.Name == "" {
        validationErrors = append(validationErrors, "Name is required")
    }
    
    if createRequest.Email == "" {
        validationErrors = append(validationErrors, "Email must be explicitly provided")
    }
    
    if createRequest.Username == "" {
        validationErrors = append(validationErrors, "Username must be explicitly provided")
    }
    
    if createRequest.Password == "" {
        validationErrors = append(validationErrors, "Password must be explicitly provided")
    }

    // If there are validation errors, return them
    if len(validationErrors) > 0 {
        tx.Rollback()
        utils.SendJSONResponse(w, http.StatusBadRequest, nil, validationErrors, nil)
        return
    }

    // Create the employee
    if err := tx.Create(&createRequest.Employee).Error; err != nil {
        tx.Rollback()
        log.Printf("Error creating employee: %v", err)
        utils.SendJSONResponse(w, http.StatusInternalServerError, nil, "Failed to create employee", nil)
        return
    }

    // Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createRequest.Password), bcrypt.DefaultCost)
    if err != nil {
        tx.Rollback()
        log.Printf("Error hashing password: %v", err)
        utils.SendJSONResponse(w, http.StatusInternalServerError, nil, "Failed to secure password", nil)
        return
    }

    // Create user
    user := models.User{
        EmployeeID:    createRequest.Employee.ID,
        Username:      createRequest.Username,
        PasswordHash: string(hashedPassword),
        Role:         createRequest.Role, // Use role from employee data
    }

    // Save user
    if err := tx.Create(&user).Error; err != nil {
        tx.Rollback()
        log.Printf("Error creating user: %v", err)
        utils.SendJSONResponse(w, http.StatusInternalServerError, nil, "Failed to create user", nil)
        return
    }

    // Commit the transaction
    if err := tx.Commit().Error; err != nil {
        log.Printf("Error committing transaction: %v", err)
        utils.SendJSONResponse(w, http.StatusInternalServerError, nil, "Failed to complete employee creation", nil)
        return
    }

    // Prepare response metadata
    meta := map[string]interface{}{
        "message": "Employee and user created successfully",
    }

    // Send standardized response
    utils.SendJSONResponse(w, http.StatusCreated, createRequest.Employee, nil, meta)
}

// UpdateEmployee updates an existing employee
func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]
    
    // Convert ID to integer
    id, err := strconv.Atoi(idStr)
    if err != nil {
        log.Printf("Invalid employee ID: %v", err)
        utils.SendJSONResponse(w, http.StatusBadRequest, nil, "Invalid employee ID", nil)
        return
    }
    
    var existingEmployee models.Employee
    if err := db.First(&existingEmployee, id).Error; err != nil {
        if gorm.IsRecordNotFoundError(err) {
            utils.SendJSONResponse(w, http.StatusNotFound, nil, "Employee not found", nil)
        } else {
            log.Printf("Error finding employee: %v", err)
            utils.SendJSONResponse(w, http.StatusInternalServerError, nil, "Failed to find employee", nil)
        }
        return
    }
    
    var updatedEmployee models.Employee
    if err := json.NewDecoder(r.Body).Decode(&updatedEmployee); err != nil {
        log.Printf("Error decoding request body: %v", err)
        utils.SendJSONResponse(w, http.StatusBadRequest, nil, "Invalid request body", nil)
        return
    }
    
    // Preserve ID and only update other fields
    updatedEmployee.ID = existingEmployee.ID
    
    if err := db.Save(&updatedEmployee).Error; err != nil {
        log.Printf("Error updating employee: %v", err)
        utils.SendJSONResponse(w, http.StatusInternalServerError, nil, "Failed to update employee", nil)
        return
    }
    
    utils.SendJSONResponse(w, http.StatusOK, updatedEmployee, nil, map[string]interface{}{
        "message": "Employee updated successfully",
    })
}

// DeleteEmployee removes an employee by ID
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]
    
    // Convert ID to integer
    id, err := strconv.Atoi(idStr)
    if err != nil {
        log.Printf("Invalid employee ID: %v", err)
        utils.SendJSONResponse(w, http.StatusBadRequest, nil, "Invalid employee ID", nil)
        return
    }
    
    result := db.Delete(&models.Employee{}, id)
    if result.Error != nil {
        log.Printf("Error deleting employee: %v", result.Error)
        utils.SendJSONResponse(w, http.StatusInternalServerError, nil, "Failed to delete employee", nil)
        return
    }
    
    if result.RowsAffected == 0 {
        utils.SendJSONResponse(w, http.StatusNotFound, nil, "Employee not found", nil)
        return
    }
    
    utils.SendJSONResponse(w, http.StatusOK, nil, nil, map[string]interface{}{
        "message": "Employee deleted successfully",
    })
}

func generateEmail(name, lastName string) string {
    // Remove spaces and convert to lowercase
    nameParts := strings.Fields(name)
    var first, last string
    
    if len(nameParts) > 0 {
        first = strings.ToLower(nameParts[0])
    }
    
    if len(nameParts) > 1 {
        last = strings.ToLower(nameParts[len(nameParts)-1])
    }
    
    // Simple email generation
    if last != "" {
        return fmt.Sprintf("%s.%s@company.com", first, last)
    }
    
    return fmt.Sprintf("%s@company.com", first)
}