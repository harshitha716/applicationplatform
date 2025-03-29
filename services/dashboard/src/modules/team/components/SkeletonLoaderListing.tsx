import React, { FC } from 'react';
import SkeletonElement from 'components/skeletons/SkeletonElement';

type SkeletonLoaderListingPropsType = {
  columns?: number;
  length?: number;
};

const SkeletonLoaderListing: FC<SkeletonLoaderListingPropsType> = ({ columns = 3, length = 3 }) => {
  return (
    <div className='flex flex-col gap-4'>
      {Array.from({ length: length }).map((_, index) => (
        <div
          key={index}
          className={`grid grid-cols-${columns} gap-4 w-full items-center h-10 border-b-0.5 border-DIVIDER_GRAY`}
        >
          {Array.from({ length: columns }).map((_, colIndex) => (
            <SkeletonElement key={colIndex} className='h-4 rounded-md w-1/3' />
          ))}
        </div>
      ))}
    </div>
  );
};

export default SkeletonLoaderListing;
