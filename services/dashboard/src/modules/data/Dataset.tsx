import React, { FC, useEffect, useMemo, useRef, useState } from 'react';
import {
  CellEditRequestEvent,
  ColDef,
  ColumnMovedEvent,
  FillEndEvent,
  IServerSideDatasource,
  IServerSideGetRowsParams,
  IServerSideGetRowsRequest,
} from 'ag-grid-community';
import { AgGridReact } from 'ag-grid-react';
import {
  useGetDatasetFilterConfigQuery,
  useLazyGetActionStatusQuery,
  useLazyGetDatasetDataQuery,
  useUpdateDatasetDataMutation,
} from 'apis/dataset';
import { ZAMP_LOGO_LOADER } from 'constants/lottie/zamp-logo-loader';
import { ROUTES_PATH } from 'constants/routeConfig';
import { useOnClickOutside } from 'hooks';
import { useAppDispatch, useAppSelector } from 'hooks/toolkit';
import usePolling from 'hooks/usePolling';
import DatasetHistory from 'modules/data/components/datasetHistory/index';
import ExportDataset from 'modules/data/components/exportDataset';
import ImportDataset from 'modules/data/components/importDataset/index';
import TableSchemaAlignmentStatus from 'modules/data/components/importDataset/TableSchemaAlignmentStatus';
import { LOADER_STATUS } from 'modules/data/data.types';
import {
  formatColumnLevelStats,
  formatColumns,
  formatDrilldownFilters,
  formatUrlFilters,
  getColumnOrderingVisibilityForCurrentDataset,
  getEncodedRequestWithAggregations,
  getFilters,
} from 'modules/data/data.utils';
import Notification from 'modules/data/Notification';
import RowPropertiesSideDrawer from 'modules/data/RowProperties';
import RulesListingSideDrawer from 'modules/data/RulesListing';
import { PAGE_CURRENCY_OPTIONS } from 'modules/page/pages.constants';
import SingleSelectFilter from 'modules/widgets/components/SingleSelectFilter';
import { useSearchParams } from 'next/navigation';
import { useRouter } from 'next/router';
import { RootState } from 'store';
import { addBreadcrumb } from 'store/slices/layout-configs';
import {
  DatasetActionStatusResponseType,
  DatasetDataResponseType,
  DatasetUpdateResponseType,
} from 'types/api/dataset.types';
import { MapAny } from 'types/commonTypes';
import { FilterModelType, LogicalOperatorType } from 'types/components/table.type';
import { checkIsObjectEmpty, cn, snakeCaseToSentenceCase } from 'utils/common';
import { getFromLocalStorage, LOCAL_STORAGE_KEYS, setToLocalStorage } from 'utils/localstorage';
import CustomHeader from 'components/common/table/CustomHeader';
import DatasetTable from 'components/common/table/DatasetTable';
import DisplayOptions from 'components/common/table/DisplayOptions';
import { getEncodedRequest } from 'components/common/table/table.utils';
import { toast } from 'components/common/toast/Toast';
import { TOAST_MESSAGES } from 'components/common/toast/toast.constants';
import CommonWrapper from 'components/commonWrapper';
import { SkeletonTypes } from 'components/commonWrapper/commonWrapper.types';
import DynamicLottiePlayer from 'components/DynamicLottiePlayer';
import { FILTER_TYPES } from 'components/filter/filter.types';
import FiltersWrapper from 'components/filter/filterMenu/FiltersWrapper';
import { CONDITION_OPERATOR_TYPE } from 'components/filter/filters.constants';
import { filtersContextActions, useFiltersContextStore, withFiltersContext } from 'components/filter/filters.context';
type DatasetByIdProps = {
  id: string;
  drilldownFilters?: FilterModelType;
};

const DatasetById: FC<DatasetByIdProps> = ({ id, drilldownFilters }) => {
  const filters = useSearchParams().get('filters');
  const currency = useSearchParams().get('currency') ?? 'local';
  const appDispatch = useAppDispatch();
  const breadcrumbStack = useAppSelector((state: RootState) => state.layoutConfig.breadcrumbStack);

  const {
    data: filterConfigData,
    refetch: refetchFilterConfig,
    isFetching,
    isError,
    isUninitialized,
  } = useGetDatasetFilterConfigQuery(
    {
      datasetId: id as string,
    },
    {
      skip: !id,
      refetchOnMountOrArgChange: true,
    },
  );
  const showFileImports = filterConfigData?.config?.is_file_import_enabled;
  const [updateDatasetData] = useUpdateDatasetDataMutation();
  const [getActionStatus] = useLazyGetActionStatusQuery();
  const [columns, setColumns] = useState<ColDef[]>([]);
  const [isPolling, setIsPolling] = useState<boolean>(false);
  const [totalRows, setTotalRows] = useState<number>(0);
  const [columnId, setColumnId] = useState<string>('');
  const [isRulesListingSideDrawerOpen, setIsRulesListingSideDrawerOpen] = useState(false);
  const [rowPropertiesData, setRowPropertiesData] = useState<MapAny>();
  const [exportsDatasetQuery, setExportsDatasetQuery] = useState<string>('');
  const [datasetTitle, setDatasetTitle] = useState<string>('');
  const [fxCurrency, setFxCurrency] = useState<string[]>([currency]);
  const [initiatedActionIds, setInitiatedActionIds] = useState<string[]>([]);
  const [isNoRowsOverlayVisible, setIsNoRowsOverlayVisible] = useState<boolean>(false);
  const [cachedDatasetData, setCachedDatasetData] = useState<DatasetDataResponseType>();
  const [columnLevelStats, setColumnLevelStats] = useState<MapAny>();
  const [hiddenColumnFilters, setHiddenColumnFilters] = useState<MapAny>();

  const firstLoadDone = useRef(false); // Track if first load is done

  const [showAiTransformationStatus, setShowAiTransformationStatus] = useState<{
    open: boolean;
    status: string;
    title: string;
    description: string;
  }>({
    open: false,
    status: LOADER_STATUS.LOADING,
    title: '',
    description: '',
  });
  const { startPolling } = usePolling();
  const [getDatasetData, { data: datasetData }] = useLazyGetDatasetDataQuery();

  const {
    dispatch,
    state: { selectedFilters, filtersConfig },
  } = useFiltersContextStore();

  const serverSideDatasource: IServerSideDatasource = useMemo(() => {
    return {
      getRows: (parameters: IServerSideGetRowsParams): void => {
        const queryConfig = getEncodedRequest(
          parameters.request,
          fxCurrency?.[0],
          false,
          false,
          false,
          hiddenColumnFilters,
        );

        const filterModel = parameters?.request?.filterModel;
        const isDefaultFilters = checkIsObjectEmpty(filterModel ?? {})
          ? false
          : Object.values(filterModel ?? {}).every((filter) => filter?.isDefault);

        removeCellFocus();
        setExportsDatasetQuery(queryConfig);
        if (!firstLoadDone.current || isDefaultFilters) {
          // Use Cached Data for First Load
          firstLoadDone.current = true; // Mark first load as done
          if (drilldownFilters?.conditions === null) {
            setIsNoRowsOverlayVisible(true);
            parameters.success({
              rowData: [],
              rowCount: 0,
            });
          } else if (!checkIsObjectEmpty(cachedDatasetData)) {
            const totalCount = cachedDatasetData?.data?.total_count ?? 0;

            setIsNoRowsOverlayVisible(totalCount === 0);
            parameters.success({
              rowData: cachedDatasetData?.data?.rows ?? [],
              ...(parameters.request.startRow === 0 ? { rowCount: totalCount } : {}),
            });
          }
        } else {
          getDatasetData({
            datasetId: id as string,
            query_config: queryConfig,
          })
            .unwrap()
            .then((response) => {
              if (parameters.request.startRow === 0) {
                setDatasetTitle(response?.title);
                setTotalRows(response?.data?.total_count);
                setIsNoRowsOverlayVisible(response?.data?.total_count === 0);
                dispatch({
                  type: filtersContextActions.SET_TOTAL_ROWS,
                  payload: { totalRows: response?.data?.total_count },
                });
              }
              parameters.success({
                rowData: response?.data?.rows,
                ...(parameters.request.startRow === 0 ? { rowCount: response?.data?.total_count } : {}),
              });
            })
            .catch(() => {
              parameters.fail();
            });
        }
      },
    };
  }, [getDatasetData, id, fxCurrency, cachedDatasetData, drilldownFilters]);

  const router = useRouter();
  const tableRef = useRef<AgGridReact>(null);
  const datasetTableRef = useRef<HTMLDivElement>(null);

  const removeCellFocus = () => {
    tableRef.current?.api?.clearCellSelection();
    tableRef.current?.api?.clearFocusedCell();
  };

  const handleSuccessfulUpdate = (data: DatasetUpdateResponseType, showPolling = true) => {
    if (showPolling) setIsPolling(true);
    setInitiatedActionIds((prev) => [...prev, data.action_id]);
    startPolling({
      fn: () =>
        getActionStatus({ datasetId: id as string, params: { action_ids: [...initiatedActionIds, data.action_id] } }),
      validate: (data: DatasetActionStatusResponseType[]) => {
        return data.filter((item) => !item.is_completed).length === 0;
      },
      interval: 30000,
      maxAttempts: 50,
    }).then(() => {
      setIsPolling(false);
      toast.success(TOAST_MESSAGES.SUCCESS_TAGGING_COMPLETED);
      tableRef.current?.api?.refreshServerSide();
      refetchFilterConfig();
    });
  };

  const updateApi = ({
    rowId,
    field,
    newValue,
    operator = CONDITION_OPERATOR_TYPE.EQUAL,
  }: {
    rowId: string | string[];
    field: string;
    newValue: string;
    operator?: CONDITION_OPERATOR_TYPE;
  }) => {
    updateDatasetData({
      datasetId: id as string,
      data: {
        filters: {
          logical_operator: LogicalOperatorType.OperatorLogicalAnd,
          conditions: [
            {
              column: '_zamp_id',
              value: rowId,
              operator: operator,
            },
          ],
        },
        update: {
          column: field as string,
          value: newValue,
        },
      },
    })
      .unwrap()
      .then((response) => handleSuccessfulUpdate(response, false));
  };

  const onCellEditRequest = (event: CellEditRequestEvent) => {
    const { colDef, newValue, data, source, node } = event;
    const { field } = colDef;
    const updatedRow = { ...event.data, [field as string]: newValue };

    // Optimistic update
    node.setData(updatedRow);

    if (source === 'edit') updateApi({ rowId: data?._zamp_id as string, field: field as string, newValue });
  };

  const onFillEnd = (event: FillEndEvent) => {
    const { finalRange } = event;
    const { startRow, endRow, startColumn } = finalRange;

    const startIndex = startRow?.rowIndex as number;
    const endIndex = endRow?.rowIndex as number;
    const field = startColumn?.getColId();
    const rowIds: string[] = [];
    let newValue = '';
    let loopStartIndex = startIndex;
    let loopEndIndex = endIndex;

    if (startIndex > endIndex) {
      loopStartIndex = endIndex;
      loopEndIndex = startIndex;
    }
    for (let i = loopStartIndex; i <= loopEndIndex; i++) {
      const row = tableRef.current?.api?.getDisplayedRowAtIndex(i);

      rowIds.push(row?.data?._zamp_id as string);
      if (i === startIndex) {
        newValue = row?.data?.[field as string] as string;
      }
    }

    updateApi({
      rowId: rowIds,
      field,
      newValue,
      operator: CONDITION_OPERATOR_TYPE.IN,
    });
  };

  const handleDrilldownClick = (data: MapAny) => {
    appDispatch(addBreadcrumb('Drilldown'));
    router.push(ROUTES_PATH.DRILLDOWN.replace(':datasetId', id as string).replace(':rowId', data?._zamp_id as string));
  };

  const handleRowPropertiesClick = (data: MapAny) => {
    setRowPropertiesData(data);
  };

  const handleRulesListingSideDrawerOpen = (columnId: string) => {
    setIsRulesListingSideDrawerOpen(true);
    setColumnId(columnId);
  };

  useEffect(() => {
    if (filterConfigData?.data?.length && !isFetching && !isUninitialized) {
      const columns = formatColumns(
        filterConfigData?.data,
        false,
        id as string,
        handleSuccessfulUpdate,
        tableRef,
        handleRulesListingSideDrawerOpen,
      );

      if (columns?.length > 0) {
        setColumns(columns);
        dispatch({
          type: filtersContextActions.SET_FILTERS_CONFIG,
          payload: {
            filtersConfig: filterConfigData?.data
              ?.filter((item) => !item?.metadata?.is_hidden)
              ?.map((column) => ({
                key: column.column,
                label: column.alias ?? snakeCaseToSentenceCase(column?.column),
                values: column.options,
                type: column.type,
              })),
          },
        });
        if (filters) {
          firstLoadDone.current = false;
          dispatch({
            type: filtersContextActions.INITIALIZE_DEFAULT_FILTERS,
            payload: { selectedFilters: getFilters(filters, filterConfigData.data) ?? {} },
          });
        }

        if (drilldownFilters) {
          firstLoadDone.current = false;
          const { selectedDrilldownFilters, hiddenDrilldownFilters } = formatDrilldownFilters(
            drilldownFilters,
            filterConfigData?.data,
          );

          if (!checkIsObjectEmpty(hiddenDrilldownFilters)) setHiddenColumnFilters(hiddenDrilldownFilters);
          if (!checkIsObjectEmpty(selectedDrilldownFilters))
            dispatch({
              type: filtersContextActions.INITIALIZE_DEFAULT_FILTERS,
              payload: { selectedFilters: selectedDrilldownFilters },
            });
        }
        const amountRangeColumns = columns
          ?.filter((column) => column?.headerComponentParams?.filterType === FILTER_TYPES.AMOUNT_RANGE)
          ?.map((column) => column?.field)
          ?.filter((column) => column !== undefined);

        if (amountRangeColumns?.length > 0) {
          getDatasetData({
            datasetId: id as string,
            query_config: getEncodedRequestWithAggregations(amountRangeColumns),
          })
            .unwrap()
            .then((response) => {
              setColumnLevelStats(formatColumnLevelStats(response?.data?.rows?.[0]));
            })
            .catch(() => {
              setColumnLevelStats(undefined);
            });
        }
      }
    }
  }, [filterConfigData?.data, filters, id, drilldownFilters, isFetching, isUninitialized]);

  useEffect(() => {
    tableRef.current?.api?.setFilterModel(selectedFilters);
  }, [selectedFilters, fxCurrency]);

  useEffect(() => {
    if (isNoRowsOverlayVisible) {
      tableRef.current?.api?.showNoRowsOverlay();
    } else {
      tableRef.current?.api?.hideOverlay();
    }
  }, [isNoRowsOverlayVisible]);

  const handleFilterChange = (value: string[]) => {
    setFxCurrency(value);
  };

  useEffect(() => {
    if (datasetTitle && (breadcrumbStack?.length === 0 || !breadcrumbStack?.includes(datasetTitle))) {
      appDispatch(addBreadcrumb(datasetTitle));
    }
  }, [datasetTitle, breadcrumbStack]);

  const handleRefetchDataset = () => {
    getDatasetData({
      datasetId: id as string,
      query_config: exportsDatasetQuery,
    });
  };

  const handleColumnMoved = (event: ColumnMovedEvent) => {
    const columnOrderingFromLocalStorage = getColumnOrderingVisibilityForCurrentDataset(id);
    const latestColumns = event?.api?.getColumns() ?? [];
    const { column, toIndex = 0 } = event;

    if (!column) return;
    const columnOrderingVisibility: { colId: string; isVisible: boolean }[] = columnOrderingFromLocalStorage?.length
      ? columnOrderingFromLocalStorage
      : latestColumns.map((column) => ({
          colId: column.getColId(),
          isVisible: column.isVisible(),
        }));

    const movedColumn = columnOrderingVisibility.find((item) => item.colId === column?.getColId()) ?? {};
    const fromIndex = columnOrderingVisibility.findIndex((item) => item.colId === column?.getColId());

    if (fromIndex === toIndex) return;
    let finalList: { colId?: string; isVisible?: boolean }[] = [];

    if (fromIndex < toIndex) {
      const zeroToOldIndex = columnOrderingVisibility.slice(0, fromIndex) ?? [];
      const oldIndexToNewIndex = columnOrderingVisibility.slice(fromIndex + 1, toIndex + 1) ?? [];
      const newIndexToEnd = columnOrderingVisibility.slice(toIndex + 1) ?? [];

      finalList = [...zeroToOldIndex, ...oldIndexToNewIndex, movedColumn, ...newIndexToEnd];
    } else {
      const endToOldIndex = columnOrderingVisibility.slice(fromIndex + 1) ?? [];
      const oldIndexToNewIndex = columnOrderingVisibility.slice(toIndex, fromIndex) ?? [];
      const newIndexToStart = columnOrderingVisibility.slice(0, toIndex) ?? [];

      finalList = [...newIndexToStart, movedColumn, ...oldIndexToNewIndex, ...endToOldIndex];
    }
    const currentColumnOrderingVisibility = JSON.parse(
      getFromLocalStorage(LOCAL_STORAGE_KEYS.COLUMN_ORDERING_VISIBILITY) ?? '{}',
    );

    setToLocalStorage(
      LOCAL_STORAGE_KEYS.COLUMN_ORDERING_VISIBILITY,
      JSON.stringify({ ...currentColumnOrderingVisibility, [id]: finalList }),
    );
  };

  useEffect(() => {
    firstLoadDone.current = false;
    if (drilldownFilters?.conditions === null) return;
    const urlFilters = formatUrlFilters(filters ?? '');
    const queryConfig = getEncodedRequest(
      {} as IServerSideGetRowsRequest,
      fxCurrency?.[0],
      false,
      false,
      false,
      hiddenColumnFilters,
      urlFilters ?? drilldownFilters,
    );

    getDatasetData({
      datasetId: id as string,
      query_config: queryConfig,
    })
      .unwrap()
      .then((response) => {
        setDatasetTitle(response?.title);
        setTotalRows(response?.data?.total_count);
        setCachedDatasetData(response);
        dispatch({
          type: filtersContextActions.SET_TOTAL_ROWS,
          payload: { totalRows: response?.data?.total_count },
        });
      });
  }, [filters, drilldownFilters, id]);

  useOnClickOutside(datasetTableRef, removeCellFocus);

  return (
    <>
      <CommonWrapper
        className={cn('h-full', {
          'flex flex-col items-center justify-center': isFetching,
        })}
        isLoading={isFetching}
        isError={isError}
        skeletonType={SkeletonTypes.CUSTOM}
        refetchFunction={refetchFilterConfig}
        loader={
          <div className='flex justify-center items-center h-[calc(100vh-200px)] w-full z-50 bg-white'>
            <DynamicLottiePlayer
              src={ZAMP_LOGO_LOADER}
              className='lottie-player h-[140px]'
              autoplay
              loop
              keepLastFrame
            />
          </div>
        }
      >
        <div className='flex items-center justify-between pr-8 z-1000'>
          <div className='flex items-center py-3'>
            <FiltersWrapper label='Filter' filterConfig={filtersConfig ?? []} />
          </div>
          <div className='relative flex items-center gap-2.5'>
            <Notification isPolling={isPolling} />
            <TableSchemaAlignmentStatus
              showAiTransformationStatus={showAiTransformationStatus}
              setShowAiTransformationStatus={setShowAiTransformationStatus}
            />
            <ExportDataset
              query={exportsDatasetQuery}
              datasetId={id as string}
              hasFilters={!!Object.keys(selectedFilters)?.length}
            />
            {showFileImports && (
              <ImportDataset
                onRefetch={handleRefetchDataset}
                setShowAiTransformationStatus={setShowAiTransformationStatus}
              />
            )}
            <DatasetHistory />
            <DisplayOptions tableRef={tableRef} datasetId={id as string} />
            <div className='flex items-center gap-2'>
              <div className='border-r border-GRAY_400 h-7'></div>
              <SingleSelectFilter
                onFilterChange={handleFilterChange}
                value={fxCurrency}
                filterKey='fx_currency'
                label='Currency'
                options={PAGE_CURRENCY_OPTIONS}
              />
            </div>
          </div>
        </div>

        <div className='z-10 w-full h-full' ref={datasetTableRef}>
          <DatasetTable
            tableRef={tableRef}
            columns={columns}
            serverSideDatasource={serverSideDatasource}
            columnConfig={{ enableRowGroup: true, enableValue: true, headerComponent: CustomHeader }}
            totalRows={totalRows}
            onCellEditRequest={onCellEditRequest}
            onFillEnd={onFillEnd}
            onRowPropertiesClick={handleRowPropertiesClick}
            onColumnMoved={handleColumnMoved}
            columnLevelStats={columnLevelStats}
            {...(datasetData?.data?.config?.is_drilldown_enabled ? { onDrilldownClick: handleDrilldownClick } : {})}
          />
        </div>
      </CommonWrapper>
      {isRulesListingSideDrawerOpen && (
        <RulesListingSideDrawer
          column={columnId}
          onClose={() => setIsRulesListingSideDrawerOpen(false)}
          datasetId={id as string}
          handleSuccessfulUpdate={handleSuccessfulUpdate}
        />
      )}
      {rowPropertiesData && (
        <RowPropertiesSideDrawer
          data={rowPropertiesData}
          onClose={() => setRowPropertiesData(undefined)}
          datasetId={id as string}
          isDrillDownEnabled={datasetData?.data?.config?.is_drilldown_enabled}
          columns={columns}
        />
      )}
    </>
  );
};

export default withFiltersContext(DatasetById);
