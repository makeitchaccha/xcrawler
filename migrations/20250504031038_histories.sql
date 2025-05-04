-- +goose Up
-- +goose StatementBegin
ALTER TABLE history RENAME TO histories;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE histories RENAME TO history;
-- +goose StatementEnd
