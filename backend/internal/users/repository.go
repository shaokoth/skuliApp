package users

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// Repository-level sentinel errors.
var (
	ErrNotFound       = errors.New("user not found")
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrEditConflict   = errors.New("edit conflict")
)

const queryTimeout = 3 * time.Second

// Repository is the data-access contract for users.
type Repository interface {
	Insert(ctx context.Context, u *User) error
	GetByID(ctx context.Context, schoolID, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, schoolID, id int64) error
	List(ctx context.Context, f ListFilter) ([]*User, error)
}

// PostgresRepository implements Repository against PostgreSQL.
type PostgresRepository struct {
	db *sql.DB
}

// NewRepository returns a Postgres-backed user repository.
func NewRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Insert(ctx context.Context, u *User) error {
	const query = `
		INSERT INTO users (school_id, role, first_name, last_name, email, phone, password_hash, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at, version`

	args := []any{u.SchoolID, u.Role, u.FirstName, u.LastName, u.Email, u.Phone, u.PasswordHash, u.Active}

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt, &u.Version)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && strings.Contains(pgErr.ConstraintName, "email") {
			return ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, schoolID, id int64) (*User, error) {
	const query = `
		SELECT id, school_id, role, first_name, last_name, email, phone,
		       password_hash, active, last_login_at, created_at, updated_at, version
		FROM users
		WHERE school_id = $1 AND id = $2`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	var u User
	err := r.db.QueryRowContext(ctx, query, schoolID, id).Scan(
		&u.ID, &u.SchoolID, &u.Role, &u.FirstName, &u.LastName, &u.Email, &u.Phone,
		&u.PasswordHash, &u.Active, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt, &u.Version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *PostgresRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	const query = `
		SELECT id, school_id, role, first_name, last_name, email, phone,
		       password_hash, active, last_login_at, created_at, updated_at, version
		FROM users
		WHERE email = $1`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	var u User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.SchoolID, &u.Role, &u.FirstName, &u.LastName, &u.Email, &u.Phone,
		&u.PasswordHash, &u.Active, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt, &u.Version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

// Update applies changes using optimistic concurrency on the version column.
func (r *PostgresRepository) Update(ctx context.Context, u *User) error {
	const query = `
		UPDATE users
		SET first_name = $1, last_name = $2, phone = $3, active = $4,
		    updated_at = now(), version = version + 1
		WHERE school_id = $5 AND id = $6 AND version = $7
		RETURNING version`

	args := []any{u.FirstName, u.LastName, u.Phone, u.Active, u.SchoolID, u.ID, u.Version}

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&u.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}
		return err
	}
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, schoolID, id int64) error {
	const query = `DELETE FROM users WHERE school_id = $1 AND id = $2`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, schoolID, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// List returns users for a school, filtered by role and a trigram search on
// name/email. The WHERE clause targets indexed columns (school_id, role) and a
// GIN trigram index backs the ILIKE search.
func (r *PostgresRepository) List(ctx context.Context, f ListFilter) ([]*User, error) {
	const query = `
		SELECT id, school_id, role, first_name, last_name, email, phone,
		       password_hash, active, last_login_at, created_at, updated_at, version
		FROM users
		WHERE school_id = $1
		  AND ($2 = '' OR role = $2)
		  AND ($3 = '' OR first_name ILIKE '%' || $3 || '%'
		               OR last_name  ILIKE '%' || $3 || '%'
		               OR email      ILIKE '%' || $3 || '%')
		ORDER BY id
		LIMIT $4 OFFSET $5`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, f.SchoolID, f.Role, f.Search, f.Limit(), f.Offset())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(
			&u.ID, &u.SchoolID, &u.Role, &u.FirstName, &u.LastName, &u.Email, &u.Phone,
			&u.PasswordHash, &u.Active, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt, &u.Version,
		); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
