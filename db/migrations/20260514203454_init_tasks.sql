-- +goose Up
CREATE TABLE tasks(
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE tasks;
