import { FC } from 'react';
import { SIZE_TYPES } from 'types/common/components';
import { BUTTON_TYPES } from 'types/components/button.type';
import { cn } from 'utils/common';
import { Button } from 'components/common/button/Button';
import { OverlayFooterProps } from 'components/common/SideDrawer/sideDrawer.types';

// TODO: needs to be updated to use the new design system
const OverlayFooter: FC<OverlayFooterProps> = ({
  onBack,
  onNext,
  nextButtonClassName = 'tw-min-w-[106px]',
  backButtonClassName = 'tw-min-w-[106px]',
  nextButtonTitle = 'Next',
  backButtonTitle = 'Back',
  bottomBar,
  isBackButtonLoading,
  isNextButtonLoading,
  isNextButtonDisabled,
  nextButtonIconProps,
  nextButtonIconPosition,
  footerClassName = '',
  nextButtonSize = SIZE_TYPES.SMALL,
  backButtonSize = SIZE_TYPES.SMALL,
}) =>
  (bottomBar || onBack || onNext) && (
    <div
      className={cn(
        'tw-border-t tw-px-4 tw-py-6 tw-border-ORANGE_100 tw-flex tw-justify-end tw-items-center tw-gap-3 tw-min-h-[86px] tw-bg-white',
        footerClassName,
      )}
    >
      {bottomBar ? (
        bottomBar
      ) : (
        <>
          {!!onBack && (
            <Button
              size={backButtonSize}
              type={BUTTON_TYPES.SECONDARY}
              className={backButtonClassName}
              id='SIDE_DRAWER_BACK_BUTTON'
              onClick={onBack}
              isLoading={isBackButtonLoading}
            >
              {backButtonTitle}
            </Button>
          )}
          {!!onNext && (
            <Button
              size={nextButtonSize}
              className={nextButtonClassName}
              id='SIDE_DRAWER_NEXT_BUTTON'
              onClick={onNext}
              isLoading={isNextButtonLoading}
              disabled={isNextButtonDisabled}
              iconProps={nextButtonIconProps}
              iconPosition={nextButtonIconPosition}
            >
              {nextButtonTitle}
            </Button>
          )}
        </>
      )}
    </div>
  );

export default OverlayFooter;
