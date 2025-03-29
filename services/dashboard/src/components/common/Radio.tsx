import React, { FC, useEffect, useState } from 'react';
import { RADIO_STATE_STYLES } from 'constants/radio.constants';
import { RADIO_STATE_TYPES, RADIO_TYPES, RadioProps } from 'types/common/components/radio';
import { defaultFn } from 'types/commonTypes';

export const Radio: FC<RadioProps> = ({
  wrapperClassName = 'flex items-center justify-center w-12 h-12 rounded-full z-10',
  radioClassName = 'w-5 h-5 rounded-full z-20',
  radioSelectedClassName = 'border-[7px]',
  radioDefaultClassName = 'border',
  wrapperStyle = '',
  radioStyle = '',
  radioSelectedStyle = '',
  radioDefaultStyle = '',
  checked = false,
  isDisabled = false,
  onSelect = defaultFn,
  isClearable = true,
  id,
}) => {
  const [isSelected, setIsSelected] = useState<boolean>(checked);

  const stateStyles = isSelected
    ? RADIO_STATE_STYLES[RADIO_TYPES.SELECTED]
    : RADIO_STATE_STYLES[RADIO_TYPES.UNSELECTED];

  const radioStylesByState = isDisabled
    ? stateStyles[RADIO_STATE_TYPES.DISABLED]
    : `${stateStyles[RADIO_STATE_TYPES.ENABLED]} ${stateStyles[RADIO_STATE_TYPES.HOVER]} ${
        stateStyles[RADIO_STATE_TYPES.PRESSED]
      }`;

  useEffect(() => {
    setIsSelected((isSelected) => (checked !== isSelected ? checked : isSelected));
  }, [checked]);

  const handleSelect = () => {
    if (!isDisabled) {
      const currentStatus = !isSelected;

      setIsSelected(currentStatus);
      onSelect(currentStatus);
    }
  };

  return (
    <div className={`${wrapperClassName} ${wrapperStyle}`} data-testid={`radio-wrapper-${id}`}>
      <div
        data-testid={`radio-${id}`}
        role='presentation'
        className={`${radioStylesByState} ${radioClassName} ${radioStyle} ${
          isSelected
            ? `${radioSelectedClassName} ${radioSelectedStyle}`
            : `${radioDefaultClassName} ${radioDefaultStyle}`
        }`}
        onClick={isClearable || !isSelected ? handleSelect : defaultFn}
      ></div>
    </div>
  );
};
