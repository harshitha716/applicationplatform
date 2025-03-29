import React, { FC, useEffect, useState } from 'react';
import { cn } from 'utils/common';
import { TOGGLE_SWITCH_SLIDER, TOGGLE_SWITCH_STYLES } from 'components/common/toggleSwitch/toggleSwitch.constants';
import {
  TOGGLE_SWITCH_STATE_TYPES,
  TOGGLE_SWITCH_TYPES,
  ToggleSwitchProps,
} from 'components/common/toggleSwitch/toggleSwitch.types';

const ToggleSwitch: FC<ToggleSwitchProps> = ({
  checked = false,
  onChange,
  label = '',
  disabled = false,
  toggleClassName = 'relative outline-none h-3 w-5 rounded-full transition-all duration-200',
  sliderClassName = 'absolute top-[2px] rounded-full w-2 h-2 transition-all duration-200',
  sliderStyle = '',
  toggleStyle = '',
  wrapperClassName = '',
  labelClassName = '',
  eventCallback,
  id,
  controlled = false,
}) => {
  const [active, setActive] = useState<boolean>(checked);

  const stateStyles = active
    ? TOGGLE_SWITCH_STYLES[TOGGLE_SWITCH_TYPES.SELECTED]
    : TOGGLE_SWITCH_STYLES[TOGGLE_SWITCH_TYPES.UNSELECTED];

  const sliderStyles = active
    ? TOGGLE_SWITCH_SLIDER[TOGGLE_SWITCH_TYPES.SELECTED]
    : TOGGLE_SWITCH_SLIDER[TOGGLE_SWITCH_TYPES.UNSELECTED];

  const toggleSwitchStylesByState = disabled
    ? stateStyles[TOGGLE_SWITCH_STATE_TYPES.DISABLED]
    : `${stateStyles[TOGGLE_SWITCH_STATE_TYPES.ENABLED]} ${stateStyles[TOGGLE_SWITCH_STATE_TYPES.HOVER]}`;

  const toggleSwitchSliderStylesByState = disabled
    ? sliderStyles[TOGGLE_SWITCH_STATE_TYPES.DISABLED]
    : `${sliderStyles[TOGGLE_SWITCH_STATE_TYPES.ENABLED]} `;

  const toggle = () => {
    if (disabled) return;
    onChange(!active);
    eventCallback?.('TOGGLE_SWITCH_CHANGE', { id, checked: !checked });
    if (!controlled) setActive(!active);
  };

  useEffect(() => {
    setActive(checked);
  }, [checked]);

  return (
    <div className={`${wrapperClassName}`} data-testid={`toggle-switch-wrapper-${id}`}>
      <button
        role='switch'
        aria-checked={active}
        className={cn(toggleSwitchStylesByState, toggleClassName, toggleStyle)}
        onClick={toggle}
        data-testid={`toggle-switch-btn-${id}`}
      >
        <span
          className={cn(sliderClassName, sliderStyle, toggleSwitchSliderStylesByState)}
          data-testid={`toggle-switch-slider-${id}`}
        />
      </button>
      {label && (
        <label className={labelClassName} data-testid={`toggle-switch-label-${id}`}>
          {label}
        </label>
      )}
    </div>
  );
};

export default ToggleSwitch;
