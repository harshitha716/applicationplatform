import { FC } from 'react';
import { COLORS } from 'constants/colors';
import Image from 'next/image';
import { SIZE_TYPES } from 'types/common/components';
import { defaultFnType } from 'types/commonTypes';
import { BUTTON_TYPES } from 'types/components/button.type';
import { cn } from 'utils/common';
import { Button } from 'components/common/button/Button';
import { Tooltip, TooltipPositions } from 'components/common/tooltip';
import { SvgSpriteLoaderProps } from 'components/SvgSpriteLoader';

interface TooltipButtonPropsType {
  tooltipBodyClassName?: string;
  tooltipBodyOverrideClassName?: string;
  tooltipBody?: string;
  className?: string;
  tooltipColor?: string;
  tooltipPosition?: TooltipPositions;
  buttonIcon?: SvgSpriteLoaderProps;
  buttonTitle?: string;
  buttonType?: BUTTON_TYPES;
  buttonSize?: SIZE_TYPES;
  id: string;
  onClick?: defaultFnType;
  isLoading?: boolean;
  disabled?: boolean;
  buttonDisabled?: boolean;
  imageIconSrc?: string;
}

const TooltipButton: FC<TooltipButtonPropsType> = ({
  tooltipBodyClassName = '',
  tooltipBodyOverrideClassName = 'f-12-300 tw-py-1 tw-px-2 tw-rounded-sm tw-whitespace-nowrap',
  tooltipBody = '',
  className = '',
  tooltipColor = COLORS.BLACK,
  tooltipPosition = TooltipPositions.BOTTOM,
  buttonIcon,
  buttonTitle = '',
  buttonType = BUTTON_TYPES.SECONDARY,
  buttonSize = SIZE_TYPES.SMALL,
  id = '',
  onClick,
  isLoading = false,
  disabled = false,
  buttonDisabled = false,
  imageIconSrc,
}) => {
  return (
    <Tooltip
      tooltipBody={tooltipBody}
      position={tooltipPosition}
      color={tooltipColor}
      disabled={disabled}
      tooltipBodyClassName={cn(tooltipBodyOverrideClassName, tooltipBodyClassName)}
    >
      <Button
        type={buttonType}
        id={id}
        size={buttonSize}
        className={className}
        iconProps={buttonIcon}
        onClick={onClick}
        isLoading={isLoading}
        disabled={buttonDisabled}
      >
        {imageIconSrc && <Image alt='' src={imageIconSrc} width={14} height={14} />}
        {buttonTitle}
      </Button>
    </Tooltip>
  );
};

export default TooltipButton;
