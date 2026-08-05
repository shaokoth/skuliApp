package teachers

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
	ErrNotFound            = errors.New("teacher not found")
	ErrDuplicateEmployeeNo = errors.New("duplicate employee number")
	ErrEditConflict        = errors.New("edit conflict")
)

const queryTimeout = 3 * time.Second

// Repository is the data-access contract for teachers.
type Repository interface {
	Insert(ctx context.Context, t *Teacher) error
	GetByID(ctx context.Context, schoolID, id int64) (*Teacher, error)
	Update(ctx context.Context, t *Teacher) error
	Delete(ctx context.Context, schoolID, id int64) error
	List(ctx context.Context, f ListFilter) ([]*Teacher, error)
}

// PostgresRepository implements Repository against PostgreSQL.
type PostgresRepository struct {
	db *sql.DB
}

// NewRepository returns a Postgres-backed teacher repository.
func NewRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Insert(ctx context.Context, t *Teacher) error {
	const query = `
		INSERT INTO teachers (
			school_id, user_id, employee_number, first_name, last_name,
			email, phone, gender, qualification, hire_date, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		RETURNING id, created_at, updated_at, version`

	args := []any{
		t.SchoolID, t.UserID, t.EmployeeNumber, t.FirstName, t.LastName,
		t.Email, t.Phone, t.Gender, t.Qualification, t.HireDate, t.Status,
	}

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt, &t.Version)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && strings.Contains(pgErr.ConstraintName, "employee") {
			return ErrDuplicateEmployeeNo
		}
		return err
	}
	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, schoolID, id int64) (*Teacher, error) {
	const query = `
		SELECT id, school_id, user_id, employee_number, first_name, last_name,
		       email, phone, gender, qualification, hire_date, status,
		       created_at, updated_at, version
		FROM teachers
		WHERE school_id = $1 AND id = $2`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	t, err := scanTeacher(r.db.QueryRowContext(ctx, query, schoolID, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return t, nil
}

func (r *PostgresRepository) Update(ctx context.Context, t *Teacher) error {
	const query = `
		UPDATE teachers
		SET first_name = $1, last_name = $2, email = $3, phone = $4,
		    qualification = $5, status = $6, updated_at = now(), version = version + 1
		WHERE school_id = $7 AND id = $8 AND version = $9
		RETURNING version`

	args := []any{
		t.FirstName, t.LastName, t.Email, t.Phone, t.Qualification,
		t.Status, t.SchoolID, t.ID, t.Version,
	}

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&t.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}
		return err
	}
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, schoolID, id int64) error {
	const query = `DELETE FROM teachers WHERE school_id = $1 AND id = $2`

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

// List returns teachers for a school, filtered by status and a search on
// name/employee number. Filter columns (school_id, status) are indexed.
func (r *PostgresRepository) List(ctx context.Context, f ListFilter) ([]*Teacher, error) {
	const query = `
		SELECT id, school_id, user_id, employee_number, first_name, last_name,
		       email, phone, gender, qualification, hire_date, status,
		       created_at, updated_at, version
		FROM teachers
		WHERE school_id = $1
		  AND ($2 = '' OR status = $2)
		  AND ($3 = '' OR first_name ILIKE '%' || $3 || '%'
		               OR last_name  ILIKE '%' || $3 || '%'
		               OR employee_number ILIKE '%' || $3 || '%')
		ORDER BY id
		LIMIT $4 OFFSET $5`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, f.SchoolID, f.Status, f.Search, f.Limit(), f.Offset())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []*Teacher{}
	for rows.Next() {
		t, err := scanTeacher(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanTeacher(row scanner) (*Teacher, error) {
	var t Teacher
	err := row.Scan(
		&t.ID, &t.SchoolID, &t.UserID, &t.EmployeeNumber, &t.FirstName, &t.LastName,
		&t.Email, &t.Phone, &t.Gender, &t.Qualification, &t.HireDate, &t.Status,
		&t.CreatedAt, &t.UpdatedAt, &t.Version,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
