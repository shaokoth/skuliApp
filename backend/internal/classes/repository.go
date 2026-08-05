package classes

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// Repository-level sentinel errors.
var (
	ErrNotFound       = errors.New("class not found")
	ErrDuplicateClass = errors.New("duplicate class for this academic year")
	ErrEditConflict   = errors.New("edit conflict")
)

const queryTimeout = 3 * time.Second

// Repository is the data-access contract for classes.
type Repository interface {
	Insert(ctx context.Context, c *Class) error
	GetByID(ctx context.Context, schoolID, id int64) (*Class, error)
	Update(ctx context.Context, c *Class) error
	Delete(ctx context.Context, schoolID, id int64) error
	List(ctx context.Context, f ListFilter) ([]*Class, error)
}

// PostgresRepository implements Repository against PostgreSQL.
type PostgresRepository struct {
	db *sql.DB
}

// NewRepository returns a Postgres-backed class repository.
func NewRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Insert(ctx context.Context, c *Class) error {
	const query = `
		INSERT INTO classes (school_id, name, grade_level, section, class_teacher_id, capacity, academic_year)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at, version`

	args := []any{c.SchoolID, c.Name, c.GradeLevel, c.Section, c.ClassTeacherID, c.Capacity, c.AcademicYear}

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt, &c.Version)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrDuplicateClass
		}
		return err
	}
	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, schoolID, id int64) (*Class, error) {
	const query = `
		SELECT id, school_id, name, grade_level, section, class_teacher_id,
		       capacity, academic_year, created_at, updated_at, version
		FROM classes
		WHERE school_id = $1 AND id = $2`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	c, err := scanClass(r.db.QueryRowContext(ctx, query, schoolID, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return c, nil
}

func (r *PostgresRepository) Update(ctx context.Context, c *Class) error {
	const query = `
		UPDATE classes
		SET name = $1, grade_level = $2, section = $3, class_teacher_id = $4,
		    capacity = $5, academic_year = $6, updated_at = now(), version = version + 1
		WHERE school_id = $7 AND id = $8 AND version = $9
		RETURNING version`

	args := []any{c.Name, c.GradeLevel, c.Section, c.ClassTeacherID, c.Capacity, c.AcademicYear, c.SchoolID, c.ID, c.Version}

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&c.Version)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrDuplicateClass
		}
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}
		return err
	}
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, schoolID, id int64) error {
	const query = `DELETE FROM classes WHERE school_id = $1 AND id = $2`

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

// List returns classes for a school, filtered by academic year and name search.
func (r *PostgresRepository) List(ctx context.Context, f ListFilter) ([]*Class, error) {
	const query = `
		SELECT id, school_id, name, grade_level, section, class_teacher_id,
		       capacity, academic_year, created_at, updated_at, version
		FROM classes
		WHERE school_id = $1
		  AND ($2 = '' OR academic_year = $2)
		  AND ($3 = '' OR name ILIKE '%' || $3 || '%')
		ORDER BY id
		LIMIT $4 OFFSET $5`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, f.SchoolID, f.AcademicYear, f.Search, f.Limit(), f.Offset())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []*Class{}
	for rows.Next() {
		c, err := scanClass(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanClass(row scanner) (*Class, error) {
	var c Class
	err := row.Scan(
		&c.ID, &c.SchoolID, &c.Name, &c.GradeLevel, &c.Section, &c.ClassTeacherID,
		&c.Capacity, &c.AcademicYear, &c.CreatedAt, &c.UpdatedAt, &c.Version,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
