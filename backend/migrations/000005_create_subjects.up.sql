CREATE TABLE IF NOT EXISTS subjects (
    id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    school_id  bigint NOT NULL REFERENCES schools (id) ON DELETE CASCADE,
    name       text NOT NULL,
    code       text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    version    integer NOT NULL DEFAULT 1
);

CREATE INDEX idx_subjects_school_id ON subjects (school_id);

-- Codes are unique per school when present (blank codes are allowed to repeat).
CREATE UNIQUE INDEX subjects_school_id_code_key ON subjects (school_id, code) WHERE code <> '';
