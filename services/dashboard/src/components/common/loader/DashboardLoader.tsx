import React, { FC } from 'react';
import { DASHBOARD_LOADER } from 'constants/lottie/dashboard_loader';
import { cn } from 'utils/common';
import SkeletonElement from 'components/common/skeletons/SkeletonElement';
import DynamicLottiePlayer from 'components/DynamicLottiePlayer';

type DashboardLoaderPropsType = {
  isFadingOut: boolean;
};

const DashboardLoader: FC<DashboardLoaderPropsType> = ({ isFadingOut }) => {
  return (
    <div
      className={cn(
        'fixed inset-0 flex flex-col justify-between items-center transition duration-500 ease-out z-1000 bg-BACKGROUND_GRAY_1',
        isFadingOut ? 'opacity-100 pointer-events-none  translate-x-40' : 'opacity-100 translate-x-0',
      )}
    >
      <div className='h-12 py-4 px-8 flex items-center justify-between w-full border-2 border-GRAY_400 bg-white'>
        <div className='flex gap-4'>
          <SkeletonElement elementCount={3} className='w-4 h-4 rounded-sm bg-GRAY_400' />
          <SkeletonElement elementCount={1} className='w-20 h-4 rounded-sm bg-GRAY_400' />
        </div>
        <span className='w-60 h-5 bg-GRAY_400 rounded-sm relative overflow-hidden before:absolute before:inset-0 before:bg-gradient-to-r before:from-transparent before:via-white/60 before:to-transparent before:animate-[shimmer-skeleton_1.5s_infinite] before:w-full before:h-full'></span>

        <div className='flex gap-4'>
          <SkeletonElement elementCount={1} className='w-8 h-5 rounded-sm bg-GRAY_400' />
          <SkeletonElement elementCount={1} className='w-12 h-5 rounded-sm bg-GRAY_400' />
        </div>
      </div>
      <div className={cn('transition-transform duration-500 delay-150', isFadingOut ? 'scale-150' : 'scale-100')}>
        <DynamicLottiePlayer
          src={DASHBOARD_LOADER}
          className='lottie-player'
          autoplay
          keepLastFrame
          style={{ height: '400px' }}
        />
      </div>
      <div></div>
    </div>
  );
};

export default DashboardLoader;
