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
    var employee models.Employee

    // Decode request body
    if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
        log.Printf("Error decoding request body: %v", err)
        utils.SendJSONResponse(w, http.StatusBadRequest, nil, "Invalid request body", nil)
        return
    }

    // Validate required fields
    if employee.Name == "" {
        utils.SendJSONResponse(w, http.StatusBadRequest, nil, "Name is required", nil)
        return
    }

    // Set default values or generate additional fields if needed
    if employee.Email == "" {
        // Generate email if not provided
        employee.Email = generateEmail(employee.Name, "")
    }

    // Create the employee
    if err := db.Create(&employee).Error; err != nil {
        log.Printf("Error creating employee: %v", err)
        utils.SendJSONResponse(w, http.StatusInternalServerError, nil, "Failed to create employee", nil)
        return
    }

    // Prepare response metadata
    meta := map[string]interface{}{
        "message": "Employee created successfully",
    }

    // Send standardized response
    utils.SendJSONResponse(w, http.StatusCreated, employee, nil, meta)
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