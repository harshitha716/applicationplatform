import { FC, memo } from 'react';
import { cn } from 'utils/common';

interface PivotAutoGroupHeaderPropsType {
  title: string;
  isSingleValue: boolean;
}

const PivotAutoGroupHeader: FC<PivotAutoGroupHeaderPropsType> = ({ title, isSingleValue = false }) => {
  return (
    <div
      className={cn(
        'bg-white w-full h-full f-18-450 p-6 flex items-start border-b-0.5 border-b-GRAY_400 border-r-0.5 border-r-GRAY_400',
        isSingleValue && 'items-center',
      )}
    >
      {title}
    </div>
  );
};

export default memo(PivotAutoGroupHeader);
