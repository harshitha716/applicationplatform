-- general table for file uploads
CREATE TABLE IF NOT EXISTS app.file_uploads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    file_type TEXT NOT NULL,
    organization_id UUID NOT NULL REFERENCES app.organizations(organization_id),
    uploaded_by_user_id UUID NOT NULL REFERENCES app.users(user_id),
    presigned_url TEXT NOT NULL,
    expiry TIMESTAMP WITH TIME ZONE NOT NULL,
    storage_provider TEXT NOT NULL,
    storage_bucket TEXT NOT NULL,
    storage_file_path TEXT,
    status TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS app.dataset_file_uploads (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	dataset_id uuid NOT NULL REFERENCES app.datasets(dataset_id),
	file_upload_id uuid NOT NULL REFERENCES app.file_uploads(id),
	status TEXT NOT NULL,
	metadata jsonb NOT NULL
);