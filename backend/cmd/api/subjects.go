package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/skuliapp/backend/internal/subjects"
	"github.com/skuliapp/backend/pkg/validator"
)

// createSubject creates a new subject in the caller's school.
func (app *application) createSubject(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
		Code string `json:"code"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	au, _ := app.contextGetUser(r)
	subject, err := app.subjects.Create(r.Context(), subjects.CreateSubjectInput{
		SchoolID: au.SchoolID,
		Name:     input.Name,
		Code:     input.Code,
	})
	if err != nil {
		app.handleSubjectError(w, r, err)
		return
	}

	headers := http.Header{}
	headers.Set("Location", fmt.Sprintf("/v1/subjects/%d", subject.ID))
	if err := app.writeJSON(w, http.StatusCreated, envelope{"subject": subject}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showSubject returns a single subject in the caller's school.
func (app *application) showSubject(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	au, _ := app.contextGetUser(r)
	subject, err := app.subjects.GetByID(r.Context(), au.SchoolID, id)
	if err != nil {
		app.handleSubjectError(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"subject": subject}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listSubjects returns a filtered, paginated list of subjects.
func (app *application) listSubjects(w http.ResponseWriter, r *http.Request) {
	au, _ := app.contextGetUser(r)
	qs := r.URL.Query()

	filter := subjects.ListFilter{
		SchoolID: au.SchoolID,
		Search:   app.readString(qs, "search", ""),
		Page:     app.readInt(qs, "page", 1),
		PageSize: app.readInt(qs, "page_size", 20),
	}

	list, err := app.subjects.List(r.Context(), filter)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"subjects": list}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateSubject applies a partial update to a subject.
func (app *application) updateSubject(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Name *string `json:"name"`
		Code *string `json:"code"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	au, _ := app.contextGetUser(r)
	subject, err := app.subjects.Update(r.Context(), au.SchoolID, id, subjects.UpdateSubjectInput{
		Name: input.Name,
		Code: input.Code,
	})
	if err != nil {
		app.handleSubjectError(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"subject": subject}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteSubject removes a subject from the caller's school.
func (app *application) deleteSubject(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	au, _ := app.contextGetUser(r)
	if err := app.subjects.Delete(r.Context(), au.SchoolID, id); err != nil {
		app.handleSubjectError(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"message": "subject successfully deleted"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// handleSubjectError maps service/repository errors to HTTP responses.
func (app *application) handleSubjectError(w http.ResponseWriter, r *http.Request, err error) {
	var ve *validator.ValidationError
	switch {
	case errors.As(err, &ve):
		app.failedValidationResponse(w, r, ve.Errors)
	case errors.Is(err, subjects.ErrDuplicateCode):
		app.failedValidationResponse(w, r, map[string]string{"code": "a subject with this code already exists"})
	case errors.Is(err, subjects.ErrNotFound):
		app.notFoundResponse(w, r)
	case errors.Is(err, subjects.ErrEditConflict):
		app.editConflictResponse(w, r)
	default:
		app.serverErrorResponse(w, r, err)
	}
}
