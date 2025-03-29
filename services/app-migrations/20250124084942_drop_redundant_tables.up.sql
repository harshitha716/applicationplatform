DELETE FROM app.resource_types WHERE name = 'organization_invitation';
DELETE FROM app.resource_types WHERE name = 'workspace';
DELETE FROM app.resource_types WHERE name = 'user';
DELETE FROM app.resource_types WHERE name = 'sheet';

DROP TABLE IF EXISTS app.organization_memberships;

DROP TABLE IF EXISTS app.organization_membership_invitations;

DROP TABLE IF EXISTS app.resource_access;

DROP TABLE IF EXISTS app.resource_access_audience_types;

DROP TABLE IF EXISTS app.resource_access_level_types;
