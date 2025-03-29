import React from 'react';
import { MonthsConfig } from 'constants/date.constants';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { getYearList } from 'components/common/dateRangePicker/dateRangePicker.utils';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

export const DateRangePickerNavigator = (
  currFocusedDate: Date,
  changeShownDate: (
    value: string | number | Date,
    mode?: 'set' | 'setYear' | 'setMonth' | 'monthOffset' | undefined,
  ) => void,
) => {
  const currentMonth = MonthsConfig[currFocusedDate.getMonth()].value;

  const yearsList = getYearList();

  return (
    <div className=' flex justify-between'>
      <div className=' flex w-fit border border-DIVIDER_SAIL_2'>
        <select
          className='f-12-400 appearance-none px-2 py-1 bg-BG_GRAY_2 focus:outline-none border border-r-DIVIDER_SAIL_2 border-y-0 border-l-0 outline-none cursor-pointer'
          value={currentMonth}
          onChange={(e) => changeShownDate(e.target.value, 'setMonth')}
        >
          {MonthsConfig.map((month) => (
            <option key={month.short} value={month.value}>
              {month.short}
            </option>
          ))}
        </select>

        <select
          className='f-12-400 appearance-none bg-BG_GRAY_2 focus:outline-none px-2 py-1 border-none  outline-none cursor-pointer'
          value={currFocusedDate?.getFullYear()}
          onChange={(e) => changeShownDate(e.target.value, 'setYear')}
        >
          {yearsList.map((year, index) => (
            <option key={index} value={year}>
              {year}
            </option>
          ))}
        </select>
      </div>

      <div className='flex '>
        <button className=' text-DIVIDER_SAIL_4  mr-3' onClick={() => changeShownDate(-1, 'monthOffset')}>
          <SvgSpriteLoader id='chevron-left' iconCategory={ICON_SPRITE_TYPES.ARROWS} width={16} height={16} />
        </button>
        <button className=' text-DIVIDER_SAIL_4' onClick={() => changeShownDate(1, 'monthOffset')}>
          <SvgSpriteLoader id='chevron-right' iconCategory={ICON_SPRITE_TYPES.ARROWS} width={16} height={16} />
        </button>
      </div>
    </div>
  );
};
