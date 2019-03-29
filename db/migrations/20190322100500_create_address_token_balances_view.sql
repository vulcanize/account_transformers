-- +goose Up
CREATE OR REPLACE VIEW accounts.address_token_balances AS
  SELECT
    accounts.addresses.address AS address_hash,
    accounts.watched_contracts.contract AS token_contract_address_hash,
    public.headers.block_number,
    ((SELECT COALESCE(SUM(amount),0) FROM accounts.token_value_transfers
                        WHERE accounts.token_value_transfers.block_number <= public.headers.block_number
                        AND accounts.token_value_transfers.dst = accounts.addresses.address
                        AND accounts.token_value_transfers.contract = accounts.watched_contracts.contract) -
    (SELECT COALESCE(SUM(amount),0) FROM accounts.token_value_transfers
                        WHERE accounts.token_value_transfers.block_number <= public.headers.block_number
                        AND accounts.token_value_transfers.src = accounts.addresses.address
                        AND accounts.token_value_transfers.contract = accounts.watched_contracts.contract)) AS "value"
  FROM accounts.token_value_transfers, accounts.addresses, public.headers, accounts.watched_contracts
  GROUP BY accounts.addresses.address, accounts.watched_contracts.contract, public.headers.block_number;

-- +goose Down
DROP VIEW accounts.address_token_balances;