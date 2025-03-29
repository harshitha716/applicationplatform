import React from 'react';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

export type IconProps = {
  id: string;
  category?: ICON_SPRITE_TYPES;
  customIcon?: React.ReactNode;
  size?: number;
  className?: string;
};

const Icon = ({ id, category, customIcon, size, className }: IconProps) => {
  return (
    <div className=''>
      {customIcon ? (
        customIcon
      ) : (
        <SvgSpriteLoader id={id ?? ''} height={size} width={size} iconCategory={category} className={className} />
      )}
    </div>
  );
};

export default Icon;
