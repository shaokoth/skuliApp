package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/skuliapp/backend/internal/users"
)

// routes registers every endpoint on a gorilla/mux router and wraps it with the
// panic-recovery and request-logging middleware.
func (app *application) routes() http.Handler {
	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)
	r.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	r.HandleFunc("/v1/healthcheck", app.healthcheck).Methods(http.MethodGet)

	// Authentication (public).
	r.HandleFunc("/v1/auth/login", app.login).Methods(http.MethodPost)

	// Role groups.
	admin := []users.Role{users.RoleSuperAdmin, users.RoleSchoolAdmin}
	staff := []users.Role{users.RoleSuperAdmin, users.RoleSchoolAdmin, users.RolePrincipal, users.RoleTeacher}
	superAdmin := []users.Role{users.RoleSuperAdmin}

	// Schools (platform administration).
	r.HandleFunc("/v1/schools", app.requireRole(app.createSchool, superAdmin...)).Methods(http.MethodPost)
	r.HandleFunc("/v1/schools", app.requireRole(app.listSchools, superAdmin...)).Methods(http.MethodGet)
	r.HandleFunc("/v1/schools/{id:[0-9]+}", app.requireRole(app.showSchool, superAdmin...)).Methods(http.MethodGet)
	r.HandleFunc("/v1/schools/{id:[0-9]+}", app.requireRole(app.updateSchool, superAdmin...)).Methods(http.MethodPatch)
	r.HandleFunc("/v1/schools/{id:[0-9]+}", app.requireRole(app.deleteSchool, superAdmin...)).Methods(http.MethodDelete)

	// Users.
	r.HandleFunc("/v1/users", app.requireRole(app.createUser, admin...)).Methods(http.MethodPost)
	r.HandleFunc("/v1/users", app.requireRole(app.listUsers, admin...)).Methods(http.MethodGet)
	r.HandleFunc("/v1/users/{id:[0-9]+}", app.requireAuth(app.showUser)).Methods(http.MethodGet)
	r.HandleFunc("/v1/users/{id:[0-9]+}", app.requireRole(app.updateUser, admin...)).Methods(http.MethodPatch)
	r.HandleFunc("/v1/users/{id:[0-9]+}", app.requireRole(app.deleteUser, admin...)).Methods(http.MethodDelete)

	// Students.
	r.HandleFunc("/v1/students", app.requireRole(app.createStudent, admin...)).Methods(http.MethodPost)
	r.HandleFunc("/v1/students", app.requireRole(app.listStudents, staff...)).Methods(http.MethodGet)
	r.HandleFunc("/v1/students/{id:[0-9]+}", app.requireRole(app.showStudent, staff...)).Methods(http.MethodGet)
	r.HandleFunc("/v1/students/{id:[0-9]+}", app.requireRole(app.updateStudent, admin...)).Methods(http.MethodPatch)
	r.HandleFunc("/v1/students/{id:[0-9]+}", app.requireRole(app.deleteStudent, admin...)).Methods(http.MethodDelete)

	// Teachers.
	r.HandleFunc("/v1/teachers", app.requireRole(app.createTeacher, admin...)).Methods(http.MethodPost)
	r.HandleFunc("/v1/teachers", app.requireRole(app.listTeachers, staff...)).Methods(http.MethodGet)
	r.HandleFunc("/v1/teachers/{id:[0-9]+}", app.requireRole(app.showTeacher, staff...)).Methods(http.MethodGet)
	r.HandleFunc("/v1/teachers/{id:[0-9]+}", app.requireRole(app.updateTeacher, admin...)).Methods(http.MethodPatch)
	r.HandleFunc("/v1/teachers/{id:[0-9]+}", app.requireRole(app.deleteTeacher, admin...)).Methods(http.MethodDelete)

	// Classes.
	r.HandleFunc("/v1/classes", app.requireRole(app.createClass, admin...)).Methods(http.MethodPost)
	r.HandleFunc("/v1/classes", app.requireRole(app.listClasses, staff...)).Methods(http.MethodGet)
	r.HandleFunc("/v1/classes/{id:[0-9]+}", app.requireRole(app.showClass, staff...)).Methods(http.MethodGet)
	r.HandleFunc("/v1/classes/{id:[0-9]+}", app.requireRole(app.updateClass, admin...)).Methods(http.MethodPatch)
	r.HandleFunc("/v1/classes/{id:[0-9]+}", app.requireRole(app.deleteClass, admin...)).Methods(http.MethodDelete)

	// Subjects.
	r.HandleFunc("/v1/subjects", app.requireRole(app.createSubject, admin...)).Methods(http.MethodPost)
	r.HandleFunc("/v1/subjects", app.requireRole(app.listSubjects, staff...)).Methods(http.MethodGet)
	r.HandleFunc("/v1/subjects/{id:[0-9]+}", app.requireRole(app.showSubject, staff...)).Methods(http.MethodGet)
	r.HandleFunc("/v1/subjects/{id:[0-9]+}", app.requireRole(app.updateSubject, admin...)).Methods(http.MethodPatch)
	r.HandleFunc("/v1/subjects/{id:[0-9]+}", app.requireRole(app.deleteSubject, admin...)).Methods(http.MethodDelete)

	// Attendance.
	r.HandleFunc("/v1/attendance", app.requireRole(app.createAttendance, staff...)).Methods(http.MethodPost)
	r.HandleFunc("/v1/attendance", app.requireRole(app.listAttendance, staff...)).Methods(http.MethodGet)
	r.HandleFunc("/v1/attendance/{id:[0-9]+}", app.requireRole(app.showAttendance, staff...)).Methods(http.MethodGet)
	r.HandleFunc("/v1/attendance/{id:[0-9]+}", app.requireRole(app.updateAttendance, staff...)).Methods(http.MethodPatch)
	r.HandleFunc("/v1/attendance/{id:[0-9]+}", app.requireRole(app.deleteAttendance, admin...)).Methods(http.MethodDelete)

	r.Use(app.recoverPanic, app.logRequest)
	return r
}
