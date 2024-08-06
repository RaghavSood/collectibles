-- +goose Up
-- +goose StatementBegin
CREATE TABLE telegram_subscription (
    chat_id BIGINT NOT NULL,
    scope TEXT NOT NULL,
    slug TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE UNIQUE INDEX telegram_subscription_chat_id_scope_slug_idx ON telegram_subscription (chat_id, scope, slug);
CREATE INDEX telegram_subscription_active_idx ON telegram_subscription (active);
CREATE INDEX telegram_subscription_chat_id_idx ON telegram_subscription (chat_id);
CREATE INDEX telegram_subscription_scope_idx ON telegram_subscription (scope, slug);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE telegram_subscription;
-- +goose StatementEnd
