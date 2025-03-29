import React, { FC } from 'react';
import { COLORS } from 'constants/colors';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { IMPORT_ALLOWED_FILE_FORMATS } from 'modules/data/components/importDataset/importData.constants';
import { ImportFileWrapperPropsType } from 'modules/data/components/importDataset/importData.types';
import ImportFile from 'modules/data/components/importDataset/ImportFile';
import Popup from 'components/common/popup/Popup';

const ImportFileWrapper: FC<ImportFileWrapperPropsType> = ({
  isOpen,
  onClose,
  setStartPollingPreview,
  setRawData,
  fileName,
  setFileName,
}) => {
  const handleClosePopup = () => {
    onClose();
  };

  return (
    <Popup
      title='Import Data'
      titleClassName='f-16-600 text-GRAY_950'
      iconCategory={ICON_SPRITE_TYPES.GENERAL}
      iconId='x-close'
      iconColor={COLORS.TEXT_PRIMARY}
      isOpen={isOpen}
      onClose={onClose}
      popupWrapperClassName='bg-white rounded-t-2.5 min-w-[464px] pl-6 py-4.5 border border-GRAY_400 border-b-0'
      closeOnClickOutside={false}
      showIcon
    >
      <div className='bg-white p-1.5 pt-0 rounded-2.5 rounded-t-none border border-GRAY_400 border-t-0'>
        <ImportFile
          fileName={fileName}
          setFileName={setFileName}
          setRawData={setRawData}
          setStartPollingPreview={setStartPollingPreview}
          acceptedFormats={IMPORT_ALLOWED_FILE_FORMATS}
          className='flex flex-col justify-center items-center bg-BG_GRAY_1 border border-dashed border-GRAY_400 min-h-[220px] rounded-md focus:border-black focus:border-solid cursor-pointer'
          onClosePopup={handleClosePopup}
        />
      </div>
    </Popup>
  );
};

export default ImportFileWrapper;
