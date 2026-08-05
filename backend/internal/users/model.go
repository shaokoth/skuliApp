package users

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Role enumerates the platform RBAC roles.
type Role string

const (
	RoleSuperAdmin  Role = "super_admin"
	RoleSchoolAdmin Role = "school_admin"
	RolePrincipal   Role = "principal"
	RoleTeacher     Role = "teacher"
	RoleAccountant  Role = "accountant"
	RoleParent      Role = "parent"
	RoleStudent     Role = "student"
)

// Roles is the full set of valid roles, used for validation.
var Roles = []Role{
	RoleSuperAdmin, RoleSchoolAdmin, RolePrincipal,
	RoleTeacher, RoleAccountant, RoleParent, RoleStudent,
}

// User is a login-capable account scoped to a single school (tenant).
type User struct {
	ID           int64      `json:"id"`
	SchoolID     int64      `json:"school_id"`
	Role         Role       `json:"role"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Email        string     `json:"email"`
	Phone        string     `json:"phone"`
	PasswordHash []byte     `json:"-"`
	Active       bool       `json:"active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Version      int32      `json:"-"`
}

// SetPassword hashes the plaintext password and stores it on the user.
func (u *User) SetPassword(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), 12)
	if err != nil {
		return err
	}
	u.PasswordHash = hash
	return nil
}

// PasswordMatches reports whether the plaintext matches the stored hash.
func (u *User) PasswordMatches(plaintext string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(plaintext))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
