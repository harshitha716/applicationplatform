-- Add dataset_id column to widget_instances table
ALTER TABLE app.widget_instances ADD COLUMN dataset_id UUID REFERENCES app.datasets(dataset_id);