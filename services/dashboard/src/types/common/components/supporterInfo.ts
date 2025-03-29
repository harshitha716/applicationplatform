import { MouseEventHandler } from 'react';
import { ICON_SPRITE_TYPES } from 'constants/icons';

export enum SUPPORT_INFO_TYPES {
  GUIDE = 'GUIDE',
  ERROR = 'ERROR',
  CUSTOM = 'CUSTOM',
}

export interface SvgSpriteLoaderProps {
  width?: number;
  height?: number;
  fillColor?: string;
  color?: string;
  iconCategory?: ICON_SPRITE_TYPES;
  id: string;
  viewBox?: string;
  domain?: string;
  dataCache?: string;
  version?: number;
  className?: string;
  onClick?: MouseEventHandler<HTMLDivElement>;
  customSpriteUrl?: string;
}
