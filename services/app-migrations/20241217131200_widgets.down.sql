-- Drop indexes
DROP INDEX IF EXISTS idx_widget_instances_deleted_at;
DROP INDEX IF EXISTS idx_widget_instances_sheet_id_widget_id;

-- Drop tables
DROP TABLE IF EXISTS app.widget_instances;
DROP TABLE IF EXISTS app.widgets;
DROP TABLE IF EXISTS app.widget_types;