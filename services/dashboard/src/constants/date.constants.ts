import { format } from 'date-fns';
import { OptionsType } from 'types/commonTypes';

export const MONTH_NAME = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

export const TIMEZONES = {
  EST: 'EST',
  UTC: 'UTC',
};

export const DATE_FORMATS = {
  YYYYMMDD: 'yyyy-MM-dd',
  DOLLLLYYYY: 'do LLLL, yyyy',
  ddMMMyyyy: 'dd MMM yyyy',
  YYYYMMDD_HHMMSS: 'yyyy-MM-dd HH:mm:ss',
  HHMM: 'HH:mm',
  MMddyyyy: 'MM/dd/yyyy',
  dd_MMM_yyyy: 'dd MMM yyyy',
  d_MMM_yyyy: 'd MMM yyyy',
  MMM_yyyy: 'MMM yyyy',
  YYYY: 'yyyy',
  DD: 'dd',
  QQ_yyyy: `'Q'Q yyyy`,
  EEE: 'EEE',
  ddMMyyyyHHmmss: 'dd/MM/yyyy, HH:mm:ss',
};

export const VALID_DATE_FORMATS = Object.values(DATE_FORMATS);

export const MONTHS_FULL = [
  'January',
  'February',
  'March',
  'April',
  'May',
  'June',
  'July',
  'August',
  'September',
  'October',
  'November',
  'December',
];

export const formatDateToZFormat = (date: Date) => {
  return `${format(date, 'yyyy-MM-dd')}T${format(date, 'HH:mm:ss.SSS')}Z`;
};

export const getUtcDate = (data: string | number) => {
  const inputDate = new Date(data);

  // Manually format using the extracted UTC parts with date-fns
  return new Date(
    inputDate.getUTCFullYear(),
    inputDate.getUTCMonth(),
    inputDate.getUTCDate(),
    inputDate.getUTCHours(),
    inputDate.getUTCMinutes(),
    inputDate.getUTCSeconds(),
  );
};

export enum DATE_FILTER_CATEGORIES {
  ALL_TIME = 'ALL_TIME',
  TODAY = 'TODAY',
  CUSTOM_DATE_RANGE = 'CUSTOM_DATE_RANGE',
  LAST_30_DAYS = 'LAST_30_DAYS',
  THIS_MONTH = 'THIS_MONTH',
  LAST_MONTH = 'LAST_MONTH',
  THIS_QUARTER = 'THIS_QUARTER',
  LAST_QUARTER = 'LAST_QUARTER',
  THIS_YEAR = 'THIS_YEAR',
  LAST_YEAR = 'LAST_YEAR',
}

export const DATE_FILTER_OPTIONS: OptionsType[] = [
  { label: 'All Time', value: DATE_FILTER_CATEGORIES.ALL_TIME },
  { label: 'Custom Date Range', value: DATE_FILTER_CATEGORIES.CUSTOM_DATE_RANGE },
  { label: 'Last 30 Days', value: DATE_FILTER_CATEGORIES.LAST_30_DAYS },
  { label: 'This Month', value: DATE_FILTER_CATEGORIES.THIS_MONTH },
  { label: 'Last Month', value: DATE_FILTER_CATEGORIES.LAST_MONTH },
  { label: 'This Quarter', value: DATE_FILTER_CATEGORIES.THIS_QUARTER },
  { label: 'Last Quarter', value: DATE_FILTER_CATEGORIES.LAST_QUARTER },
  { label: 'This Year', value: DATE_FILTER_CATEGORIES.THIS_YEAR },
  { label: 'Last Year', value: DATE_FILTER_CATEGORIES.LAST_YEAR },
];

export type dateFilterValueType = { start: Date | undefined; end: Date | undefined; periodicity?: PERIODICITY_TYPES };
// export type dateFilterValueType = { start: Date | null; end: Date | null };

export type RangeType = {
  startDate?: Date | undefined;
  endDate?: Date | undefined;
  color?: string | undefined;
  key?: string | undefined;
  autoFocus?: boolean | undefined;
  disabled?: boolean | undefined;
  showDateDisplay?: boolean | undefined;
};

export interface DateRangeValue {
  type: DATE_RANGE_TYPES;
  value: number | Date;
  year: number;
  label: string;
}

export enum DATE_RANGE_TYPES {
  DAY = 'day',
  MONTH = 'month',
  QUARTER = 'quarter',
  YEAR = 'year',
}

export const DATE_RANGE_TABS = [
  { label: 'Day', value: DATE_RANGE_TYPES.DAY },
  { label: 'Month', value: DATE_RANGE_TYPES.MONTH },
  { label: 'Quarter', value: DATE_RANGE_TYPES.QUARTER },
  { label: 'Year', value: DATE_RANGE_TYPES.YEAR },
];

export const MonthsConfig = [
  { label: 'January', value: 0, short: 'Jan' },
  { label: 'February', value: 1, short: 'Feb' },
  { label: 'March', value: 2, short: 'Mar' },
  { label: 'April', value: 3, short: 'Apr' },
  { label: 'May', value: 4, short: 'May' },
  { label: 'June', value: 5, short: 'Jun' },
  { label: 'July', value: 6, short: 'Jul' },
  { label: 'August', value: 7, short: 'Aug' },
  { label: 'September', value: 8, short: 'Sep' },
  { label: 'October', value: 9, short: 'Oct' },
  { label: 'November', value: 10, short: 'Nov' },
  { label: 'December', value: 11, short: 'Dec' },
];

export const DaysConfig = [
  { short: 'Sun', value: 0, label: 'Sunday' },
  { short: 'Mon', value: 1, label: 'Monday' },
  { short: 'Tue', value: 2, label: 'Tuesday' },
  { short: 'Wed', value: 3, label: 'Wednesday' },
  { short: 'Thu', value: 4, label: 'Thursday' },
  { short: 'Fri', value: 5, label: 'Friday' },
  { short: 'Sat', value: 6, label: 'Saturday' },
];

export const QuartersConfig = [
  { short: 'Q1', value: 0 },
  { short: 'Q2', value: 1 },
  { short: 'Q3', value: 2 },
  { short: 'Q4', value: 3 },
];

export enum DateRangeKeys {
  START_DATE = 'startDate',
  END_DATE = 'endDate',
}

export enum PERIODICITY_TYPES {
  DAILY = 'day',
  WEEKLY = 'week',
  MONTHLY = 'month',
  QUARTERLY = 'quarter',
  YEARLY = 'year',
}

export const PERIODICITY_OPTIONS: OptionsType[] = [
  { label: 'Daily', value: PERIODICITY_TYPES.DAILY },
  { label: 'Weekly', value: PERIODICITY_TYPES.WEEKLY },
  { label: 'Monthly', value: PERIODICITY_TYPES.MONTHLY },
  { label: 'Quarterly', value: PERIODICITY_TYPES.QUARTERLY },
  { label: 'Yearly', value: PERIODICITY_TYPES.YEARLY },
];
