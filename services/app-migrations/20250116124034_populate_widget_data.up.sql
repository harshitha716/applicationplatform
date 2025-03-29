-- Insert Bar Chart Widget
WITH formatted_json AS (
  SELECT '{
    "mappings": {
        "x_axis": {
            "name": "X-axis",
            "type": "dimension",
            "required": true,
            "description": "Categories to compare",
            "allowed_types": ["string", "datetime"]
        },
        "y_axis": {
            "name": "Y-axis",
            "type": "measure",
            "required": true,
            "description": "Values to show",
            "allowed_types": ["number"],
            "allowed_aggregations": ["sum", "avg", "count", "max", "min"]
        },
        "group_by": {
            "max_fields": 2,
            "allowed_types": ["string", "datetime"],
            "required": false
        }
    },
    "supports_group_by": true
}'::jsonb AS template_schema
)
INSERT INTO app.widgets (
    widget_id,
    name,
    type,
    template_schema,
    created_at,
    updated_at,
    deleted_at
) 
SELECT 
    '11111111-1111-1111-1111-111111111111'::uuid,
    'Bar Chart',
    'bar_chart',
    template_schema,
    NOW(),
    NOW(),
    NULL
FROM formatted_json;

-- Insert Pie Chart Widget
WITH formatted_json AS (
  SELECT '{
    "mappings": {
        "slices": {
            "name": "Slices",
            "type": "dimension",
            "required": true,
            "description": "Categories to compare",
            "allowed_types": ["string", "datetime"]
        },
        "values": {
            "name": "Values",
            "type": "measure",
            "required": true,
            "description": "Values to show",
            "allowed_types": ["number"],
            "allowed_aggregations": ["sum", "avg", "count", "max", "min"]
        }
    },
    "supports_group_by": false
}'::jsonb AS template_schema
)
INSERT INTO app.widgets (
    widget_id,
    name,
    type,
    template_schema,
    created_at,
    updated_at,
    deleted_at
) 
SELECT 
    '33333333-3333-3333-3333-333333333333'::uuid,
    'Pie Chart',
    'pie_chart',
    template_schema,
    NOW(),
    NOW(),
    NULL
FROM formatted_json;

-- Insert Line Chart Widget
WITH formatted_json AS (
  SELECT '{
    "mappings": {
        "x_axis": {
            "name": "X-axis",
            "type": "dimension",
            "required": true,
            "description": "Categories to compare",
            "allowed_types": ["string", "datetime"]
        },
        "y_axis": {
            "name": "Y-axis",
            "type": "measure",
            "required": true,
            "description": "Values to show",
            "allowed_types": ["number"],
            "allowed_aggregations": ["sum", "avg", "count", "max", "min"]
        },
        "group_by": {
            "max_fields": 2,
            "allowed_types": ["string", "datetime"],
            "required": false
        }
    },
    "supports_group_by": true
}'::jsonb AS template_schema
)
INSERT INTO app.widgets (
    widget_id,
    name,
    type,
    template_schema,
    created_at,
    updated_at,
    deleted_at
) 
SELECT 
    '22222222-2222-2222-2222-222222222222'::uuid,
    'Line Chart',
    'line_chart',
    template_schema,
    NOW(),
    NOW(),
    NULL
FROM formatted_json;