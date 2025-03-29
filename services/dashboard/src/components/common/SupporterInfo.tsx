import React, { FC, ReactElement } from 'react';
import { SUPPORT_INFO_TYPES } from 'types/common/components/input/input.types';
import SvgSpriteLoader, { SvgSpriteLoaderProps } from 'components/SvgSpriteLoader';

export interface SupporterInfoProps {
  type?: SUPPORT_INFO_TYPES;
  icon?: SvgSpriteLoaderProps;
  text?: string | ReactElement | null;
  className?: string;
  textClass?: string;
  showSupportInfo?: boolean;
}

const TYPE_PROP = {
  [SUPPORT_INFO_TYPES.GUIDE]: {
    iconByType: null,
    textClassByType: 'f-12-300 text-GRAY_600',
  },
  [SUPPORT_INFO_TYPES.ERROR]: {
    iconByType: null,
    textClassByType: 'f-12-300 text-RED_PRIMARY',
  },
  [SUPPORT_INFO_TYPES.CUSTOM]: {
    iconByType: null,
    iconColorByType: '',
    textClassByType: '',
  },
};

export const SupporterInfo: FC<SupporterInfoProps> = ({
  type = SUPPORT_INFO_TYPES.GUIDE,
  icon = null,
  text = null,
  className = 'flex items-center',
  textClass = null,
}) => {
  const { iconByType, textClassByType } = TYPE_PROP[type] || {};
  const typeIcon = icon ?? iconByType ?? null;
  const typeTextClass = textClass ?? textClassByType ?? '';

  return icon?.id || text ? (
    <div className={`${className} ${typeTextClass}`}>
      {typeIcon?.id && (
        <div className={typeIcon.className ?? 'mr-2 w-3 h-3'}>
          <SvgSpriteLoader color={icon?.color} {...typeIcon} id={icon?.id ?? typeIcon?.id ?? ''} />
        </div>
      )}
      {text && <div className={typeTextClass}>{text}</div>}
    </div>
  ) : null;
};
