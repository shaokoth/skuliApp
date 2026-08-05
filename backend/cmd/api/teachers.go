package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/skuliapp/backend/internal/teachers"
	"github.com/skuliapp/backend/pkg/validator"
)

// createTeacher adds a new teacher to the caller's school.
func (app *application) createTeacher(w http.ResponseWriter, r *http.Request) {
	var input struct {
		EmployeeNumber string `json:"employee_number"`
		FirstName      string `json:"first_name"`
		LastName       string `json:"last_name"`
		Email          string `json:"email"`
		Phone          string `json:"phone"`
		Gender         string `json:"gender"`
		Qualification  string `json:"qualification"`
		HireDate       string `json:"hire_date"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	hired, ok := parseDate(input.HireDate)
	if !ok {
		app.failedValidationResponse(w, r, map[string]string{"hire_date": "must be a valid date (YYYY-MM-DD)"})
		return
	}

	au, _ := app.contextGetUser(r)
	teacher, err := app.teachers.Create(r.Context(), teachers.CreateTeacherInput{
		SchoolID:       au.SchoolID,
		EmployeeNumber: input.EmployeeNumber,
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Email:          input.Email,
		Phone:          input.Phone,
		Gender:         input.Gender,
		Qualification:  input.Qualification,
		HireDate:       hired,
	})
	if err != nil {
		app.handleTeacherError(w, r, err)
		return
	}

	headers := http.Header{}
	headers.Set("Location", fmt.Sprintf("/v1/teachers/%d", teacher.ID))
	if err := app.writeJSON(w, http.StatusCreated, envelope{"teacher": teacher}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showTeacher returns a single teacher in the caller's school.
func (app *application) showTeacher(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	au, _ := app.contextGetUser(r)
	teacher, err := app.teachers.GetByID(r.Context(), au.SchoolID, id)
	if err != nil {
		app.handleTeacherError(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"teacher": teacher}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listTeachers returns a filtered, paginated list of teachers.
func (app *application) listTeachers(w http.ResponseWriter, r *http.Request) {
	au, _ := app.contextGetUser(r)
	qs := r.URL.Query()

	filter := teachers.ListFilter{
		SchoolID: au.SchoolID,
		Status:   app.readString(qs, "status", ""),
		Search:   app.readString(qs, "search", ""),
		Page:     app.readInt(qs, "page", 1),
		PageSize: app.readInt(qs, "page_size", 20),
	}

	list, err := app.teachers.List(r.Context(), filter)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"teachers": list}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateTeacher applies a partial update to a teacher.
func (app *application) updateTeacher(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		FirstName     *string `json:"first_name"`
		LastName      *string `json:"last_name"`
		Email         *string `json:"email"`
		Phone         *string `json:"phone"`
		Qualification *string `json:"qualification"`
		Status        *string `json:"status"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	au, _ := app.contextGetUser(r)
	teacher, err := app.teachers.Update(r.Context(), au.SchoolID, id, teachers.UpdateTeacherInput{
		FirstName:     input.FirstName,
		LastName:      input.LastName,
		Email:         input.Email,
		Phone:         input.Phone,
		Qualification: input.Qualification,
		Status:        input.Status,
	})
	if err != nil {
		app.handleTeacherError(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"teacher": teacher}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteTeacher removes a teacher from the caller's school.
func (app *application) deleteTeacher(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	au, _ := app.contextGetUser(r)
	if err := app.teachers.Delete(r.Context(), au.SchoolID, id); err != nil {
		app.handleTeacherError(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"message": "teacher successfully deleted"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// handleTeacherError maps service/repository errors to HTTP responses.
func (app *application) handleTeacherError(w http.ResponseWriter, r *http.Request, err error) {
	var ve *validator.ValidationError
	switch {
	case errors.As(err, &ve):
		app.failedValidationResponse(w, r, ve.Errors)
	case errors.Is(err, teachers.ErrDuplicateEmployeeNo):
		app.failedValidationResponse(w, r, map[string]string{"employee_number": "a teacher with this employee number already exists"})
	case errors.Is(err, teachers.ErrNotFound):
		app.notFoundResponse(w, r)
	case errors.Is(err, teachers.ErrEditConflict):
		app.editConflictResponse(w, r)
	default:
		app.serverErrorResponse(w, r, err)
	}
}
