-- +goose Up
-- +goose StatementBegin
CREATE TABLE block_notification_queue (
    block_height NUMERIC NOT NULL,
    chain TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE block_notification_queue;
-- +goose StatementEnd
