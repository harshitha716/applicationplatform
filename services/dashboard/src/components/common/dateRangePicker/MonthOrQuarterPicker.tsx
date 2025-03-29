import React, { useEffect, useState } from 'react';
import { RangeFocus } from 'react-date-range';
import {
  DATE_RANGE_TYPES,
  DateRangeKeys,
  DateRangeValue,
  MonthsConfig,
  QuartersConfig,
} from 'constants/date.constants';
import { MapAny } from 'types/commonTypes';
import { getYearList } from 'components/common/dateRangePicker/dateRangePicker.utils';

interface MonthOrQuarterPickerProps {
  onSelect: (value: DateRangeValue) => void;
  currentValueStart: DateRangeValue | null;
  currentValueEnd: DateRangeValue | null;
  searchValue: DateRangeValue | null;
  type: DATE_RANGE_TYPES;
  focusedRange: RangeFocus;
  focusedInput?: DateRangeKeys;
}

export const MonthOrQuarterPicker: React.FC<MonthOrQuarterPickerProps> = ({
  onSelect,
  currentValueStart,
  currentValueEnd,
  type,
  searchValue,
  focusedRange,
  focusedInput,
}) => {
  const yearsList = getYearList();
  const containerRef = React.useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!searchValue) return;

    const container = containerRef.current;

    if (!container) return;

    const selectedYear = container.querySelector('.selected-year');

    if (!selectedYear) return;

    selectedYear.scrollIntoView({
      behavior: 'smooth',
      block: 'center',
    });
  }, [searchValue]);

  const onSelectValue = ({ year, config }: MapAny) => {
    const value: DateRangeValue = {
      type,
      year,
      value: config?.value,
      label: config?.short,
    };

    onSelect(value);
  };

  const isSelected = (year: number, value: number) => {
    return (
      (year === currentValueStart?.year && value === currentValueStart?.value && type === currentValueStart?.type) ||
      (year === currentValueEnd?.year && value === currentValueEnd?.value && type === currentValueEnd?.type)
    );
  };

  const isPartiallySelected = (year: number, value: number) => {
    return value === searchValue?.value && year === searchValue?.year && type === searchValue?.type;
  };

  const config = type === DATE_RANGE_TYPES.MONTH ? MonthsConfig : QuartersConfig;

  // ---------- Range selection-----------------
  const isRangeSelectionInProgress = focusedRange[1] === 1;

  const [lastHoveredValue, setLastHoveredValue] = useState<DateRangeValue | null>(null);

  const onMouseEnter = (year: number, config: MapAny) => {
    if (!isRangeSelectionInProgress) return;

    const value: DateRangeValue = {
      type,
      year,
      value: config.value,
      label: config.short,
    };

    setLastHoveredValue(value);
  };

  useEffect(() => {
    if (focusedRange[1] === 0) {
      setLastHoveredValue(null);
    }
  }, [focusedRange]);

  const isWithinRange = (startRange: DateRangeValue, endRange: DateRangeValue, year: number, curentConfig: MapAny) => {
    const minYear = startRange?.year <= endRange?.year ? startRange : endRange;
    const maxYear = startRange?.year > endRange?.year ? startRange : endRange;

    /**
     * config.year must fall between startRange.year and endRange.year and if the year is same then must fall between
     * startRange.value and endRange.value.
     *  */
    if (year > minYear?.year && year < maxYear?.year) {
      return true;
    }

    if (year === minYear?.year && year === maxYear?.year) {
      const minVal = Math.min(startRange?.value as number, endRange?.value as number);
      const maxVal = Math.max(startRange?.value as number, endRange?.value as number);

      return curentConfig.value >= minVal && curentConfig.value <= maxVal;
    }

    if (year === minYear?.year) {
      return curentConfig.value >= minYear?.value;
    }

    if (year === maxYear?.year) {
      return curentConfig.value <= maxYear?.value;
    }

    return false;
  };

  const shouldHighlightCell = (year: number, config: MapAny) => {
    const isStartEqualToEnd =
      currentValueStart &&
      currentValueEnd &&
      currentValueStart?.year === currentValueEnd?.year &&
      currentValueStart?.value === currentValueEnd?.value;

    /**
     * When range is already selected
     */
    if (!isStartEqualToEnd && currentValueEnd && currentValueStart) {
      if (currentValueEnd?.type !== type || currentValueStart?.type !== currentValueEnd?.type) return false;

      return isWithinRange(currentValueStart, currentValueEnd, year, config);
    }

    if (!lastHoveredValue) return false;

    /**
     * When range selection is in progress
     */
    const currentValue = focusedInput === DateRangeKeys.START_DATE ? currentValueEnd : currentValueStart;

    if (!currentValue || !lastHoveredValue) return false;

    return isWithinRange(lastHoveredValue, currentValue, year, config);
  };

  return (
    <div className='' ref={containerRef} onMouseLeave={() => setLastHoveredValue(null)}>
      {yearsList.map((year, index) => {
        return (
          <div className=' flex flex-col items-start' key={index}>
            <div className='f-12-500 mb-2 text-GRAY_1000'>{year}</div>
            <div className=' flex flex-wrap gap-2 mb-4'>
              {config.map((config, index) => {
                return (
                  <div
                    key={index}
                    className={` ${
                      isSelected(year, config?.value)
                        ? 'bg-BLUE_700 text-white border-GRAY_400'
                        : shouldHighlightCell(year, config)
                          ? 'bg-BLUE_50 '
                          : isPartiallySelected(year, config?.value)
                            ? 'border-BLUE_700 selected-year'
                            : 'hover:border-BLUE_700 hover:selected-year bg-BG_GRAY_2'
                    }  cursor-pointer f-12-500 w-13.5 flex items-center justify-center py-[4.5px]  rounded-sm border `}
                    onClick={() => onSelectValue({ year, config: config })}
                    onMouseEnter={() => onMouseEnter(year, config)}
                  >
                    {config?.short}
                  </div>
                );
              })}
            </div>
          </div>
        );
      })}
    </div>
  );
};
