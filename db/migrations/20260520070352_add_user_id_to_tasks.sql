-- +goose Up
ALTER TABLE tasks ADD COLUMN user_id INTEGER;

ALTER TABLE tasks
ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;

ALTER TABLE tasks ALTER COLUMN user_id SET NOT NULL;

-- +goose Down
ALTER TABLE tasks DROP CONSTRAINT fk_user_id;

ALTER TABLE tasks DROP COLUMN user_id;