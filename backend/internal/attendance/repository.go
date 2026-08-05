package attendance

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// Repository-level sentinel errors.
var (
	ErrNotFound      = errors.New("attendance record not found")
	ErrDuplicateMark = errors.New("attendance already marked for this student on this date")
)

const queryTimeout = 3 * time.Second

// Repository is the data-access contract for attendance records.
type Repository interface {
	Insert(ctx context.Context, a *Attendance) error
	GetByID(ctx context.Context, schoolID, id int64) (*Attendance, error)
	Update(ctx context.Context, a *Attendance) error
	Delete(ctx context.Context, schoolID, id int64) error
	List(ctx context.Context, f ListFilter) ([]*Attendance, error)
}

// PostgresRepository implements Repository against PostgreSQL.
type PostgresRepository struct {
	db *sql.DB
}

// NewRepository returns a Postgres-backed attendance repository.
func NewRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Insert(ctx context.Context, a *Attendance) error {
	const query = `
		INSERT INTO attendance (school_id, student_id, class_id, date, status, remark, marked_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	args := []any{a.SchoolID, a.StudentID, a.ClassID, a.Date, a.Status, a.Remark, a.MarkedBy}

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrDuplicateMark
		}
		return err
	}
	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, schoolID, id int64) (*Attendance, error) {
	const query = `
		SELECT id, school_id, student_id, class_id, date, status, remark, marked_by, created_at, updated_at
		FROM attendance
		WHERE school_id = $1 AND id = $2`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	a, err := scanAttendance(r.db.QueryRowContext(ctx, query, schoolID, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return a, nil
}

// Update corrects the status/remark of an existing mark.
func (r *PostgresRepository) Update(ctx context.Context, a *Attendance) error {
	const query = `
		UPDATE attendance
		SET status = $1, remark = $2, updated_at = now()
		WHERE school_id = $3 AND id = $4
		RETURNING updated_at`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, a.Status, a.Remark, a.SchoolID, a.ID).Scan(&a.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, schoolID, id int64) error {
	const query = `DELETE FROM attendance WHERE school_id = $1 AND id = $2`

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

// List returns attendance records for a school, filtered by class, student,
// status and date. All filter columns are covered by indexes.
func (r *PostgresRepository) List(ctx context.Context, f ListFilter) ([]*Attendance, error) {
	const query = `
		SELECT id, school_id, student_id, class_id, date, status, remark, marked_by, created_at, updated_at
		FROM attendance
		WHERE school_id = $1
		  AND ($2::bigint IS NULL OR class_id = $2)
		  AND ($3::bigint IS NULL OR student_id = $3)
		  AND ($4 = '' OR status = $4)
		  AND ($5::date IS NULL OR date = $5::date)
		ORDER BY date DESC, id
		LIMIT $6 OFFSET $7`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, f.SchoolID, f.ClassID, f.StudentID, f.Status, f.Date, f.Limit(), f.Offset())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []*Attendance{}
	for rows.Next() {
		a, err := scanAttendance(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanAttendance(row scanner) (*Attendance, error) {
	var a Attendance
	err := row.Scan(
		&a.ID, &a.SchoolID, &a.StudentID, &a.ClassID, &a.Date, &a.Status,
		&a.Remark, &a.MarkedBy, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
