import React, { FC, useRef, useState } from 'react';
import { useSelector } from 'react-redux';
import { useGetFileImportHistoryQuery } from 'apis/dataset';
import { useOnClickOutside } from 'hooks';
import HistoryBulkLoaders from 'modules/data/components/datasetHistory/HistoryBulkLoaders';
import HistoryEmptyState from 'modules/data/components/datasetHistory/HistoryEmptyState';
import HistoryList from 'modules/data/components/datasetHistory/HistoryList';
import { ImportFileHistoryPropsType } from 'modules/data/components/importDataset/importData.types';
import { useRouter } from 'next/router';
import { RootState } from 'store';

const ImportFileHistory: FC<ImportFileHistoryPropsType> = ({ onClose }) => {
  const importFileHistoryRef = useRef<HTMLDivElement>(null);
  const [isHoveredLoaders, setIsHoveredLoaders] = useState(false);
  const datasetBulkLoaders = useSelector((state: RootState) => state?.user?.datasetBulkLoaders) || [];
  const router = useRouter();
  const datasetId = router?.query?.id as string;
  const { data } = useGetFileImportHistoryQuery({ datasetId });
  const fileImportHistoryData = data?.file_uploads || [];

  useOnClickOutside(importFileHistoryRef, onClose);

  return (
    <>
      <div className='fixed w-screen !h-[calc(100vh-136px)] z-1000 top-[94px] left-0 flex justify-end'>
        {!datasetBulkLoaders?.length && !fileImportHistoryData?.length ? (
          <div className='absolute right-8 top-0 z-50' ref={importFileHistoryRef}>
            <HistoryEmptyState />
          </div>
        ) : (
          <div
            ref={importFileHistoryRef}
            className='h-full overflow-y-scroll [&::-webkit-scrollbar]:hidden'
            onMouseLeave={() => setIsHoveredLoaders(false)}
          >
            <div className='flex flex-col h-auto'>
              <HistoryBulkLoaders
                isHoveredLoaders={isHoveredLoaders}
                setIsHoveredLoaders={setIsHoveredLoaders}
                datasetBulkLoaders={datasetBulkLoaders}
              />
              <HistoryList isHoveredLoaders={isHoveredLoaders} fileImportHistoryData={fileImportHistoryData} />
            </div>
          </div>
        )}
      </div>
    </>
  );
};

export default ImportFileHistory;
