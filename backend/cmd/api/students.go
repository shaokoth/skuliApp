package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/skuliapp/backend/internal/students"
	"github.com/skuliapp/backend/pkg/validator"
)

const dateLayout = "2006-01-02"

// parseDate parses an optional YYYY-MM-DD date; empty is allowed (zero time).
func parseDate(s string) (time.Time, bool) {
	if s == "" {
		return time.Time{}, true
	}
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

// createStudent admits a new student to the caller's school.
func (app *application) createStudent(w http.ResponseWriter, r *http.Request) {
	var input struct {
		AdmissionNumber string `json:"admission_number"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		DateOfBirth     string `json:"date_of_birth"`
		Gender          string `json:"gender"`
		ClassID         *int64 `json:"class_id"`
		ParentID        *int64 `json:"parent_id"`
		Address         string `json:"address"`
		Phone           string `json:"phone"`
		Email           string `json:"email"`
		EnrollmentDate  string `json:"enrollment_date"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	dob, ok := parseDate(input.DateOfBirth)
	if !ok {
		app.failedValidationResponse(w, r, map[string]string{"date_of_birth": "must be a valid date (YYYY-MM-DD)"})
		return
	}
	enrolled, ok := parseDate(input.EnrollmentDate)
	if !ok {
		app.failedValidationResponse(w, r, map[string]string{"enrollment_date": "must be a valid date (YYYY-MM-DD)"})
		return
	}

	au, _ := app.contextGetUser(r)
	student, err := app.students.Create(r.Context(), students.CreateStudentInput{
		SchoolID:        au.SchoolID,
		AdmissionNumber: input.AdmissionNumber,
		FirstName:       input.FirstName,
		LastName:        input.LastName,
		DateOfBirth:     dob,
		Gender:          input.Gender,
		ClassID:         input.ClassID,
		ParentID:        input.ParentID,
		Address:         input.Address,
		Phone:           input.Phone,
		Email:           input.Email,
		EnrollmentDate:  enrolled,
	})
	if err != nil {
		app.handleStudentError(w, r, err)
		return
	}

	headers := http.Header{}
	headers.Set("Location", fmt.Sprintf("/v1/students/%d", student.ID))
	if err := app.writeJSON(w, http.StatusCreated, envelope{"student": student}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showStudent returns a single student in the caller's school.
func (app *application) showStudent(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	au, _ := app.contextGetUser(r)
	student, err := app.students.GetByID(r.Context(), au.SchoolID, id)
	if err != nil {
		app.handleStudentError(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"student": student}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listStudents returns a filtered, paginated list of students.
func (app *application) listStudents(w http.ResponseWriter, r *http.Request) {
	au, _ := app.contextGetUser(r)
	qs := r.URL.Query()

	filter := students.ListFilter{
		SchoolID: au.SchoolID,
		ClassID:  app.readOptionalInt64(qs, "class_id"),
		Status:   app.readString(qs, "status", ""),
		Search:   app.readString(qs, "search", ""),
		Page:     app.readInt(qs, "page", 1),
		PageSize: app.readInt(qs, "page_size", 20),
	}

	list, err := app.students.List(r.Context(), filter)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"students": list}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateStudent applies a partial update to a student.
func (app *application) updateStudent(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Gender    *string `json:"gender"`
		ClassID   *int64  `json:"class_id"`
		ParentID  *int64  `json:"parent_id"`
		Address   *string `json:"address"`
		Phone     *string `json:"phone"`
		Email     *string `json:"email"`
		Status    *string `json:"status"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	au, _ := app.contextGetUser(r)
	student, err := app.students.Update(r.Context(), au.SchoolID, id, students.UpdateStudentInput{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Gender:    input.Gender,
		ClassID:   input.ClassID,
		ParentID:  input.ParentID,
		Address:   input.Address,
		Phone:     input.Phone,
		Email:     input.Email,
		Status:    input.Status,
	})
	if err != nil {
		app.handleStudentError(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"student": student}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteStudent removes a student from the caller's school.
func (app *application) deleteStudent(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	au, _ := app.contextGetUser(r)
	if err := app.students.Delete(r.Context(), au.SchoolID, id); err != nil {
		app.handleStudentError(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"message": "student successfully deleted"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// handleStudentError maps service/repository errors to HTTP responses.
func (app *application) handleStudentError(w http.ResponseWriter, r *http.Request, err error) {
	var ve *validator.ValidationError
	switch {
	case errors.As(err, &ve):
		app.failedValidationResponse(w, r, ve.Errors)
	case errors.Is(err, students.ErrDuplicateAdmissionNo):
		app.failedValidationResponse(w, r, map[string]string{"admission_number": "a student with this admission number already exists"})
	case errors.Is(err, students.ErrNotFound):
		app.notFoundResponse(w, r)
	case errors.Is(err, students.ErrEditConflict):
		app.editConflictResponse(w, r)
	default:
		app.serverErrorResponse(w, r, err)
	}
}
