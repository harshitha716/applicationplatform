-- 1) Remove foreign key constraint from widget_instances to widgets
ALTER TABLE app.widget_instances DROP CONSTRAINT widget_instances_widget_id_fkey;

-- Drop the existing index that includes widget_id
DROP INDEX app.idx_widget_instances_sheet_id_widget_id;

-- 4) Make type as the primary key of the widgets table
ALTER TABLE app.widgets DROP CONSTRAINT widgets_pkey;
ALTER TABLE app.widgets ADD PRIMARY KEY (type);

-- 5) Add widget_type column and populate it
ALTER TABLE app.widget_instances ADD COLUMN widget_type TEXT;
UPDATE app.widget_instances wi 
SET widget_type = w.type 
FROM app.widgets w 
WHERE wi.widget_id = w.widget_id;

-- Make widget_type NOT NULL after populating
ALTER TABLE app.widget_instances ALTER COLUMN widget_type SET NOT NULL;

-- 6) Remove widget_id column
ALTER TABLE app.widget_instances DROP COLUMN widget_id;
ALTER TABLE app.widgets DROP COLUMN widget_id;

-- 7) Add foreign key constraint from widget_instances to widget_types
ALTER TABLE app.widget_instances 
ADD CONSTRAINT widget_instances_widget_type_fkey 
FOREIGN KEY (widget_type) REFERENCES app.widget_types(name);

-- Create new index to replace the old one
CREATE INDEX idx_widget_instances_sheet_id_widget_type 
ON app.widget_instances(sheet_id, widget_type);

-- 8) Add default_filters column
ALTER TABLE app.widget_instances ADD COLUMN default_filters JSONB DEFAULT '{}'::jsonb;

