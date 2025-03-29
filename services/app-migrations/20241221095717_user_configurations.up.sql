CREATE TABLE IF NOT EXISTS app.user_configurations (
    user_configuration_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES app.users (user_id) ON DELETE CASCADE,
    resource_type TEXT REFERENCES app.resource_types(name),
    resource_id uuid NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    configuration JSONB NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_user_configurations_user_id ON app.user_configurations (user_id);
CREATE INDEX IF NOT EXISTS idx_user_configurations_resource_type_resource_id ON app.user_configurations (resource_type, resource_id);