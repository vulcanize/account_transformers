-- +goose Up
CREATE TABLE accounts.contract_addresses (
  contract VARCHAR(42) PRIMARY KEY
);

-- +goose Down
DROP TABLE accounts.contract_addresses;