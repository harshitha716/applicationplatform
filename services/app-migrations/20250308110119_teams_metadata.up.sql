ALTER TABLE app.teams ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}';
ALTER TABLE app.teams ADD COLUMN IF NOT EXISTS created_by uuid NOT NULL REFERENCES app.users(user_id);
ALTER TABLE app.teams ALTER COLUMN team_id SET DEFAULT gen_random_uuid();
ALTER TABLE app.teams DROP CONSTRAINT IF EXISTS unique_organization_id_name;
ALTER TABLE app.teams ADD CONSTRAINT unique_organization_id_name UNIQUE (organization_id, name);

ALTER TABLE app.team_memberships ADD COLUMN IF NOT EXISTS created_by uuid NOT NULL REFERENCES app.users(user_id);
ALTER TABLE app.team_memberships ALTER COLUMN team_membership_id SET DEFAULT gen_random_uuid();
ALTER TABLE app.team_memberships DROP CONSTRAINT IF EXISTS unique_team_id_user_id;
ALTER TABLE app.team_memberships ADD CONSTRAINT unique_team_id_user_id UNIQUE (team_id, user_id);

