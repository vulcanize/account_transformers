-- +goose Up
CREATE TABLE accounts.eth_balances (
  id          SERIAL PRIMARY KEY,
  address     VARCHAR(66),
  eth_balance NUMERIC NOT NULL,
  block_id    INTEGER,
  CONSTRAINT account_uc UNIQUE (address, block_id),
  CONSTRAINT block_id_fk FOREIGN KEY (block_id)
  REFERENCES public.blocks (id),
  CONSTRAINT address_fk FOREIGN KEY (address)
  REFERENCES accounts.addresses (address)
);

-- +goose Down
DROP TABLE accounts.eth_balances;