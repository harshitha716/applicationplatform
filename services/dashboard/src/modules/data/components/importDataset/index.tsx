import React, { FC, useEffect, useState } from 'react';
import { useDispatch } from 'react-redux';
import { useLazyGetActionStatusQuery, useLazyGetAiTransformationQuery } from 'apis/dataset';
import { COINS_STACKED_05 } from 'constants/icons';
import usePolling from 'hooks/usePolling';
import ImportedDataPreview from 'modules/data/components/importDataset/dataPreview';
import { AI_TRANSFORMATION_STATUS } from 'modules/data/components/importDataset/importData.constants';
import {
  ImportDatasetPropsType,
  StartPollingPreviewType,
} from 'modules/data/components/importDataset/importData.types';
import ImportFileWrapper from 'modules/data/components/importDataset/ImportFileWrapper';
import Image from 'next/image';
import { useRouter } from 'next/router';
import { addDatasetBulkLoaders, removeDatasetBulkLoader } from 'store/slices/user';
import { DatasetActionStatusResponseType, RawMetadata, TransformationPreviewMetadata } from 'types/api/dataset.types';
import { cn } from 'utils/common';
import { Tooltip, TooltipPositions } from 'components/common/tooltip';

const ImportDataset: FC<ImportDatasetPropsType> = ({ setShowAiTransformationStatus, onRefetch }) => {
  const router = useRouter();
  const dispatch = useDispatch();
  const datasetId = router?.query?.id as string;
  const { startPolling } = usePolling();
  const [getActionStatus] = useLazyGetActionStatusQuery();
  const [getAiTransformation] = useLazyGetAiTransformationQuery();
  const [rawData, setRawData] = useState<RawMetadata | null>(null);
  const [startAiTransformation, setStartAiTransformation] = useState<boolean>(false);
  const [isImportFilePopupOpen, setIsImportFilePopupOpen] = useState<boolean>(false);
  const [mappedData, setMappedData] = useState<TransformationPreviewMetadata | null>(null);
  const [startPollingPreview, setStartPollingPreview] = useState<StartPollingPreviewType>({
    check: false,
    actionId: '',
    fileUploadId: '',
  });
  const [fileName, setFileName] = useState<string | null>(null);

  const handleOpenImportFilePopup = () => setIsImportFilePopupOpen(true);
  const handleCloseImportFilePopup = () => setIsImportFilePopupOpen(false);

  const handleReset = () => {
    setRawData(null);
    setMappedData(null);
    setStartAiTransformation(false);
    handleCloseImportFilePopup();
  };

  const handleShowAiTransformationStatus = () => {
    setShowAiTransformationStatus({
      open: true,
      status: AI_TRANSFORMATION_STATUS.STATUS_LOADING.status,
      title: AI_TRANSFORMATION_STATUS.STATUS_LOADING.title,
      description: AI_TRANSFORMATION_STATUS.STATUS_LOADING.description,
    });

    dispatch(
      addDatasetBulkLoaders({
        id: startPollingPreview?.actionId,
        status: AI_TRANSFORMATION_STATUS.STATUS_LOADING.status,
        title: AI_TRANSFORMATION_STATUS.STATUS_LOADING.title,
        description: AI_TRANSFORMATION_STATUS.STATUS_LOADING.description,
      }),
    );

    startPolling({
      fn: () => getActionStatus({ datasetId, params: { action_ids: [startPollingPreview?.actionId] } }),
      validate: (data: DatasetActionStatusResponseType[]) => data.filter((item) => !item.is_completed).length === 0,
      interval: 3000,
      maxAttempts: 50,
    })
      .then(() => {
        return getAiTransformation({ file_upload_id: startPollingPreview?.fileUploadId }).unwrap();
      })
      .then((data) => {
        dispatch(removeDatasetBulkLoader(startPollingPreview?.actionId));
        setShowAiTransformationStatus({
          open: true,
          status: AI_TRANSFORMATION_STATUS.STATUS_SUCCESS.status,
          title: AI_TRANSFORMATION_STATUS.STATUS_SUCCESS.title,
          description: AI_TRANSFORMATION_STATUS.STATUS_SUCCESS.description,
        });
        setIsImportFilePopupOpen(true);
        setStartAiTransformation(true);

        setMappedData({
          data_preview: data?.data_preview,
        });
      })
      .catch(() => {
        handleReset();
        dispatch(removeDatasetBulkLoader(startPollingPreview?.actionId));
        setShowAiTransformationStatus({
          open: true,
          status: AI_TRANSFORMATION_STATUS.STATUS_ERROR.status,
          title: AI_TRANSFORMATION_STATUS.STATUS_ERROR.title,
          description: AI_TRANSFORMATION_STATUS.STATUS_ERROR.description,
        });
      });
  };

  useEffect(() => {
    if (startPollingPreview?.check) {
      handleShowAiTransformationStatus();
    }
  }, [startPollingPreview]);

  return (
    <div className='z-1000'>
      {isImportFilePopupOpen && !startAiTransformation ? (
        <ImportFileWrapper
          fileName={fileName}
          setFileName={setFileName}
          isOpen={isImportFilePopupOpen}
          setRawData={setRawData}
          setStartPollingPreview={setStartPollingPreview}
          onReset={handleReset}
          onClose={handleCloseImportFilePopup}
        />
      ) : isImportFilePopupOpen && startAiTransformation ? (
        <ImportedDataPreview
          fileName={fileName}
          onReset={handleReset}
          rawData={rawData}
          mappedData={mappedData}
          startAiTransformation={startAiTransformation}
          setShowAiTransformationStatus={setShowAiTransformationStatus}
          fileUploadId={startPollingPreview?.fileUploadId}
          onRefetch={onRefetch}
        />
      ) : null}
      <Tooltip
        tooltipBody='Import Data'
        position={TooltipPositions.BOTTOM}
        tooltipBodyClassName='f-12-300 rounded-md whitespace-nowrap z-[1000] bg-black text-GRAY_200'
        className='z-1 h-full w-full'
        tooltipBodystyle='f-10-400'
      >
        <div className={cn('p-1 hover:bg-GRAY_100 !rounded cursor-pointer', isImportFilePopupOpen && 'bg-GRAY_100')}>
          <Image
            src={COINS_STACKED_05}
            alt='coins-stacked-05'
            width={14}
            height={14}
            className='text-GRAY_900'
            onClick={handleOpenImportFilePopup}
          />
        </div>
      </Tooltip>
    </div>
  );
};

export default ImportDataset;
