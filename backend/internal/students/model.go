package students

import "time"

// Student status values.
const (
	StatusActive      = "active"
	StatusGraduated   = "graduated"
	StatusTransferred = "transferred"
	StatusSuspended   = "suspended"
	StatusWithdrawn   = "withdrawn"
)

// Statuses lists the valid student statuses.
var Statuses = []string{StatusActive, StatusGraduated, StatusTransferred, StatusSuspended, StatusWithdrawn}

// Student is an enrolled pupil belonging to a school (tenant).
type Student struct {
	ID              int64     `json:"id"`
	SchoolID        int64     `json:"school_id"`
	UserID          *int64    `json:"user_id,omitempty"`
	AdmissionNumber string    `json:"admission_number"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	DateOfBirth     time.Time `json:"date_of_birth"`
	Gender          string    `json:"gender"`
	ClassID         *int64    `json:"class_id,omitempty"`
	ParentID        *int64    `json:"parent_id,omitempty"`
	Address         string    `json:"address"`
	Phone           string    `json:"phone"`
	Email           string    `json:"email"`
	PhotoURL        string    `json:"photo_url"`
	EnrollmentDate  time.Time `json:"enrollment_date"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Version         int32     `json:"-"`
}
