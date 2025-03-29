import { FC, useMemo, useRef, useState } from 'react';
import {
  useDeleteAudienceFromOrganizationAccessMutation,
  useGetAudiencesByOrganisationIdQuery,
  usePatchChangeAudienceRoleInOrganizationMutation,
} from 'apis/people';
import { useOnClickOutside } from 'hooks';
import { useAppSelector } from 'hooks/toolkit';
import RemoveFromTeamPopup from 'modules/team/components/RemoveFromTeamPopup';
import { TEAM_MEMBERS_PRIVILEGES_LIST } from 'modules/team/people.constants';
import { MembersRolePropsType, TeamMemberAccessPrivilegesType } from 'modules/team/people.types';
import { RootState } from 'store';
import { accessPermissionForPeople } from 'utils/accessPermission/accessPermission';
import { PERMISSION_MESSAGES } from 'utils/accessPermission/accessPermission.constants';
import { PERMISSION_TYPES } from 'utils/accessPermission/accessPermission.types';
import { convertEmailUsernameToName, getUserNameFromEmail } from 'utils/common';
import AsyncDropdown from 'components/asyncDropdown/AsyncDropdown';
import { toast } from 'components/common/toast/Toast';
import { TOAST_MESSAGES } from 'components/common/toast/toast.constants';

const MembersRole: FC<MembersRolePropsType> = ({ value, member = false }) => {
  const { user_id, privilege, userEmail } = value;
  const role = TEAM_MEMBERS_PRIVILEGES_LIST.find((role) => role?.value === privilege);
  const userName = useMemo(() => convertEmailUsernameToName(getUserNameFromEmail(userEmail ?? '')), [userEmail]);
  const [isOpenRemoveFromTeamPopup, setIsOpenRemoveFromTeamPopup] = useState<boolean>(false);
  const [isHoveredDropdown, setIsHoveredDropdown] = useState<boolean>(false);
  const [openChangeRoleDropdown, setOpenChangeRoleDropdown] = useState<boolean>(false);
  const [selectedRole, setSelectedRole] = useState<TeamMemberAccessPrivilegesType>(
    role as TeamMemberAccessPrivilegesType,
  );
  const [changeRole] = usePatchChangeAudienceRoleInOrganizationMutation();
  const [deleteAudience, { isLoading: isLoadingDeleteAudience }] = useDeleteAudienceFromOrganizationAccessMutation();
  const organizationId = useAppSelector((state: RootState) => state?.user?.user?.orgs?.[0]?.organization_id) ?? '';
  const { refetch: refetchAudiencesByOrganizationId } = useGetAudiencesByOrganisationIdQuery(
    { organizationId },
    { skip: !organizationId, refetchOnMountOrArgChange: false },
  );
  const dropdownRef = useRef<HTMLDivElement>(null);
  const checkPermission = accessPermissionForPeople() && member;

  const handleOpenChangeRoleDropdown = () => {
    setOpenChangeRoleDropdown(true);
  };

  const handleCloseChangeRoleDropdown = () => {
    setOpenChangeRoleDropdown(false);
  };

  const handleRoleChange = (selectedOption: TeamMemberAccessPrivilegesType) => {
    if (!checkPermission) {
      toast.error(PERMISSION_MESSAGES[PERMISSION_TYPES.ROLE_CHANGE]);

      return;
    } else {
      changeRole({
        organizationId: organizationId,
        body: {
          user_id: user_id,
          role: selectedOption?.value,
        },
      })
        .unwrap()
        .then(() => {
          setSelectedRole(selectedOption);
          setOpenChangeRoleDropdown(false);
          setIsHoveredDropdown(false);
          refetchAudiencesByOrganizationId();
          toast.success(TOAST_MESSAGES.SUCCESS_AUDIENCE_ROLE_CHANGED);
        })
        .catch((err) => {
          toast.error(err?.data?.error || TOAST_MESSAGES.FAILED_AUDIENCE_ROLE_CHANGED);
        });
    }
  };

  const handleOpenRemoveFromTeamPopup = () => {
    setIsOpenRemoveFromTeamPopup(true);
  };

  const handleCloseRemoveFromTeamPopup = () => {
    setIsOpenRemoveFromTeamPopup(false);
  };

  const handleDeleteAudience = () => {
    if (!checkPermission) {
      toast.error(PERMISSION_MESSAGES[PERMISSION_TYPES.DELETE]);

      return;
    } else {
      deleteAudience({
        organizationId: organizationId,
        body: {
          user_id: user_id,
        },
      })
        .unwrap()
        .then(() => {
          handleCloseRemoveFromTeamPopup();
          refetchAudiencesByOrganizationId();
          toast.success(`Removed ${userName} successfully`);
        })
        .catch((err) => {
          handleCloseRemoveFromTeamPopup();
          toast.error(err?.data?.error || TOAST_MESSAGES.FAILED_AUDIENCE_DELETED);
        });
    }
  };

  useOnClickOutside(dropdownRef, handleCloseChangeRoleDropdown);

  return (
    <div className='w-full h-full text-left'>
      {checkPermission ? (
        <div className='relative w-fit'>
          <AsyncDropdown
            onOpen={handleOpenChangeRoleDropdown}
            onClose={handleCloseChangeRoleDropdown}
            isOpen={openChangeRoleDropdown}
            onDelete={handleOpenRemoveFromTeamPopup}
            onChange={(role) => handleRoleChange(role as TeamMemberAccessPrivilegesType)}
            options={TEAM_MEMBERS_PRIVILEGES_LIST}
            selectedValue={selectedRole}
            defaultValue={role as TeamMemberAccessPrivilegesType}
            showDelete
            showSelectedIcon
            isHoveredDropdown={isHoveredDropdown}
            setIsHoveredDropdown={setIsHoveredDropdown}
            parentWrapperClassName='pl-2'
            wrapperClassName='w-[200px]'
            selectedOptionClassName='!bg-GRAY_100 !py-2.5'
          />
        </div>
      ) : (
        <span className='flex justify-between items-start f-12-400 text-GRAY_1000 pl-2 py-3 pr-2'>{role?.label}</span>
      )}
      <RemoveFromTeamPopup
        isOpen={isOpenRemoveFromTeamPopup}
        onClose={handleCloseRemoveFromTeamPopup}
        isLoading={isLoadingDeleteAudience}
        onDelete={handleDeleteAudience}
        feature='remove-access-from-dataset'
        warningDescription={`${userName} will be immediately removed from the organization and lose all access`}
      />
    </div>
  );
};

export default MembersRole;
