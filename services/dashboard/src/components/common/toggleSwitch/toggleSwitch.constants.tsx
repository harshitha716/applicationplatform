import { TOGGLE_SWITCH_STATE_TYPES, TOGGLE_SWITCH_TYPES } from 'components/common/toggleSwitch/toggleSwitch.types';

export const TOGGLE_SWITCH_STYLES = {
  [TOGGLE_SWITCH_TYPES.SELECTED]: {
    [TOGGLE_SWITCH_STATE_TYPES.ENABLED]: 'bg-GRAY_1000 box-content border-0.5 border-GRAY_1000',
    [TOGGLE_SWITCH_STATE_TYPES.HOVER]: '',
    [TOGGLE_SWITCH_STATE_TYPES.DISABLED]: 'bg-GRAY_100 cursor-not-allowed',
  },
  [TOGGLE_SWITCH_TYPES.UNSELECTED]: {
    [TOGGLE_SWITCH_STATE_TYPES.ENABLED]: 'bg-BORDER_GRAY_400 box-content border-BORDER_GRAY_400',
    [TOGGLE_SWITCH_STATE_TYPES.HOVER]: ' border-0.5 hover:border-GRAY_1000',
    [TOGGLE_SWITCH_STATE_TYPES.DISABLED]: 'bg-GRAY_100 cursor-not-allowed',
  },
};

export const TOGGLE_SWITCH_SLIDER = {
  [TOGGLE_SWITCH_TYPES.SELECTED]: {
    [TOGGLE_SWITCH_STATE_TYPES.ENABLED]: 'right-[3px] bg-white',
    [TOGGLE_SWITCH_STATE_TYPES.DISABLED]: 'bg-BORDER_GRAY_400 right-[3px]',
  },
  [TOGGLE_SWITCH_TYPES.UNSELECTED]: {
    [TOGGLE_SWITCH_STATE_TYPES.ENABLED]: 'left-[3px] bg-white',
    [TOGGLE_SWITCH_STATE_TYPES.DISABLED]: 'bg-BORDER_GRAY_400 left-[3px]',
  },
};
