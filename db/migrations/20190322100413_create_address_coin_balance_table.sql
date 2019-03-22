-- +goose Up
CREATE TABLE accounts.address_coin_balances (
  id                          BIGINT PRIMARY KEY,
  address_hash                BYTEA NOT NULL,
  block_number                BIGINT NOT NULL,
  value                       NUMERIC(100,0),
  value_fetched_at            TIMESTAMP WITHOUT TIME ZONE,
  inserted_at                 TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  updated_at                  TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX address_coin_balances_address_hash_block_number_index ON accounts.address_coin_balances (address_hash, block_number);
CREATE UNIQUE INDEX unfetched_balances ON accounts.address_coin_balances (address_hash, block_number) WHERE (value_fetched_at IS NULL);

-- +goose Down
DROP TABLE accounts.address_coin_balances;
