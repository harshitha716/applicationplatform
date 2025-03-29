import React, { FC, useCallback, useEffect, useRef, useState } from 'react';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { useOnClickOutside } from 'hooks';
import { defaultFn, defaultFnType, MapAny } from 'types/commonTypes';
import { TooltipPositions } from 'components/common/tooltip';
import { FilterConfigType } from 'components/filter/filter.types';
import { getFilterValueForKey } from 'components/filter/filter.utils';
import FilterControlButton from 'components/filter/FilterControlButton';
import ClearFiltersConfirmationPopup from 'components/filter/filterMenu/ClearFiltersConfirmationPopup';
import FilterDropdown from 'components/filter/filterMenu/FilterDropdown';
import FiltersMenu from 'components/filter/filterMenu/FiltersMenu';
import { FILTER_KEYS } from 'components/filter/filters.constants';
import { filtersContextActions, useFiltersContextStore } from 'components/filter/filters.context';

interface FiltersContainerProps {
  onClearAllFilters?: defaultFnType;
  onClearRules?: defaultFnType;
  onOpenAdvancedSearch?: defaultFnType;
  persistId?: string;
  onSetTotalSelectedFilters?: (val: number) => void;
  filterConfig: FilterConfigType[];
  className?: string;
  allowActions?: boolean;
  controlClassName?: string;
  allowClear?: boolean;
  label?: string;
  showResetFilters?: boolean;
  isPeriodicityEnabled?: boolean;
  isRightAligned?: boolean;
}

const FiltersContainer: FC<FiltersContainerProps> = ({
  onClearAllFilters = defaultFn,
  persistId,
  onSetTotalSelectedFilters,
  filterConfig,
  className = 'px-6',
  allowActions = true,
  controlClassName = '',
  allowClear = true,
  label = 'Add Filters',
  isPeriodicityEnabled = false,
  isRightAligned = false,
}) => {
  const [shouldShowConfirmationPopup, setShouldShowConfirmationPopup] = useState(false);
  const {
    dispatch,
    state: { selectedFilters, selectedFiltersInUI, currentPageFilters },
  } = useFiltersContextStore();

  const [filtersList, setFiltersList] = useState<FilterConfigType[]>([]);

  const onAddFilterToFiltersList = (filterKey: string, list: FilterConfigType[], value: FilterConfigType) => {
    const filterItemIndex = list.findIndex((item: FilterConfigType) => item?.key === filterKey);

    if (filterItemIndex === -1) {
      list.push(value);

      return;
    }

    list[filterItemIndex] = value;
  };

  const onRemoveFiltersWithoutKeys = (list: FilterConfigType[], selectedFiltersInUI: MapAny) => {
    const keys = Object.keys(selectedFiltersInUI);

    for (let i = list.length - 1; i >= 0; i--) {
      const filter = list[i];

      if (!keys?.includes(filter?.key)) {
        list?.splice(i, 1);
      }
    }
  };

  const onSetFiltersList = useCallback(() => {
    const list = [...filtersList];

    const selectedFilters = selectedFiltersInUI;

    for (const key in selectedFilters) {
      const value: any = getFilterValueForKey(key as FILTER_KEYS, filterConfig, selectedFilters);

      onAddFilterToFiltersList(key, list, value);
    }

    onRemoveFiltersWithoutKeys(list, selectedFiltersInUI);
    onSetTotalSelectedFilters?.(list.length);

    setFiltersList(list);
  }, [selectedFiltersInUI, selectedFilters]);

  useEffect(() => {
    onSetFiltersList();
  }, [selectedFiltersInUI]);

  const handleResetFilters = () => {
    setShouldShowConfirmationPopup(false);

    dispatch({
      type: filtersContextActions.RESET_ALL_FILTERS,
      payload: { shouldClearDate: false },
    });

    onClearAllFilters?.();
  };

  const onAddEmptyFilter = (filterKey: string) => {
    dispatch({
      type: filtersContextActions.ADD_EMPTY_STATE_FILTER,
      payload: { filterKey },
    });
  };

  const onRemoveFilter = (filterKey: string) => {
    dispatch({
      type: filtersContextActions.REMOVE_FILTER,
      payload: { filterKey },
    });
  };

  const confirmationPopupRef = useRef<HTMLDivElement>(null);
  const confirmationPopupControlRef = useRef<HTMLButtonElement>(null);

  useOnClickOutside(confirmationPopupRef, () => {
    setShouldShowConfirmationPopup(false);
  }, [confirmationPopupControlRef]);

  return (
    <div id={`${persistId}_FILTERS_CONTAINER`} className={`flex items-center flex-wrap gap-2 z-50 ${className}`}>
      {filtersList.map((filter, index) => (
        <FilterDropdown
          key={index}
          index={index}
          filter={filter}
          onRemoveFilter={allowActions ? onRemoveFilter : null}
          allowActions={allowActions}
          isFilterSelected={selectedFilters[filter?.key]}
          controlClassName={controlClassName}
          allowClear={allowClear}
          isPeriodicityEnabled={isPeriodicityEnabled}
          isRightAligned={isRightAligned}
        />
      ))}

      {allowActions && !filtersList?.length && <FiltersMenu label={label} onAddFilter={onAddEmptyFilter} />}

      {allowActions && filtersList?.length > 0 ? (
        <>
          <FiltersMenu
            tooltipText='Add Filters'
            currentPageFilters={currentPageFilters}
            onAddFilter={onAddEmptyFilter}
          />

          <div className='relative'>
            <FilterControlButton
              tooltipText='Remove all filters'
              tooltipPosition={TooltipPositions.TOP}
              onClick={() => setShouldShowConfirmationPopup(!shouldShowConfirmationPopup)}
              buttonRef={confirmationPopupControlRef}
              icon='x-close'
              iconCategory={ICON_SPRITE_TYPES.GENERAL}
              id='clear-all-filters'
            >
              {shouldShowConfirmationPopup ? (
                <ClearFiltersConfirmationPopup
                  containerRef={confirmationPopupRef}
                  onClick={handleResetFilters}
                  onCancel={() => setShouldShowConfirmationPopup(false)}
                  className='absolute left-0 z-9999'
                />
              ) : null}
            </FilterControlButton>
          </div>
        </>
      ) : null}
    </div>
  );
};

export default FiltersContainer;
