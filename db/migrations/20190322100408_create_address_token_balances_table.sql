-- +goose Up
CREATE TABLE accounts.address_token_balances (
  id                          SERIAL PRIMARY KEY,
  address_hash                BYTEA NOT NULL,
  block_number                BIGINT NOT NULL,
  token_contract_address_hash BYTEA NOT NULL,
  value                       NUMERIC,
  value_fetched_at            TIMESTAMP WITHOUT TIME ZONE,
  inserted_at                 TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  updated_at                  TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX address_token_balances_address_hash_token_contract_address_hash ON accounts.address_token_balances (address_hash, token_contract_address_hash, block_number);
CREATE UNIQUE INDEX unfetched_token_balances ON accounts.address_token_balances (address_hash, token_contract_address_hash, block_number) WHERE (value_fetched_at IS NULL);

-- +goose Down
DROP TABLE accounts.address_token_balances;
