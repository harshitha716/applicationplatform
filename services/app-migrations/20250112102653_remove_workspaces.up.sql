ALTER TABLE app.pages DROP COLUMN workspace_id;

DROP TABLE app.workspaces;

ALTER TABLE app.pages ADD COLUMN organization_id UUID NOT NULL REFERENCES app.organizations(organization_id);