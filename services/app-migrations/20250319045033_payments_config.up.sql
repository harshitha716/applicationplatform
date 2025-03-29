CREATE TABLE IF NOT EXISTS app.payments_config (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id uuid NOT NULL REFERENCES app.organizations(organization_id) ON DELETE CASCADE,
    accounts_dataset_id uuid NOT NULL,
    mapping_config jsonb NOT NULL,
    status TEXT NOT NULL,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_payments_config_organization_id ON app.payments_config (organization_id);


-- insert enum entries
INSERT INTO app.resource_types (name, table_name) VALUES('payments', 'payments_config');

-- add payments resource privileges
INSERT INTO "app"."resource_privileges" ("privilege", "resource_type", "description") VALUES
('admin','payments', 'Admin of payments'),
('initiator','payments', 'Initiator of payments'),
('viewer','payments', 'Viewer of payments');
