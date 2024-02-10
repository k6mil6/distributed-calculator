CREATE TABLE IF NOT EXISTS subexpressions (
    id SERIAL PRIMARY KEY,
    expression_id UUID NOT NULL,
    worker_id INT,
    subexpression TEXT NOT NULL,
    is_taken BOOLEAN NOT NULL DEFAULT FALSE,
    is_done BOOLEAN NOT NULL DEFAULT FALSE,
    timeout INTERVAL NOT NULL DEFAULT '00:01:00',
    result FLOAT
);