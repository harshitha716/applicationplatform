import { ZAMP_ICON } from 'constants/icons';
import Image from 'next/image';
import { SIZE_TYPES } from 'types/common/components';
import { defaultFnType } from 'types/commonTypes';
import { BUTTON_TYPES } from 'types/components/button.type';
import { Button } from 'components/common/button/Button';

type MembershipRequestedProps = {
  text: string;
  subText: string;
  userEmail: string;
  actionItems: {
    text: string;
    onClick: defaultFnType;
  }[];
};

export const MembershipRequested = (props: MembershipRequestedProps) => {
  const { text, subText, userEmail, actionItems } = props;

  return (
    <div className='w-screen h-screen flex flex-col bg-white justify-center items-center'>
      <div>
        <Image
          width={60}
          height={60}
          alt='zamp logo'
          className='w-8 align-middle cursor-pointer'
          src={ZAMP_ICON}
          priority={true}
        />
      </div>
      <div className='flex flex-col items-center justify-center'>
        <span className='f-16-600 mt-4'>{text}</span>
        <span className='f-13-400 text-GRAY_600 mt-4'>{subText}</span>
        <span className='f-13-400 text-GRAY_600 mt-4'>You are logged in as</span>
        <span className='f-13-600 text-GRAY_950 mt-1'>{userEmail}</span>
      </div>
      <div className='flex gap-2.5 mt-6'>
        {actionItems.map((actionItem) => (
          <Button
            key={actionItem.text}
            type={BUTTON_TYPES.SECONDARY}
            id='send-user-invite-btn'
            size={SIZE_TYPES.SMALL}
            onClick={actionItem.onClick}
          >
            {actionItem.text}
          </Button>
        ))}
      </div>
    </div>
  );
};
