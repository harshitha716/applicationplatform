import { MouseEventHandler } from 'react';
import clsx, { ClassValue } from 'clsx';
import { CHIP_COLORS } from 'constants/colors';
import { SCREEN_BREAKPOINTS } from 'constants/common.constants';
import { DATE_FILTER_CATEGORIES, DATE_FILTER_OPTIONS } from 'constants/date.constants';
import { format, startOfYear } from 'date-fns';
import { twMerge } from 'tailwind-merge';
import { DateFilterValueType } from 'components/filter/DateRangeFilter';

declare type MapAny = Record<string, any>;

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export const checkIsObjectEmpty = (obj?: MapAny, ignoreKeys?: string[]) => {
  if (!obj) return true;

  if (typeof obj === 'object' && !Object.keys(obj).length) return true;

  if (!ignoreKeys?.length && typeof obj === 'object' && !Object.keys(obj).length) return true;

  if (ignoreKeys?.length && typeof obj === 'object') {
    const keys = Object.keys(obj);

    if (JSON.stringify(keys.sort()) === JSON.stringify(ignoreKeys?.sort())) return true;
  }

  return false;
};

export const stopPropagationAction: MouseEventHandler<Element> = (event) => {
  event?.stopPropagation();
};

export const isArrayOrObject = (value: unknown): boolean => {
  return Array.isArray(value) || typeof value === 'object';
};

// eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/explicit-module-boundary-types
export function debounce<T extends (...args: any[]) => any>(func: T, wait: number) {
  let timeout: ReturnType<typeof setTimeout> | null;

  return function (this: ThisParameterType<T>, ...args: Parameters<T>) {
    if (timeout) clearTimeout(timeout);
    timeout = setTimeout(() => {
      timeout = null;
      func.apply(this, args);
    }, wait);
  };
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/explicit-module-boundary-types
export function doDebounce<T extends (...args: any[]) => any>(func: T, wait: number) {
  let timeout: ReturnType<typeof setTimeout> | null;

  return function (this: ThisParameterType<T>, ...args: Parameters<T>) {
    if (!timeout) {
      func.apply(this, args);
      timeout = setTimeout(() => {
        timeout = null;
      }, wait);
    }
  };
}

export const getStartOfYear = (year: number) => {
  return startOfYear(new Date(year, 0, 1)); // January 1st of the specified year
};

export const isUserInUS = function () {
  const userLocale = navigator.language;

  return userLocale.endsWith('-US');
};

export const getDateRangeTitle = (dateRangeFilter: DateFilterValueType, showSingleDate?: boolean): string => {
  const start = dateRangeFilter?.start_date;
  const end = dateRangeFilter?.end_date;

  if (dateRangeFilter?.date_category === DATE_FILTER_CATEGORIES.CUSTOM_DATE_RANGE && start && end) {
    if (showSingleDate && start?.toDateString() === end?.toDateString()) {
      return format(start, 'dd MMM yyyy');
    }

    return `${format(start, 'dd MMM yyyy')} - ${format(end, 'dd MMM yyyy')}`;
  }

  const dateRangeCategory =
    DATE_FILTER_OPTIONS.find((category) => category.value === dateRangeFilter?.date_category) ?? DATE_FILTER_OPTIONS[0];

  return `Date range - ${dateRangeCategory?.label}`;
};

/**
 * Inject the dynamic parameters in the url from a parameter object
 * @param url
 * @param params
 * @returns
 */
export const formRequestUrlWithParams = (url: string, params: MapAny) => {
  let formattedUrl = url;

  Object.keys(params).forEach((key) => {
    formattedUrl = formattedUrl.replace(`{{${key}}}`, params[key]);
  });

  return formattedUrl;
};

export function isCamelCase(str: string) {
  // Regular expression to match camelCase
  const camelCaseRegex = /^[a-z]+([A-Z][a-z]*)*$/;

  return camelCaseRegex.test(str);
}

export function camelCaseToNormalText(camelCaseStr: string) {
  const isCamelCaseString = isCamelCase(camelCaseStr);

  if (!isCamelCaseString) return camelCaseStr;

  return camelCaseStr
    ?.replace(/([A-Z])/g, ' $1') // Insert a space before uppercase letters
    ?.replace(/^./, (str) => str.toUpperCase()); // Capitalize the first letter
}

/**
 * Format the number to a comma separated number
 * @param num 1000000
 * @returns 1,000,000
 */
export const getCommaSeparatedNumber = (num?: number, precision = 0) => {
  return num === undefined || num === null
    ? '-'
    : num.toLocaleString('en-US', {
        maximumFractionDigits: precision,
        minimumFractionDigits: precision,
      });
};

export const capitalizeFirstLetter = (str: string) => {
  if (!str) return '';

  return str.charAt(0).toUpperCase() + str.slice(1);
};

export const getColorValue = () => Math.floor(Math.random() * 64) + 80;

export const getRandomColor = () => `rgb(${getColorValue()}, ${getColorValue()}, ${getColorValue()}`;

export const getFirstLetters = (str: string) =>
  str
    ?.split(' ')
    .map((word, index) => {
      if (index > 1 || !word.length) return null;
      else return word[0].toUpperCase();
    })
    .join('');

export const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

/**
 * Check if the email is valid
 * @param email string admin@zamp.ai
 * @returns boolean true
 */
export const isValidEmail = (email: string) => {
  return emailRegex.test(email);
};

/**
 * Get the domain from the email
 * @param email string admin@zamp.ai
 * @returns string zamp.ai
 */
export const getDomainFromEmail = (email: string) => {
  return email.split('@')[1];
};

/**
 * Get the username from the email
 * @param email string admin@zamp.ai
 * @returns string admin
 */
export const getUserNameFromEmail = (email: string) => {
  if (!email) return '';

  return email.split('@')[0];
};

/**
 * Convert the email username to name
 * @param emailUsername string admin.zamp
 * @returns string Admin Zamp
 */
export const convertEmailUsernameToName = (emailUsername: string) => {
  return emailUsername
    .split('.')
    .map((name) => capitalizeFirstLetter(name))
    .join(' ');
};

export function isValidDate(dateString: string) {
  // Try to parse the string into a Date object
  const date = new Date(dateString);

  // Check if the date is invalid or not
  return !isNaN(date.getTime());
}

/**
 * Shuffle the array position of elements
 * @param array
 * @returns
 */
export function shuffleArray(array: any[]) {
  for (let i = array.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1)); // Random index between 0 and i

    [array[i], array[j]] = [array[j], array[i]]; // Swap elements
  }

  return array;
}

/**
 * Format the number to a short format
 * @param number
 * @returns string number with suffix
 */
export function formatNumber(
  value = 0,
  precision: number = 1,
  allowSuffix: boolean = true,
  getSuffix: boolean = false,
): string {
  const suffixes = [
    { threshold: 1000000000, suffix: 'B', valueString: 'Billions' },
    { threshold: 1000000, suffix: 'M', valueString: 'Millions' },
    { threshold: 1000, suffix: 'K', valueString: 'Thousands' },
  ];

  if (getSuffix) {
    for (const { threshold, valueString } of suffixes) {
      if (value >= threshold) {
        return valueString;
      }
    }
  }

  for (const { threshold, suffix } of suffixes) {
    if (value >= threshold) {
      return (value / threshold).toFixed(precision).replace(/\.00$/, '') + (allowSuffix ? suffix : '');
    }
  }

  return value.toString();
}

/**
 * Get all data values for the given keys
 * @param data
 * @param keys
 * @returns return the max value from the data for the given keys
 */
export const getMaxValue = (data: MapAny[], keys: string[]) => {
  const maxValue = Math.max(...data.flatMap((item) => keys.map((key) => item[key] || 0)));

  return maxValue;
};

/**
 * Calling cyclicIterator() returns a new iterator instance each time.
 * The iterator function remembers its current index using a closure.
 * When it reaches the last element, it loops back to the start.
 * @param arr
 * @returns a new iterator instance each time
 */
export const cyclicIterator = (arr: any[]) => {
  let index = 0;

  return () => {
    const value = arr[index];

    index = (index + 1) % arr.length;

    return value;
  };
};

/**
 * Get a color from the predefined list of colors using a cyclic iterator
 * @param colorArray
 * @returns a new iterator instance each time
 */
export const getChipColor = (colorArray: string[]) => cyclicIterator(colorArray);

/**
 * Get a color from CHIP_COLORS using a cyclic iterator
 * @returns a new iterator instance each time
 */
export const getTagColor: () => string = cyclicIterator(CHIP_COLORS);

/**
 * Get the leading path after the '/' from the current URL.
 * @param path https://zamp.ai/datasets/12345678
 * @returns {string} /datasets
 */
export const getLeadingPathFromURL = (path: string) => {
  if (!path || path === '/') return '/';

  const pathSegments = path.split('/').filter(Boolean);

  return `/${pathSegments[0]}`;
};

export function trimString(str: string, maxLength: number) {
  if (str.length > maxLength) {
    return str.slice(0, maxLength - 3) + '...'; // Slice the string to maxLength minus 3 for the ellipsis
  }

  return str;
}

/**
 * Copy the text to the clipboard
 * @param copyText string
 */
export const copyToClipBoard = (copyText: string) => {
  if ('clipboard' in navigator) {
    navigator.clipboard.writeText(copyText);
  } else {
    const textField = document.createElement('textarea');

    textField.innerText = copyText;
    document.body.appendChild(textField);
    textField.select();
    document.execCommand('copy');
    textField.remove();
  }
};

export const extractFileNameFromUrl = (url: string): string => {
  return url.split('/').pop()?.split('?')[0] ?? '';
};

export const fetchFileBlob = async (url: string) => {
  const file = await fetch(url)
    .then((response) => response.blob())
    .then((data) => data);

  return file;
};

/**
 * @param {string | File} file
 * @param {Function} setIsLoading
 * @description Download File or from url
 */
export const downloadFile = async (
  file: string | File | null,
  setIsLoading?: (flag: boolean) => void,
  overrideFileName?: string,
) => {
  try {
    if (!file) return;

    let fileBlob: Blob;
    let fileName: string;

    if (typeof file === 'string') {
      fileName = extractFileNameFromUrl(file);
      fileBlob = await fetchFileBlob(file).then((data) => data);
    } else {
      fileName = file.name;
      fileBlob = file;
    }

    const aTag = document.createElement('a');
    const objUrl = URL.createObjectURL(fileBlob);

    aTag.setAttribute('href', objUrl);
    if (fileName) aTag.setAttribute('download', overrideFileName ?? fileName);
    document.body.appendChild(aTag);
    aTag.click();
    aTag.remove();
    URL.revokeObjectURL(objUrl);
  } finally {
    setIsLoading?.(false);
  }
};

export const getPastDateByNumberOfDays = (numberOfDays: number) => {
  const pastDate = new Date();

  pastDate.setDate(pastDate.getDate() - numberOfDays);

  return pastDate;
};

/**
 * Formats a given value as a currency string.
 *
 * @param value - The value to be formatted. If the value is not a number, it will be treated as 0.
 * @returns The formatted currency string with a dollar sign, commas as thousand separators, and two decimal places.
 */
export const formatCurrencyValue = (value: any): string => {
  const numValue = isNaN(Number(value)) ? 0 : Number(value);

  return `$${numValue.toFixed(2).replace(/\B(?=(\d{3})+(?!\d))/g, ',')}`;
};

/**
 *
 * @param text e.g. Abcd Efgh
 * @returns sentence casew e.g. Abcd efgh
 */
export const sentenceCase = (str: string) => {
  if (str === undefined || str === null || str === '') return '';

  // Convert the input to a string
  str = str.toString();

  // Capitalize the first letter and leave the rest unchanged
  return str.charAt(0).toUpperCase() + str.slice(1);
};

/**
 * Converts a snake_case string to Sentence case.
 *
 * @param str - The snake_case string to be converted.
 * @returns The converted string in Sentence case.
 */
export const snakeCaseToSentenceCase = (str: string) => {
  return sentenceCase(str?.split('_').join(' '));
};

/**
 * Creates a date object from a UTC string
 * @param date
 * @returns date object
 */
export const createDateObjectFromUTCString = (date: string | Date) => {
  if (date instanceof Date) return date;

  const dateParts = date?.split('T');
  const timeParts = dateParts?.[1]?.split(':');
  const dateObject = new Date(dateParts?.[0]);

  const adjustedDate = new Date(
    Date.UTC(dateObject.getUTCFullYear(), dateObject.getUTCMonth(), dateObject.getUTCDate()),
  );

  adjustedDate.setHours(Number(timeParts?.[0]), Number(timeParts?.[1]), 0, 0);

  return adjustedDate;
};

/**
 * Formats a count with the appropriate singular or plural form.
 *
 * @param count - The number according to which the plural or singular form is to be used.
 * @param word - The word to be used if the count is 1.
 * @param pluralWord - The plural word to be used if the count is greater than 1 else the default plural word will be `${word}s`.
 * @returns The formatted string with the appropriate singular or plural form.
 */
export const formatPlural = (count: number, word: string, pluralWord?: string) => {
  return `${count} ${count > 1 ? (pluralWord ?? `${word}s`) : word}`;
};

export const checkScreenBreakpoint = (width: number, height: number) =>
  width && height ? width < SCREEN_BREAKPOINTS.MIN_WIDTH || height < SCREEN_BREAKPOINTS.MIN_HEIGHT : false;

/**
 * Validates email
 * @param email
 * @returns boolean
 */
export const validateEmail = (email: string): boolean => {
  const testEmailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

  return testEmailRegex.test(email);
};

/*
 * Check if the value is of type object or array or both
 * @param value
 * @param type
 * @returns boolean
 */
export const checkObjOrArrType = (value: unknown, type: 'object' | 'array' | 'both'): boolean => {
  if (type === 'object') return typeof value === 'object' && value !== null && !Array.isArray(value);
  if (type === 'array') return Array.isArray(value);

  return typeof value === 'object' && value !== null; // 'both' case (checks for object or array)
};
