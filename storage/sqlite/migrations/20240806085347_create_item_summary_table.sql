-- +goose Up
-- +goose StatementBegin
CREATE TABLE item_summary_c (
  sku TEXT NOT NULL PRIMARY KEY,
  serial TEXT,
  series_name TEXT NOT NULL,
  series_slug TEXT NOT NULL,
  tvl INTEGER NOT NULL,
  unspent INTEGER NOT NULL,
  spent INTEGER NOT NULL,
  total_received INTEGER NOT NULL,
  total_spent INTEGER NOT NULL,
  unfunded INTEGER NOT NULL,
  redeemed INTEGER NOT NULL,
  unredeemed INTEGER NOT NULL
);

CREATE INDEX item_summary_c_series_slug ON item_summary_c(series_slug);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE item_summary_c;
-- +goose StatementEnd
