import { PostAddTeamToAudiencePayload, PostTeamsByOrganizationIdPayload } from 'modules/team/people.types';

export type AudiencesByOrganisationIdRequest = {
  organizationId: string;
};

export type AudiencesByOrganisationIdResponse = {
  user: {
    email: string;
    user_id: string;
    name: string;
  };
  privilege: string;
  resource_audience_type: string;
  resource_audience_id: string;
};

export type InvitedAudiencesByOrganisationIdResponse = {
  name: string;
  email: string;
  privilege: string;
};

export type PostAudiencesInviteData = {
  invitations: {
    email: string;
    role: string;
  }[];
};

export type PatchChangeAudienceRoleInOrganizationType = {
  organizationId: string;
  body: { user_id: string; role: string };
};

export type DeleteAudienceFromOrganizationAccessType = { organizationId: string; body: { user_id: string } };

export type GetMembershipRequestsByOrganizationIdRequest = { organizationId: string };

export type GetMembershipRequestsByOrganizationIdResponse = {
  id: string;
  organization_id: string;
  user_id: string;
  created_at: string;
  updated_at: string;
  deleted_at: string;
  status: string;
}[];

export type GetTeamsByOrganizationIdResponseType = {
  team_id: string;
  organization_id: string;
  name: string;
  description: string;
  metadata: {
    color_hex_code: string;
  };
  team_memberships: [
    {
      team_membership_id: string;
      team_id: string;
      user_id: string;
    },
  ];
}[];
export type GetTeamsByOrganizationIdRequestType = {
  organizationId: string;
};

export type PostTeamsByOrganizationIdRequestType = {
  organizationId: string;
  payload: PostTeamsByOrganizationIdPayload;
};

export type PostTeamsByOrganizationIdResponseType = {
  team_id: string;
};

export type PostAddTeamToAudienceRequestType = {
  organizationId: string;
  teamId: string;
  payload: PostAddTeamToAudiencePayload;
};

export type RemoveTeamFromAudienceRequestType = {
  organizationId: string;
  teamId: string;
  payload: { team_id: string; team_membership_id: string };
};
