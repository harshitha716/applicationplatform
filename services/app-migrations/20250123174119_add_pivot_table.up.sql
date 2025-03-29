-- Update Bar Chart Widget
UPDATE app.widgets
SET template_schema = '{
    "mappings": {
        "dataset_id": {
            "required": true
        },
        "fields": {
            "x_axis": {
                "name": "X-axis",
                "type": "dimension",
                "required": true,
                "description": "Categories to compare",
                "allowed_types": ["string", "datetime"],
                "cardinality": "single"
            },
            "y_axis": {
                "name": "Y-axis",
                "type": "measure",
                "required": true,
                "description": "Values to show",
                "allowed_types": ["number"],
                "allowed_aggregations": ["sum", "avg", "count", "max", "min"],
                "cardinality": "single"
            },
            "group_by": {
                "max_fields": 2,
                "allowed_types": ["string", "datetime"],
                "required": false,
                "cardinality": "multiple"
            }
        }
    },
    "mappings_cardinality": "single",
    "supports_group_by": true
}'::jsonb
WHERE type = 'bar_chart';

-- Update Pie Chart Widget
UPDATE app.widgets
SET template_schema = '{
    "mappings": {
        "dataset_id": {
            "required": true
        },
        "fields": {
            "slices": {
                "name": "Slices",
                "type": "dimension",
                "required": true,
                "description": "Categories to compare",
                "allowed_types": ["string", "datetime"],
                "cardinality": "single"
            },
            "values": {
                "name": "Values",
                "type": "measure",
                "required": true,
                "description": "Values to show",
                "allowed_types": ["number"],
                "allowed_aggregations": ["sum", "avg", "count", "max", "min"],
                "cardinality": "single"
            }
        }
    },
    "mappings_cardinality": "single",
    "supports_group_by": false
}'::jsonb
WHERE type = 'pie_chart';

-- Update Line Chart Widget
UPDATE app.widgets
SET template_schema = '{
    "mappings": {
        "dataset_id": {
            "required": true
        },
        "fields": {
            "x_axis": {
                "name": "X-axis",
                "type": "dimension",
                "required": true,
                "description": "Categories to compare",
                "allowed_types": ["string", "datetime"],
                "cardinality": "single"
            },
            "y_axis": {
                "name": "Y-axis",
                "type": "measure",
                "required": true,
                "description": "Values to show",
                "allowed_types": ["number"],
                "allowed_aggregations": ["sum", "avg", "count", "max", "min"],
                "cardinality": "single"
            },
            "group_by": {
                "max_fields": 2,
                "allowed_types": ["string", "datetime"],
                "required": false,
                "cardinality": "multiple"
            }
        }
    },
    "mappings_cardinality": "single",
    "supports_group_by": true
}'::jsonb
WHERE type = 'line_chart';


-- Insert Pivot Table Widget
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
    '44444444-4444-4444-4444-444444444444',  -- Fixed UUID matching the seeds pattern
    'Pivot Table',
    'pivot_table',
    '{
        "mappings": {
            "dataset_id": {
                "required": true
            },
            "fields": {
                "rows": {
                    "name": "Rows",
                    "type": "dimension",
                    "required": true,
                    "description": "Row headers",
                    "allowed_types": ["string", "datetime"],
                    "cardinality": "multiple"
                },
                "columns": {
                    "name": "Columns",
                    "type": "dimension",
                    "required": true,
                    "description": "Column headers",
                    "allowed_types": ["string", "datetime"],
                    "cardinality": "multiple"
                },
                "values": {
                    "name": "Values",
                    "type": "measure",
                    "required": true,
                    "description": "Values to aggregate",
                    "allowed_types": ["number"],
                    "allowed_aggregations": ["sum", "avg", "count", "max", "min"],
                    "cardinality": "multiple"
                }
            }
        },
        "mappings_cardinality": "multiple",
        "supports_group_by": false
    }'::jsonb,
    NOW(),
    NOW(),
    NULL
);

-- Dropping DatasetId column, it will be only stored within the mappings
ALTER TABLE app.widget_instances DROP COLUMN IF EXISTS dataset_id;