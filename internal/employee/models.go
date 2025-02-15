package employee

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Employee struct {
	gorm.Model
	Name                     string    `json:"name"`
	Designation              string    `json:"designation"` // e.g., Doctor, Nurse, Admin
	Salary                   float64   `json:"salary"`
	Email                    string    `json:"email" gorm:"unique"`
	PhoneNumber              string    `json:"phone_number"`
	HireDate                 time.Time `json:"hire_date"` // Use time.Time for actual date handling
	DateOfBirth              time.Time `json:"date_of_birth"`
	Gender                   string    `json:"gender"`
	Address                  string    `json:"address"`
	Country                  string    `json:"country"`
	State                    string    `json:"state"`
	MaritalStatus            string    `json:"marital_status"`
	Children                 int       `json:"children"`
	EmergencyContact         string    `json:"emergency_contact"`
	EmergencyContactRelation string    `json:"emergency_contact_relation"`
}
