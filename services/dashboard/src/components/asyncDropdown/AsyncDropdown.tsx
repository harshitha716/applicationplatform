import React, { FC, useEffect, useRef, useState } from 'react';
import { COLORS } from 'constants/colors';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { useOnClickOutside } from 'hooks';
import { cn } from 'utils/common';
import { AsyncDropdownPropsType } from 'components/asyncDropdown/asyncDropdown.types';
import { KEY_CODES } from 'components/multiSelectInput/multiSelectInput.types';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

const AsyncDropdown: FC<AsyncDropdownPropsType> = ({
  onOpen,
  onClose,
  isOpen,
  onDelete,
  onChange,
  options,
  selectedValue,
  showDelete,
  isHoveredDropdown,
  setIsHoveredDropdown,
  wrapperClassName,
  parentWrapperClassName,
  showSelectedIcon,
  selectedOptionClassName,
  isOverflowStyle,
}) => {
  const [dropdownTop, setDropdownTop] = useState(0);
  const buttonRef = useRef<HTMLDivElement>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);

  const handleToggleDropdown = () => {
    if (isOpen) {
      onClose();
    } else {
      onOpen();
    }
  };

  useOnClickOutside(dropdownRef, onClose);

  useEffect(() => {
    if (isOpen && buttonRef?.current) {
      const rect = buttonRef?.current.getBoundingClientRect();

      setDropdownTop(rect.bottom - 40);
    }
  }, [isOpen, options.length]);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      const keyEvent = e.key;

      if (keyEvent === KEY_CODES.ESCAPE && isOpen) {
        onClose();
      }
    };

    document.addEventListener('keydown', handleKeyDown);

    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, onClose]);

  return (
    <div
      ref={dropdownRef}
      onMouseEnter={() => setIsHoveredDropdown && setIsHoveredDropdown(true)}
      onMouseLeave={() => setIsHoveredDropdown && setIsHoveredDropdown(false)}
    >
      <div
        className={cn(
          'flex justify-between items-center py-3 pl-4 gap-0 cursor-pointer h-10 f-12-400',
          parentWrapperClassName,
        )}
        onClick={handleToggleDropdown}
        ref={buttonRef}
      >
        {selectedValue?.label}
        {typeof isHoveredDropdown !== 'undefined' && (
          <div className='ml-1'>
            <SvgSpriteLoader
              id={isOpen ? 'chevron-up' : 'chevron-down'}
              iconCategory={ICON_SPRITE_TYPES.ARROWS}
              width={12}
              height={12}
              color={isHoveredDropdown ? COLORS.GRAY_1000 : COLORS.WHITE}
            />
          </div>
        )}
      </div>
      {isOpen && (
        <div
          className={cn(
            'flex flex-col border border-GRAY_50 rounded-md p-1 absolute right-0 max-w-[170px] min-w-max bg-white z-1000 shadow-tableFilterMenu',
            wrapperClassName,
          )}
          style={{
            top: `${isOverflowStyle ? dropdownTop : 40}px`,
          }}
        >
          {options.map((role) => (
            <div
              key={role.value}
              className={cn(
                'flex flex-col py-2 pl-2.5 pr-2 hover:bg-GRAY_100 cursor-pointer rounded-md',
                role.value === selectedValue?.value && selectedOptionClassName,
              )}
              onClick={() => onChange(role)}
            >
              <span className='flex justify-between items-start f-12-500 text-GRAY_1000'>
                {role?.label}
                {showSelectedIcon && role?.value === selectedValue?.value && (
                  <SvgSpriteLoader
                    id='check'
                    iconCategory={ICON_SPRITE_TYPES.GENERAL}
                    width={14}
                    height={14}
                    color={COLORS.GRAY_900}
                  />
                )}
              </span>
              {!!role?.desc && <span className='f-10-500 text-GRAY_700 mt-1.5'>{role?.desc}</span>}
            </div>
          ))}
          {showDelete && (
            <span
              className='flex gap-1.5 items-center f-12-500 text-RED_700 py-2 px-2.5 border-t border-DIVIDER_GRAY cursor-pointer'
              onClick={onDelete}
            >
              <SvgSpriteLoader
                id='trash-04'
                iconCategory={ICON_SPRITE_TYPES.GENERAL}
                width={12}
                height={12}
                color={COLORS.RED_700}
              />
              Remove
            </span>
          )}
        </div>
      )}
    </div>
  );
};

export default AsyncDropdown;
