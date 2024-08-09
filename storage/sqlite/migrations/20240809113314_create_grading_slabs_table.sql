-- +goose Up
-- +goose StatementBegin
CREATE TABLE grading_slabs (
  sku TEXT NOT NULL,
  service TEXT NOT NULL,
  identifier TEXT,
  grade TEXT,
  view_link TEXT
);

CREATE INDEX grading_slabs_sku_idx ON grading_slabs (sku);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE grading_slabs;
-- +goose StatementEnd
