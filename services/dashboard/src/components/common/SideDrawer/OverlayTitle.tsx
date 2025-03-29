import { FC } from 'react';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { cn } from 'utils/common';
import { OverlayTitleProps } from 'components/common/SideDrawer/sideDrawer.types';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

// TODO: needs to be updated to use the new design system
const OverlayTitle: FC<OverlayTitleProps> = ({
  topBar,
  title,
  hideCloseButton,
  step,
  subtitle,
  onClose,
  headerClassName = '',
  closeButtonClassName,
  titleClassName = '',
  subtitleClassName = '',
  closeButtonDimensions = { width: 24, height: 24 },
}) =>
  (topBar || title || !hideCloseButton) && (
    <div
      className={cn('min-h-[56px] px-4 border-b f-16-300 flex justify-between items-center py-2.5', headerClassName)}
    >
      {topBar ? (
        topBar
      ) : (
        <div className='grow'>
          <div className='f-14-500 gap-1 flex items-center '>
            <div className={titleClassName}>{title}</div>
            {step && <div className='text-GRAY_1000 uppercase f-12-300'>step {step}</div>}
          </div>
          {!!subtitle && <div className={cn('f-11-300 text-GRAY_600 mt-1', subtitleClassName)}>{subtitle}</div>}
        </div>
      )}

      {!!onClose && !hideCloseButton && (
        <div className={cn('p-2 rounded-full cursor-pointer', closeButtonClassName)} onClick={onClose}>
          <SvgSpriteLoader id='x-close' iconCategory={ICON_SPRITE_TYPES.GENERAL} {...closeButtonDimensions} />
        </div>
      )}
    </div>
  );

export default OverlayTitle;
