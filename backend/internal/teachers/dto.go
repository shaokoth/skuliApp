package teachers

import "time"

// CreateTeacherInput is the payload accepted by the service to add a teacher.
type CreateTeacherInput struct {
	SchoolID       int64
	EmployeeNumber string
	FirstName      string
	LastName       string
	Email          string
	Phone          string
	Gender         string
	Qualification  string
	HireDate       time.Time
}

// UpdateTeacherInput carries optional fields for a partial update.
type UpdateTeacherInput struct {
	FirstName     *string
	LastName      *string
	Email         *string
	Phone         *string
	Qualification *string
	Status        *string
}

// ListFilter narrows and paginates a teacher listing within a school.
type ListFilter struct {
	SchoolID int64
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
