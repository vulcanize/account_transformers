-- +goose Up
CREATE TABLE accounts.checked_logs (
  id                 SERIAL PRIMARY KEY,
  log_id           INTEGER UNIQUE NOT NULL REFERENCES public.logs (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE accounts.checked_logs;