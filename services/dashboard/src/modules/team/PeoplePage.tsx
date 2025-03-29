import { useEffect, useMemo, useState } from 'react';
import { useGetAudiencesByOrganisationIdQuery, useGetInvitedAudiencesByOrganisationIdQuery } from 'apis/people';
import { debounce } from 'hooks';
import { useAppSelector } from 'hooks/toolkit';
import PeopleHeader from 'modules/team/components/PeopleHeader';
import PeopleTabs from 'modules/team/components/PeopleTabs';
import { RootState } from 'store';
import { convertEmailUsernameToName, getUserNameFromEmail } from 'utils/common';

const PeoplePage = () => {
  const organizationId = useAppSelector((state: RootState) => state?.user?.user?.orgs?.[0]?.organization_id) ?? '';
  const { data: teamMembersData, isLoading: isLoadingTeamMembersData } = useGetAudiencesByOrganisationIdQuery(
    { organizationId },
    { skip: !organizationId, refetchOnMountOrArgChange: false },
  );
  const { data: invitedTeamMembersData, isLoading: isLoadingInvitedTeamMembersData } =
    useGetInvitedAudiencesByOrganisationIdQuery({ organizationId }, { skip: !organizationId });
  const [search, setSearch] = useState('');
  const [filteredTeamMembers, setFilteredTeamMembers] = useState(teamMembersData);
  const [filteredInvitedMembers, setFilteredInvitedMembers] = useState(invitedTeamMembersData);

  const debouncedFilterTeamMembers = useMemo(
    () =>
      debounce((searchValue: string) => {
        setFilteredTeamMembers(
          teamMembersData?.filter((member) =>
            convertEmailUsernameToName(getUserNameFromEmail(member?.user?.email))
              .toLowerCase()
              .startsWith(searchValue.toLowerCase()),
          ),
        );
      }, 300),
    [teamMembersData],
  );

  const debouncedFilterInvitedMembers = useMemo(
    () =>
      debounce((searchValue: string) => {
        setFilteredInvitedMembers(
          invitedTeamMembersData?.filter((member) =>
            convertEmailUsernameToName(getUserNameFromEmail(member?.email))
              .toLowerCase()
              .startsWith(searchValue.toLowerCase()),
          ),
        );
      }, 300),
    [invitedTeamMembersData],
  );

  useEffect(() => {
    debouncedFilterTeamMembers(search);
    debouncedFilterInvitedMembers(search);
  }, [search, debouncedFilterTeamMembers, debouncedFilterInvitedMembers]);

  return (
    <div className='p-10 w-full h-full'>
      <PeopleHeader search={search} setSearch={setSearch} teamMembersData={teamMembersData ?? []} />
      <PeopleTabs
        filteredTeamMembers={filteredTeamMembers ?? []}
        isLoadingTeamMembersData={isLoadingTeamMembersData}
        filteredInvitedMembers={filteredInvitedMembers ?? []}
        isLoadingInvitedTeamMembersData={isLoadingInvitedTeamMembersData}
      />
    </div>
  );
};

export default PeoplePage;
