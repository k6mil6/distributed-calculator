CREATE TABLE subexpressions (
    id BIGSERIAL PRIMARY KEY,
    expression_id BIGINT NOT NULL,
    subexpression TEXT NOT NULL,
    is_taken BOOLEAN NOT NULL DEFAULT FALSE,
    is_done BOOLEAN NOT NULL DEFAULT FALSE,
);