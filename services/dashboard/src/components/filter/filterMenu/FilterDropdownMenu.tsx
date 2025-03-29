import React, { FC } from 'react';
import { defaultFnType, MapAny, OptionsType } from 'types/commonTypes';
import { FILTER_TYPES, FilterConfigType } from 'components/filter/filter.types';
import { AG_GRID_FILTER_TYPES } from 'components/filter/filters.constants';
import { useFiltersContextStore } from 'components/filter/filters.context';

interface FilterDropdownMenuProps {
  className?: string;
  filterKey: string;
  filterType: FILTER_TYPES;
  onClose?: defaultFnType;
  isOpen?: boolean;
  id?: string;
  onSelect?: defaultFnType;
  onChange?: defaultFnType;
  filter?: FilterConfigType;
  forView?: string;
  showSearch?: boolean;
  label?: string;
  allowClear?: boolean;
  filterComponentProps?: MapAny;
  disableFutureDate?: boolean;
  operatorOptions?: OptionsType[];
  isPeriodicityEnabled?: boolean;
  onFilterChange?: (value: string[]) => void;
}

const FilterDropdownMenu: FC<FilterDropdownMenuProps> = ({
  filterKey,
  filterType,
  filterComponentProps = {},
  isPeriodicityEnabled = false,
  onFilterChange,
  isOpen,
  label,
}) => {
  const {
    state: { filtersConfig },
  } = useFiltersContextStore();

  const values = filtersConfig?.find((filter) => filter?.key === filterKey)?.values || [];

  const FilterMenuComponent = AG_GRID_FILTER_TYPES[filterType];

  return FilterMenuComponent ? (
    <FilterMenuComponent
      column={{ colId: filterKey }}
      values={values || []}
      key={filterKey}
      className='w-full'
      isPeriodicityEnabled={isPeriodicityEnabled}
      onFilterChange={onFilterChange}
      isOpen={isOpen}
      label={label}
      {...filterComponentProps}
    />
  ) : null;
};

export default FilterDropdownMenu;
