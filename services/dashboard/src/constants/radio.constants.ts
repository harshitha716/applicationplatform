import { RADIO_STATE_TYPES, RADIO_TYPES } from 'types/common/components/radio';

export const RADIO_STATE_STYLES = {
  [RADIO_TYPES.SELECTED]: {
    [RADIO_STATE_TYPES.ENABLED]: 'cursor-pointer border-TEXT_PRIMARY',
    [RADIO_STATE_TYPES.HOVER]: 'hover:bg-white hover:ring-[14px] hover:ring-LIGHT_PRIMARY_1',
    [RADIO_STATE_TYPES.PRESSED]: 'active:bg-TEXT_PRIMARY hover:ring-LIGHT_PRIMARY_2',
    [RADIO_STATE_TYPES.DISABLED]: 'cursor-not-allowed border-DIVIDER_GRAY bg-white',
  },
  [RADIO_TYPES.UNSELECTED]: {
    [RADIO_STATE_TYPES.ENABLED]: 'cursor-pointer border-DIVIDER_GRAY bg-white',
    [RADIO_STATE_TYPES.HOVER]: 'hover:bg-white hover:border-TEXT_PRIMARY hover:ring-[14px] hover:ring-LIGHT_PRIMARY_1',
    [RADIO_STATE_TYPES.PRESSED]: 'active:bg-TEXT_PRIMARY hover:ring-LIGHT_PRIMARY_2',
    [RADIO_STATE_TYPES.DISABLED]: 'cursor-not-allowed border-DIVIDER_GRAY bg-BASE_PRIMARY',
  },
};
