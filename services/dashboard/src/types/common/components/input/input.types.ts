import React, { HTMLInputTypeAttribute, ReactNode } from 'react';
import { SIZE_TYPES } from 'types/common/components';
import { defaultFnType, EventCallbackType } from 'types/commonTypes';
import { SupporterInfoProps } from 'components/common/SupporterInfo';
import { SvgSpriteLoaderProps } from 'components/SvgSpriteLoader';

export enum SUPPORT_INFO_TYPES {
  GUIDE = 'GUIDE',
  ERROR = 'ERROR',
  CUSTOM = 'CUSTOM',
}
export interface InputTagProps {
  id?: string;
  name?: string;
  value?: string | number;
  type?: HTMLInputTypeAttribute;
  placeholder?: string;
  error?: string;
  maxLength?: number;
  disabled?: boolean;
  readOnly?: boolean;
  style?: React.CSSProperties;
  autocomplete?: string;
  inputTagWrapperClassName?: string;
  inputClassName?: string;
  inputSizeClassName?: string;
  errorClass?: string;
  inputRef?: React.RefObject<HTMLInputElement>;
  onChange?: (evt: React.ChangeEvent<HTMLInputElement>) => void;
  onKeyPress?: (evt: React.KeyboardEvent<HTMLInputElement>) => void;
  onKeyDown?: (evt: React.KeyboardEvent<HTMLInputElement>) => void;
  onBlur?: (evt: React.FocusEvent<HTMLInputElement>) => void;
  onFocus?: (evt: React.FocusEvent<HTMLInputElement>) => void;
  onDeleteTag?: (index: number) => void;
  eventId?: string;
  eventCallback?: EventCallbackType;
  eventCallbackDelay?: number;
  inputFontClassName?: string;
  autoFocus?: boolean;
  tabIndex?: number;
  onKeyUp?: React.KeyboardEventHandler<HTMLInputElement>;
  noBorders?: boolean;
  overrideInputBgClassName?: string;
  inputRoundedClassName?: string;
  inputTagBorderClassName?: string;
  isMulti?: boolean;
  tags?: string[];
  customTags?: ReactNode;
  onEnterKey?: (evt: React.KeyboardEvent<HTMLInputElement>) => void;
  inputPillsWrapperClasses?: string;
  focusClassNames?: string;
  cursorClassname?: string;
  customPaddingClassName?: string;
}

export interface InputFieldProps extends InputTagProps {
  size?: SIZE_TYPES;
  inputTagProps?: InputTagProps;
  leadingIconClassName?: string;
  trailingIconClassName?: string;
  inputFieldWrapperClassName?: string;
  handleLeadingAction?: defaultFnType;
  handleTrailingAction?: defaultFnType;
  leadingIconProps?: Partial<SvgSpriteLoaderProps>;
  trailingIconProps?: Partial<SvgSpriteLoaderProps>;
  leadingNode?: ReactNode;
}

export interface InputProps extends InputTagProps, InputFieldProps {
  supporterInfoProps?: SupporterInfoProps;
  showSupporterInfo?: boolean;
  inputWrapperClassName?: string;
  className?: string;
  label?: string;
  description?: string;
  placeholder?: string;
  labelClassName?: string;
  labelOverrideClassName?: string;
  required?: boolean;
  testId?: string;
}
