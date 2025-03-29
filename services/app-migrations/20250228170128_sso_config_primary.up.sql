ALTER TABLE app.organization_sso_configs ADD COLUMN is_primary BOOLEAN NOT NULL DEFAULT TRUE;

-- change default to false
ALTER TABLE app.organization_sso_configs ALTER COLUMN is_primary SET DEFAULT FALSE;
