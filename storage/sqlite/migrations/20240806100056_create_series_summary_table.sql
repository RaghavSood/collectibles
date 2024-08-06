-- +goose Up
-- +goose StatementBegin
CREATE TABLE series_summary_c (
  slug TEXT NOT NULL PRIMARY KEY,
  name TEXT NOT NULL,
  item_count INTEGER NOT NULL,
  tvl INTEGER NOT NULL,
  unfunded INTEGER NOT NULL,
  redeemed INTEGER NOT NULL,
  unredeemed INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE series_summary_c;
-- +goose StatementEnd
