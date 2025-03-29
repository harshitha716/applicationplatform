ALTER TABLE app.teams DROP CONSTRAINT IF EXISTS unique_organization_id_name;
ALTER TABLE app.team_memberships DROP CONSTRAINT IF EXISTS unique_team_id_user_id;