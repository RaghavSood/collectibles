-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX creators_slug_unique_index ON creators (slug);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX creators_slug_unique_index;
-- +goose StatementEnd
