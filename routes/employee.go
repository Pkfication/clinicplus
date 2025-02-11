package routes

import (
    "clinicplus/models"
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/jinzhu/gorm"
)

func RegisterRoutes(r *mux.Router) {
    r.HandleFunc("/employees", GetEmployees).Methods("GET")
    r.HandleFunc("/employees/{id}", GetEmployee).Methods("GET")
    r.HandleFunc("/employees", CreateEmployee).Methods("POST")
    r.HandleFunc("/employees/{id}", UpdateEmployee).Methods("PUT")
    r.HandleFunc("/employees/{id}", DeleteEmployee).Methods("DELETE")
}

// Define handler functions (GetEmployees, GetEmployee, CreateEmployee, UpdateEmployee, DeleteEmployee)
// Assuming db is a global variable for the database connection
var db *gorm.DB

func SetDB(database *gorm.DB) {
    db = database
}

// GetEmployees retrieves all employees
func GetEmployees(w http.ResponseWriter, r *http.Request) {
    var employees []models.Employee
    if err := db.Find(&employees).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(employees)
}

// GetEmployee retrieves a single employee by ID
func GetEmployee(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    var employee models.Employee
    if err := db.First(&employee, vars["id"]).Error; err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(employee)
}

// CreateEmployee creates a new employee
func CreateEmployee(w http.ResponseWriter, r *http.Request) {
    var employee models.Employee
    if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    if err := db.Create(&employee).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(employee)
}

// UpdateEmployee updates an existing employee
func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    var employee models.Employee
    if err := db.First(&employee, vars["id"]).Error; err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    db.Save(&employee)
    json.NewEncoder(w).Encode(employee)
}


// DeleteEmployee deletes an employee by ID
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    if err := db.Delete(&models.Employee{}, vars["id"]).Error; err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}