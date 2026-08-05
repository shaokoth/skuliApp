package students

import "time"

// CreateStudentInput is the payload accepted by the service to admit a student.
type CreateStudentInput struct {
	SchoolID        int64
	AdmissionNumber string
	FirstName       string
	LastName        string
	DateOfBirth     time.Time
	Gender          string
	ClassID         *int64
	ParentID        *int64
	Address         string
	Phone           string
	Email           string
	EnrollmentDate  time.Time
}

// UpdateStudentInput carries optional fields for a partial update.
type UpdateStudentInput struct {
	FirstName *string
	LastName  *string
	Gender    *string
	ClassID   *int64
	ParentID  *int64
	Address   *string
	Phone     *string
	Email     *string
	Status    *string
}

// ListFilter narrows and paginates a student listing within a school.
type ListFilter struct {
	SchoolID int64
	ClassID  *int64
	Status   string
	Search   string
	Page     int
	PageSize int
}

// Limit returns the page size, defaulting/capping to sane bounds.
func (f ListFilter) Limit() int {
	switch {
	case f.PageSize <= 0:
		return 20
	case f.PageSize > 100:
		return 100
	default:
		return f.PageSize
	}
}

// Offset returns the SQL offset for the requested page.
func (f ListFilter) Offset() int {
	if f.Page <= 1 {
		return 0
	}
	return (f.Page - 1) * f.Limit()
}
