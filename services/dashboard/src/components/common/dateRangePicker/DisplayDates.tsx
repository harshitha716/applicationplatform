import React, { useEffect } from 'react';
import { DateRangeKeys } from 'constants/date.constants';
import { defaultFnType } from 'types/commonTypes';
import { getFromLocalStorage, LOCAL_STORAGE_KEYS, setToLocalStorage } from 'utils/localstorage';
import { getPlacehoderDate } from 'components/common/dateRangePicker/dateRangePicker.utils';
import { DateSearch } from 'components/common/dateRangePicker/DateSearch';

interface DisplayDatesProps {
  setFocusedInput: (value: DateRangeKeys) => void;
  startDate: string;
  endDate: string;
  onChange: (value: string) => void;
  onApply: defaultFnType;
  focusedInput: DateRangeKeys;
  currentTab: string;
  isSingle: boolean;
}

export const DisplayDates: React.FC<DisplayDatesProps> = ({
  setFocusedInput,
  startDate,
  endDate,
  onChange,
  onApply,
  focusedInput,
  currentTab,
  isSingle,
}) => {
  const [startDateSearchValue, setStartDateSearchValue] = React.useState<string>(startDate);
  const [endDateSearchValue, setEndDateSearchValue] = React.useState<string>(endDate);

  const handleSearchChange = (value: string, key: DateRangeKeys) => {
    if (key === DateRangeKeys.START_DATE) {
      setStartDateSearchValue(value);
    } else {
      setEndDateSearchValue(value);
    }

    onChange(value);
  };

  const onApplySearchValue = (type: DateRangeKeys) => {
    if (type === DateRangeKeys.START_DATE) {
      setEndDateSearchValue('');
    }

    onApply();
  };

  useEffect(() => {
    setStartDateSearchValue(startDate);
    setEndDateSearchValue(endDate);
  }, [startDate, endDate]);

  const placeholderDate = getPlacehoderDate();

  const placeholderSeenCount = getFromLocalStorage(LOCAL_STORAGE_KEYS.DATE_PLACEHOLDER_SEEN) ?? 0;

  const shouldShowInfo = Number(placeholderSeenCount) <= 3;

  useEffect(() => {
    if (!shouldShowInfo) {
      return;
    }

    setToLocalStorage(LOCAL_STORAGE_KEYS.DATE_PLACEHOLDER_SEEN, `${Number(placeholderSeenCount) + 1}`);
  }, []);

  return (
    <div className='flex flex-col w-full px-3 pt-3 pb-4 '>
      <div className='mb-2 flex flex-col  items-start' onClick={() => setFocusedInput(DateRangeKeys.START_DATE)}>
        <div className='flex justify-between w-full'>
          <div className='f-12-400 mb-1.5 text-GRAY_500'>{isSingle ? 'Enter date' : 'From'} </div>
          {shouldShowInfo ? (
            <div className='f-11-300 mb-1.5 text-GRAY_500'>{`Try: Next month, Q4 or ${placeholderDate}`}</div>
          ) : null}
        </div>

        <DateSearch
          value={startDateSearchValue}
          onChange={(e) => handleSearchChange(e?.target?.value, DateRangeKeys.START_DATE)}
          onClear={() => handleSearchChange('', DateRangeKeys.START_DATE)}
          onApply={() => onApplySearchValue(DateRangeKeys.START_DATE)}
          id={DateRangeKeys.START_DATE}
          focusedInput={focusedInput}
          currentTab={currentTab}
          shouldShowInfo={shouldShowInfo}
        />
      </div>

      {!isSingle && (
        <div className='flex flex-col  items-start' onClick={() => setFocusedInput(DateRangeKeys.END_DATE)}>
          <div className='f-12-400 mb-1.5 text-GRAY_500'>To</div>
          <DateSearch
            value={endDateSearchValue}
            onChange={(e) => handleSearchChange(e?.target?.value, DateRangeKeys.END_DATE)}
            onApply={() => onApplySearchValue(DateRangeKeys.END_DATE)}
            id={DateRangeKeys.END_DATE}
            onClear={() => handleSearchChange('', DateRangeKeys.END_DATE)}
            focusedInput={focusedInput}
            currentTab={currentTab}
          />
        </div>
      )}
    </div>
  );
};
