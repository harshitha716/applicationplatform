import { FC } from 'react';
import { COLORS } from 'constants/colors';
import { RECIPIENT_CARD_ACTION_ITEMS } from 'modules/payments/payments.constant';
import { SIZE_TYPES } from 'types/common/components';
import { defaultFn } from 'types/commonTypes';
import TooltipButton from 'components/common/button/TooltipButton';
import { TooltipPositions } from 'components/common/tooltip';

const RecipientCard: FC = () => {
  return (
    <div className='flex items-center justify-between px-1.5 py-1 hover:bg-GRAY_50 cursor-pointer rounded-md hover:z-50'>
      <div className='flex items-center gap-1.5'>
        <div className='w-6 h-6 flex items-center justify-center rounded-full bg-BLUE_200 f-12-500'>SS</div>
        <div>
          <div className='f-12-500'>Satabdi S</div>
          <div className='f-11-400 text-GRAY_700'>Account Name 4453</div>
        </div>
      </div>
      <div className='flex items-center gap-2.5'>
        <div className='f-11-400'>3 Accounts</div>
        {RECIPIENT_CARD_ACTION_ITEMS.map((item, index) => (
          <TooltipButton
            key={item.id}
            id='recipient-card-action'
            tooltipBodyClassName='z-50'
            onClick={defaultFn}
            tooltipBody={item.tooltipBody}
            className='border-none'
            tooltipColor={COLORS.BLACK}
            buttonSize={SIZE_TYPES.XSMALL}
            tooltipPosition={
              index === RECIPIENT_CARD_ACTION_ITEMS.length - 1 ? TooltipPositions.LEFT : TooltipPositions.BOTTOM
            }
            buttonIcon={item.icon}
          />
        ))}
      </div>
    </div>
  );
};

export default RecipientCard;
