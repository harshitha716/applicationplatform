import { NavigationItemSchema } from 'types/config';

export const ROUTES_PATH = {
  HOME: '/',
  LOGIN: '/login',
  DATA: '/datasets',
  TEAM: '/team',
  DRILLDOWN: '/drilldown/:datasetId/:rowId',
  DATASET: '/datasets/:datasetId',
  PAGES: '/pages/',
  NO_ACCESS: '/no-access',
  ADMIN: '/admin',
  PAYMENTS: '/payments',
};

export const getPageRouteById = (pageId: string) => {
  return `${ROUTES_PATH.PAGES}${pageId}`;
};

export const getDatasetRouteById = (datasetId: string) => {
  return `${ROUTES_PATH.DATA}/${datasetId}`;
};

export const LOGIN_URLS = [ROUTES_PATH.LOGIN];

export const SIDEBAR_ITEMS: NavigationItemSchema[] = [
  {
    label: 'Data',
    iconId: 'coins-stacked-04',
    path: ROUTES_PATH.DATA,
  },
  {
    label: 'Payments',
    iconId: 'send-01',
    path: ROUTES_PATH.PAYMENTS,
    isHidden: true,
  },
  {
    label: 'Team',
    iconId: 'users-02',
    path: ROUTES_PATH.TEAM,
  },
];
