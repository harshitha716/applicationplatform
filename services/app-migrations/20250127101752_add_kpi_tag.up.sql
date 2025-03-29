INSERT INTO app.widgets (
    widget_id,
    name,
    type,
    template_schema,
    created_at,
    updated_at,
    deleted_at
) 
VALUES (
    '55555555-5555-5555-5555-555555555555',
    'KPI',
    'kpi',
    '{
        "mappings": {
            "dataset_id": { "required": true },
            "fields": {
                "primary_value": {
                    "name": "Primary Value",
                    "type": "measure",
                    "required": true,
                    "description": "Main value to display",
                    "allowed_types": ["number"],
                    "allowed_aggregations": ["sum", "avg", "max", "min", "count", "count_distinct"],
                    "cardinality": "single"
                },
                "comparison_value": {
                    "name": "Comparison Value",
                    "type": "measure",
                    "required": false,
                    "description": "Value to compare against (e.g., previous period)",
                    "allowed_types": ["number"],
                    "allowed_aggregations": ["sum", "avg", "max", "min"]
                },
                "comparison_period": {
                    "name": "Comparison Period",
                    "type": "dimension",
                    "required": false,
                    "description": "Time period to compare against",
                    "allowed_types": ["string"],
                    "options": [
                        {"value": "previous_week", "label": "Previous Week"},
                        {"value": "previous_month", "label": "Previous Month"},
                        {"value": "previous_year", "label": "Previous Year"},
                        {"value": "custom", "label": "Custom Date Range"}
                    ],
                    "cardinality": "single"
                }
            }
        },
        "mappings_cardinality": "single",
        "supports_group_by": false
    }'::jsonb,
    NOW(),
    NOW(),
    NULL
);