CREATE TABLE IF NOT EXISTS "app"."audit_logs" (
    "audit_log_id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "kind" TEXT NOT NULL,
    "organization_id" UUID NOT NULL REFERENCES "app"."organizations" ("organization_id"),
    "ip_address" TEXT NOT NULL,
    "user_email" TEXT NOT NULL,
    "user_agent" TEXT NOT NULL,
    "resource_type" TEXT NOT NULL REFERENCES "app"."resource_types" ("name"),
    "resource_id" UUID NOT NULL,
    "event_name" TEXT NOT NULL,
    "payload" JSONB NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now()
);
