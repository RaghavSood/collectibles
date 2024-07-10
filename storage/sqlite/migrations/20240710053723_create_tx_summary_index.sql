-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_outpoints_script_spending_block_height_time ON outpoints(script, spending_block_height, spending_block_time);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_outpoints_script_spending_block_height_time;
-- +goose StatementEnd
