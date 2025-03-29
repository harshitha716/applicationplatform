CREATE TABLE IF NOT EXISTS app.organization_sso_configs (
    organization_sso_config_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id uuid NOT NULL REFERENCES app.organizations(organization_id),
    sso_provider_id text NOT NULL,
    sso_provider_name text NOT NULL,
    sso_config jsonb NOT NULL,
    email_domain text NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

DROP TRIGGER IF EXISTS insert_user_into_organizations ON app.users;