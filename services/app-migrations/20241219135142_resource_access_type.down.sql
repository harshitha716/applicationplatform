

ALTER TABLE app.resource_access DROP COLUMN access_level_type;

DROP TABLE app.resource_access_level_types;

DROP VIEW app.users_with_traits;