import { EventCallbackType } from 'types/commonTypes';

export enum TOGGLE_SWITCH_TYPES {
  SELECTED = 'SELECTED',
  UNSELECTED = 'UNSELECTED',
}

export enum TOGGLE_SWITCH_STATE_TYPES {
  ENABLED = 'ENABLED',
  HOVER = 'HOVER',
  PRESSED = 'PRESSED',
  DISABLED = 'DISABLED',
}

export interface ToggleSwitchProps {
  checked?: boolean;
  id: string;
  onChange: (state: boolean) => void;
  label?: string;
  disabled?: boolean;
  toggleClassName?: string;
  sliderClassName?: string;
  sliderStyle?: string;
  toggleStyle?: string;
  eventCallback?: EventCallbackType;
  wrapperClassName?: string;
  labelClassName?: string;
  controlled?: boolean;
}
