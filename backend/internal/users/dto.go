package users

// CreateUserInput is the payload accepted by the service to create a user.
// SchoolID is set by the handler from the authenticated tenant context.
type CreateUserInput struct {
	SchoolID  int64
	Role      string
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Password  string
}

// UpdateUserInput carries optional fields for a partial update.
type UpdateUserInput struct {
	FirstName *string
	LastName  *string
	Phone     *string
	Active    *bool
}

// Credentials is a login attempt.
type Credentials struct {
	Email    string
	Password string
}

// ListFilter narrows and paginates a user listing within a school.
type ListFilter struct {
	SchoolID int64
	Role     string
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
