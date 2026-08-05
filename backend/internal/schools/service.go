package schools

import (
	"context"
	"strings"

	"github.com/skuliapp/backend/pkg/validator"
)

// Service is the business-logic contract for schools.
type Service interface {
	Create(ctx context.Context, in CreateSchoolInput) (*School, error)
	GetByID(ctx context.Context, id int64) (*School, error)
	Update(ctx context.Context, id int64, in UpdateSchoolInput) (*School, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, f ListFilter) ([]*School, error)
}

type service struct {
	repo Repository
}

// NewService wires a Service to its repository.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, in CreateSchoolInput) (*School, error) {
	v := validator.New()
	validateCreate(v, in)
	if !v.Valid() {
		return nil, &validator.ValidationError{Errors: v.Errors}
	}

	school := &School{
		Name:    strings.TrimSpace(in.Name),
		Code:    strings.ToLower(strings.TrimSpace(in.Code)),
		Email:   strings.ToLower(strings.TrimSpace(in.Email)),
		Phone:   strings.TrimSpace(in.Phone),
		Address: strings.TrimSpace(in.Address),
		LogoURL: strings.TrimSpace(in.LogoURL),
		Active:  true,
	}

	if err := s.repo.Insert(ctx, school); err != nil {
		return nil, err
	}
	return school, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (*School, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) Update(ctx context.Context, id int64, in UpdateSchoolInput) (*School, error) {
	school, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if in.Name != nil {
		school.Name = strings.TrimSpace(*in.Name)
	}
	if in.Email != nil {
		school.Email = strings.ToLower(strings.TrimSpace(*in.Email))
	}
	if in.Phone != nil {
		school.Phone = strings.TrimSpace(*in.Phone)
	}
	if in.Address != nil {
		school.Address = strings.TrimSpace(*in.Address)
	}
	if in.LogoURL != nil {
		school.LogoURL = strings.TrimSpace(*in.LogoURL)
	}
	if in.Active != nil {
		school.Active = *in.Active
	}

	v := validator.New()
	v.Check(validator.NotBlank(school.Name), "name", "must be provided")
	if school.Email != "" {
		v.Check(validator.Matches(school.Email, validator.EmailRX), "email", "must be a valid email address")
	}
	if !v.Valid() {
		return nil, &validator.ValidationError{Errors: v.Errors}
	}

	if err := s.repo.Update(ctx, school); err != nil {
		return nil, err
	}
	return school, nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) List(ctx context.Context, f ListFilter) ([]*School, error) {
	return s.repo.List(ctx, f)
}

func validateCreate(v *validator.Validator, in CreateSchoolInput) {
	v.Check(validator.NotBlank(in.Name), "name", "must be provided")
	v.Check(validator.MaxRunes(in.Name, 200), "name", "must not exceed 200 characters")
	v.Check(validator.NotBlank(in.Code), "code", "must be provided")
	v.Check(validator.MaxRunes(in.Code, 50), "code", "must not exceed 50 characters")
	if in.Email != "" {
		v.Check(validator.Matches(in.Email, validator.EmailRX), "email", "must be a valid email address")
	}
}
