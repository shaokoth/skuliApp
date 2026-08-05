package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/skuliapp/backend/internal/attendance"
	"github.com/skuliapp/backend/pkg/validator"
)

// createAttendance marks a student's attendance for a class on a date.
func (app *application) createAttendance(w http.ResponseWriter, r *http.Request) {
	var input struct {
		StudentID int64  `json:"student_id"`
		ClassID   int64  `json:"class_id"`
		Date      string `json:"date"`
		Status    string `json:"status"`
		Remark    string `json:"remark"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	date, ok := parseDate(input.Date)
	if !ok {
		app.failedValidationResponse(w, r, map[string]string{"date": "must be a valid date (YYYY-MM-DD)"})
		return
	}

	au, _ := app.contextGetUser(r)
	markedBy := au.UserID
	record, err := app.attendance.Create(r.Context(), attendance.CreateAttendanceInput{
		SchoolID:  au.SchoolID,
		StudentID: input.StudentID,
		ClassID:   input.ClassID,
		Date:      date,
		Status:    input.Status,
		Remark:    input.Remark,
		MarkedBy:  &markedBy,
	})
	if err != nil {
		app.handleAttendanceError(w, r, err)
		return
	}

	headers := http.Header{}
	headers.Set("Location", fmt.Sprintf("/v1/attendance/%d", record.ID))
	if err := app.writeJSON(w, http.StatusCreated, envelope{"attendance": record}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showAttendance returns a single attendance record.
func (app *application) showAttendance(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	au, _ := app.contextGetUser(r)
	record, err := app.attendance.GetByID(r.Context(), au.SchoolID, id)
	if err != nil {
		app.handleAttendanceError(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"attendance": record}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listAttendance returns filtered, paginated attendance records.
func (app *application) listAttendance(w http.ResponseWriter, r *http.Request) {
	au, _ := app.contextGetUser(r)
	qs := r.URL.Query()

	var datePtr *time.Time
	if ds := app.readString(qs, "date", ""); ds != "" {
		d, ok := parseDate(ds)
		if !ok {
			app.failedValidationResponse(w, r, map[string]string{"date": "must be a valid date (YYYY-MM-DD)"})
			return
		}
		datePtr = &d
	}

	filter := attendance.ListFilter{
		SchoolID:  au.SchoolID,
		ClassID:   app.readOptionalInt64(qs, "class_id"),
		StudentID: app.readOptionalInt64(qs, "student_id"),
		Status:    app.readString(qs, "status", ""),
		Date:      datePtr,
		Page:      app.readInt(qs, "page", 1),
		PageSize:  app.readInt(qs, "page_size", 50),
	}

	list, err := app.attendance.List(r.Context(), filter)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"attendance": list}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateAttendance corrects a mark's status/remark.
func (app *application) updateAttendance(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Status *string `json:"status"`
		Remark *string `json:"remark"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	au, _ := app.contextGetUser(r)
	record, err := app.attendance.Update(r.Context(), au.SchoolID, id, attendance.UpdateAttendanceInput{
		Status: input.Status,
		Remark: input.Remark,
	})
	if err != nil {
		app.handleAttendanceError(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"attendance": record}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteAttendance removes an attendance record.
func (app *application) deleteAttendance(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	au, _ := app.contextGetUser(r)
	if err := app.attendance.Delete(r.Context(), au.SchoolID, id); err != nil {
		app.handleAttendanceError(w, r, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"message": "attendance record successfully deleted"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// handleAttendanceError maps service/repository errors to HTTP responses.
func (app *application) handleAttendanceError(w http.ResponseWriter, r *http.Request, err error) {
	var ve *validator.ValidationError
	switch {
	case errors.As(err, &ve):
		app.failedValidationResponse(w, r, ve.Errors)
	case errors.Is(err, attendance.ErrDuplicateMark):
		app.failedValidationResponse(w, r, map[string]string{"date": "attendance already marked for this student on this date"})
	case errors.Is(err, attendance.ErrNotFound):
		app.notFoundResponse(w, r)
	default:
		app.serverErrorResponse(w, r, err)
	}
}
