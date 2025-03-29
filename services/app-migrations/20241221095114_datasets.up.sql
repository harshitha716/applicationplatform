CREATE TABLE IF NOT EXISTS app.dataset_types (
    name TEXT NOT NULL PRIMARY KEY,
    description TEXT NOT NULL
);

INSERT INTO app.dataset_types (name, description) VALUES 
    ('bronze', 'Bronze dataset'),
    ('source', 'Source dataset'),
    ('materialised', 'Materialised dataset');

CREATE TABLE IF NOT EXISTS app.datasets (
    dataset_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    type TEXT NOT NULL REFERENCES app.dataset_types (name) ON DELETE CASCADE,
    organization_id uuid NOT NULL REFERENCES app.organizations (organization_id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by uuid NOT NULL REFERENCES app.users (user_id) ON DELETE CASCADE,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    metadata JSONB NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_datasets_dataset_id ON app.datasets (dataset_id);