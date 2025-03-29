import React from 'react';
import SkeletonElement from 'components/skeletons/SkeletonElement';

const SkeletonLoaderSidebarPages = () => {
  return (
    <div className='flex flex-col gap-3 w-full ml-1'>
      <SkeletonElement className='h-4 rounded bg-GRAY_400 w-1/3' />
      <SkeletonElement className='h-4 rounded bg-GRAY_400 w-1/2' />
      <SkeletonElement className='h-4 rounded bg-GRAY_400 w-2/5' />
    </div>
  );
};

export default SkeletonLoaderSidebarPages;
