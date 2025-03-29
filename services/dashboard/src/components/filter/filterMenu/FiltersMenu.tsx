import React, { FC, useRef, useState } from 'react';
import { useOnClickOutside } from 'hooks';
import { POSITION_TYPES } from 'types/common/components';
import { TooltipPositions } from 'components/common/tooltip';
import { FilterConfigType } from 'components/filter/filter.types';
import FilterControlButton from 'components/filter/FilterControlButton';
import SelectFilterMenuItem from 'components/filter/filterMenu/SelectFilterMenuItem';
import { useFiltersContextStore } from 'components/filter/filters.context';

interface FiltersMenuProps {
  filtersList?: Record<string, FilterConfigType>;
  onAddFilter: (filterKey: string) => void;
  label?: string;
  tooltipText?: string;
  currentPageFilters?: string[];
}

const FiltersMenu: FC<FiltersMenuProps> = ({ onAddFilter, label, tooltipText, currentPageFilters }) => {
  const {
    state: { filtersConfig },
  } = useFiltersContextStore();

  const [isOpen, setIsOpen] = useState<boolean>(false);
  const menuRef = useRef<HTMLDivElement>(null);
  const controlRef = useRef<HTMLDivElement>(null);

  useOnClickOutside(menuRef, () => {
    setIsOpen(false);
  }, [controlRef]);

  const toggleMenu = () => {
    setIsOpen((prev) => !prev);
  };

  const getMenuPlacement = () => {
    if (menuRef.current) {
      const { left } = menuRef.current.getBoundingClientRect();
      const isMenuCutOff = left + 200 > window.innerWidth;

      return isMenuCutOff ? POSITION_TYPES.LEFT : POSITION_TYPES.RIGHT;
    }

    return POSITION_TYPES.RIGHT;
  };

  const onAddfilter = (filterKey: string) => {
    onAddFilter(filterKey);
    toggleMenu();
  };

  return (
    <div className='relative'>
      <div ref={controlRef}>
        <FilterControlButton
          onClick={toggleMenu}
          tooltipPosition={TooltipPositions.TOP}
          tooltipText={tooltipText}
          id='add-filters'
        >
          {label}
        </FilterControlButton>
      </div>
      <SelectFilterMenuItem
        menuRef={menuRef}
        isOpen={isOpen}
        getMenuPlacement={getMenuPlacement}
        filtersConfig={filtersConfig ?? []}
        onAddFilter={onAddfilter}
        currentPageFilters={currentPageFilters ?? []}
      />
    </div>
  );
};

export default FiltersMenu;
