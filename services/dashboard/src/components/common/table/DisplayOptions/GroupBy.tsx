import React, { useEffect, useState } from 'react';
import { AgGridReact } from 'ag-grid-react';
import { DRAG_ICON, ICON_SPRITE_TYPES } from 'constants/icons';
import Image from 'next/image';
import { SIZE_TYPES } from 'types/common/components';
import { defaultFnType } from 'types/commonTypes';
import Input from 'components/common/input';
import { MenuWrapper } from 'components/common/MenuWrapper';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

type GroupByProps = {
  onClose: defaultFnType;
  tableRef: React.RefObject<AgGridReact>;
};

const GroupBy: React.FC<GroupByProps> = ({ tableRef, onClose }) => {
  // State to manage grouped and available columns
  const [searchTerm, setSearchTerm] = useState('');
  const [groupedColumns, setGroupedColumns] = useState<string[]>([]);
  const [availableColumns, setAvailableColumns] = useState<string[]>([]);

  const handleDragStart = (column: string) => (event: React.DragEvent) => {
    event.dataTransfer.setData('text/plain', column);
  };

  const handleDropOnGroup = (event: React.DragEvent) => {
    const data = event.dataTransfer.getData('text/plain');
    const column: string = data;
    const latestColumns = tableRef?.current?.api?.getColumns() ?? [];
    const currentColumn = latestColumns.find((col) => col.getColDef()?.headerName === column);

    setGroupedColumns((prev) => (prev?.includes(column) ? prev : [...(prev ?? []), column]));
    setAvailableColumns((prev) => prev?.filter((col) => col !== column));
    tableRef?.current?.api?.applyColumnState({
      state: [{ colId: currentColumn?.getColId() ?? '', rowGroup: true, hide: true }],
    });
  };

  const handleDropOnAvailable = (event: React.DragEvent) => {
    const data = event.dataTransfer.getData('text/plain');

    const column: string = data;
    const latestColumns = tableRef?.current?.api?.getColumns() ?? [];
    const currentColumn = latestColumns.find((col) => col.getColDef()?.headerName === column);

    setAvailableColumns((prev) => (prev?.includes(column) ? prev : [...(prev ?? []), column]));
    setGroupedColumns((prev) => prev?.filter((col) => col !== column));
    tableRef?.current?.api?.applyColumnState({
      state: [{ colId: currentColumn?.getColId() ?? '', rowGroup: false }],
    });
  };

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    const latestColumns = tableRef?.current?.api?.getColumns() ?? [];
    const columnNames = latestColumns
      .map((col) => col.getColDef()?.headerName)
      .filter((column) => column !== undefined);

    setSearchTerm(value);
    if (value) {
      const filteredColumns = columnNames
        ?.filter((column) => column?.toLowerCase().includes(value.toLowerCase()))
        .filter((column) => column !== undefined);

      setAvailableColumns(filteredColumns);
    } else {
      setAvailableColumns(columnNames);
    }
  };

  const handleReset = () => {
    const latestColumns = tableRef?.current?.api?.getColumns() ?? [];
    const columnNames = latestColumns
      .map((col) => col.getColDef()?.headerName)
      .filter((column) => column !== undefined);

    setGroupedColumns([]);
    setAvailableColumns(columnNames);
    tableRef?.current?.api?.setRowGroupColumns([]);
  };

  const handleRemoveGroupedColumn = (column: string) => {
    setGroupedColumns((prev) => prev?.filter((col) => col !== column));
    setAvailableColumns((prev) => [...(prev ?? []), column]);
    const latestColumns = tableRef?.current?.api?.getColumns() ?? [];
    const currentColumn = latestColumns.find((col) => col.getColDef()?.headerName === column);

    tableRef?.current?.api?.applyColumnState({
      state: [{ colId: currentColumn?.getColId() ?? '', rowGroup: false, hide: false }],
    });
  };

  useEffect(() => {
    const latestColumns = tableRef?.current?.api?.getColumns() ?? [];
    const groupedColumns = tableRef?.current?.api?.getRowGroupColumns() ?? [];
    const groupedColumnNames = groupedColumns
      .map((col) => col.getColDef()?.headerName)
      .filter((column) => column !== undefined);

    const columnNames = latestColumns
      .map((col) => col.getColDef()?.headerName)
      .filter((column) => column !== undefined)
      ?.filter((col) => !groupedColumnNames.includes(col));

    setGroupedColumns(groupedColumnNames);
    setAvailableColumns(columnNames);
  }, [tableRef]);

  return (
    <MenuWrapper
      id='group-by'
      className='!absolute z-10 right-0 mt-1 min-w-[376px] min-h-[344px] h-fit'
      childrenWrapperClassName='!overflow-visible !min-h-[344px] h-fit !max-h-fit'
    >
      <div className='px-3 py-1'>
        <div className='flex items-center gap-1.5 py-2'>
          <SvgSpriteLoader
            id='arrow-narrow-left'
            iconCategory={ICON_SPRITE_TYPES.ARROWS}
            width={12}
            height={12}
            className='cursor-pointer'
            onClick={onClose}
          />
          <div className='f-12-500 text-GRAY_1000'>Group By</div>
        </div>
        {/* Grouped Columns */}
        <div
          onDragOver={(event) => event.preventDefault()}
          onDrop={handleDropOnGroup}
          className='border border-GRAY_500 rounded-md p-2.5 bg-BG_GRAY_2 min-h-[70px]'
        >
          <div className='f-12-400 text-GRAY_700 mb-3'>Drag columns here to group by</div>
          <div className='flex gap-1 flex-wrap overflow-y-auto max-h-[100px] overflow-x-visible'>
            {groupedColumns.map((col, index) => (
              <div className='flex gap-1 items-center' key={col}>
                <div
                  key={col}
                  className='border border-GRAY_400 rounded-md px-2 py-1 text-GRAY_1000 bg-white f-12-500 flex items-center gap-1.5'
                  draggable
                  onDragStart={handleDragStart(col)}
                >
                  <Image src={DRAG_ICON} width={14} height={14} alt='drag icon' />
                  <div>{col}</div>
                  <SvgSpriteLoader
                    id='x-close'
                    iconCategory={ICON_SPRITE_TYPES.GENERAL}
                    width={12}
                    height={12}
                    onClick={() => handleRemoveGroupedColumn(col)}
                    className='cursor-pointer'
                  />
                </div>
                {index < groupedColumns.length - 1 && (
                  <SvgSpriteLoader id='chevron-right' iconCategory={ICON_SPRITE_TYPES.ARROWS} width={12} height={12} />
                )}
              </div>
            ))}
          </div>
        </div>
        <Input
          placeholder='Search Columns...'
          size={SIZE_TYPES.XSMALL}
          noBorders
          focusClassNames='mt-3 mb-2 !pl-0'
          onChange={handleSearch}
          value={searchTerm}
          autoFocus
        />
        {/* Available Columns */}
        <div
          onDragOver={(event) => event.preventDefault()}
          onDrop={handleDropOnAvailable}
          className='flex flex-wrap gap-1.5 overflow-y-auto max-h-[150px] overflow-x-visible pb-2'
        >
          {availableColumns?.map((col) => (
            <div
              key={col}
              draggable
              onDragStart={handleDragStart(col)}
              className='border border-GRAY_400 rounded-md px-2 py-1 w-fit cursor-move text-GRAY_900 f-12-400 hover:bg-BG_GRAY_2'
            >
              {col}
            </div>
          ))}
        </div>
      </div>
      <div className='w-full flex flex-row-reverse f-12-500 text-GRAY_1000 py-2.5 px-3 border-t border-GRAY_400 absolute bottom-0 bg-white rounded-b-md'>
        <div className='cursor-pointer' onClick={handleReset}>
          Reset
        </div>
      </div>
    </MenuWrapper>
  );
};

export default GroupBy;
