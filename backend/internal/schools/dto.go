package schools

// CreateSchoolInput is the payload accepted by the service to onboard a school.
type CreateSchoolInput struct {
	Name    string
	Code    string
	Email   string
	Phone   string
	Address string
	LogoURL string
}

// UpdateSchoolInput carries optional fields for a partial update.
type UpdateSchoolInput struct {
	Name    *string
	Email   *string
	Phone   *string
	Address *string
	LogoURL *string
	Active  *bool
}

// ListFilter paginates and searches the school listing.
type ListFilter struct {
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
