-- +goose Up
CREATE TABLE accounts.watched_contracts (
  contract BYTEA PRIMARY KEY
);

-- +goose Down
DROP TABLE accounts.addresses;