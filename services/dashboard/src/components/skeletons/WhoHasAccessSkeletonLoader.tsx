import React from 'react';
import SkeletonElement from 'components/skeletons/SkeletonElement';

const WhoHasAccessSkeletonLoader = () => {
  return (
    <div className='flex justify-between w-full pl-2 pr-1 mt-2'>
      <SkeletonElement className='h-4 w-32 rounded-md bg-GRAY_50' />
      <SkeletonElement className='h-4 w-16 rounded-md bg-GRAY_50' />
    </div>
  );
};

export default WhoHasAccessSkeletonLoader;
