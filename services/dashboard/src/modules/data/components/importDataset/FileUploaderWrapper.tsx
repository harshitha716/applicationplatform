import React, { FC, useEffect, useRef, useState } from 'react';
import { REQUEST_TYPES } from 'apis/apiEndpoint.constants';
import { useGetPreviewTransformationMutation, useGetSignedUrlMutation } from 'apis/dataset';
import { useAppSelector } from 'hooks/toolkit';
import {
  FILE_IMPORT_STATUS_MSG,
  FILE_SIZE,
  FileExtensionToInputFormatMapping,
  FileExtensionToTypeMap,
  FileMimeType,
} from 'modules/data/components/importDataset/importData.constants';
import { FileUploaderWrapperPropsType } from 'modules/data/components/importDataset/importData.types';
import { getFileType } from 'modules/data/components/importDataset/importData.utils';
import { useRouter } from 'next/router';
import { RootState } from 'store';
import { toast } from 'components/common/toast/Toast';

const FileUploaderWrapper: FC<FileUploaderWrapperPropsType> = ({
  acceptedFormats,
  fileName,
  setFileName,
  isFileUploading = false,
  maxSize = FILE_SIZE.TWO_MB,
  filesSelected,
  Component,
  className,
  onClosePopup,
  setStartPollingPreview,
  onFileSelect,
  setRawData,
}) => {
  const router = useRouter();
  const isDatasetRoute = !!router?.query?.id;
  const hiddenFileInput = useRef<HTMLInputElement>(null);
  const user_id = useAppSelector((state: RootState) => state?.user)?.user?.user_id;
  const [getSignedUrl] = useGetSignedUrlMutation();
  const [getPreviewTransformation, { isSuccess }] = useGetPreviewTransformationMutation();
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [fileUploaderKey, setFileUploaderKey] = useState<number>(0);

  const handleChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const filesToUpload: File | null = event?.target?.files?.[0] ?? null;

    setFileName(filesToUpload?.name ?? null);
    handleUpload(filesToUpload);
  };

  const handleUpload = (filesToUpload: File | null) => {
    if (filesToUpload) {
      if (filesToUpload?.size > maxSize) {
        const err = `${FILE_IMPORT_STATUS_MSG.FILE_SIZE_EXCEEDED} ${maxSize / FILE_SIZE.ONE_MB}MB`;

        setError(err);
      } else {
        const fileExtension: string = filesToUpload?.name?.split('.')?.pop()?.toLowerCase() ?? '';
        const fileName = user_id + '_' + Date.now() + '.' + (FileMimeType[getFileType(filesToUpload)] ?? fileExtension);
        const fileType = FileMimeType[getFileType(filesToUpload)] ?? fileExtension;

        const signedUrlPayload = {
          file_name: fileName,
          file_type: fileType,
        };

        const acceptedFormatsArr = acceptedFormats.split(',').map((item) => item.trim());
        const isAllowedFormat = acceptedFormatsArr.includes(FileExtensionToInputFormatMapping[fileExtension]);

        if (isAllowedFormat) {
          setError(null);
          setIsLoading(true);
          onFileSelect(null);

          getSignedUrl(signedUrlPayload)
            .unwrap()
            .then(async (data: any) => {
              const upload_url = data?.upload_url;
              const file_upload_id = data?.file_upload_id;

              if (upload_url) {
                const xhr = new XMLHttpRequest();

                xhr.open(REQUEST_TYPES.PUT, upload_url, true);
                xhr.setRequestHeader(
                  'Content-Type',
                  getFileType(filesToUpload) || FileExtensionToTypeMap[fileExtension],
                );

                xhr.onload = function () {
                  if (xhr.status === 200) {
                    onFileSelect({
                      rawFile: filesToUpload,
                      identifier: file_upload_id,
                      url: upload_url,
                      fileName: filesToUpload.name,
                      downloadableUrl: upload_url,
                    });
                    triggerPreviewTransformation(file_upload_id);
                  } else {
                    setIsLoading(false);
                    setError(FILE_IMPORT_STATUS_MSG.FILE_UPLOAD_COMMON_ERROR);
                  }
                };

                xhr.onerror = function () {
                  setIsLoading(false);
                  setError(FILE_IMPORT_STATUS_MSG.FILE_UPLOAD_COMMON_ERROR);
                };

                xhr.send(filesToUpload);
              }
            })
            .catch(() => {
              setIsLoading(false);
              setRawData(null);
              setError(FILE_IMPORT_STATUS_MSG.FILE_UPLOAD_COMMON_ERROR);
            });
        } else {
          const err = `${FILE_IMPORT_STATUS_MSG.FILE_TYPE_INVALID} ${acceptedFormats}`;

          setRawData(null);
          setError(err);
        }
      }
    }
  };

  const triggerPreviewTransformation = async (file_upload_id: string) => {
    const oldDatasetImportPayload = {
      file_upload_id,
      dataset_id: router?.query?.id as string,
    };
    const newDatasetImportPayload = { file_upload_id };
    const payload = isDatasetRoute ? oldDatasetImportPayload : newDatasetImportPayload;

    getPreviewTransformation(payload)
      .unwrap()
      .then((data) => {
        if (data?.dataset_action_id) {
          setIsLoading(false);
          setStartPollingPreview({ check: true, actionId: data?.dataset_action_id, fileUploadId: file_upload_id });

          setTimeout(() => {
            onClosePopup();
          }, 1500);
        }
      })
      .catch(() => {
        setIsLoading(false);
        setRawData(null);
        toast.error(FILE_IMPORT_STATUS_MSG.PREVIEW_DATA_FAILED);
      });
  };

  const handleClick = () => {
    hiddenFileInput?.current?.click();
  };

  useEffect(() => {
    if (error) {
      setFileUploaderKey((prev) => prev + 1);
    }
  }, [error]);

  return (
    <Component
      onClick={handleClick}
      isLoading={isLoading || isFileUploading}
      isUploading={isLoading || isFileUploading}
      error={error}
      onFileDrop={handleUpload}
      filesSelected={filesSelected}
      supportedFile={acceptedFormats}
      className={className}
      fileName={fileName}
      setFileName={setFileName}
      isSuccess={isSuccess}
      handleChange={handleChange}
      indexKey={fileUploaderKey}
    >
      {!filesSelected && (
        <input
          type='file'
          ref={hiddenFileInput}
          onChange={handleChange}
          style={{ display: 'none' }}
          accept={acceptedFormats}
        />
      )}
    </Component>
  );
};

export default FileUploaderWrapper;
