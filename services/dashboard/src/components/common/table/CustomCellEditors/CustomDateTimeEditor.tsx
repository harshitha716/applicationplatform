import React from 'react';
import { DATE_FILTER_CATEGORIES } from 'constants/date.constants';
import { MapAny } from 'types/commonTypes';
import DateRangeFilter, { DateFilterValueType } from 'components/filter/DateRangeFilter';

const CustomDateTimeEditor = (props: MapAny) => {
  const { value, onValueChange, stopEditing } = props;

  const onChange = (value: DateFilterValueType) => {
    onValueChange(value?.start_date?.toISOString());
    stopEditing();
  };

  return (
    <div className='fixed top-[210px]'>
      <DateRangeFilter
        onChange={onChange}
        value={{
          date_category: DATE_FILTER_CATEGORIES.CUSTOM_DATE_RANGE,
          start_date: new Date(value),
          end_date: new Date(value),
        }}
        disabled={false}
        controlClassName='px-2 py-1.5 border-DIVIDER_SAIL_2 rounded-lg h-auto mr-3 cursor-pointer'
        isSingle={true}
        disableFutureDate={false}
      />
    </div>
  );
};

export default CustomDateTimeEditor;
