-- +goose Up
CREATE TABLE accounts.value_transfer_events (
  id         SERIAL PRIMARY KEY,
  name       VARCHAR NOT NULL CHECK (name <> ''),
  dst        VARCHAR(66) NOT NULL,
  src        VARCHAR(66) NOT NULL,
  amount     NUMERIC NOT NULL,
  contract   VARCHAR(66) NOT NULL,
  log_id     INTEGER NOT NULL,
  CONSTRAINT log_id_fk FOREIGN KEY (log_id)
  REFERENCES public.logs (id)
);

-- +goose Down
DROP TABLE accounts.value_transfer_events;
