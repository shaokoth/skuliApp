package students

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
	ErrNotFound             = errors.New("student not found")
	ErrDuplicateAdmissionNo = errors.New("duplicate admission number")
	ErrEditConflict         = errors.New("edit conflict")
)

const queryTimeout = 3 * time.Second

// Repository is the data-access contract for students.
type Repository interface {
	Insert(ctx context.Context, s *Student) error
	GetByID(ctx context.Context, schoolID, id int64) (*Student, error)
	Update(ctx context.Context, s *Student) error
	Delete(ctx context.Context, schoolID, id int64) error
	List(ctx context.Context, f ListFilter) ([]*Student, error)
}

// PostgresRepository implements Repository against PostgreSQL.
type PostgresRepository struct {
	db *sql.DB
}

// NewRepository returns a Postgres-backed student repository.
func NewRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Insert(ctx context.Context, s *Student) error {
	const query = `
		INSERT INTO students (
			school_id, user_id, admission_number, first_name, last_name,
			date_of_birth, gender, class_id, parent_id, address, phone, email,
			photo_url, enrollment_date, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
		RETURNING id, created_at, updated_at, version`

	args := []any{
		s.SchoolID, s.UserID, s.AdmissionNumber, s.FirstName, s.LastName,
		s.DateOfBirth, s.Gender, s.ClassID, s.ParentID, s.Address, s.Phone, s.Email,
		s.PhotoURL, s.EnrollmentDate, s.Status,
	}

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.Version)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && strings.Contains(pgErr.ConstraintName, "admission") {
			return ErrDuplicateAdmissionNo
		}
		return err
	}
	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, schoolID, id int64) (*Student, error) {
	const query = `
		SELECT id, school_id, user_id, admission_number, first_name, last_name,
		       date_of_birth, gender, class_id, parent_id, address, phone, email,
		       photo_url, enrollment_date, status, created_at, updated_at, version
		FROM students
		WHERE school_id = $1 AND id = $2`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	s, err := scanStudent(r.db.QueryRowContext(ctx, query, schoolID, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return s, nil
}

func (r *PostgresRepository) Update(ctx context.Context, s *Student) error {
	const query = `
		UPDATE students
		SET first_name = $1, last_name = $2, gender = $3, class_id = $4,
		    parent_id = $5, address = $6, phone = $7, email = $8, status = $9,
		    updated_at = now(), version = version + 1
		WHERE school_id = $10 AND id = $11 AND version = $12
		RETURNING version`

	args := []any{
		s.FirstName, s.LastName, s.Gender, s.ClassID, s.ParentID,
		s.Address, s.Phone, s.Email, s.Status, s.SchoolID, s.ID, s.Version,
	}

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

func (r *PostgresRepository) Delete(ctx context.Context, schoolID, id int64) error {
	const query = `DELETE FROM students WHERE school_id = $1 AND id = $2`

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

// List returns students for a school, filtered by class/status and a search on
// name/admission number. Filter columns (school_id, class_id, status) are indexed.
func (r *PostgresRepository) List(ctx context.Context, f ListFilter) ([]*Student, error) {
	const query = `
		SELECT id, school_id, user_id, admission_number, first_name, last_name,
		       date_of_birth, gender, class_id, parent_id, address, phone, email,
		       photo_url, enrollment_date, status, created_at, updated_at, version
		FROM students
		WHERE school_id = $1
		  AND ($2::bigint IS NULL OR class_id = $2)
		  AND ($3 = '' OR status = $3)
		  AND ($4 = '' OR first_name ILIKE '%' || $4 || '%'
		               OR last_name  ILIKE '%' || $4 || '%'
		               OR admission_number ILIKE '%' || $4 || '%')
		ORDER BY id
		LIMIT $5 OFFSET $6`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, f.SchoolID, f.ClassID, f.Status, f.Search, f.Limit(), f.Offset())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []*Student{}
	for rows.Next() {
		s, err := scanStudent(rows)
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

// scanner abstracts *sql.Row and *sql.Rows for a shared scan helper.
type scanner interface {
	Scan(dest ...any) error
}

func scanStudent(row scanner) (*Student, error) {
	var s Student
	err := row.Scan(
		&s.ID, &s.SchoolID, &s.UserID, &s.AdmissionNumber, &s.FirstName, &s.LastName,
		&s.DateOfBirth, &s.Gender, &s.ClassID, &s.ParentID, &s.Address, &s.Phone, &s.Email,
		&s.PhotoURL, &s.EnrollmentDate, &s.Status, &s.CreatedAt, &s.UpdatedAt, &s.Version,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
