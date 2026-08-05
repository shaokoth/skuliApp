package classes

import (
	"context"
	"strings"

	"github.com/skuliapp/backend/pkg/validator"
)

// Service is the business-logic contract for classes.
type Service interface {
	Create(ctx context.Context, in CreateClassInput) (*Class, error)
	GetByID(ctx context.Context, schoolID, id int64) (*Class, error)
	Update(ctx context.Context, schoolID, id int64, in UpdateClassInput) (*Class, error)
	Delete(ctx context.Context, schoolID, id int64) error
	List(ctx context.Context, f ListFilter) ([]*Class, error)
}

type service struct {
	repo Repository
}

// NewService wires a Service to its repository.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, in CreateClassInput) (*Class, error) {
	v := validator.New()
	validateCreate(v, in)
	if !v.Valid() {
		return nil, &validator.ValidationError{Errors: v.Errors}
	}

	class := &Class{
		SchoolID:       in.SchoolID,
		Name:           strings.TrimSpace(in.Name),
		GradeLevel:     strings.TrimSpace(in.GradeLevel),
		Section:        strings.TrimSpace(in.Section),
		ClassTeacherID: in.ClassTeacherID,
		Capacity:       in.Capacity,
		AcademicYear:   strings.TrimSpace(in.AcademicYear),
	}

	if err := s.repo.Insert(ctx, class); err != nil {
		return nil, err
	}
	return class, nil
}

func (s *service) GetByID(ctx context.Context, schoolID, id int64) (*Class, error) {
	return s.repo.GetByID(ctx, schoolID, id)
}

func (s *service) Update(ctx context.Context, schoolID, id int64, in UpdateClassInput) (*Class, error) {
	class, err := s.repo.GetByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}

	if in.Name != nil {
		class.Name = strings.TrimSpace(*in.Name)
	}
	if in.GradeLevel != nil {
		class.GradeLevel = strings.TrimSpace(*in.GradeLevel)
	}
	if in.Section != nil {
		class.Section = strings.TrimSpace(*in.Section)
	}
	if in.ClassTeacherID != nil {
		class.ClassTeacherID = in.ClassTeacherID
	}
	if in.Capacity != nil {
		class.Capacity = *in.Capacity
	}
	if in.AcademicYear != nil {
		class.AcademicYear = strings.TrimSpace(*in.AcademicYear)
	}

	v := validator.New()
	v.Check(validator.NotBlank(class.Name), "name", "must be provided")
	v.Check(class.Capacity >= 0, "capacity", "must not be negative")
	if !v.Valid() {
		return nil, &validator.ValidationError{Errors: v.Errors}
	}

	if err := s.repo.Update(ctx, class); err != nil {
		return nil, err
	}
	return class, nil
}

func (s *service) Delete(ctx context.Context, schoolID, id int64) error {
	return s.repo.Delete(ctx, schoolID, id)
}

func (s *service) List(ctx context.Context, f ListFilter) ([]*Class, error) {
	return s.repo.List(ctx, f)
}

func validateCreate(v *validator.Validator, in CreateClassInput) {
	v.Check(in.SchoolID > 0, "school_id", "must be provided")
	v.Check(validator.NotBlank(in.Name), "name", "must be provided")
	v.Check(validator.MaxRunes(in.Name, 100), "name", "must not exceed 100 characters")
	v.Check(in.Capacity >= 0, "capacity", "must not be negative")
}
