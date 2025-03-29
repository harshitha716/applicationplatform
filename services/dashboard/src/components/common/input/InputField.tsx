import React, { FC, memo } from 'react';
import { SIZE_TYPES } from 'types/common/components';
import { InputFieldProps } from 'types/common/components/input/input.types';
import { defaultFn } from 'types/commonTypes';
import InputTag from 'components/common/input/InputTag';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

const sizeClassName = {
  //TODO: Update other classes once design component is ready
  [SIZE_TYPES.XLARGE]: {
    inputClassBySize: 'h-[60px] py-6 pl-16 pr-[58px]',
    inputClassWithoutIcons: 'h-[60px] px-4 py-4.5',
    inputClassWithoutLeadingIcon: 'h-[60px] py-6 pl-4.5 pr-[58px]',
    inputClassWithoutTrailingIcon: 'h-[60px] py-6 pl-16 pr-4.5',
    leadingIconClassBySize: 'left-6 w-6 h-6',
    trailingIconClassBySize: 'right-6 w-6 h-6',
    inputFontClassName: 'f-28-400',
  },
  [SIZE_TYPES.LARGE]: {
    inputClassBySize: 'h-[48px] py-6 pl-16 pr-[58px]',
    inputClassWithoutIcons: 'h-[48px] px-3 py-4.5',
    inputClassWithoutLeadingIcon: 'h-[48px] py-6 pl-4.5 pr-[58px]',
    inputClassWithoutTrailingIcon: 'h-[48px] py-6 pl-16 pr-4.5',
    leadingIconClassBySize: 'left-6 w-6 h-6',
    trailingIconClassBySize: 'right-6 w-6 h-6',
    inputFontClassName: 'f-16-400',
  },
  [SIZE_TYPES.MEDIUM]: {
    inputClassBySize: 'h-[40px] py-4 pl-[60px] pr-[54px]',
    inputClassWithoutIcons: 'h-[40px] p-3',
    inputClassWithoutLeadingIcon: 'h-[40px] py-4 pl-2.5 pr-9',
    inputClassWithoutTrailingIcon: 'h-[40px] py-4 pl-9 pr-2.5',
    leadingIconClassBySize: 'left-3 w-4 h-4',
    trailingIconClassBySize: 'right-3 w-4 h-4',
    inputFontClassName: 'f-14-400 ',
  },
  [SIZE_TYPES.SMALL]: {
    inputClassBySize: 'h-[32px] py-3 pl-14 pr-[50px]',
    inputClassWithoutIcons: 'h-[32px] py-3 px-3',
    inputClassWithoutLeadingIcon: 'h-[32px] py-3 pl-2 pr-9',
    inputClassWithoutTrailingIcon: 'h-[32px] py-3 pl-9 pr-2',
    leadingIconClassBySize: 'left-3 w-4 h-4',
    trailingIconClassBySize: 'right-3 w-4 h-4',
    inputFontClassName: 'f-12-400 ',
  },
  [SIZE_TYPES.XSMALL]: {
    inputClassBySize: 'h-7 py-1.5 pl-2 pr-[30px]',
    inputClassWithoutIcons: 'h-7 py-2.5 px-2',
    inputClassWithoutLeadingIcon: 'h-7 py-1.5 pl-2 pr-6',
    inputClassWithoutTrailingIcon: 'h-7 py-2 pl-6 pr-2',
    leadingIconClassBySize: 'left-2 w-3 h-3',
    trailingIconClassBySize: 'right-2 w-3 h-3',
    inputFontClassName: 'f-13-300',
  },
};

const InputField: FC<InputFieldProps> = ({
  size = SIZE_TYPES.MEDIUM,
  leadingIconClassName = 'absolute flex justify-center items-center',
  trailingIconClassName = 'absolute flex justify-center items-center',
  inputFieldWrapperClassName = 'flex items-center relative ',
  handleLeadingAction = defaultFn,
  handleTrailingAction = defaultFn,
  leadingIconProps = {},
  trailingIconProps = {},
  leadingNode = null,
  ...rest
}) => {
  const {
    inputClassBySize,
    inputClassWithoutIcons,
    inputClassWithoutLeadingIcon,
    inputClassWithoutTrailingIcon,
    leadingIconClassBySize,
    trailingIconClassBySize,
  } = sizeClassName[size];
  const inputSizeClassName =
    leadingIconProps.id && trailingIconProps.id
      ? inputClassBySize
      : leadingIconProps.id
        ? inputClassWithoutTrailingIcon
        : trailingIconProps.id
          ? inputClassWithoutLeadingIcon
          : inputClassWithoutIcons;
  const leadingIconSizeClass = `${leadingIconClassName} ${leadingIconClassBySize}`;
  const trailingIconSizeClass = `${trailingIconClassName} ${trailingIconClassBySize}`;

  return (
    <div className={inputFieldWrapperClassName}>
      {leadingIconProps?.id && leadingIconProps?.iconCategory && (
        <div className={leadingIconSizeClass} role='button' onClick={handleLeadingAction}>
          <SvgSpriteLoader
            id={leadingIconProps?.id}
            iconCategory={leadingIconProps?.iconCategory}
            {...leadingIconProps}
          />
        </div>
      )}

      {leadingNode}

      <InputTag
        inputSizeClassName={inputSizeClassName}
        inputFontClassName={sizeClassName[size].inputFontClassName}
        {...rest}
      />

      {trailingIconProps?.id && trailingIconProps?.iconCategory && (
        <div className={trailingIconSizeClass} role='button' onClick={handleTrailingAction}>
          <SvgSpriteLoader
            id={trailingIconProps?.id}
            iconCategory={trailingIconProps?.iconCategory}
            {...trailingIconProps}
          />
        </div>
      )}
    </div>
  );
};

export default memo(InputField);
