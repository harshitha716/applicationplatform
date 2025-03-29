import { ReactElement } from 'react';
import { MapAny } from 'types/commonTypes';
import { SvgSpriteLoaderProps } from 'components/SvgSpriteLoader';

export enum SIZE_TYPES {
  XLARGE = 'XLARGE',
  LARGE = 'LARGE',
  MEDIUM = 'MEDIUM',
  SMALL = 'SMALL',
  XSMALL = 'XSMALL',
}

export enum POSITION_TYPES {
  LEFT = 'LEFT',
  RIGHT = 'RIGHT',
  BOTTOM = 'BOTTOM',
  TOP = 'TOP',
}

export interface TextProps {
  textClass?: string;
  children: string | ReactElement;
  id?: string;
}

export type EventCallbackType = (id: string, payload: MapAny) => void;

export interface MenuItem {
  label: string;
  value: string | number;
  labelColor?: string;
  color?: string;
  icons?: string;
  isDisabled?: boolean;
  id?: string;
  code?: string;
  isHidden?: boolean;
  iconProps?: SvgSpriteLoaderProps;
  metadata?: MapAny;
  description?: string;
}

export enum TAB_TYPES {
  FILLED = 'FILLED',
  OUTLINE = 'OUTLINE',
  FILLED_OUTLINED = 'FILLED_OUTLINED',
  UNDERLINE = 'UNDERLINE',
}
