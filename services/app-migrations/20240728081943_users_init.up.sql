CREATE TABLE IF NOT EXISTS app.users (
    user_id uuid PRIMARY KEY REFERENCES kratos.identities (id) ON DELETE CASCADE,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- create a function to set a trigger on kratos.identity such that every identity gets registered in the users table
CREATE OR REPLACE FUNCTION kratos.insert_identity_into_users()
RETURNS TRIGGER AS $$
BEGIN
    -- Insert a corresponding entry into the users table
    INSERT INTO app.users (user_id)
    VALUES (NEW.id);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_insert_kratos_identity
AFTER INSERT ON kratos.identities
FOR EACH ROW
EXECUTE FUNCTION kratos.insert_identity_into_users();
