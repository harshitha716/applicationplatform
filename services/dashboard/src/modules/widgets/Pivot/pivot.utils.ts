import { ColDef, IRowNode, RowStyle } from 'ag-grid-community';
import { DATE_FORMATS, PERIODICITY_TYPES } from 'constants/date.constants';
import { format, isValid, parse } from 'date-fns';
import PivotColGroupHeader from 'modules/widgets/Pivot/components/PivotColGroupHeader';
import PivotColHeader from 'modules/widgets/Pivot/components/PivotColHeader';
import { GROUPING_COL_NAME_PREFIX, NESTING_LEVEL_INFIX, PIVOT_REF } from 'modules/widgets/Pivot/pivot.constants';
import {
  ColumnFilterConfig,
  PIVOT_DATA_TYPES,
  PivotColumnMetadata,
  UNTAGGED_TAGS,
  UNTAGGED_TAGS_FRONTEND_MAPPING,
} from 'modules/widgets/Pivot/pivot.types';
import { getFormattedDateWithPeriodicity } from 'modules/widgets/widgets.constant';
import { getDateRangeWithPeriodicity } from 'modules/widgets/widgets.utils';
import {
  AGGREGATION_TYPES,
  PivotTableWidgetInstanceType,
  WIDGET_TYPES,
  WidgetDataResponseType,
  WidgetInstanceType,
} from 'types/api/widgets.types';
import { MapAny } from 'types/commonTypes';
import { capitalizeFirstLetter, formatCurrencyValue, snakeCaseToSentenceCase } from 'utils/common';
import { FILTER_TYPES } from 'components/filter/filter.types';
import { CONDITION_OPERATOR_TYPE } from 'components/filter/filters.constants';

export const backendConfig = {
  styleConfig: {
    rowStyles: [
      // {
      //   conditions: [
      //     { level: 0 }, // Condition for level 0
      //   ],
      //   operator: 'AND', // All conditions must be met
      //   style: { backgroundColor: 'red', color: 'blue' },
      // },
      // {
      //   conditions: [{ level: 1 }], // Condition for level 1
      //   operator: 'AND',
      //   style: { backgroundColor: 'blue' },
      // },
    ],
    cellStyles: [
      // {
      //   field: 'Metric',
      //   conditions: [{ equals: 'Opening Balance' }], // Condition for CARD_BRAND equals 'Visa'
      //   operator: 'AND',
      //   style: { backgroundColor: 'yellow', color: 'white' },
      // },
      // {
      //   field: 'Actual',
      //   conditions: [
      //     { greaterThan: 20 }, // Condition for value greater than 20
      //   ],
      //   operator: 'AND',
      //   style: { backgroundColor: 'yellow', color: 'red' },
      // },
      // {
      //   field: 'SubMetric',
      //   conditions: [{ default: true }], // Default condition
      //   operator: 'AND',
      //   style: { border: '1px solid green' },
      // },
    ],
  },
};

export const evaluateConditions = (conditions: MapAny[], groupingLevel: number, value: string, operator = 'AND') => {
  const results = conditions?.map((condition) => {
    if (condition?.level !== undefined && groupingLevel === condition?.level) {
      return true;
    }
    if (condition?.equals !== undefined && value === condition?.equals) {
      return true;
    }
    if (condition?.greaterThan !== undefined && value > condition?.greaterThan) {
      return true;
    }
    if (condition?.default) {
      return true;
    }
    if (condition?.type === 'dateLessThanToday') {
      const today = new Date().toISOString().split('T')[0];

      return new Date(value).toISOString().split('T')[0] < today;
    }

    return false;
  });

  return operator === 'AND' ? results?.every(Boolean) : results?.some(Boolean);
};

export const getDynamicRowStyle = (rowStyles: MapAny[], groupingLevel: number, value: string): RowStyle => {
  for (const rule of rowStyles) {
    if (evaluateConditions(rule?.conditions, groupingLevel, value, rule?.operator)) {
      return rule?.style;
    }
  }

  return {};
};

export const getDynamicCellStyle = (cellStyles: MapAny[], field: string, groupingLevel: number, value: string) => {
  for (const rule of cellStyles) {
    if (rule?.field === field && evaluateConditions(rule?.conditions, groupingLevel, value, rule?.operator)) {
      return rule?.style;
    }
  }

  return {}; // Default cell style
};

export const parseType = (type: PIVOT_DATA_TYPES, value: string | number | boolean, periodicity: PERIODICITY_TYPES) => {
  switch (type) {
    case PIVOT_DATA_TYPES.DATE:
    case PIVOT_DATA_TYPES.TIMESTAMP:
      return getFormattedDateWithPeriodicity(periodicity, value as string);
    case PIVOT_DATA_TYPES.NUMBER:
    case PIVOT_DATA_TYPES.AMOUNT: {
      const number = Number(value);

      return isNaN(number) ? 0 : number;
    }
    case PIVOT_DATA_TYPES.BANK:
    case PIVOT_DATA_TYPES.TAG:
      return value === UNTAGGED_TAGS.UNTAGGED ? UNTAGGED_TAGS_FRONTEND_MAPPING.UNTAGGED : value;
    case PIVOT_DATA_TYPES.COUNTRY:
    case PIVOT_DATA_TYPES.STATUS:
      return value;
    case PIVOT_DATA_TYPES.BOOLEAN:
      return value === 'true' || value === true;
    default:
      return value;
  }
};

export const getGroupingColName = (groupingLevel: number) => {
  return `${GROUPING_COL_NAME_PREFIX}${groupingLevel}`;
};

export const getNestedGroupingColName = (columnName: string, hierarchy: number) => {
  if (hierarchy === -1) {
    return `__${columnName}${NESTING_LEVEL_INFIX}`;
  }

  return `__${columnName}${NESTING_LEVEL_INFIX}${hierarchy}`;
};

// type to represent the column metadata that powers a pivot
// this is used to transform the pivot data into a format that can be used by ag-grid
// extract the column metadata from powering a pivot from the widget instance details
export const getPivotColumns = (
  wInstanceDetails: Extract<WidgetInstanceType, { widget_type: WIDGET_TYPES.PIVOT_TABLE }>,
  wInstanceData: WidgetDataResponseType,
) => {
  const { data_mappings } = wInstanceDetails;
  const pivotColumns: PivotColumnMetadata[] = [];

  // iterate over each mapping in the widget instance details; each mapping is a stack in the stacked pivot
  data_mappings?.mappings?.forEach((mapping, mappingIndex) => {
    const { fields, ref } = mapping;

    const { columns, values } = fields;

    // iterate over each column in the columns array
    // each column is a pivot column
    columns?.forEach((col) => {
      pivotColumns?.push({
        kind: 'pivot',
        name: col.alias ? col.alias : col.column,
        dataType: col.type as 'string' | 'number' | 'date',
        sourceName: col?.column,
        mappingName: ref,
        alias: col.alias ? col.alias : col.column,
      });
    });

    // iterate over each value in the values array
    // each value is an aggregate column
    values?.forEach((val) => {
      pivotColumns?.push({
        kind: 'aggregate',
        name: val?.alias ? val?.alias : val?.column,
        dataType: val?.type as 'string' | 'number' | 'date',
        aggregation: val?.aggregation,
        sourceName: val?.column,
        mappingName: ref,
        alias: val?.alias ? val?.alias : val?.column,
      });
    });

    // if the mapping has no rows, we create a default row with the mapping name; the mapping name becomes the row group name (eg: Closing Balance)
    const mappingRows = [];

    if (data_mappings.mappings.length === 1) {
      if (fields?.rows) {
        mappingRows.push(...fields.rows);
      } else {
        mappingRows.push({
          column: mapping?.ref,
          type: 'string',
          alias: mapping?.ref,
        });
      }
    } else {
      mappingRows.push({ column: mapping?.ref, type: 'string', alias: '' }, ...(fields?.rows || []));
    }

    let currentLevel = 0;
    const colNameMapping: Record<
      string,
      {
        name: string;
        heirarchy: number;
        dataType: string;
        hasChildren: boolean;
        sourceName: string;
        mappingName: string;
        alias: string;
      }
    > = {};

    // iterate over each row in the rows array
    // we normalize the rows of every mapping to the same format
    mappingRows?.forEach((row) => {
      currentLevel += 1;
      const { column, type, alias } = row;

      // check if the row has hierarchy; if it does, we need to determine the depth of the hierarchy;
      // hirarchy is determined by _LEVEL_<n> suffix and the order of the columns in the row set
      const hasHierarchy = wInstanceData?.result[mappingIndex]?.columns?.find((c) =>
        c.column_name.startsWith(getNestedGroupingColName(alias ? alias : column, -1)),
      );

      if (hasHierarchy) {
        const depth = wInstanceData?.result[mappingIndex]?.columns?.filter((c) =>
          c.column_name.startsWith(getNestedGroupingColName(alias ? alias : column, -1)),
        )?.length;

        // iterate over each level of the hierarchy
        Array(depth)
          .fill(null)
          .forEach((_, colIndex: number) => {
            const colName = getNestedGroupingColName(alias ? alias : column, colIndex + 1);

            colNameMapping[colName] = {
              name: getGroupingColName(currentLevel + colIndex),
              heirarchy: currentLevel + colIndex,
              dataType: type,
              hasChildren: colIndex < depth - 1,
              sourceName: colName,
              mappingName: ref,
              alias: row?.alias || column,
            };
          });
        currentLevel += depth - 1;
      } else {
        colNameMapping[alias ? alias : column] = {
          name: getGroupingColName(currentLevel),
          heirarchy: currentLevel,
          dataType: type,
          hasChildren: false,
          sourceName: alias ? alias : column,
          mappingName: ref,
          alias: row?.alias || column,
        };
      }
    });

    Object.entries(colNameMapping).forEach(([, colData]) => {
      pivotColumns.push({
        kind: 'group',
        ...colData,
        dataType: colData?.dataType as 'string' | 'number' | 'date',
        alias: colData?.alias,
        maxHeirarchy: currentLevel,
      });
    });
  });

  return pivotColumns;
};

// transform the pivot data into a format that can be used by ag-grid
export const getPivotData = (
  pivotColumns: PivotColumnMetadata[],
  wInstanceData: WidgetDataResponseType,
  periodicity: PERIODICITY_TYPES,
) => {
  // Array to store transformed rows
  const rows: MapAny[] = [];

  // Process each result set in the widget data
  wInstanceData.result.forEach((resultSet) => {
    const resultRows = resultSet?.data;

    // Transform each row in the result set
    resultRows.forEach((row) => {
      // Create a copy of the row to transform
      const transformedRow = { ...row };

      // Process each field in the row
      Object.entries(row)?.forEach(([key, value]) => {
        // Get the mapping name from the NAME field
        const mappingName = transformedRow[PIVOT_REF];

        // Find matching pivot column based on mapping name and source column
        const pivotColumn = pivotColumns?.find(
          (col) => col?.mappingName === mappingName && (col?.sourceName === key || col?.alias === key),
        );

        if (pivotColumn) {
          // If matching pivot column found, transform the value using its data type
          transformedRow[pivotColumn?.name] = parseType(pivotColumn?.dataType as PIVOT_DATA_TYPES, value, periodicity);
        } else {
          if (key === PIVOT_REF) {
            // Special handling for REF field - find pivot column by source name
            const transformedColumn = pivotColumns?.find((col) => col?.sourceName === value);

            if (transformedColumn) {
              transformedRow[transformedColumn?.name] = value;
            }
          } else {
            // Keep original key-value pair if no transformation needed
            transformedRow[key] = value;
          }
        }
      });

      // Add transformed row to results
      rows.push(transformedRow);
    });
  });

  return rows;
};

export type ColumnContext = {
  name: string;
  alias: string;
};

export const getPivotColDefs = (
  pivotColumns: PivotColumnMetadata[],
): { coldefs: ColDef[]; columnContextMapping: Record<string, Record<string, ColumnContext>> } => {
  const columnContextMapping: Record<string, Record<string, ColumnContext>> = {};

  pivotColumns.forEach((col) => {
    if (!columnContextMapping[col?.mappingName]) {
      columnContextMapping[col?.mappingName] = {};
    }
    columnContextMapping[col?.mappingName][col?.name] = {
      name: col?.sourceName,
      alias: col?.alias,
    };
  });

  const coldefs: ColDef[] = pivotColumns
    .filter((col, index, self) => self?.findIndex((t) => t?.name === col?.name) === index)
    .map((col) => {
      switch (col.kind) {
        case 'group':
          return {
            field: col?.name,
            rowGroup: true,
            context: col,
            cellStyle: (params) =>
              getDynamicCellStyle(backendConfig.styleConfig.cellStyles, col?.name, params?.node?.level, params?.value),
          };
        case 'pivot':
          return {
            field: col?.name,
            pivot: true,
            pivotComparator: pivotComparator,
            headerComponent: PivotColGroupHeader,
            context: col,
            sortable: false,
          };
        case 'aggregate': {
          return {
            field: col?.name,
            aggFunc: AGGREGATION_TYPES.SUM,
            valueFormatter: (params) => formatCurrencyValue(params?.value),
            headerComponent: PivotColHeader,
            sortable: false,
            cellStyle: (params) =>
              getDynamicCellStyle(
                backendConfig.styleConfig.cellStyles,
                snakeCaseToSentenceCase(col.name),
                params.node?.level,
                params?.value,
              ),
            context: col,
          };
        }
      }
    });

  return { coldefs, columnContextMapping };
};

export type AGGridPivotNode<T extends AGGridPivotNode<T>> = {
  key?: string | null;
  parent?: T | null;
  childrenAfterGroup?: T[] | null;
};
export const flattenChildrenAfterGroup = <T extends AGGridPivotNode<T>>(node: T): T[] => {
  const { childrenAfterGroup } = node;

  if (!childrenAfterGroup) {
    return [];
  }

  return [
    ...childrenAfterGroup,
    ...childrenAfterGroup.map((child) => flattenChildrenAfterGroup(child).flat()).flat(),
  ] as T[];
};

export const shouldAllowExpandingRow = <T extends AGGridPivotNode<T>>(node: T) => {
  const flattenedChildren = flattenChildrenAfterGroup(node);

  // return false if the node has no children
  if (flattenedChildren?.length === 0) {
    return false;
  }

  // return false if the node has all children with key=""
  if (flattenedChildren.every((child) => child?.key === '' || child?.key === null || child?.key === undefined)) {
    return false;
  }

  return true;
};

export const getFilterContext = (
  widgetInstance: PivotTableWidgetInstanceType,
): Record<string, ColumnFilterConfig[]> => {
  const { data_mappings } = widgetInstance;

  const columnFilterConfigs: Record<string, ColumnFilterConfig[]> = {};

  data_mappings?.mappings?.forEach((mapping) => {
    if (!columnFilterConfigs[mapping?.ref]) {
      columnFilterConfigs[mapping?.ref] = [];
    }

    const { fields } = mapping;

    fields?.rows?.forEach((row) => {
      columnFilterConfigs[mapping?.ref].push({
        column: row?.column,
        filterType: row?.drilldown_filter_type as FILTER_TYPES,
        type: row?.drilldown_filter_operator as CONDITION_OPERATOR_TYPE,
      } as ColumnFilterConfig);
    });

    fields?.columns?.forEach((col) => {
      columnFilterConfigs[mapping?.ref].push({
        column: col?.column,
        filterType: col?.drilldown_filter_type as FILTER_TYPES,
        type: col?.drilldown_filter_operator as CONDITION_OPERATOR_TYPE,
      } as ColumnFilterConfig);
    });
  });

  return columnFilterConfigs;
};

// utility function to execute a function on a node and all its parents
const execOnNodeTillParent = (node: IRowNode | null, exec: (node: IRowNode) => void) => {
  if (!node) return;

  exec(node);

  if (node?.level > 0) {
    execOnNodeTillParent(node?.parent || null, exec);
  }
};

export const getRowLevelFilters = (
  rowColumnFilters: ColumnFilterConfig[],
  rowGroupMapping: Record<string, ColumnContext>,
  currentNode: IRowNode,
) => {
  // holds the row level filters
  const rowLevelFilters: Record<string, any> = {};

  // build filters by traversing every node in the tree from current node to the top node
  const exec = (node: IRowNode) => {
    const nodeColumnName = rowGroupMapping[node.field || ''];

    if (!nodeColumnName) return;

    let filterColumnName = nodeColumnName?.alias || nodeColumnName?.name;

    const rowColumnFilterConfig = rowColumnFilters?.find((col) => {
      if (isTagColumn(nodeColumnName?.name)) {
        const { name } = unwrapTagColumn(nodeColumnName?.name);

        if (col?.column === name) {
          filterColumnName = nodeColumnName?.name;

          return true;
        }
      }

      return col?.column === nodeColumnName?.alias || col?.column === nodeColumnName?.name;
    });

    if (rowColumnFilterConfig) {
      const _rowColumnFilterConfig = { ...rowColumnFilterConfig };

      rowLevelFilters[filterColumnName] = _rowColumnFilterConfig;

      switch (rowColumnFilterConfig?.filterType) {
        case FILTER_TYPES.MULTI_SELECT:
          rowLevelFilters[filterColumnName].values = [node?.key];
          break;
        case FILTER_TYPES.SEARCH:
          rowLevelFilters[filterColumnName].values = [node?.key];
          break;
      }
    }
  };

  // execute the function on the current node and all its parents
  execOnNodeTillParent(currentNode, exec);

  return rowLevelFilters;
};

// TODO: Break this function into smaller functions
export const getColumnLevelFilters = (
  columnColumnFilters: ColumnFilterConfig[],
  colDefs: ColDef[],
  mappingColumnContext: Record<string, ColumnContext>,
  currentDateConfig: {
    periodicity: PERIODICITY_TYPES;
    widgetSelectedFilter: Record<string, any>;
  },
  currentPivotKeys: string[],
): Record<string, any> => {
  // holds the column level filters
  const columnLevelFilters: Record<string, any> = {};

  // for all coldefs in the pivot, build the column level filters
  colDefs?.forEach((colDef, i) => {
    // get the label (alias or name) of the column
    const colLabel = colDef.context?.alias || colDef.context?.name;

    // find the filter config for the label; either the column name or the column alias
    const columnColumnFilterConfig = columnColumnFilters?.find((col) => {
      return col?.column === colLabel || mappingColumnContext[colLabel]?.name === col?.column;
    });

    // get the pivot key for the column
    // this is the key of the column in the pivot
    const pivotKey = currentPivotKeys[i];

    if (columnColumnFilterConfig) {
      columnLevelFilters[columnColumnFilterConfig?.column] = columnColumnFilterConfig;
      switch (columnColumnFilterConfig.filterType) {
        case FILTER_TYPES.MULTI_SELECT: {
          columnLevelFilters[columnColumnFilterConfig.column].values = [pivotKey];
          break;
        }
        case FILTER_TYPES.SEARCH: {
          columnLevelFilters[columnColumnFilterConfig.column].values = [pivotKey];
          break;
        }
        case FILTER_TYPES.DATE_RANGE: {
          const [updatedDateFrom, updatedDateTo] = getDateRangeWithPeriodicity(
            currentDateConfig.periodicity,
            pivotKey,
            currentDateConfig.widgetSelectedFilter[columnColumnFilterConfig.column]?.dateFrom,
            currentDateConfig.widgetSelectedFilter[columnColumnFilterConfig.column]?.dateTo,
          );

          columnLevelFilters[columnColumnFilterConfig.column].dateFrom = updatedDateFrom;
          columnLevelFilters[columnColumnFilterConfig.column].dateTo = updatedDateTo;
          break;
        }
      }
    }
  });

  return columnLevelFilters;
};

// used to map the ref to the dataset id
export const getWidgetMappingDatasets = (widgetInstance: PivotTableWidgetInstanceType): Record<string, string> => {
  const { data_mappings } = widgetInstance;

  const mappingDatasets: Record<string, string> = {};

  data_mappings?.mappings?.forEach((mapping) => {
    mappingDatasets[mapping?.ref] = mapping?.dataset_id;
  });

  return mappingDatasets;
};

// utility function to get the top node of the tree
export const getTopNode = (node: IRowNode): IRowNode => {
  if (!node?.parent || node?.level === 0) return node;

  return getTopNode(node.parent);
};

const isTagColumn = (colName: string) => {
  return colName?.startsWith('__') && colName?.includes(NESTING_LEVEL_INFIX);
};

// remove the prefix and the infix --- for example __tag_LEVEL_1 should return tag and 1
export const unwrapTagColumn = (colName: string): { name: string; hierarchy: number } => {
  // remove the prefix and the infix
  const nameWithoutPrefix = colName?.substring(2);
  const name = nameWithoutPrefix?.substring(0, nameWithoutPrefix?.indexOf(NESTING_LEVEL_INFIX));
  const hierarchy = parseInt(nameWithoutPrefix.split(NESTING_LEVEL_INFIX).pop() || '');

  return { name, hierarchy };
};

export const concatTagFilters = (filters: Record<string, any>) => {
  const concatenatedFilters = { ...filters };

  // heirarchy <> filter value
  const tagFilters: Record<string, any> = {};

  let tagName = '';

  Object.entries(filters)?.forEach(([key, value]) => {
    if (isTagColumn(key)) {
      const { name, hierarchy } = unwrapTagColumn(key);

      tagName = name;
      tagFilters[hierarchy] = value?.values?.[0];
      delete concatenatedFilters[key];
    }
  });

  const sortedTagFilters = Object.entries(tagFilters)?.sort((a, b) => parseInt(a[0]) - parseInt(b[0]));

  if (tagName && sortedTagFilters?.length > 0) {
    concatenatedFilters[tagName] = {
      filterType: FILTER_TYPES.SEARCH,
      type: CONDITION_OPERATOR_TYPE.STARTS_WITH,
      values: [sortedTagFilters?.map(([, value]) => value).join('.')],
    };
  }

  return concatenatedFilters;
};

export const formatRowTitleValue = (value: string): string => {
  if (!isNaN(Date.parse(value))) {
    return getFormattedDateWithPeriodicity(PERIODICITY_TYPES.DAILY, value);
  }

  if (value?.includes('_')) {
    return snakeCaseToSentenceCase(value);
  }

  return capitalizeFirstLetter(value);
};

const parseDate = (dateStr: string): Date | null => {
  if (!dateStr) return null;

  // Handle date ranges like "1-7 Jan 2025" by extracting the first date
  const rangeMatch = dateStr?.match(/^(\d+)-\d+\s([A-Za-z]+)\s(\d{4})$/);

  if (rangeMatch) {
    dateStr = `${rangeMatch[1]} ${rangeMatch[2]} ${rangeMatch[3]}`;
  }

  const formats = [DATE_FORMATS.d_MMM_yyyy, DATE_FORMATS.MMM_yyyy, DATE_FORMATS.QQ_yyyy, DATE_FORMATS.YYYY];

  for (const format of formats) {
    const parsedDate = parse(dateStr, format, new Date());

    if (isValid(parsedDate)) {
      return parsedDate;
    }
  }

  return null;
};

const pivotComparator = (str1: string, str2: string): number => {
  const dateA = parseDate(str1);
  const dateB = parseDate(str2);

  return dateA && dateB ? dateA.getTime() - dateB.getTime() : dateA ? -1 : dateB ? 1 : str1.localeCompare(str2);
};

export const formatColGroupHeaderDisplayName = (displayName: string) => {
  const match = displayName?.match(/^(\d{1,2} \w{3} \d{4})$/);

  if (match) {
    const dateString = match[1];
    const parsedDate = parse(dateString, DATE_FORMATS.d_MMM_yyyy, new Date());

    if (isValid(parsedDate)) {
      const dayOfWeek = format(parsedDate, DATE_FORMATS.EEE);

      return { mainText: dateString, suffix: dayOfWeek };
    }
  }

  const formattedText = displayName?.includes('_')
    ? snakeCaseToSentenceCase(displayName)
    : capitalizeFirstLetter(displayName);

  return { mainText: formattedText, suffix: '' };
};
