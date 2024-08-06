-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS message_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    chat_id INTEGER,
    message TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS message_logs;
-- +goose StatementEnd
