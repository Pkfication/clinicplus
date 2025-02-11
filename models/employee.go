package models

import (
    "github.com/jinzhu/gorm"
)

type Employee struct {
    gorm.Model
    Name         string  `json:"name"`
    Role         string  `json:"role"` // e.g., Doctor, Nurse, Admin
    Salary       float64 `json:"salary"`
    Email        string  `json:"email" gorm:"unique"`
    PhoneNumber  string  `json:"phone_number"`
    HireDate     string  `json:"hire_date"` // Use time.Time for actual date handling
    LeaveBalance int     `json:"leave_balance"`
}