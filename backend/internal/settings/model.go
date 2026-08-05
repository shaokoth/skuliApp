package settings

import "time"

// Setting is a per-school key/value configuration entry.
type Setting struct {
	ID        int64     `json:"id"`
	SchoolID  int64     `json:"school_id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}
