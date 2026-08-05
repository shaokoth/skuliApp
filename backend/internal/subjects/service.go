package subjects

import (
	"context"
	"strings"

	"github.com/skuliapp/backend/pkg/validator"
)

// Service is the business-logic contract for subjects.
type Service interface {
	Create(ctx context.Context, in CreateSubjectInput) (*Subject, error)
	GetByID(ctx context.Context, schoolID, id int64) (*Subject, error)
	Update(ctx context.Context, schoolID, id int64, in UpdateSubjectInput) (*Subject, error)
	Delete(ctx context.Context, schoolID, id int64) error
	List(ctx context.Context, f ListFilter) ([]*Subject, error)
}

type service struct {
	repo Repository
}

// NewService wires a Service to its repository.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, in CreateSubjectInput) (*Subject, error) {
	v := validator.New()
	v.Check(in.SchoolID > 0, "school_id", "must be provided")
	v.Check(validator.NotBlank(in.Name), "name", "must be provided")
	v.Check(validator.MaxRunes(in.Name, 100), "name", "must not exceed 100 characters")
	v.Check(validator.MaxRunes(in.Code, 20), "code", "must not exceed 20 characters")
	if !v.Valid() {
		return nil, &validator.ValidationError{Errors: v.Errors}
	}

	subject := &Subject{
		SchoolID: in.SchoolID,
		Name:     strings.TrimSpace(in.Name),
		Code:     strings.ToUpper(strings.TrimSpace(in.Code)),
	}

	if err := s.repo.Insert(ctx, subject); err != nil {
		return nil, err
	}
	return subject, nil
}

func (s *service) GetByID(ctx context.Context, schoolID, id int64) (*Subject, error) {
	return s.repo.GetByID(ctx, schoolID, id)
}

func (s *service) Update(ctx context.Context, schoolID, id int64, in UpdateSubjectInput) (*Subject, error) {
	subject, err := s.repo.GetByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}

	if in.Name != nil {
		subject.Name = strings.TrimSpace(*in.Name)
	}
	if in.Code != nil {
		subject.Code = strings.ToUpper(strings.TrimSpace(*in.Code))
	}

	v := validator.New()
	v.Check(validator.NotBlank(subject.Name), "name", "must be provided")
	if !v.Valid() {
		return nil, &validator.ValidationError{Errors: v.Errors}
	}

	if err := s.repo.Update(ctx, subject); err != nil {
		return nil, err
	}
	return subject, nil
}

func (s *service) Delete(ctx context.Context, schoolID, id int64) error {
	return s.repo.Delete(ctx, schoolID, id)
}

func (s *service) List(ctx context.Context, f ListFilter) ([]*Subject, error) {
	return s.repo.List(ctx, f)
}
