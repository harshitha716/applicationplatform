import { useMemo, useState } from 'react';
import { DATASET_TABLE, ICON_SPRITE_TYPES } from 'constants/icons';
import { ACCOUNT_DATASET_OPTIONS } from 'modules/payments/connect-account/connect-account-dummydata';
import Image from 'next/image';
import Input from 'components/common/input';

interface SelectAccountDatasetProps {
  onSelectDataset: (dataset: string) => void;
}

const SelectAccountDataset = ({ onSelectDataset }: SelectAccountDatasetProps) => {
  const [search, setSearch] = useState('');

  const filteredDatasetOptions = useMemo(() => {
    return ACCOUNT_DATASET_OPTIONS.filter((item) => item.label.toLowerCase().includes(search.toLowerCase()));
  }, [search]);

  return (
    <div>
      <Input
        autoFocus
        placeholder='Search dataset...'
        inputClassName='border-none w-full focus:outline-none focus:border-none focus:shadow-none !px-0'
        value={search}
        trailingIconProps={
          search
            ? {
                id: 'x',
                iconCategory: ICON_SPRITE_TYPES.GENERAL,
                onClick: () => setSearch(''),
              }
            : undefined
        }
        onChange={(e) => {
          if (e?.target?.value !== undefined) {
            setSearch(e.target.value);
          }
        }}
      />
      {filteredDatasetOptions.map((item) => (
        <div
          onClick={() => onSelectDataset(item.value)}
          key={item.value}
          className='flex items-center gap-2.5 px-2.5 py-3 rounded-lg hover:bg-GRAY_100 cursor-pointer'
        >
          <Image src={DATASET_TABLE} alt='dataset-table' width={20} height={20} />
          <div className='f-13-500'>{item.label}</div>
        </div>
      ))}
      {filteredDatasetOptions.length === 0 && (
        <div className='f-12-450 text-GRAY_700 mb-6.5 block animate-pulse'>No results found</div>
      )}
    </div>
  );
};

export default SelectAccountDataset;
