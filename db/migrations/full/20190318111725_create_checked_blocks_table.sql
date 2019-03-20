-- +goose Up
CREATE TABLE accounts.checked_blocks (
  id                 SERIAL PRIMARY KEY,
  block_id           INTEGER UNIQUE NOT NULL REFERENCES public.blocks (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE accounts.checked_blocks;