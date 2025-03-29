import { PERIODICITY_TYPES } from 'constants/date.constants';
import { FilterDefaultValueType, SheetFilterType } from 'types/api/pagesApi.types';
import { MapAny } from 'types/commonTypes';
import { getPastDateByNumberOfDays } from 'utils/common';
import { FILTER_TYPES } from 'components/filter/filter.types';
import { CONDITION_OPERATOR_TYPE } from 'components/filter/filters.constants';

export const getFormattedSheetsFiltersConfig = (filter: SheetFilterType) => {
  return {
    key: filter?.targets?.[0]?.column,
    label: filter?.name,
    values: filter?.options,
    datatype: filter?.data_type,
    type: filter?.filter_type,
    targets: filter?.targets,
    widgetsInScope: filter?.widgets_in_scope,
  };
};

const getFilterDefaultValue = (filter: FilterDefaultValueType, filterType: FILTER_TYPES) => {
  switch (filterType) {
    case FILTER_TYPES.SEARCH:
      return {
        filterType: filterType,
        type: filter.operator,
        filter: filter?.value?.[0],
      };
    case FILTER_TYPES.AMOUNT_RANGE:
      return {
        filterType: filterType,
        type: filter.operator,
        filter: filter?.value?.[0],
        filterTo: filter?.value?.[1],
      };
    case FILTER_TYPES.DATE_RANGE:
      return {
        filterType: filterType,
        dateFrom: filter?.value?.[0] ?? getPastDateByNumberOfDays(30).toISOString(),
        dateTo: filter?.value?.[1] ?? getPastDateByNumberOfDays(0).toISOString(),
        periodicity: PERIODICITY_TYPES.DAILY,
        type: filter?.operator,
      };
    case FILTER_TYPES.MULTI_SELECT:
      return {
        filterType: filterType,
        type: filter?.operator,
        values: filter?.value,
      };
  }
};

export const getDefaultFilterValues = (filters: SheetFilterType[]) => {
  const defaultFilters: MapAny = {};

  filters.forEach((filter) => {
    if (filter?.default_value) {
      defaultFilters[filter?.targets?.[0]?.column] = getFilterDefaultValue(filter?.default_value, filter?.filter_type);
    } else if (filter?.filter_type === FILTER_TYPES.DATE_RANGE) {
      defaultFilters[filter?.targets?.[0]?.column] = getFilterDefaultValue(
        { value: [], operator: CONDITION_OPERATOR_TYPE.IN_BETWEEN },
        filter?.filter_type,
      );
    }
  });

  return defaultFilters;
};
