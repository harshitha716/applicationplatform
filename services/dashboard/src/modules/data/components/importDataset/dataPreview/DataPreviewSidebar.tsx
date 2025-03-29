import React, { FC } from 'react';
import { usePostAiTransformationConfirmMutation } from 'apis/dataset';
import { COLORS } from 'constants/colors';
import {
  AI_TRANSFORMATION_STATUS,
  FILE_IMPORT_STATUS_MSG,
} from 'modules/data/components/importDataset/importData.constants';
import { DataPreviewSidebarPropsType } from 'modules/data/components/importDataset/importData.types';
import { useRouter } from 'next/router';
import { SIZE_TYPES } from 'types/common/components';
import { BUTTON_TYPES } from 'types/components/button.type';
import { Button } from 'components/common/button/Button';
import { toast } from 'components/common/toast/Toast';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

const DataPreviewSidebar: FC<DataPreviewSidebarPropsType> = ({
  fileName,
  onReset,
  fileUploadId,
  onRefetch,
  setShowAiTransformationStatus,
}) => {
  const router = useRouter();
  const datasetId = router?.query?.id as string;
  const [postAiTransformationConfirm, { isLoading: isLoadingPostAiTransformationConfirm }] =
    usePostAiTransformationConfirmMutation();

  const handleConfirmImport = () => {
    postAiTransformationConfirm({ file_upload_id: fileUploadId, dataset_id: datasetId })
      .unwrap()
      .then(() => {
        toast.success(FILE_IMPORT_STATUS_MSG.FILE_IMPORT_AFTER_AI_SUCCESS);
        setShowAiTransformationStatus({
          open: true,
          status: AI_TRANSFORMATION_STATUS.STATUS_INGESTION_ONGOING.status,
          title: AI_TRANSFORMATION_STATUS.STATUS_INGESTION_ONGOING.title,
          description: AI_TRANSFORMATION_STATUS.STATUS_INGESTION_ONGOING.description,
        });
        onRefetch();
        onReset();
      })
      .catch(() => {
        toast.error(FILE_IMPORT_STATUS_MSG.FILE_IMPORT_DATA_FAILED);
      });
  };

  return (
    <div className='flex flex-col justify-between h-full '>
      <div className='flex flex-col px-6 pt-6'>
        <span className='f-16-600'>Import Data</span>
        <div className='flex flex-col gap-2 w-full'>
          <div className='flex justify-start items-center mt-6 gap-1.5'>
            <SvgSpriteLoader id='file-06' width={14} height={14} color={COLORS.GRAY_1000} />
            <div className='flex w-full justify-between'>
              <span className='f-12-400'>{fileName}</span>
              <SvgSpriteLoader id='check' width={14} height={14} color={COLORS.GREEN_PRIMARY} />
            </div>
          </div>
          <div className='w-full bg-GREEN_700 h-1 rounded-lg'></div>
        </div>
      </div>
      <div className='p-6 pb-12 flex justify-between items-center border-t border-GRAY_400'>
        <span onClick={onReset} className='f-13-500 text-GRAY_1000 cursor-pointer'>
          Discard
        </span>
        <Button
          id='import-confirm-import'
          className='tw-min-w-[70px]'
          size={SIZE_TYPES.SMALL}
          type={BUTTON_TYPES.PRIMARY}
          isLoading={isLoadingPostAiTransformationConfirm}
          onClick={handleConfirmImport}
        >
          Import
        </Button>
      </div>
    </div>
  );
};

export default DataPreviewSidebar;
