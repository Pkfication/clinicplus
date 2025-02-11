package models

import (
    "github.com/jinzhu/gorm"
)

type Leave struct {
    gorm.Model
    EmployeeID uint   `json:"employee_id"`
    StartDate  string `json:"start_date"`
    EndDate    string `json:"end_date"`
    LeaveType  string `json:"leave_type"` // e.g., Sick, Vacation
    Status     string `json:"status"`     // e.g., Pending, Approved, Rejected
}