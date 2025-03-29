import React, { FC } from 'react';
import { COLORS } from 'constants/colors';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { MenuOptionProps } from 'types/common/components/dropdown/dropdown.types';
import { defaultFn } from 'types/commonTypes';
import { cn } from 'utils/common';
import { CheckBox } from 'components/common/Checkbox';
import { Radio } from 'components/common/Radio';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

const leadingIconClassBySize = 'w-6 h-6';
const leadingIconClass = 'flex justify-center items-center mr-4';
const leadingIconSizeClass = `${leadingIconClass} ${leadingIconClassBySize}`;
const selectedIconClass = 'ml-auto';

export const MenuOption: FC<MenuOptionProps> = ({
  innerProps = {},
  isSelected = false,
  label,
  isMulti = false,
  data,
  spriteSelectedIcon = '',
  selectedIcon = null,
  containerClass = 'hover:bg-BASE_SECONDARY first:rounded-t-[10px] last:rounded-b-[10px]',
  contentWrapper = 'pl-2 py-3 w-full',
  wrapperClass = 'p-2 h-16 flex items-center',
  spriteSelectedIconColor = COLORS.GREEN_SECONDARY,
  onClick = defaultFn,
  labelOverrideClassName = 'f-16-400',
  checkboxClassName = 'pr-[33px] pl-[15px] h-12',
  checkboxDisplayContainerClassName = 'top-[15px] left-[15px]',
  disabled = false,
  isRadio = false,
  radioWrapperStyle = '',
  radioDefaultStyle = '',
  radioSelectedStyle = '',
  radioStyle = '',
  showSelectedIcon,
}) => (
  <div
    onClick={disabled ? undefined : onClick}
    className={cn('flex flex-row items-start justify-between', containerClass)}
    data-testid={`menu-option-container-${innerProps.id}`}
  >
    <div className='absolute -translate-x-1 p-1 right-0'>
      <SvgSpriteLoader
        id='check'
        iconCategory={ICON_SPRITE_TYPES.GENERAL}
        width={15}
        height={15}
        color={isSelected && showSelectedIcon ? COLORS.GRAY_900 : COLORS.TRANSPARENT}
      />
    </div>
    <div className={wrapperClass} data-testid={`menu-option-wrapper-${innerProps.id}`}>
      {isMulti && (
        <CheckBox
          checked={isSelected}
          onPress={defaultFn}
          id={innerProps.id}
          className={checkboxClassName}
          displayContainerClassName={checkboxDisplayContainerClassName}
          disabled={disabled}
        />
      )}
      {isRadio && (
        <Radio
          checked={isSelected}
          onSelect={defaultFn}
          id={innerProps.id}
          wrapperStyle={radioWrapperStyle}
          radioDefaultStyle={radioDefaultStyle}
          radioSelectedStyle={radioSelectedStyle}
          radioStyle={radioStyle}
        />
      )}
      <div
        className={cn(`${contentWrapper} ${isMulti ? '' : 'py-1'} flex flex-col justify-start items-start`)}
        data-testid={`menu-option-content-${innerProps.id}`}
      >
        {data?.spriteIcon && (
          <div className={leadingIconSizeClass} data-testid={`menu-option-leading-icon-wrapper-${innerProps.id}`}>
            <SvgSpriteLoader id={data.spriteIcon} height={15} width={15} />
          </div>
        )}
        {data?.icon ?? null}

        {typeof label === 'string' ? (
          <span className={cn('text-TEXT_PRIMARY', labelOverrideClassName)} id={`menu-option-${innerProps.id}`}>
            {label}
          </span>
        ) : (
          label
        )}
        {data?.desc && (
          <span className='f-10-500 text-GRAY_700 mt-1.5' id={`menu-option-${innerProps.id}`}>
            {data?.desc}
          </span>
        )}
      </div>
      {!isMulti && spriteSelectedIcon && isSelected && (
        <div
          className={cn(leadingIconSizeClass, selectedIconClass)}
          data-testid={`menu-option-selected-icon-wrapper-${innerProps.id}`}
        >
          <SvgSpriteLoader id={spriteSelectedIcon} color={spriteSelectedIconColor} />
        </div>
      )}
      {!isMulti && isSelected && selectedIcon}
    </div>
  </div>
);
