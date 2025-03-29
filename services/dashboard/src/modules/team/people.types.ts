import { AudiencesByOrganisationIdResponse, InvitedAudiencesByOrganisationIdResponse } from 'types/api/people.types';
import { defaultFnType } from 'types/commonTypes';

export enum TEAM_TABS_TYPES {
  TEAM_MEMBERS = 'team_members',
  INVITED_MEMBERS = 'invited_members',
}

export const TeamTabsList = [
  { label: 'Team members', value: TEAM_TABS_TYPES.TEAM_MEMBERS },
  { label: 'Invited', value: TEAM_TABS_TYPES.INVITED_MEMBERS },
];

export type InviteMembersPopupPropsType = {
  isOpen: boolean;
  onClose?: defaultFnType;
  teamMembersData: AudiencesByOrganisationIdResponse[];
};

export type TeamMembersPrivilegeType = {
  label: string;
  value: TEAM_MEMBERS_PRIVILEGES;
};

export type RemoveFromTeamPopupPropsType = {
  name?: string;
  isOpen: boolean;
  onClose?: defaultFnType;
  onDelete?: defaultFnType;
  isLoading?: boolean;
  feature?: string;
  warningDescription: string;
};

export type MembersEmailPropsType = {
  value: string;
};

export type InvitedMembersListingPropsType = {
  isLoadingInvitedTeamMembersData: boolean;
  data: InvitedAudiencesByOrganisationIdResponse[];
};

export type MembersNamePropsType = {
  value: string;
  member?: boolean;
};

export type MembersRolePropsType = {
  value: { user_id: string; privilege: string; userEmail?: string };
  member?: boolean;
};

export enum TEAM_MEMBERS_PRIVILEGES {
  SYSTEM_ADMIN = 'system_admin',
  MEMBER = 'member',
  REMOVE = 'remove',
}

export type TeamMemberAccessPrivilegesType = {
  label: string;
  value: TEAM_MEMBERS_PRIVILEGES;
};

export type EmptyStateListingPropsType = {
  title?: string;
};

export type TeamMembersListingPropsType = {
  isLoadingTeamMembersData: boolean;
  data: AudiencesByOrganisationIdResponse[];
};

export type CustomTeamsDropdownPropsType = {
  search?: string;
  optionRefs: React.MutableRefObject<(HTMLDivElement | null)[]>;
  optionList: { value: string; label: string; color?: string }[];
  isLoadingOptionsList: boolean;
  hoveredOptionIndex: number;
  setHoveredOptionIndex: (index: number) => void;
  onSelectOption: (option: { value: string; label: string; color?: string }) => void;
  transformLabel?: (label: string) => string;
  handleKeyDown: (e: React.KeyboardEvent<HTMLDivElement>) => void;
  randomColor: string | null;
  onKeyDown: (e: React.KeyboardEvent<HTMLDivElement>) => void;
  onCloseDropdown: defaultFnType;
};

export type SelectedItemsType = {
  value: string;
  label: string;
  valid: boolean;
  role?: string;
  color: string;
  resource_audience_type?: string;
  resource_audience_id?: string;
  validationMessage?: string;
};

export type MembersTeamPropsType = {
  organizationId: string;
  userId: string;
  teamsData: {
    team_id: string;
    name: string;
    description: string;
    metadata: { color_hex_code: string };
  }[];
  userMappedTeams: {
    value: string;
    label: string;
    valid: boolean;
    color?: string;
    isNew?: boolean;
    teamId?: string;
    teamMembershipId?: string;
  }[];
};

export type PostTeamsByOrganizationIdPayload = {
  name: string;
  description: string;
  color_hex_code: string;
};

export type PostAddTeamToAudiencePayload = {
  user_id: string;
  team_id: string;
};
