package main

import (
	"log/slog"

	"github.com/skuliapp/backend/internal/attendance"
	"github.com/skuliapp/backend/internal/classes"
	"github.com/skuliapp/backend/internal/config"
	"github.com/skuliapp/backend/internal/schools"
	"github.com/skuliapp/backend/internal/students"
	"github.com/skuliapp/backend/internal/subjects"
	"github.com/skuliapp/backend/internal/teachers"
	"github.com/skuliapp/backend/internal/users"
	"github.com/skuliapp/backend/pkg/jwt"
)

// application is the dependency container. Every handler is a method on it, so
// each handler gets exactly the collaborators it needs via these fields.
type application struct {
	config config.Config
	logger *slog.Logger
	jwt    *jwt.Manager

	// service-layer dependencies (interfaces, not concrete types)
	schools    schools.Service
	users      users.Service
	students   students.Service
	teachers   teachers.Service
	classes    classes.Service
	subjects   subjects.Service
	attendance attendance.Service
}
