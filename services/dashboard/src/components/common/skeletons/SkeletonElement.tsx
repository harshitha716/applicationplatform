import React, { FC } from 'react';
import { SKELETON_ELEMENT_SHAPES } from 'components/common/skeletons/skeletons.types';

interface SkeletonElementProps {
  elementCount?: number;
  className?: string;
  shape?: (typeof SKELETON_ELEMENT_SHAPES)[keyof typeof SKELETON_ELEMENT_SHAPES];
}

const SHAPE_STYLE = {
  [SKELETON_ELEMENT_SHAPES.CIRCLE]: 'rounded-full',
};

const SkeletonElement: FC<SkeletonElementProps> = ({ elementCount = 1, className = '', shape = null }) => {
  const elements = Array(elementCount)?.fill('');

  return elements?.map((_, index) => {
    return (
      <span
        className={`block animate-pulse bg-BASE_PRIMARY ${className} ${shape ? (SHAPE_STYLE[shape] ?? '') : ''}`}
        key={index}
      >
        &zwnj;
      </span>
    );
  });
};

export default SkeletonElement;
