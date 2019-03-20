-- +goose Up
CREATE TRIGGER accounts.addresses_trigger
  AFTER INSERT OF address ON accounts.addresses
  EXECUTE PROCEDURE accounts.update_checked_blocks(address);


-- +goose Down
DROP TRIGGER accounts.addresses_trigger;