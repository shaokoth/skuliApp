package reports

import "time"

// Report types.
const (
	TypeAttendanceSummary = "attendance_summary"
	TypeFinanceSummary    = "finance_summary"
	TypeExamResults       = "exam_results"
	TypeEnrollment        = "enrollment"
)

// Report is metadata for a generated report artefact.
type Report struct {
	ID          int64     `json:"id"`
	SchoolID    int64     `json:"school_id"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Params      string    `json:"params"` // JSON-encoded parameters
	GeneratedBy int64     `json:"generated_by"`
	FileURL     string    `json:"file_url"`
	CreatedAt   time.Time `json:"created_at"`
}
