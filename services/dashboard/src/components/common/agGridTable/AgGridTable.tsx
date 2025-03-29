import React, { FC, RefObject, useEffect, useMemo, useState } from 'react';
import { ColDef, GridApi } from 'ag-grid-community';
import { MapAny } from 'types/commonTypes';
import {
  AG_GRID_CELL_CLASSNAME,
  AG_GRID_CELL_STYLE,
  AG_GRID_CELL_TEXT_WRAP_STYLE,
  AG_GRID_LAST_CELL_STYLE,
  COLUMN_HEADER_HEIGHT,
  ROW_HEIGHT,
} from 'components/common/agGridTable/agGridTable.constants';
import AgGridTableHeader from 'components/common/agGridTable/components/AgGridTableHeader';
import { AgGridTableHeaderActionProps } from 'components/common/agGridTable/components/AgGridTableHeaderAction';
import AgGridTooltipComponent from 'components/common/agGridTable/components/AgGridTooltipComponent';
import SpreadsheetGrid from 'components/common/agGridTable/SpreadsheetGrid';

export interface ColumnDef extends ColDef {
  textWrap?: boolean;
}

interface ColumnHeaderMap {
  [key: string]: {
    isCustomHeader?: boolean;
    CustomHeaderComponent?: FC<any>;
    isActionEnabled?: boolean;
  } & AgGridTableHeaderActionProps;
}
interface AgGridTableProps {
  columnDefs: ColumnDef[];
  columnHeaderMap?: ColumnHeaderMap;
  data?: MapAny[];
  componentProps?: MapAny;
  gridProps?: MapAny;
  refreshCellsCount?: number;
  wrapperClassName?: string;
  wrapperStyle?: MapAny;
  headerWrapperClassName?: string;
  headerClassName?: string;
  hideGridOnZeroData?: boolean;
  onReady?: (params: MapAny) => void;
  resetDataSourceCount?: number;
  setDataSourceWithFilters?: (gridApi: RefObject<GridApi>) => void;
}

const AgGridTable: FC<AgGridTableProps> = ({
  columnDefs: columnDefsData,
  columnHeaderMap = {},
  data = [],
  componentProps = {},
  gridProps = {},
  refreshCellsCount,
  wrapperClassName = '',
  wrapperStyle = {},
  headerWrapperClassName = '',
  headerClassName = 'f-12-600 text-GRAY_1000',
  hideGridOnZeroData,
  onReady,
  resetDataSourceCount,
  setDataSourceWithFilters,
}) => {
  const [columnDefs, setColumnDefs] = useState<any[]>([]);
  const [, setHorizontalScrollCount] = useState(0);

  const handleHorizontalScroll = () => {
    setHorizontalScrollCount((prev) => prev + 1);
  };

  useEffect(() => {
    const updatedColumnDefs = columnDefsData?.map((columnData, columnIndex) => {
      const { textWrap, cellStyle = {}, cellClass, ...column } = columnData || {};
      const CellComponent = column?.cellRenderer ?? null;
      const isLastCellInRow = columnIndex === columnDefsData?.length - 1;

      return {
        ...column,
        ...(textWrap ? { autoHeight: true } : {}),
        cellStyle: {
          ...AG_GRID_CELL_STYLE,
          ...(isLastCellInRow ? AG_GRID_LAST_CELL_STYLE : {}),
          ...(textWrap ? AG_GRID_CELL_TEXT_WRAP_STYLE : {}),
          ...cellStyle,
        },
        cellClass: cellClass ?? AG_GRID_CELL_CLASSNAME,
        cellRenderer: CellComponent
          ? (props: any) => <CellComponent {...props} componentProps={componentProps} />
          : undefined,
      };
    });

    setColumnDefs(updatedColumnDefs);
  }, [refreshCellsCount]);

  const CustomHeaderComponent: FC<MapAny> = useMemo(() => {
    const Component = (props: MapAny) => (
      <AgGridTableHeader
        {...props}
        columnDefs={columnDefs}
        columnHeaderMap={columnHeaderMap}
        headerWrapperClassName={headerWrapperClassName}
        headerClassName={headerClassName}
      />
    );

    Component.displayName = 'CustomHeaderComponentInner';

    return Component;
  }, [columnDefs?.length]);

  const defaultColDef = {
    headerComponent: CustomHeaderComponent,
    tooltipComponent: AgGridTooltipComponent,
  };

  return (
    <div className={`w-full h-full ${wrapperClassName}`} style={wrapperStyle}>
      <SpreadsheetGrid
        columnDefs={columnDefs}
        rowData={data}
        defaultColDef={defaultColDef}
        shouldScrollToRightOnLoad={false}
        gridProps={{
          headerHeight: COLUMN_HEADER_HEIGHT,
          rowHeight: ROW_HEIGHT,
          ...gridProps,
        }}
        refreshCellsCount={refreshCellsCount}
        hideGridOnZeroData={hideGridOnZeroData}
        onHorizontalScroll={handleHorizontalScroll}
        onReady={onReady}
        resetDataSourceCount={resetDataSourceCount}
        setDataSourceWithFilters={setDataSourceWithFilters}
      />
    </div>
  );
};

export default AgGridTable;
