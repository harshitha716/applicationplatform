-- table to store organization invitations
CREATE TABLE IF NOT EXISTS app.organization_invitations (
    organization_invitation_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES app.organizations(organization_id),
    email TEXT NOT NULL,
    privilege TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    invited_by UUID NOT NULL REFERENCES app.users(user_id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    email_retry_count INT NOT NULL DEFAULT 0,
    email_sent_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS app.organization_invitation_statuses (
    organization_invitation_id UUID NOT NULL REFERENCES app.organization_invitations(organization_invitation_id),
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (organization_invitation_id)
);

-- table to store requests by users to join an organization
CREATE TABLE IF NOT EXISTS app.organization_membership_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES app.organizations(organization_id),
    user_id UUID NOT NULL REFERENCES app.users(user_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    status TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);