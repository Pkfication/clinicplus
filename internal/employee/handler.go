// internal/employee/handler.go
package employee

import (
	"clinicplus/internal/shared/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type EmployeeHandler struct {
	service EmployeeService
}

func NewEmployeeHandler(service EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: service}
}

func (h *EmployeeHandler) GetEmployees(w http.ResponseWriter, r *http.Request) {
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

	// Fetch employees from the service
	employees, total, err := h.service.GetEmployees(page, limit)
	if err != nil {
		log.Printf("Error fetching employees: %v", err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, nil, "Failed to retrieve employees", nil)
		return
	}

	// Prepare metadata
	meta := map[string]interface{}{
		"page":        page,
		"limit":       limit,
		"total":       total,
		"total_pages": (total + limit - 1) / limit,
	}

	// Send response
	utils.SendJSONResponse(w, http.StatusOK, employees, nil, meta)
}

func (h *EmployeeHandler) GetEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid employee ID: %v", err)
		utils.SendJSONResponse(w, http.StatusBadRequest, nil, "Invalid employee ID", nil)
		return
	}

	employee, err := h.service.GetEmployee(id)
	if err != nil {
		if err.Error() == "employee not found" {
			utils.SendJSONResponse(w, http.StatusNotFound, nil, "Employee not found", nil)
		} else {
			log.Printf("Error fetching employee: %v", err)
			utils.SendJSONResponse(w, http.StatusInternalServerError, nil, "Failed to retrieve employee", nil)
		}
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, employee, nil, nil)
}

func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var createRequest struct {
		Employee
	}

	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		log.Printf("Error decoding request body: %v", err)
		utils.SendJSONResponse(w, http.StatusBadRequest, nil, "Invalid request body", nil)
		return
	}

	employee, err := h.service.CreateEmployee(createRequest.Employee)
	if err != nil {
		log.Printf("Error creating employee: %v", err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, nil, "Failed to create employee", nil)
		return
	}

	meta := map[string]interface{}{
		"message": "Employee and user created successfully",
	}

	utils.SendJSONResponse(w, http.StatusCreated, employee, nil, meta)
}

func (h *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid employee ID: %v", err)
		utils.SendJSONResponse(w, http.StatusBadRequest, nil, "Invalid employee ID", nil)
		return
	}

	var updatedEmployee Employee
	if err := json.NewDecoder(r.Body).Decode(&updatedEmployee); err != nil {
		log.Printf("Error decoding request body: %v", err)
		utils.SendJSONResponse(w, http.StatusBadRequest, nil, "Invalid request body", nil)
		return
	}

	employee, err := h.service.UpdateEmployee(id, updatedEmployee)
	if err != nil {
		if err.Error() == "employee not found" {
			utils.SendJSONResponse(w, http.StatusNotFound, nil, "Employee not found", nil)
		} else {
			log.Printf("Error updating employee: %v", err)
			utils.SendJSONResponse(w, http.StatusInternalServerError, nil, "Failed to update employee", nil)
		}
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, employee, nil, map[string]interface{}{
		"message": "Employee updated successfully",
	})
}

func (h *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid employee ID: %v", err)
		utils.SendJSONResponse(w, http.StatusBadRequest, nil, "Invalid employee ID", nil)
		return
	}

	if err := h.service.DeleteEmployee(id); err != nil {
		if err.Error() == "employee not found" {
			utils.SendJSONResponse(w, http.StatusNotFound, nil, "Employee not found", nil)
		} else {
			log.Printf("Error deleting employee: %v", err)
			utils.SendJSONResponse(w, http.StatusInternalServerError, nil, "Failed to delete employee", nil)
		}
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, nil, nil, map[string]interface{}{
		"message": "Employee deleted successfully",
	})
}
