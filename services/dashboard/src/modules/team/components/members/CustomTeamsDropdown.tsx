import React, { FC, useEffect, useRef } from 'react';
import { COLORS } from 'constants/colors';
import { CustomTeamsDropdownPropsType } from 'modules/team/people.types';
import { cn } from 'utils/common';
import CommonWrapper from 'components/commonWrapper';
import { SkeletonTypes } from 'components/commonWrapper/commonWrapper.types';
import { ArrayListOption, KEY_CODES } from 'components/multiSelectInput/multiSelectInput.types';
import OptionsListSkeletonLoader from 'components/multiSelectInput/OptionsListSkeletonLoader';

const CustomTeamsDropdown: FC<CustomTeamsDropdownPropsType> = ({
  search,
  optionRefs,
  optionList,
  isLoadingOptionsList,
  hoveredOptionIndex,
  setHoveredOptionIndex,
  onSelectOption,
  onKeyDown,
  transformLabel,
  randomColor,
  onCloseDropdown,
}) => {
  const dropdownRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      const keyEvent = e.key;

      if (keyEvent === KEY_CODES.ESCAPE) {
        onCloseDropdown();
      }
    };

    document.addEventListener('keydown', handleKeyDown);

    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [onKeyDown]);

  return (
    <div className='absolute left-0 bg-white max-w-48 w-fit p-1 f-10-500 text-GRAY_700 rounded-md border border-GRAY_400 mt-1 z-10 shadow-tableFilterMenu'>
      <span className='flex pt-2 pb-1.5 px-1.5 whitespace-nowrap'>Select a team or create one</span>
      <div
        className='flex flex-col w-full max-h-[200px] overflow-y-auto [&::-webkit-scrollbar]:hidden outline-none'
        ref={dropdownRef}
        tabIndex={0}
        onKeyDown={onKeyDown}
        onClick={(e) => e.stopPropagation()}
      >
        <CommonWrapper
          skeletonType={SkeletonTypes.CUSTOM}
          isLoading={isLoadingOptionsList}
          loader={<OptionsListSkeletonLoader />}
        >
          {optionList?.map((option, index) => (
            <div
              key={index}
              ref={(el) => {
                if (optionRefs?.current) {
                  optionRefs.current[index] = el;
                }
              }}
              onMouseEnter={() => setHoveredOptionIndex(index)}
              onClick={() => {
                const optionWithNew = option as ArrayListOption & { isNew?: boolean };

                if (optionWithNew?.isNew && !search?.trim()) return;
                onSelectOption(option);
              }}
            >
              {(option as ArrayListOption & { isNew?: boolean })?.isNew && !!search?.length ? (
                <div
                  className={cn('w-full px-1.5 py-1 hover:bg-GRAY_50 rounded-md cursor-pointer', {
                    'bg-GRAY_50': (hoveredOptionIndex === null && index === 0) || hoveredOptionIndex === index,
                  })}
                >
                  <div className='f-12-400 flex flex-wrap items-center text-GRAY_1000 gap-1 rounded-md cursor-pointer min-h-5 px-1.5'>
                    <span> Create team :</span>
                    {search && (
                      <span
                        className='px-1.5 py-0.5 rounded cursor-pointer text-black w-fit h-fit text-wrap'
                        style={{ backgroundColor: randomColor ?? COLORS.WHITE }}
                      >
                        {option?.label}
                      </span>
                    )}
                  </div>
                </div>
              ) : search?.length === 0 && !(option as ArrayListOption & { isNew?: boolean })?.isNew ? (
                <div
                  className={cn('w-full px-1.5 py-1 hover:bg-GRAY_50 rounded-md cursor-pointer', {
                    'bg-GRAY_50': (hoveredOptionIndex === null && index === 0) || hoveredOptionIndex === index,
                  })}
                >
                  <span
                    className='f-12-400 text-GRAY_1000 flex px-1.5 py-0.5 w-fit rounded'
                    style={{ backgroundColor: option?.color ?? COLORS.WHITE }}
                  >
                    {transformLabel ? transformLabel(option?.label) : option?.label}
                  </span>
                </div>
              ) : search?.length !== 0 && !(option as ArrayListOption & { isNew?: boolean })?.isNew ? (
                <div
                  className={cn('w-full px-1.5 py-1 hover:bg-GRAY_50 rounded-md cursor-pointer', {
                    'bg-GRAY_50': (hoveredOptionIndex === null && index === 0) || hoveredOptionIndex === index,
                  })}
                >
                  <span
                    className='f-12-400 text-GRAY_1000 flex px-1.5 py-0.5 w-fit rounded'
                    style={{ backgroundColor: option?.color ?? COLORS.WHITE }}
                  >
                    {transformLabel ? transformLabel(option?.label) : option?.label}
                  </span>
                </div>
              ) : null}
            </div>
          ))}
        </CommonWrapper>
      </div>
    </div>
  );
};

export default CustomTeamsDropdown;
