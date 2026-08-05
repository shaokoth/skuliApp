CREATE TABLE IF NOT EXISTS users (
    id            bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    school_id     bigint NOT NULL REFERENCES schools (id) ON DELETE CASCADE,
    role          text NOT NULL,
    first_name    text NOT NULL,
    last_name     text NOT NULL,
    email         text NOT NULL,
    phone         text NOT NULL DEFAULT '',
    password_hash bytea NOT NULL,
    active        boolean NOT NULL DEFAULT true,
    last_login_at timestamptz,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    version       integer NOT NULL DEFAULT 1,
    CONSTRAINT users_email_key UNIQUE (email),
    CONSTRAINT users_role_check CHECK (role IN (
        'super_admin', 'school_admin', 'principal',
        'teacher', 'accountant', 'parent', 'student'))
);

-- Tenant scoping and role filters hit these composite btree indexes.
CREATE INDEX idx_users_school_id ON users (school_id);
CREATE INDEX idx_users_school_role ON users (school_id, role);

-- Trigram indexes back the ILIKE '%term%' search in the users repository.
CREATE INDEX idx_users_first_name_trgm ON users USING gin (first_name gin_trgm_ops);
CREATE INDEX idx_users_last_name_trgm ON users USING gin (last_name gin_trgm_ops);
CREATE INDEX idx_users_email_trgm ON users USING gin (email gin_trgm_ops);
