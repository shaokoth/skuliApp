CREATE TABLE IF NOT EXISTS students (
    id               bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    school_id        bigint NOT NULL REFERENCES schools (id) ON DELETE CASCADE,
    user_id          bigint REFERENCES users (id) ON DELETE SET NULL,
    admission_number text NOT NULL,
    first_name       text NOT NULL,
    last_name        text NOT NULL,
    date_of_birth    date NOT NULL,
    gender           text NOT NULL DEFAULT '',
    class_id         bigint REFERENCES classes (id) ON DELETE SET NULL,
    parent_id        bigint REFERENCES parents (id) ON DELETE SET NULL,
    address          text NOT NULL DEFAULT '',
    phone            text NOT NULL DEFAULT '',
    email            text NOT NULL DEFAULT '',
    photo_url        text NOT NULL DEFAULT '',
    enrollment_date  date NOT NULL DEFAULT CURRENT_DATE,
    status           text NOT NULL DEFAULT 'active',
    created_at       timestamptz NOT NULL DEFAULT now(),
    updated_at       timestamptz NOT NULL DEFAULT now(),
    version          integer NOT NULL DEFAULT 1,
    CONSTRAINT students_school_id_admission_number_key UNIQUE (school_id, admission_number),
    CONSTRAINT students_status_check CHECK (status IN (
        'active', 'graduated', 'transferred', 'suspended', 'withdrawn'))
);

-- Tenant scoping plus the class/status list filters are all covered by indexes.
CREATE INDEX idx_students_school_id ON students (school_id);
CREATE INDEX idx_students_school_class ON students (school_id, class_id);
CREATE INDEX idx_students_school_status ON students (school_id, status);
CREATE INDEX idx_students_parent_id ON students (parent_id);

-- Trigram indexes back the ILIKE '%term%' search in the students repository.
CREATE INDEX idx_students_first_name_trgm ON students USING gin (first_name gin_trgm_ops);
CREATE INDEX idx_students_last_name_trgm ON students USING gin (last_name gin_trgm_ops);
CREATE INDEX idx_students_admission_number_trgm ON students USING gin (admission_number gin_trgm_ops);
