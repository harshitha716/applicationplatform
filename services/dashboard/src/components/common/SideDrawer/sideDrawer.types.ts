import { ReactElement } from 'react';
import { POSITION_TYPES, SIZE_TYPES } from 'types/common/components';
import { defaultFnType } from 'types/commonTypes';
import { ICON_POSITION_TYPES } from 'types/components/button.type';
import { SvgSpriteLoaderProps } from 'components/SvgSpriteLoader';

export enum SIDE_DRAWER_TYPES {
  PRIMARY = 'PRIMARY',
  SECONDARY = 'SECONDARY',
}

export interface OverlayTitleProps {
  topBar?: ReactElement | string;
  subtitle?: string | ReactElement;
  step?: string;
  title?: string | ReactElement;
  hideCloseButton?: boolean;
  onClose: defaultFnType;
  headerClassName?: string;
  closeButtonClassName?: string;
  closeButtonDimensions?: { width: number; height: number };
  titleClassName?: string;
  subtitleClassName?: string;
}

export interface OverlayFooterProps {
  bottomBar?: ReactElement;
  onBack?: defaultFnType;
  onNext?: defaultFnType;
  isNextButtonDisabled?: boolean;
  nextButtonClassName?: string;
  backButtonClassName?: string;
  nextButtonTitle?: string | ReactElement;
  backButtonTitle?: string | ReactElement;
  isBackButtonLoading?: boolean;
  isNextButtonLoading?: boolean;
  nextButtonIconProps?: SvgSpriteLoaderProps;
  nextButtonIconPosition?: ICON_POSITION_TYPES;
  footerClassName?: string;
  nextButtonSize?: SIZE_TYPES;
  backButtonSize?: SIZE_TYPES;
}

export interface SideDrawerProps extends OverlayTitleProps, OverlayFooterProps {
  isOpen: boolean;
  children: React.ReactNode;
  className?: string;
  closeOnClickOutside?: boolean;
  stackPosition?: number;
  backdropClassName?: string;
  id: string;
  size?: SIZE_TYPES;
  childrenWrapperClassName?: string;
  animateOnClose?: boolean;
  titleClassName?: string;
  subtitleClassName?: string;
  position?: POSITION_TYPES;
  type?: SIDE_DRAWER_TYPES;
}
