-- +goose Up
-- +goose StatementBegin
ALTER TABLE outpoints ADD COLUMN chain TEXT NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
