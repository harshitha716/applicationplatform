import React from 'react';

const LoadingWidthAnimation = () => {
  return (
    <div className='relative'>
      <div className='w-4 border border-GRAY_400 rounded-full'></div>
      <div className='absolute top-0 w-2 border border-GRAY_1000 rounded-full animate-width'></div>
    </div>
  );
};

export default LoadingWidthAnimation;
