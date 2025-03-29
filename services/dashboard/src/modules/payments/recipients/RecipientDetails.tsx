import { FC } from 'react';
import RecipientAccountCard from 'modules/payments/recipients/components/RecipientAccountCard';
import { RECIPIENT_ACCOUNT_DETAILS } from 'modules/payments/recipients/recipient.dummy';
import { SIZE_TYPES } from 'types/common/components';
import { defaultFnType } from 'types/commonTypes';
import { BUTTON_TYPES, ICON_POSITION_TYPES } from 'types/components/button.type';
import { Button } from 'components/common/button/Button';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

type RecipientDetailsProps = {
  onBack: defaultFnType;
};

const RecipientDetails: FC<RecipientDetailsProps> = ({ onBack }) => {
  return (
    <div className='px-4.5 py-6.5 flex flex-col gap-8 overflow-y-scroll'>
      <div className='flex items-center gap-3'>
        <SvgSpriteLoader id='arrow-narrow-left' size={14} onClick={onBack} />
        <div className='flex items-center gap-2.5'>
          <div className='w-6 h-6 flex items-center justify-center rounded-full bg-BLUE_200 f-12-500'>SS</div>
          <div>
            <div className='f-16-600'>Satabdi S</div>
          </div>
        </div>
      </div>
      <div className='flex flex-col gap-2.5'>
        {RECIPIENT_ACCOUNT_DETAILS.map((item, index) => (
          <div key={index} className='flex items-center gap-4'>
            <div className='f-12-400 text-GRAY_700 w-[150px]'>{item.label}</div>
            <div className='f-11-400'>{item.value}</div>
          </div>
        ))}
      </div>
      <div className='flex flex-col gap-3'>
        <div className=' flex justify-between items-center f-13-500'>
          Accounts
          <Button
            id='add-account'
            size={SIZE_TYPES.XSMALL}
            type={BUTTON_TYPES.SECONDARY}
            iconPosition={ICON_POSITION_TYPES.LEFT}
            iconProps={{
              id: 'plus',
              size: 14,
            }}
          >
            Add
          </Button>
        </div>
        {Array.from({ length: 3 }).map((_, index) => (
          <RecipientAccountCard key={index} />
        ))}
      </div>
    </div>
  );
};

export default RecipientDetails;
