import { useCallback, useEffect, useMemo, useRef } from 'react';
import {
  CellDoubleClickedEvent,
  CellStyleModule,
  ClientSideRowModelApiModule,
  ClientSideRowModelModule,
  ColDef,
  ColGroupDef,
  ColumnApiModule,
  ColumnAutoSizeModule,
  ColumnMenuModule,
  ColumnsToolPanelModule,
  ContextMenuModule,
  CsvExportModule,
  FiltersToolPanelModule,
  GridApi,
  GridStateModule,
  GroupCellRendererParams,
  ModuleRegistry,
  PivotModule,
  RenderApiModule,
  RowApiModule,
  RowGroupingPanelModule,
  RowStyleModule,
  ScrollApiModule,
  ValidationModule,
} from 'ag-grid-enterprise';
import { AgGridReact } from 'ag-grid-react';
import { PERIODICITY_TYPES } from 'constants/date.constants';
import { ROUTES_PATH } from 'constants/routeConfig';
import PivotCell from 'modules/widgets/Pivot/components/PivotCell';
import PivotColGroupHeader from 'modules/widgets/Pivot/components/PivotColGroupHeader';
import PivotConfigDropdown from 'modules/widgets/Pivot/components/PivotConfigDropdown';
import PivotRowTitle from 'modules/widgets/Pivot/components/PivotRowTitle';
import PinnedColHeader from 'modules/widgets/Pivot/PinnedColHeader';
import {
  COL_MIN_WIDTH,
  GRAND_ROW_TOTAL_POSITION,
  PINNED_COL_WIDTH,
  PINNED_DIRECTION,
  PIVOT_GRID_OPTIONS,
  PIVOT_GROUP_HEADER_HEIGHT,
  PIVOT_HEADER_HEIGHT,
  PIVOT_TABLE_THEME_PARAMS,
} from 'modules/widgets/Pivot/pivot.constants';
import {
  ColumnFilterConfig,
  ParentFilters,
  PivotContext,
  UNTAGGED_TAGS_FRONTEND_MAPPING,
} from 'modules/widgets/Pivot/pivot.types';
import {
  concatTagFilters,
  getColumnLevelFilters,
  getFilterContext,
  getPivotColDefs,
  getPivotColumns,
  getPivotData,
  getRowLevelFilters,
  getTopNode,
  getWidgetMappingDatasets,
  shouldAllowExpandingRow,
} from 'modules/widgets/Pivot/pivot.utils';
import { getDefaultFilterByDatasetId } from 'modules/widgets/widgets.utils';
import { useRouter } from 'next/navigation';
import { WIDGET_TYPES, WidgetDataResponseType, WidgetInstanceType } from 'types/api/widgets.types';
import { MapAny, OptionsType } from 'types/commonTypes';
import { myTheme } from 'components/common/table/table.constants';
import { getDataTableTheme } from 'components/common/table/table.utils';
import { CONDITION_OPERATOR_TYPE } from 'components/filter/filters.constants';

ModuleRegistry.registerModules([CellStyleModule]);
ModuleRegistry.registerModules([
  ClientSideRowModelModule,
  ColumnsToolPanelModule,
  ColumnMenuModule,
  ContextMenuModule,
  ScrollApiModule,
  RenderApiModule,
  PivotModule,
  ColumnApiModule,
  ClientSideRowModelApiModule,
  FiltersToolPanelModule,
  RowGroupingPanelModule,
  CellStyleModule,
  RowStyleModule,
  RowApiModule,
  ValidationModule,
  GridStateModule,
  ColumnAutoSizeModule,
  CsvExportModule,
]);

type StackedPivotProps = {
  widgetInstanceDetails: Extract<WidgetInstanceType, { widget_type: WIDGET_TYPES.PIVOT_TABLE }>;
  widgetData: WidgetDataResponseType;
  groupWidgetsOptions: OptionsType[];
  onWidgetChange: (widgetId: string) => void;
  currentWidgetSelectedFilter: MapAny;
  periodicity: PERIODICITY_TYPES;
  activeWidget: string;
  handleWidgetHeightChange: (height: number, isSingleHeader: boolean) => void;
  defaultCurrency: string;
};

const StackedPivot = ({
  widgetInstanceDetails,
  widgetData,
  groupWidgetsOptions,
  onWidgetChange,
  currentWidgetSelectedFilter,
  periodicity,
  activeWidget,
  handleWidgetHeightChange,
  defaultCurrency,
}: StackedPivotProps) => {
  const router = useRouter();
  const gridApi = useRef<GridApi | null>(null);
  const customTheme = useMemo(() => getDataTableTheme({ ...PIVOT_TABLE_THEME_PARAMS, ...{} }), []);
  const { title, display_config } = widgetInstanceDetails;
  const gridContainerRef = useRef<HTMLDivElement>(null);

  const handleExportAgGridData = () => {
    gridApi.current?.exportDataAsCsv({ fileName: title, allColumns: true });
  };

  const handleExpandAll = useCallback(() => {
    if (gridApi.current) {
      gridApi.current?.forEachNode((node) => {
        if (node?.group && shouldAllowExpandingRow(node)) {
          node?.setExpanded(true);
        }
      });
    }
  }, []);

  const handleCollapseAll = useCallback(() => {
    if (gridApi.current) {
      gridApi.current?.collapseAll();
    }
  }, []);

  const { colDef, rowData, columnContextMapping } = useMemo(() => {
    const pivotCols = getPivotColumns(widgetInstanceDetails, widgetData);
    const { coldefs, columnContextMapping } = getPivotColDefs(pivotCols);

    return {
      colDef: coldefs,
      rowData: getPivotData(pivotCols, widgetData, periodicity),
      columnContextMapping,
    };
  }, [widgetInstanceDetails, widgetData, periodicity]);

  const pivotContext: PivotContext = useMemo(
    () => ({
      filterContext: getFilterContext(widgetInstanceDetails),
      widgetMappingDatasets: getWidgetMappingDatasets(widgetInstanceDetails),
      columnContextMapping,
    }),
    [widgetInstanceDetails, columnContextMapping],
  );

  const isSingleHeader = useMemo(() => colDef.filter((col) => 'aggFunc' in col).length === 1, [colDef]);

  const defaultColDef = useMemo<ColDef>(
    () => ({
      flex: 1,
      minWidth: COL_MIN_WIDTH,
      enableValue: true,
      enableRowGroup: true,
      enablePivot: true,
      resizable: false,
      cellRenderer: ({ valueFormatted, node, api, column }: GroupCellRendererParams) => {
        return (
          <PivotCell
            value={valueFormatted ?? ''}
            column={column}
            api={api}
            node={node}
            currency={defaultCurrency ?? widgetData?.currency}
            maxGroupingLevel={colDef?.filter((col) => col.rowGroup).length - 1}
            showPercentage={display_config?.show_percentages}
          />
        );
      },
    }),
    [widgetInstanceDetails, display_config, colDef, widgetData],
  );

  const autoGroupColumnDef = useMemo<ColDef>(
    () => ({
      minWidth: PINNED_COL_WIDTH,
      resizable: true,
      pinned: 'left',
      lockPinned: true,
      lockPosition: 'left',
      headerComponent: PinnedColHeader,
      suppressMovable: true,
      headerComponentParams: {
        title,
        isSingleHeader,
        groupWidgetsOptions,
        onWidgetChange,
        widgetType: WIDGET_TYPES.PIVOT_TABLE,
        activeWidget,
        isPortalNeeded: true,
        handleCollapseAll,
        handleExpandAll,
      },
      cellRenderer: (props: GroupCellRendererParams) => {
        return (
          <PivotRowTitle
            node={props?.node}
            value={props?.value}
            maxGroupingLevel={colDef?.filter((col) => col?.rowGroup)?.length - 1}
            displayConfig={display_config}
          />
        );
      },
      cellRendererParams: {
        suppressCount: true,
        suppressPadding: true,
      },
    }),
    [
      widgetInstanceDetails,
      isSingleHeader,
      colDef,
      groupWidgetsOptions,
      onWidgetChange,
      title,
      display_config,
      activeWidget,
      handleExportAgGridData,
      handleCollapseAll,
      handleExpandAll,
    ],
  );

  const processPivotResultColGroupDef = useMemo(() => {
    return (colGroupDef: ColGroupDef) => {
      colGroupDef.headerGroupComponent = PivotColGroupHeader;
      colGroupDef.headerGroupComponentParams = {
        isSingleHeader,
      };
    };
  }, [isSingleHeader]);

  const mergeFilters = (currentFilters: ParentFilters, defaultFilters: ParentFilters) => {
    const mergedFilters: ParentFilters = {};

    Object?.keys({ ...currentFilters, ...defaultFilters })?.forEach((key) => {
      const currentValues = currentFilters[key]?.values || [];
      const defaultValues = defaultFilters[key]?.values || [];

      const updatedCurrentFilter = currentFilters[key] ? { ...currentFilters[key] } : undefined;

      if (updatedCurrentFilter?.targets) {
        delete updatedCurrentFilter?.targets;
      }

      if (updatedCurrentFilter && defaultFilters[key]) {
        mergedFilters[key] = {
          ...updatedCurrentFilter,
          values: currentValues?.filter((value: string) => defaultValues?.includes(value)),
        };
      } else {
        mergedFilters[key] = updatedCurrentFilter || defaultFilters[key];
      }
    });

    return mergedFilters;
  };

  const navigateToDataset = (datasetId: string | null, filters: ParentFilters) => {
    const defaultFilters = getDefaultFilterByDatasetId(widgetInstanceDetails?.data_mappings?.mappings, datasetId ?? '');

    const currentColumnName = Object.keys(currentWidgetSelectedFilter)?.[0];

    const targetDatasetIdColumnName = currentWidgetSelectedFilter[currentColumnName]?.targets?.find(
      (item: MapAny) => item?.dataset_id === datasetId,
    )?.column;

    if (currentColumnName !== targetDatasetIdColumnName) {
      const firstKey = Object.keys(currentWidgetSelectedFilter)?.[0];

      currentWidgetSelectedFilter = { [targetDatasetIdColumnName]: currentWidgetSelectedFilter[firstKey] };
    }

    const query = {
      ...mergeFilters(currentWidgetSelectedFilter, defaultFilters),
      ...filters,
    };

    const path = ROUTES_PATH.DATASET.replace(':datasetId', datasetId ?? '');

    router.push(`${path}?filters=${JSON.stringify(query)}`);
  };

  const handleDrilldown = (params: CellDoubleClickedEvent<MapAny[], PivotContext>) => {
    const { node, colDef: currentColDef } = params;

    if (node?.level === -1 || currentColDef?.pinned === PINNED_DIRECTION.LEFT) return;

    // extract the context from the params
    const context: PivotContext = params.context;

    // extract the current mapping ref from the colDef
    let currentRef = currentColDef?.context?.mappingName;

    // if there are more than one mappings, then the current ref is the top node
    if (Object.keys(context?.columnContextMapping).length > 1) {
      currentRef = getTopNode(node)?.key;
    }
    if (!currentRef) return;

    // extract the column filters for the current ref
    const currentRefColumnFilters: ColumnFilterConfig[] = context?.filterContext?.[currentRef];

    if (!currentRefColumnFilters) return;

    // extract the column context mapping for the current ref
    // this mapping holds the mapping of identifiers set in AGGrid against the column context (name, alias)
    const currentRefColumnContextMapping = context?.columnContextMapping[currentRef];

    if (!currentRefColumnContextMapping) return;

    // extract the row level filters for the currently clickedn ode
    const rowLevelFilters = getRowLevelFilters(currentRefColumnFilters, currentRefColumnContextMapping, node);

    // get the pivot columns that the current cell belongs to
    const pivotColumns = params.api.getPivotColumns().map((col) => col.getColDef());

    // get the column level filters for the current cell
    const columnLevelFilters = getColumnLevelFilters(
      currentRefColumnFilters,
      pivotColumns,
      currentRefColumnContextMapping,
      {
        periodicity,
        widgetSelectedFilter: currentWidgetSelectedFilter,
      },
      currentColDef.pivotKeys || [],
    );

    // merge the row level and column level filters
    let widgetFilter: ParentFilters = concatTagFilters({
      ...rowLevelFilters,
      ...columnLevelFilters,
    });

    // extract the dataset id from the context
    const datasetId = context?.widgetMappingDatasets?.[currentRef];

    if (!datasetId) return;

    if (widgetFilter?.tags?.values?.includes(UNTAGGED_TAGS_FRONTEND_MAPPING.UNTAGGED)) {
      widgetFilter = {
        ...widgetFilter,
        tags: {
          ...widgetFilter?.tags,
          type: CONDITION_OPERATOR_TYPE.IS_NULL,
          values: widgetFilter?.tags?.values?.map((item) =>
            item === UNTAGGED_TAGS_FRONTEND_MAPPING.UNTAGGED ? '' : item,
          ),
        },
      };
    }

    // navigate to the dataset
    navigateToDataset(datasetId, widgetFilter);
  };

  const handleScrollToRightEnd = () => {
    if (gridApi?.current) {
      const allColumns = gridApi.current?.getDisplayedCenterColumns();

      if (allColumns?.length > 0) {
        const lastColumn = allColumns[allColumns?.length - 1];

        gridApi.current?.ensureColumnVisible(lastColumn, 'auto');
      }
    }
  };

  const onGridReady = useCallback((params: { api: GridApi }) => {
    gridApi.current = params.api;

    setTimeout(() => {
      handleScrollToRightEnd();
    }, 0);
  }, []);

  useEffect(() => {
    const observer = new ResizeObserver(() => {
      if (gridContainerRef?.current) {
        handleWidgetHeightChange(gridContainerRef.current.clientHeight, isSingleHeader);
      }
    });

    if (gridContainerRef?.current) {
      observer.observe(gridContainerRef?.current);
    }

    return () => {
      observer.disconnect();
    };
  }, []);

  return (
    <div className='h-fit w-full relative pivot group' ref={gridContainerRef}>
      <PivotConfigDropdown handleExportAgGridData={handleExportAgGridData} />
      <AgGridReact
        onGridReady={onGridReady}
        theme={customTheme ?? myTheme}
        domLayout='autoHeight'
        context={pivotContext}
        rowData={rowData}
        columnDefs={colDef}
        defaultColDef={defaultColDef}
        autoGroupColumnDef={autoGroupColumnDef}
        pivotGroupHeaderHeight={isSingleHeader ? 93 : PIVOT_GROUP_HEADER_HEIGHT}
        pivotHeaderHeight={isSingleHeader ? 0 : PIVOT_HEADER_HEIGHT}
        grandTotalRow={display_config?.show_column_aggregations ? GRAND_ROW_TOTAL_POSITION : undefined}
        processPivotResultColGroupDef={processPivotResultColGroupDef}
        onCellDoubleClicked={handleDrilldown}
        {...PIVOT_GRID_OPTIONS}
      />
    </div>
  );
};

export default StackedPivot;
