package classes

// CreateClassInput is the payload accepted by the service to create a class.
type CreateClassInput struct {
	SchoolID       int64
	Name           string
	GradeLevel     string
	Section        string
	ClassTeacherID *int64
	Capacity       int
	AcademicYear   string
}

// UpdateClassInput carries optional fields for a partial update.
type UpdateClassInput struct {
	Name           *string
	GradeLevel     *string
	Section        *string
	ClassTeacherID *int64
	Capacity       *int
	AcademicYear   *string
}

// ListFilter narrows and paginates a class listing within a school.
type ListFilter struct {
	SchoolID     int64
	AcademicYear string
	Search       string
	Page         int
	PageSize     int
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
