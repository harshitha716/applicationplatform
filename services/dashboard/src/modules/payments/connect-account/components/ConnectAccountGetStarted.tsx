import { FC } from 'react';
import { CONNECT_ACCOUNT } from 'constants/icons';
import Image from 'next/image';
import { SIZE_TYPES } from 'types/common/components';
import { BUTTON_TYPES, ICON_POSITION_TYPES } from 'types/components/button.type';
import { Button } from 'components/common/button/Button';

interface ConnectAccountGetStartedType {
  onNextStep: (step: number) => void;
}

const ConnectAccountGetStarted: FC<ConnectAccountGetStartedType> = ({ onNextStep }) => {
  return (
    <div className='w-[312px] pt-[166px]'>
      <Image src={CONNECT_ACCOUNT} alt='connect-account-get-started' width={76} height={76} className='mb-6' />
      <div className='f-16-550'>Payments</div>
      <div className='f-13-450 text-GRAY_700 mt-2 mb-6'>
        Move money anywhere in the world—seamlessly. Send single or bulk payments in a few clicks, save templates for
        faster transactions, and set up custom access controls and approval policies to stay in full control. Fast,
        secure, and built for teams like yours
      </div>
      <div className='flex gap-2'>
        <Button
          id='connect-account'
          onClick={() => onNextStep(1)}
          size={SIZE_TYPES.SMALL}
          iconProps={{ id: 'link-04', size: 14, className: 'text-white' }}
          iconPosition={ICON_POSITION_TYPES.LEFT}
        >
          Connect Accounts
        </Button>
        <Button id='guid' size={SIZE_TYPES.SMALL} type={BUTTON_TYPES.SECONDARY}>
          Guide
        </Button>
      </div>
    </div>
  );
};

export default ConnectAccountGetStarted;
