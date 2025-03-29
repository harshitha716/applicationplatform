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
    '66666666-6666-6666-6666-666666666666'::uuid,
    'Donut Chart',
    'donut_chart',
    template_schema,
    NOW(),
    NOW(),
    NULL
FROM formatted_json;