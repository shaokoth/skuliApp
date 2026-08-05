package subjects

import "time"

// Subject is a teachable subject offered by a school.
type Subject struct {
	ID        int64     `json:"id"`
	SchoolID  int64     `json:"school_id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int32     `json:"-"`
}
