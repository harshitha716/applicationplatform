import React, { FC, useRef, useState } from 'react';
import { AgGridReact } from 'ag-grid-react';
import { useOnClickOutside } from 'hooks';
import { SIZE_TYPES } from 'types/common/components';
import { BUTTON_TYPES, ICON_POSITION_TYPES } from 'types/components/button.type';
import { Button } from 'components/common/button/Button';
import { MenuWrapper } from 'components/common/MenuWrapper';
import ColumnListing from 'components/common/table/DisplayOptions/ColumnListing';
import DisplayOptionItem from 'components/common/table/DisplayOptions/DisplayOptionItem';
import GroupBy from 'components/common/table/DisplayOptions/GroupBy';
import { DisplayOptionsList } from 'components/common/table/table.constants';
import { DISPLAY_OPTIONS } from 'components/common/table/table.types';

type DisplayOptionsProps = {
  tableRef: React.RefObject<AgGridReact>;
  datasetId: string;
};

const DisplayOptions: FC<DisplayOptionsProps> = ({ tableRef, datasetId }) => {
  const menuRef = useRef<HTMLDivElement>(null);
  const [isOpen, setIsOpen] = useState(false);
  const [isColumnListingOpen, setIsColumnListingOpen] = useState(false);
  const [isGroupByOpen, setIsGroupByOpen] = useState(false);
  // TODO: Implement fx later
  // const [isCurrencyOpen, setIsCurrencyOpen] = useState(false);
  // const [currency, setCurrency] = useState<string>('');

  useOnClickOutside(menuRef, () => {
    setIsOpen(false);
    setIsColumnListingOpen(false);
    setIsGroupByOpen(false);
    // TODO: Implement fx later
    // setIsCurrencyOpen(false);
  });

  const handleClick = (id: DISPLAY_OPTIONS) => {
    setIsOpen(false);
    switch (id) {
      case DISPLAY_OPTIONS.COLUMNS:
        setIsColumnListingOpen(true);
        break;
      case DISPLAY_OPTIONS.GROUP_BY:
        setIsGroupByOpen(true);
        break;
      case DISPLAY_OPTIONS.CURRENCY:
        // TODO: Implement fx later
        // setIsCurrencyOpen(true);
        break;
    }
  };

  const handleCloseColumnListing = () => {
    setIsColumnListingOpen(false);
    setIsOpen(true);
  };

  const handleCloseGroupBy = () => {
    setIsGroupByOpen(false);
    setIsOpen(true);
  };

  return (
    <div className='relative' ref={menuRef}>
      <Button
        id='display-options'
        onClick={() => setIsOpen(!isOpen)}
        type={BUTTON_TYPES.SECONDARY}
        size={SIZE_TYPES.XSMALL}
        iconPosition={ICON_POSITION_TYPES.LEFT}
        iconProps={{
          id: 'settings-04',
        }}
      >
        Display
      </Button>
      {isOpen && (
        <MenuWrapper
          id='display-options'
          className='!absolute z-10 p-1 right-0 mt-1 w-[180px]'
          childrenWrapperClassName='text-GRAY_900 !overflow-y-auto'
        >
          {DisplayOptionsList.map((option) => (
            <DisplayOptionItem key={option.id} {...option} onClick={handleClick} />
          ))}
        </MenuWrapper>
      )}
      {isColumnListingOpen && (
        <ColumnListing tableRef={tableRef} onClose={handleCloseColumnListing} datasetId={datasetId} />
      )}
      {isGroupByOpen && <GroupBy onClose={handleCloseGroupBy} tableRef={tableRef} />}
    </div>
  );
};

export default DisplayOptions;
