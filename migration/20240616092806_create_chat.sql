-- +goose Up
-- +goose StatementBegin
CREATE TABLE chats(
    chat_id BIGINT PRIMARY KEY
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE chats;
-- +goose StatementEnd
