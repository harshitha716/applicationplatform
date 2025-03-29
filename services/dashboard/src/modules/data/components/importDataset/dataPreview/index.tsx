import React, { FC } from 'react';
import DataPreviewContent from 'modules/data/components/importDataset/dataPreview/DataPreviewContent';
import DataPreviewSidebar from 'modules/data/components/importDataset/dataPreview/DataPreviewSidebar';
import { ImportedDataPreviewPropsType } from 'modules/data/components/importDataset/importData.types';

const ImportedDataPreview: FC<ImportedDataPreviewPropsType> = ({
  onReset,
  rawData,
  mappedData,
  startAiTransformation,
  setShowAiTransformationStatus,
  fileUploadId,
  fileName,
  onRefetch,
}) => {
  const handleReset = () => {
    onReset();
  };

  if (startAiTransformation) {
    return (
      <div className='fixed w-screen h-screen z-1000 top-0 left-0 bg-GRAY_70'>
        <div className='h-screen bg-white mt-7 rounded-2.5 border border-t border-GRAY_400 flex overflow-hidden'>
          <div className='sticky w-1/3 h-full border-r border-GRAY_400'>
            <DataPreviewSidebar
              fileName={fileName}
              fileUploadId={fileUploadId}
              setShowAiTransformationStatus={setShowAiTransformationStatus}
              onReset={handleReset}
              onRefetch={onRefetch}
            />
          </div>
          <div className='w-2/3 h-full'>
            <DataPreviewContent mappedData={mappedData} rawData={rawData} />
          </div>
        </div>
      </div>
    );
  }

  return null;
};

export default ImportedDataPreview;
