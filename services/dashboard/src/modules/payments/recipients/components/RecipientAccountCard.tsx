import { useState } from 'react';
import { RECIPIENT_ACCOUNTS } from 'modules/payments/recipients/recipient.dummy';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

const RecipientAccountCard = () => {
  const [isDetailsOpen, setIsDetailsOpen] = useState(false);

  const accountDetails = RECIPIENT_ACCOUNTS[0];

  return (
    <div
      className='w-full rounded-md border border-BORDER_GRAY_400 overflow-hidden'
      onClick={() => setIsDetailsOpen((prev) => !prev)}
    >
      <div className='px-2 py-2.5 flex items-center gap-1.5 bg-BACKGROUND_GRAY_2 f-11-400 cursor-pointer'>
        <SvgSpriteLoader id='bank' size={14} />
        Account Name {accountDetails.account_name}
      </div>
      {isDetailsOpen && (
        <div className='px-2.5 py-3 flex flex-col gap-3.5'>
          {accountDetails.Account_details.map((detail, index) => (
            <div key={index} className='flex items-center gap-4'>
              <div className='f-12-400 text-GRAY_700 w-[150px]'>{detail.label}</div>
              <div className='f-11-400'>{detail.value}</div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default RecipientAccountCard;
