import React from 'react';
import { SIZE_TYPES } from 'types/common/components';
import { BUTTON_TYPES } from 'types/components/button.type';
import { Button } from 'components/common/button/Button';

const ChooseRoleMembers = () => {
  return (
    <div className='relative w-fit flex flex-col'>
      <Button type={BUTTON_TYPES.SECONDARY} id='send-user-invite-btn' size={SIZE_TYPES.SMALL} className='!bg-GRAY_100'>
        Share
      </Button>
      <div className='z-1000 relative'>
        <div className='absolute bottom-0 right-0 flex h-[10rem] w-[20rem] bg-red-400'>this is the dropdown</div>
      </div>
    </div>
  );
};

export default ChooseRoleMembers;
