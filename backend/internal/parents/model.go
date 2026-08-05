package parents

import "time"

// Relationship values describing a guardian's link to a student.
const (
	RelationshipFather   = "father"
	RelationshipMother   = "mother"
	RelationshipGuardian = "guardian"
)

// Parent is a guardian contact, optionally linked to a login account.
type Parent struct {
	ID           int64     `json:"id"`
	SchoolID     int64     `json:"school_id"`
	UserID       *int64    `json:"user_id,omitempty"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Occupation   string    `json:"occupation"`
	Address      string    `json:"address"`
	Relationship string    `json:"relationship"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Version      int32     `json:"-"`
}
