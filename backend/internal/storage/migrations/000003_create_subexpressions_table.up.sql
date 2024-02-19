CREATE TABLE IF NOT EXISTS subexpressions (
    id INT PRIMARY KEY,
    expression_id UUID NOT NULL,
    worker_id INT,
    subexpression TEXT NOT NULL,
    is_taken BOOLEAN NOT NULL DEFAULT FALSE,
    taken_at TIMESTAMP,
    is_being_checked BOOLEAN NOT NULL DEFAULT FALSE,
    timeout BIGINT NOT NULL DEFAULT 0,
    depends_on INTEGER[],
    result FLOAT NOT NULL DEFAULT 0,
    is_done BOOLEAN NOT NULL DEFAULT FALSE
);