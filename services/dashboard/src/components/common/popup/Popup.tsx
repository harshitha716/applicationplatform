import React, { FC } from 'react';
import { cn, stopPropagationAction } from 'utils/common';
import { PopupProps } from 'components/common/popup/popup.types';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

const Popup: FC<PopupProps> = ({
  title,
  subTitle,
  titleClassName,
  subTitleClassName,
  popupWrapperClassName,
  showIcon,
  iconCategory,
  iconId,
  iconColor,
  isOpen = false,
  children,
  className = '',
  onClose,
  closeOnClickOutside = true,
  isOverlay = true,
  wrapperClassName,
}) => {
  if (!isOpen) return null;

  return (
    <div
      className={cn(
        `transition-all duration-300 ease-in fixed w-screen h-screen z-1000 top-0 left-0 ${isOverlay ? 'bg-GRAY_70 ' : ''} ${
          isOpen ? 'opacity-1' : 'hidden opacity-0'
        }`,
      )}
      role='presentation'
      onClick={() => {
        if (closeOnClickOutside) onClose?.();
      }}
    >
      <div className={cn('w-full h-full flex items-center justify-center', wrapperClassName)}>
        <div
          className={cn(
            `transition-all duration-300 ease-in px-5 py-5 rounded-xl block ${className} ${
              isOpen ? ' translate-y-0 opacity-1' : 'translate-y-[50px] opacity-0'
            }`,
          )}
          role='presentation'
          onClick={stopPropagationAction}
        >
          <div className={cn('flex w-full justify-between items-center px-5 pt-5 pb-0', popupWrapperClassName)}>
            <div className='flex flex-col'>
              {title && <span className={cn('f-16-600 text-GRAY_950', titleClassName)}>{title}</span>}
              {subTitle && <span className={cn('f-12-400 text-GRAY_700 mt-1', subTitleClassName)}>{subTitle}</span>}
            </div>
            {showIcon && (
              <div className='p-1 cursor-pointer' onClick={onClose}>
                <SvgSpriteLoader id={iconId} iconCategory={iconCategory} width={16} height={16} color={iconColor} />
              </div>
            )}
          </div>
          {children}
        </div>
      </div>
    </div>
  );
};

export default Popup;
