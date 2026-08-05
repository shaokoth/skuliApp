CREATE TABLE IF NOT EXISTS teachers (
    id              bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    school_id       bigint NOT NULL REFERENCES schools (id) ON DELETE CASCADE,
    user_id         bigint REFERENCES users (id) ON DELETE SET NULL,
    employee_number text NOT NULL,
    first_name      text NOT NULL,
    last_name       text NOT NULL,
    email           text NOT NULL DEFAULT '',
    phone           text NOT NULL DEFAULT '',
    gender          text NOT NULL DEFAULT '',
    qualification   text NOT NULL DEFAULT '',
    hire_date       date NOT NULL,
    status          text NOT NULL DEFAULT 'active',
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now(),
    version         integer NOT NULL DEFAULT 1,
    CONSTRAINT teachers_school_id_employee_number_key UNIQUE (school_id, employee_number),
    CONSTRAINT teachers_status_check CHECK (status IN ('active', 'on_leave', 'suspended', 'resigned'))
);

CREATE INDEX idx_teachers_school_id ON teachers (school_id);
CREATE INDEX idx_teachers_school_status ON teachers (school_id, status);

CREATE INDEX idx_teachers_first_name_trgm ON teachers USING gin (first_name gin_trgm_ops);
CREATE INDEX idx_teachers_last_name_trgm ON teachers USING gin (last_name gin_trgm_ops);
CREATE INDEX idx_teachers_employee_number_trgm ON teachers USING gin (employee_number gin_trgm_ops);
