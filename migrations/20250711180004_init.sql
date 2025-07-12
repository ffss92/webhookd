-- +goose Up
-- +goose StatementBegin
CREATE TABLE "subscribers" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" TEXT NOT NULL,
    "metadata" JSONB NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX "subscribers_metadata_idx" ON "subscribers" USING GIN ("metadata");

CREATE TABLE "endpoints" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "label" TEXT NOT NULL,
    "url" TEXT NOT NULL,
    "secret" TEXT NOT NULL,
    "disabled" BOOLEAN NOT NULL DEFAULT FALSE,
    "filter_types" TEXT[] NOT NULL,
    "subscriber_id" UUID NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("subscriber_id") REFERENCES "subscribers"("id") ON DELETE CASCADE
);
CREATE INDEX "endpoints_subscriber_idx" ON "endpoints"("subscriber_id");
CREATE INDEX "endpoints_filter_types_idx" ON "endpoints" USING GIN ("filter_types");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "endpoints";
DROP TABLE "subscribers";
-- +goose StatementEnd
