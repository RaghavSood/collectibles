-- +goose Up
-- +goose StatementBegin
INSERT INTO kv (key, value) VALUES ('block_height', '850150');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM kv WHERE key = 'block_height';
-- +goose StatementEnd
