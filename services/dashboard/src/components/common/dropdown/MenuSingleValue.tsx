import React, { FC } from 'react';
import { MenuSingleValuePropsType } from 'types/common/components/dropdown/dropdown.types';
import { DROPDOWN_SIZE_STYLES } from 'components/common/dropdown/dropdown.constants';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

const MenuSingleValue: FC<MenuSingleValuePropsType> = ({
  icon,
  label,
  spriteIcon,
  value,
  showValueInControl,
  size,
  customClassNames,
}) => {
  return (
    <div className='flex items-center'>
      {spriteIcon && (
        <div className='w-6 mr-4'>
          <SvgSpriteLoader id={spriteIcon} />
        </div>
      )}
      {!showValueInControl && icon}
      <div className={customClassNames?.placeholder ?? DROPDOWN_SIZE_STYLES[size].customClassNames.placeholder}>
        {showValueInControl ? value : label}
      </div>
    </div>
  );
};

export default MenuSingleValue;
