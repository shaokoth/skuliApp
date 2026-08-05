package classes

import "time"

// Class is a group of students in a grade/section for an academic year.
type Class struct {
	ID             int64     `json:"id"`
	SchoolID       int64     `json:"school_id"`
	Name           string    `json:"name"`
	GradeLevel     string    `json:"grade_level"`
	Section        string    `json:"section"`
	ClassTeacherID *int64    `json:"class_teacher_id,omitempty"`
	Capacity       int       `json:"capacity"`
	AcademicYear   string    `json:"academic_year"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Version        int32     `json:"-"`
}
