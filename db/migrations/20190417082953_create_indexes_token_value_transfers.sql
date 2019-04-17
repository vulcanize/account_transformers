-- +goose Up
-- +goose StatementBegin
CREATE INDEX CONCURRENTLY address_src
ON accounts.token_value_transfers(src, block_number);
CREATE INDEX CONCURRENTLY address_dest
ON accounts.token_value_transfers(dest, block_number);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX address_dest;
DROP INDEX address_src;
-- +goose StatementEnd
