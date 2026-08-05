package attendance

import "time"

// CreateAttendanceInput is the payload accepted by the service to mark attendance.
type CreateAttendanceInput struct {
	SchoolID  int64
	StudentID int64
	ClassID   int64
	Date      time.Time
	Status    string
	Remark    string
	MarkedBy  *int64
}

// UpdateAttendanceInput carries optional fields for correcting a mark.
type UpdateAttendanceInput struct {
	Status *string
	Remark *string
}

// ListFilter narrows and paginates an attendance listing within a school.
type ListFilter struct {
	SchoolID  int64
	ClassID   *int64
	StudentID *int64
	Status    string
	Date      *time.Time
	Page      int
	PageSize  int
}

// Limit returns the page size, defaulting/capping to sane bounds.
func (f ListFilter) Limit() int {
	switch {
	case f.PageSize <= 0:
		return 50
	case f.PageSize > 200:
		return 200
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
