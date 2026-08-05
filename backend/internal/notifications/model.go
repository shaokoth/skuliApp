package notifications

import "time"

// Delivery channels.
const (
	ChannelInApp = "in_app"
	ChannelEmail = "email"
	ChannelSMS   = "sms"
)

// Announcement audiences.
const (
	AudienceAll      = "all"
	AudienceTeachers = "teachers"
	AudienceStudents = "students"
	AudienceParents  = "parents"
)

// Notification is an in-app/email/SMS message targeted at one user.
type Notification struct {
	ID        int64     `json:"id"`
	SchoolID  int64     `json:"school_id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Channel   string    `json:"channel"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"created_at"`
}

// Announcement is a broadcast message to an audience within a school.
type Announcement struct {
	ID          int64      `json:"id"`
	SchoolID    int64      `json:"school_id"`
	Title       string     `json:"title"`
	Body        string     `json:"body"`
	Audience    string     `json:"audience"`
	CreatedBy   int64      `json:"created_by"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Message logs an outbound SMS/email dispatch and its delivery status.
type Message struct {
	ID        int64     `json:"id"`
	SchoolID  int64     `json:"school_id"`
	Channel   string    `json:"channel"`
	Recipient string    `json:"recipient"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	Status    string    `json:"status"` // queued, sent, failed
	CreatedAt time.Time `json:"created_at"`
}
