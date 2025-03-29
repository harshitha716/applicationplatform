import { FC, useMemo } from 'react';
import Select from 'react-select';
import { COLORS } from 'constants/colors';
import { SIZE_TYPES } from 'types/common/components';
import {
  CustomDropdownIndicatorProps,
  CustomReactSelectPropsType,
} from 'types/common/components/dropdown/dropdown.types';
import { cn } from 'utils/common';
import { CustomDropdownIndicator } from 'components/common/dropdown/CustomDropdownIndicator';
import { CustomMultivalueRemove } from 'components/common/dropdown/CustomMultivalueRemove';
import { DROPDOWN_SIZE_STYLES } from 'components/common/dropdown/dropdown.constants';

const CustomReactSelect: FC<CustomReactSelectPropsType> = ({
  enableSelectAll,
  options,
  defaultValue,
  isMulti,
  defaultMenuIsOpen,
  placeholder,
  autoFocus,
  noOptionsText,
  CustomOption,
  CustomSingleValue,
  enableReset,
  MenuList,
  showLabelInControl,
  ValueContainer,
  showCountOfSelected,
  customStyles,
  size = SIZE_TYPES.MEDIUM,
  controlColor,
  menuPortalTarget,
  id,
  controlled,
  onFocus,
  handleChange,
  handleKeyDown,
  handleInputChange,
  isOptionSelected,
  isSearchable,
  disabled,
  readOnly,
  customClassNames,
  error,
  errorColor,
  addSelectAllInOptions,
  getValue,
  MultiValue,
  customClass,
  enableDelete,
  isHoveredDropdown,
  showSelectedIcon,
  customDropdownIndicatorSize,
}) => {
  const memoizedDropdownIndicator = useMemo(() => {
    const Component = (props: CustomDropdownIndicatorProps) => (
      <CustomDropdownIndicator
        {...props}
        customDropdownIndicatorSize={customDropdownIndicatorSize}
        isHoveredDropdown={isHoveredDropdown}
      />
    );

    Component.displayName = 'memoized-react-select-dropdown-indicator';

    return Component;
  }, [customDropdownIndicatorSize, isHoveredDropdown]);

  return (
    <Select
      name='select'
      options={enableSelectAll ? addSelectAllInOptions() : options}
      defaultValue={defaultValue}
      isMulti={isMulti}
      defaultMenuIsOpen={defaultMenuIsOpen}
      placeholder={placeholder}
      menuPosition='absolute'
      autoFocus={autoFocus}
      noOptionsMessage={() => noOptionsText}
      closeMenuOnSelect={!isMulti}
      hideSelectedOptions={false}
      components={{
        Option: CustomOption,
        DropdownIndicator: memoizedDropdownIndicator,
        MultiValueRemove: CustomMultivalueRemove,
        SingleValue: CustomSingleValue,
        ...(enableReset ? { MenuList } : {}),
        ...(showLabelInControl ? { ValueContainer } : showCountOfSelected ? { MultiValue } : {}),
        ...(enableDelete ? { MenuList } : {}),
      }}
      styles={{
        container: (styles) => ({
          ...styles,
          ...customStyles?.container,
        }),
        option: (styles, { isSelected }) => ({
          ...styles,
          fontFamily: 'Inter',
          fontWeight: 500,
          borderRadius: '6px',
          paddingTop: '4px',
          paddingBottom: '4px',
          paddingLeft: '0',
          paddingRight: '8px',
          cursor: 'pointer',
          backgroundColor: isSelected ? (showSelectedIcon ? 'white' : COLORS.GRAY_100) : 'white',
          color: isSelected ? COLORS.GRAY_1000 : COLORS.GRAY_900,
          ':active': {
            backgroundColor: COLORS.GRAY_100,
          },
          ':hover': {
            backgroundColor: COLORS.GRAY_100,
          },
          ...DROPDOWN_SIZE_STYLES[size].customStyles.option,
          ...customStyles?.option,
        }),
        indicatorSeparator: () => ({
          appearance: 'none',
          ...customStyles?.indicatorSeparator,
        }),
        dropdownIndicator: (styles) => ({
          ...styles,
          padding: 0,
          ...customStyles?.dropdownIndicator,
        }),
        menu: (styles) => ({
          ...styles,
          fontFamily: 'Inter',
          fontWeight: 500,
          boxShadow: '0px 4px 15px 0px rgba(166, 166, 166, 0.20);',
          border: '1px solid var(--GRAY_400)',
          height: 'fit-content',
          cursor: 'pointer',
          width: 'max-content',
          position: 'absolute',
          right: 0,
          ...DROPDOWN_SIZE_STYLES[size].customStyles.menu,
          ...customStyles?.menu,
        }),
        input: (styles) => ({
          ...styles,
          'input:focus': {
            boxShadow: 'none',
          },

          margin: 0,
          padding: 0,
          ...DROPDOWN_SIZE_STYLES[size].customStyles.input,
          ...customStyles?.input,
        }),
        control: (styles, data) => {
          let backgroundColor = controlColor?.background;

          if (controlColor?.overrideBackgroundColor) {
            backgroundColor = controlColor?.overrideBackgroundColor;
          } else if (readOnly || disabled) {
            backgroundColor = COLORS.DIVIDER_GRAY_1;
          }

          return {
            ...styles,
            backgroundColor,
            fontFamily: 'Inter',
            fontWeight: 500,
            borderWidth: '1px',
            width: 'fit-content',
            borderColor: customClass?.border
              ? customClass?.border
              : error
                ? errorColor
                : data?.isFocused
                  ? COLORS.GRAY_600
                  : COLORS.GRAY_400,
            boxShadow: customClass?.focus
              ? customClass?.focus
              : data?.isFocused
                ? '0px 0px 0px 3px var(--GRAY_400)'
                : 'none',
            ':hover': {
              borderColor: customClass?.border ? customClass?.border : COLORS.GRAY_600,
            },
            ':active': {},
            ...DROPDOWN_SIZE_STYLES[size].customStyles.control,
            ...customStyles?.control,
          };
        },
        multiValue: (styles) => ({
          ...(showCountOfSelected
            ? {}
            : {
                ...styles,
                backgroundColor: COLORS.ZAMP_SECONDARY,
                borderRadius: '0px',
                padding: '4px 8px 4px 3px',
                height: '24px',
                display: 'flex',
                alignItems: 'center',
                margin: '2px 8px 2px 0px',
              }),
          ...customStyles?.multiValue,
        }),
        multiValueLabel: (styles) => ({
          ...styles,
          color: COLORS.TEXT_SECONDARY,
          padding: 0,
          fontSize: '12px',
          fontWeight: 300,
          lineHeight: '16px',
          ...customStyles?.multiValueLabel,
        }),
        multiValueRemove: (styles) => ({
          ...styles,
          padding: 0,
          height: '16px',
          ':hover': {
            backgroundColor: COLORS.ZAMP_SECONDARY,
          },
          ...customStyles?.multiValueRemove,
        }),
        singleValue: (styles) => ({
          ...styles,
          margin: 0,
          ...customStyles?.singleValue,
        }),
        valueContainer: (styles) => ({
          ...styles,
          padding: 0,
          minHeight: '0',
          ...DROPDOWN_SIZE_STYLES[size].customStyles.valueContainer,
          ...customStyles?.valueContainer,
        }),
        placeholder: (styles) => ({
          ...styles,
          marginLeft: 0,
          color: disabled ? COLORS.GRAY_700 : COLORS.TEXT_TERTIARY,
          ...customStyles?.placeholder,
        }),
        menuList: (styles) => ({
          ...styles,
          padding: 0,
          ...customStyles?.menuList,
        }),
        noOptionsMessage: (styles) => ({
          ...styles,
          color: COLORS.TEXT_TERTIARY,
          ...customStyles?.noOptionsMessage,
          ...DROPDOWN_SIZE_STYLES[size].customStyles.noOptionsMessage,
        }),
        groupHeading: (styles) => ({
          ...styles,
          color: COLORS.TEXT_TERTIARY,
          padding: '16px 24px',
          fontSize: '13px',
        }),
        menuPortal: (styles) => ({
          ...styles,
          zIndex: 2,
        }),
      }}
      classNames={{
        placeholder: () =>
          cn(
            `${customClass?.fontSize ? customClass?.fontSize : (customClassNames?.placeholder ?? DROPDOWN_SIZE_STYLES[size].customClassNames.placeholder)}`,
          ),
        menu: () => cn(`${customClassNames?.menu ?? 'bg-white border-0.5 border-DIVIDER_GRAY'}`),
        noOptionsMessage: () => cn(customClassNames?.noOptionsMessage ?? 'h-16 flex items-center justify-center'),
      }}
      onChange={handleChange}
      onKeyDown={handleKeyDown}
      isClearable={false}
      isDisabled={disabled}
      isSearchable={isSearchable && !readOnly}
      onInputChange={handleInputChange}
      onFocus={onFocus}
      {...(enableSelectAll ? { isOptionSelected: isOptionSelected } : {})}
      {...(controlled ? { value: getValue() } : {})}
      menuPortalTarget={menuPortalTarget}
      classNamePrefix={id}
      // @ts-ignore selectProps contains all props passed to react select. It's passed to each child component of react-select and takes custom props as well.
      size={size}
    />
  );
};

export default CustomReactSelect;
