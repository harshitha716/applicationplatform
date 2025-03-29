import { FC, useMemo } from 'react';
import { toast as reactToastify, ToastOptions } from 'react-toastify';
import { COLORS } from 'constants/colors';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { CustomToastPropsType } from 'components/common/toast/toast.types';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

const CustomToast: FC<CustomToastPropsType> = ({ text = '' }) => {
  return <div className='f-14-400 flex w-full items-center gap-6 text-GRAY_700 -ml-[8px] -mt-[2px]'>{text}</div>;
};

const closeToast = () => {
  reactToastify.dismiss();
};

const GetToastIcon = (type: string) => {
  return useMemo(() => {
    switch (type) {
      case 'success':
        return (
          <SvgSpriteLoader
            id='check-circle'
            iconCategory={ICON_SPRITE_TYPES.GENERAL}
            width={16}
            height={16}
            color={COLORS.GREEN_PRIMARY}
          />
        );
      case 'error':
        return (
          <SvgSpriteLoader
            id='x-circle'
            iconCategory={ICON_SPRITE_TYPES.GENERAL}
            width={16}
            height={16}
            color={COLORS.RED_PRIMARY}
          />
        );
      case 'warning':
        return (
          <SvgSpriteLoader
            id='alert-circle'
            iconCategory={ICON_SPRITE_TYPES.ALERTS_AND_FEEDBACK}
            width={16}
            height={16}
            color={COLORS.ORANGE_SECONDARY}
          />
        );
      default:
        return (
          <SvgSpriteLoader
            id='alert-circle'
            iconCategory={ICON_SPRITE_TYPES.ALERTS_AND_FEEDBACK}
            width={16}
            height={16}
            color={COLORS.GREEN_PRIMARY}
          />
        );
    }
  }, [type]);
};

const defaultToastOptions: ToastOptions = {
  position: 'bottom-right',
  autoClose: 3000,
  hideProgressBar: true,
  closeOnClick: true,
  pauseOnHover: true,
  closeButton() {
    return (
      <SvgSpriteLoader
        id='x-close'
        iconCategory={ICON_SPRITE_TYPES.GENERAL}
        width={16}
        height={16}
        color={COLORS.GRAY_900}
        onClick={closeToast}
      />
    );
  },
  icon: ({ type }) => GetToastIcon(type),
  style: {
    marginRight: '24px',
    marginBottom: '24px',
    padding: '20px',
    minHeight: '57px',
    minWidth: '420px',
    border: `1px solid ${COLORS.GRAY_400}`,
    borderRadius: '10px',
    boxShadow: `1px 2px 10px ${COLORS.TOAST_SHADOW}`,
    display: 'flex',
    alignItems: 'start',
  },
};

export const toast = {
  success: (message: string, options?: ToastOptions) => {
    reactToastify.success(<CustomToast text={message} type='success' />, { ...defaultToastOptions, ...options });
  },
  error: (message: string, options?: ToastOptions) => {
    reactToastify.error(<CustomToast text={message} type='error' />, { ...defaultToastOptions, ...options });
  },
  warn: (message: string, options?: ToastOptions) => {
    reactToastify.warn(<CustomToast text={message} type='warn' />, { ...defaultToastOptions, ...options });
  },
};
