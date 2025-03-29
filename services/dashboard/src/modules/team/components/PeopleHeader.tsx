import React, { FC, useState } from 'react';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import InviteMembersPopup from 'modules/team/InviteMembersPopup';
import { AudiencesByOrganisationIdResponse } from 'types/api/people.types';
import { SIZE_TYPES } from 'types/common/components';
import { BUTTON_TYPES } from 'types/components/button.type';
import { PERMISSION_ROLES } from 'utils/accessPermission/accessPermission.types';
import { getUserPrivilege } from 'utils/accessPermission/accessPermission.utils';
import { Button } from 'components/common/button/Button';
import Input from 'components/common/input';

type PeopleHeaderPropsType = {
  search: string;
  setSearch: (value: string) => void;
  teamMembersData: AudiencesByOrganisationIdResponse[];
};

const PeopleHeader: FC<PeopleHeaderPropsType> = ({ search, setSearch, teamMembersData }) => {
  const inputRef = React.useRef<HTMLInputElement>(null);
  const [isInviteMembersPopupOpen, setIsInviteMembersPopupOpen] = useState(false);
  const userPrivilege = getUserPrivilege();
  const checkIfMember = userPrivilege === PERMISSION_ROLES.MEMBER;

  const handleOpenInviteMembersPopup = () => {
    setIsInviteMembersPopupOpen(true);
  };
  const handleCloseInviteMembersPopup = () => {
    setIsInviteMembersPopupOpen(false);
  };

  return (
    <>
      <div className='f-20-600 text-GRAY_1000'>Team</div>
      <div className='flex justify-between items-center w-full mt-5'>
        <Input
          placeholder='Search team members'
          className='w-80'
          inputRef={inputRef}
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          leadingIconProps={{
            id: 'search-sm',
            iconCategory: ICON_SPRITE_TYPES.GENERAL,
            className: 'text-GRAY_700',
          }}
          size={SIZE_TYPES.SMALL}
        />

        <Button
          type={BUTTON_TYPES.PRIMARY}
          id='invite-user-btn'
          size={SIZE_TYPES.SMALL}
          onClick={handleOpenInviteMembersPopup}
          disabled={checkIfMember}
        >
          Invite members
        </Button>
        <InviteMembersPopup
          isOpen={isInviteMembersPopupOpen}
          onClose={handleCloseInviteMembersPopup}
          teamMembersData={teamMembersData}
        />
      </div>
    </>
  );
};

export default PeopleHeader;
