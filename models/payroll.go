package models

import (
    "github.com/jinzhu/gorm"
)

type Payroll struct {
    gorm.Model
    EmployeeID   uint    `json:"employee_id"`
    Month        int     `json:"month"`
    Year         int     `json:"year"`
    BaseSalary   float64 `json:"base_salary"`
    OvertimePay  float64 `json:"overtime_pay"`
    Deductions   float64 `json:"deductions"`
    Bonuses      float64 `json:"bonuses"`
    NetSalary    float64 `json:"net_salary"`
}