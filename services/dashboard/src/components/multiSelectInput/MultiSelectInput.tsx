import React, { FC, useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { COLORS } from 'constants/colors';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { defaultFn, MapAny } from 'types/commonTypes';
import { checkObjOrArrType, cn } from 'utils/common';
import { Dropdown } from 'components/common/dropdown';
import Input from 'components/common/input';
import CommonWrapper from 'components/commonWrapper';
import { SkeletonTypes } from 'components/commonWrapper/commonWrapper.types';
import { KEY_CODES, MultiSelectInputPropsType } from 'components/multiSelectInput/multiSelectInput.types';
import OptionsListSkeletonLoader from 'components/multiSelectInput/OptionsListSkeletonLoader';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

const MultiSelectInput: FC<MultiSelectInputPropsType> = ({
  id,
  search,
  setSearch,
  isOpen,
  placeholderText,
  roleOptions,
  inputArrayList,
  setInputArrayList,
  showValidationError,
  setShowValidationError,
  validationErrorText,
  onValidateAndAdd,
  optionsList,
  onSelectOption,
  selectOnlyFromList = false,
  transformLabel,
  isLoadingOptionsList,
  optionalOpenDropdownOptions = true,
  wrapperClassName,
  inputWrapperClassName,
  multiSelectInputClassName,
  setIsCustomInputFocused,
  customOptionsListDropdown,
  selectedRole,
  setSelectedRole,
  onCustomDeleteFn,
}) => {
  const inputRef = useRef<HTMLInputElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const dropdownOptionsRef = useRef<HTMLDivElement>(null);
  const optionRefs = useRef<(HTMLDivElement | null)[]>([]);
  const inputPlaceholderText = inputArrayList?.length > 0 ? '' : placeholderText;
  const [isInputFocused, setIsInputFocused] = useState<boolean>(false);
  const [debouncedSearch, setDebouncedSearch] = useState<string>(search);
  const [hoveredOptionIndex, setHoveredOptionIndex] = useState<number>(0);
  const [openDropdownOptions, setOpenDropdownOptions] = useState<boolean>(false);

  const handleSetInputFocus = useCallback(() => {
    setIsInputFocused(true);
    setIsCustomInputFocused?.(true);
    inputRef?.current?.focus();
    setOpenDropdownOptions(true);
  }, [setIsInputFocused, setIsCustomInputFocused, inputRef, setOpenDropdownOptions]);

  const handleSetInputBlur = useCallback(() => {
    setIsInputFocused(false);
    setIsCustomInputFocused?.(false);
    inputRef?.current?.blur();
    setOpenDropdownOptions(false);
  }, [setIsInputFocused, setIsCustomInputFocused, inputRef, setOpenDropdownOptions]);

  useEffect(() => {
    if (isOpen) {
      setHoveredOptionIndex(0);
      if (optionalOpenDropdownOptions) {
        handleSetInputFocus();
      } else {
        setIsInputFocused(true);
        setIsCustomInputFocused?.(true);
        inputRef?.current?.focus();
      }
    }
  }, [isOpen, handleSetInputFocus, optionalOpenDropdownOptions]);

  const handleClickKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    const keyEvent = e.key;

    if (keyEvent === KEY_CODES.BACKSPACE && search.trim() === '') {
      if (inputArrayList?.length > 0) {
        handleRemoveItem(inputArrayList[inputArrayList?.length - 1], inputArrayList?.length - 1);
      }
      handleSetInputFocus();

      return;
    }

    if (selectOnlyFromList) {
      if (keyEvent === KEY_CODES.ENTER || keyEvent === KEY_CODES.COMMA || keyEvent === KEY_CODES.SPACE) {
        e.preventDefault();

        const selectedOption = (filteredDropdownOptions ?? [])[hoveredOptionIndex ?? 0];

        if (selectedOption) {
          onValidateAndAdd({
            value: selectedOption?.value,
            label: selectedOption?.label,
            color: selectedOption?.color,
            type: selectedOption?.type,
            team_id: selectedOption?.team_id,
          });

          setSearch('');
        }
      } else if (keyEvent === KEY_CODES.ARROW_DOWN || keyEvent === KEY_CODES.ARROW_UP) {
        handleKeyDown(e);
      }
    } else {
      const trimmedSearch = search?.trim();

      if (
        (keyEvent === KEY_CODES.ENTER || keyEvent === KEY_CODES.COMMA || keyEvent === KEY_CODES.SPACE) &&
        trimmedSearch
      ) {
        e.preventDefault();
        onValidateAndAdd({
          value: trimmedSearch,
          label: trimmedSearch,
        });
        setSearch('');
      }
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLDivElement>) => {
    const keyEvent = e.key;

    if (!filteredDropdownOptions?.length) return;

    if (keyEvent === KEY_CODES.ARROW_DOWN) {
      e.preventDefault();
      setHoveredOptionIndex((prevIndex) => {
        const newIndex = prevIndex === null || prevIndex === filteredDropdownOptions?.length - 1 ? 0 : prevIndex + 1;

        optionRefs?.current[newIndex]?.scrollIntoView({
          behavior: 'smooth',
          block: 'nearest',
        });

        return newIndex;
      });
    } else if (keyEvent === KEY_CODES.ARROW_UP) {
      e.preventDefault();
      setHoveredOptionIndex((prevIndex) => {
        const newIndex = prevIndex === null || prevIndex === 0 ? filteredDropdownOptions?.length - 1 : prevIndex - 1;

        optionRefs?.current[newIndex]?.scrollIntoView({
          behavior: 'smooth',
          block: 'nearest',
        });

        return newIndex;
      });
    } else if (keyEvent === KEY_CODES.ENTER && hoveredOptionIndex !== null) {
      e.preventDefault();
      handleSelectDropdownOption(filteredDropdownOptions[hoveredOptionIndex]);
      handleSetInputFocus();
    }
  };

  const handleRemoveItem = useCallback(
    (item: MapAny, index: number) => {
      if (onCustomDeleteFn) {
        onCustomDeleteFn(item);

        return;
      }

      const updatedItems = inputArrayList?.filter((_, i) => i !== index);

      if (setShowValidationError) {
        setShowValidationError(updatedItems?.some((item) => !item.valid));
      }

      setInputArrayList(updatedItems);

      handleSetInputFocus();
    },
    [inputArrayList, setInputArrayList, setShowValidationError, handleSetInputFocus],
  );

  const handleClickOutside = useCallback(
    (event: MouseEvent) => {
      if (
        containerRef?.current?.contains(event.target as Node) ||
        dropdownOptionsRef?.current?.contains(event.target as Node)
      ) {
        return;
      }
      setTimeout(() => {
        setIsInputFocused(false);
        setOpenDropdownOptions(false);
        setIsCustomInputFocused?.(false);
      }, 0);
    },
    [containerRef, dropdownOptionsRef],
  );

  useEffect(() => {
    document.addEventListener('mousedown', handleClickOutside);

    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [handleClickOutside]);

  useEffect(() => {
    const debounceHandler = setTimeout(() => {
      setDebouncedSearch(search);
    }, 300);

    return () => {
      if (debounceHandler) clearTimeout(debounceHandler);
    };
  }, [search]);

  const combinedOptions = useMemo(() => optionsList ?? [], [optionsList]);

  const filteredDropdownOptions = useMemo(() => {
    if (!combinedOptions) return [];
    if (!debouncedSearch?.trim()) return combinedOptions;

    const filteredOptions = combinedOptions?.filter((option) =>
      option?.value.toLowerCase().startsWith(debouncedSearch.toLowerCase()),
    );

    setOpenDropdownOptions(filteredOptions?.length > 0);

    return filteredOptions;
  }, [combinedOptions, debouncedSearch]);

  const handleSelectDropdownOption = useCallback(
    (option: { value: string; label: string; color?: string; type?: string; team_id?: string }) => {
      onSelectOption?.(option);
      setSearch('');
      handleSetInputFocus();
    },
    [onSelectOption, setSearch, inputRef, showValidationError],
  );

  // if selectedRole is not object, finds the option from roleOptions, if yes, converts into { label, value } with a stringified version
  const selectedRoleValue = useMemo(() => {
    if (checkObjOrArrType(selectedRole, 'object')) {
      return { label: JSON.stringify(selectedRole), value: JSON.stringify(selectedRole) };
    }

    return (roleOptions ?? [])?.find((option) => option?.value === selectedRole) ?? roleOptions?.[0];
  }, [selectedRole, roleOptions]);

  useEffect(() => {
    if (debouncedSearch?.trim()) {
      setOpenDropdownOptions((prev) =>
        prev !== filteredDropdownOptions?.length > 0 ? filteredDropdownOptions?.length > 0 : prev,
      );
    }
    setHoveredOptionIndex(0);
  }, [filteredDropdownOptions, debouncedSearch]);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      const keyEvent = e.key;

      if (keyEvent === KEY_CODES.ESCAPE && isOpen) {
        handleSetInputBlur();
      } else if (keyEvent === KEY_CODES.ENTER && isOpen) {
        handleSetInputFocus();
      }
    };

    document.addEventListener('keydown', handleKeyDown);

    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [isOpen]);

  return (
    <div className='flex flex-col items-center'>
      <div
        className={cn(
          'flex justify-between items-start w-full rounded-md gap-1.5 border',
          isInputFocused ? 'border-GRAY_600 shadow-inputOutlineShadow' : 'border-GRAY_400',
          wrapperClassName,
        )}
      >
        <div
          className={cn('flex flex-wrap gap-1.5 py-3 pl-3 w-full', inputWrapperClassName)}
          ref={containerRef}
          onClick={handleSetInputFocus}
        >
          {Array.isArray(inputArrayList) &&
            inputArrayList?.length > 0 &&
            inputArrayList?.map((item, index) => (
              <div
                key={index}
                className='flex items-center gap-1 px-1.5 pr-1 py-0.5 rounded w-fit h-fit cursor-default'
                style={{
                  backgroundColor: item?.valid ? (item?.color ? item?.color : COLORS.GRAY_50) : COLORS.RED_100,
                  border: `1px solid ${item?.valid ? (item?.color !== COLORS.WHITE ? 'transparent' : COLORS.GRAY_400) : COLORS.RED_200}`,
                }}
              >
                <span className='f-12-500 text-GRAY_1000'>{item?.label}</span>
                <SvgSpriteLoader
                  id='x-close'
                  iconCategory={ICON_SPRITE_TYPES.GENERAL}
                  width={10}
                  height={10}
                  onClick={() => handleRemoveItem(item, index)}
                  color={item?.valid ? COLORS.GRAY_700 : COLORS.GRAY_900}
                  className='cursor-pointer'
                />
              </div>
            ))}
          <Input
            placeholder={inputPlaceholderText}
            type='email'
            inputRef={inputRef}
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            onKeyDown={handleClickKeyDown}
            className='flex-1 min-w-[20px] h-fit mt-[2px]'
            customPaddingClassName='p-0'
            focusClassNames='focus:outline-none focus:border-none focus:shadow-none'
            cursorClassname='cursor-default'
            inputFontClassName={multiSelectInputClassName || 'f-13-400 py-0 !rounded-none'}
            style={{
              maxWidth: '100%',
            }}
          />
        </div>
        {roleOptions && (
          <div className='flex min-w-max h-fit !cursor-pointer'>
            <Dropdown
              options={roleOptions}
              id={`${id}-multi-select-input-dropdown`}
              eventCallback={defaultFn}
              onChange={(selectedOption) => {
                setSelectedRole?.(selectedOption?.value);
                inputRef.current?.focus();
                setIsInputFocused(true);
              }}
              value={selectedRoleValue}
              defaultValue={roleOptions[0]}
              placeholder='Member'
              isSearchable={false}
              customClass={{
                focus: 'none',
                border: 'transparent',
                fontSize: 'f-12-400',
              }}
              customClassNames={{
                placeholder: 'f-12-300',
              }}
              menuOptionClasses={{
                contentWrapper: 'py-2',
              }}
              customDropdownIndicatorSize={14}
            />
          </div>
        )}
      </div>
      {openDropdownOptions && (
        <div className='w-full relative'>
          <div ref={dropdownOptionsRef} onClick={(e) => e.stopPropagation()}>
            {customOptionsListDropdown
              ? React.createElement(
                  customOptionsListDropdown as React.ElementType,
                  {
                    search,
                    optionRefs,
                    transformLabel,
                    hoveredOptionIndex,
                    isLoadingOptionsList,
                    setHoveredOptionIndex,
                    onKeyDown: handleKeyDown,
                    onCloseDropdown: handleSetInputBlur,
                    optionList: filteredDropdownOptions,
                    onSelectOption: handleSelectDropdownOption,
                    isDropdownOpenOnZeroLength: !!filteredDropdownOptions?.length,
                  } as Record<string, unknown>,
                )
              : !!combinedOptions?.length && (
                  <div className='absolute left-0 bg-white w-full p-1 f-10-500 text-GRAY_700 rounded-md border border-GRAY_400 mt-1 z-10 shadow-tableFilterMenu'>
                    <span className='flex pt-2 pb-1.5 px-1.5'>Select a team or person</span>
                    <div
                      className='flex flex-col w-full max-h-[200px] overflow-y-auto [&::-webkit-scrollbar]:hidden'
                      tabIndex={0}
                    >
                      <CommonWrapper
                        skeletonType={SkeletonTypes.CUSTOM}
                        isLoading={isLoadingOptionsList}
                        loader={<OptionsListSkeletonLoader />}
                      >
                        {filteredDropdownOptions?.map((option, index) => (
                          <div
                            key={index}
                            ref={(el) => {
                              if (optionRefs?.current) {
                                optionRefs.current[index] = el;
                              }
                            }}
                            className={cn('w-full px-1.5 py-1 hover:bg-GRAY_50 rounded-md cursor-pointer', {
                              'bg-GRAY_50':
                                (hoveredOptionIndex === null && index === 0) || hoveredOptionIndex === index,
                            })}
                            onMouseEnter={() => setHoveredOptionIndex(index)}
                            onClick={() => handleSelectDropdownOption(option)}
                            onKeyDown={handleKeyDown}
                          >
                            <span
                              className={cn(
                                'f-12-400 text-GRAY_1000 flex px-1.5 py-0.5 w-fit rounded capitalize border',
                              )}
                              style={{
                                backgroundColor: option?.color || COLORS.WHITE,
                                borderColor: option?.color ? 'transparent' : COLORS.GRAY_400,
                              }}
                            >
                              {transformLabel ? transformLabel(option?.label) : option?.label}
                            </span>
                          </div>
                        ))}
                      </CommonWrapper>
                    </div>
                  </div>
                )}
          </div>
        </div>
      )}
      {validationErrorText && showValidationError && (
        <span className='f-11-400 text-RED_700 mt-2 w-full flex text-start'>{validationErrorText}</span>
      )}
    </div>
  );
};

export default MultiSelectInput;
