-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_outpoints_script_spending_txid ON outpoints (script, spending_txid);
CREATE INDEX idx_addresses_sku_script ON addresses (sku, script);
CREATE INDEX idx_addresses_address ON addresses (address);
CREATE INDEX idx_addresses_script ON addresses (script);
CREATE INDEX idx_items_series_slug_sku ON items (series_slug, sku);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_outpoints_script_spending_txid;
DROP INDEX idx_addresses_sku_script;
DROP INDEX idx_addresses_address;
DROP INDEX idx_addresses_script;
DROP INDEX idx_items_series_slug_sku;
-- +goose StatementEnd
