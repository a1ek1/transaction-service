-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS service_state (
                                             key VARCHAR(255) PRIMARY KEY,
                                             value VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS service_state;
-- +goose StatementEnd
