import React, { FC, useEffect, useMemo, useRef, useState } from 'react';
import {
  usePostAddTeamToAudienceMutation,
  usePostAddTeamToOrganizationMutation,
  useRemoveTeamFromAudienceMutation,
} from 'apis/people';
import { COLORS, TEAMS_COLORS } from 'constants/colors';
import { useOnClickOutside } from 'hooks';
import CustomTeamsDropdown from 'modules/team/components/members/CustomTeamsDropdown';
import { TEAM_PERMISSION_TOAST_MSG } from 'modules/team/people.constants';
import {
  CustomTeamsDropdownPropsType,
  MembersTeamPropsType,
  PostAddTeamToAudiencePayload,
  PostTeamsByOrganizationIdPayload,
} from 'modules/team/people.types';
import { MapAny } from 'types/commonTypes';
import { checkIfCurrentUserIsMember } from 'utils/accessPermission/accessPermission.utils';
import { cn, cyclicIterator } from 'utils/common';
import { toast } from 'components/common/toast/Toast';
import MultiSelectInput from 'components/multiSelectInput/MultiSelectInput';

const MembersTeam: FC<MembersTeamPropsType> = ({ organizationId, teamsData, userId, userMappedTeams }) => {
  const isMember = checkIfCurrentUserIsMember();
  const teamsRowRef = useRef<HTMLDivElement>(null);
  const teamsRandomColorRef = useRef(cyclicIterator(TEAMS_COLORS));
  const [postAddTeamToOrganization] = usePostAddTeamToOrganizationMutation();
  const [postAddTeamToAudience] = usePostAddTeamToAudienceMutation();
  const [removeTeamFromAudience] = useRemoveTeamFromAudienceMutation();
  const [search, setSearch] = useState<string>('');
  const [isCustomInputFocused, setIsCustomInputFocused] = useState<boolean>(false);
  const [openFullViewTeamTags, setOpenFullViewTeamTags] = useState<boolean>(false);
  const [randomColor, setRandomColor] = useState(() => teamsRandomColorRef.current());
  const [selectedItems, setSelectedItems] = useState<
    {
      value: string;
      label: string;
      valid: boolean;
      color?: string;
      isNew?: boolean;
      teamId?: string;
      teamMembershipId?: string;
    }[]
  >(userMappedTeams);

  const handleAddTeamToOrg = async (payload: PostTeamsByOrganizationIdPayload) => {
    postAddTeamToOrganization({ organizationId, payload })
      .unwrap()
      .then((res) => {
        const teamId = res?.team_id;

        handleAddTeamToAudience({ user_id: userId, team_id: teamId });
      })
      .catch(() => {
        toast.error(TEAM_PERMISSION_TOAST_MSG.TEAM_CREATE_ERROR);
      });
  };

  const handleAddTeamToAudience = async (payload: PostAddTeamToAudiencePayload) => {
    postAddTeamToAudience({ organizationId, teamId: payload?.team_id, payload })
      .unwrap()
      .then(() => {
        toast.success(TEAM_PERMISSION_TOAST_MSG.TEAM_ASSIGN_SUCCESS);
      })
      .catch(() => {
        toast.error(TEAM_PERMISSION_TOAST_MSG.TEAM_ASSIGN_ERROR);
      });
  };

  const handleCheckIfTeamExists = (teamInfo: PostTeamsByOrganizationIdPayload) => {
    const teamId = teamsData.find(
      (team) => team?.name === teamInfo?.name && team?.metadata?.color_hex_code === teamInfo?.color_hex_code,
    )?.team_id;

    const updatedTeamInfo = {
      user_id: userId,
      team_id: teamId ?? '',
    };

    if (teamId) {
      handleAddTeamToAudience(updatedTeamInfo);
    } else {
      handleAddTeamToOrg(teamInfo);
    }
  };

  const handleRemoveAudienceFromTeam = (item: MapAny) => {
    const membershipId = item?.teamMembershipId;
    const teamId = item?.teamId;

    if (!teamId || !membershipId) {
      toast.error(TEAM_PERMISSION_TOAST_MSG.INVALID_TEAM_ERROR);

      return;
    }

    // optimistic delete
    setSelectedItems((prev) => prev?.filter((selected) => selected?.teamId !== teamId));

    const payload = {
      team_id: teamId,
      team_membership_id: membershipId,
    };

    removeTeamFromAudience({ organizationId, teamId, payload })
      .unwrap()
      .then(() => {
        toast.success(TEAM_PERMISSION_TOAST_MSG.TEAM_REMOVE_SUCCESS);
      })
      .catch(() => {
        toast.error(TEAM_PERMISSION_TOAST_MSG.TEAM_REMOVE_ERROR);
      });
  };

  const handleValidateAndAdd = ({
    value,
    color,
  }: {
    value: string;
    label: string;
    color?: string;
    isNew?: boolean;
  }) => {
    if (!value) return;

    // optimistic update
    setSelectedItems((prev) => [
      ...prev,
      {
        value,
        label: value,
        valid: true,
        color: color ?? randomColor,
        isNew: true,
      },
    ]);

    const payload = {
      name: value,
      description: '',
      color_hex_code: color ?? randomColor,
    };

    handleCheckIfTeamExists(payload);
  };

  const handleOptionSelection = (option: { value: string; label: string; color?: string; isNew?: boolean }) => {
    // optimistic update
    setSelectedItems((prev) => [
      ...prev,
      {
        value: option?.value,
        label: option?.value,
        valid: true,
        color: option?.color ?? randomColor,
        isNew: option?.isNew,
      },
    ]);

    const payload = {
      name: option?.value,
      description: '',
      color_hex_code: option?.color ?? randomColor,
    };

    handleCheckIfTeamExists(payload);
  };

  const filteredOptionListsData = [
    ...(teamsData
      ?.filter((item) => !selectedItems.some((selected) => selected?.value === item?.name))
      .map((member) => ({
        label: member?.name ?? '',
        value: member?.name ?? '',
        color: member?.metadata?.color_hex_code ?? randomColor,
        isNew: false,
      })) ?? []),
    ...[
      {
        label: search,
        value: search,
        color: randomColor,
        isNew: true,
      },
    ],
  ];

  useEffect(() => {
    setSelectedItems(userMappedTeams);
  }, [userMappedTeams]);

  useEffect(() => {
    if (!search) {
      const newColor = teamsRandomColorRef.current();

      setRandomColor(newColor);
    }
  }, [search]);

  const memoizedDropdown = useMemo(() => {
    const MemoizedDropdownComponent = (props: CustomTeamsDropdownPropsType) => (
      <CustomTeamsDropdown {...props} randomColor={randomColor} />
    );

    MemoizedDropdownComponent.displayName = 'memoized-teams-dropdown-component';

    return MemoizedDropdownComponent;
  }, [randomColor]);

  const handleCloseFullViewTeamTags = () => {
    setOpenFullViewTeamTags(false);
  };

  const handleToggleFullViewTeamTags = () => {
    setOpenFullViewTeamTags((prev) => !prev);
  };

  useOnClickOutside(teamsRowRef, handleCloseFullViewTeamTags);

  return (
    <div
      className='relative f-12-400 text-GRAY_1000 h-full flex items-center justify-start text-left py-2 px-2 overflow-visible'
      ref={teamsRowRef}
      onClick={handleToggleFullViewTeamTags}
    >
      {isMember ? (
        <div className={cn('flex flex-nowrap overflow-hidden gap-1.5', openFullViewTeamTags && 'flex-wrap')}>
          {selectedItems.map((item, index) => (
            <span
              key={index}
              className='f-12-400 text-GRAY_1000 flex px-1.5 py-0.5 w-fit rounded capitalize'
              style={{ backgroundColor: item?.color ?? COLORS.WHITE }}
            >
              {item?.label}
            </span>
          ))}
        </div>
      ) : (
        <MultiSelectInput
          id='select-team'
          search={search}
          setSearch={setSearch}
          inputArrayList={selectedItems}
          setInputArrayList={setSelectedItems}
          optionsList={filteredOptionListsData}
          customOptionsListDropdown={memoizedDropdown}
          onValidateAndAdd={handleValidateAndAdd}
          onSelectOption={handleOptionSelection}
          placeholderText='Add team'
          isOpen={false}
          wrapperClassName='border-none rounded-none shadow-none f-12-400'
          inputWrapperClassName={cn(isCustomInputFocused ? 'flex-wrap' : 'flex-nowrap', 'p-0')}
          multiSelectInputClassName='f-12-400 !rounded-none'
          setIsCustomInputFocused={setIsCustomInputFocused}
          selectOnlyFromList
          onCustomDeleteFn={handleRemoveAudienceFromTeam}
        />
      )}
    </div>
  );
};

export default MembersTeam;
