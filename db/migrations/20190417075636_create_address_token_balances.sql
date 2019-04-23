-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE VIEW accounts.address_token_balances AS (
    SELECT
        address_hash,
        token_contract_address_hash,
        block_number,
        COALESCE(SUM(amount), 0)
    FROM (
        SELECT
            xfer.src AS address_hash,
            xfer.contract AS token_contract_address_hash,
            xfer.block_number AS block_number,
            -SUM(xfer.amount) OVER (
                PARTITION BY xfer.src, xfer.contract, xfer.block_number
                ORDER BY xfer.block_number ASC
                ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
            ) AS amount
        FROM accounts.token_value_transfers AS xfer

        UNION ALL

        SELECT
            xfer.dst AS address_hash,
            xfer.contract AS token_contract_address_hash,
            xfer.block_number AS block_number,
            sum(xfer.amount) OVER (
                PARTITION BY xfer.dst, xfer.contract, xfer.block_number
                ORDER BY xfer.block_number ASC
                ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
            ) AS amount
        FROM accounts.token_value_transfers AS xfer
    ) AS x
    GROUP BY 1, 2, 3
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW accounts.address_token_balances;
-- +goose StatementEnd
