-- +goose Up
CREATE TABLE accounts.contract_addresses (
  contract BYTEA PRIMARY KEY
);

-- +goose Down
DROP TABLE accounts.contract_addresses;