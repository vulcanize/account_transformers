-- +goose Up
CREATE TABLE accounts.addresses (
  address VARCHAR(66) PRIMARY KEY
);

-- +goose Down
DROP TABLE accounts.addresses;