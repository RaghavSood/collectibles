-- +goose Up
-- +goose StatementBegin
CREATE TABLE address_summary_c (
  address TEXT NOT NULL,
  sku TEXT NOT NULL,
  series_slug TEXT NOT NULL,
  serial TEXT,
  unspent INTEGER NOT NULL,
  spent INTEGER NOT NULL,
  first_active DATETIME,
  redeemed_on DATETIME,
  total_value INTEGER
);

CREATE INDEX address_summary_c_sku ON address_summary_c(sku);
CREATE INDEX address_summary_c_series_slug ON address_summary_c(series_slug);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE address_summary_c;
-- +goose StatementEnd
