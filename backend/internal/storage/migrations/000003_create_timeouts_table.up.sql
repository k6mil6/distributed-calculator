CREATE TABLE timeouts(
    id SERIAL PRIMARY KEY,
    timeouts JSONB,
    expression_id UUID,
    CONSTRAINT fk_expression FOREIGN KEY(expression_id) REFERENCES expressions(id)
)