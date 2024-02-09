CREATE TABLE subexpressions (
    id BIGSERIAL PRIMARY KEY,
    expression_id BIGINT NOT NULL,
    worker_id BIGINT
    subexpression TEXT NOT NULL,
    is_taken BOOLEAN NOT NULL DEFAULT FALSE,
    is_done BOOLEAN NOT NULL DEFAULT FALSE,
    timeout INTERVAL NOT NULL DEFAULT '00:01:00',
    result FLOAT
);