-- +goose Up
-- +goose StatementBegin
CREATE TABLE "messages" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "type" TEXT NOT NULL,
    "data" JSONB NOT NULL,
    "tags" TEXT[] NOT NULL DEFAULT '{}',
    "subscriber_id" UUID NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("subscriber_id") REFERENCES "subscribers"("id") ON DELETE CASCADE
);
CREATE INDEX "messages_subscriber_id_idx" ON "messages"("subscriber_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "messages";
-- +goose StatementEnd
