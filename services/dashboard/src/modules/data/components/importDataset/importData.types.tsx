import { Dispatch, SetStateAction } from 'react';
import { LOADER_STATUS } from 'modules/data/data.types';
import { DatasetBulkLoadersType } from 'store/slices/user';
import { RawMetadata, TransformationPreviewMetadata, UploadFileResponseType } from 'types/api/dataset.types';
import { defaultFnType } from 'types/commonTypes';

export enum INPUT_FILE_FORMATS {
  PNG = 'image/png',
  JPEG = 'image/jpeg',
  JPG = 'image/jpg',
  DOCX = '.docx',
  PDF = '.pdf',
  BMP = '.bmp',
  CSV = '.csv',
  XLSX = '.xlsx',
  XLS = '.xls',
  BAI2 = '.bai2',
}

export enum FILE_EXTENSION {
  CSV = 'csv',
  XLSX = 'xlsx',
}

export enum STATUS_TYPES {
  PENDING = 'pending',
  SUCCESS = 'success',
  FAILED = 'failed',
}

export enum FILE_MIME {
  APPLICATION_PDF = 'application/pdf',
  IMAGE_JPEG = 'image/jpeg',
  IMAGE_PNG = 'image/png',
  IMAGE_BMP = 'image/bmp',
  TEXT_CSV = 'text/csv',
  XLSX = 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
  XLS = 'application/vnd.ms-excel',
  BAI2 = 'application/x-bai',
}

export type HistoryListPropsType = {
  fileImportHistoryData: FileHistoryDataType[];
  isHoveredLoaders: boolean;
};

export type HistoryBulkLoadersPropsType = {
  isHoveredLoaders: boolean;
  datasetBulkLoaders: DatasetBulkLoadersType[];
  setIsHoveredLoaders: (isHoveredLoaders: boolean) => void;
};

export type FileHistoryDataType = {
  id: string;
  dataset_id: string;
  file_id: string;
  file_name: string;
  file_upload_status: string;
  status?: LOADER_STATUS;
  file_upload_created_at: string;
  uploaded_by_user: {
    email: string;
    name?: string;
  };
};

export type ImportFileHistoryPropsType = {
  onClose: defaultFnType;
};

export type FileUploaderSideDrawerPropsType = {
  onClose: defaultFnType;
};

export type TableSchemaAlignmentStatusPropsType = {
  showAiTransformationStatus: { open: boolean; status: string; title: string; description: string };
  setShowAiTransformationStatus: Dispatch<
    SetStateAction<{ open: boolean; status: string; title: string; description: string }>
  >;
};

export type ImportDatasetPropsType = {
  setShowAiTransformationStatus: Dispatch<
    SetStateAction<{ open: boolean; status: string; title: string; description: string }>
  >;
  onRefetch: defaultFnType;
};

export type FileUploaderPropsType = {
  isLoading?: boolean;
  isSuccess?: boolean;
  error: string | null;
  onFileDrop: (file: File | null) => void;
  onClick: defaultFnType;
  children?: React.ReactNode;
  errorClassName?: string;
  footer?: string;
  className?: string;
  fileName: string | null;
  setFileName: (fileName: string | null) => void;
  indexKey?: number;
};

export type FileUploaderWrapperPropsType = {
  acceptedFormats: string;
  filesSelected?: string;
  isFileUploading?: boolean;
  maxSize?: number;
  Component: React.ElementType;
  className?: string;
  onClosePopup: defaultFnType;
  setRawData: Dispatch<SetStateAction<RawMetadata | null>>;
  setStartPollingPreview: Dispatch<SetStateAction<{ check: boolean; actionId: string; fileUploadId: string }>>;
  fileName: string | null;
  setFileName: (fileName: string | null) => void;
  onFileSelect: (arg: UploadFileResponseType | null) => void;
};

export type ImportFilePropsType = {
  acceptedFormats: string;
  filesSelected?: string;
  className?: string;
  file?: UploadFileResponseType | null;
  onClosePopup: defaultFnType;
  setRawData: Dispatch<SetStateAction<RawMetadata | null>>;
  setStartPollingPreview: Dispatch<SetStateAction<{ check: boolean; actionId: string; fileUploadId: string }>>;
  fileName: string | null;
  setFileName: (fileName: string | null) => void;
};

export type ImportFileWrapperPropsType = {
  onReset: defaultFnType;
  isOpen: boolean;
  onClose: defaultFnType;
  setStartPollingPreview: Dispatch<SetStateAction<{ check: boolean; actionId: string; fileUploadId: string }>>;
  setRawData: Dispatch<SetStateAction<RawMetadata | null>>;
  fileName: string | null;
  setFileName: (fileName: string | null) => void;
};

export type ImportedDataPreviewPropsType = {
  onReset: defaultFnType;
  rawData: RawMetadata | null;
  mappedData: TransformationPreviewMetadata | null;
  startAiTransformation: boolean;
  fileUploadId: string;
  fileName: string | null;
  onRefetch: defaultFnType;
  setShowAiTransformationStatus: Dispatch<
    SetStateAction<{ open: boolean; status: string; title: string; description: string }>
  >;
};

export type DataPreviewContentPropsType = {
  rawData: RawMetadata | null;
  mappedData: TransformationPreviewMetadata | null;
};

export type DataPreviewSidebarPropsType = {
  fileName: string | null;
  onReset: defaultFnType;
  fileUploadId: string;
  onRefetch: defaultFnType;
  setShowAiTransformationStatus: Dispatch<
    SetStateAction<{ open: boolean; status: string; title: string; description: string }>
  >;
};

export type StartPollingPreviewType = {
  check: boolean;
  actionId: string;
  fileUploadId: string;
};
