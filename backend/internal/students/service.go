package students

import (
	"context"
	"strings"
	"time"

	"github.com/skuliapp/backend/pkg/validator"
)

// Service is the business-logic contract for students.
type Service interface {
	Create(ctx context.Context, in CreateStudentInput) (*Student, error)
	GetByID(ctx context.Context, schoolID, id int64) (*Student, error)
	Update(ctx context.Context, schoolID, id int64, in UpdateStudentInput) (*Student, error)
	Delete(ctx context.Context, schoolID, id int64) error
	List(ctx context.Context, f ListFilter) ([]*Student, error)
}

type service struct {
	repo Repository
}

// NewService wires a Service to its repository.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, in CreateStudentInput) (*Student, error) {
	v := validator.New()
	validateCreate(v, in)
	if !v.Valid() {
		return nil, &validator.ValidationError{Errors: v.Errors}
	}

	enrolled := in.EnrollmentDate
	if enrolled.IsZero() {
		enrolled = time.Now()
	}

	student := &Student{
		SchoolID:        in.SchoolID,
		AdmissionNumber: strings.TrimSpace(in.AdmissionNumber),
		FirstName:       strings.TrimSpace(in.FirstName),
		LastName:        strings.TrimSpace(in.LastName),
		DateOfBirth:     in.DateOfBirth,
		Gender:          strings.TrimSpace(in.Gender),
		ClassID:         in.ClassID,
		ParentID:        in.ParentID,
		Address:         strings.TrimSpace(in.Address),
		Phone:           strings.TrimSpace(in.Phone),
		Email:           strings.ToLower(strings.TrimSpace(in.Email)),
		EnrollmentDate:  enrolled,
		Status:          StatusActive,
	}

	if err := s.repo.Insert(ctx, student); err != nil {
		return nil, err
	}
	return student, nil
}

func (s *service) GetByID(ctx context.Context, schoolID, id int64) (*Student, error) {
	return s.repo.GetByID(ctx, schoolID, id)
}

func (s *service) Update(ctx context.Context, schoolID, id int64, in UpdateStudentInput) (*Student, error) {
	student, err := s.repo.GetByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}

	if in.FirstName != nil {
		student.FirstName = strings.TrimSpace(*in.FirstName)
	}
	if in.LastName != nil {
		student.LastName = strings.TrimSpace(*in.LastName)
	}
	if in.Gender != nil {
		student.Gender = strings.TrimSpace(*in.Gender)
	}
	if in.ClassID != nil {
		student.ClassID = in.ClassID
	}
	if in.ParentID != nil {
		student.ParentID = in.ParentID
	}
	if in.Address != nil {
		student.Address = strings.TrimSpace(*in.Address)
	}
	if in.Phone != nil {
		student.Phone = strings.TrimSpace(*in.Phone)
	}
	if in.Email != nil {
		student.Email = strings.ToLower(strings.TrimSpace(*in.Email))
	}
	if in.Status != nil {
		student.Status = *in.Status
	}

	v := validator.New()
	v.Check(validator.NotBlank(student.FirstName), "first_name", "must be provided")
	v.Check(validator.NotBlank(student.LastName), "last_name", "must be provided")
	v.Check(validator.PermittedValue(student.Status, Statuses...), "status", "must be a valid status")
	if !v.Valid() {
		return nil, &validator.ValidationError{Errors: v.Errors}
	}

	if err := s.repo.Update(ctx, student); err != nil {
		return nil, err
	}
	return student, nil
}

func (s *service) Delete(ctx context.Context, schoolID, id int64) error {
	return s.repo.Delete(ctx, schoolID, id)
}

func (s *service) List(ctx context.Context, f ListFilter) ([]*Student, error) {
	return s.repo.List(ctx, f)
}

func validateCreate(v *validator.Validator, in CreateStudentInput) {
	v.Check(in.SchoolID > 0, "school_id", "must be provided")
	v.Check(validator.NotBlank(in.AdmissionNumber), "admission_number", "must be provided")
	v.Check(validator.MaxRunes(in.AdmissionNumber, 50), "admission_number", "must not exceed 50 characters")
	v.Check(validator.NotBlank(in.FirstName), "first_name", "must be provided")
	v.Check(validator.MaxRunes(in.FirstName, 100), "first_name", "must not exceed 100 characters")
	v.Check(validator.NotBlank(in.LastName), "last_name", "must be provided")
	v.Check(validator.MaxRunes(in.LastName, 100), "last_name", "must not exceed 100 characters")
	v.Check(!in.DateOfBirth.IsZero(), "date_of_birth", "must be provided")
	v.Check(in.DateOfBirth.Before(time.Now()), "date_of_birth", "must be in the past")
	if in.Email != "" {
		v.Check(validator.Matches(in.Email, validator.EmailRX), "email", "must be a valid email address")
	}
}
