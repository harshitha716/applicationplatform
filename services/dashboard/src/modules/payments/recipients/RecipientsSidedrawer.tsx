import React, { FC, useState } from 'react';
import AddRecipientAccount from 'modules/payments/recipients/AddRecipientAccount';
import RecipientDetails from 'modules/payments/recipients/RecipientDetails';
import RecipientsList from 'modules/payments/recipients/RecipientsList';
import { defaultFnType } from 'types/commonTypes';
import SideDrawer from 'components/common/SideDrawer/SideDrawer';
import { SIDE_DRAWER_TYPES } from 'components/common/SideDrawer/sideDrawer.types';

type RecipientsSideDrawerProps = {
  onClose: defaultFnType;
  isOpen: boolean;
};

const RecipientsSideDrawer: FC<RecipientsSideDrawerProps> = ({ onClose, isOpen }) => {
  const [onRecipientDetails, setOnRecipientDetails] = useState<string>('');
  const [isAddRecipient, setIsAddRecipient] = useState<boolean>(false);

  const renderStep = () => {
    if (onRecipientDetails) return <RecipientDetails onBack={() => setOnRecipientDetails('')} />;
    if (isAddRecipient) return <AddRecipientAccount />;

    return (
      <RecipientsList
        onRecipientDetails={(id) => setOnRecipientDetails(id)}
        onAddRecipient={() => setIsAddRecipient(true)}
      />
    );
  };

  return (
    <SideDrawer
      id='json-preview-sidebar'
      isOpen={isOpen}
      onClose={onClose}
      hideCloseButton
      type={SIDE_DRAWER_TYPES.SECONDARY}
      className='h-screen overflow-hidden '
      childrenWrapperClassName='!p-0 overflow-y-scroll'
    >
      {renderStep()}
    </SideDrawer>
  );
};

export default RecipientsSideDrawer;
