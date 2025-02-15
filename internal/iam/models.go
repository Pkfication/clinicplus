package iam

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	EmployeeID   uint   `json:"employee_id"`
	Username     string `json:"username" gorm:"unique"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"` // e.g., Admin, Manager, Employee
}
