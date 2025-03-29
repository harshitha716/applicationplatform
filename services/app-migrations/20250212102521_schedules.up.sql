CREATE TABLE IF NOT EXISTS app.schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    organization_id UUID NOT NULL,
    schedule_group VARCHAR(255) NOT NULL,
    connector_id UUID NOT NULL,
    connection_id UUID NOT NULL,
    temporal_workflow_id VARCHAR(65535) NOT NULL,
    status VARCHAR(255) NOT NULL,
    config JSONB NOT NULL,
    cron_schedule VARCHAR(255),
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_connector FOREIGN KEY (connector_id) REFERENCES app.connectors(id),
    CONSTRAINT fk_connection FOREIGN KEY (connection_id) REFERENCES app.connections(id),
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES app.organizations(organization_id)
);
