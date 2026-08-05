package finance

import "time"

// Invoice status values.
const (
	InvoicePending   = "pending"
	InvoicePartial   = "partial"
	InvoicePaid      = "paid"
	InvoiceOverdue   = "overdue"
	InvoiceCancelled = "cancelled"
)

// Payment method values.
const (
	MethodCash         = "cash"
	MethodCard         = "card"
	MethodMobileMoney  = "mobile_money"
	MethodBankTransfer = "bank_transfer"
)

// FeeStructure defines a chargeable fee, optionally tied to a class/term.
// Monetary amounts are stored in minor units (e.g. cents) to avoid float error.
type FeeStructure struct {
	ID           int64     `json:"id"`
	SchoolID     int64     `json:"school_id"`
	Name         string    `json:"name"`
	ClassID      *int64    `json:"class_id,omitempty"`
	AcademicYear string    `json:"academic_year"`
	Term         string    `json:"term"`
	Amount       int64     `json:"amount"`
	Currency     string    `json:"currency"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Version      int32     `json:"-"`
}

// Invoice is a bill issued to a student for a fee structure.
type Invoice struct {
	ID             int64     `json:"id"`
	SchoolID       int64     `json:"school_id"`
	StudentID      int64     `json:"student_id"`
	FeeStructureID int64     `json:"fee_structure_id"`
	Number         string    `json:"number"`
	AmountDue      int64     `json:"amount_due"`
	AmountPaid     int64     `json:"amount_paid"`
	DueDate        time.Time `json:"due_date"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Version        int32     `json:"-"`
}

// Payment records money received against an invoice.
type Payment struct {
	ID         int64     `json:"id"`
	SchoolID   int64     `json:"school_id"`
	InvoiceID  int64     `json:"invoice_id"`
	StudentID  int64     `json:"student_id"`
	Amount     int64     `json:"amount"`
	Method     string    `json:"method"`
	Reference  string    `json:"reference"`
	PaidAt     time.Time `json:"paid_at"`
	ReceivedBy int64     `json:"received_by"`
	CreatedAt  time.Time `json:"created_at"`
}

// Receipt is issued to acknowledge a payment.
type Receipt struct {
	ID        int64     `json:"id"`
	SchoolID  int64     `json:"school_id"`
	PaymentID int64     `json:"payment_id"`
	Number    string    `json:"number"`
	IssuedAt  time.Time `json:"issued_at"`
	CreatedAt time.Time `json:"created_at"`
}
