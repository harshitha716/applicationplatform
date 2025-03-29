import React, { FC } from 'react';
import { COLORS } from 'constants/colors';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { ResetSectionProps } from 'types/common/components/dropdown/dropdown.types';
import { cn } from 'utils/common';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

export const ResetSection: FC<ResetSectionProps> = ({ resetProps, onClickReset }) => {
  return (
    <div
      className={cn('flex py-3 pl-4 border-t border-DIVIDER_GRAY', resetProps?.resetClassName)}
      onClick={onClickReset}
    >
      <SvgSpriteLoader
        id='refresh-ccw-01'
        iconCategory={ICON_SPRITE_TYPES.ARROWS}
        height={14}
        width={14}
        color={COLORS.TEXT_PRIMARY}
      />
      <div className={cn('pl-2 f-12-400 text-GRAY_700', resetProps?.resetTextClassName)}>{resetProps?.resetText}</div>
    </div>
  );
};
