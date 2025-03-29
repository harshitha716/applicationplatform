import React, { useState } from 'react';
import { useSelector } from 'react-redux';
import { COLORS } from 'constants/colors';
import ImportFileHistory from 'modules/data/components/datasetHistory/ImportFileHistory';
import LoadingWidthAnimation from 'modules/data/components/LoadingWidthAnimation';
import { RootState } from 'store';
import { cn } from 'utils/common';
import { Tooltip, TooltipPositions } from 'components/common/tooltip';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

const DatasetHistory = () => {
  const datasetBulkLoaders = useSelector((state: RootState) => state?.user?.datasetBulkLoaders) || [];
  const [isFileImportHistoryOpen, setIsFileImportHistoryOpen] = useState<boolean>(false);

  const handleOpenFileImportHistory = () => setIsFileImportHistoryOpen(true);
  const handleCloseFileImportHistory = () => setIsFileImportHistoryOpen(false);

  return (
    <div>
      {isFileImportHistoryOpen && <ImportFileHistory onClose={handleCloseFileImportHistory} />}
      <div className='relative z-[800]'>
        <Tooltip
          tooltipBody='Activity'
          position={TooltipPositions.BOTTOM}
          tooltipBodyClassName='f-12-300 rounded-md whitespace-nowrap z-[1000] bg-black text-GRAY_200'
          className='z-1 h-full w-full'
          tooltipBodystyle='f-10-400'
        >
          <div
            className={cn('p-1 hover:bg-GRAY_100 rounded cursor-pointer', isFileImportHistoryOpen && 'bg-GRAY_100')}
            onClick={handleOpenFileImportHistory}
          >
            <SvgSpriteLoader id='clock-rewind' width={14} height={14} color={COLORS.GRAY_900} />
          </div>
        </Tooltip>
        {!!datasetBulkLoaders.length && (
          <div className='absolute bottom-px left-[3px]'>
            <LoadingWidthAnimation />
          </div>
        )}
      </div>
    </div>
  );
};

export default DatasetHistory;
