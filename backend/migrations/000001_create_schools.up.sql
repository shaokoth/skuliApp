CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE IF NOT EXISTS schools (
    id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name       text NOT NULL,
    code       text NOT NULL,
    email      text NOT NULL DEFAULT '',
    phone      text NOT NULL DEFAULT '',
    address    text NOT NULL DEFAULT '',
    logo_url   text NOT NULL DEFAULT '',
    active     boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    version    integer NOT NULL DEFAULT 1,
    CONSTRAINT schools_code_key UNIQUE (code)
);
