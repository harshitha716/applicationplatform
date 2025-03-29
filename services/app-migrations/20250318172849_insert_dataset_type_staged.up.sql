INSERT INTO app.dataset_types (name, description)
SELECT 'staged', 'Hidden dataset from Listings'
WHERE NOT EXISTS (
    SELECT 1 FROM app.dataset_types WHERE name = 'staged'
);
