CREATE TABLE IF NOT EXISTS app.rules (
    rule_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id uuid NOT NULL REFERENCES app.organizations(organization_id),
    dataset_id uuid NOT NULL REFERENCES app.datasets(dataset_id),
    "column" TEXT NOT NULL,
    value TEXT NOT NULL,
    filter_config JSONB NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    priority INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    created_by uuid NOT NULL REFERENCES app.users(user_id),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_by uuid NOT NULL REFERENCES app.users(user_id),
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by uuid REFERENCES app.users(user_id)
);

CREATE UNIQUE INDEX idx_rules_organization_dataset_column_value ON app.rules (organization_id, dataset_id, "column", priority) WHERE deleted_at IS NULL;