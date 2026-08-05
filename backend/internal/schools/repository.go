package schools

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
	ErrNotFound      = errors.New("school not found")
	ErrDuplicateCode = errors.New("duplicate school code")
	ErrEditConflict  = errors.New("edit conflict")
)

const queryTimeout = 3 * time.Second

// Repository is the data-access contract for schools (the tenant root).
type Repository interface {
	Insert(ctx context.Context, s *School) error
	GetByID(ctx context.Context, id int64) (*School, error)
	Update(ctx context.Context, s *School) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, f ListFilter) ([]*School, error)
}

// PostgresRepository implements Repository against PostgreSQL.
type PostgresRepository struct {
	db *sql.DB
}

// NewRepository returns a Postgres-backed school repository.
func NewRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Insert(ctx context.Context, s *School) error {
	const query = `
		INSERT INTO schools (name, code, email, phone, address, logo_url, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at, version`

	args := []any{s.Name, s.Code, s.Email, s.Phone, s.Address, s.LogoURL, s.Active}

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.Version)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && strings.Contains(pgErr.ConstraintName, "code") {
			return ErrDuplicateCode
		}
		return err
	}
	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int64) (*School, error) {
	const query = `
		SELECT id, name, code, email, phone, address, logo_url, active,
		       created_at, updated_at, version
		FROM schools
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	s, err := scanSchool(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return s, nil
}

func (r *PostgresRepository) Update(ctx context.Context, s *School) error {
	const query = `
		UPDATE schools
		SET name = $1, email = $2, phone = $3, address = $4, logo_url = $5,
		    active = $6, updated_at = now(), version = version + 1
		WHERE id = $7 AND version = $8
		RETURNING version`

	args := []any{s.Name, s.Email, s.Phone, s.Address, s.LogoURL, s.Active, s.ID, s.Version}

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&s.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}
		return err
	}
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM schools WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, id)
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

func (r *PostgresRepository) List(ctx context.Context, f ListFilter) ([]*School, error) {
	const query = `
		SELECT id, name, code, email, phone, address, logo_url, active,
		       created_at, updated_at, version
		FROM schools
		WHERE ($1 = '' OR name ILIKE '%' || $1 || '%' OR code ILIKE '%' || $1 || '%')
		ORDER BY id
		LIMIT $2 OFFSET $3`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, f.Search, f.Limit(), f.Offset())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []*School{}
	for rows.Next() {
		s, err := scanSchool(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanSchool(row scanner) (*School, error) {
	var s School
	err := row.Scan(
		&s.ID, &s.Name, &s.Code, &s.Email, &s.Phone, &s.Address, &s.LogoURL,
		&s.Active, &s.CreatedAt, &s.UpdatedAt, &s.Version,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
