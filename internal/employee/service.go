// internal/employee/service.go
package employee

import (
	"errors"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

type EmployeeService interface {
	GetEmployees(page, limit int) ([]Employee, int, error)
	GetEmployee(id int) (*Employee, error)
	CreateEmployee(employee Employee) (*Employee, error)
	UpdateEmployee(id int, employee Employee) (*Employee, error)
	DeleteEmployee(id int) error
	SearchEmployees(query string, page, limit int) ([]Employee, int, error)

	ClockIn(employeeID uint, shiftID uint) (*Attendance, error)
	ClockOut(employeeID uint, shiftID uint) (*Attendance, error)

	CreateShift(shift Shift) (*Shift, error)
	GetShift(id uint) (*Shift, error)
	UpdateShift(id uint, shift Shift) (*Shift, error)
	DeleteShift(id uint) error
	AssignShift(employeeID uint, shiftID uint, startDate time.Time, endDate time.Time) (*EmployeeShift, error)
}

type employeeService struct {
	db *gorm.DB
}

func NewEmployeeService(db *gorm.DB) EmployeeService {
	return &employeeService{db: db}
}

func (s *employeeService) GetEmployees(page, limit int) ([]Employee, int, error) {
	var employees []Employee
	offset := (page - 1) * limit

	if err := s.db.Offset(offset).Limit(limit).Find(&employees).Error; err != nil {
		log.Printf("Error fetching employees: %v", err)
		return nil, 0, err
	}

	var total int
	s.db.Model(&Employee{}).Count(&total)

	return employees, total, nil
}

func (s *employeeService) GetEmployee(id int) (*Employee, error) {
	var employee Employee
	if err := s.db.First(&employee, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New("employee not found")
		}
		log.Printf("Error fetching employee: %v", err)
		return nil, err
	}
	return &employee, nil
}

func (s *employeeService) CreateEmployee(employee Employee) (*Employee, error) {
	// Start a transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		log.Printf("Error starting transaction: %v", tx.Error)
		return nil, tx.Error
	}

	// Create the employee
	if err := tx.Create(&employee).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating employee: %v", err)
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction: %v", err)
		return nil, err
	}

	return &employee, nil
}

func (s *employeeService) UpdateEmployee(id int, employee Employee) (*Employee, error) {
	var existingEmployee Employee
	if err := s.db.First(&existingEmployee, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New("employee not found")
		}
		log.Printf("Error finding employee: %v", err)
		return nil, err
	}

	// Preserve ID and update other fields
	employee.ID = existingEmployee.ID

	if err := s.db.Save(&employee).Error; err != nil {
		log.Printf("Error updating employee: %v", err)
		return nil, err
	}

	return &employee, nil
}

func (s *employeeService) DeleteEmployee(id int) error {
	result := s.db.Delete(&Employee{}, id)
	if result.Error != nil {
		log.Printf("Error deleting employee: %v", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("employee not found")
	}

	return nil
}

func (s *employeeService) SearchEmployees(query string, page, limit int) ([]Employee, int, error) {
	var employees []Employee
	searchPattern := "%" + query + "%"

	// Calculate offset
	offset := (page - 1) * limit

	err := s.db.Where(
		"name ILIKE ? OR designation ILIKE ? OR email ILIKE ? OR phone_number ILIKE ? OR address ILIKE ? OR country ILIKE ? OR state ILIKE ? OR marital_status ILIKE ? OR emergency_contact ILIKE ?",
		searchPattern, searchPattern, searchPattern, searchPattern, searchPattern,
		searchPattern, searchPattern, searchPattern, searchPattern, // 9 parameters
	).Offset(offset).Limit(limit).Find(&employees).Error

	if err != nil {
		log.Printf("Error searching employees: %v", err)
		return nil, 0, err
	}

	// Count total matching employees
	var total int
	s.db.Model(&Employee{}).Where(
		"name ILIKE ? OR designation ILIKE ? OR email ILIKE ? OR phone_number ILIKE ? OR address ILIKE ? OR country ILIKE ? OR state ILIKE ? OR marital_status ILIKE ? OR emergency_contact ILIKE ?",
		searchPattern, searchPattern, searchPattern, searchPattern, searchPattern,
		searchPattern, searchPattern, searchPattern, searchPattern, // 9 parameters
	).Count(&total)

	return employees, total, nil
}

// ClockIn allows an employee to clock in for a specific shift
func (s *employeeService) ClockIn(employeeID uint, shiftID uint) (*Attendance, error) {
	var attendance Attendance
	now := time.Now().UTC() // Get the current time in UTC
	date := now.Truncate(24 * time.Hour)

	// Check if the employee is already clocked in for today and the specified shift
	err := s.db.Where("employee_id = ? AND date = ? AND shift_id = ?", employeeID, date, shiftID).First(&attendance).Error
	if err == nil {
		return nil, errors.New("employee is already clocked in for today for this shift")
	}

	// Create a new attendance record
	attendance = Attendance{
		EmployeeID:  employeeID,
		ShiftID:     shiftID, // Associate the attendance with the shift
		Date:        date,
		ClockInTime: now,
		Status:      "Present",
	}

	if err := s.db.Create(&attendance).Error; err != nil {
		log.Printf("Error clocking in employee: %v", err)
		return nil, err
	}

	return &attendance, nil
}

// ClockOut allows an employee to clock out for a specific shift
func (s *employeeService) ClockOut(employeeID uint, shiftID uint) (*Attendance, error) {
	var attendance Attendance
	now := time.Now().UTC() // Get the current time in UTC
	date := now.Truncate(24 * time.Hour)

	// Find the attendance record for today and the specified shift
	err := s.db.Where("employee_id = ? AND date = ? AND shift_id = ?", employeeID, date, shiftID).First(&attendance).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New("no clock-in record found for today for this shift")
		}
		log.Printf("Error fetching attendance record: %v", err)
		return nil, err
	}

	// Update the clock-out time and status
	attendance.ClockOutTime = now
	attendance.Status = "Present" // Adjust status if needed

	if err := s.db.Save(&attendance).Error; err != nil {
		log.Printf("Error clocking out employee: %v", err)
		return nil, err
	}

	return &attendance, nil
}

// CreateShift creates a new shift with overlapping constraints
func (s *employeeService) CreateShift(shift Shift) (*Shift, error) {
	// Check for overlapping shifts
	var existingShifts []Shift
	err := s.db.Where("start_time < ? AND end_time > ? AND id != ?", shift.EndTime, shift.StartTime, shift.ID).Find(&existingShifts).Error
	if err != nil {
		log.Printf("Error checking for overlapping shifts: %v", err)
		return nil, err
	}

	if len(existingShifts) > 0 {
		return nil, errors.New("shift overlaps with existing shifts")
	}

	if err := s.db.Create(&shift).Error; err != nil {
		log.Printf("Error creating shift: %v", err)
		return nil, err
	}
	return &shift, nil
}

// GetShift retrieves a shift by ID
func (s *employeeService) GetShift(id uint) (*Shift, error) {
	var shift Shift
	if err := s.db.First(&shift, id).Error; err != nil {
		log.Printf("Error fetching shift: %v", err)
		return nil, err
	}
	return &shift, nil
}

// UpdateShift updates an existing shift
func (s *employeeService) UpdateShift(id uint, shift Shift) (*Shift, error) {
	shift.ID = id // Ensure the ID is set for the update
	if err := s.db.Save(&shift).Error; err != nil {
		log.Printf("Error updating shift: %v", err)
		return nil, err
	}
	return &shift, nil
}

// DeleteShift deletes a shift by ID
func (s *employeeService) DeleteShift(id uint) error {
	if err := s.db.Delete(&Shift{}, id).Error; err != nil {
		log.Printf("Error deleting shift: %v", err)
		return err
	}
	return nil
}

// AssignShift assigns a shift to an employee
func (s *employeeService) AssignShift(employeeID uint, shiftID uint, startDate time.Time, endDate time.Time) (*EmployeeShift, error) {
	// Create a new EmployeeShift record
	employeeShift := EmployeeShift{
		EmployeeID: employeeID,
		ShiftID:    shiftID,
		StartDate:  startDate,
		EndDate:    endDate,
	}

	// Check for overlapping shifts
	var existingShifts []EmployeeShift
	err := s.db.Where("employee_id = ? AND shift_id = ? AND ((start_date <= ? AND end_date >= ?) OR (start_date >= ? AND start_date <= ?))",
		employeeID, shiftID, endDate, startDate, startDate, endDate).Find(&existingShifts).Error
	if err != nil {
		log.Printf("Error checking for overlapping shifts: %v", err)
		return nil, err
	}

	if len(existingShifts) > 0 {
		return nil, errors.New("shift overlaps with existing shifts for this employee")
	}

	// Save the new EmployeeShift record
	if err := s.db.Create(&employeeShift).Error; err != nil {
		log.Printf("Error assigning shift: %v", err)
		return nil, err
	}

	return &employeeShift, nil
}
