import React, { useEffect, useRef } from 'react';
import { DateRangePicker, RangeFocus } from 'react-date-range';
import { DATE_RANGE_TYPES, DateRangeValue, RangeType } from 'constants/date.constants';
import { DateRangePickerNavigator } from 'components/common/dateRangePicker/DateRangePickerNavigator';

interface DateRangePickerWrapperProps {
  searchValue: DateRangeValue | null;
  disabled: boolean;
  id: string;
  onRangeChange: (range: any) => void;
  range: RangeType;
  focusedRange: RangeFocus;
  disableFutureDate?: boolean;
}

export const DateRangePickerWrapper: React.FC<DateRangePickerWrapperProps> = ({
  searchValue,
  disabled,
  id,
  onRangeChange,
  range,
  focusedRange,
  disableFutureDate,
}) => {
  const renderDayContent = (day: Date) => {
    const currentShownMonth =
      dateRangePickerRef?.current?.dateRange?.calendar?.state?.focusedDate?.getMonth() ?? day.getMonth();

    const isNotInCurrentMonth = day.getMonth() !== currentShownMonth;

    const dayString = day.toDateString();

    const isStart = range?.startDate?.toDateString() === dayString;
    const isEnd = range?.endDate?.toDateString() === dayString;
    const isSearchValue =
      searchValue?.type === DATE_RANGE_TYPES.DAY && (searchValue?.value as Date)?.toDateString() === dayString;
    const isToday = new Date().toDateString() === dayString;
    const isFutureDate = day > new Date();

    return (
      <div
        className={`w-full h-full flex justify-center items-center rounded-full  ${
          isStart ? 'bg-BLUE_700 !text-white' : ''
        } ${isEnd ? 'bg-BLUE_700 !text-white' : ''}
            ${isSearchValue ? 'border-BLUE_700 border is-searched' : ''}
            ${!isStart && !isEnd && !isSearchValue && isToday ? 'border border-DIVIDER_SAIL_2' : ''}
            ${isNotInCurrentMonth || (disableFutureDate && isFutureDate) ? ' text-GRAY_500' : 'text-black'}`}
      >
        {day.getDate()}
      </div>
    );
  };

  const dateRangePickerRef = useRef<any>(null);

  useEffect(() => {
    if (searchValue?.type !== DATE_RANGE_TYPES.DAY || !dateRangePickerRef?.current) return;

    dateRangePickerRef?.current?.dateRange?.calendar?.focusToDate(searchValue?.value as Date);
  }, [searchValue, dateRangePickerRef?.current]);

  return (
    <div className={disabled ? 'pointer-events-none' : ''} data-testid={`date-range-menu-custom-picker-wrapper-${id}`}>
      <DateRangePicker
        ref={dateRangePickerRef}
        ranges={[range]}
        showMonthAndYearPickers={false}
        onChange={onRangeChange}
        focusedRange={focusedRange}
        showDateDisplay={false}
        rangeColors={[]}
        staticRanges={[]}
        inputRanges={[]}
        showMonthArrow={false}
        direction='horizontal'
        dayContentRenderer={renderDayContent}
        className='!w-[258px]'
        navigatorRenderer={DateRangePickerNavigator}
        fixedHeight={true}
        maxDate={disableFutureDate ? new Date() : undefined}
      />
    </div>
  );
};
