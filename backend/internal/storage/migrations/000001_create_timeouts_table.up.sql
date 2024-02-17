CREATE TABLE IF NOT EXISTS timeouts (
    id SERIAL PRIMARY KEY,
    timeouts_values JSONB NOT NULL
);