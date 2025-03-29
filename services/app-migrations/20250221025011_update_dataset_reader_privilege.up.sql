
INSERT INTO app.resource_privileges (privilege, resource_type, description)
SELECT 'viewer', 'dataset', 'View dataset schema and data'
WHERE NOT EXISTS (
    SELECT 1 FROM app.resource_privileges WHERE privilege = 'viewer' AND resource_type = 'dataset'
);

UPDATE app.resource_audience_policies
SET privilege = 'viewer'
WHERE privilege = 'data_reader';
