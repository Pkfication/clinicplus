// internal/employee/service.go
package employee

import (
	"errors"
	"log"

	"github.com/jinzhu/gorm"
)

type EmployeeService interface {
	GetEmployees(page, limit int) ([]Employee, int, error)
	GetEmployee(id int) (*Employee, error)
	CreateEmployee(employee Employee) (*Employee, error)
	UpdateEmployee(id int, employee Employee) (*Employee, error)
	DeleteEmployee(id int) error
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
