-- +goose Up
-- +goose StatementBegin
CREATE TABLE transactions (
                              id UUID PRIMARY KEY,
                              "from" TEXT NOT NULL,
                              "to" TEXT NOT NULL,
                              amount INT NOT NULL,
                              created_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
-- +goose StatementEnd
