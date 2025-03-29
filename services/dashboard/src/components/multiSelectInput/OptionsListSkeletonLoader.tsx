import React from 'react';
import SkeletonElement from 'components/skeletons/SkeletonElement';

const OptionsListSkeletonLoader = () => {
  return (
    <div className='pl-1'>
      <SkeletonElement className='h-4 w-40 rounded-md bg-GRAY_50' />
    </div>
  );
};

export default OptionsListSkeletonLoader;
