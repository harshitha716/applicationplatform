CREATE TABLE IF NOT EXISTS "app"."dataset_actions" (
    "id" UUID PRIMARY KEY,
    "action_id" TEXT NOT NULL,
    "action_type" TEXT NOT NULL,
    "dataset_id" UUID NOT NULL REFERENCES "app"."datasets" ("dataset_id"),
    "organization_id" UUID NOT NULL REFERENCES "app"."organizations" ("organization_id"),
    "status" TEXT NOT NULL,
    "config" JSONB NOT NULL,
    "action_by" UUID NOT NULL REFERENCES "app"."users" ("user_id"),
    "started_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "completed_at" TIMESTAMPTZ
);