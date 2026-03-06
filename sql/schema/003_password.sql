-- +goose Up
Alter TABLE users
ADD COLUMN hashed_password TEXT NOT NULL DEFAULT 'unset';

-- +goose Down
Alter TABLE users
DROP COLUMN hashed_password;