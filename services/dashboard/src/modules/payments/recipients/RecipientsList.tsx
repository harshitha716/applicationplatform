import React, { FC, useState } from 'react';
import RecipientCard from 'modules/payments/recipients/components/RecipientCard';
import { SIZE_TYPES } from 'types/common/components';
import { defaultFnType } from 'types/commonTypes';
import { BUTTON_TYPES, ICON_POSITION_TYPES } from 'types/components/button.type';
import { Button } from 'components/common/button/Button';
import Input from 'components/common/input';

type RecipientsListProps = {
  onRecipientDetails: (recipientDetails: string) => void;
  onAddRecipient: defaultFnType;
};

const RecipientsList: FC<RecipientsListProps> = ({ onRecipientDetails, onAddRecipient }) => {
  const [search, setSearch] = useState('');

  return (
    <div className='p-6'>
      <div className='w-full flex items-center justify-between'>
        <div className='f-16-600 mb-4'>Recipients</div>
        <Button
          size={SIZE_TYPES.SMALL}
          type={BUTTON_TYPES.SECONDARY}
          iconPosition={ICON_POSITION_TYPES.LEFT}
          id='add-recipient'
          onClick={onAddRecipient}
          iconProps={{
            id: 'plus',
            size: 14,
          }}
        >
          Add
        </Button>
      </div>
      <Input
        type='text'
        placeholder='Search...'
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        className='mb-4'
        inputClassName='border-none !px-0'
      />
      <div className='flex flex-col gap-2'>
        {Array.from({ length: 10 }).map((_, index) => (
          <div key={index} onClick={() => onRecipientDetails('Siddharth')} className='hover:z-1000'>
            <RecipientCard />
          </div>
        ))}
      </div>
    </div>
  );
};

export default RecipientsList;
