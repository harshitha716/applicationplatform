-- enum table for kind of resources that can be accessed
CREATE TABLE IF NOT EXISTS app.resource_types (
    name TEXT NOT NULL PRIMARY KEY,
    table_name TEXT
);

-- insert enum entries
INSERT INTO app.resource_types (name, table_name) VALUES 
    ('user', 'users'),
    ('organization', 'organizations'),
    ('organization_membership', 'organization_memberships'),
    ('organization_invitation', 'organization_membership_invitations'),
    ('workspace', 'workspaces'),
    ('sheet', 'sheets'),
    ('page', 'pages'),
    ('dataset', 'datasets');

-- enum table for kind of audience types that can access resources
CREATE TABLE IF NOT EXISTS app.resource_access_audience_types (
    name TEXT NOT NULL PRIMARY KEY
);

-- insert enum entries
INSERT INTO app.resource_access_audience_types (name) 
VALUES 
    ('user'),
    ('organization');

-- create resource access table that defines the access rights of different audiences to different resources
CREATE TABLE IF NOT EXISTS app.resource_access (
    resource_access_id TEXT PRIMARY KEY,
    resource_type TEXT REFERENCES app.resource_types(name),
    resource_id uuid NOT NULL,
 
    access_audience_type TEXT REFERENCES app.resource_access_audience_types(name),
    access_audience_id uuid NOT NULL,

    metadata jsonb NOT NULL DEFAULT '{}',

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
