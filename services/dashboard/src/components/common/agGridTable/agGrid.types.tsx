import { MapAny } from 'types/commonTypes';

export enum DataAlign {
  LEFT,
  CENTER,
  RIGHT,
}

export enum WorkspaceColumnDisplayTypes {
  COLOURED_NUMBERS = 'coloured_numbers',
  COMPARISON = 'comparison',
  TAGS = 'tags',
  CHANGE_RATE = 'change_rate',
  COLOURED_NUMBERS_WITH_CURRENCY = 'coloured_numbers_with_currency',
  CURRENCY_AMOUNT = 'currency_amount',
  DATE_BY_PERIODICITY = 'date_by_periodicity',
  EDITABLE_AMOUNT = 'editable_amount',
  EDITABLE_DATE = 'editable_date',
  FILE_DOWNLOAD = 'file_download',
  FORMATTED_DATE = 'formatted_date',
  MONTH_FROM_TIMESTAMP = 'month_from_timestamp',
  RISK_GRADE = 'risk_grade',
  STATUS_BY_COUNT = 'status_by_count',
  PROCESS_COUNT = 'process_count',
  STATUS_CHIP = 'status_chip',
  USER_NAME = 'user_name',
  DATE_RANGE = 'date_range',
  SNAKE_TO_PROPER_FORMATTING = 'snake_to_proper_formatting',
  ACTION_COMMENTS = 'action_comments',
  SYNC_STATUS = 'sync_status',
}

export interface WorkspaceDisplayColumns {
  column_name: string;
  display_name: string;
  is_hidden: boolean;
  display_type?: WorkspaceColumnDisplayTypes;
  metadata?: MapAny;
  dynamic_visibility?: {
    type: string;
    key: string;
    value: string;
  }[];
  Component?: React.ElementType;
}
