package teachers

import "time"

// Teacher employment status values.
const (
	StatusActive    = "active"
	StatusOnLeave   = "on_leave"
	StatusSuspended = "suspended"
	StatusResigned  = "resigned"
)

// Statuses lists the valid teacher statuses.
var Statuses = []string{StatusActive, StatusOnLeave, StatusSuspended, StatusResigned}

// Teacher is a staff member who teaches, scoped to a school (tenant).
type Teacher struct {
	ID             int64     `json:"id"`
	SchoolID       int64     `json:"school_id"`
	UserID         *int64    `json:"user_id,omitempty"`
	EmployeeNumber string    `json:"employee_number"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Gender         string    `json:"gender"`
	Qualification  string    `json:"qualification"`
	HireDate       time.Time `json:"hire_date"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Version        int32     `json:"-"`
}
