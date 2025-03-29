import React, { FC } from 'react';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { defaultFnType } from 'types/commonTypes';
import { cn } from 'utils/common';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

export interface MenuWrapperProps {
  children: React.ReactNode;
  className?: string;
  resetText?: string;
  resetClassName?: string;
  resetTextClassName?: string;
  onReset?: defaultFnType;
  id: string;
  childrenWrapperClassName?: string;
}
export const MenuWrapper: FC<MenuWrapperProps> = ({
  children,
  className = '',
  resetText = 'Reset filters',
  resetClassName = '',
  resetTextClassName = '',
  onReset,
  id,
  childrenWrapperClassName = '',
}) => {
  const handleReset = () => {
    onReset?.();
  };

  return (
    <div
      className={cn('bg-white relative z-1 shadow-menuList rounded-md border-0.5 border-GRAY_500', className)}
      data-testid={`menu-wrapper-${id}`}
    >
      <div
        className={`max-h-[300px] overflow-y-scroll ${childrenWrapperClassName}`}
        data-testid={`menu-wrapper-children-${id}`}
      >
        {children}
      </div>
      {!!onReset && (
        <div
          className={`flex py-3 pl-4 border-t border-GRAY_400 ${resetClassName}`}
          onClick={handleReset}
          data-testid={`menu-wrapper-reset-${id}`}
        >
          <SvgSpriteLoader id='refresh-ccw-01' iconCategory={ICON_SPRITE_TYPES.ARROWS} height={14} width={14} />
          <div className={`pl-2 f-12-400 ${resetTextClassName}`} data-testid={`menu-wrapper-reset-text-${id}`}>
            {resetText}
          </div>
        </div>
      )}
    </div>
  );
};
