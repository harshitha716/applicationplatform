import React, { FC, useCallback, useRef, useState } from 'react';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { debounce, useOnClickOutside } from 'hooks';
import { SIZE_TYPES } from 'types/common/components';
import { OptionsType } from 'types/commonTypes';
import { camelCaseToNormalText } from 'utils/common';
import Input from 'components/common/input';
import { Tooltip } from 'components/common/tooltip';
import { FILTER_TYPES } from 'components/filter/filter.types';
import { AMOUNT_RANGE_FILTER_OPTIONS, CONDITION_OPERATOR_TYPE } from 'components/filter/filters.constants';
import { filtersContextActions, useFiltersContextStore } from 'components/filter/filters.context';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

interface AmountRangeFilterMenuItemProps {
  column: { colId: string };
  values: string[];
  className?: string;
  label?: string;
}

const AmountRangeFilterMenuItem: FC<AmountRangeFilterMenuItemProps> = ({ column, className, label }) => {
  const ref = useRef(null);
  const columnId = column?.colId;
  const {
    state: { selectedFilters },
    dispatch,
  } = useFiltersContextStore();
  const currentOperator = AMOUNT_RANGE_FILTER_OPTIONS.find(
    (option) => option.value === selectedFilters[columnId]?.type,
  );
  const [startValue, setStartValue] = useState(selectedFilters[columnId]?.filter || '');
  const [endValue, setEndValue] = useState(selectedFilters[columnId]?.filterTo || '');
  const [isOpen, setIsOpen] = useState(false);
  const [selectedOperator, setSelectedOperator] = useState<OptionsType>(
    currentOperator ?? AMOUNT_RANGE_FILTER_OPTIONS[0],
  );

  const setFilter = (operator: string, startValue: string, endValue: string) => {
    const condition =
      operator === CONDITION_OPERATOR_TYPE.IN_BETWEEN
        ? endValue !== '' && startValue !== ''
        : operator === CONDITION_OPERATOR_TYPE.IS_NULL
          ? true
          : startValue !== '';

    dispatch({
      type: filtersContextActions.SET_SELECTED_FILTERS,
      payload: {
        selectedFilters: {
          [columnId]: condition
            ? {
                filterType: FILTER_TYPES.AMOUNT_RANGE,
                type: operator,
                filter: startValue,
                filterTo: endValue,
              }
            : {},
        },
      },
    });
  };

  const handleSetValues = useCallback(
    debounce((operator: string, startValue: string, endValue: string) => {
      setFilter(operator, startValue, endValue);
    }, 800),
    [],
  );

  const onChange = (isStart: boolean, value: string) => {
    if (isStart) setStartValue(value);
    else setEndValue(value);

    handleSetValues(selectedOperator?.value as string, isStart ? value : startValue, isStart ? endValue : value);
  };

  const onOperatorChange = (option: OptionsType) => {
    setSelectedOperator(option);
    handleSetValues(option?.value as string, startValue, endValue);
  };

  const onClear = () => {
    setStartValue('');
    setEndValue('');
    setFilter(selectedOperator?.value as string, '', '');
  };

  useOnClickOutside(ref, () => setIsOpen(false));

  return (
    <div
      className={`px-2.5 py-2 w-[250px] min-w-[250px] border-0.5 border-GRAY_500 rounded-md bg-white shadow-tableFilterMenu ${className}`}
    >
      <div className='flex text-GRAY_600 items-center gap-1 w-full z-80 mb-2'>
        <div className='f-11-400 text-GRAY_700 whitespace-nowrap text-ellipsis overflow-hidden'>
          {label || camelCaseToNormalText(columnId)}
        </div>
        <div
          className='flex items-center gap-[2px] cursor-pointer relative select-none grow mr-4'
          onClick={() => setIsOpen(!isOpen)}
        >
          <div className='f-11-500 text-BLUE_700 max-w-[110px] whitespace-nowrap text-ellipsis overflow-hidden'>
            {selectedOperator?.label || 'is equal to'}
          </div>
          <SvgSpriteLoader id='chevron-down' iconCategory={ICON_SPRITE_TYPES.ARROWS} height={12} width={12} />
          {isOpen && (
            <div
              ref={ref}
              className='p-1 z-10 absolute top-full left-0 w-[256px] bg-white text-GRAY_900 border border-GRAY_400 shadow-tableFilterMenu rounded-md'
            >
              {AMOUNT_RANGE_FILTER_OPTIONS.map((option) => (
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
        <Tooltip
          tooltipBody={`condition set to “is blank”`}
          tooltipBodyClassName='f-12-300 px-3 py-1.5 rounded-md whitespace-nowrap z-999 bg-black text-white'
          className='z-1 !cursor-not-allowed'
          disabled={selectedOperator?.value !== CONDITION_OPERATOR_TYPE.IS_NULL}
        >
          <Input
            size={SIZE_TYPES.XSMALL}
            value={startValue}
            placeholder='type a value...'
            onChange={(e) => onChange(true, e.target.value)}
            disabled={selectedOperator?.value === CONDITION_OPERATOR_TYPE.IS_NULL}
            autoFocus
          />
        </Tooltip>
        {selectedOperator?.value === CONDITION_OPERATOR_TYPE.IN_BETWEEN && (
          <span className='f-11-400 text-GRAY_700 select-none'>and</span>
        )}
        {selectedOperator?.value === CONDITION_OPERATOR_TYPE.IN_BETWEEN && (
          <Input
            size={SIZE_TYPES.XSMALL}
            value={endValue}
            placeholder='type a value...'
            onChange={(e) => onChange(false, e.target.value)}
          />
        )}
      </div>
    </div>
  );
};

export default AmountRangeFilterMenuItem;
