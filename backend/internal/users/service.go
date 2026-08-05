package users

import (
	"context"
	"errors"
	"strings"

	"github.com/skuliapp/backend/pkg/validator"
)

// ErrInvalidCredentials is returned when authentication fails.
var ErrInvalidCredentials = errors.New("invalid credentials")

// Service is the business-logic contract for users. Handlers depend on this
// interface, not the concrete implementation.
type Service interface {
	Create(ctx context.Context, in CreateUserInput) (*User, error)
	GetByID(ctx context.Context, schoolID, id int64) (*User, error)
	Update(ctx context.Context, schoolID, id int64, in UpdateUserInput) (*User, error)
	Delete(ctx context.Context, schoolID, id int64) error
	List(ctx context.Context, f ListFilter) ([]*User, error)
	Authenticate(ctx context.Context, c Credentials) (*User, error)
}

type service struct {
	repo Repository
}

// NewService wires a Service to its repository.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, in CreateUserInput) (*User, error) {
	v := validator.New()
	validateCreate(v, in)
	if !v.Valid() {
		return nil, &validator.ValidationError{Errors: v.Errors}
	}

	u := &User{
		SchoolID:  in.SchoolID,
		Role:      Role(in.Role),
		FirstName: strings.TrimSpace(in.FirstName),
		LastName:  strings.TrimSpace(in.LastName),
		Email:     strings.ToLower(strings.TrimSpace(in.Email)),
		Phone:     strings.TrimSpace(in.Phone),
		Active:    true,
	}
	if err := u.SetPassword(in.Password); err != nil {
		return nil, err
	}
	if err := s.repo.Insert(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *service) GetByID(ctx context.Context, schoolID, id int64) (*User, error) {
	return s.repo.GetByID(ctx, schoolID, id)
}

func (s *service) Update(ctx context.Context, schoolID, id int64, in UpdateUserInput) (*User, error) {
	u, err := s.repo.GetByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}

	if in.FirstName != nil {
		u.FirstName = strings.TrimSpace(*in.FirstName)
	}
	if in.LastName != nil {
		u.LastName = strings.TrimSpace(*in.LastName)
	}
	if in.Phone != nil {
		u.Phone = strings.TrimSpace(*in.Phone)
	}
	if in.Active != nil {
		u.Active = *in.Active
	}

	v := validator.New()
	v.Check(validator.NotBlank(u.FirstName), "first_name", "must be provided")
	v.Check(validator.NotBlank(u.LastName), "last_name", "must be provided")
	if !v.Valid() {
		return nil, &validator.ValidationError{Errors: v.Errors}
	}

	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *service) Delete(ctx context.Context, schoolID, id int64) error {
	return s.repo.Delete(ctx, schoolID, id)
}

func (s *service) List(ctx context.Context, f ListFilter) ([]*User, error) {
	return s.repo.List(ctx, f)
}

func (s *service) Authenticate(ctx context.Context, c Credentials) (*User, error) {
	u, err := s.repo.GetByEmail(ctx, strings.ToLower(strings.TrimSpace(c.Email)))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	ok, err := u.PasswordMatches(c.Password)
	if err != nil {
		return nil, err
	}
	if !ok || !u.Active {
		return nil, ErrInvalidCredentials
	}
	return u, nil
}

func validateCreate(v *validator.Validator, in CreateUserInput) {
	v.Check(in.SchoolID > 0, "school_id", "must be provided")
	v.Check(validator.NotBlank(in.FirstName), "first_name", "must be provided")
	v.Check(validator.MaxRunes(in.FirstName, 100), "first_name", "must not exceed 100 characters")
	v.Check(validator.NotBlank(in.LastName), "last_name", "must be provided")
	v.Check(validator.MaxRunes(in.LastName, 100), "last_name", "must not exceed 100 characters")
	v.Check(validator.Matches(in.Email, validator.EmailRX), "email", "must be a valid email address")
	v.Check(validator.PermittedValue(Role(in.Role), Roles...), "role", "must be a valid role")
	v.Check(validator.NotBlank(in.Password), "password", "must be provided")
	v.Check(validator.MinRunes(in.Password, 8), "password", "must be at least 8 characters")
	v.Check(validator.MaxRunes(in.Password, 72), "password", "must not exceed 72 characters")
}
