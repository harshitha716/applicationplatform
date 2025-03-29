CREATE TABLE IF NOT EXISTS "app"."teams" (
    "team_id" UUID PRIMARY KEY,
    "organization_id" UUID NOT NULL REFERENCES "app"."organizations" ("organization_id"),
    "name" TEXT NOT NULL,
    "description" TEXT,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "deleted_at" TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS "app"."team_memberships" (
    "team_membership_id" UUID PRIMARY KEY,
    "team_id" UUID NOT NULL REFERENCES "app"."teams" ("team_id"),
    "user_id" UUID NOT NULL REFERENCES "app"."users" ("user_id"),
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "deleted_at" TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS "app".roles (
    "role_id" UUID PRIMARY KEY,
    "name" TEXT NOT NULL,
    "description" TEXT,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "deleted_at" TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS "app".role_assignments(
    "role_assignment_id" UUID PRIMARY KEY,
    "role_id" UUID NOT NULL REFERENCES "app".roles ("role_id"),
    "assignment_audience_type" TEXT NOT NULL, -- team or user
    "assignment_audience_id" UUID NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "deleted_at" TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS "app"."resource_audience_types" (
    "name" TEXT PRIMARY KEY
);
INSERT INTO "app"."resource_audience_types" ("name") VALUES ('team'), ('user'), ('role'), ('organization');

CREATE TABLE IF NOT EXISTS "app"."resource_privileges" (
    "privilege" TEXT NOT NULL,
    "resource_type" TEXT NOT NULL REFERENCES "app".resource_types("name"),
    "description" TEXT,
    PRIMARY KEY ("privilege", "resource_type")
);

INSERT INTO "app"."resource_privileges" ("privilege", "resource_type", "description") VALUES
('member', 'organization', 'Member of the organization'),
('system_admin', 'organization', 'System admin who can add other members to the organization'),
('admin', 'page', 'Admin of a page'),
('viewer', 'page', 'Viewer of a page'),
('admin', 'dataset', 'Admin can do all actions on a dataset'),
('data_reader', 'dataset', 'Data reader can view the dataset');


CREATE TABLE IF NOT EXISTS "app"."resource_audience_policies" (
    "resource_audience_policy_id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "resource_audience_type" TEXT NOT NULL REFERENCES "app".resource_audience_types("name"),
    "resource_audience_id" UUID NOT NULL,
    "privilege" TEXT NOT NULL,
    "resource_type" TEXT NOT NULL,
    "resource_id" UUID NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "deleted_at" TIMESTAMPTZ,
    "metadata" JSONB DEFAULT '{}',
    FOREIGN KEY ("privilege", "resource_type") REFERENCES "app".resource_privileges("privilege", "resource_type")
);

-- flattened resource_audience_policies_table
-- View to flatten the resource_audience for checking users' resource access easily
CREATE OR REPLACE VIEW app.flattened_resource_audience_policies AS
    -- Get users
    SELECT 
        ra.resource_id,
        ra.resource_audience_type,
        ra.resource_audience_id AS user_id,  -- Direct user access
        ra.resource_type,
        ra.privilege,
        ra.created_at,
        ra.updated_at,
        ra.deleted_at
    FROM 
        app.resource_audience_policies ra
    WHERE 
        ra.resource_audience_type = 'user'
    
    UNION ALL

    -- Flatten team members
    SELECT 
        ra.resource_id,
        ra.resource_audience_type,
        tm.user_id AS user_id,
        ra.resource_type,
        ra.privilege,
        ra.created_at,
        ra.updated_at,
        ra.deleted_at
    FROM 
        app.resource_audience_policies ra
    JOIN 
        app.team_memberships tm 
        ON ra.resource_audience_type = 'team' 
        AND ra.resource_audience_id = tm.team_id

    UNION ALL

    -- Flatten organization members
    SELECT 
        ra.resource_id,
        ra.resource_audience_type,
        jra.resource_audience_id AS user_id,
        ra.resource_type,
        ra.privilege,
        ra.created_at,
        ra.updated_at,
        ra.deleted_at
    FROM
        app.resource_audience_policies ra
    JOIN
        app.resource_audience_policies jra
        ON ra.resource_audience_type = 'organization'
        AND jra.resource_type = 'organization'
        AND jra.resource_id = ra.resource_audience_id
        AND jra.resource_audience_type = 'user';

    -- todo handle further level of flattening when organization has team/role audiences

-- alter user trigger to add organization membership through resource_audience

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


    -- Insert an entry into the resource_audience_policies table with the right organization id and the user id
    INSERT INTO app.resource_audience_policies (resource_audience_type, resource_audience_id, privilege, resource_type, resource_id)
    VALUES ('user', NEW.user_id, 'system_admin', 'organization', new_org_id);

    RETURN NEW;

END;
$$ LANGUAGE plpgsql;


-- back fill resource_audience_polcies for existing users
INSERT INTO app.resource_audience_policies (resource_audience_type, resource_audience_id, privilege, resource_type, resource_id)
SELECT 'user', user_id, 'system_admin', 'organization', organization_id
FROM app.users
JOIN app.organization_memberships
ON app.users.user_id = app.organization_memberships.member_id;
