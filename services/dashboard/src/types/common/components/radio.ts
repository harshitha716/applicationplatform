export enum RADIO_TYPES {
  SELECTED = 'SELECTED',
  UNSELECTED = 'UNSELECTED',
}

export enum RADIO_STATE_TYPES {
  ENABLED = 'ENABLED',
  HOVER = 'HOVER',
  PRESSED = 'PRESSED',
  DISABLED = 'DISABLED',
}

export interface RadioProps {
  wrapperClassName?: string;
  radioClassName?: string;
  radioSelectedClassName?: string;
  radioDefaultClassName?: string;
  wrapperStyle?: string;
  radioStyle?: string;
  radioSelectedStyle?: string;
  radioDefaultStyle?: string;
  isDisabled?: boolean;
  checked?: boolean;
  onSelect?: (value?: boolean) => void;
  isClearable?: boolean;
  id: string;
}
