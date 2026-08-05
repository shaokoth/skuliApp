package exams

import "time"

// Exam is an assessment event (e.g. a mid-term) within an academic year.
type Exam struct {
	ID           int64     `json:"id"`
	SchoolID     int64     `json:"school_id"`
	Name         string    `json:"name"`
	Term         string    `json:"term"`
	AcademicYear string    `json:"academic_year"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Published    bool      `json:"published"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Version      int32     `json:"-"`
}
