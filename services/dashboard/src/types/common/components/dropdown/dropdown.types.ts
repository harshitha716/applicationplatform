import React, { ComponentType, CSSProperties, ReactElement, ReactNode, RefCallback } from 'react';
import {
  ActionMeta,
  CSSObjectWithLabel,
  DropdownIndicatorProps,
  MenuListProps,
  MultiValue,
  MultiValueProps,
  OptionProps,
  SingleValue,
  SingleValueProps,
  ValueContainerProps,
} from 'react-select';
import { SIZE_TYPES } from 'types/common/components';
import { defaultFnType, EventCallbackType, MapAny } from 'types/commonTypes';
import { DROPDOWN_SIZE_STYLES } from 'components/common/dropdown/dropdown.constants';
import { LabelProps } from 'components/common/Label';
import { SupporterInfoProps } from 'components/common/SupporterInfo';

export type DropdownProps = {
  options: OptionsType[];
  dropdownIndicator?: React.ReactNode;
  closeIcon?: React.ReactNode;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  onChange?: (selected: any) => void;
  showLabel?: boolean;
  wrapperClass?: string;
  defaultMenuIsOpen?: boolean;
  labelProps?: LabelProps;
  selectFieldWrapperClass?: string;
  error?: boolean;
  errorColor?: string;
  placeholder?: string;
  noOptionsText?: React.ReactNode;
  isMulti?: boolean;
  autoFocus?: boolean;
  spriteSelectedIcon?: string;
  selectedIcon?: React.ReactNode;
  isSearchable?: boolean;
  spriteSelectedIconColor?: string;
  spriteCloseIcon?: string;
  spriteCloseIconColor?: string;
  showSupporterInfo?: boolean;
  supporterInfoProps?: SupporterInfoProps;
  customStyles?: DropdownCustomStyles;
  customClassNames?: DropdownCustomClassNames;
  menuOptionClasses?: MenuOptionClassesProps;
  defaultValue?: OptionsType | OptionsType[] | null;
  handleInputChange?: (inputValue: string) => void;
  handleKeyDown?: (evt: React.KeyboardEvent<HTMLInputElement>) => void;
  onFocus?: (evt: React.FocusEvent<HTMLInputElement>) => void;
  controlled?: boolean;
  value?: OptionsType | OptionsType[] | null;
  enableSelectAll?: boolean;
  showCountOfSelected?: boolean;
  countSelectedSuffix?: string;
  disabled?: boolean;
  readOnly?: boolean;
  id: string;
  eventCallback: EventCallbackType;
  showValueInControl?: boolean;
  controlFocusedColor?: string;
  showLabelInControl?: boolean;
  controlColor?: {
    focused: string;
    background: string;
    overrideBackgroundColor: string;
  };
  tooltipBodyClassName?: string;
  enableReset?: boolean;
  resetProps?: {
    resetClassName: string;
    resetTextClassName: string;
    resetText: string;
  };
  onReset?: defaultFnType;
  size?: SIZE_TYPES;
  menuPortalTarget?: HTMLElement | null;
  customClass?: {
    focus?: string;
    border?: string;
    fontSize?: string;
  };
  enableDelete?: boolean;
  onClickDelete?: defaultFnType;
  isHoveredDropdown?: boolean;
  showSelectedIcon?: boolean;
  customDropdownIndicatorSize?: number;
};

export interface CustomDropdownIndicatorProps extends DropdownIndicatorProps<OptionsType> {
  isHoveredDropdown?: boolean;
  customDropdownIndicatorSize?: number;
}

export interface DropdownCustomClassNames {
  placeholder?: string;
  menu?: string;
  noOptionsMessage?: string;
}

export interface DropdownCustomStyles {
  option?: CSSProperties;
  indicatorSeparator?: CSSProperties;
  menu?: CSSProperties;
  input?: CSSProperties;
  control?: CSSObjectWithLabel;
  multiValue?: CSSProperties;
  multiValueLabel?: CSSProperties;
  multiValueRemove?: CSSProperties;
  singleValue?: CSSProperties;
  valueContainer?: CSSProperties;
  placeholder?: CSSProperties;
  menuList?: CSSProperties;
  noOptionsMessage?: CSSProperties;
  dropdownIndicator?: CSSProperties;
  container?: CSSProperties;
}

export interface MenuOptionProps extends MenuOptionClassesProps {
  innerRef?: RefCallback<HTMLDivElement>;
  innerProps?: MapAny;
  isSelected?: boolean;
  label?: string | ReactNode;
  isMulti?: boolean;
  data?: OptionsType;
  spriteSelectedIcon?: string;
  selectedIcon?: React.ReactNode;
  spriteSelectedIconColor?: string;
  onClick?: defaultFnType;
  eventCallback: EventCallbackType;
  labelOverrideClassName?: string;
  checkboxClassName?: string;
  checkboxDisplayContainerClassName?: string;
  disabled?: boolean;
  isRadio?: boolean;
  radioWrapperStyle?: string;
  radioDefaultStyle?: string;
  radioSelectedStyle?: string;
  radioStyle?: string;
  showSelectedIcon?: boolean;
}

export interface MenuOptionClassesProps {
  contentWrapper?: string;
  wrapperClass?: string;
  containerClass?: string;
  labelOverrideClassName?: string;
}

export interface OptionsType {
  label?: React.ReactNode;
  value: string | number;
  id?: string;
  spriteIcon?: string;
  icon?: React.ReactNode;
  isDisabled?: boolean;
  metadata?: MapAny;
  options?: OptionsType[];
  desc?: string;
}

export interface ChipProps {
  value: string | number;
  label?: string | ReactElement;
  closeIcon?: React.ReactNode;
  onClick?: defaultFnType;
  handleRemoveValue?: defaultFnType;
  className?: string;
  bgClassName?: string;
}

export interface SelectedOptionsProps {
  optionSelected: OptionsType[];
  closeIcon: React.ReactNode;
  handleRemoveValue: defaultFnType;
}

export type MenuSingleValuePropsType = {
  icon?: React.ReactNode;
  label?: string;
  spriteIcon?: string;
  value?: string | number;
  showValueInControl?: boolean;
  size: keyof typeof DROPDOWN_SIZE_STYLES;
  customClassNames?: {
    placeholder?: string;
  };
};

export type SelectedCountTooltipPropsType = {
  tooltipBodyClassName: string;
  value: { value: string; label: string }[];
};

export type CustomReactSelectPropsType = DropdownProps & {
  CustomOption: React.ComponentType<OptionProps<OptionsType>>;
  CustomSingleValue: React.ComponentType<SingleValueProps<OptionsType>>;
  handleChange: (
    selected: MultiValue<OptionsType> | SingleValue<OptionsType>,
    actionMeta: ActionMeta<OptionsType>,
  ) => void;
  isOptionSelected: (option: OptionsType) => boolean;
  addSelectAllInOptions: () => OptionsType[];
  getValue: () => OptionsType[];
  MenuList: ComponentType<MenuListProps<OptionsType>>;
  ValueContainer: ComponentType<ValueContainerProps<OptionsType>>;
  MultiValue: ComponentType<MultiValueProps<OptionsType>>;
  enableDelete?: boolean;
  isHoveredDropdown?: boolean;
};

export type ValueContainerContentProps = {
  labelProps: LabelProps;
  value: { value: string; label: string }[];
  showCountOfSelected: boolean;
  tooltipBodyClassName: string;
};

export type ResetSectionProps = {
  resetProps?: {
    resetClassName?: string;
    resetText?: string;
    resetTextClassName?: string;
  };
  onClickReset: defaultFnType;
};
