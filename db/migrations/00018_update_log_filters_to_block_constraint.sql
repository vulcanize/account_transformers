-- +goose Up
ALTER TABLE log_filters
  DROP CONSTRAINT log_filters_from_block_check1;

ALTER TABLE log_filters
  ADD CONSTRAINT log_filters_to_block_check CHECK (to_block >= 0);


-- +goose Down
ALTER TABLE log_filters
  DROP CONSTRAINT log_filters_to_block_check;

ALTER TABLE log_filters
  ADD CONSTRAINT log_filters_from_block_check1 CHECK (to_block >= 0);
