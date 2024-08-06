-- +goose Up
-- +goose StatementBegin
CREATE TABLE item_address_summary_c (
  sku TEXT NOT NULL,
  serial TEXT,
  addresses TEXT NOT NULL,
  first_active DATETIME,
  redeemed_on DATETIME,
  total_value INTEGER NOT NULL,
  series_name TEXT NOT NULL,
  series_slug TEXT NOT NULL
);

CREATE INDEX item_address_summary_c_sku ON item_address_summary_c(sku);
CREATE INDEX item_address_summary_c_series_slug ON item_address_summary_c(series_slug);
CREATE INDEX item_address_summary_c_redeemed_on ON item_address_summary_c(redeemed_on);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE item_address_summary_c;
-- +goose StatementEnd
