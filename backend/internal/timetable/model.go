package timetable

import "time"

// Slot is one scheduled period: a subject taught to a class by a teacher.
type Slot struct {
	ID        int64     `json:"id"`
	SchoolID  int64     `json:"school_id"`
	ClassID   int64     `json:"class_id"`
	SubjectID int64     `json:"subject_id"`
	TeacherID int64     `json:"teacher_id"`
	DayOfWeek int       `json:"day_of_week"` // 1=Monday .. 7=Sunday
	StartTime string    `json:"start_time"`  // "08:00"
	EndTime   string    `json:"end_time"`    // "08:45"
	Room      string    `json:"room"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
