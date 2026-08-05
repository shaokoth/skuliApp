package main

import (
	"flag"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/skuliapp/backend/internal/attendance"
	"github.com/skuliapp/backend/internal/classes"
	"github.com/skuliapp/backend/internal/config"
	"github.com/skuliapp/backend/internal/schools"
	"github.com/skuliapp/backend/internal/students"
	"github.com/skuliapp/backend/internal/subjects"
	"github.com/skuliapp/backend/internal/teachers"
	"github.com/skuliapp/backend/internal/users"
	"github.com/skuliapp/backend/migrations"
	"github.com/skuliapp/backend/pkg/database"
	"github.com/skuliapp/backend/pkg/jwt"
)

func main() {
	var cfg config.Config

	// Render (and most PaaS) inject the listen port via the PORT env var.
	defaultPort := 4000
	if p, err := strconv.Atoi(os.Getenv("PORT")); err == nil && p > 0 {
		defaultPort = p
	}
	autoMigrate := os.Getenv("SKULI_AUTO_MIGRATE") == "true"

	flag.IntVar(&cfg.Port, "port", defaultPort, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.DB.DSN, "db-dsn", os.Getenv("SKULI_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.DB.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DB.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.DB.MaxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.StringVar(&cfg.JWT.Secret, "jwt-secret", os.Getenv("SKULI_JWT_SECRET"), "JWT signing secret")
	flag.StringVar(&cfg.JWT.Issuer, "jwt-issuer", "skuliapp", "JWT issuer")
	flag.DurationVar(&cfg.JWT.TTL, "jwt-ttl", 24*time.Hour, "JWT time-to-live")
	flag.BoolVar(&autoMigrate, "auto-migrate", autoMigrate, "run pending DB migrations on startup")

	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if cfg.JWT.Secret == "" {
		logger.Error("missing JWT secret (set -jwt-secret flag or SKULI_JWT_SECRET env var)")
		os.Exit(1)
	}

	db, err := database.Open(cfg.DB)
	if err != nil {
		logger.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	if autoMigrate {
		logger.Info("running database migrations")
		if err := migrations.Up(db); err != nil {
			logger.Error("migrations failed", "error", err)
			os.Exit(1)
		}
		logger.Info("database migrations up to date")
	}

	app := &application{
		config:   cfg,
		logger:   logger,
		jwt:      jwt.NewManager(cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.TTL),
		schools:  schools.NewService(schools.NewRepository(db)),
		users:    users.NewService(users.NewRepository(db)),
		students: students.NewService(students.NewRepository(db)),
		teachers: teachers.NewService(teachers.NewRepository(db)),

		classes:    classes.NewService(classes.NewRepository(db)),
		subjects:   subjects.NewService(subjects.NewRepository(db)),
		attendance: attendance.NewService(attendance.NewRepository(db)),
	}

	if err := app.serve(); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
