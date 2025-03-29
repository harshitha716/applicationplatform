INSERT INTO app.dataset_types (name, description)
SELECT 'bronze', 'Bronze dataset'
WHERE NOT EXISTS (
    SELECT 1 FROM app.dataset_types WHERE name = 'bronze'
);
