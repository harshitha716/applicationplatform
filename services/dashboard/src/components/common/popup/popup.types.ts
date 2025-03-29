import { ICON_SPRITE_TYPES } from 'constants/icons';
import { defaultFnType } from 'types/commonTypes';

export type PopupProps = {
  title?: string;
  subTitle?: string;
  titleClassName?: string;
  subTitleClassName?: string;
  popupWrapperClassName?: string;
  showIcon?: boolean;
  iconCategory?: ICON_SPRITE_TYPES;
  iconId: string;
  iconColor?: string;
  isOpen: boolean;
  children: any;
  className?: string;
  onClose?: defaultFnType;
  closeOnClickOutside?: boolean;
  isOverlay?: boolean;
  wrapperClassName?: string;
};
