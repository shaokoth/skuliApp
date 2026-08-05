package attendance

import "time"

// Attendance mark values.
const (
	StatusPresent = "present"
	StatusAbsent  = "absent"
	StatusLate    = "late"
	StatusExcused = "excused"
)

// Statuses lists the valid attendance marks.
var Statuses = []string{StatusPresent, StatusAbsent, StatusLate, StatusExcused}

// Attendance records a single student's presence for a class on a date.
type Attendance struct {
	ID        int64     `json:"id"`
	SchoolID  int64     `json:"school_id"`
	StudentID int64     `json:"student_id"`
	ClassID   int64     `json:"class_id"`
	Date      time.Time `json:"date"`
	Status    string    `json:"status"`
	Remark    string    `json:"remark"`
	MarkedBy  *int64    `json:"marked_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
