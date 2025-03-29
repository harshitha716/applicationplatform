import React, { ReactElement } from 'react';
import { defaultFn } from 'types/commonTypes';
import { cn } from 'utils/common';
import { CHECKBOX_STATE_TYPES, CHECKBOX_TYPES } from 'components/common/Checkbox/checkbox.constant';

const checkmarkStyles = `after:content-[''] after:hidden after:left-1 after:top-px after:w-1 after:h-2 after:border-r-2 after:border-b-2 after:absolute after:border-solid after:border-white after:rounded-[1px] after:rotate-45 peer-checked:after:block`;
const CHECKBOX_STATE_STYLES = {
  [CHECKBOX_TYPES.SELECTED]: {
    [CHECKBOX_STATE_TYPES.ENABLED]: 'bg-GRAY_1000 cursor-pointer border border-GRAY_600',
    [CHECKBOX_STATE_TYPES.HOVER]: '',
    [CHECKBOX_STATE_TYPES.DISABLED]: 'bg-TEXT_TERTIARY cursor-not-allowed',
  },
  [CHECKBOX_TYPES.UNSELECTED]: {
    [CHECKBOX_STATE_TYPES.ENABLED]: 'bg-white border border-GRAY_400',
    [CHECKBOX_STATE_TYPES.HOVER]: 'hover:bg-GRAY_200 hover:border-GRAY_200',
    [CHECKBOX_STATE_TYPES.DISABLED]: 'bg-DIVIDER_GRAY border-2 border-TEXT_TERTIARY cursor-not-allowed',
  },
};

export interface CheckBoxProps {
  checked: boolean;
  onPress?: React.MouseEventHandler<HTMLInputElement>;
  disabled?: boolean;
  id: string;
  className?: string;
  displayContainerClassName?: string;
  isCustomCheckMark?: boolean;
  customCheckMarkClassName?: string;
  customCheckMark?: ReactElement;
}

export const CheckBox: React.FC<CheckBoxProps> = ({
  checked,
  id,
  onPress,
  disabled = false,
  className = '',
  displayContainerClassName = '',
  isCustomCheckMark = false,
  customCheckMarkClassName = '',
  customCheckMark,
}) => {
  const handlePress: React.MouseEventHandler<HTMLInputElement> = (e) => {
    onPress?.(e);
  };

  const stateStyles = checked
    ? CHECKBOX_STATE_STYLES[CHECKBOX_TYPES.SELECTED]
    : CHECKBOX_STATE_STYLES[CHECKBOX_TYPES.UNSELECTED];

  const checkBoxStylesByState = disabled
    ? stateStyles[CHECKBOX_STATE_TYPES.DISABLED]
    : `${stateStyles[CHECKBOX_STATE_TYPES.ENABLED]} ${stateStyles[CHECKBOX_STATE_TYPES.HOVER]}`;

  return (
    <div className={`h-3.5 w-3.5 flex items-center relative ${className}`} data-testid={`checkbox-wrapper-${id}`}>
      <input
        type='checkbox'
        value=''
        data-checkboxid='check-box'
        checked={checked}
        onClick={handlePress}
        onChange={defaultFn}
        disabled={disabled}
        id={id}
        className='absolute opacity-0 cursor-pointer h-0 w-0 peer'
        role='checkbox'
      />
      <span
        onClick={disabled ? undefined : handlePress}
        data-checkboxid='check-box'
        className={cn(
          'absolute top-0 left-0 h-3.5 w-3.5 rounded',
          isCustomCheckMark ? customCheckMarkClassName : checkmarkStyles,
          displayContainerClassName,
          checkBoxStylesByState,
        )}
        data-testid={`checkbox-span-${id}`}
      >
        {customCheckMark ?? null}
      </span>
    </div>
  );
};
