-- +goose Up
CREATE TABLE accounts.addresses (
  address VARCHAR(42) PRIMARY KEY
);

-- +goose Down
DROP TABLE accounts.addresses;
