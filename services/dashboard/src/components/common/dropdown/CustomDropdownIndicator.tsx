import React, { useMemo } from 'react';
import { components } from 'react-select';
import { COLORS } from 'constants/colors';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { SIZE_TYPES } from 'types/common/components';
import { CustomDropdownIndicatorProps } from 'types/common/components/dropdown/dropdown.types';
import { DROPDOWN_SIZE_STYLES } from 'components/common/dropdown/dropdown.constants';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

export const CustomDropdownIndicator = (props: CustomDropdownIndicatorProps) => {
  const { selectProps = {}, isHoveredDropdown, customDropdownIndicatorSize } = props;

  // @ts-ignore selectProps contains all props passed to react select. It's passed to each child component of react-select and takes custom props as well.
  const { size, menuIsOpen } = selectProps;

  const ChevronIcon = useMemo(
    () => (
      <div className='ml-1'>
        <SvgSpriteLoader
          id={menuIsOpen ? 'chevron-up' : 'chevron-down'}
          iconCategory={ICON_SPRITE_TYPES.ARROWS}
          width={
            customDropdownIndicatorSize
              ? customDropdownIndicatorSize
              : DROPDOWN_SIZE_STYLES[size as SIZE_TYPES].dropdownIndicatorProps.width
          }
          height={
            customDropdownIndicatorSize
              ? customDropdownIndicatorSize
              : DROPDOWN_SIZE_STYLES[size as SIZE_TYPES].dropdownIndicatorProps.height
          }
          color={COLORS.GRAY_900}
        />
      </div>
    ),
    [menuIsOpen, customDropdownIndicatorSize, size],
  );

  return (
    <components.DropdownIndicator {...props}>
      {typeof isHoveredDropdown !== 'undefined' ? (
        <div className='ml-1'>
          <SvgSpriteLoader
            id={menuIsOpen ? 'chevron-up' : 'chevron-down'}
            iconCategory={ICON_SPRITE_TYPES.ARROWS}
            width={12}
            height={12}
            color={isHoveredDropdown ? COLORS.GRAY_1000 : COLORS.WHITE}
          />
        </div>
      ) : (
        ChevronIcon
      )}
    </components.DropdownIndicator>
  );
};
