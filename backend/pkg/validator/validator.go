// Package validator provides a small, reusable input-validation helper used by
// the service layer to collect field-level errors before touching the database.
package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

// EmailRX is a pragmatic email pattern (HTML5 spec).
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Validator collects validation errors keyed by field name.
type Validator struct {
	Errors map[string]string
}

// New returns an initialised Validator.
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid reports whether no errors have been recorded.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError records an error for a field unless one already exists.
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check adds an error message when ok is false.
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// NotBlank reports whether a string is non-empty after trimming.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxRunes reports whether a string has at most n runes.
func MaxRunes(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// MinRunes reports whether a string has at least n runes.
func MinRunes(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// PermittedValue reports whether value is one of the permitted values.
func PermittedValue[T comparable](value T, permitted ...T) bool {
	return slices.Contains(permitted, value)
}

// Matches reports whether a string matches a regular expression.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Unique reports whether all values in the slice are distinct.
func Unique[T comparable](values []T) bool {
	seen := make(map[T]bool)
	for _, value := range values {
		seen[value] = true
	}
	return len(values) == len(seen)
}
