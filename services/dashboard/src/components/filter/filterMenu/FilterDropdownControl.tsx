import React, { FC } from 'react';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { defaultFnType } from 'types/commonTypes';
import { cn } from 'utils/common';
import { FilterConfigType } from 'components/filter/filter.types';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

interface FilterDropdownControlProps {
  onClick?: defaultFnType;
  className?: string;
  filterConfig: FilterConfigType;
  onClear?: ((filterKey: string) => void) | null;
  controlClassName?: string;
  isMenuDropdownOpen?: boolean;
  allowClear?: boolean;
  isOpen?: boolean;
}

const FilterDropdownControl: FC<FilterDropdownControlProps> = ({
  onClick,
  filterConfig,
  className = '',
  controlClassName = '',
  onClear,
  isMenuDropdownOpen,
  allowClear,
  isOpen,
}) => {
  const handleRemoveFilter = (e: React.MouseEvent) => {
    if (allowClear) {
      e.stopPropagation();
      onClear?.(filterConfig.key);
    }
  };

  return (
    <div
      data-testid={`filter-control-${filterConfig?.key}`}
      className={`cursor-pointer relative ${className}`}
      onClick={onClick}
    >
      <div
        className={`select-none rounded h-[26px] flex items-center gap-1.5 border hover:border-DIVIDER_SAIL_4 border-DIVIDER_SAIL_3 px-1.5 py-1.5 w-fit bg-white ${
          isMenuDropdownOpen ? 'border-DIVIDER_SAIL_4' : ''
        } ${controlClassName}`}
      >
        <div className='f-12-400 text-GRAY_900'>{filterConfig?.label}</div>
        <div className='f-12-500 text-GRAY_1000 whitespace-nowrap'>{filterConfig?.title}</div>
        <div onClick={handleRemoveFilter}>
          {allowClear ? (
            <SvgSpriteLoader
              id='x-close'
              iconCategory={ICON_SPRITE_TYPES.GENERAL}
              width={12}
              height={12}
              className={'text-GRAY_700 mt-0.5'}
            />
          ) : (
            <SvgSpriteLoader
              id='chevron-down'
              iconCategory={ICON_SPRITE_TYPES.ARROWS}
              width={16}
              height={16}
              className={cn('text-GRAY_700 transition-transform duration-300', isOpen ? 'rotate-180' : 'rotate-0')}
            />
          )}
        </div>
      </div>
    </div>
  );
};

export default FilterDropdownControl;
