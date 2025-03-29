import { useState } from 'react';
import DatasetColumnsMapping from 'modules/payments/connect-account/components/DatasetColumnsMapping';
import SelectAccountDataset from 'modules/payments/connect-account/components/SelectAccountDataset';
import { CONNECT_ACCOUNT_DESCRIPTION, CONNECT_ACCOUNT_TITLE } from 'modules/payments/payments.constant';
import SvgSpriteLoader from 'components/SvgSpriteLoader';
interface ConnectAccountSelectDatasetProps {
  stateChange: (state: number) => void;
}

const ConnectAccountSelectDataset = ({ stateChange }: ConnectAccountSelectDatasetProps) => {
  const [selectedDataset, setSelectedDataset] = useState<string>('');

  return (
    <div className='w-[436px] pt-[76px]'>
      <SvgSpriteLoader id='arrow-narrow-left' size={16} className='mb-3.5' onClick={() => stateChange(0)} />
      <div className='f-16-550 text-GRAY_1000 mt-4.5 mb-2'>{CONNECT_ACCOUNT_TITLE}</div>
      <div className='f-12-450 text-GRAY_700 mb-[26px]'>{CONNECT_ACCOUNT_DESCRIPTION}</div>
      {selectedDataset ? (
        <DatasetColumnsMapping onSelectDataset={setSelectedDataset} />
      ) : (
        <SelectAccountDataset onSelectDataset={setSelectedDataset} />
      )}
    </div>
  );
};

export default ConnectAccountSelectDataset;
