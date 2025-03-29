import React, { FC, useCallback, useEffect, useRef, useState } from 'react';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { useOnClickOutside } from 'hooks';
import { SIZE_TYPES } from 'types/common/components';
import { OptionsType } from 'types/commonTypes';
import { camelCaseToNormalText, debounce } from 'utils/common';
import Input from 'components/common/input';
import { FILTER_TYPES } from 'components/filter/filter.types';
import { SEARCH_FILTER_OPTIONS } from 'components/filter/filters.constants';
import { filtersContextActions, useFiltersContextStore } from 'components/filter/filters.context';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

interface SearchFilterMenuItemProps {
  column: { colId: string };
  values: string[];
  className?: string;
  isOpen?: boolean;
  label?: string;
}

const SearchFilterMenuItem: FC<SearchFilterMenuItemProps> = ({ column, className, isOpen = false, label }) => {
  const ref = useRef(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const columnId = column?.colId;
  const {
    state: { selectedFilters },
    dispatch,
  } = useFiltersContextStore();
  const currentOperatorValue = selectedFilters[columnId]?.type;
  const currentOperator = SEARCH_FILTER_OPTIONS.find((option) => option.value === currentOperatorValue);
  const [searchValue, setSearchValue] = useState(selectedFilters[columnId]?.filter || '');
  const [isConditionOptionsOpen, setIsConditionOptionsOpen] = useState(false);
  const [selectedOperator, setSelectedOperator] = useState<OptionsType>(currentOperator ?? SEARCH_FILTER_OPTIONS[0]);

  const setFilter = (operator: string, searchValue: string) => {
    dispatch({
      type: filtersContextActions.SET_SELECTED_FILTERS,
      payload: {
        selectedFilters: {
          [columnId]: searchValue
            ? {
                filterType: FILTER_TYPES.SEARCH,
                type: operator,
                filter: searchValue,
              }
            : {},
        },
      },
    });
  };

  const handleSetValues = useCallback(
    debounce((operator: string, searchValue: string) => {
      setFilter(operator, searchValue);
    }, 800),
    [],
  );

  const onChange = (value: string) => {
    setSearchValue(value);
    handleSetValues(selectedOperator?.value as string, value);
  };

  const onOperatorChange = (option: OptionsType) => {
    setSelectedOperator(option);
    handleSetValues(option?.value as string, searchValue);
  };

  const onClear = () => {
    setSearchValue('');
    setFilter(selectedOperator?.value as string, '');
  };

  useOnClickOutside(ref, () => setIsConditionOptionsOpen(false));

  useEffect(() => {
    if (inputRef.current) {
      inputRef.current.focus();
    }
  }, [isOpen]);

  return (
    <div
      className={`px-2.5 py-2 min-w-[218px] border-0.5 border-GRAY_500 rounded-md bg-white shadow-tableFilterMenu ${className}`}
    >
      <div className='flex text-GRAY_600 items-center gap-1 w-full z-80 mb-2'>
        <div className='f-11-400 text-GRAY_700  whitespace-nowrap'>{label || camelCaseToNormalText(columnId)}</div>
        <div
          className='flex items-center gap-[2px] cursor-pointer relative select-none grow mr-2'
          onClick={() => setIsConditionOptionsOpen(!isOpen)}
        >
          <div className='f-11-500 text-BLUE_700 max-w-[110px] whitespace-nowrap text-ellipsis overflow-hidden'>
            {selectedOperator?.label || 'is equal to'}
          </div>
          <SvgSpriteLoader id='chevron-down' iconCategory={ICON_SPRITE_TYPES.ARROWS} height={12} width={12} />
          {isConditionOptionsOpen && (
            <div
              ref={ref}
              className='p-1 z-10 absolute top-full left-0 w-[256px] bg-white text-GRAY_900 border border-GRAY_400 shadow-tableFilterMenu rounded-md'
            >
              {SEARCH_FILTER_OPTIONS.map((option) => (
                <div
                  className='hover:bg-GRAY_100 f-12-500 py-2 px-2.5 rounded-md'
                  key={option.value}
                  onClick={() => onOperatorChange(option)}
                >
                  {option.label}
                </div>
              ))}
            </div>
          )}
        </div>
        <div className='flex justify-end text-GRAY_700 cursor-pointer'>
          <SvgSpriteLoader
            id='refresh-ccw-01'
            iconCategory={ICON_SPRITE_TYPES.ARROWS}
            height={14}
            width={14}
            onClick={onClear}
          />
        </div>
      </div>
      <div className='flex flex-col gap-2'>
        <Input
          inputRef={inputRef}
          size={SIZE_TYPES.XSMALL}
          value={searchValue}
          placeholder='type a value...'
          onChange={(e) => onChange(e.target.value)}
        />
      </div>
    </div>
  );
};

export default SearchFilterMenuItem;
