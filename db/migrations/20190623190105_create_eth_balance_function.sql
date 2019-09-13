-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION eth_balance(account VARCHAR(66), block BIGINT) RETURNS NUMERIC AS $$
BEGIN
  SELECT accounts.state_accounts.balance
    FROM accounts.state_accounts, public.headers
    WHERE public.headers.block_number <= block
    AND accounts.state_accounts.header_id = public.headers.id
    AND accounts.state_accounts.account_key = account
    ORDER BY public.headers.block_number DESC LIMIT 1;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd
-- +goose Down
DROP FUNCTION accounts.eth_balance;