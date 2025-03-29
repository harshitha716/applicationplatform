-- create organizations table
CREATE TABLE IF NOT EXISTS app.organizations (
    organization_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    owner_id UUID NOT NULL REFERENCES app.users(user_id),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

-- create organization memberships table
CREATE TABLE IF NOT EXISTS app.organization_memberships (
    organization_membership_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES app.organizations(organization_id),
    member_id UUID NOT NULL REFERENCES app.users(user_id),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    status TEXT NOT NULL
);

-- create organization membership invitations table
CREATE TABLE IF NOT EXISTS app.organization_membership_invitations (
    organization_id UUID NOT NULL REFERENCES app.organizations(organization_id),
    invited_by_id UUID NOT NULL REFERENCES app.users(user_id),
    invitee_id UUID NOT NULL REFERENCES app.users(user_id),
    invited_at TIMESTAMPTZ DEFAULT now(),
    acceptance_status TEXT NOT NULL,
    acceptance_time TIMESTAMPTZ
);

-- create a function to set a trigger on users such that every user gets registered in the organizations table
CREATE OR REPLACE FUNCTION app.insert_user_into_organizations()
RETURNS TRIGGER AS $$
DECLARE
    new_org_id uuid;
BEGIN

    -- define org_id variable uuid

    -- Insert a corresponding entry into the organizations table with the right user id and the default organization name and get the organization id
    INSERT INTO app.organizations (name, owner_id)
    VALUES ('Default', NEW.user_id)
    RETURNING organization_id INTO new_org_id;

    -- Insert an entry into the organization_memberships table with the right organization id and the user id
    INSERT INTO app.organization_memberships (organization_id, member_id, status)
    VALUES (new_org_id, NEW.user_id, 'active');

    RETURN NEW;

END;
$$ LANGUAGE plpgsql;

-- create a trigger to call the function above
CREATE TRIGGER insert_user_into_organizations
AFTER INSERT
ON app.users
FOR EACH ROW
EXECUTE FUNCTION app.insert_user_into_organizations();


