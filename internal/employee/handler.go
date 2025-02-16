// internal/employee/handler.go
package employee

import (
	"clinicplus/internal/shared/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

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

func (h *EmployeeHandler) SearchEmployees(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1   // default page
	limit := 10 // default limit

	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	employees, total, err := h.service.SearchEmployees(query, page, limit)
	if err != nil {
		log.Printf("Error searching employees: %v", err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, nil, "Failed to search employees", nil)
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

func (h *EmployeeHandler) ClockInEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		log.Printf("Invalid employee ID: %v", err)
		utils.SendJSONResponse(w, http.StatusBadRequest, nil, "Invalid employee ID", nil)
		return
	}

	var request struct {
		ShiftID uint `json:"shift_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Error decoding request body: %v", err)
		utils.SendJSONResponse(w, http.StatusBadRequest, nil, "Invalid request body", nil)
		return
	}

	attendance, err := h.service.ClockIn(uint(employeeID), request.ShiftID)
	if err != nil {
		log.Printf("Error clocking in employee: %v", err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, nil, "Failed to clock in", nil)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, attendance, nil, nil)
}

func (h *EmployeeHandler) ClockOutEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		log.Printf("Invalid employee ID: %v", err)
		utils.SendJSONResponse(w, http.StatusBadRequest, nil, "Invalid employee ID", nil)
		return
	}

	var request struct {
		ShiftID uint `json:"shift_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Error decoding request body: %v", err)
		utils.SendJSONResponse(w, http.StatusBadRequest, nil, "Invalid request body", nil)
		return
	}

	attendance, err := h.service.ClockOut(uint(employeeID), request.ShiftID)
	if err != nil {
		log.Printf("Error clocking out employee: %v", err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, nil, "Failed to clock out", nil)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, attendance, nil, nil)
}

func (h *EmployeeHandler) CreateShift(w http.ResponseWriter, r *http.Request) {
	var shift Shift
	if err := json.NewDecoder(r.Body).Decode(&shift); err != nil {
		log.Printf("Error decoding request body: %v", err)
		utils.SendJSONResponse(w, http.StatusBadRequest, nil, "Invalid request body", nil)
		return
	}

	createdShift, err := h.service.CreateShift(shift)
	if err != nil {
		log.Printf("Error creating shift: %v", err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, nil, "Failed to create shift", nil)
		return
	}

	utils.SendJSONResponse(w, http.StatusCreated, createdShift, nil, nil)
}

func (h *EmployeeHandler) GetShift(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid shift ID: %v", err)
		utils.SendJSONResponse(w, http.StatusBadRequest, nil, "Invalid shift ID", nil)
		return
	}

	shift, err := h.service.GetShift(uint(id))
	if err != nil {
		log.Printf("Error fetching shift: %v", err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, nil, "Failed to retrieve shift", nil)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, shift, nil, nil)
}

func (h *EmployeeHandler) UpdateShift(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid shift ID: %v", err)
		utils.SendJSONResponse(w, http.StatusBadRequest, nil, "Invalid shift ID", nil)
		return
	}

	var shift Shift
	if err := json.NewDecoder(r.Body).Decode(&shift); err != nil {
		log.Printf("Error decoding request body: %v", err)
		utils.SendJSONResponse(w, http.StatusBadRequest, nil, "Invalid request body", nil)
		return
	}

	updatedShift, err := h.service.UpdateShift(uint(id), shift)
	if err != nil {
		log.Printf("Error updating shift: %v", err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, nil, "Failed to update shift", nil)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, updatedShift, nil, nil)
}

func (h *EmployeeHandler) DeleteShift(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid shift ID: %v", err)
		utils.SendJSONResponse(w, http.StatusBadRequest, nil, "Invalid shift ID", nil)
		return
	}

	if err := h.service.DeleteShift(uint(id)); err != nil {
		log.Printf("Error deleting shift: %v", err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, nil, "Failed to delete shift", nil)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, nil, "Shift deleted successfully", nil)
}

func (h *EmployeeHandler) AssignShift(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		log.Printf("Invalid employee ID: %v", err)
		utils.SendJSONResponse(w, http.StatusBadRequest, nil, "Invalid employee ID", nil)
		return
	}

	var request struct {
		ShiftID   uint      `json:"shift_id"`
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Error decoding request body: %v", err)
		utils.SendJSONResponse(w, http.StatusBadRequest, nil, "Invalid request body", nil)
		return
	}

	assignedShift, err := h.service.AssignShift(uint(employeeID), request.ShiftID, request.StartDate, request.EndDate)
	if err != nil {
		log.Printf("Error assigning shift: %v", err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, nil, "Failed to assign shift", nil)
		return
	}

	utils.SendJSONResponse(w, http.StatusCreated, assignedShift, nil, nil)
}
