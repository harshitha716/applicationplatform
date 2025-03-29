import React, { useEffect, useRef } from 'react';
import { RangeFocus } from 'react-date-range';
import { DATE_RANGE_TYPES, DateRangeKeys, DateRangeValue, RangeType } from 'constants/date.constants';
import { DateRangePickerWrapper } from 'components/common/dateRangePicker/DateRangePickerWrapper';
import { MonthOrQuarterPicker } from 'components/common/dateRangePicker/MonthOrQuarterPicker';
import { YearPicker } from 'components/common/dateRangePicker/YearPicker';

interface DateUnitTabDisplayProps {
  currentValueStart: DateRangeValue | null;
  currentValueEnd: DateRangeValue | null;
  onSetCurrentValue: (value: DateRangeValue) => void;
  currentTab: string;
  searchValue: DateRangeValue | null;
  disabled: boolean;
  id: string;
  handleRangeChange: (range: any) => void;
  range: RangeType;
  focusedRange: RangeFocus;
  focusedInput?: DateRangeKeys;
  disableFutureDate?: boolean;
}

export const DateUnitTabDisplay: React.FC<DateUnitTabDisplayProps> = ({
  currentValueStart,
  currentValueEnd,
  onSetCurrentValue,
  currentTab,
  searchValue,
  disabled,
  id,
  handleRangeChange,
  range,
  focusedRange,
  focusedInput,
  disableFutureDate,
}) => {
  const dateRangePickerRef = useRef<any>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (searchValue?.type !== DATE_RANGE_TYPES.DAY || !dateRangePickerRef?.current) return;

    dateRangePickerRef?.current?.dateRange?.calendar?.focusToDate(searchValue?.value as Date);
  }, [searchValue, dateRangePickerRef?.current]);

  useEffect(() => {
    if (containerRef?.current) {
      containerRef.current.scrollTop = containerRef?.current?.scrollHeight;
    }
  }, [currentTab]);

  return (
    <div
      className='h-[calc(100%-65px)] pt-3 px-3 border-t border-GRAY_400 w-full overflow-y-auto date-unit-container'
      ref={containerRef}
    >
      {currentTab === DATE_RANGE_TYPES.MONTH ? (
        <MonthOrQuarterPicker
          onSelect={onSetCurrentValue}
          currentValueStart={currentValueStart}
          currentValueEnd={currentValueEnd}
          searchValue={searchValue}
          type={DATE_RANGE_TYPES.MONTH}
          focusedRange={focusedRange}
        />
      ) : currentTab === DATE_RANGE_TYPES.QUARTER ? (
        <MonthOrQuarterPicker
          onSelect={onSetCurrentValue}
          currentValueStart={currentValueStart}
          currentValueEnd={currentValueEnd}
          searchValue={searchValue}
          type={DATE_RANGE_TYPES.QUARTER}
          focusedRange={focusedRange}
          focusedInput={focusedInput}
        />
      ) : currentTab === DATE_RANGE_TYPES.YEAR ? (
        <YearPicker
          onSelect={onSetCurrentValue}
          currentValueStart={currentValueStart}
          currentValueEnd={currentValueEnd}
          searchValue={searchValue}
          focusedRange={focusedRange}
        />
      ) : (
        <DateRangePickerWrapper
          searchValue={searchValue}
          disabled={disabled}
          id={id}
          onRangeChange={handleRangeChange}
          range={range}
          focusedRange={focusedRange}
          disableFutureDate={disableFutureDate}
        />
      )}
    </div>
  );
};
