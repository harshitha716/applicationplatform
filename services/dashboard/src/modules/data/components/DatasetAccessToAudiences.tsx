import React, { FC, useRef, useState } from 'react';
import {
  useDeleteAudienceFromDatasetAccessMutation,
  useGetAudiencesByDatasetIdQuery,
  usePatchChangeAudienceRoleInDatasetMutation,
} from 'apis/dataset';
import { COLORS } from 'constants/colors';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { useOnClickOutside } from 'hooks';
import { CHANGE_ACCESS_PRIVILEGES_LIST, DATASET_ACCESS_PRIVILEGES_LIST } from 'modules/data/data.constants';
import { DatasetAccessPrivilegesType, DatasetAccessToAudiencesPropsType } from 'modules/data/data.types';
import RemoveFromTeamPopup from 'modules/team/components/RemoveFromTeamPopup';
import { ResourceAudienceType } from 'types/api/auth.types';
import { accessPermissionForDataset } from 'utils/accessPermission/accessPermission';
import { PERMISSION_MESSAGES } from 'utils/accessPermission/accessPermission.constants';
import { PERMISSION_TYPES } from 'utils/accessPermission/accessPermission.types';
import { checkIfCurrentUser } from 'utils/accessPermission/accessPermission.utils';
import { cn, convertEmailUsernameToName, getUserNameFromEmail } from 'utils/common';
import AsyncDropdown from 'components/asyncDropdown/AsyncDropdown';
import Avatar from 'components/common/avatar';
import { toast } from 'components/common/toast/Toast';
import { TOAST_MESSAGES } from 'components/common/toast/toast.constants';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

const DatasetAccessToAudiences: FC<DatasetAccessToAudiencesPropsType> = ({
  resource_type,
  privilege,
  datasetId,
  resource_audience_id,
  resource_audience_type,
  user,
  userPrivilege,
  orgName,
  customerName,
  teamInfo,
}) => {
  const dropdownRef = useRef<HTMLDivElement>(null);
  const role = DATASET_ACCESS_PRIVILEGES_LIST.find((role) => role.value === privilege);
  const [isOpenRemoveFromTeamPopup, setIsOpenRemoveFromTeamPopup] = useState<boolean>(false);
  const [isHoveredDropdown, setIsHoveredDropdown] = useState<boolean>(false);
  const [openChangeRoleDropdown, setOpenChangeRoleDropdown] = useState<boolean>(false);
  const [selectedRole, setSelectedRole] = useState<DatasetAccessPrivilegesType>(role as DatasetAccessPrivilegesType);
  const { refetch: refetchAudiencesByDatasetId } = useGetAudiencesByDatasetIdQuery({ datasetId }, { skip: !datasetId });
  const [changeRole] = usePatchChangeAudienceRoleInDatasetMutation();
  const [deleteAudience] = useDeleteAudienceFromDatasetAccessMutation();
  const checkIfUser = checkIfCurrentUser(user?.email ?? '');
  const checkIfResourceTypeOrg = resource_audience_type === ResourceAudienceType.ORGANIZATION;
  const checkIfResourceTypeTeam = resource_audience_type === ResourceAudienceType.TEAM;
  const userName = checkIfResourceTypeOrg
    ? orgName
    : checkIfResourceTypeTeam
      ? teamInfo?.name
      : convertEmailUsernameToName(getUserNameFromEmail(user?.email || resource_audience_type)) || 'Unknown';
  const customAvatarWord = (checkIfResourceTypeOrg ? customerName : userName) || 'Unknown';
  const checkPermission = accessPermissionForDataset(userPrivilege);
  const showRoleChangeDropdown = checkPermission && !checkIfResourceTypeOrg;

  const handleOpenChangeRoleDropdown = () => {
    setOpenChangeRoleDropdown(true);
  };

  const handleCloseChangeRoleDropdown = () => {
    setOpenChangeRoleDropdown(false);
  };

  const handleRoleChange = (selectedOption: DatasetAccessPrivilegesType) => {
    if (!checkPermission) {
      toast.error(PERMISSION_MESSAGES[PERMISSION_TYPES.ROLE_CHANGE]);

      return;
    } else {
      changeRole({
        datasetId: datasetId,
        body: {
          audience_id: resource_audience_id,
          role: selectedOption?.value,
        },
      })
        .unwrap()
        .then(() => {
          setSelectedRole(selectedOption);
          setOpenChangeRoleDropdown(false);
          setIsHoveredDropdown(false);
          refetchAudiencesByDatasetId();
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
        datasetId: datasetId,
        body: {
          audience_id: resource_audience_id,
        },
      })
        .unwrap()
        .then(() => {
          handleCloseRemoveFromTeamPopup();
          refetchAudiencesByDatasetId();
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
    <>
      <div className='f-12-400 pl-2 bg-white flex justify-between items-center'>
        <div className='flex items-center justify-start'>
          <div className='flex items-start justify-start gap-x-1 w-[140px]'>
            <div className='flex items-center gap-1'>
              {checkIfResourceTypeTeam ? (
                <div>
                  <SvgSpriteLoader id='users-02' width={14} height={14} color={COLORS.GRAY_1000} className='mr-0.5' />
                </div>
              ) : (
                <div className='w-fit'>
                  <Avatar
                    name={customAvatarWord}
                    backgroundColor={COLORS.GRAY_1000}
                    className='w-4 h-4 rounded-full text-white f-8-400 flex items-center justify-center'
                  />
                </div>
              )}
              <div
                className={cn(
                  'flex justify-center items-center gap-1',
                  checkIfResourceTypeTeam && 'px-1.5 py-0.5 rounded',
                )}
                style={{
                  backgroundColor: checkIfResourceTypeTeam ? teamInfo?.color : 'transparent',
                }}
              >
                {userName}
                <span className='f-12-400 text-GRAY_700'>{checkIfUser && '(You)'}</span>
              </div>
            </div>
          </div>
          <div className='hidden text-wrap flex-wrap break-words whitespace-normal items-center justify-start gap-1 w-[100px]'>
            {checkPermission && (
              <>
                <SvgSpriteLoader
                  id='coins-stacked-04'
                  iconCategory={ICON_SPRITE_TYPES.FINANCE_AND_ECOMMERCE}
                  width={12}
                  height={12}
                  color={COLORS.GRAY_1000}
                  className='mr-1'
                />
                {resource_type}
              </>
            )}
          </div>
        </div>
        {showRoleChangeDropdown ? (
          <AsyncDropdown
            onOpen={handleOpenChangeRoleDropdown}
            onClose={handleCloseChangeRoleDropdown}
            isOpen={openChangeRoleDropdown}
            onDelete={handleOpenRemoveFromTeamPopup}
            onChange={(role) => handleRoleChange(role as DatasetAccessPrivilegesType)}
            options={CHANGE_ACCESS_PRIVILEGES_LIST}
            selectedValue={selectedRole}
            defaultValue={role as DatasetAccessPrivilegesType}
            showDelete
            showSelectedIcon
            isHoveredDropdown={isHoveredDropdown}
            setIsHoveredDropdown={setIsHoveredDropdown}
            isOverflowStyle
          />
        ) : (
          <span
            className={cn(
              'flex justify-between items-start f-12-400 text-GRAY_1000 pl-4 py-3 pr-2',
              !showRoleChangeDropdown && 'pr-4 text-GRAY_600',
            )}
          >
            {role?.label}
          </span>
        )}
      </div>
      <RemoveFromTeamPopup
        isOpen={isOpenRemoveFromTeamPopup}
        onClose={handleCloseRemoveFromTeamPopup}
        onDelete={handleDeleteAudience}
        feature='remove-access-from-dataset'
        warningDescription={`${userName} will be immediately removed from ${resource_type} and lose all access`}
      />
    </>
  );
};

export default DatasetAccessToAudiences;
