-- widget types enum table
CREATE TABLE IF NOT EXISTS app.widget_types (
    name TEXT NOT NULL PRIMARY KEY,
    description TEXT NOT NULL
);

INSERT INTO app.widget_types (name, description) VALUES
    ('bar_chart', 'Bar chart widget'),
    ('line_chart', 'Line chart widget'),
    ('area_chart', 'Area chart widget'),
    ('pie_chart', 'Pie chart widget'),
    ('donut_chart', 'Donut chart widget'),
    ('kpi', 'KPI widget'),
    ('table', 'Table widget'),
    ('pivot_table', 'Pivot table widget');

-- widgets table
CREATE TABLE IF NOT EXISTS app.widgets (
    widget_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    type TEXT NOT NULL REFERENCES app.widget_types(name),
    template_schema JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- widget instances table
CREATE TABLE IF NOT EXISTS app.widget_instances (
    widget_instance_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    widget_id UUID NOT NULL REFERENCES app.widgets(widget_id),
    sheet_id UUID NOT NULL REFERENCES app.sheets(sheet_id),
    title TEXT NOT NULL,
    data_mappings JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_widget_instances_sheet_id_widget_id 
ON app.widget_instances(sheet_id, widget_id);

CREATE INDEX IF NOT EXISTS idx_widget_instances_deleted_at 
ON app.widget_instances(deleted_at);