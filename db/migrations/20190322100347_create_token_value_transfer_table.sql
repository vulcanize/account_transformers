-- +goose Up
CREATE TABLE accounts.token_value_transfers (
  id           SERIAL PRIMARY KEY,
  header_id    INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  block_number BIGINT NOT NULL,
  name         VARCHAR NOT NULL CHECK (name <> ''),
  dst          VARCHAR(66),
  src          VARCHAR(66),
  amount       NUMERIC,
  contract     VARCHAR(66) NOT NULL,
  log_idx      INTEGER NOT NULL,
  tx_idx       INTEGER NOT NULL,
  raw_log      JSONB,
  UNIQUE (header_id, tx_idx, log_idx)
);

-- +goose Down
DROP TABLE accounts.token_value_transfers;
