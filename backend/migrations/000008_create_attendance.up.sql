CREATE TABLE IF NOT EXISTS attendance (
    id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    school_id  bigint NOT NULL REFERENCES schools (id) ON DELETE CASCADE,
    student_id bigint NOT NULL REFERENCES students (id) ON DELETE CASCADE,
    class_id   bigint NOT NULL REFERENCES classes (id) ON DELETE CASCADE,
    date       date NOT NULL,
    status     text NOT NULL DEFAULT 'present',
    remark     text NOT NULL DEFAULT '',
    marked_by  bigint REFERENCES users (id) ON DELETE SET NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT attendance_status_check CHECK (status IN ('present', 'absent', 'late', 'excused')),
    CONSTRAINT attendance_school_student_date_key UNIQUE (school_id, student_id, date)
);

-- One student is marked once per day; lookups are by class+date or by student.
CREATE INDEX idx_attendance_school_id ON attendance (school_id);
CREATE INDEX idx_attendance_school_class_date ON attendance (school_id, class_id, date);
CREATE INDEX idx_attendance_student ON attendance (student_id);
CREATE INDEX idx_attendance_date ON attendance (date);
