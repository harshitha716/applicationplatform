import React, { useEffect, useState } from 'react';
import { RangeFocus } from 'react-date-range';
import {
  DATE_FILTER_CATEGORIES,
  DATE_FILTER_OPTIONS,
  DATE_RANGE_TYPES,
  dateFilterValueType,
  PERIODICITY_OPTIONS,
  PERIODICITY_TYPES,
  RangeType,
} from 'constants/date.constants';
import { EventCallbackType } from 'types/common/components';
import { MapAny, OptionsType } from 'types/commonTypes';
import { cn } from 'utils/common';
import DateRangeMenu from 'components/common/dateRangePicker/DateRangeMenu';
import { MenuWrapper } from 'components/common/MenuWrapper';

interface DateFilterProps {
  onChange: (category: string, value: dateFilterValueType) => void;
  defaultValue: { category: string } & dateFilterValueType;
  className?: string;
  id: string;
  menuWrapperClassName?: string;
  disabled?: boolean;
  customRangeOptions?: OptionsType[];
  isMenuOpen?: boolean;
  isSingle?: boolean;
  dateFormat?: string;
  customTabValues?: DATE_RANGE_TYPES[];
  disableFutureDate?: boolean;
  isPeriodicityEnabled?: boolean;
}

export interface DateRangeComponentProps {
  className?: string;
  eventCallback?: EventCallbackType;
  onCategorySelect: (value: OptionsType, updatedRange?: dateFilterValueType) => void;
  defaultCategoryValue: OptionsType;
  selectedCategory: DATE_FILTER_CATEGORIES;
  onDateChange: (changedDate: Date, type: string) => void;
  range: RangeType;
  onRangeChange: (value: MapAny) => void;
  menuWrapperClassName?: string;
  id: string;
  disabled?: boolean;
  customRangeOptions?: OptionsType[];
  dateFormat?: string;
}

const DateRangePicker: React.FC<DateFilterProps> = ({
  onChange,
  defaultValue,
  id,
  disabled = false,
  customRangeOptions,
  isSingle = false, // If only single date select is allowed
  dateFormat,
  customTabValues,
  disableFutureDate,
  isPeriodicityEnabled,
}) => {
  const [selectedCategory, setSelectedCategory] = useState<DATE_FILTER_CATEGORIES>(
    defaultValue?.category ? (defaultValue?.category as DATE_FILTER_CATEGORIES) : DATE_FILTER_CATEGORIES.ALL_TIME,
  );
  const [selectedPeriodicity, setSelectedPeriodicity] = useState(PERIODICITY_OPTIONS[0]);
  const [range, setRange] = useState<RangeType>({
    startDate: defaultValue?.start ?? undefined,
    endDate: defaultValue?.end ?? undefined,
    key: 'selection',
    showDateDisplay: false,
  });

  const dateFilterOptions = customRangeOptions ?? DATE_FILTER_OPTIONS;

  const handleCategorySelect = (value: OptionsType, updatedRange: dateFilterValueType) => {
    const typedValue = value?.value as DATE_FILTER_CATEGORIES;

    setSelectedCategory(typedValue);
    onChange(typedValue, updatedRange);
  };

  const [focusedRange, setFocusedRange] = useState<RangeFocus>([0, 0]);

  const onSetFocusedRange = (range: RangeFocus) => {
    setFocusedRange(range);
  };

  const resetFilter = () => {
    setSelectedCategory(DATE_FILTER_CATEGORIES.ALL_TIME);
    setRange({
      startDate: undefined,
      endDate: undefined,
      key: 'selection',
      showDateDisplay: false,
    });

    onChange(DATE_FILTER_CATEGORIES.ALL_TIME, {
      start: undefined,
      end: undefined,
    });

    if (focusedRange[1] !== 0) {
      onSetFocusedRange([0, 0]);
    }
  };

  const handleRangeChange = (value: MapAny, updatedCategory?: DATE_FILTER_CATEGORIES) => {
    setRange(value.selection);

    /** Sets whether range selection is in progress
     * 0, 0 means no range selection is in progress
     * 0, 1 means range selection is in progress
     */
    if (!isSingle) {
      onSetFocusedRange(focusedRange[1] === 0 ? [0, 1] : [0, 0]);

      if (focusedRange[1] === 0) {
        return;
      }
    }

    const updatedDateRange: dateFilterValueType = {
      start: value.selection.startDate,
      end: value.selection.endDate,
      periodicity: selectedPeriodicity?.value as PERIODICITY_TYPES,
    };

    if (updatedCategory) {
      setSelectedCategory(updatedCategory);
    }
    onChange(updatedCategory ?? selectedCategory, updatedDateRange);
  };

  const handleDateChange = (range: RangeType) => {
    setSelectedCategory(DATE_FILTER_CATEGORIES.CUSTOM_DATE_RANGE);
    setRange(range);
    const updatedDateRange: dateFilterValueType = {
      start: range.startDate,
      end: range.endDate,
      periodicity: selectedPeriodicity?.value as PERIODICITY_TYPES,
    };

    onChange(DATE_FILTER_CATEGORIES.CUSTOM_DATE_RANGE, updatedDateRange);
  };

  const handlePeriodicityChange = (value: OptionsType) => {
    setSelectedPeriodicity(value);

    const updatedDateRange: dateFilterValueType = {
      start: range.startDate,
      end: range.endDate,
      periodicity: value?.value as PERIODICITY_TYPES,
    };

    onChange(DATE_FILTER_CATEGORIES.CUSTOM_DATE_RANGE, updatedDateRange);
  };

  const defaultCategoryValue =
    dateFilterOptions?.find((each) => each.value === selectedCategory) ?? dateFilterOptions?.[0];

  useEffect(() => {
    setSelectedCategory(defaultValue?.category as DATE_FILTER_CATEGORIES);

    setRange({
      startDate: defaultValue?.start ?? new Date(),
      endDate: defaultValue?.end ?? new Date(),
      key: 'selection',
      showDateDisplay: false,
    });
  }, [defaultValue]);

  return (
    <MenuWrapper
      id={`${id}_DATE_RANGE_FILTER`}
      childrenWrapperClassName={cn(
        '!overflow-visible !w-[284px]',
        isSingle
          ? '!h-[390px] !max-h-125'
          : isPeriodicityEnabled
            ? '!h-[590px] !max-h-[590px]'
            : '!h-[480px] !max-h-125',
      )}
    >
      <DateRangeMenu
        onDateChange={handleDateChange}
        range={range}
        onSetFocusedRange={onSetFocusedRange}
        onRangeChange={handleRangeChange}
        onCategorySelect={handleCategorySelect}
        disabled={disabled}
        defaultCategoryValue={defaultCategoryValue}
        id={id}
        customRangeOptions={dateFilterOptions}
        focusedRange={focusedRange}
        resetFilter={resetFilter}
        isSingle={isSingle}
        dateFormat={dateFormat}
        customTabValues={customTabValues}
        disableFutureDate={disableFutureDate}
        isPeriodicityEnabled={isPeriodicityEnabled}
        selectedPeriodicity={selectedPeriodicity}
        onPeriodicityChange={handlePeriodicityChange}
      />
    </MenuWrapper>
  );
};

export default DateRangePicker;
