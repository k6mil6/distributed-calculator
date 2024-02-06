CREATE TABLE expressions (
    id BIGSERIAL PRIMARY KEY,
    expression TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    is_taken BOOLEAN NOT NULL DEFAULT TRUE
);