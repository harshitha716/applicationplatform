import { defaultFnType } from 'types/commonTypes';

export type TopBarPropsType = {
  isSidebarOpen: boolean;
  onSidebarToggle: defaultFnType;
};

export const enum SHARE_BTN_ALLOWED_ROUTES {
  PAGES = '/pages/',
  DATASETS = '/datasets/',
  DATASET = '/datasets',
}
