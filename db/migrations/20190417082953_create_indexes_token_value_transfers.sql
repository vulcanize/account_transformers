-- +goose Up
-- +goose StatementBegin
CREATE INDEX address_src
ON accounts.token_value_transfers(src, block_number);
CREATE INDEX address_dst
ON accounts.token_value_transfers(dst, block_number);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX address_dst;
DROP INDEX address_src;
-- +goose StatementEnd
