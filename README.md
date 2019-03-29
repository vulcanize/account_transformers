[VulcanizeDB](https://github.com/vulcanize/maker-vulcanizedb) transformers for watching ETH and token value transfers

This repo contains transformers for indexing of ETH and token balances for account addresses.

[This](https://github.com/vulcanize/account_transformers/tree/master/transformers/account/light) 
transformer works by filtering through all eth logs that have one of [these topic0](https://github.com/vulcanize/account_transformers/blob/master/transformers/account/shared/constants/signatures.go#L56).
These events are unpacked and converted to generic "Value Transfer" [records](https://github.com/vulcanize/account_transformers/blob/master/transformers/account/light/models/model.go#L19).
Token balances for user accounts are then constructed  as views on these records.
It then polls an archival Eth node to retrieve balances for these accounts and generate eth balance records.

# Setup 

1. Setup VulcanizeDB
1. Switch to `account_transformer_staging` branch
1. Run `lightSync`
1. Setup config for `account_transformer`
1. Setup Postgraphile to expose the `accounts` schema
1. Run `composeAndExecute` using our `account_transformer` config

These transformers are run as plugin to VulcanizeDB's `composeAndExecute` command.

To begin, setup VulcanizeDB as described [here](https://github.com/vulcanize/maker-vulcanizedb#project-setup).
Currently, this transformer needs to be run from the `account_transformer_staging` [branch](https://github.com/vulcanize/maker-vulcanizedb/tree/account_transformer_staging)
of VulcanizeDB, so switch to that branch before building.

Once vulcanizeDB is setup and built, run vulcanizeDB in `lightSync` [mode](https://github.com/vulcanize/maker-vulcanizedb#alternatively-sync-in-light-mode)
to begin syncing headers into Postgres. It is vital that this sync process begins at a block before the `account.start` field below.

Once `lightSync` has begun, we can run the `composeAndExecute` command to compose and execute our account transformer. To
do so, we use a normal `compose` [config](https://github.com/vulcanize/maker-vulcanizedb#contractwatcher) with two additional parameter maps:

```toml
[contract]
    addresses = [
        "0x0000000000085d4780B73119b644AE5ecd22b376",
    ]
    [contract.equivalents]
    0x0000000000085d4780B73119b644AE5ecd22b376 = [
       "0x8dd5fbCe2F6a956C3022bA3663759011Dd51e73E"
    ]

[account]
    start     = 0
    addresses = [
        "0x48E78948C80e9f8F53190DbDF2990f9a69491ef4",
        "0x009C1E8674038605C5AE33C74f13bC528E1222B5"
    ]
```
`contract.addresses` are a list of the token addresses we want to track balances for, these addresses are used create token balance views.
This can be updated at runtime by adding new contract addresses to the `accounts.watched_contracts` table in Postgres:

```postgresql
CREATE TABLE accounts.watched_contracts (
  contract BYTEA PRIMARY KEY
);
```
`contract.equivalents` is used to manually map contract addresses which represent the same token and need to be tracked
as such. For example, TrueUSD as shown above has a proxy contract `0x0000000000085d4780B73119b644AE5ecd22b376`
contract that was recently upgraded to from a direct implementation at `0x8dd5fbCe2F6a956C3022bA3663759011Dd51e73E`. 
These two addresses do not emit each other's events and so to track the balance of TrueUSD we are configuring our
transformer to watch events emitted from both these addresses as though they all belong to `0x0000000000085d4780B73119b644AE5ecd22b376`.

`account.start` is used to specify when to begin watching events and producing token and eth balance records for the user accounts,
this needs to be set to a block lower than the deployment block of any tokens we want to track. Additionally, this block number must fall within
the contiguous set of unchecked_headers (this is important if we need to restart a sync, we will need to restart from the lowest unchecked header)

`account.addresses` is used to specify which user account addresses we want to track and index ETH balance and token balance
records for. This can be updated at runtime by adding new addresses to the `accounts.addresses` table in Postgres: 

```postgresql
CREATE TABLE accounts.addresses (
  address BYTEA PRIMARY KEY
);
```

Currently, this config's `ipcPath` needs to point to an archival node endpoint in order to track ETH balances, this will be deprecated
by the use of state diff data in the near future.

To expose the transformed data over Postgraphile, we need to modify our Postgraphile [config.ts](https://github.com/vulcanize/maker-vulcanizedb/blob/staging/postgraphile/src/server/config.ts#L42)
and also [config.spec.ts](https://github.com/vulcanize/maker-vulcanizedb/blob/staging/postgraphile/spec/server/config.spec.ts) to include the "accounts" schema
e.g. (`["public", "accounts"]`). After this, we should be able to expose graphQL endpoints as usual.


# Output

The transformer processes value transfer events from all contract addresses into uniform records of the form:

```postgresql
CREATE TABLE accounts.token_value_transfers (
  id           SERIAL PRIMARY KEY,
  header_id    INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  block_number BIGINT NOT NULL,
  name         VARCHAR NOT NULL CHECK (name <> ''),
  dst          VARCHAR(42),
  src          VARCHAR(42),
  amount       NUMERIC,
  contract     VARCHAR(42) NOT NULL,
  log_idx      INTEGER NOT NULL,
  tx_idx       INTEGER NOT NULL,
  raw_log      JSONB,
  UNIQUE (header_id, tx_idx, log_idx)
);
```

A view on a join of these records with the `accounts.addresses` table, the `accounts.watched_contracts` table, and
the `public.headers` table is used to construct our users' token balance records:

```postgresql
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
```

Which produces a view equivalent to the below table:

```postgresql
CREATE TABLE accounts.address_token_balances (
  id                          SERIAL PRIMARY KEY,
  address_hash                BYTEA NOT NULL,
  block_number                BIGINT NOT NULL,
  token_contract_address_hash BYTEA NOT NULL,
  value                       NUMERIC,
  UNIQUE (address_hash, block_number, token_contract_address_hash)
);
```

Additionally, for each user account it is configured with, it fetches their ETH balances and persists them as coin
balance records of the form:

```postgresql
CREATE TABLE accounts.address_coin_balances (
  id                          SERIAL PRIMARY KEY,
  address_hash                BYTEA NOT NULL,
  block_number                BIGINT NOT NULL,
  value                       NUMERIC(100,0),
  value_fetched_at            TIMESTAMP WITHOUT TIME ZONE,
  inserted_at                 TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  updated_at                  TIMESTAMP WITHOUT TIME ZONE,
  UNIQUE (address_hash, block_number)
);
```

It also syncs the users (and only the users) transactions and receipts into the core vDB `light_sync_transactions` table

```postgresql
CREATE TABLE light_sync_transactions (
  id          SERIAL PRIMARY KEY,
  header_id   INTEGER NOT NULL REFERENCES headers(id) ON DELETE CASCADE,
  hash        VARCHAR(66),
  gaslimit    NUMERIC,
  gasprice    NUMERIC,
  input_data  BYTEA,
  nonce       NUMERIC,
  raw         BYTEA,
  tx_from     VARCHAR(44),
  tx_index    INTEGER,
  tx_to       VARCHAR(44),
  "value"     NUMERIC,
  UNIQUE (header_id, hash)
);
```

```postgresql
CREATE TABLE light_sync_receipts(
  id                  SERIAL PRIMARY KEY,
  transaction_id      INTEGER NOT NULL REFERENCES light_sync_transactions(id) ON DELETE CASCADE,
  header_id           INTEGER NOT NULL REFERENCES headers(id) ON DELETE CASCADE,
  contract_address    VARCHAR(42),
  cumulative_gas_used NUMERIC,
  gas_used            NUMERIC,
  state_root          VARCHAR(66),
  status              INTEGER,
  tx_hash             VARCHAR(66),
  UNIQUE(header_id, transaction_id)
);
```

# Contributing
If you notice a value transfer type event is missing from the [ones we are already tracking](https://github.com/vulcanize/account_transformers/blob/master/transformers/account/shared/constants/signatures.go#L56),
please feel free to submit a PR to introduce the event or submit an issue to note it for inclusion.
