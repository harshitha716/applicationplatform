import React, { FC, RefObject, useMemo } from 'react';
import { GridApi, IDatasource, IGetRowsParams } from 'ag-grid-community';
import { MapAny } from 'types/commonTypes';
import { WorkspaceDisplayColumns } from 'components/common/agGridTable/agGrid.types';
import AgGridTable from 'components/common/agGridTable/AgGridTable';
import { getColumnDefs } from 'components/common/agGridTable/agGridTable.utils';

interface InfiniteScrollAgGridProps {
  columns: WorkspaceDisplayColumns[];
  getRows: (params: IGetRowsParams, selectedFilters?: MapAny) => void;
  limit: number;
  id: string;
  columnHeaderMap: MapAny;
  resetDataSourceCount?: number;
  selectedFilters?: MapAny;
  wrapperClassName?: string;
}

const InfiniteScrollAgGrid: FC<InfiniteScrollAgGridProps> = ({
  columns,
  getRows,
  limit,
  id,
  columnHeaderMap,
  resetDataSourceCount,
  selectedFilters,
  wrapperClassName,
}) => {
  const onReady = (params: MapAny) => {
    const dataSource: IDatasource = {
      getRows: getRows,
    };

    params.api.setGridOption('datasource', dataSource);
  };

  const setDataSourceWithFilters = (gridApi: RefObject<GridApi>) => {
    const dataSource: IDatasource = {
      getRows: (params) => getRows(params, selectedFilters),
    };

    gridApi.current?.setGridOption('datasource', dataSource);
  };

  const columnDefs = useMemo(() => getColumnDefs(columns), [columns]);

  return (
    <AgGridTable
      columnDefs={columnDefs}
      gridProps={{
        suppressCellFocus: false,
        rowHeight: 50,
        headerHeight: 34,
        noRowsOverlayComponentParams: {
          isLoading: false,
          title: id,
        },
        rowModelType: 'infinite',
        cacheBlockSize: limit,
        maxBlocksInCache: 10,
        infiniteInitialRowCount: 10,
      }}
      headerWrapperClassName='bg-white'
      headerClassName='uppercase f-8-400 text-GRAY_600'
      hideGridOnZeroData={false}
      onReady={onReady}
      columnHeaderMap={columnHeaderMap}
      resetDataSourceCount={resetDataSourceCount}
      setDataSourceWithFilters={setDataSourceWithFilters}
      wrapperClassName={wrapperClassName}
    />
  );
};

export default InfiniteScrollAgGrid;
