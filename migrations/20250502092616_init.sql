-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id BIGINT PRIMARY KEY,
    username VARCHAR(255) NOT NULL
);

CREATE TABLE history (
    id BIGINT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL ,
    favorites_count INT NOT NULL,
    tweet_count INT NOT NULL,
    FOREIGN KEY (id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE history;
DROP TABLE users;
-- +goose StatementEnd
