import React, { DragEventHandler, FC, useState } from 'react';
import { COLORS } from 'constants/colors';
import { KEYBOARD_KEYS } from 'constants/shortcuts';
import { FILE_IMPORT_STATUS_MSG } from 'modules/data/components/importDataset/importData.constants';
import { FileUploaderPropsType } from 'modules/data/components/importDataset/importData.types';
import { SIZE_TYPES } from 'types/common/components';
import { BUTTON_TYPES } from 'types/components/button.type';
import { cn } from 'utils/common';
import { Button } from 'components/common/button/Button';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

const FileUploader: FC<FileUploaderPropsType> = ({
  isLoading,
  isSuccess,
  error,
  onFileDrop,
  onClick,
  children,
  errorClassName,
  footer,
  className,
  fileName,
  setFileName,
  indexKey,
}) => {
  const errorTitle = error ?? FILE_IMPORT_STATUS_MSG.FILE_UPLOAD_COMMON_ERROR;
  const [isDragOver, setIsDragOver] = useState<boolean>(false);

  const handleFileDrop: DragEventHandler<HTMLDivElement> = (event) => {
    if (isLoading) return null;

    event?.preventDefault();
    event?.stopPropagation();

    const files = event?.dataTransfer?.files?.[0];

    setFileName(files?.name);

    setIsDragOver(false);
    onFileDrop(files);
  };

  const handleKeyDown = (event: React.KeyboardEvent<HTMLDivElement>) => {
    if (event.key === KEYBOARD_KEYS.ENTER) onClick();
  };

  const renderContent = () => {
    if (isLoading) {
      return (
        <div className='w-3/5 h-full'>
          <div className='flex w-full justify-start gap-1.5 py-1'>
            <SvgSpriteLoader id='file-06' width={14} height={14} color={COLORS.GRAY_700} />
            <span className='f-12-400 text-GRAY_700'>{fileName}</span>
          </div>
          <span className='relative flex bg-GRAY_400 h-1 w-full rounded-xl mt-2 overflow-hidden'>
            <span className='absolute left-0 h-full w-1/2 bg-black animate-slide  rounded-xl'></span>
          </span>
        </div>
      );
    }

    if (isSuccess) {
      return (
        <div className='w-3/5'>
          <div className='flex items-center justify-between w-full gap-4'>
            <div className='flex w-full justify-start gap-1.5 py-1'>
              <SvgSpriteLoader id='file-06' width={14} height={14} color={COLORS.GRAY_1000} />
              <span className='f-12-400 text-GRAY_1000'>{fileName}</span>
            </div>
            <SvgSpriteLoader id='check' width={14} height={14} color={COLORS.GREEN_PRIMARY} />
          </div>
          <span className='flex bg-GREEN_700 h-1 w-full rounded-xl mt-2'></span>
        </div>
      );
    }

    return (
      <div key={indexKey} className='relative flex flex-col justify-center items-center w-full'>
        <Button
          id='UPLOAD_FILE_BUTTON'
          className='mt-4 h-fit'
          size={SIZE_TYPES.SMALL}
          type={BUTTON_TYPES.SECONDARY}
          isLoading={false}
        >
          Browse files
        </Button>
        <span className='f-12-400 rounded-2.5 text-GRAY_700  mt-1.5'>
          {footer ?? 'Or drag and drop .csv, .xslx files here'}
        </span>
        {!!error && (
          <span
            className={cn(
              error &&
                'absolute flex grow justify-center items-center w-[calc(100%+16px)] !text-RED_800 bg-white f-14-400 p-5 -left-2 top-[160px] rounded-2.5 border border-GRAY_400 shadow-tableFilterMenu',
              errorClassName,
            )}
          >
            {errorTitle}
          </span>
        )}
      </div>
    );
  };

  return (
    <>
      <div
        className={cn(
          'relative flex flex-col justify-center items-center bg-BG_GRAY_1 border border-dashed border-GRAY_400 min-h-[220px] rounded-md focus:border-black focus:border-solid cursor-pointer',
          isLoading ? 'cursor-not-allowed' : 'cursor-pointer',
          className,
          isDragOver && 'bg-GRAY_100 border-GRAY_500',
        )}
        onDrop={handleFileDrop}
        onClick={onClick}
        onKeyDown={handleKeyDown}
        onDragOver={(event) => event.preventDefault()}
        onDragEnter={() => setIsDragOver(true)}
        onDragExit={() => setIsDragOver(false)}
        onDragLeave={() => setIsDragOver(false)}
      >
        <div onDragOver={() => setIsDragOver(true)} className='flex flex-col items-center text-center w-full'>
          {renderContent()}
        </div>
        {children}
      </div>
    </>
  );
};

export default FileUploader;
