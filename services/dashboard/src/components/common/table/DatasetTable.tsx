import React, { FC } from 'react';
import {
  CellDoubleClickedEvent,
  CellEditRequestEvent,
  ColDef,
  ColumnMovedEvent,
  ColumnVisibleEvent,
  FillEndEvent,
  IServerSideDatasource,
} from 'ag-grid-community';
import { AgGridReact } from 'ag-grid-react';
import { MapAny } from 'types/commonTypes';
import Table from 'components/common/table';

interface DatasetTableProps {
  tableRef?: React.RefObject<AgGridReact>;
  columns: MapAny[];
  serverSideDatasource?: IServerSideDatasource;
  columnConfig?: ColDef;
  totalRows?: number;
  onCellDoubleClicked?: (event: CellDoubleClickedEvent) => void;
  rows?: MapAny[];
  onColumnVisible?: (event: ColumnVisibleEvent) => void;
  onCellEditRequest?: (event: CellEditRequestEvent) => void;
  onFillEnd?: (event: FillEndEvent) => void;
  onDrilldownClick?: (data: MapAny) => void;
  onRowPropertiesClick?: (data: MapAny) => void;
  onColumnMoved?: (event: ColumnMovedEvent) => void;
  columnLevelStats?: MapAny;
}

const DatasetTable: FC<DatasetTableProps> = ({
  columns,
  rows,
  onCellDoubleClicked,
  tableRef,
  totalRows,
  serverSideDatasource,
  columnConfig,
  onColumnVisible,
  onCellEditRequest,
  onFillEnd,
  onDrilldownClick,
  onRowPropertiesClick,
  onColumnMoved,
  columnLevelStats,
}) => {
  return (
    <div id='dataset-table'>
      <Table
        tableRef={tableRef}
        columns={columns}
        rows={rows}
        columnConfig={columnConfig}
        onCellDoubleClicked={onCellDoubleClicked}
        totalRows={totalRows}
        serverSideDatasource={serverSideDatasource}
        onCellEditRequest={onCellEditRequest}
        showStatusBar
        enableCellSelection
        onColumnVisible={onColumnVisible}
        onFillEnd={onFillEnd}
        onDrilldownClick={onDrilldownClick}
        onRowPropertiesClick={onRowPropertiesClick}
        onColumnMoved={onColumnMoved}
        columnLevelStats={columnLevelStats}
      />
    </div>
  );
};

export default DatasetTable;
