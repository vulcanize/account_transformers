-- +goose Up
ALTER TABLE transactions
  DROP CONSTRAINT fk_test;

ALTER TABLE transactions
  ADD CONSTRAINT blocks_fk
FOREIGN KEY (block_id)
REFERENCES blocks (id)
ON DELETE CASCADE;


-- +goose Down
ALTER TABLE transactions
  DROP CONSTRAINT blocks_fk;

ALTER TABLE transactions
  ADD CONSTRAINT fk_test
FOREIGN KEY (block_id)
REFERENCES blocks (id);
