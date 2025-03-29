import { defaultFnType, OptionsType } from 'types/commonTypes';

export type AsyncDropdownPropsType = {
  onOpen: defaultFnType;
  onClose: defaultFnType;
  isOpen: boolean;
  onDelete?: defaultFnType;
  onChange: (role: OptionsType) => void;
  options: OptionsType[];
  selectedValue: OptionsType;
  defaultValue: OptionsType;
  showDelete?: boolean;
  isHoveredDropdown?: boolean;
  setIsHoveredDropdown?: (isHoveredDropdown: boolean) => void;
  wrapperClassName?: string;
  parentWrapperClassName?: string;
  showSelectedIcon?: boolean;
  selectedOptionClassName?: string;
  isOverflowStyle?: boolean;
};
