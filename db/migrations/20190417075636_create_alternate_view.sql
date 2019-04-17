-- +goose Up
-- +goose StatementBegin
create or replace view accounts.address_token_balances_two as (
    select
        address_hash,
        token_contract_address_hash,
        block_number,
        sum(amount)
    from (
        select
            xfer.src as address_hash,
            xfer.contract as token_contract_address_hash,
            xfer.block_number as block_number,
            -sum(xfer.amount) over (
                partition by xfer.src, xfer.contract, xfer.block_number
                order by xfer.block_number asc
                rows between unbounded preceding and current row
            ) as amount
        from accounts.token_value_transfers as xfer

        union all

        select
            xfer.dst as address_hash,
            xfer.contract as token_contract_address_hash,
            xfer.block_number as block_number,
            sum(xfer.amount) over (
                partition by xfer.dst, xfer.contract, xfer.block_number
                order by xfer.block_number asc
                rows between unbounded preceding and current row
            ) as amount
        from accounts.token_value_transfers as xfer
    ) as x
    group by 1, 2, 3
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW address_token_balances_two;
-- +goose StatementEnd
