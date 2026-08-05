package migrations

import (
	"database/sql"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	pgxmigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// newMigrator builds a *migrate.Migrate over the embedded SQL files, running
// against the supplied database handle (opened with the pgx stdlib driver).
func newMigrator(db *sql.DB) (*migrate.Migrate, error) {
	source, err := iofs.New(FS, ".")
	if err != nil {
		return nil, err
	}
	driver, err := pgxmigrate.WithInstance(db, &pgxmigrate.Config{})
	if err != nil {
		return nil, err
	}
	return migrate.NewWithInstance("iofs", source, "pgx5", driver)
}

// Up applies all pending migrations. It is a no-op (nil error) when the schema
// is already current, and is safe to call concurrently (advisory-locked).
func Up(db *sql.DB) error {
	m, err := newMigrator(db)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

// Down rolls back all migrations.
func Down(db *sql.DB) error {
	m, err := newMigrator(db)
	if err != nil {
		return err
	}
	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

// Version reports the current migration version and whether it is dirty.
func Version(db *sql.DB) (version uint, dirty bool, err error) {
	m, err := newMigrator(db)
	if err != nil {
		return 0, false, err
	}
	return m.Version()
}
