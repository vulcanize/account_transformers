-- +goose Up
CREATE SCHEMA accounts;

-- +goose Down
DROP SCHEMA accounts CASCADE;
