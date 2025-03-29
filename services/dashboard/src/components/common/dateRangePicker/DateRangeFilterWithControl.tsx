import React, { FC, useRef, useState } from 'react';
import { DATE_RANGE_TYPES } from 'constants/date.constants';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { useOnClickOutside } from 'hooks';
import { OptionsType } from 'types/commonTypes';
import { getDateRangeTitle } from 'utils/common';
import DateRangeFilter, { DateFilterValueType } from 'components/filter/DateRangeFilter';
import FilterControlButton from 'components/filter/FilterControlButton';

interface DateRangeFilterWithControlProps {
  value: DateFilterValueType;
  onChange: (value: DateFilterValueType) => void;
  disabled?: boolean;
  className?: string;
  controlClassName?: string;
  customRangeOptions?: OptionsType[];
  isSingle?: boolean;
  showSingleDate?: boolean;
  customTabValues?: DATE_RANGE_TYPES[];
  id?: string;
  disableFutureDate?: boolean;
}

const DateRangeFilterWithControl: FC<DateRangeFilterWithControlProps> = ({
  value,
  onChange,
  disabled = false,
  className = '',
  controlClassName = '',
  customRangeOptions,
  isSingle = false,
  showSingleDate = false,
  customTabValues,
  id = '',
  disableFutureDate,
}) => {
  const [isDateRangeOpen, setIsDateRangeOpen] = useState(false);
  const dateRangeRef = useRef<HTMLDivElement>(null);

  useOnClickOutside(dateRangeRef, () => {
    setIsDateRangeOpen(false);
  });

  const onToggleDateRange = () => {
    setIsDateRangeOpen(!isDateRangeOpen);
  };

  return (
    <div ref={dateRangeRef} className={`relative ${className}`}>
      <FilterControlButton
        onClick={onToggleDateRange}
        icon='calendar'
        id={`DATE_RANGE_BUTTON_${id}`}
        iconCategory={ICON_SPRITE_TYPES.TIME}
        isSelected={isDateRangeOpen}
        childrenWrapperClassName='!text-xs'
        className={`!mb-0 ${controlClassName}`}
      >
        {getDateRangeTitle(value, showSingleDate)}
      </FilterControlButton>

      {isDateRangeOpen ? (
        <div className='absolute top-10 z-999 right-0'>
          <DateRangeFilter
            onChange={onChange}
            value={value}
            id={`DATE_RANGE_FILTER_${id}`}
            disabled={disabled}
            controlClassName='px-2 py-1.5 border-DIVIDER_SAIL_2 rounded-lg h-auto mr-3 cursor-pointer'
            customRangeOptions={customRangeOptions}
            isSingle={isSingle}
            customTabValues={customTabValues}
            disableFutureDate={disableFutureDate}
          />
        </div>
      ) : null}
    </div>
  );
};

export default DateRangeFilterWithControl;
