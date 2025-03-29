import React, { FC, useEffect, useState } from 'react';
import { RangeFocus } from 'react-date-range';
import {
  DATE_FILTER_CATEGORIES,
  DATE_FORMATS,
  DATE_RANGE_TABS,
  DATE_RANGE_TYPES,
  dateFilterValueType,
  DateRangeKeys,
  DateRangeValue,
  PERIODICITY_OPTIONS,
  RangeType,
} from 'constants/date.constants';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { format } from 'date-fns';
import { EventCallbackType, MenuItem, TAB_TYPES } from 'types/common/components';
import { MapAny, OptionsType } from 'types/commonTypes';
import { cn } from 'utils/common';
import { searchDateRange } from 'components/common/dateRangePicker/dateRangePicker.utils';
import { DateUnitTabDisplay } from 'components/common/dateRangePicker/DateUnitTabDisplay';
import { DisplayDates } from 'components/common/dateRangePicker/DisplayDates';
import { Tabs } from 'components/common/tabs/Tabs';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

interface DateRangeMenuProps {
  className?: string;
  eventCallback?: EventCallbackType;
  onCategorySelect: (value: OptionsType, updatedRange: dateFilterValueType) => void;
  onDateChange: (range: RangeType) => void;
  range: RangeType;
  onRangeChange: (value: MapAny, updatedCategory?: DATE_FILTER_CATEGORIES) => void;
  menuWrapperClassName?: string;
  disabled?: boolean;
  defaultCategoryValue: OptionsType;
  id: string;
  customRangeOptions?: OptionsType[];
  focusedRange: RangeFocus;
  onSetFocusedRange: (range: RangeFocus) => void;
  resetFilter: () => void;
  isSingle?: boolean;
  dateFormat?: string;
  customTabValues?: DATE_RANGE_TYPES[];
  disableFutureDate?: boolean;
  isPeriodicityEnabled?: boolean;
  selectedPeriodicity?: OptionsType;
  onPeriodicityChange?: (value: OptionsType) => void;
}

const DateRangeMenu: FC<DateRangeMenuProps> = ({
  range,
  onRangeChange,
  onDateChange,
  disabled = false,
  id,
  focusedRange,
  onSetFocusedRange,
  resetFilter,
  isSingle = false,
  dateFormat = DATE_FORMATS.ddMMMyyyy,
  customTabValues = [],
  disableFutureDate,
  isPeriodicityEnabled = false,
  selectedPeriodicity,
  onPeriodicityChange,
}) => {
  const [currentValueStart, setCurrentValueStart] = useState<DateRangeValue | null>(null);
  const [currentValueEnd, setCurrentValueEnd] = useState<DateRangeValue | null>(null);
  const [searchValue, setSearchValue] = useState<DateRangeValue | null>(null);

  const [focusedInput, setFocusedInput] = useState<DateRangeKeys>(DateRangeKeys.START_DATE);

  const dateRangeTabs = customTabValues?.length
    ? DATE_RANGE_TABS?.filter((tab) => {
        return customTabValues?.includes(tab?.value as DATE_RANGE_TYPES);
      })
    : DATE_RANGE_TABS;
  const [currentTab, setCurrentTab] = useState<string>(dateRangeTabs?.[0].value);

  const [startDateDisplay, setStartDateDisplay] = useState<string>(
    range?.startDate ? format(range?.startDate, dateFormat) : '',
  );
  const [endDateDisplay, setEndDateDisplay] = useState<string>(
    range?.endDate ? format(range?.endDate, dateFormat) : '',
  );

  useEffect(() => {
    if (focusedRange[1] === 1) {
      onSetFocusedRange([0, 0]);
    }
  }, [currentTab]);

  useEffect(() => {
    if (!range) {
      return;
    }

    if (range.startDate) {
      setCurrentValueStart({
        type: DATE_RANGE_TYPES.DAY,
        year: range?.startDate?.getFullYear(),
        value: range?.startDate,
        label: range?.startDate?.toDateString(),
      });
    }

    if (range.endDate) {
      setCurrentValueEnd({
        type: DATE_RANGE_TYPES.DAY,
        year: range?.endDate?.getFullYear(),
        value: range?.endDate,
        label: range?.endDate?.toDateString(),
      });
    }

    setFocusedInput(DateRangeKeys.START_DATE);
  }, []);

  const handleRangeChange = (value: MapAny) => {
    onRangeChange(value, DATE_FILTER_CATEGORIES.CUSTOM_DATE_RANGE);

    if (
      focusedInput === DateRangeKeys.START_DATE &&
      value?.selection?.startDate.toString() === value?.selection?.endDate.toString()
    ) {
      setCurrentValueStart({
        type: DATE_RANGE_TYPES.DAY,
        year: value?.selection?.startDate?.getFullYear(),
        value: value?.selection?.startDate,
        label: value?.selection?.startDate?.toDateString(),
      });
      setCurrentValueEnd(null);
      setFocusedInput(DateRangeKeys.END_DATE);
      setStartDateDisplay(
        value?.selection?.startDate ? format(value.selection.startDate as Date, DATE_FORMATS.ddMMMyyyy) : '',
      );
      setEndDateDisplay('');
    } else {
      setCurrentValueStart({
        type: DATE_RANGE_TYPES.DAY,
        year: value?.selection?.startDate?.getFullYear(),
        value: value?.selection?.startDate,
        label: value?.selection?.startDate?.toDateString(),
      });
      setCurrentValueEnd({
        type: DATE_RANGE_TYPES.DAY,
        year: value?.selection?.endDate?.getFullYear(),
        value: value?.selection?.endDate,
        label: value?.selection?.endDate?.toDateString(),
      });
      setFocusedInput(DateRangeKeys.START_DATE);
      setStartDateDisplay(
        value?.selection?.startDate ? format(value.selection.startDate as Date, DATE_FORMATS.ddMMMyyyy) : '',
      );
      setEndDateDisplay(
        value?.selection?.endDate ? format(value?.selection?.endDate as Date, DATE_FORMATS.ddMMMyyyy) : '',
      );
    }
  };

  const handleTabSelect = (selected?: MenuItem) => {
    if (!selected) return;

    setCurrentTab(selected.value as string);
    setFocusedInput(DateRangeKeys.START_DATE);
  };

  const getValidDate = (value: Date) => {
    let date: Date | undefined = new Date(value.toString());

    date = date instanceof Date && !isNaN(date?.getTime()) ? date : undefined;

    return date;
  };

  const handleDateChange = (value: DateRangeValue, type: string) => {
    let update = new Date();
    const today = new Date();
    const newRange = { ...range };

    if (type === DateRangeKeys.START_DATE) {
      let endDateUpdate = new Date();

      switch (value?.type) {
        case DATE_RANGE_TYPES.MONTH:
          update.setMonth(value?.value as number);
          update.setDate(1);
          update.setFullYear(value?.year);

          endDateUpdate.setMonth((value?.value as number) + 1);
          endDateUpdate.setDate(0);
          endDateUpdate.setFullYear(value?.year);
          break;
        case DATE_RANGE_TYPES.QUARTER: {
          const endMonth = ((value?.value as number) + 1) * 3 - 3;

          update.setMonth(endMonth);
          update.setDate(1);
          update.setFullYear(value?.year);

          endDateUpdate.setMonth(endMonth + 3);
          endDateUpdate.setDate(0);
          endDateUpdate.setFullYear(value?.year);
          break;
        }
        case DATE_RANGE_TYPES.YEAR:
          update = new Date(value?.value as Date);
          update.setMonth(0);
          update.setDate(1);

          endDateUpdate = new Date(value?.value as Date);
          endDateUpdate.setMonth(11);
          endDateUpdate.setDate(31);
          break;

        case DATE_RANGE_TYPES.DAY:
          update = new Date(value?.value as Date);
          break;
      }

      newRange.startDate = getValidDate(disableFutureDate && update > today ? today : update);
      newRange.endDate = getValidDate(disableFutureDate && endDateUpdate > today ? today : endDateUpdate);

      setStartDateDisplay(
        format(disableFutureDate && update > today ? today : (update as Date), DATE_FORMATS.ddMMMyyyy),
      );
      setEndDateDisplay('');
    } else {
      switch (value?.type) {
        case DATE_RANGE_TYPES.MONTH:
          update.setMonth((value?.value as number) + 1);
          update.setDate(0);
          update.setFullYear(value?.year);

          break;
        case DATE_RANGE_TYPES.QUARTER: {
          const endMonth = ((value?.value as number) + 1) * 3;

          update.setMonth(endMonth);
          update.setDate(0);
          update.setFullYear(value?.year);
          break;
        }
        case DATE_RANGE_TYPES.YEAR:
          update = new Date(value?.value as Date);
          update.setFullYear(update.getFullYear() + 1);
          update.setDate(update.getDate() - 1);
          break;

        case DATE_RANGE_TYPES.DAY:
          update = new Date(value?.value as Date);
          break;
      }

      newRange.endDate = getValidDate(disableFutureDate && update > today ? today : update);
      setEndDateDisplay(format(disableFutureDate && update > today ? today : (update as Date), DATE_FORMATS.ddMMMyyyy));
      setFocusedInput(DateRangeKeys.START_DATE);
    }

    onDateChange(newRange);
  };

  const onClearSearch = () => {
    setSearchValue(null);
  };

  const onApplySearchValue = () => {
    if (!searchValue) return;

    if (isSingle && searchValue?.type !== DATE_RANGE_TYPES.DAY) {
      return;
    }

    onSetCurrentValue(searchValue);
    onClearSearch();
  };

  /**
   * Sets the value of the current date range based on the year, quarter or month tabs
   *
   */
  const onSetCurrentValue = (value: DateRangeValue) => {
    if (focusedInput === DateRangeKeys.START_DATE) {
      setCurrentValueStart(value);
      setCurrentValueEnd(null);

      handleDateChange(value, DateRangeKeys.START_DATE);

      setFocusedInput(DateRangeKeys.END_DATE);
      onSetFocusedRange([0, 1]);
    } else {
      if (currentValueStart && value.year <= currentValueStart?.year && value.value < currentValueStart?.value) {
        setCurrentValueStart(value);
        setCurrentValueEnd(currentValueStart);
        handleDateChange(value, DateRangeKeys.START_DATE);
        handleDateChange(currentValueStart, DateRangeKeys.END_DATE);

        onSetFocusedRange([0, 0]);
      } else {
        setCurrentValueEnd(value);

        handleDateChange(value, DateRangeKeys.END_DATE);

        onSetFocusedRange([0, 0]);
      }
    }

    onClearSearch();
  };

  const handleSearchChange = (searchInput: string) => {
    const value = searchDateRange(searchInput);

    if (!value) {
      setSearchValue(null);

      return;
    }

    setSearchValue(value);

    if (!isSingle) {
      setCurrentTab(value?.type);
    }
  };

  const onResetFilter = () => {
    setEndDateDisplay('');
    setStartDateDisplay('');
    setCurrentValueStart(null);
    setCurrentValueEnd(null);

    resetFilter();
  };

  return (
    <div className='h-full'>
      <div className='flex  overflow-hidden h-full'>
        <div className='flex-1 shadow-dateContainer w-full' data-testid={`date-range-menu-custom-${id}`}>
          {!isSingle && (
            <div className='border-b border-GRAY_400 mx-3 pt-3 flex w-auto justify-between items-center'>
              <div className=''>
                <Tabs
                  customSelectedIndex={dateRangeTabs?.findIndex((tab) => tab.value === currentTab)}
                  list={dateRangeTabs}
                  onSelect={handleTabSelect}
                  wrapperStyle='border-white !w-auto'
                  tabItemWrapperStyle='!w-auto'
                  id='ACCOUNTS_TABS'
                  scrollWrapperClassName='pb-0'
                  type={TAB_TYPES.UNDERLINE}
                />
              </div>
              <div className='cursor-pointer' onClick={onResetFilter} data-testid='reset-date-range'>
                <SvgSpriteLoader id='refresh-ccw-01' iconCategory={ICON_SPRITE_TYPES.ARROWS} width={14} height={14} />
              </div>
            </div>
          )}
          <div
            className={cn(
              `items-start flex flex-col`,
              isSingle ? 'h-full' : isPeriodicityEnabled ? 'h-[calc(100%-44px)]' : 'h-[calc(100%-38.5px)]',
            )}
          >
            {/* ----- Display value of startDate and endDate -------- */}
            <DisplayDates
              startDate={startDateDisplay}
              endDate={endDateDisplay}
              onChange={handleSearchChange}
              setFocusedInput={(type: DateRangeKeys) => setFocusedInput(type)}
              onApply={onApplySearchValue}
              focusedInput={focusedInput}
              currentTab={currentTab}
              isSingle={isSingle}
            />
            {/* ----- Tabs based on unit of time -------- */}
            <DateUnitTabDisplay
              currentTab={currentTab}
              onSetCurrentValue={onSetCurrentValue}
              currentValueStart={currentValueStart}
              currentValueEnd={currentValueEnd}
              searchValue={searchValue}
              handleRangeChange={handleRangeChange}
              disabled={disabled}
              id={id}
              range={range}
              focusedRange={focusedRange}
              focusedInput={focusedInput}
              disableFutureDate={disableFutureDate}
            />
            {isPeriodicityEnabled && (
              <div className='border-t border-GRAY_400 p-3 pb-4'>
                <div className='f-13-500 mb-1.5'>Periodicity</div>
                <div className='flex items-center gap-1.5 flex-wrap'>
                  {PERIODICITY_OPTIONS.map((item) => (
                    <div
                      className={cn(
                        'f-13-500 border rounded px-2 py-1 f-12-400 cursor-pointer',
                        selectedPeriodicity?.value === item.value ? 'bg-BG_GRAY_2 border-GRAY_500' : 'border-GRAY_400',
                      )}
                      key={item?.value}
                      onClick={() => onPeriodicityChange?.(item)}
                    >
                      {item?.label}
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default DateRangeMenu;
