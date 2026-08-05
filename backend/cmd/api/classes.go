package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/skuliapp/backend/internal/classes"
	"github.com/skuliapp/backend/pkg/validator"
)

// createClass creates a new class in the caller's school.
func (app *application) createClass(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name           string `json:"name"`
		GradeLevel     string `json:"grade_level"`
		Section        string `json:"section"`
		ClassTeacherID *int64 `json:"class_teacher_id"`
		Capacity       int    `json:"capacity"`
		AcademicYear   string `json:"academic_year"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	au, _ := app.contextGetUser(r)
	class, err := app.classes.Create(r.Context(), classes.CreateClassInput{
		SchoolID:       au.SchoolID,
		Name:           input.Name,
		GradeLevel:     input.GradeLevel,
		Section:        input.Section,
		ClassTeacherID: input.ClassTeacherID,
		Capacity:       input.Capacity,
		AcademicYear:   input.AcademicYear,
	})
	if err != nil {
		app.handleClassError(w, r, err)
		return
	}

	headers := http.Header{}
	headers.Set("Location", fmt.Sprintf("/v1/classes/%d", class.ID))
	if err := app.writeJSON(w, http.StatusCreated, envelope{"class": class}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showClass returns a single class in the caller's school.
func (app *application) showClass(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	au, _ := app.contextGetUser(r)
	class, err := app.classes.GetByID(r.Context(), au.SchoolID, id)
	if err != nil {
		app.handleClassError(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"class": class}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listClasses returns a filtered, paginated list of classes.
func (app *application) listClasses(w http.ResponseWriter, r *http.Request) {
	au, _ := app.contextGetUser(r)
	qs := r.URL.Query()

	filter := classes.ListFilter{
		SchoolID:     au.SchoolID,
		AcademicYear: app.readString(qs, "academic_year", ""),
		Search:       app.readString(qs, "search", ""),
		Page:         app.readInt(qs, "page", 1),
		PageSize:     app.readInt(qs, "page_size", 20),
	}

	list, err := app.classes.List(r.Context(), filter)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"classes": list}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateClass applies a partial update to a class.
func (app *application) updateClass(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Name           *string `json:"name"`
		GradeLevel     *string `json:"grade_level"`
		Section        *string `json:"section"`
		ClassTeacherID *int64  `json:"class_teacher_id"`
		Capacity       *int    `json:"capacity"`
		AcademicYear   *string `json:"academic_year"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	au, _ := app.contextGetUser(r)
	class, err := app.classes.Update(r.Context(), au.SchoolID, id, classes.UpdateClassInput{
		Name:           input.Name,
		GradeLevel:     input.GradeLevel,
		Section:        input.Section,
		ClassTeacherID: input.ClassTeacherID,
		Capacity:       input.Capacity,
		AcademicYear:   input.AcademicYear,
	})
	if err != nil {
		app.handleClassError(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"class": class}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteClass removes a class from the caller's school.
func (app *application) deleteClass(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	au, _ := app.contextGetUser(r)
	if err := app.classes.Delete(r.Context(), au.SchoolID, id); err != nil {
		app.handleClassError(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"message": "class successfully deleted"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// handleClassError maps service/repository errors to HTTP responses.
func (app *application) handleClassError(w http.ResponseWriter, r *http.Request, err error) {
	var ve *validator.ValidationError
	switch {
	case errors.As(err, &ve):
		app.failedValidationResponse(w, r, ve.Errors)
	case errors.Is(err, classes.ErrDuplicateClass):
		app.failedValidationResponse(w, r, map[string]string{"name": "a class with this name already exists for the academic year"})
	case errors.Is(err, classes.ErrNotFound):
		app.notFoundResponse(w, r)
	case errors.Is(err, classes.ErrEditConflict):
		app.editConflictResponse(w, r)
	default:
		app.serverErrorResponse(w, r, err)
	}
}
