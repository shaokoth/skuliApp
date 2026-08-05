package library

import "time"

// Book loan status values.
const (
	LoanBorrowed = "borrowed"
	LoanReturned = "returned"
	LoanOverdue  = "overdue"
	LoanLost     = "lost"
)

// Book is a catalogue title with a number of physical copies.
type Book struct {
	ID        int64     `json:"id"`
	SchoolID  int64     `json:"school_id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	ISBN      string    `json:"isbn"`
	Category  string    `json:"category"`
	Copies    int       `json:"copies"`
	Available int       `json:"available"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int32     `json:"-"`
}

// BookLoan tracks a copy borrowed by a user.
type BookLoan struct {
	ID         int64      `json:"id"`
	SchoolID   int64      `json:"school_id"`
	BookID     int64      `json:"book_id"`
	BorrowerID int64      `json:"borrower_id"`
	BorrowedAt time.Time  `json:"borrowed_at"`
	DueDate    time.Time  `json:"due_date"`
	ReturnedAt *time.Time `json:"returned_at,omitempty"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
