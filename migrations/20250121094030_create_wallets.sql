-- +goose Up
-- +goose StatementBegin
CREATE TABLE wallets (
                         id UUID PRIMARY KEY,
                         amount INT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wallets;
-- +goose StatementEnd
