import React, { FC } from 'react';
import {
  DATE_FILTER_CATEGORIES,
  DATE_RANGE_TYPES,
  dateFilterValueType,
  PERIODICITY_TYPES,
} from 'constants/date.constants';
import { MapAny, OptionsType } from 'types/commonTypes';
import DateRangePicker from 'components/common/dateRangePicker';

export type DateFilterValueType = {
  start_date: Date | null;
  end_date: Date | null;
  date_category?: DATE_FILTER_CATEGORIES;
  periodicity?: PERIODICITY_TYPES;
};

const DateRangeFilter: FC<{
  value?: DateFilterValueType;
  disabled?: boolean;
  onChange?: (value: DateFilterValueType) => void;
  id?: string;
  menuWrapperClassName?: string;
  customRangeOptions?: OptionsType[];
  controlClassName?: string;
  withControl?: boolean;
  isSingle?: boolean;
  customTabValues?: DATE_RANGE_TYPES[];
  disableFutureDate?: boolean;
  isPeriodicityEnabled?: boolean;
}> = ({
  value,
  disabled = false,
  onChange,
  id,
  menuWrapperClassName = 'right-0',
  customRangeOptions,
  isSingle = false,
  customTabValues,
  disableFutureDate,
  isPeriodicityEnabled,
}) => {
  const handleChange = (category: string, value: dateFilterValueType) => {
    const customCondition = category === DATE_FILTER_CATEGORIES.CUSTOM_DATE_RANGE && value?.start && value?.end;
    const allCondition = category === DATE_FILTER_CATEGORIES.ALL_TIME && !value?.start && !value?.end;

    if (
      customCondition ||
      allCondition ||
      (category &&
        category !== DATE_FILTER_CATEGORIES.CUSTOM_DATE_RANGE &&
        category !== DATE_FILTER_CATEGORIES.ALL_TIME)
    ) {
      onChange?.({
        date_category: category as DATE_FILTER_CATEGORIES,
        start_date: value?.start || null,
        end_date: value?.end || null,
        periodicity: value?.periodicity,
      });
    }
  };

  const dateRangeProps: MapAny = {};

  if (customRangeOptions) {
    dateRangeProps.customRangeOptions = customRangeOptions;
  }

  const defaultValue = {
    start: value?.start_date ?? undefined,
    end: value?.end_date ?? undefined,
    category: value?.date_category || DATE_FILTER_CATEGORIES.CUSTOM_DATE_RANGE,
  };

  return (
    <>
      <DateRangePicker
        onChange={handleChange}
        menuWrapperClassName={menuWrapperClassName}
        defaultValue={defaultValue}
        id={`${id}_DATE_RANGE_FILTER`}
        disabled={disabled}
        isSingle={isSingle}
        customTabValues={customTabValues}
        disableFutureDate={disableFutureDate}
        isPeriodicityEnabled={isPeriodicityEnabled}
        {...dateRangeProps}
      />
    </>
  );
};

export default DateRangeFilter;
