// Command migrate applies or rolls back database schema migrations using the
// embedded SQL files (golang-migrate iofs source + pgx/v5 database driver).
package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/skuliapp/backend/migrations"
)

func main() {
	var (
		dsn    string
		action string
	)
	flag.StringVar(&dsn, "db-dsn", os.Getenv("SKULI_DB_DSN"), "PostgreSQL DSN")
	flag.StringVar(&action, "action", "up", "migration action: up | down | version")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if dsn == "" {
		logger.Error("missing database DSN (set -db-dsn flag or SKULI_DB_DSN env var)")
		os.Exit(1)
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Error("opening database failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	switch action {
	case "up":
		if err := migrations.Up(db); err != nil {
			logger.Error("migration failed", "action", action, "error", err)
			os.Exit(1)
		}
	case "down":
		if err := migrations.Down(db); err != nil {
			logger.Error("migration failed", "action", action, "error", err)
			os.Exit(1)
		}
	case "version":
		version, dirty, err := migrations.Version(db)
		if err != nil {
			logger.Error("reading version failed", "error", err)
			os.Exit(1)
		}
		logger.Info("migration status", "version", version, "dirty", dirty)
		return
	default:
		logger.Error("unknown action (use up | down | version)", "action", action)
		os.Exit(1)
	}

	logger.Info("migrations applied", "action", action)
}
