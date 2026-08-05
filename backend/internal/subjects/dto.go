package subjects

// CreateSubjectInput is the payload accepted by the service to create a subject.
type CreateSubjectInput struct {
	SchoolID int64
	Name     string
	Code     string
}

// UpdateSubjectInput carries optional fields for a partial update.
type UpdateSubjectInput struct {
	Name *string
	Code *string
}

// ListFilter narrows and paginates a subject listing within a school.
type ListFilter struct {
	SchoolID int64
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
