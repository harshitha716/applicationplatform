
-- create a view that joins users table with kratos.identities table and returns user_id and email, name from kratos.identities.traits
CREATE OR REPLACE VIEW app.users_with_traits AS
SELECT
    u.user_id,
    i.traits->>'email' AS email,  -- Extracting email from jsonb traits
    COALESCE(i.traits->>'name', '') AS name, -- Use empty string if name is null
    i.traits  -- Including the full traits as jsonb
FROM
    app.users u
JOIN
    kratos.identities i ON u.user_id = i.id;

-- Add a resource_access_level_types table
CREATE TABLE IF NOT EXISTS app.resource_access_level_types (
    name TEXT NOT NULL PRIMARY KEY
);

-- INSERT current allowed access levels
INSERT INTO app.resource_access_level_types (name) VALUES ('read'),('write');

-- Add a resource_access_level_type column to the app.resource_access table that references app.resource_access_level_types
ALTER TABLE app.resource_access
ADD COLUMN access_level_type TEXT REFERENCES app.resource_access_level_types(name);