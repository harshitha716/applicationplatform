import React, { FC } from 'react';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { DISPLAY_OPTIONS } from 'components/common/table/table.types';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

export type DisplayOptionItemProps = {
  id: DISPLAY_OPTIONS;
  label: string;
  iconId: string;
  onClick?: (id: DISPLAY_OPTIONS) => void;
  value?: string;
};

const DisplayOptionItem: FC<DisplayOptionItemProps> = ({ id, label, iconId, onClick, value }) => {
  return (
    <div
      key={id}
      className='flex items-center justify-between py-2 px-2.5 hover:bg-GRAY_100 group cursor-pointer rounded-md'
      onClick={() => onClick?.(id)}
    >
      <div className='flex items-center gap-1.5'>
        <SvgSpriteLoader id={iconId} width={12} height={12} />
        <div className='f-12-500'>{label}</div>
      </div>
      <div className='flex items-center gap-1.5'>
        {value && <span className='text-GRAY_700 f-12-400'>{value}</span>}
        <SvgSpriteLoader
          id='arrow-narrow-right'
          iconCategory={ICON_SPRITE_TYPES.ARROWS}
          width={12}
          height={12}
          className='group-hover:opacity-100 opacity-0'
        />
      </div>
    </div>
  );
};

export default DisplayOptionItem;
