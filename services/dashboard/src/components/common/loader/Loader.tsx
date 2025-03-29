/**Temp file not part of design system */
import React from 'react';
import { SIZE } from 'constants/common.constants';

interface LoaderProps {
  className?: string;
  size?: SIZE;
}

export const SIZE_CLASSNAMES = {
  [SIZE.XSMALL]:
    'w-3 h-3 border-[2px] border-transparent border-t-[2px] border-r-[2px] border-l-[2px] border-t-blue-500 border-r-blue-500 border-l-blue-500 rounded-full',
  [SIZE.SMALL]:
    'w-6 h-6 border-[3px] border-transparent border-t-[3px] border-r-[3px] border-l-[3px] border-t-blue-500 border-r-blue-500 border-l-blue-500 rounded-full',
  [SIZE.MEDIUM]:
    'w-8 h-8 border-[4px] border-transparent border-t-[4px] border-r-[4px] border-l-[4px] border-t-blue-500 border-r-blue-500 border-l-blue-500 rounded-full',
  [SIZE.LARGE]:
    'w-12 h-12 border-[5px] border-transparent border-t-[5px] border-r-[5px] border-l-[5px] border-t-blue-500 border-r-blue-500 border-l-blue-500 rounded-full',
  [SIZE.XLARGE]:
    'w-12 h-12 border-[5px] border-transparent border-t-[5px] border-r-[5px] border-l-[5px] border-t-blue-500 border-r-blue-500 border-l-blue-500 rounded-full',
};

export const Loader: React.FC<LoaderProps> = ({ className = '', size = SIZE.MEDIUM }) => {
  return <div className={`loader ${SIZE_CLASSNAMES[size]} animate-spin ${className}`}></div>;
};
