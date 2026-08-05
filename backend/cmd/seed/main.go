// Command seed bootstraps a school and its first super-admin user so the API
// can be used end-to-end. Run once against a freshly migrated database.
package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/skuliapp/backend/internal/config"
	"github.com/skuliapp/backend/internal/schools"
	"github.com/skuliapp/backend/internal/users"
	"github.com/skuliapp/backend/pkg/database"
)

func main() {
	var (
		dsn        string
		schoolName string
		schoolCode string
		firstName  string
		lastName   string
		email      string
		password   string
	)

	flag.StringVar(&dsn, "db-dsn", os.Getenv("SKULI_DB_DSN"), "PostgreSQL DSN")
	flag.StringVar(&schoolName, "school-name", "", "Name of the school to create")
	flag.StringVar(&schoolCode, "school-code", "", "Unique code/subdomain for the school")
	flag.StringVar(&firstName, "admin-first-name", "Super", "Super admin first name")
	flag.StringVar(&lastName, "admin-last-name", "Admin", "Super admin last name")
	flag.StringVar(&email, "admin-email", "", "Super admin email")
	flag.StringVar(&password, "admin-password", os.Getenv("SKULI_SEED_PASSWORD"), "Super admin password (min 8 chars)")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if dsn == "" || schoolName == "" || schoolCode == "" || email == "" || password == "" {
		logger.Error("missing required flags: -db-dsn, -school-name, -school-code, -admin-email, -admin-password")
		os.Exit(1)
	}

	db, err := database.Open(config.DB{
		DSN:          dsn,
		MaxOpenConns: 5,
		MaxIdleConns: 5,
		MaxIdleTime:  5 * time.Minute,
	})
	if err != nil {
		logger.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	schoolSvc := schools.NewService(schools.NewRepository(db))
	userSvc := users.NewService(users.NewRepository(db))

	school, err := schoolSvc.Create(ctx, schools.CreateSchoolInput{
		Name: schoolName,
		Code: schoolCode,
	})
	if err != nil {
		logger.Error("failed to create school", "error", err)
		os.Exit(1)
	}
	logger.Info("school created", "id", school.ID, "code", school.Code)

	admin, err := userSvc.Create(ctx, users.CreateUserInput{
		SchoolID:  school.ID,
		Role:      string(users.RoleSuperAdmin),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	})
	if err != nil {
		logger.Error("failed to create super admin", "error", err)
		os.Exit(1)
	}

	logger.Info("super admin created", "id", admin.ID, "email", admin.Email, "school_id", school.ID)
	logger.Info("seed complete")
}
