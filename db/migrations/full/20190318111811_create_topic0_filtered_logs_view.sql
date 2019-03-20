-- +goose Up
CREATE VIEW accounts.topic0_filtered_logs AS
  SELECT
    topic0_filters.name,
    logs.id,
    block_number,
    logs.address,
    tx_hash,
    index,
    logs.topic0,
    logs.topic1,
    logs.topic2,
    logs.topic3,
    data,
    receipt_id
  FROM accounts.topic0_filters
    CROSS JOIN block_stats
    JOIN logs ON logs.topic0 = topic0_filters.topic0
                 AND logs.block_number >= coalesce(topic0_filters.from_block, block_stats.min_block)
                 AND logs.block_number <= coalesce(topic0_filters.to_block, block_stats.max_block)
    WHERE (topic0_filters.topic1 = logs.topic1 OR topic0_filters.topic1 ISNULL)
        AND (topic0_filters.topic2 = logs.topic2 OR topic0_filters.topic2 ISNULL)
        AND (topic0_filters.topic3 = logs.topic3 OR topic0_filters.topic3 ISNULL);

-- +goose Down
DROP VIEW accounts.topic0_filtered_logs;
