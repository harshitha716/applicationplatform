import React, { FC, memo } from 'react';
import { MapAny } from 'types/commonTypes';
import HeaderAction from 'components/common/agGridTable/components/AgGridTableHeaderAction';

const AgGridTableHeader: FC<MapAny> = (props: any) => {
  const { columnDefs, columnHeaderMap, headerWrapperClassName, headerClassName } = props;
  const colId = props?.column?.colId;
  const columnIndex = columnDefs?.findIndex((col: MapAny) => col.field === colId);
  const column = columnDefs?.[columnIndex];
  const { isCustomHeader, CustomHeaderComponent, isActionEnabled, ...headerActionProps } =
    columnHeaderMap?.[column?.field] || {};
  const isLastCellInRow = props?.column?.colId === columnDefs?.[columnDefs?.length - 1]?.field;

  return (
    <div
      className={`w-full h-full flex items-center justify-between p-2 ${
        isLastCellInRow ? 'border-r' : ''
      } w-full border-y border-l border-DIVIDER_SAIL_2 ${headerWrapperClassName}`}
    >
      {isCustomHeader && !!CustomHeaderComponent && <CustomHeaderComponent {...props} />}
      {!isCustomHeader && <div className={headerClassName}>{column?.headerName ?? ''}</div>}

      {isActionEnabled && <HeaderAction {...headerActionProps} id={colId} columnIndex={columnIndex} />}
    </div>
  );
};

export default memo(AgGridTableHeader);
