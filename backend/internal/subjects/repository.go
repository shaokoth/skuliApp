package subjects

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
	ErrNotFound      = errors.New("subject not found")
	ErrDuplicateCode = errors.New("duplicate subject code")
	ErrEditConflict  = errors.New("edit conflict")
)

const queryTimeout = 3 * time.Second

// Repository is the data-access contract for subjects.
type Repository interface {
	Insert(ctx context.Context, s *Subject) error
	GetByID(ctx context.Context, schoolID, id int64) (*Subject, error)
	Update(ctx context.Context, s *Subject) error
	Delete(ctx context.Context, schoolID, id int64) error
	List(ctx context.Context, f ListFilter) ([]*Subject, error)
}

// PostgresRepository implements Repository against PostgreSQL.
type PostgresRepository struct {
	db *sql.DB
}

// NewRepository returns a Postgres-backed subject repository.
func NewRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Insert(ctx context.Context, s *Subject) error {
	const query = `
		INSERT INTO subjects (school_id, name, code)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at, version`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, s.SchoolID, s.Name, s.Code).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.Version)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && strings.Contains(pgErr.ConstraintName, "code") {
			return ErrDuplicateCode
		}
		return err
	}
	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, schoolID, id int64) (*Subject, error) {
	const query = `
		SELECT id, school_id, name, code, created_at, updated_at, version
		FROM subjects
		WHERE school_id = $1 AND id = $2`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	s, err := scanSubject(r.db.QueryRowContext(ctx, query, schoolID, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return s, nil
}

func (r *PostgresRepository) Update(ctx context.Context, s *Subject) error {
	const query = `
		UPDATE subjects
		SET name = $1, code = $2, updated_at = now(), version = version + 1
		WHERE school_id = $3 AND id = $4 AND version = $5
		RETURNING version`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, s.Name, s.Code, s.SchoolID, s.ID, s.Version).Scan(&s.Version)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && strings.Contains(pgErr.ConstraintName, "code") {
			return ErrDuplicateCode
		}
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}
		return err
	}
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, schoolID, id int64) error {
	const query = `DELETE FROM subjects WHERE school_id = $1 AND id = $2`

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

// List returns subjects for a school, filtered by a name/code search.
func (r *PostgresRepository) List(ctx context.Context, f ListFilter) ([]*Subject, error) {
	const query = `
		SELECT id, school_id, name, code, created_at, updated_at, version
		FROM subjects
		WHERE school_id = $1
		  AND ($2 = '' OR name ILIKE '%' || $2 || '%' OR code ILIKE '%' || $2 || '%')
		ORDER BY name
		LIMIT $3 OFFSET $4`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, f.SchoolID, f.Search, f.Limit(), f.Offset())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []*Subject{}
	for rows.Next() {
		s, err := scanSubject(rows)
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

func scanSubject(row scanner) (*Subject, error) {
	var s Subject
	err := row.Scan(&s.ID, &s.SchoolID, &s.Name, &s.Code, &s.CreatedAt, &s.UpdatedAt, &s.Version)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
