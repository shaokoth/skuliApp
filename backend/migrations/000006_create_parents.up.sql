CREATE TABLE IF NOT EXISTS parents (
    id           bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    school_id    bigint NOT NULL REFERENCES schools (id) ON DELETE CASCADE,
    user_id      bigint REFERENCES users (id) ON DELETE SET NULL,
    first_name   text NOT NULL,
    last_name    text NOT NULL,
    email        text NOT NULL DEFAULT '',
    phone        text NOT NULL DEFAULT '',
    occupation   text NOT NULL DEFAULT '',
    address      text NOT NULL DEFAULT '',
    relationship text NOT NULL DEFAULT '',
    created_at   timestamptz NOT NULL DEFAULT now(),
    updated_at   timestamptz NOT NULL DEFAULT now(),
    version      integer NOT NULL DEFAULT 1
);

CREATE INDEX idx_parents_school_id ON parents (school_id);
CREATE INDEX idx_parents_phone ON parents (phone);
