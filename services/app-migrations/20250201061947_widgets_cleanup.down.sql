-- 8) Remove default_filters column
ALTER TABLE app.widget_instances DROP COLUMN IF EXISTS default_filters;

-- Remove the new index
DROP INDEX app.idx_widget_instances_sheet_id_widget_type;

-- 7) Remove foreign key constraint to widget_types
ALTER TABLE app.widget_instances DROP CONSTRAINT widget_instances_widget_type_fkey;

-- 6) Add widget_id column back
ALTER TABLE app.widget_instances ADD COLUMN widget_id INTEGER;

-- Populate widget_id from widgets table based on type
UPDATE app.widget_instances wi 
SET widget_id = w.widget_id 
FROM app.widgets w 
WHERE wi.widget_type = w.type;

-- Make widget_id NOT NULL after populating
ALTER TABLE app.widget_instances ALTER COLUMN widget_id SET NOT NULL;

-- 5) Remove widget_type column
ALTER TABLE app.widget_instances DROP COLUMN widget_type;

-- 4) Restore original primary key on widgets
ALTER TABLE app.widgets DROP CONSTRAINT widgets_pkey;
ALTER TABLE app.widgets ADD PRIMARY KEY (widget_id);

-- Recreate the original index
CREATE INDEX idx_widget_instances_sheet_id_widget_id 
ON app.widget_instances(sheet_id, widget_id);

-- 1) Restore foreign key constraint
ALTER TABLE app.widget_instances 
ADD CONSTRAINT widget_instances_widget_id_fkey 
FOREIGN KEY (widget_id) REFERENCES app.widgets(widget_id);