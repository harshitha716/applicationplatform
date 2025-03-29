import { PAGE_ACCESS_PRIVILEGES } from 'modules/page/pages.types';

export const PAGE_ACCESS_PRIVILEGES_LIST = [
  {
    label: 'Admin',
    value: PAGE_ACCESS_PRIVILEGES.ADMIN,
  },
  {
    label: 'Viewer',
    value: PAGE_ACCESS_PRIVILEGES.VIEWER,
  },
];

export const CHANGE_PAGE_ACCESS_PRIVILEGES_LIST = [
  {
    label: 'Admin',
    value: PAGE_ACCESS_PRIVILEGES.ADMIN,
    desc: 'Can manage and share dataset',
  },
  {
    label: 'Viewer',
    value: PAGE_ACCESS_PRIVILEGES.VIEWER,
    desc: 'Can view data only',
  },
];

export const PAGE_CURRENCY_OPTIONS = [
  'local',
  'TRY',
  'ILS',
  'LBP',
  'NZD',
  'VND',
  'HUF',
  'PKR',
  'MXN',
  'BHD',
  'KZT',
  'QAR',
  'DKK',
  'PLN',
  'SGD',
  'EGP',
  'JOD',
  'USD',
  'GBP',
  'AZN',
  'GEL',
  'SEK',
  'INR',
  'ALL',
  'KWD',
  'CAD',
  'HKD',
  'OMR',
  'CHF',
  'AED',
  'NOK',
  'ISK',
  'BRL',
  'AUD',
  'RSD',
  'DZD',
  'SAR',
  'JPY',
  'MAD',
  'CZK',
  'IQD',
  'EUR',
  'UZS',
];

export const CURRENCY_SYMBOLS = {
  USD: '$',
};
