import React, { FC } from 'react';
import { GOOGLE_ICON } from 'constants/icons';
import Image from 'next/image';
import { cn } from 'utils/common';

type LoginButtonPropsType = {
  loading: boolean;
  onClick: (e?: React.MouseEvent<HTMLButtonElement>) => void;
};

const LoginButton: FC<LoginButtonPropsType> = ({ loading, onClick }) => {
  return (
    <button
      id='google-login'
      type='submit'
      className={cn(
        'relative bg-BG_GRAY_3 h-12 w-full mt-4 rounded-md',
        loading ? '!cursor-not-allowed' : '!cursor-pointer',
      )}
      onClick={onClick}
    >
      <div
        className={cn(
          'color-transition before:transform before:translate-x-0 before:bg-BG_GRAY_3 after:transform after:-translate-x-1/2 relative h-full w-full overflow-hidden rounded-md before:absolute before:top-0 before:h-full before:w-full before:transition-transform before:duration-[3000ms] before:ease-in-out before:rounded-[6px] after:absolute after:top-0 after:h-full after:w-full after:transition-transform after:duration-[3000ms] after:ease-in-out after:rounded-[6px] after:bg-BG_GRAY_4',
          { active: loading },
        )}
      ></div>
      <div className='absolute -top-[6px] right-40 text-white'>
        <span
          className={cn(
            '!translate-y-5 flex justify-center items-center gap-1.5 f-14-500',
            loading
              ? '!login-btn-scale-100 !opacity-100 !login-btn-opacity-300-easeInOut '
              : 'login-btn-scale-20 opacity-0 login-btn-opacity-300-easeInOut',
          )}
        >
          Signing in with
          <Image src={GOOGLE_ICON} alt='google_icon' width={20} height={20} />
        </span>
        <span
          className={cn(
            'flex justify-center items-center f-14-500 scale -translate-y-[2px]',
            loading
              ? 'login-btn-scale-20 opacity-0 login-btn-opacity-300-easeInOut'
              : 'login-btn-scale-100 opacity-100 login-btn-opacity-300-easeInOut',
          )}
        >
          Login
        </span>
      </div>
    </button>
  );
};

export default LoginButton;
