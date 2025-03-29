import { useEffect, useState } from 'react';
import { DATASET_TABLE } from 'constants/icons';
import {
  ACCOUNT_DATASET_COLUMNS_MAPPING,
  ACCOUNT_DATASET_COLUMNS_MAPPING_OPTIONS,
} from 'modules/payments/connect-account/connect-account-dummydata';
import Image from 'next/image';
import { SIZE_TYPES } from 'types/common/components';
import { defaultFn, MapAny } from 'types/commonTypes';
import { Button } from 'components/common/button/Button';
import { Dropdown } from 'components/common/dropdown';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

interface DatasetColumnsMappingProps {
  onSelectDataset: (dataset: string) => void;
}

const DatasetColumnsMapping = ({ onSelectDataset }: DatasetColumnsMappingProps) => {
  const [isLoading, setIsLoading] = useState(false);
  const [selectedColumnMapping, setSelectedColumnMapping] = useState<MapAny>({});

  // TODO: remove this
  console.log('selectedColumnMapping', selectedColumnMapping);

  useEffect(() => {
    setTimeout(() => {
      setIsLoading(false);
    }, 1000);
  }, []);

  const handleChangeColumnMapping = (columnName: string, selectedOption: any) => {
    setSelectedColumnMapping((prev) => ({ ...prev, [columnName]: selectedOption }));
  };

  return (
    <div>
      <div className='flex items-center gap-1.5 bg-GRAY_100 rounded-lg px-3 py-2.5 w-full mb-6'>
        <Image src={DATASET_TABLE} alt='dataset-table' width={20} height={20} />
        <div className='f-13-500 grow'>Accounts</div>
        <SvgSpriteLoader id='x-close' size={14} className='text-GRAY_900' onClick={() => onSelectDataset('')} />
      </div>
      {isLoading ? (
        <div className='f-12-450 text-GRAY_700 mb-6.5 block animate-pulse'>
          Verifying connection, this might take a second...
        </div>
      ) : (
        <div className=''>
          <div className='f-13-450 text-GRAY_700 text-left flex gap-4'>
            <div className='grow'>Required Attribute</div>
            <div className='w-65'>Dataset column</div>
          </div>
          {ACCOUNT_DATASET_COLUMNS_MAPPING.map((item) => (
            <div key={item.value} className='my-4.5 flex w-full'>
              <div className='f-12-500 flex items-center justify-between gap-2 min-w-[178px]'>
                <div>{item.label}</div>
                <SvgSpriteLoader id='arrow-right' size={16} className='text-GRAY_900 mx-3' />
              </div>
              <div className='w-65'>
                <Dropdown
                  placeholder='Select a column'
                  customStyles={{
                    control: {
                      width: '100%',
                      fontSize: '15px',
                      fontWeight: '400',
                    },
                    menu: {
                      width: '100%',
                    },
                  }}
                  size={SIZE_TYPES.XSMALL}
                  options={ACCOUNT_DATASET_COLUMNS_MAPPING_OPTIONS}
                  id={`${item.value}-multi-select-input-dropdown`}
                  eventCallback={defaultFn}
                  onChange={(value) => handleChangeColumnMapping(item.value, value)}
                  customClass={{
                    fontSize: 'f-12-450',
                  }}
                  customClassNames={{
                    placeholder: 'f-12-500',
                  }}
                  menuOptionClasses={{
                    contentWrapper: 'py-2',
                  }}
                  customDropdownIndicatorSize={14}
                />
              </div>
            </div>
          ))}
          <Button id='connect-account-button' size={SIZE_TYPES.MEDIUM} className='w-[56px] !mt-[26px]'>
            Done
          </Button>
        </div>
      )}
    </div>
  );
};

export default DatasetColumnsMapping;
