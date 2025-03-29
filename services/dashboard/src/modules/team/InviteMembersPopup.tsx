import React, { FC, useCallback, useEffect, useMemo, useState } from 'react';
import {
  useGetInvitedAudiencesByOrganisationIdQuery,
  usePostInviteAudiencesByOrganisationIdMutation,
} from 'apis/people';
import { COLORS } from 'constants/colors';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { useAppSelector } from 'hooks/toolkit';
import { TEAM_MEMBERS_PRIVILEGES_LIST } from 'modules/team/people.constants';
import { InviteMembersPopupPropsType, TEAM_MEMBERS_PRIVILEGES } from 'modules/team/people.types';
import { RootState } from 'store';
import { SIZE_TYPES } from 'types/common/components';
import { BUTTON_TYPES } from 'types/components/button.type';
import { accessPermissionForPeople } from 'utils/accessPermission/accessPermission';
import { PERMISSION_MESSAGES, VALIDATION_ERROR_MESSAGES } from 'utils/accessPermission/accessPermission.constants';
import { PERMISSION_ROLES, PERMISSION_TYPES } from 'utils/accessPermission/accessPermission.types';
import { getUserEmail, getUserPrivilege } from 'utils/accessPermission/accessPermission.utils';
import { validateEmail } from 'utils/common';
import { Button } from 'components/common/button/Button';
import Popup from 'components/common/popup/Popup';
import { toast } from 'components/common/toast/Toast';
import { TOAST_MESSAGES } from 'components/common/toast/toast.constants';
import MultiSelectInput from 'components/multiSelectInput/MultiSelectInput';
import { ArrayListOption } from 'components/multiSelectInput/multiSelectInput.types';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

const InviteMembersPopup: FC<InviteMembersPopupPropsType> = ({ isOpen, onClose, teamMembersData }) => {
  const user_email = getUserEmail();
  const userPrivilege = getUserPrivilege();
  const checkPermission = accessPermissionForPeople();
  const placeholderText = 'Share with people and teams';
  const organizationId = useAppSelector((state: RootState) => state?.user?.user?.orgs?.[0]?.organization_id) ?? '';
  const [validationErrorText, setValidationErrorText] = useState<string>('');
  const [showValidationError, setShowValidationError] = useState<boolean>(true);
  const [multiSelectInstances, setMultiSelectInstances] = useState<number[]>([0]);
  const [searchValues, setSearchValues] = useState<{ [key: number]: string }>({});
  const [pendingEntryByInstance, setPendingEntryByInstance] = useState<{ [key: number]: string }>({});
  const [selectedItemsByInstance, setSelectedItemsByInstance] = useState<{ [key: number]: ArrayListOption[] }>({});
  const [selectedRoleByInstance, setSelectedRoleByInstance] = useState<{ [key: number]: TEAM_MEMBERS_PRIVILEGES }>({});
  const hasEmptySearchValue = useMemo(() => Object.values(searchValues).some((value) => value === ''), [searchValues]);
  const hasNonEmptySelectedItems = useMemo(
    () => Object.values(selectedItemsByInstance)?.some((item) => item?.length !== 0),
    [selectedItemsByInstance],
  );
  const disableAddBtn = useMemo(
    () => multiSelectInstances?.length === TEAM_MEMBERS_PRIVILEGES_LIST.length,
    [multiSelectInstances],
  );
  const isInvitable = useMemo(() => {
    if (showValidationError || userPrivilege === PERMISSION_ROLES.MEMBER) return false;
    const hasValidSearch = !hasEmptySearchValue || hasNonEmptySelectedItems;

    return hasValidSearch && multiSelectInstances?.length > 0;
  }, [showValidationError, userPrivilege, hasEmptySearchValue, hasNonEmptySelectedItems, multiSelectInstances]);
  const [postInviteAudiences, { isLoading: postInviteAudiencesIsLoading }] =
    usePostInviteAudiencesByOrganisationIdMutation();
  const { data: invitedTeamMembersData, refetch: refetchAudiencesByOrganizationId } =
    useGetInvitedAudiencesByOrganisationIdQuery(
      { organizationId },
      { skip: !organizationId, refetchOnMountOrArgChange: false },
    );

  const handleSearchChange = (id: number, value: string) => {
    setSearchValues((prev) => ({
      ...prev,
      [id]: value,
    }));

    setPendingEntryByInstance((prev) => (prev[id] === value ? prev : { ...prev, [id]: value }));
  };

  const handleCloseInviteMembersPopup = () => {
    onClose?.();
    setShowValidationError(false);
    setSelectedItemsByInstance({});
    setSearchValues({});
    setMultiSelectInstances([0]);
    setPendingEntryByInstance({});
    setSelectedRoleByInstance({});
  };

  const handleInviteMembers = () => {
    if (!checkPermission) {
      return toast.error(PERMISSION_MESSAGES[PERMISSION_TYPES.INVITE]);
    }

    const finalSelectedItemsByInstance = { ...selectedItemsByInstance };

    Object.entries(pendingEntryByInstance)?.forEach(([idStr, value]) => {
      if (!value?.trim()) return;
      const id = Number(idStr);
      const { isValid, message } = validateAndGetUserDetails(value);
      const role = selectedRoleByInstance[id] || TEAM_MEMBERS_PRIVILEGES_LIST[0].value;

      if (!finalSelectedItemsByInstance[id]) {
        finalSelectedItemsByInstance[id] = [];
      }

      finalSelectedItemsByInstance[id].push({
        value,
        label: value,
        valid: isValid,
        role,
        color: isValid ? COLORS.WHITE : COLORS.RED_100,
        validationMessage: message,
      });
    });

    const invitations = Object.values(finalSelectedItemsByInstance)
      .flat()
      .filter((item) => item?.value)
      .map(({ value, role }) => ({ email: value, role: role ?? TEAM_MEMBERS_PRIVILEGES_LIST[0].value }));

    postInviteAudiences({ organizationId, body: { invitations } })
      .unwrap()
      .then(() => {
        refetchAudiencesByOrganizationId();
        toast.success(TOAST_MESSAGES.SUCCESS_AUDIENCE_INVITED);
        handleCloseInviteMembersPopup();
      })
      .catch((err) => {
        toast.error(err?.data?.error || TOAST_MESSAGES.FAILED_AUDIENCE_INVITED);
      });
  };

  const validateAndGetUserDetails = useCallback(
    (value: string) => {
      const isValid = validateEmail(value);

      if (!isValid) {
        return { isValid: false, message: VALIDATION_ERROR_MESSAGES.INVALID_EMAIL };
      }

      const isAlreadyInvited = invitedTeamMembersData?.some((item) => item?.email === value);

      if (isAlreadyInvited) {
        return { isValid: false, message: VALIDATION_ERROR_MESSAGES.USER_ALREADY_HAS_ACCESS };
      }

      const isOrgMember = teamMembersData?.some((item) => item?.user?.email === value);

      if (isOrgMember) {
        return { isValid: false, message: VALIDATION_ERROR_MESSAGES.USER_ALREADY_IN_ORG };
      }

      if (value === user_email) {
        return { isValid: false, message: VALIDATION_ERROR_MESSAGES.CANNOT_ADD_SELF };
      }

      return { isValid };
    },
    [invitedTeamMembersData, teamMembersData, user_email],
  );

  const handleValidateAndAdd = (id: number, { value }: { value: string; label: string; color?: string }) => {
    const instanceRole = selectedRoleByInstance[id] || TEAM_MEMBERS_PRIVILEGES_LIST[0].value;
    const existingEmails = new Set(selectedItemsByInstance[id]?.map((item) => item.value) || []);
    const uniqueEntries = new Set<string>();
    const splitUsingRegex = /[, ]+/;

    const validatedEntries = value
      .split(splitUsingRegex)
      .map((email) => email?.trim())
      .filter(Boolean)
      .filter((email) => !existingEmails?.has(email) && !uniqueEntries?.has(email))
      .map((email) => {
        uniqueEntries.add(email);
        const { isValid, message } = validateAndGetUserDetails(email);

        return {
          value: email,
          label: email,
          valid: isValid,
          role: instanceRole,
          color: isValid ? COLORS.WHITE : COLORS.RED_100,
          validationMessage: message,
        };
      });

    if (!validatedEntries?.length) return;

    setSelectedItemsByInstance((prev) => ({
      ...prev,
      [id]: [...(prev[id] || []), ...validatedEntries],
    }));

    const firstInvalidEntry = validatedEntries?.find((item) => !item?.valid);

    setShowValidationError(!!firstInvalidEntry);
    setValidationErrorText(firstInvalidEntry?.validationMessage ?? '');
  };

  const updateSelectedRoles = useCallback(() => {
    setSelectedRoleByInstance((prev) => {
      const updatedRoles = { ...prev };

      multiSelectInstances?.forEach((id) => {
        if (!(id in updatedRoles)) {
          updatedRoles[id] = TEAM_MEMBERS_PRIVILEGES_LIST[0].value;
        }
      });

      return updatedRoles;
    });
  }, [multiSelectInstances]);

  const validateSearchAndSelectedItems = (
    searchValues: { [key: number]: string },
    selectedItemsByInstance: { [key: number]: ArrayListOption[] },
    setShowValidationError: React.Dispatch<React.SetStateAction<boolean>>,
    setValidationErrorText: React.Dispatch<React.SetStateAction<string>>,
  ) => {
    let hasInvalidEntry = false;
    let firstErrorMessage = '';

    Object.entries(searchValues).forEach(([idStr, search]) => {
      const id = Number(idStr);

      if (search?.trim() !== '' && !selectedItemsByInstance[id]?.some((item) => item?.value === search)) {
        const { isValid, message } = validateAndGetUserDetails(search);

        if (!isValid) {
          hasInvalidEntry = true;
          firstErrorMessage = message || '';
        }
      }
    });

    if (!hasInvalidEntry) {
      Object.values(selectedItemsByInstance)?.forEach((items) => {
        if (!Array.isArray(items)) return;

        const invalidItem = items?.find((item) => !item?.valid);

        if (invalidItem) {
          hasInvalidEntry = true;
          firstErrorMessage = invalidItem.validationMessage || 'Invalid item found';
        }
      });
    }

    setShowValidationError(hasInvalidEntry);
    setValidationErrorText(firstErrorMessage);
  };

  useEffect(() => {
    updateSelectedRoles();
  }, [updateSelectedRoles]);

  useEffect(() => {
    const debounceSearchHandler = setTimeout(() => {
      validateSearchAndSelectedItems(
        searchValues,
        selectedItemsByInstance,
        setShowValidationError,
        setValidationErrorText,
      );
    }, 150);

    return () => clearTimeout(debounceSearchHandler);
  }, [searchValues, selectedItemsByInstance]);

  return (
    <Popup
      isOpen={isOpen}
      showIcon
      title='Invite Members'
      subTitle='Type or paste mail addresses, separated by spaces or commas'
      titleClassName='f-16-600 text-GRAY_950'
      iconCategory={ICON_SPRITE_TYPES.GENERAL}
      iconId='x-close'
      iconColor={COLORS.TEXT_PRIMARY}
      onClose={handleCloseInviteMembersPopup}
      popupWrapperClassName='bg-white rounded-t-3.5 border border-b-0 border-GRAY_400'
      closeOnClickOutside={false}
    >
      <div className='flex flex-col rounded-b-3.5 w-[458px] bg-white border border-t-0 border-GRAY_400'>
        <div className='flex flex-col px-4 py-6 gap-2'>
          <div className='flex flex-col'>
            <div className='flex flex-col gap-2'>
              {multiSelectInstances.map((id) => (
                <MultiSelectInput
                  key={id}
                  id={`invite-members-${id}`}
                  search={searchValues[id] || ''}
                  setSearch={(value) => handleSearchChange(id, value)}
                  isOpen={isOpen}
                  placeholderText={placeholderText}
                  roleOptions={TEAM_MEMBERS_PRIVILEGES_LIST}
                  inputArrayList={selectedItemsByInstance[id] || []}
                  setInputArrayList={(items) =>
                    setSelectedItemsByInstance((prev) => ({
                      ...prev,
                      [id]: items as ArrayListOption[],
                    }))
                  }
                  onValidateAndAdd={({ value, label, color }) => handleValidateAndAdd(id, { value, label, color })}
                  selectedRole={selectedRoleByInstance[id] ?? TEAM_MEMBERS_PRIVILEGES_LIST[0].value}
                  setSelectedRole={(role) =>
                    setSelectedRoleByInstance((prev) => ({
                      ...prev,
                      [id]: role as TEAM_MEMBERS_PRIVILEGES,
                    }))
                  }
                />
              ))}
            </div>
            {validationErrorText && showValidationError && (
              <span className='f-11-400 text-RED_700 mt-2 w-full flex text-start'>{validationErrorText}</span>
            )}
          </div>

          <Button
            type={BUTTON_TYPES.SECONDARY}
            id='add-duplicate-multiselect-user-invite'
            size={SIZE_TYPES.SMALL}
            className='mt-2 w-fit'
            disabled={disableAddBtn}
            onClick={() => setMultiSelectInstances((prev) => [...prev, prev?.length])}
          >
            <div className='flex gap-0.5'>
              <SvgSpriteLoader
                id='plus'
                width={14}
                height={14}
                color={disableAddBtn ? COLORS.GRAY_700 : COLORS.GRAY_1000}
              />
              Add
            </div>
          </Button>
        </div>

        <div className='flex justify-end border-t border-GRAY_200 py-4 px-5 w-full'>
          <Button
            type={BUTTON_TYPES.PRIMARY}
            id='send-user-invite'
            size={SIZE_TYPES.MEDIUM}
            disabled={!isInvitable}
            onClick={handleInviteMembers}
            isLoading={postInviteAudiencesIsLoading}
          >
            Send invite
          </Button>
        </div>
      </div>
    </Popup>
  );
};

export default InviteMembersPopup;
