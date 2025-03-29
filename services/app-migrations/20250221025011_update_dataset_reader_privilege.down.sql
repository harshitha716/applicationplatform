UPDATE app.resource_audience_policies
SET privilege = 'data_reader'
WHERE privilege = 'viewer';

UPDATE app.resource_privileges
SET privilege = 'data_reader'
WHERE privilege = 'viewer';
