package grading

import "time"

// Mark is a student's score for one subject in one exam.
type Mark struct {
	ID        int64     `json:"id"`
	SchoolID  int64     `json:"school_id"`
	ExamID    int64     `json:"exam_id"`
	StudentID int64     `json:"student_id"`
	SubjectID int64     `json:"subject_id"`
	Score     float64   `json:"score"`
	MaxScore  float64   `json:"max_score"`
	Grade     string    `json:"grade"`
	Remark    string    `json:"remark"`
	EnteredBy int64     `json:"entered_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int32     `json:"-"`
}

// ReportCard aggregates a student's performance for an exam, with GPA and rank.
type ReportCard struct {
	ID         int64     `json:"id"`
	SchoolID   int64     `json:"school_id"`
	StudentID  int64     `json:"student_id"`
	ExamID     int64     `json:"exam_id"`
	TotalScore float64   `json:"total_score"`
	Average    float64   `json:"average"`
	GPA        float64   `json:"gpa"`
	Grade      string    `json:"grade"`
	Rank       int       `json:"rank"`
	Remark     string    `json:"remark"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// GradeScale defines a grade band and its GPA points for a school.
type GradeScale struct {
	ID       int64   `json:"id"`
	SchoolID int64   `json:"school_id"`
	Grade    string  `json:"grade"`
	MinScore float64 `json:"min_score"`
	MaxScore float64 `json:"max_score"`
	Points   float64 `json:"points"`
}
