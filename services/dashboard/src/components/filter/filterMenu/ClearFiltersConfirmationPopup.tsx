import React, { FC } from 'react';
import { defaultFnType } from 'types/commonTypes';
import { cn } from 'utils/common';

interface ClearFiltersConfirmationPopupProps {
  className?: string;
  onClick: defaultFnType;
  onCancel: defaultFnType;
  containerRef: React.RefObject<HTMLDivElement>;
}

const ClearFiltersConfirmationPopup: FC<ClearFiltersConfirmationPopupProps> = ({
  className = '',
  onClick,
  onCancel,
  containerRef,
}) => {
  return (
    <div
      ref={containerRef}
      className={cn('p-4 z-30 bg-white border-0.5 border-GRAY_400 rounded-md top-full mt-1', className)}
    >
      <div className='mb-3'>Remove all filters?</div>

      <div className='flex'>
        <button
          className='hover:border-DIVIDER_SAIL_4 border border-DIVIDER_SAIL_2 outline-none rounded-lg p-1.5 min-w-17.5 mr-3'
          onClick={onClick}
          data-testid='clear-filters-confirmation-popup-yes'
        >
          Yes
        </button>
        <button
          className='hover:border-DIVIDER_SAIL_4 border border-DIVIDER_SAIL_2 outline-none rounded-lg p-1.5 min-w-17.5'
          onClick={onCancel}
        >
          No
        </button>
      </div>
    </div>
  );
};

export default ClearFiltersConfirmationPopup;
