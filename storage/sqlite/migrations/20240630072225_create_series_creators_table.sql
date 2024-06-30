-- +goose Up
-- +goose StatementBegin
CREATE TABLE series_creators (
  series_slug TEXT NOT NULL,
  creator_slug TEXT NOT NULL,
  PRIMARY KEY (series_slug, creator_slug),
  FOREIGN KEY (series_slug) REFERENCES series (slug) ON DELETE CASCADE,
  FOREIGN KEY (creator_slug) REFERENCES creators (slug) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE series_creators;
-- +goose StatementEnd
