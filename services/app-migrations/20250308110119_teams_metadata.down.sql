ALTER TABLE app.teams DROP COLUMN IF EXISTS metadata;
ALTER TABLE app.teams DROP COLUMN IF EXISTS created_by;
ALTER TABLE app.teams ALTER COLUMN team_id DROP DEFAULT;

ALTER TABLE app.team_memberships DROP COLUMN IF EXISTS created_by;
ALTER TABLE app.team_memberships ALTER COLUMN team_membership_id DROP DEFAULT;


