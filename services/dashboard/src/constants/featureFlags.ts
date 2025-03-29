import { ENVIRONMENT } from 'constants/common.constants';

export const getClientSideId = (environment: string) => {
  switch (environment) {
    case 'staging':
      return '657145b8cc5a09104bbb584a';
    case 'production':
      return '6569c10a60139c0f61aa5cc2';
    default:
      return '657145d1b551a5101ba31c81';
  }
};

export const LAUNCH_DARKLY_CLIENT_SIDE_ID = getClientSideId(ENVIRONMENT);

export enum FEATURE_FLAGS {
  PEOPLE_MEMBERSHIP_REQUESTS = 'people-membership-requests',
  ADMIN_PAGE = 'admin-page',
}
