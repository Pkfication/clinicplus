package models

import (
    "github.com/jinzhu/gorm"
)

type Attendance struct {
    gorm.Model
    EmployeeID   uint   `json:"employee_id"`
    CheckInTime  string `json:"check_in_time"` // Use time.Time for actual date handling
    CheckOutTime string `json:"check_out_time"`
    Date         string `json:"date"`
    Status       string `json:"status"` // e.g., Present, Absent, Leave
}