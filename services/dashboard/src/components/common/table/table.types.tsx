export enum CUSTOM_COLUMNS_TYPE {
  TAG = 'tags',
}

export enum DISPLAY_OPTIONS {
  COLUMNS = 'columns',
  GROUP_BY = 'group_by',
  CURRENCY = 'currency',
}

export enum VALUE_FORMAT_TYPE {
  ROUND_OFF = 'round_off',
  DATE_TIME = 'date_time',
  PREFIX = 'prefix',
  COLUMN_PREFIX = 'column_prefix',
}

export type ColumnVisibility = {
  colId: string;
  isVisible: boolean;
};
