package models

import (
    "github.com/jinzhu/gorm"
)

type Expense struct {
    gorm.Model
    EmployeeID uint    `json:"employee_id"`
    Date       string  `json:"date"`
    Amount     float64 `json:"amount"`
    Category   string  `json:"category"` // e.g., Utilities, Supplies
    Description string  `json:"description"`
    Status     string  `json:"status"` // e.g., Pending, Approved, Rejected
}