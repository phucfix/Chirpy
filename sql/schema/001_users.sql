-- +goose Up
CREATE TABLE users (
id UUID PRIMARY KEY,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
email TEXT UNIQUE DEFAULT NULL
);

-- +goose Down
DROP TABLE users;
