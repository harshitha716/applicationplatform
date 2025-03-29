import React, { useEffect } from 'react';
import { COLORS } from 'constants/colors';
import { DateRangeKeys } from 'constants/date.constants';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { KEYBOARD_KEYS } from 'constants/shortcuts';
import { SIZE_TYPES } from 'types/common/components';
import { defaultFnType } from 'types/commonTypes';
import { getPlacehoderDate } from 'components/common/dateRangePicker/dateRangePicker.utils';
import Input from 'components/common/input';

interface DateSearchProps {
  value: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onApply: () => void;
  onClear: defaultFnType;
  id: DateRangeKeys;
  focusedInput: DateRangeKeys;
  currentTab: string;
  shouldShowInfo?: boolean;
}

export const DateSearch: React.FC<DateSearchProps> = ({
  value,
  onChange,
  onApply,
  id,
  focusedInput,
  shouldShowInfo = false,
  onClear,
}) => {
  const inputRef = React.useRef<HTMLInputElement>(null);

  const onKeyUp = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === KEYBOARD_KEYS.ENTER) {
      e.preventDefault();
      e.stopPropagation();

      onApply();
    }
  };

  useEffect(() => {
    if (focusedInput === id) {
      inputRef.current?.select();
    }
  }, [focusedInput]);

  const placeholderDate = getPlacehoderDate();

  const placeholder =
    id === DateRangeKeys.START_DATE
      ? !shouldShowInfo
        ? `Try: Next month, Q4 or ${placeholderDate}`
        : 'Add start date'
      : 'Add end date';

  return (
    <Input
      inputRef={inputRef}
      className='w-full'
      size={SIZE_TYPES.XSMALL}
      id={`date-search-${id}`}
      inputFieldWrapperClassName=' w-full relative !placeholder-GRAY_500 f-12-500'
      placeholder={placeholder}
      value={value}
      onChange={onChange}
      onKeyUp={onKeyUp}
      noBorders
      trailingIconProps={{
        iconCategory: ICON_SPRITE_TYPES.GENERAL,
        id: 'x-close',
        className: `absolute bottom-5 ${value.length ? '' : 'hidden'}`,
        width: 14,
        height: 14,
        color: COLORS.GRAY_400,
        onClick: onClear,
      }}
    />
  );
};
