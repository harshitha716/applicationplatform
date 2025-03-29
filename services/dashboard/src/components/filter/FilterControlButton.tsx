import React, { FC, PropsWithChildren } from 'react';
import { SIZE } from 'constants/common.constants';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { cn } from 'utils/common';
import { Loader } from 'components/common/loader/Loader';
import { Tooltip, TooltipPositions } from 'components/common/tooltip';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

interface FilterControlButtonProps extends PropsWithChildren {
  onClick: (event: React.MouseEvent<HTMLButtonElement>) => void;
  tooltipText?: string;
  icon?: string;
  iconCategory?: ICON_SPRITE_TYPES;
  iconColor?: string;
  buttonRef?: React.RefObject<HTMLButtonElement>;
  isSelected?: boolean;
  childrenWrapperClassName?: string;
  className?: string;
  isLoading?: boolean;
  disabled?: boolean;
  id?: string;
  tooltipPosition?: TooltipPositions;
}

const FilterControlButton: FC<FilterControlButtonProps> = ({
  onClick,
  tooltipText = '',
  tooltipPosition = TooltipPositions.BOTTOM,
  icon = 'filter-lines',
  iconCategory = ICON_SPRITE_TYPES.GENERAL,
  buttonRef,
  children,
  isSelected = false,
  childrenWrapperClassName = '',
  className = '',
  isLoading = false,
  id = '',
  disabled = false,
}) => {
  const onButtonClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    if (disabled || isLoading) {
      return;
    }

    onClick(event);
  };

  return (
    <Tooltip
      tooltipBody={tooltipText}
      color='{TMS_COLORS.GRAY_200}'
      tooltipBodyClassName='f-12-450 px-3 py-1.5 rounded-md whitespace-nowrap z-999 bg-black text-GRAY_200'
      className='z-1'
      disabled={disabled}
      position={tooltipPosition}
    >
      <button
        className={cn(
          'border border-GRAY_400 rounded px-2 py-1.5 w-fit outline-none flex items-center h-[26px] text-GRAY_1000',
          className,
          isSelected ? 'bg-DIVIDER_SAIL_1' : '',
          disabled ? 'opacity-50' : 'hover:border-DIVIDER_SAIL_4',
        )}
        onClick={onButtonClick}
        ref={buttonRef}
        data-testid={`filter-control-button-${id}`}
        disabled={disabled}
      >
        {isLoading ? (
          <Loader size={SIZE.XSMALL} className='m-auto' />
        ) : (
          <>
            <SvgSpriteLoader id={icon} iconCategory={iconCategory} width={12} height={12} />
            {!!children && (
              <span className={`f-12-500 ${typeof children === 'string' ? 'ml-1' : ''} ${childrenWrapperClassName}`}>
                {children}
              </span>
            )}
          </>
        )}
      </button>
    </Tooltip>
  );
};

export default FilterControlButton;
