package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/skuliapp/backend/internal/schools"
	"github.com/skuliapp/backend/pkg/validator"
)

// createSchool onboards a new school (super admin only).
func (app *application) createSchool(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string `json:"name"`
		Code    string `json:"code"`
		Email   string `json:"email"`
		Phone   string `json:"phone"`
		Address string `json:"address"`
		LogoURL string `json:"logo_url"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	school, err := app.schools.Create(r.Context(), schools.CreateSchoolInput{
		Name:    input.Name,
		Code:    input.Code,
		Email:   input.Email,
		Phone:   input.Phone,
		Address: input.Address,
		LogoURL: input.LogoURL,
	})
	if err != nil {
		app.handleSchoolError(w, r, err)
		return
	}

	headers := http.Header{}
	headers.Set("Location", fmt.Sprintf("/v1/schools/%d", school.ID))
	if err := app.writeJSON(w, http.StatusCreated, envelope{"school": school}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showSchool returns a single school.
func (app *application) showSchool(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	school, err := app.schools.GetByID(r.Context(), id)
	if err != nil {
		app.handleSchoolError(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"school": school}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listSchools returns a filtered, paginated list of schools.
func (app *application) listSchools(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()

	filter := schools.ListFilter{
		Search:   app.readString(qs, "search", ""),
		Page:     app.readInt(qs, "page", 1),
		PageSize: app.readInt(qs, "page_size", 20),
	}

	list, err := app.schools.List(r.Context(), filter)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"schools": list}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateSchool applies a partial update to a school.
func (app *application) updateSchool(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Name    *string `json:"name"`
		Email   *string `json:"email"`
		Phone   *string `json:"phone"`
		Address *string `json:"address"`
		LogoURL *string `json:"logo_url"`
		Active  *bool   `json:"active"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	school, err := app.schools.Update(r.Context(), id, schools.UpdateSchoolInput{
		Name:    input.Name,
		Email:   input.Email,
		Phone:   input.Phone,
		Address: input.Address,
		LogoURL: input.LogoURL,
		Active:  input.Active,
	})
	if err != nil {
		app.handleSchoolError(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"school": school}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteSchool removes a school and (via ON DELETE CASCADE) its records.
func (app *application) deleteSchool(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	if err := app.schools.Delete(r.Context(), id); err != nil {
		app.handleSchoolError(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"message": "school successfully deleted"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// handleSchoolError maps service/repository errors to HTTP responses.
func (app *application) handleSchoolError(w http.ResponseWriter, r *http.Request, err error) {
	var ve *validator.ValidationError
	switch {
	case errors.As(err, &ve):
		app.failedValidationResponse(w, r, ve.Errors)
	case errors.Is(err, schools.ErrDuplicateCode):
		app.failedValidationResponse(w, r, map[string]string{"code": "a school with this code already exists"})
	case errors.Is(err, schools.ErrNotFound):
		app.notFoundResponse(w, r)
	case errors.Is(err, schools.ErrEditConflict):
		app.editConflictResponse(w, r)
	default:
		app.serverErrorResponse(w, r, err)
	}
}
