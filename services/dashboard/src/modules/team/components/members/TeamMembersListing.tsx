import { FC, useMemo } from 'react';
import { useGetAudiencesByOrganisationIdQuery, useGetTeamsByOrganizationIdQuery } from 'apis/people';
import { useAppSelector } from 'hooks/toolkit';
import EmptyStateListing from 'modules/team/components/EmptyStateListing';
import MembersEmail from 'modules/team/components/members/MembersEmail';
import MembersName from 'modules/team/components/members/MembersName';
import MembersRole from 'modules/team/components/members/MembersRole';
import MembersTeam from 'modules/team/components/members/MembersTeam';
import SkeletonLoaderListing from 'modules/team/components/SkeletonLoaderListing';
import { TEAM_MEMBERS_LISTING_COLUMN_DEFS } from 'modules/team/people.constants';
import { TeamMembersListingPropsType } from 'modules/team/people.types';
import { RootState } from 'store';
import CommonWrapper from 'components/commonWrapper';
import { SkeletonTypes } from 'components/commonWrapper/commonWrapper.types';

const TeamMembersListing: FC<TeamMembersListingPropsType> = ({ data, isLoadingTeamMembersData }) => {
  const organizationId = useAppSelector((state: RootState) => state?.user?.user?.orgs?.[0]?.organization_id) ?? '';
  const { data: teamMembersData } = useGetAudiencesByOrganisationIdQuery(
    { organizationId },
    { skip: !organizationId, refetchOnMountOrArgChange: false },
  );
  const hasData = (teamMembersData?.length ?? 0) > 0;
  const { data: teamsData } = useGetTeamsByOrganizationIdQuery(
    { organizationId },
    { skip: !organizationId, refetchOnMountOrArgChange: false },
  );

  const userMappedTeamsMap = useMemo(() => {
    if (!teamsData || !data) return new Map();

    return new Map(
      data.map((row) => {
        const mappedTeams = teamsData
          .filter((team) => team?.team_memberships?.some((membership) => membership?.user_id === row?.user?.user_id))
          .map((team) => {
            const membership = team?.team_memberships?.find((membership) => membership?.user_id === row?.user?.user_id);

            return {
              value: team?.name,
              label: team?.name,
              valid: true,
              color: team?.metadata?.color_hex_code,
              isNew: false,
              teamId: team?.team_id,
              teamMembershipId: membership?.team_membership_id,
            };
          })
          .reverse();

        return [row?.user?.user_id, mappedTeams];
      }),
    );
  }, [teamsData, data]);

  return hasData || isLoadingTeamMembersData ? (
    <>
      <div className='grid grid-cols-4 gap-4 border-b-0.5 border-DIVIDER_GRAY'>
        {TEAM_MEMBERS_LISTING_COLUMN_DEFS.map((column, index) => (
          <div key={index} className='py-2 px-2'>
            <span className='text-left f-11-400 text-GRAY_700'>{column.headerName}</span>
          </div>
        ))}
      </div>
      <CommonWrapper
        isLoading={isLoadingTeamMembersData}
        skeletonType={SkeletonTypes.CUSTOM}
        loader={<SkeletonLoaderListing columns={4} />}
      >
        <div className='overflow-y-auto h-[calc(100vh-270px)] [&::-webkit-scrollbar]:hidden'>
          {data?.map((row, index) => {
            const userMappedTeams = userMappedTeamsMap.get(row?.user?.user_id) ?? [];

            return (
              <div key={index} className='grid grid-cols-4 gap-4 border-b-0.5 border-DIVIDER_GRAY'>
                <MembersName value={row?.user?.email} member />
                <MembersEmail value={row?.user?.email} />
                <MembersRole
                  value={{ user_id: row?.user?.user_id, privilege: row?.privilege, userEmail: row?.user?.email }}
                  member
                />
                <MembersTeam
                  organizationId={organizationId}
                  teamsData={teamsData ?? []}
                  userId={row?.user?.user_id}
                  userMappedTeams={userMappedTeams}
                />
              </div>
            );
          })}
        </div>
      </CommonWrapper>
    </>
  ) : (
    <EmptyStateListing title='No team members were added' />
  );
};

export default TeamMembersListing;
