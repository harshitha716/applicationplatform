ALTER TABLE app.widget_instances
ADD COLUMN display_config jsonb DEFAULT '{}'::jsonb;