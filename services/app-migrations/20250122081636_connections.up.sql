CREATE TABLE IF NOT EXISTS app.connections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    connector_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    organization_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    is_deleted BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (connector_id) REFERENCES app.connectors(id),
    FOREIGN KEY (organization_id) REFERENCES app.organizations(organization_id)
);

CREATE INDEX idx_connections_organization_id ON app.connections(organization_id);

INSERT INTO "app"."resource_types" ("name", "table_name") VALUES ('connection', 'connections');

INSERT INTO "app"."resource_privileges" ("privilege", "resource_type", "description") VALUES
('admin', 'connection', 'Admin can do all actions on a connection'),
('viewer', 'connection', 'Viewer can view a connection');
