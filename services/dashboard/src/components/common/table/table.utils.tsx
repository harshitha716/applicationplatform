import { IServerSideGetRowsRequest, themeQuartz } from 'ag-grid-community';
import { DATE_FORMATS } from 'constants/date.constants';
import { format } from 'date-fns';
import { MapAny } from 'types/commonTypes';
import {
  AggregationType,
  FilterModelType,
  FilterType,
  GroupByType,
  LogicalOperatorType,
  OrderByType,
  OrderType,
  RequestType,
} from 'types/components/table.type';
import { checkIsObjectEmpty } from 'utils/common';
import {
  AggregationFunctionMap,
  ArrayFilters,
  LogicalOperatorMap,
  PAGE_SIZE,
} from 'components/common/table/table.constants';
import { FILTER_TYPES } from 'components/filter/filter.types';
import { CONDITION_OPERATOR_TYPE } from 'components/filter/filters.constants';

const getFiltersFromGroupKeys = (request: IServerSideGetRowsRequest): FilterType[] => {
  const { groupKeys, rowGroupCols } = request;

  if (!groupKeys?.length || !rowGroupCols?.length) {
    return [];
  }

  return groupKeys?.map((key, index) => ({
    column: rowGroupCols?.[index]?.id,
    operator: CONDITION_OPERATOR_TYPE.EQUAL,
    value: key,
  }));
};

export const getConditionValues = (condition: MapAny): FilterType | null => {
  switch (condition.filterType) {
    case FILTER_TYPES.AMOUNT_RANGE:
      if (condition.type === CONDITION_OPERATOR_TYPE.IN_BETWEEN) {
        if (condition.filterTo !== '' && condition.filter !== '')
          return {
            column: condition.colId,
            operator: condition.type,
            value: [Number(condition.filter), Number(condition.filterTo)],
          };
        else return null;
      } else if (condition.type === CONDITION_OPERATOR_TYPE.IS_NULL) {
        return {
          column: condition.colId,
          operator: condition.type,
          value: '',
        };
      } else if (condition.filter !== '') {
        return {
          column: condition.colId,
          operator: condition.type,
          value: Number(condition.filter),
        };
      } else return null;
    case FILTER_TYPES.MULTI_SELECT:
      if (condition?.values?.length) {
        return {
          column: condition.colId,
          operator: condition.type,
          value: condition.values,
        };
      } else if (condition.type === CONDITION_OPERATOR_TYPE.IS_NULL) {
        return {
          column: condition.colId,
          operator: condition.type,
          value: '',
        };
      } else return null;
    case FILTER_TYPES.DATE_RANGE:
      if (condition.dateFrom && condition.dateTo) {
        const startDate = new Date(condition.dateFrom);

        startDate.setHours(0, 0, 0, 0);
        const endDate = new Date(condition.dateTo);

        endDate.setHours(23, 59, 59, 999);

        return {
          column: condition.colId,
          operator: condition.type,
          value: [format(startDate, DATE_FORMATS.YYYYMMDD_HHMMSS), format(endDate, DATE_FORMATS.YYYYMMDD_HHMMSS)],
        };
      } else return null;
    case FILTER_TYPES.SEARCH:
      return {
        column: condition.colId,
        operator: condition.type,
        value: ArrayFilters.includes(condition.type) ? [condition.filter] : condition.filter,
      };
    case FILTER_TYPES.ARRAY_SEARCH:
      return {
        column: condition?.colId,
        operator: condition?.type,
        value: condition?.value?.split(','),
      };
    default:
      return null;
  }
};

const parseCondition = (condition: MapAny): FilterType | null => {
  if (condition.conditions) {
    return {
      logical_operator: LogicalOperatorMap[condition.type] || LogicalOperatorType.OperatorLogicalAnd,
      conditions: condition.conditions.map((cond: MapAny) => parseCondition(cond)),
    };
  } else {
    return getConditionValues(condition);
  }
};

export const convertToFilterModel = (input: MapAny | null): FilterModelType | null => {
  if (!input) {
    return null;
  } else if (input.filterType === 'join') {
    return {
      logical_operator: LogicalOperatorMap[input.type] || LogicalOperatorType.OperatorLogicalAnd,
      conditions: input.conditions
        .map((condition: MapAny) => parseCondition(condition))
        .filter((condition: MapAny) => condition !== null),
    };
  } else if (input.conditions) {
    return {
      logical_operator: LogicalOperatorMap[input.operator] || LogicalOperatorType.OperatorLogicalAnd,
      conditions: input.conditions
        .map((condition: MapAny) => parseCondition(condition))
        .filter((condition: MapAny) => condition !== null),
    };
  } else {
    const keys = Object.keys(input);

    if (keys.length) {
      const formattedConditions = keys.map((key) => ({ colId: key, ...input?.[key] }));
      const conditions = formattedConditions
        .map((condition: MapAny) => parseCondition(condition))
        .filter((condition: MapAny | null) => condition !== null);

      if (conditions.length) {
        return {
          logical_operator: LogicalOperatorType.OperatorLogicalAnd,
          conditions,
        };
      } else return null;
    }

    return null;
  }
};

export const getFilterModelFromGroupAndFilterModel = (
  request: IServerSideGetRowsRequest,
  hiddenColumnFilters?: MapAny,
): FilterModelType | null => {
  const filtersFromGroup = getFiltersFromGroupKeys(request);
  const filterModel = checkIsObjectEmpty(hiddenColumnFilters ?? {})
    ? request.filterModel
    : checkIsObjectEmpty(request.filterModel ?? {})
      ? (hiddenColumnFilters ?? null)
      : { ...request.filterModel, ...hiddenColumnFilters };
  const filtersFromFilterModel = convertToFilterModel(filterModel);

  if (filtersFromGroup.length) {
    return {
      logical_operator: LogicalOperatorType.OperatorLogicalAnd,
      conditions: filtersFromFilterModel ? [...filtersFromGroup, filtersFromFilterModel] : filtersFromGroup,
    };
  }

  return filtersFromFilterModel;
};

const getGroupByColumns = (request: IServerSideGetRowsRequest): GroupByType[] => {
  const { rowGroupCols, groupKeys } = request;
  const rowGroupsToBeUsed = groupKeys?.length ? rowGroupCols.slice(groupKeys.length) : rowGroupCols;

  if (rowGroupsToBeUsed?.length) {
    return [
      {
        column: rowGroupsToBeUsed[0]?.id,
        alias: rowGroupsToBeUsed[0]?.id,
      },
    ];
  }

  return [];
};

const getAggregations = (
  request: IServerSideGetRowsRequest,
  useAlias?: boolean,
  ignoreGroupCheck?: boolean,
): AggregationType[] => {
  const { valueCols, rowGroupCols, groupKeys } = request;

  if (rowGroupCols?.length === groupKeys?.length && !ignoreGroupCheck) {
    return [];
  }

  return valueCols.map((item) => ({
    column: item.id,
    alias: useAlias ? item.displayName : item.id,
    function: AggregationFunctionMap[item?.aggFunc ?? 'sum'],
  }));
};

const getOrderByColumns = (request: IServerSideGetRowsRequest): OrderByType[] => {
  const { sortModel } = request;

  return sortModel?.map((item) => ({
    column: item.colId,
    order: item.sort as OrderType,
  }));
};

const formatRequest = (
  request: IServerSideGetRowsRequest,
  fx_currency?: string | undefined,
  useAlias?: boolean,
  ignoreGroupCheck?: boolean,
  disableTotalCount?: boolean,
  hiddenColumnFilters?: MapAny,
  drilldownFilters?: MapAny,
): RequestType => {
  const { endRow } = request;

  return {
    filters: drilldownFilters ?? getFilterModelFromGroupAndFilterModel(request, hiddenColumnFilters),
    aggregations: getAggregations(request, useAlias, ignoreGroupCheck),
    group_by: getGroupByColumns(request),
    order_by: getOrderByColumns(request),
    pagination: {
      page: endRow ? Math.ceil(endRow / PAGE_SIZE) : 1,
      page_size: PAGE_SIZE,
    },
    get_total_records: !disableTotalCount,
    fx_currency: !fx_currency || fx_currency === 'local' ? undefined : fx_currency,
  };
};

export const encodeRequest = (request: RequestType): string => {
  const jsonString = JSON.stringify(request);

  return jsonString;
};

export const getEncodedRequest = (
  request: IServerSideGetRowsRequest,
  fx_currency?: string,
  useAlias?: boolean,
  ignoreGroupCheck?: boolean,
  disableTotalCount?: boolean,
  hiddenColumnFilters?: MapAny,
  drilldownFilters?: MapAny,
): string => {
  const formattedRequest = formatRequest(
    request,
    fx_currency,
    useAlias,
    ignoreGroupCheck,
    disableTotalCount,
    hiddenColumnFilters,
    drilldownFilters,
  );
  const encodedRequest = encodeRequest(formattedRequest);

  return encodedRequest;
};

export const getDataTableTheme = (params: MapAny) => themeQuartz.withParams(params);
