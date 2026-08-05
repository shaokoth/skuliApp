CREATE TABLE IF NOT EXISTS classes (
    id               bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    school_id        bigint NOT NULL REFERENCES schools (id) ON DELETE CASCADE,
    name             text NOT NULL,
    grade_level      text NOT NULL DEFAULT '',
    section          text NOT NULL DEFAULT '',
    class_teacher_id bigint REFERENCES teachers (id) ON DELETE SET NULL,
    capacity         integer NOT NULL DEFAULT 0,
    academic_year    text NOT NULL DEFAULT '',
    created_at       timestamptz NOT NULL DEFAULT now(),
    updated_at       timestamptz NOT NULL DEFAULT now(),
    version          integer NOT NULL DEFAULT 1,
    CONSTRAINT classes_school_id_name_year_key UNIQUE (school_id, name, academic_year)
);

CREATE INDEX idx_classes_school_id ON classes (school_id);
CREATE INDEX idx_classes_class_teacher_id ON classes (class_teacher_id);
