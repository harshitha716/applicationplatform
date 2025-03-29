import { FC, MouseEvent, useEffect, useState } from 'react';
import { KEYBOARD_KEYS } from 'constants/shortcuts';
import useKeyDown from 'hooks/useKeyDown';
import { POSITION_TYPES, SIZE_TYPES } from 'types/common/components';
import { cn } from 'utils/common';
import OverlayFooter from 'components/common/SideDrawer/OverlayFooter';
import OverlayTitle from 'components/common/SideDrawer/OverlayTitle';
import { SIDE_DRAWER_TYPES, SideDrawerProps } from 'components/common/SideDrawer/sideDrawer.types';

const stackedClassNames: Record<number, string> = {
  0: '',
  1: 'w-sideDrawer1 h-sideDrawer1',
  2: 'w-sideDrawer2 h-sideDrawer2',
  3: 'w-sideDrawer3 h-sideDrawer3',
  4: 'w-sideDrawer4 h-sideDrawer4',
  5: 'w-sideDrawer5 h-sideDrawer5',
};

const SIZE_CLASSNAMES: Record<SIZE_TYPES, string> = {
  [SIZE_TYPES.XSMALL]: 'w-sideDrawerXSmall',
  [SIZE_TYPES.SMALL]: 'w-sideDrawer',
  [SIZE_TYPES.MEDIUM]: 'w-sideDrawerMedium',
  [SIZE_TYPES.LARGE]: 'w-sideDrawerLarge',
  [SIZE_TYPES.XLARGE]: 'w-sideDrawerLarge',
};

const SIDE_DRAWER_TYPES_CLASSNAMES: Record<
  SIDE_DRAWER_TYPES,
  { wrapper: string; children: string; backdropClassName: string }
> = {
  [SIDE_DRAWER_TYPES.PRIMARY]: {
    wrapper: 'shadow-sideDrawer bg-white border',
    children: 'bg-white',
    backdropClassName: 'pt-12',
  },
  [SIDE_DRAWER_TYPES.SECONDARY]: {
    wrapper: 'p-2.5 backdrop-blur-sm',
    children: 'rounded-xl bg-white w-full h-full sideDrawerInner border-[0.5px] border-GRAY_500 shadow-sideDrawerInner',
    backdropClassName: '',
  },
};

const POSITION_CLASSNAMES = {
  [POSITION_TYPES.RIGHT]: {
    mountClassName: 'animate-rightSideDrawer-mount',
    common: '',
    unmountClassName: 'animate-rightSideDrawer-unMount',
  },
  [POSITION_TYPES.BOTTOM]: {
    mountClassName: 'animate-bottomSideDrawer-mount',
    unmountClassName: 'animate-bottomSideDrawer-unMount',
    common: '!w-screen max-h-[90vh]',
  },
  [POSITION_TYPES.TOP]: { mountClassName: '', unmountClassName: '', common: '' },
  [POSITION_TYPES.LEFT]: { mountClassName: '', unmountClassName: '', common: '' },
};

const SideDrawer: FC<SideDrawerProps> = ({
  isOpen = false,
  type = SIDE_DRAWER_TYPES.PRIMARY,
  size = SIZE_TYPES.SMALL,
  position = POSITION_TYPES.RIGHT,
  onClose,
  children,
  className = '',
  closeButtonClassName = '',
  closeOnClickOutside = true,
  stackPosition = 0,
  title = '',
  subtitle = '',
  backdropClassName = '',
  hideCloseButton = false,
  headerClassName = '',
  id = 'GENERAL_SIDEDRAWER',
  onBack,
  onNext,
  titleClassName = '',
  subtitleClassName = '',
  nextButtonClassName = 'min-w-[106px]',
  backButtonClassName = 'min-w-[106px]',
  nextButtonTitle = 'Next',
  backButtonTitle = 'Back',
  bottomBar,
  topBar,
  step = '',
  isNextButtonLoading = false,
  isNextButtonDisabled = false,
  childrenWrapperClassName = 'overflow-scroll ',
  nextButtonIconProps,
  nextButtonIconPosition,
  footerClassName = '',
  closeButtonDimensions = { width: 24, height: 24 },
  animateOnClose = true,
  nextButtonSize = SIZE_TYPES.SMALL,
  backButtonSize = SIZE_TYPES.SMALL,
}) => {
  const [isMount, setIsMount] = useState(true);
  const handleClose = () => {
    if (animateOnClose) setIsMount(false);
    setTimeout(() => onClose(), 200);
  };

  useKeyDown(handleClose, KEYBOARD_KEYS.ESCAPE);

  const handleClickOutside = (e: MouseEvent<HTMLDivElement>) => {
    e?.stopPropagation();
    if (closeOnClickOutside) handleClose();
  };

  useEffect(() => {
    if (isOpen) setIsMount(true);
  }, [isOpen]);

  if (!isOpen) return null;

  const sidebarWidthClasses = SIZE_CLASSNAMES[size];

  const { mountClassName, unmountClassName, common } = POSITION_CLASSNAMES[position as POSITION_TYPES];

  return (
    <div
      className={cn(
        'h-full fixed w-screen z-1000 top-0 left-0 items-center animate-opacity',
        isOpen ? '' : 'hidden',
        backdropClassName,
        SIDE_DRAWER_TYPES_CLASSNAMES[type].backdropClassName,
      )}
      role='presentation'
      onClick={(e: MouseEvent<HTMLDivElement>) => handleClickOutside(e)}
      data-testid={`side-drawer-${id}`}
    >
      <div
        className={cn(
          isMount ? mountClassName : unmountClassName,
          common,
          ' -right-[100vw] w-screen flex flex-col absolute transition-all h-screen',
          stackedClassNames[stackPosition],
          sidebarWidthClasses,
          className,
          SIDE_DRAWER_TYPES_CLASSNAMES[type].wrapper,
        )}
        role='presentation'
        onClick={(e: MouseEvent) => e.stopPropagation()}
      >
        <div className={cn(' bg-white w-full h-full', SIDE_DRAWER_TYPES_CLASSNAMES[type].children)}>
          <OverlayTitle
            closeButtonDimensions={closeButtonDimensions}
            topBar={topBar}
            title={title}
            hideCloseButton={hideCloseButton}
            step={step}
            subtitle={subtitle}
            onClose={handleClose}
            headerClassName={headerClassName}
            closeButtonClassName={closeButtonClassName}
            titleClassName={titleClassName}
            subtitleClassName={subtitleClassName}
          />
          <div className={cn('p-4 h-full', childrenWrapperClassName)}>{children}</div>
          <OverlayFooter
            onBack={onBack}
            onNext={onNext}
            nextButtonClassName={nextButtonClassName}
            backButtonClassName={backButtonClassName}
            nextButtonTitle={nextButtonTitle}
            backButtonTitle={backButtonTitle}
            bottomBar={bottomBar}
            isNextButtonLoading={isNextButtonLoading}
            isNextButtonDisabled={isNextButtonDisabled}
            nextButtonIconProps={nextButtonIconProps}
            nextButtonIconPosition={nextButtonIconPosition}
            footerClassName={footerClassName}
            nextButtonSize={nextButtonSize}
            backButtonSize={backButtonSize}
          />
        </div>
      </div>
    </div>
  );
};

export default SideDrawer;
