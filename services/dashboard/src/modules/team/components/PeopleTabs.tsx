import { FC, useState } from 'react';
import InvitedMembersListing from 'modules/team/components/members/InvitedMembersListing';
import TeamMembersListing from 'modules/team/components/members/TeamMembersListing';
import { TEAM_TABS_TYPES, TeamTabsList } from 'modules/team/people.types';
import { AudiencesByOrganisationIdResponse, InvitedAudiencesByOrganisationIdResponse } from 'types/api/people.types';
import { MenuItem, TAB_TYPES } from 'types/common/components';
import { checkIfCurrentUserIsMember } from 'utils/accessPermission/accessPermission.utils';
import { Tabs } from 'components/common/tabs/Tabs';

type PeopleTabsPropsType = {
  filteredTeamMembers: AudiencesByOrganisationIdResponse[];
  isLoadingTeamMembersData: boolean;
  filteredInvitedMembers: InvitedAudiencesByOrganisationIdResponse[];
  isLoadingInvitedTeamMembersData: boolean;
};

const PeopleTabs: FC<PeopleTabsPropsType> = ({
  filteredTeamMembers,
  isLoadingTeamMembersData,
  filteredInvitedMembers,
  isLoadingInvitedTeamMembersData,
}) => {
  const [selectedTab, setSelectedTab] = useState<TEAM_TABS_TYPES>(TEAM_TABS_TYPES.TEAM_MEMBERS);
  const checkIfSystemAdmin = !checkIfCurrentUserIsMember();

  const handleTabSelect = (item?: MenuItem) => {
    setSelectedTab(item?.value as TEAM_TABS_TYPES);
  };

  const renderTeamListing = () => {
    switch (selectedTab) {
      case TEAM_TABS_TYPES.TEAM_MEMBERS:
        return <TeamMembersListing data={filteredTeamMembers} isLoadingTeamMembersData={isLoadingTeamMembersData} />;
      case TEAM_TABS_TYPES.INVITED_MEMBERS:
        return (
          <InvitedMembersListing
            data={filteredInvitedMembers}
            isLoadingInvitedTeamMembersData={isLoadingInvitedTeamMembersData}
          />
        );
      default:
        return null;
    }
  };

  return (
    <>
      <div className='my-4'>
        {checkIfSystemAdmin && (
          <Tabs list={TeamTabsList} id='team-tabs' type={TAB_TYPES.UNDERLINE} onSelect={handleTabSelect} />
        )}
      </div>
      {renderTeamListing()}
    </>
  );
};

export default PeopleTabs;
