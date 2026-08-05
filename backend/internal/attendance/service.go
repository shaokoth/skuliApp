package attendance

import (
	"context"
	"strings"

	"github.com/skuliapp/backend/pkg/validator"
)

// Service is the business-logic contract for attendance.
type Service interface {
	Create(ctx context.Context, in CreateAttendanceInput) (*Attendance, error)
	GetByID(ctx context.Context, schoolID, id int64) (*Attendance, error)
	Update(ctx context.Context, schoolID, id int64, in UpdateAttendanceInput) (*Attendance, error)
	Delete(ctx context.Context, schoolID, id int64) error
	List(ctx context.Context, f ListFilter) ([]*Attendance, error)
}

type service struct {
	repo Repository
}

// NewService wires a Service to its repository.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, in CreateAttendanceInput) (*Attendance, error) {
	v := validator.New()
	v.Check(in.SchoolID > 0, "school_id", "must be provided")
	v.Check(in.StudentID > 0, "student_id", "must be provided")
	v.Check(in.ClassID > 0, "class_id", "must be provided")
	v.Check(!in.Date.IsZero(), "date", "must be provided")
	v.Check(validator.PermittedValue(in.Status, Statuses...), "status", "must be a valid attendance status")
	if !v.Valid() {
		return nil, &validator.ValidationError{Errors: v.Errors}
	}

	record := &Attendance{
		SchoolID:  in.SchoolID,
		StudentID: in.StudentID,
		ClassID:   in.ClassID,
		Date:      in.Date,
		Status:    in.Status,
		Remark:    strings.TrimSpace(in.Remark),
		MarkedBy:  in.MarkedBy,
	}

	if err := s.repo.Insert(ctx, record); err != nil {
		return nil, err
	}
	return record, nil
}

func (s *service) GetByID(ctx context.Context, schoolID, id int64) (*Attendance, error) {
	return s.repo.GetByID(ctx, schoolID, id)
}

func (s *service) Update(ctx context.Context, schoolID, id int64, in UpdateAttendanceInput) (*Attendance, error) {
	record, err := s.repo.GetByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}

	if in.Status != nil {
		record.Status = *in.Status
	}
	if in.Remark != nil {
		record.Remark = strings.TrimSpace(*in.Remark)
	}

	v := validator.New()
	v.Check(validator.PermittedValue(record.Status, Statuses...), "status", "must be a valid attendance status")
	if !v.Valid() {
		return nil, &validator.ValidationError{Errors: v.Errors}
	}

	if err := s.repo.Update(ctx, record); err != nil {
		return nil, err
	}
	return record, nil
}

func (s *service) Delete(ctx context.Context, schoolID, id int64) error {
	return s.repo.Delete(ctx, schoolID, id)
}

func (s *service) List(ctx context.Context, f ListFilter) ([]*Attendance, error) {
	return s.repo.List(ctx, f)
}
