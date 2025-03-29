import React, { memo, MouseEventHandler } from 'react';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import dynamic from 'next/dynamic';
import { cn } from 'utils/common';
import { SPRITE_CATEGORY_BY_ID } from 'components/SvgSpriteLoader/svgSpriteLoader.constants';
import 'external-svg-loader';

export interface SvgSpriteLoaderProps {
  width?: number;
  size?: number;
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

const SvgSpriteLoader: React.FC<SvgSpriteLoaderProps> = ({
  size,
  width = 20,
  height = 20,
  viewBox = '0 0 24 24',
  fillColor = 'none',
  color,
  iconCategory,
  id = 'alert-circle',
  version = 0,
  domain = 'https://assets.zamp.finance',
  dataCache = '',
  className = '',
  onClick,
  customSpriteUrl,
}) => {
  const category = iconCategory ?? SPRITE_CATEGORY_BY_ID[id];

  return (
    <div
      onClick={onClick}
      className={cn(className, onClick && 'cursor-pointer')}
      data-testid={`svg-sprite-loader-${id}`}
    >
      <svg
        id={id}
        viewBox={viewBox}
        width={size ?? width}
        height={size ?? height}
        fill={fillColor}
        color={color}
        data-src={customSpriteUrl ?? `${domain}/sprites/v${version}/${category}.svg#${id}`}
        data-cache={dataCache}
      ></svg>
    </div>
  );
};

// Export with memoization
const MemoizedSvgSpriteLoader = memo(SvgSpriteLoader);

// Disable SSR using dynamic import
export default dynamic(() => Promise.resolve(MemoizedSvgSpriteLoader), { ssr: false });
