import React, { FC, useEffect } from 'react';
import { TableSchemaAlignmentStatusPropsType } from 'modules/data/components/importDataset/importData.types';
import { LOADER_STATUS } from 'modules/data/data.types';
import { cn } from 'utils/common';
import StatusIndicator from 'components/common/StatusIndicator';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

const TableSchemaAlignmentStatus: FC<TableSchemaAlignmentStatusPropsType> = ({
  showAiTransformationStatus,
  setShowAiTransformationStatus,
}) => {
  const handleCloseTransformationStatus = () => {
    setShowAiTransformationStatus({
      open: false,
      status: LOADER_STATUS.SUCCESS,
      title: '',
      description: '',
    });
  };

  useEffect(() => {
    if (showAiTransformationStatus?.open) {
      const timer = setTimeout(handleCloseTransformationStatus, 5000);

      return () => clearTimeout(timer);
    }
  }, [showAiTransformationStatus, setShowAiTransformationStatus]);

  if (!showAiTransformationStatus?.open) return null;

  return (
    <div className='animate-slideInOut absolute flex justify-center items-start gap-5 right-0 top-8 border-[0.5px] border-GRAY_500 rounded-2.5 py-3 px-5 w-[300px] bg-white z-1000 shadow-tableFilterMenu'>
      <div className='flex gap-3 items-start'>
        <StatusIndicator status={showAiTransformationStatus?.status as LOADER_STATUS} />
        <div className='flex flex-col'>
          <span className='f-13-500 text-GRAY_1000'>{showAiTransformationStatus?.title}</span>
          <span
            className={cn(
              'f-11-400 text-GRAY_700 mt-1',
              showAiTransformationStatus?.status === LOADER_STATUS.ERROR && 'text-RED_800',
            )}
          >
            {showAiTransformationStatus?.description}
          </span>
        </div>
      </div>
      <SvgSpriteLoader
        id='x-close'
        width={16}
        height={16}
        onClick={handleCloseTransformationStatus}
        className='text-GRAY_800 hover:text-GRAY_1000 cursor-pointer'
      />
    </div>
  );
};

export default TableSchemaAlignmentStatus;
