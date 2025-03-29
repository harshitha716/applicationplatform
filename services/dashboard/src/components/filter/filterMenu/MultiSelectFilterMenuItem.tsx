import React, { ChangeEvent, FC, useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { SIZE_TYPES } from 'types/common/components';
import { OptionsType } from 'types/commonTypes';
import { camelCaseToNormalText, cn, debounce } from 'utils/common';
import { CheckBox } from 'components/common/Checkbox';
import Input from 'components/common/input';
import { Tooltip } from 'components/common/tooltip';
import { FILTER_TYPES } from 'components/filter/filter.types';
import { CONDITION_OPERATOR_TYPE, MULTI_SELECT_FILTER_OPTIONS } from 'components/filter/filters.constants';
import { filtersContextActions, useFiltersContextStore } from 'components/filter/filters.context';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

interface MultiSelectFilterMenuItemProps {
  column: { colId: string };
  values: string[];
  className?: string;
  LabelComponent?: (item: string) => React.ReactNode;
  operatorOptions?: OptionsType[];
  isOpen?: boolean;
  showSelectAll?: boolean;
  label?: string;
}

const MultiSelectFilterMenuItem: FC<MultiSelectFilterMenuItemProps> = ({
  column,
  values,
  className,
  LabelComponent,
  operatorOptions = MULTI_SELECT_FILTER_OPTIONS,
  isOpen = false,
  showSelectAll = false,
  label,
}) => {
  const inputRef = useRef<HTMLInputElement>(null);
  const listRef = useRef<HTMLDivElement>(null);
  const [hasScrolled, setHasScrolled] = useState(false);

  const columnId = column?.colId;
  const {
    state: { selectedFilters },
    dispatch,
  } = useFiltersContextStore();

  const currentOperator =
    operatorOptions.find((option) => option.value === selectedFilters[columnId]?.type) || operatorOptions[0];
  const [selectedValues, setSelectedValues] = useState<string[]>(selectedFilters[columnId]?.values || []);
  const [inputValue, setInputValue] = useState('');
  const [isConditionOpen, setIsConditionOpen] = useState(false);
  const [selectedOperator, setSelectedOperator] = useState<OptionsType>(currentOperator);
  const [isSelectAll, setIsSelectAll] = useState(false);

  const setFilter = useCallback(
    (operator: string, updatedValues: string[]) => {
      dispatch({
        type: filtersContextActions.SET_SELECTED_FILTERS,
        payload: {
          selectedFilters: {
            [columnId]: {
              filterType: FILTER_TYPES.MULTI_SELECT,
              type: operator,
              values: updatedValues,
            },
          },
        },
      });
    },
    [dispatch, columnId],
  );

  const handleSetValues = useCallback(
    debounce((operator: string, updatedValues: string[]) => {
      setFilter(operator, updatedValues);
    }, 800),
    [setFilter],
  );

  const onSearchChange = (event: ChangeEvent<HTMLInputElement>) => {
    setInputValue(event.target.value);
  };

  const onChange = (value: string) => {
    const updatedValues = selectedValues.includes(value)
      ? selectedValues.filter((item) => item !== value)
      : [...selectedValues, value];

    setSelectedValues(updatedValues);
    handleSetValues(selectedOperator.value as string, updatedValues);
  };

  const onReset = () => {
    setSelectedValues([]);
    setInputValue('');
    setIsSelectAll(false);
    setFilter(selectedOperator.value as string, []);
  };

  const onOperatorChange = (option: OptionsType) => {
    setSelectedOperator(option);
    const newValues = option.value === CONDITION_OPERATOR_TYPE.IS_NULL ? [] : selectedValues;

    setSelectedValues(newValues);
    handleSetValues(option.value as string, newValues);
  };

  const handleScroll = () => {
    if (listRef.current) {
      setHasScrolled(listRef.current.scrollTop > 0);
    }
  };

  useEffect(() => {
    if (inputRef.current && isOpen) {
      inputRef.current.focus();
    }
  }, [isOpen]);

  const filteredValues = useMemo(() => {
    const lowerCasedInput = inputValue.toLowerCase();

    return values.filter((item) => item && item.toLowerCase().includes(lowerCasedInput));
  }, [values, inputValue]);

  useEffect(() => {
    setIsSelectAll(filteredValues.length > 0 && filteredValues.every((item) => selectedValues.includes(item)));
  }, [filteredValues, selectedValues]);

  const onSelectAll = () => {
    const newSelectedValues = isSelectAll
      ? selectedValues.filter((val) => !filteredValues.includes(val))
      : Array.from(new Set([...selectedValues, ...filteredValues]));

    setSelectedValues(newSelectedValues);
    handleSetValues(selectedOperator.value as string, newSelectedValues);
  };

  return (
    <div
      className={cn(
        'flex flex-col gap-2 bg-white pt-2 pb-1 w-[218px] border-0.5 border-GRAY_500 rounded-md shadow-tableFilterMenu max-h-[330px] min-w-[230px]',
        className,
      )}
    >
      <div className='flex text-GRAY_600 items-center gap-1 w-full z-80 px-2.5'>
        <div className='f-11-400 text-GRAY_700 whitespace-nowrap text-ellipsis overflow-hidden'>
          {label || camelCaseToNormalText(columnId)}
        </div>
        <div
          className='flex items-center gap-[2px] cursor-pointer relative select-none grow'
          onClick={() => setIsConditionOpen(!isConditionOpen)}
        >
          <div className='f-11-500 text-BLUE_700 max-w-[110px] whitespace-nowrap text-ellipsis overflow-hidden'>
            {selectedOperator?.label || 'is equal to'}
          </div>
          <SvgSpriteLoader
            id='chevron-down'
            className={cn(
              'text-GRAY_700 transition-transform duration-300',
              isConditionOpen ? 'rotate-180' : 'rotate-0',
            )}
            height={12}
            width={12}
          />
          {isConditionOpen && (
            <div className='p-1 z-10 absolute top-full left-0 bg-white text-GRAY_900 border border-GRAY_400 shadow-tableFilterMenu rounded-md'>
              {operatorOptions.map((option) => (
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
            onClick={onReset}
          />
        </div>
      </div>
      <div className='px-2.5'>
        <Input
          size={SIZE_TYPES.XSMALL}
          inputRef={inputRef}
          value={inputValue}
          placeholder='type a value...'
          onChange={onSearchChange}
          autoFocus
        />
      </div>
      {showSelectAll && (
        <div
          onClick={() => onSelectAll()}
          className='flex items-center gap-2 justify-between mx-1 py-2 px-2.5 cursor-pointer select-none rounded hover:bg-GRAY_100'
        >
          <div className='f-12-400 text-GRAY_1000'>Select All</div>
          <div className='min-w-[14px]'>
            <CheckBox checked={isSelectAll} id='checkbox-1' />
          </div>
        </div>
      )}
      <div
        className={cn(
          'flex flex-col h-full overflow-y-auto overflow-x-hidden px-1 [&::-webkit-scrollbar]:hidden',
          hasScrolled && 'border-t border-GRAY_400',
        )}
        ref={listRef}
        onScroll={handleScroll}
      >
        {!!filteredValues?.length &&
          filteredValues.map((item) => (
            <div
              key={item}
              onClick={() => onChange(item)}
              className='flex items-center gap-2 justify-between py-2 px-2.5 cursor-pointer select-none rounded hover:bg-GRAY_100'
            >
              {LabelComponent ? LabelComponent(item) : <div className='f-12-400 text-GRAY_1000'>{String(item)}</div>}
              <Tooltip
                tooltipBody={`condition set to “${operatorOptions.find((option) => option.value === CONDITION_OPERATOR_TYPE.IS_NULL)?.label}”`}
                tooltipBodyClassName='f-12-300 px-3 py-1.5 rounded-md z-999 bg-black text-white w-28'
                className='z-1'
                disabled={selectedOperator?.value !== CONDITION_OPERATOR_TYPE.IS_NULL}
              >
                <div className='min-w-[14px]'>
                  <CheckBox
                    checked={selectedValues?.includes(item)}
                    id='checkbox-1'
                    disabled={selectedOperator?.value === CONDITION_OPERATOR_TYPE.IS_NULL}
                  />
                </div>
              </Tooltip>
            </div>
          ))}
      </div>
    </div>
  );
};

export default MultiSelectFilterMenuItem;
