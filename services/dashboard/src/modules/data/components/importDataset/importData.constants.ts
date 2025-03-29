import { FILE_MIME, INPUT_FILE_FORMATS } from 'modules/data/components/importDataset/importData.types';
import { LOADER_STATUS } from 'modules/data/data.types';

export const IMPORT_ALLOWED_FILE_FORMATS = `${INPUT_FILE_FORMATS.CSV},${INPUT_FILE_FORMATS.XLSX},${INPUT_FILE_FORMATS.XLS}`;

export const FILE_SIZE = {
  ONE_MB: 1000000,
  TWO_MB: 2000000,
  THREE_MB: 3000000,
  TWENTY_MB: 20000000,
  HUNDRED_MB: 100000000,
  TWO_HUNDRED_MB: 200000000,
  THREE_HUNDRED_MB: 300000000,
  FIVE_HUNDRED_MB: 500000000,
  ONE_GB: 1000000000,
};

export const FileMimeType: Record<string, string> = {
  [FILE_MIME.APPLICATION_PDF]: 'pdf',
  [FILE_MIME.IMAGE_JPEG]: 'jpeg',
  [FILE_MIME.IMAGE_PNG]: 'png',
  [FILE_MIME.IMAGE_BMP]: 'bmp',
  [FILE_MIME.TEXT_CSV]: 'csv',
};

export const FileExtensionToTypeMap: Record<string, string> = {
  pdf: FILE_MIME.APPLICATION_PDF,
  jpeg: FILE_MIME.IMAGE_JPEG,
  png: FILE_MIME.IMAGE_PNG,
  bmp: FILE_MIME.IMAGE_BMP,
  csv: FILE_MIME.TEXT_CSV,
};

export const FileExtensionToInputFormatMapping: Record<string, INPUT_FILE_FORMATS> = {
  pdf: INPUT_FILE_FORMATS.PDF,
  jpeg: INPUT_FILE_FORMATS.JPEG,
  png: INPUT_FILE_FORMATS.PNG,
  bmp: INPUT_FILE_FORMATS.BMP,
  csv: INPUT_FILE_FORMATS.CSV,
  jpg: INPUT_FILE_FORMATS.JPG,
  xls: INPUT_FILE_FORMATS.XLS,
  xlsx: INPUT_FILE_FORMATS.XLSX,
  bai2: INPUT_FILE_FORMATS.BAI2,
};

export const AI_TRANSFORMATION_STATUS = {
  STATUS_LOADING: {
    status: LOADER_STATUS.LOADING,
    title: 'Aligning file with table schema',
    description: 'In progress',
  },
  STATUS_ERROR: {
    status: LOADER_STATUS.ERROR,
    title: 'Error',
    description: 'Import failed due to schema mismatch',
  },
  STATUS_SUCCESS: {
    status: LOADER_STATUS.SUCCESS,
    title: 'File transformed successfully',
    description: 'Done',
  },
  STATUS_INGESTION_ONGOING: {
    status: LOADER_STATUS.LOADING,
    title: 'Data will be ingested in your dataset in few minutes',
    description: 'In progress',
  },
};

export enum FILE_IMPORT_STATUS_MSG {
  FILE_SIZE_EXCEEDED = 'File size cannot exceed more than ',
  FILE_TYPE_INVALID = 'Invalid file type, please upload ',
  FILE_FORMAT_NOT_SUPPORTED = 'File format not supported',
  FILE_UPLOAD_FAILED = 'Uploading failed!',
  FILE_UPLOAD_COMMON_ERROR = 'Your file couldn’t be processed, you can try again with a new file.',
  FILE_IMPORT_DATA_FAILED = 'Failed to import dataset',
  PREVIEW_DATA_FAILED = 'Error generating preview data',
  FILE_IMPORT_AFTER_AI_SUCCESS = 'File imported successfully',
}

export const enum DATA_PREVIEW_TABS_TYPES {
  FORMATTED = 'formatted',
  ORIGINAL = 'original',
}

export const dataPreviewTabItemList = [
  {
    value: DATA_PREVIEW_TABS_TYPES.FORMATTED,
    label: 'Formatted',
  },
  {
    value: DATA_PREVIEW_TABS_TYPES.ORIGINAL,
    label: 'Original',
  },
];
