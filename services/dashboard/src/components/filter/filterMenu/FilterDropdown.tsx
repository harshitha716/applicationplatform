import React, { FC, useMemo, useRef, useState } from 'react';
import { useOnClickOutside } from 'hooks';
import { MapAny } from 'types/commonTypes';
import { cn } from 'utils/common';
import { FILTER_TYPES, FilterConfigType } from 'components/filter/filter.types';
import FilterControl from 'components/filter/filterMenu/FilterDropdownControl';
import FilterDropdownMenu from 'components/filter/filterMenu/FilterDropdownMenu';

interface FilterDropdownProps {
  index: number;
  filter: FilterConfigType;
  onRemoveFilter?: ((filterKey: string) => void) | null;
  isFilterSelected: boolean;
  props?: MapAny;
  controlClassName?: string;
  allowClear?: boolean;
  allowActions: boolean;
  isPeriodicityEnabled?: boolean;
  onFilterChange?: (value: string[]) => void;
  closeOnSelect?: boolean;
  isRightAligned?: boolean;
}

const FilterDropdown: FC<FilterDropdownProps> = ({
  index,
  filter,
  onRemoveFilter,
  isFilterSelected,
  props = {},
  controlClassName = '',
  allowClear = true,
  allowActions = true,
  isPeriodicityEnabled = false,
  onFilterChange,
  closeOnSelect = false,
  isRightAligned = false,
}) => {
  const [isOpen, setIsOpen] = useState<boolean>(!isFilterSelected && allowActions);
  const controlRef = useRef<HTMLDivElement>(null);
  const menuRef = useRef<HTMLDivElement>(null);

  useOnClickOutside(menuRef, () => {
    setIsOpen(false);
  }, [controlRef]);

  const getMenuPlacement = useMemo(() => {
    if (menuRef.current) {
      const { left } = menuRef.current.getBoundingClientRect();

      return left + 300 > window.innerWidth;
    }

    return false;
  }, [menuRef, isOpen]);

  const onClick = () => {
    setIsOpen((prev) => !prev);
  };

  const onChange = (value: string[]) => {
    if (closeOnSelect) {
      setIsOpen(false);
    }

    onFilterChange?.(value);
  };

  return (
    <div key={index} className='relative w-fit'>
      <div ref={controlRef}>
        <FilterControl
          filterConfig={filter}
          key={index}
          isMenuDropdownOpen={false}
          onClick={onClick}
          onClear={onRemoveFilter}
          controlClassName={controlClassName}
          allowClear={allowClear}
          isOpen={isOpen}
        />
      </div>
      <div
        ref={menuRef}
        className={cn(
          `absolute top-full mt-1.5 w-fit shadow-dropdown transition-all duration-500 z-50 min-w-[218px]`,
          controlClassName,
          isRightAligned ? 'right-0' : getMenuPlacement ? 'right-0' : 'left-0',
          isOpen ? '' : 'max-h-0 overflow-hidden border-0',
        )}
      >
        <FilterDropdownMenu
          forView='filters'
          filterKey={filter?.key}
          filterType={filter?.type as FILTER_TYPES}
          label={filter?.label}
          className='min-w-[200px] w-full'
          isOpen={isOpen}
          onClose={() => setIsOpen(false)}
          allowClear={allowClear}
          isPeriodicityEnabled={isPeriodicityEnabled}
          onFilterChange={onChange}
          {...props}
        />
      </div>
    </div>
  );
};

export default FilterDropdown;
