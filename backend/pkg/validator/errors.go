package validator

// ValidationError wraps field-level validation failures so the service layer
// can return them as a normal error and handlers can render a 422 response.
type ValidationError struct {
	Errors map[string]string
}

func (e *ValidationError) Error() string {
	return "validation failed"
}
