CREATE TABLE IF NOT EXISTS expressions (
    id UUID PRIMARY KEY,
    expression TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    timeouts_id INT,
    is_taken BOOLEAN NOT NULL DEFAULT FALSE,
    result FLOAT NOT NULL DEFAULT 0,
    is_done BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY(timeouts_id) REFERENCES timeouts(id)
);