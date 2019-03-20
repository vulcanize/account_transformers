-- +goose Up
CREATE OR REPLACE FUNCTION accounts.update_checked_blocks(address VARCHAR(66)) RETURNS VOID AS $$
BEGIN
  ALTER TABLE accounts.checked_blocks ADD COLUMN IF NOT EXISTS address INTEGER NOT NULL DEFAULT 0;
END;
$$ LANGUAGE plpgsql


-- +goose Down
DROP FUNCTION accounts.update_checked_blocks;
