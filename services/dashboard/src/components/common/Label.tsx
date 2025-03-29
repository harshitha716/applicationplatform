import React, { FC, ReactElement } from 'react';

export interface LabelProps {
  title?: string | ReactElement | number | null;
  description?: string | ReactElement | number | null;
  showTitle?: boolean;
  showDescription?: boolean;
  wrapperClassName?: string;
  titleClassName?: string;
  descriptionClassName?: string;
}

export const Label: FC<LabelProps> = ({
  title = null,
  description = null,
  wrapperClassName = 'w-full',
  titleClassName = 'f-14-500 text-GRAY_700 mb-[4px]',
  descriptionClassName = 'f-12-300 text-GRAY_600',
}) =>
  title || descriptionClassName ? (
    <div className={wrapperClassName}>
      {title && <div className={titleClassName}>{title}</div>}
      {description && <div className={descriptionClassName}>{description}</div>}
    </div>
  ) : null;
