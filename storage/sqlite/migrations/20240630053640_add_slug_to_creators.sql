-- +goose Up
-- +goose StatementBegin
ALTER TABLE creators ADD COLUMN slug TEXT NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
