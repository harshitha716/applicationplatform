-- workspaces table
CREATE TABLE IF NOT EXISTS app.workspaces (
    workspace_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    organization_id UUID NOT NULL REFERENCES app.organizations(organization_id),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- pages table
CREATE TABLE IF NOT EXISTS app.pages (
    page_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    workspace_id UUID NOT NULL REFERENCES app.workspaces(workspace_id),
    fractional_index DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- function and trigger to set the fractional index on insert
CREATE OR REPLACE FUNCTION auto_increment_page_fractional_index()
RETURNS TRIGGER AS $$
BEGIN
    -- Increment fractional_index for the given workspace_id
    NEW.fractional_index := COALESCE(
        (SELECT MAX(fractional_index) + 1.0 
         FROM app.pages 
         -- WHERE workspace_id = NEW.workspace_id
        ),
        1.0 -- Assign 1.0 for the first page in the workspace
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_page_fractional_index
BEFORE INSERT ON app.pages
FOR EACH ROW
WHEN (NEW.fractional_index = 0.0)
EXECUTE FUNCTION auto_increment_page_fractional_index();

-- sheets table
CREATE TABLE IF NOT EXISTS app.sheets (
    sheet_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    page_id UUID NOT NULL REFERENCES app.pages(page_id),
    fractional_index DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- function and trigger to set the fractional index on insert
CREATE OR REPLACE FUNCTION auto_increment_sheet_fractional_index()
RETURNS TRIGGER AS $$
BEGIN
    -- Calculate the next fractional_index within the same page_id
    NEW.fractional_index := COALESCE(
        (SELECT MAX(fractional_index) + 1.0
         FROM app.sheets
         WHERE page_id = NEW.page_id),
        1.0 -- Assign 1.0 for the first sheet in the page
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_sheet_fractional_index
BEFORE INSERT ON app.sheets
FOR EACH ROW
WHEN (NEW.fractional_index = 0.0)
EXECUTE FUNCTION auto_increment_sheet_fractional_index();