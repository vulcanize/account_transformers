-- +goose Up
CREATE TABLE accounts.addresses (
  address BYTEA PRIMARY KEY
);

-- +goose Down
DROP TABLE accounts.addresses;
