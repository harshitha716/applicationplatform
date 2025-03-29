import { TAB_TYPES } from 'types/common/components';

export const TAB_STYLES = {
  [TAB_TYPES.FILLED]: {
    tabItemSelectedClassName: 'text-GRAY_1000 bg-GRAY_100 transition-all transform',
    tabItemDefaultClassName: 'text-GRAY_700 bg-white',
    tabItemGapClassName: 'mr-3',
    tabItemClassName: 'py-1 px-2 rounded',
  },
  [TAB_TYPES.OUTLINE]: {
    tabItemSelectedClassName: 'text-GRAY_1000 border border-GRAY_400',
    tabItemDefaultClassName: 'text-GRAY_1000 border border-transparent',
    tabItemGapClassName: 'mr-1.5',
    tabItemClassName: 'py-2 px-3 rounded-lg',
  },
  [TAB_TYPES.FILLED_OUTLINED]: {
    tabItemSelectedClassName: 'text-GRAY_1000 border border-GRAY_400 bg-white transition-all transform',
    tabItemDefaultClassName: 'text-GRAY_900 bg-GRAY_100 border border-transparent',
    tabItemGapClassName: 'mr-0',
    tabItemClassName: 'py-1.5 px-6 rounded',
  },
  [TAB_TYPES.UNDERLINE]: {
    tabItemSelectedClassName: 'text-GRAY_1000 border-b-2 border-GRAY_1000',
    tabItemDefaultClassName: 'text-GRAY_700 border-b-2 border-transparent',
    tabItemGapClassName: 'mr-5',
    tabItemClassName: 'py-1.5 px-1 rounded-none',
  },
};
