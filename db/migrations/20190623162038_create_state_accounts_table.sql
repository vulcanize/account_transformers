-- +goose Up
CREATE TABLE accounts.state_accounts (
  id           SERIAL PRIMARY KEY,
  header_id    INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  account_key  VARCHAR(66) NOT NULL,
  balance      NUMERIC NOT NULL,
  root         VARCHAR(66) NOT NULL,
  nonce        INTEGER NOT NULL,
  code_hash    BYTEA,
  UNIQUE (header_id, account_key)
);

-- +goose Down
DROP TABLE accounts.state_accounts;
