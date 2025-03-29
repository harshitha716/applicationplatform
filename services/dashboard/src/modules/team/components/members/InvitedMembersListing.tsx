import { FC } from 'react';
import { useGetInvitedAudiencesByOrganisationIdQuery } from 'apis/people';
import { useAppSelector } from 'hooks/toolkit';
import EmptyStateListing from 'modules/team/components/EmptyStateListing';
import MembersEmail from 'modules/team/components/members/MembersEmail';
import MembersName from 'modules/team/components/members/MembersName';
import MembersRole from 'modules/team/components/members/MembersRole';
import SkeletonLoaderListing from 'modules/team/components/SkeletonLoaderListing';
import { INVITE_TEAM_MEMBERS_LISTING_COLUMN_DEFS } from 'modules/team/people.constants';
import { InvitedMembersListingPropsType } from 'modules/team/people.types';
import { RootState } from 'store';
import CommonWrapper from 'components/commonWrapper';
import { SkeletonTypes } from 'components/commonWrapper/commonWrapper.types';

const InvitedMembersListing: FC<InvitedMembersListingPropsType> = ({ data, isLoadingInvitedTeamMembersData }) => {
  const organizationId = useAppSelector((state: RootState) => state?.user?.user?.orgs?.[0]?.organization_id) ?? '';
  const { data: invitedTeamMembersData } = useGetInvitedAudiencesByOrganisationIdQuery(
    { organizationId },
    { skip: !organizationId, refetchOnMountOrArgChange: false },
  );
  const hasData = (invitedTeamMembersData?.length ?? 0) > 0;

  return hasData || isLoadingInvitedTeamMembersData ? (
    <>
      <div className='grid grid-cols-3 gap-4 border-b-0.5 border-DIVIDER_GRAY'>
        {INVITE_TEAM_MEMBERS_LISTING_COLUMN_DEFS.map((column, index) => (
          <div key={index} className='py-2 px-2'>
            <span className='text-left f-11-400 text-GRAY_700'>{column.headerName}</span>
          </div>
        ))}
      </div>
      <CommonWrapper
        isLoading={isLoadingInvitedTeamMembersData}
        skeletonType={SkeletonTypes.CUSTOM}
        loader={<SkeletonLoaderListing />}
      >
        <div className='overflow-y-auto h-[calc(100vh-270px)] [&::-webkit-scrollbar]:hidden'>
          {data.map((row, index) => (
            <div key={index} className='grid grid-cols-3 gap-4 border-b-0.5 border-DIVIDER_GRAY'>
              <MembersName value={row?.email} />
              <MembersEmail value={row?.email} />
              <MembersRole value={{ user_id: '', privilege: row?.privilege }} />
            </div>
          ))}
        </div>
      </CommonWrapper>
    </>
  ) : (
    <EmptyStateListing title='No pending invitations' />
  );
};

export default InvitedMembersListing;
