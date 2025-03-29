import React from 'react';
import SkeletonElement from 'components/skeletons/SkeletonElement';

const SkeletonLoaderFileHistory = () => {
  return (
    <div className='flex flex-col flex-wrap justify-start items-start w-full '>
      {Array.from({ length: 3 }).map((_, index) => (
        <div key={index} className='grid grid-cols-1 gap-4 w-full items-center border-b border-GRAY_400 py-3.5'>
          <div className='flex flex-col justify-start rounded-md w-full'>
            <div className='flex justify-between items-center w-full'>
              <SkeletonElement key={index} className='h-7 w-30 rounded-md bg-GRAY_500' />
              <SkeletonElement key={index} className='h-4 w-4 rounded-sm bg-GRAY_500' />
            </div>
            <SkeletonElement key={index} className='h-3 w-50 rounded-md bg-GRAY_500 mt-1' />
          </div>
        </div>
      ))}
    </div>
  );
};

export default SkeletonLoaderFileHistory;
