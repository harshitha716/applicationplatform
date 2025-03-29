export const API_VERSION = process.env.NEXT_PUBLIC_API_VERSION;
export const APP_NAME = process.env.NEXT_PUBLIC_APP_NAME;
export const ENVIRONMENT = process.env.NEXT_PUBLIC_ENVIRONMENT;
export const AZURE_CLIENT_ID = process.env.NEXT_PUBLIC_AZURE_CLIENT_ID ?? '';
export const AZURE_AUTHORITY = process.env.NEXT_PUBLIC_AZURE_AUTHORITY ?? '';
export const AZURE_REDIRECT = process.env.NEXT_PUBLIC_AZURE_REDIRECT || '/';
export const AZURE_SCOPE = process.env.NEXT_PUBLIC_AZURE_SCOPE ?? '';
export const VELT_KEY = process.env.NEXT_PUBLIC_VELT_KEY ?? '';
export const CSRF_TOKEN_KEY = 'X-ZAMP-CSRF';
export const PLATFORM_HEADER_KEY = 'X-Platform';
export const CANARY_HEADER_KEY = 'X-Canary';
export const REQUEST_TIMEOUT = 40000;
export const ABORT_ERROR = 'AbortError: signal is aborted without reason';
export const PLATFORM_TMS = 'TMS';

export enum APITags {
  GET_USER = 'GET_USER',
  GET_PEOPLE_INVITATIONS = 'GET_PEOPLE_INVITATIONS',
  GET_PEOPLE_TEAM_MEMBERS = 'GET_PEOPLE_TEAM_MEMBERS',
  GET_ALL_TEAMS = 'GET_ALL_TEAMS',
  GET_DATASET_LISTING = 'GET_DATASET_LISTING',
}
export const API_TAGS = Object.values(APITags);

const getApiDomain = (environment: string = '') => {
  switch (environment) {
    case 'production':
      return 'https://api.zamp.ai';
    case 'staging':
      return 'https://api-stg.zamp.ai';
    case 'development':
      return 'https://api-dev.zamp.ai';
    default:
      return 'http://localhost:8080';
  }
};

export const ERROR_TOKENS = {
  INVALID_TOKEN: 'INVALID_TOKEN',
  FAILED_TO_INITIATE_LAUNCHDARKLY: 'FAILED_TO_INITIATE_LAUNCHDARKLY',
  MISSING_TOKEN: 'MISSING_TOKEN',
  PAGE_BREAK: 'PAGE_BREAK',
  PAGE_404: 'PAGE_404',
  NO_PERMISSIONS_PAGE: 'NO_PERMISSIONS_PAGE',
  USER_WITH_NO_PERMISSIONS: 'USER_WITH_NO_PERMISSIONS',
  CLIENT_INVALID_API_CALL: 'CLIENT_INVALID_API_CALL',
  CSV_PARSING_ERROR: 'CSV_PARSING_ERROR',
  CLIENT_SIDE_EXCEPTION: 'CLIENT_SIDE_EXCEPTION',
};

export const SESSION_EXPIRY_TOKENS = [ERROR_TOKENS.INVALID_TOKEN, ERROR_TOKENS.MISSING_TOKEN];

export const API_DOMAIN = getApiDomain(ENVIRONMENT);
