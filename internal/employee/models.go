package employee

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Employee struct {
	gorm.Model
	Name                     string          `json:"name"`
	Designation              string          `json:"designation"` // e.g., Doctor, Nurse, Admin
	Salary                   float64         `json:"salary"`
	Email                    string          `json:"email" gorm:"unique"`
	PhoneNumber              string          `json:"phone_number"`
	HireDate                 time.Time       `json:"hire_date"` // Use time.Time for actual date handling
	DateOfBirth              time.Time       `json:"date_of_birth"`
	Gender                   string          `json:"gender"`
	Address                  string          `json:"address"`
	Country                  string          `json:"country"`
	State                    string          `json:"state"`
	MaritalStatus            string          `json:"marital_status"`
	Children                 int             `json:"children"`
	EmergencyContact         string          `json:"emergency_contact"`
	EmergencyContactRelation string          `json:"emergency_contact_relation"`
	Shifts                   []EmployeeShift `gorm:"foreignkey:EmployeeID"`
}

type Shift struct {
	gorm.Model
	Name      string          `json:"name"`
	StartTime time.Time       `json:"start_time"`
	EndTime   time.Time       `json:"end_time"`
	Employees []EmployeeShift `gorm:"foreignkey:ShiftID"` // Relationship with EmployeeShift
}

type Attendance struct {
	gorm.Model
	EmployeeID   uint      `gorm:"not null" json:"employee_id"` // Foreign key
	ShiftID      uint      `gorm:"not null" json:"shift_id"`    // Foreign key
	Date         time.Time `gorm:"type:date;not null" json:"date"`
	ClockInTime  time.Time `json:"clock_in_time"`
	ClockOutTime time.Time `json:"clock_out_time"`
	Status       string    `gorm:"not null" json:"status"` // e.g., Present/Absent

	Employee Employee `gorm:"foreignkey:EmployeeID"` // Relationship with Employee
	Shift    Shift    `gorm:"foreignkey:ShiftID"`    // Relationship with Shift
}

type EmployeeShift struct {
	gorm.Model
	EmployeeID uint      `gorm:"not null" json:"employee_id"` // Foreign key
	ShiftID    uint      `gorm:"not null" json:"shift_id"`    // Foreign key
	StartDate  time.Time `gorm:"type:date;not null" json:"start_date"`
	EndDate    time.Time `gorm:"type:date;not null" json:"end_date"`

	Employee Employee `gorm:"foreignkey:EmployeeID"`
	Shift    Shift    `gorm:"foreignkey:ShiftID"` // Relationship with Shift

	UniqueIndex string `gorm:"unique_index:idx_employee_shift;employee_id,shift_id,start_date,end_date"`
}
