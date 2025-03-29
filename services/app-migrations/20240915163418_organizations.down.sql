DROP TRIGGER IF EXISTS insert_user_into_organizations
ON app.users CASCADE;

DROP TABLE IF EXISTS app.organization_membership_invitations;

DROP TABLE IF EXISTS app.organization_memberships;

DROP TABLE IF EXISTS  app.organizations;

