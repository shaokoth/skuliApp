package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/skuliapp/backend/internal/users"
	"github.com/skuliapp/backend/pkg/validator"
)

// login authenticates a user and returns a signed JWT.
func (app *application) login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.users.Authenticate(r.Context(), users.Credentials{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		if errors.Is(err, users.ErrInvalidCredentials) {
			app.invalidCredentialsResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	token, expiresAt, err := app.jwt.Generate(user.ID, user.SchoolID, string(user.Role))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"token":      token,
		"expires_at": expiresAt.Format(time.RFC3339),
		"user":       user,
	}
	if err := app.writeJSON(w, http.StatusOK, env, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// createUser provisions a new user within the caller's school.
func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Role      string `json:"role"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		Password  string `json:"password"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	au, _ := app.contextGetUser(r)
	user, err := app.users.Create(r.Context(), users.CreateUserInput{
		SchoolID:  au.SchoolID,
		Role:      input.Role,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Phone:     input.Phone,
		Password:  input.Password,
	})
	if err != nil {
		app.handleUserError(w, r, err)
		return
	}

	headers := http.Header{}
	headers.Set("Location", fmt.Sprintf("/v1/users/%d", user.ID))
	if err := app.writeJSON(w, http.StatusCreated, envelope{"user": user}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showUser returns a single user in the caller's school.
func (app *application) showUser(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	au, _ := app.contextGetUser(r)
	user, err := app.users.GetByID(r.Context(), au.SchoolID, id)
	if err != nil {
		app.handleUserError(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listUsers returns a filtered, paginated list of users.
func (app *application) listUsers(w http.ResponseWriter, r *http.Request) {
	au, _ := app.contextGetUser(r)
	qs := r.URL.Query()

	filter := users.ListFilter{
		SchoolID: au.SchoolID,
		Role:     app.readString(qs, "role", ""),
		Search:   app.readString(qs, "search", ""),
		Page:     app.readInt(qs, "page", 1),
		PageSize: app.readInt(qs, "page_size", 20),
	}

	list, err := app.users.List(r.Context(), filter)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"users": list}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateUser applies a partial update to a user.
func (app *application) updateUser(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Phone     *string `json:"phone"`
		Active    *bool   `json:"active"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	au, _ := app.contextGetUser(r)
	user, err := app.users.Update(r.Context(), au.SchoolID, id, users.UpdateUserInput{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.Phone,
		Active:    input.Active,
	})
	if err != nil {
		app.handleUserError(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteUser removes a user from the caller's school.
func (app *application) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	au, _ := app.contextGetUser(r)
	if err := app.users.Delete(r.Context(), au.SchoolID, id); err != nil {
		app.handleUserError(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"message": "user successfully deleted"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// handleUserError maps service/repository errors to HTTP responses.
func (app *application) handleUserError(w http.ResponseWriter, r *http.Request, err error) {
	var ve *validator.ValidationError
	switch {
	case errors.As(err, &ve):
		app.failedValidationResponse(w, r, ve.Errors)
	case errors.Is(err, users.ErrDuplicateEmail):
		app.failedValidationResponse(w, r, map[string]string{"email": "a user with this email address already exists"})
	case errors.Is(err, users.ErrNotFound):
		app.notFoundResponse(w, r)
	case errors.Is(err, users.ErrEditConflict):
		app.editConflictResponse(w, r)
	default:
		app.serverErrorResponse(w, r, err)
	}
}
