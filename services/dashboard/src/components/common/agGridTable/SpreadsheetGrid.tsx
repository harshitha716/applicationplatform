import React, { Dispatch, FC, Ref, RefObject, SetStateAction, useCallback, useEffect, useRef, useState } from 'react';
import { Column, GridApi } from 'ag-grid-community';
import { AgGridReact } from 'ag-grid-react';
import { debounce } from 'hooks';
import { defaultFn, MapAny } from 'types/commonTypes';
import 'ag-grid-enterprise';
import 'ag-grid-community/styles/ag-grid.css';

interface SpreadsheetGridProps {
  tableWrapperRef?: Ref<HTMLDivElement>;
  scrollToRowId?: string;
  dividerRows?: MapAny[];
  columnDefs: MapAny[];
  rowData: MapAny[];
  defaultColDef?: MapAny;
  gridProps?: MapAny;
  shouldScrollToRightOnLoad?: boolean;
  scrollToRightEnd?: boolean;
  refreshCellsCount?: number;
  scrollToCellData?: MapAny;
  suppressCellFocus?: boolean;
  onScrollToCell?: () => void;
  setFirstColumnIdInViewport?: Dispatch<SetStateAction<string>>;
  onScrollToLeftEnd?: () => void;
  onScrollToRightEnd?: () => void;
  onHorizontalScroll?: () => void;
  onCellClicked?: (event: any) => void;
  onCellMouseOver?: (event: any) => void;
  onCellMouseOut?: (event: any) => void;
  isLoading?: boolean;
  onReady?: (params: any) => void;
  hideGridOnZeroData?: boolean;
  onSortChanged?: (event: any) => void;
  resetDataSourceCount?: number;
  setDataSourceWithFilters?: (gridApi: RefObject<GridApi>) => void;
}

const SpreadsheetGrid: FC<SpreadsheetGridProps> = ({
  tableWrapperRef,
  scrollToRowId,
  dividerRows,
  columnDefs,
  rowData,
  defaultColDef,
  gridProps,
  shouldScrollToRightOnLoad = true,
  scrollToRightEnd = false,
  refreshCellsCount,
  scrollToCellData,
  suppressCellFocus = false,
  onScrollToCell,
  setFirstColumnIdInViewport,
  onScrollToLeftEnd,
  onScrollToRightEnd,
  onHorizontalScroll = defaultFn,
  onCellClicked,
  onCellMouseOver,
  onCellMouseOut,
  onReady,
  onSortChanged,
  hideGridOnZeroData = true,
  resetDataSourceCount,
  setDataSourceWithFilters,
}) => {
  const gridApi = useRef<GridApi | null>(null);

  const [lastViewedColumn, setLastViewedColumn] = useState<Column<any> | undefined>(undefined);
  let isEnsureColumnVisibleScroll = false;

  const handleRefreshCells = () => {
    gridApi?.current?.refreshCells({ force: true });
  };

  //--------------------  Scroll to the last viewed column --------------------
  useEffect(() => {
    if (lastViewedColumn) gridApi?.current?.ensureColumnVisible(lastViewedColumn, 'middle');
  }, [columnDefs?.length]);
  //--------------------  Scroll to the last viewed column --------------------

  //--------------------  Scroll to the specified row index --------------------
  useEffect(() => {
    if (!!scrollToRowId && gridApi?.current) {
      const row = dividerRows?.find((rowData: MapAny) => rowData?.row?.reference_id === scrollToRowId);
      const rowIndex = row?.index;

      if (rowIndex === undefined) return;

      setTimeout(() => {
        gridApi?.current?.ensureIndexVisible?.(rowIndex, 'top');
      }, 0);
    }
  }, [scrollToRowId]);
  //--------------------  Scroll to the specified row index --------------------

  //--------------------  Horizontal scroll event handler --------------------

  const findFirstColumnIdInViewPort = () => {
    const horizontalPixelRange = gridApi?.current?.getHorizontalPixelRange?.();
    const { left = 0, right = 0 } = horizontalPixelRange ?? {};

    const allColumns = gridApi?.current?.getDisplayedCenterColumns?.() ?? [];

    let columnAtLeft = null;
    let minDistance = Infinity;
    let accumulatedWidth = 0;

    for (const column of allColumns) {
      const columnWidth = column?.getActualWidth?.();
      const columnRightPixel = accumulatedWidth + columnWidth;

      // Check if the column is within the visible range
      if (columnRightPixel >= left && accumulatedWidth <= right) {
        // Calculate the distance of the column's left edge from the left boundary of the viewport
        const distanceToLeft = Math.abs(accumulatedWidth - left);

        // Update the leftmost column if this column is closer to the left boundary
        if (distanceToLeft < minDistance) {
          minDistance = distanceToLeft;
          columnAtLeft = column;
        }
      }

      // Update the accumulated width for the next iteration
      accumulatedWidth += columnWidth;
    }

    return columnAtLeft?.getColDef?.()?.field ?? '';
  };

  const handleHorizontalScroll = (left: number) => {
    if (isEnsureColumnVisibleScroll) {
      isEnsureColumnVisibleScroll = false;

      return;
    }

    if (left === 0) {
      const lastViewedColumn =
        columnDefs?.filter((col) => !col?.pinned && !(col?.field === 'skeleton'))?.[0] ?? undefined;

      setLastViewedColumn(lastViewedColumn?.field);
      onScrollToLeftEnd?.();
    }
    setFirstColumnIdInViewport?.(findFirstColumnIdInViewPort());
  };

  //--------------------  Horizontal scroll event handler --------------------

  const handleScroll = (event: any) => {
    const { left, direction } = event;

    if (direction === 'horizontal') {
      handleHorizontalScroll(left);
      debounce(onHorizontalScroll, 50)();
    }
  };

  //--------------------  Scroll to the rightmost column by default --------------------
  const handleScrollToRightEnd = () => {
    isEnsureColumnVisibleScroll = true;
    gridApi?.current?.ensureColumnVisible?.(columnDefs?.[columnDefs?.length - 1]?.field, 'end');
  };
  //--------------------  Scroll to the rightmost column by default --------------------

  useEffect(() => {
    if (!!columnDefs?.length && scrollToRightEnd) {
      handleScrollToRightEnd();
      onScrollToRightEnd?.();
    }
  }, [columnDefs?.length]);

  //--------------------  Scroll to cell by rowId, columnId --------------------
  const handleScrollToCell = () => {
    if (!!scrollToCellData?.columnId && !!scrollToCellData?.rowId) {
      const rowIndex = rowData?.findIndex((row) => row?.id === scrollToCellData?.rowId);

      if (rowIndex >= 0) {
        setTimeout(() => {
          gridApi?.current?.ensureIndexVisible?.(rowIndex, 'top');
        }, 0);
      }

      const column = columnDefs?.find((col) => col?.field === scrollToCellData?.columnId);

      if (column) {
        setTimeout(() => {
          gridApi?.current?.ensureColumnVisible?.(column?.field, 'middle');
        }, 0);
      }

      onScrollToCell?.();
    }
  };
  //--------------------  Scroll to cell by rowId, columnId --------------------

  //--------------------  Cell copy handlers --------------------

  const getCellValue = (event: any): string => {
    // Get the rendered cell element
    const cellElement = event?.event?.target as HTMLElement;

    // Extract text content from the cell
    let cellValue = cellElement?.innerText || cellElement?.textContent || '';

    // Trim whitespace
    cellValue = cellValue?.trim();

    // If cell is empty, fall back to the raw value
    if (!cellValue) {
      cellValue = String(event?.value || '');
    }

    return cellValue;
  };

  const onCellKeyDown = useCallback((event: any) => {
    // Check if the pressed keys are 'Ctrl + C' (Windows/Linux) or 'Cmd + C' (Mac)
    if ((event?.event?.ctrlKey || event?.event?.metaKey) && event?.event?.key === 'c') {
      const cellValue = getCellValue?.(event);

      // Copy to clipboard
      navigator?.clipboard?.writeText?.(cellValue)?.catch((err) => {
        console.error('Failed to copy cell value: ', err);
      });

      // Prevent default copy behavior
      event?.event?.preventDefault?.();
    }
  }, []);
  //--------------------  Cell copy handlers --------------------

  useEffect(() => {
    if (refreshCellsCount) handleRefreshCells();
  }, [refreshCellsCount]);

  const onGridReady = (params: any) => {
    gridApi.current = params?.api;

    onReady?.(params);

    if (!!scrollToCellData?.columnId && !!scrollToCellData?.rowId) handleScrollToCell();
    else if (shouldScrollToRightOnLoad) {
      setTimeout(() => {
        handleScrollToRightEnd();
      }, 0);
    }
  };

  useEffect(() => {
    if (resetDataSourceCount) setDataSourceWithFilters?.(gridApi);
  }, [resetDataSourceCount]);

  return (columnDefs?.length === 0 || rowData?.length === 0) && hideGridOnZeroData ? null : (
    <div className='sensitive' style={{ height: '100%', width: '100%' }} ref={tableWrapperRef}>
      <AgGridReact
        columnDefs={columnDefs}
        defaultColDef={defaultColDef}
        rowData={rowData}
        suppressHorizontalScroll={false}
        suppressScrollOnNewData={true}
        suppressMovableColumns={true}
        enableCellTextSelection={true}
        suppressCellFocus={suppressCellFocus}
        suppressHeaderFocus={true}
        suppressContextMenu={true}
        onCellMouseOver={onCellMouseOver}
        onCellMouseOut={onCellMouseOut}
        onCellDoubleClicked={onCellClicked}
        onCellKeyDown={onCellKeyDown}
        onGridReady={onGridReady}
        onBodyScroll={handleScroll}
        onSortChanged={onSortChanged}
        {...gridProps}
      />
    </div>
  );
};

export default SpreadsheetGrid;
