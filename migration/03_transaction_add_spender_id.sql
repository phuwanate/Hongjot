-- +goose Up
-- +goose StatementBegin
ALTER TABLE "transaction" ADD COLUMN spender_id INT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "transaction" DROP COLUMN spender_id;
-- +goose StatementEnd
