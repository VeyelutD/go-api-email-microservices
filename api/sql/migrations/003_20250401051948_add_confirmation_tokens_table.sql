-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_confirmation_tokens
(
    id         BIGSERIAL PRIMARY KEY,
    email      TEXT      NOT NULL,
    token      TEXT      NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (email) references users (email) ON DELETE CASCADE,
    UNIQUE (email)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_confirmation_tokens;
-- +goose StatementEnd
