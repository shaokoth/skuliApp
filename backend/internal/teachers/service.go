package teachers

import (
	"context"
	"strings"
	"time"

	"github.com/skuliapp/backend/pkg/validator"
)

// Service is the business-logic contract for teachers.
type Service interface {
	Create(ctx context.Context, in CreateTeacherInput) (*Teacher, error)
	GetByID(ctx context.Context, schoolID, id int64) (*Teacher, error)
	Update(ctx context.Context, schoolID, id int64, in UpdateTeacherInput) (*Teacher, error)
	Delete(ctx context.Context, schoolID, id int64) error
	List(ctx context.Context, f ListFilter) ([]*Teacher, error)
}

type service struct {
	repo Repository
}

// NewService wires a Service to its repository.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, in CreateTeacherInput) (*Teacher, error) {
	v := validator.New()
	validateCreate(v, in)
	if !v.Valid() {
		return nil, &validator.ValidationError{Errors: v.Errors}
	}

	hired := in.HireDate
	if hired.IsZero() {
		hired = time.Now()
	}

	teacher := &Teacher{
		SchoolID:       in.SchoolID,
		EmployeeNumber: strings.TrimSpace(in.EmployeeNumber),
		FirstName:      strings.TrimSpace(in.FirstName),
		LastName:       strings.TrimSpace(in.LastName),
		Email:          strings.ToLower(strings.TrimSpace(in.Email)),
		Phone:          strings.TrimSpace(in.Phone),
		Gender:         strings.TrimSpace(in.Gender),
		Qualification:  strings.TrimSpace(in.Qualification),
		HireDate:       hired,
		Status:         StatusActive,
	}

	if err := s.repo.Insert(ctx, teacher); err != nil {
		return nil, err
	}
	return teacher, nil
}

func (s *service) GetByID(ctx context.Context, schoolID, id int64) (*Teacher, error) {
	return s.repo.GetByID(ctx, schoolID, id)
}

func (s *service) Update(ctx context.Context, schoolID, id int64, in UpdateTeacherInput) (*Teacher, error) {
	teacher, err := s.repo.GetByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}

	if in.FirstName != nil {
		teacher.FirstName = strings.TrimSpace(*in.FirstName)
	}
	if in.LastName != nil {
		teacher.LastName = strings.TrimSpace(*in.LastName)
	}
	if in.Email != nil {
		teacher.Email = strings.ToLower(strings.TrimSpace(*in.Email))
	}
	if in.Phone != nil {
		teacher.Phone = strings.TrimSpace(*in.Phone)
	}
	if in.Qualification != nil {
		teacher.Qualification = strings.TrimSpace(*in.Qualification)
	}
	if in.Status != nil {
		teacher.Status = *in.Status
	}

	v := validator.New()
	v.Check(validator.NotBlank(teacher.FirstName), "first_name", "must be provided")
	v.Check(validator.NotBlank(teacher.LastName), "last_name", "must be provided")
	v.Check(validator.PermittedValue(teacher.Status, Statuses...), "status", "must be a valid status")
	if teacher.Email != "" {
		v.Check(validator.Matches(teacher.Email, validator.EmailRX), "email", "must be a valid email address")
	}
	if !v.Valid() {
		return nil, &validator.ValidationError{Errors: v.Errors}
	}

	if err := s.repo.Update(ctx, teacher); err != nil {
		return nil, err
	}
	return teacher, nil
}

func (s *service) Delete(ctx context.Context, schoolID, id int64) error {
	return s.repo.Delete(ctx, schoolID, id)
}

func (s *service) List(ctx context.Context, f ListFilter) ([]*Teacher, error) {
	return s.repo.List(ctx, f)
}

func validateCreate(v *validator.Validator, in CreateTeacherInput) {
	v.Check(in.SchoolID > 0, "school_id", "must be provided")
	v.Check(validator.NotBlank(in.EmployeeNumber), "employee_number", "must be provided")
	v.Check(validator.MaxRunes(in.EmployeeNumber, 50), "employee_number", "must not exceed 50 characters")
	v.Check(validator.NotBlank(in.FirstName), "first_name", "must be provided")
	v.Check(validator.MaxRunes(in.FirstName, 100), "first_name", "must not exceed 100 characters")
	v.Check(validator.NotBlank(in.LastName), "last_name", "must be provided")
	v.Check(validator.MaxRunes(in.LastName, 100), "last_name", "must not exceed 100 characters")
	if in.Email != "" {
		v.Check(validator.Matches(in.Email, validator.EmailRX), "email", "must be a valid email address")
	}
}
