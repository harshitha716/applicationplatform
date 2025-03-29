import { SIZE_TYPES } from 'types/common/components';
import { MapAny } from 'types/commonTypes';
import { SvgSpriteLoaderProps } from 'components/SvgSpriteLoader';

export enum BUTTON_TYPES {
  PRIMARY = 'PRIMARY',
  SECONDARY = 'SECONDARY',
  TEXT_NAV = 'TEXT_NAV',
  DANGER = 'DANGER',
  SHARE = 'SHARE',
}

export enum ICON_POSITION_TYPES {
  LEFT = 'LEFT',
  RIGHT = 'RIGHT',
}

export enum BUTTON_STATE_TYPES {
  COMMON = 'COMMON',
  DEFAULT = 'DEFAULT',
  HOVER = 'HOVER',
  PRESSED = 'PRESSED',
  DISABLED = 'DISABLED',
  LOADING = 'LOADING',
}

export type EventCallbackType = (id: string, payload: MapAny) => void;

export interface ButtonProps {
  className?: string;
  disabled?: boolean;
  size?: SIZE_TYPES;
  state?: BUTTON_STATE_TYPES;
  onClick?: (e?: React.MouseEvent<HTMLButtonElement>) => void;
  isLoading?: boolean;
  defaultLoader?: boolean;
  id: string;
  type?: BUTTON_TYPES;
  customLeadingIcon?: React.ReactNode;
  customTrailingIcon?: React.ReactNode;
  children?: React.ReactNode | string | null;
  childrenClassName?: string;
  loader?: React.ReactNode | string | null;
  tabIndex?: number;
  iconPosition?: ICON_POSITION_TYPES;
  iconProps?: SvgSpriteLoaderProps;
  textSizeOverrideClassName?: string;
  customAttributes?: MapAny;
  onMouseEnter?: (e: React.MouseEvent<HTMLButtonElement>) => void;
  onMouseLeave?: (e: React.MouseEvent<HTMLButtonElement>) => void;
}
