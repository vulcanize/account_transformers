[VulcanizeDB](https://github.com/vulcanize/maker-vulcanizedb) transformers for watching ETH and token value transfers

This repo contains transformers for indexing of ETH and token balances for account addresses. [This](https://github.com/vulcanize/account_transformers/tree/master/transformers/account/light) transformer 
works by filtering and indexing all token transfer events (of a set of [known types](https://github.com/vulcanize/account_transformers/tree/master/transformers/account/shared/constants)) from all contract addresses,
compiling this information into token balance records for provided account addresses. It then polls an archival Eth node to
retrieve balances for these accounts and generate eth balance records.

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
to begin syncing headers into Postgres. It is vital that this sync process begins at block 0.

Once `lightSync` has begun, we can run the `composeAndExecute` command to compose and execute our account transformer. To
do so, we use a normal `compose` [config](https://github.com/vulcanize/maker-vulcanizedb#contractwatcher) with two additional parameter maps:

```toml
[equivalents]
    0x0000000000085d4780B73119b644AE5ecd22b376 = [
       "0x8dd5fbCe2F6a956C3022bA3663759011Dd51e73E"
    ]

[account]
    addresses = [
        "0x48E78948C80e9f8F53190DbDF2990f9a69491ef4",
        "0x009C1E8674038605C5AE33C74f13bC528E1222B5"
    ]
```

The first, `equivalents` is used to manually map contract addresses which represent the same token and need to be tracked
as such. This is the case as in the example with TrueUSD, where `0x0000000000085d4780B73119b644AE5ecd22b376` is a proxy
contract that was upgraded to from the direct implementation at `0x8dd5fbCe2F6a956C3022bA3663759011Dd51e73E`, meaning events
before the upgrade were emitted from `0x8dd5fbCe2F6a956C3022bA3663759011Dd51e73E` whereas they are now emitted from `0x0000000000085d4780B73119b644AE5ecd22b376`
and we need to know to treat them equivalently in order to properly index TrueUSD token balances.

The second, `account` is used to specify which user account addresses we want to track and index ETH balance and token balance
records for. This can also be updated at runtime by adding new addresses to the `accounts.addresses` table in Postgres.

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
  dst          VARCHAR(66),
  src          VARCHAR(66),
  amount       NUMERIC,
  contract     VARCHAR(66) NOT NULL,
  log_idx      INTEGER NOT NULL,
  tx_idx       INTEGER NOT NULL,
  raw_log      JSONB,
  UNIQUE (header_id, tx_idx, log_idx)
);
```

For each user account it is configured with, it then filters through these at a given blockheight to tally up token
balances for those accounts and persist them as token balance records of the form:

```postgresql
CREATE TABLE accounts.address_token_balances (
  id                          SERIAL PRIMARY KEY,
  address_hash                BYTEA NOT NULL,
  block_number                BIGINT NOT NULL,
  token_contract_address_hash BYTEA NOT NULL,
  value                       NUMERIC,
  value_fetched_at            TIMESTAMP WITHOUT TIME ZONE,
  inserted_at                 TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  updated_at                  TIMESTAMP WITHOUT TIME ZONE,
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

