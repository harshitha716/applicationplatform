import React from 'react';
import { COLORS } from 'constants/colors';
import { SIZE_TYPES } from 'types/common/components';
import { defaultFn } from 'types/commonTypes';
import { BUTTON_STATE_TYPES, BUTTON_TYPES, ButtonProps, ICON_POSITION_TYPES } from 'types/components/button.type';
import { cn, doDebounce } from 'utils/common';
import ProgressBar from 'components/common/RingProgress';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

const BUTTON_STATE_STYLES = {
  [BUTTON_TYPES.PRIMARY]: {
    [BUTTON_STATE_TYPES.COMMON]: '',
    [BUTTON_STATE_TYPES.DEFAULT]: 'bg-GRAY_1000 text-white',
    [BUTTON_STATE_TYPES.HOVER]: 'hover:bg-GRAY_950 hover:text-white',
    [BUTTON_STATE_TYPES.PRESSED]: 'active:bg-GRAY_950 active:text-white',
    [BUTTON_STATE_TYPES.DISABLED]: 'disabled:cursor-not-allowed disabled:bg-GRAY_100 disabled:text-GRAY_700',
    [BUTTON_STATE_TYPES.LOADING]: 'bg-GRAY_700 !cursor-not-allowed',
  },
  [BUTTON_TYPES.SECONDARY]: {
    [BUTTON_STATE_TYPES.COMMON]: 'border',
    [BUTTON_STATE_TYPES.DEFAULT]: 'outline-0 bg-white text-GRAY_1000 border-BORDER_GRAY_400',
    [BUTTON_STATE_TYPES.HOVER]: 'hover:bg-BG_GRAY_2',
    [BUTTON_STATE_TYPES.PRESSED]: 'active:bg-GRAY_400',
    [BUTTON_STATE_TYPES.DISABLED]: 'disabled:cursor-not-allowed disabled:bg-BG_GRAY_2 disabled:text-GRAY_700',
    [BUTTON_STATE_TYPES.LOADING]: '!cursor-not-allowed',
  },
  [BUTTON_TYPES.TEXT_NAV]: {
    [BUTTON_STATE_TYPES.COMMON]: 'bg-transparent',
    [BUTTON_STATE_TYPES.DEFAULT]: 'text-GRAY_600',
    [BUTTON_STATE_TYPES.HOVER]: 'hover:underline hover:text-GRAY_700',
    [BUTTON_STATE_TYPES.PRESSED]: 'active:underline active:text-GRAY_700',
    [BUTTON_STATE_TYPES.DISABLED]: 'disabled:cursor-not-allowed !text-GRAY_500 !hover:no-underline',
    [BUTTON_STATE_TYPES.LOADING]: '',
  },
  [BUTTON_TYPES.DANGER]: {
    [BUTTON_STATE_TYPES.COMMON]: '',
    [BUTTON_STATE_TYPES.DEFAULT]: 'bg-RED_700 text-white',
    [BUTTON_STATE_TYPES.HOVER]: 'hover:bg-RED_600 hover:text-white',
    [BUTTON_STATE_TYPES.PRESSED]: 'active:bg-RED_600 active:text-white',
    [BUTTON_STATE_TYPES.DISABLED]: 'disabled:cursor-not-allowed disabled:bg-RED_100 disabled:text-white',
    [BUTTON_STATE_TYPES.LOADING]: 'bg-RED_500 !cursor-not-allowed',
  },
  [BUTTON_TYPES.SHARE]: {
    [BUTTON_STATE_TYPES.COMMON]: '',
    [BUTTON_STATE_TYPES.DEFAULT]: 'outline-0 !bg-transparent !text-GRAY_1000 !border-BORDER_GRAY_400',
    [BUTTON_STATE_TYPES.HOVER]: 'hover:bg-BG_GRAY_2',
    [BUTTON_STATE_TYPES.PRESSED]: '!bg-GRAY_400',
    [BUTTON_STATE_TYPES.DISABLED]: 'disabled:cursor-not-allowed disabled:bg-BG_GRAY_2 disabled:text-GRAY_700',
    [BUTTON_STATE_TYPES.LOADING]: '!cursor-not-allowed',
  },
};

const ICON_SIZE_BY_TYPE = {
  [SIZE_TYPES.XLARGE]: 22,
  [SIZE_TYPES.LARGE]: 16,
  [SIZE_TYPES.MEDIUM]: 16,
  [SIZE_TYPES.SMALL]: 14,
  [SIZE_TYPES.XSMALL]: 12,
};

const BUTTON_SIZE_STYLES = {
  [SIZE_TYPES.XLARGE]: {
    wrapperClassBySize: 'h-14 py-4 px-6 gap-1.5',
    wrapperWithOnlySingleIcon: 'h-14 w-14 p-4',
    textClassDefaultBySize: 'f-16-400',
    textClassWithLeftIcons: 'pr-2',
    textClassWithRightIcons: 'pl-2',
    textClassWithoutTrailingIcon: 'pl-1',
  },
  [SIZE_TYPES.LARGE]: {
    wrapperClassBySize: 'h-10 py-3.5 px-3 gap-1.5',
    wrapperWithOnlySingleIcon: 'h-14 w-14 p-4',
    textClassDefaultBySize: 'f-14-500',
    textClassWithLeftIcons: 'pr-2',
    textClassWithRightIcons: 'pl-2',
    textClassWithoutTrailingIcon: 'pl-1',
  },
  [SIZE_TYPES.MEDIUM]: {
    wrapperClassBySize: 'h-8 py-2 px-3.5 gap-1.5',
    wrapperWithOnlySingleIcon: 'h-8 w-11 p-3',
    textClassDefaultBySize: 'f-13-500',
    textClassWithLeftIcons: 'pr-1.5',
    textClassWithRightIcons: 'pl-1.5',
    textClassWithoutTrailingIcon: 'pl-1',
  },
  [SIZE_TYPES.SMALL]: {
    wrapperClassBySize: 'h-7 py-[7px] px-3 gap-1',
    wrapperWithOnlySingleIcon: 'h-9 w-9 p-[9.5px]',
    textClassDefaultBySize: 'f-12-500',
    textClassWithLeftIcons: 'pr-1',
    textClassWithRightIcons: 'pl-1',
    textClassWithoutTrailingIcon: 'pl-2',
  },
  [SIZE_TYPES.XSMALL]: {
    wrapperClassBySize: 'h-[26px] py-1.5 px-2 gap-1',
    wrapperWithOnlySingleIcon: 'h-7 w-7 p-[6.5px]',
    textClassDefaultBySize: 'f-11-500',
    textClassWithLeftIcons: 'pr-1',
    textClassWithRightIcons: 'pl-1',
    textClassWithoutTrailingIcon: 'pl-1.5',
  },
};

const BUTTON_ALIGN_STYLES = {
  wrapperDefaultAlignment: 'flex items-center !justify-center',
  wrapperAlignmentWithoutLeadingIcon: 'flex items-center justify-between',
};

const BUTTON_DEFAULT_STYLES = 'cursor-pointer disabled:cursor-not-allowed overflow-clip rounded-md';

export const Button: React.FC<ButtonProps> = ({
  type = BUTTON_TYPES.PRIMARY,
  className = '',
  disabled = false,
  onClick = defaultFn,
  size = SIZE_TYPES.LARGE,
  state = BUTTON_STATE_TYPES.DEFAULT,
  isLoading = false,
  id = '',
  textSizeOverrideClassName = '',
  children = null,
  childrenClassName = '',
  loader = null,
  customLeadingIcon = null,
  customTrailingIcon = null,
  tabIndex = 0,
  iconPosition = ICON_POSITION_TYPES.RIGHT,
  iconProps,
  customAttributes,
  onMouseEnter = defaultFn,
  onMouseLeave = defaultFn,
  ...rest
}) => {
  const {
    wrapperClassBySize,
    wrapperWithOnlySingleIcon,
    textClassDefaultBySize,
    textClassWithLeftIcons,
    textClassWithRightIcons,
  } = BUTTON_SIZE_STYLES?.[size] || {};

  const { wrapperDefaultAlignment, wrapperAlignmentWithoutLeadingIcon } = BUTTON_ALIGN_STYLES;

  const wrapperSizeClass = children ? wrapperClassBySize : wrapperWithOnlySingleIcon;

  const wrapperAlignmentClass = `${wrapperDefaultAlignment} ${children ? wrapperAlignmentWithoutLeadingIcon : ''}`;

  const iconSize = ICON_SIZE_BY_TYPE[size as keyof typeof SIZE_TYPES];

  const textSizeClass = `${textSizeOverrideClassName ? textSizeOverrideClassName : textClassDefaultBySize} ${
    iconProps ? (iconPosition === ICON_POSITION_TYPES.RIGHT ? textClassWithRightIcons : textClassWithLeftIcons) : ''
  }`;

  const handleButtonClick = (e?: React.MouseEvent<HTMLButtonElement>) => {
    if (!isLoading && !disabled) {
      onClick(e);
    }
  };

  const debouncedClick = doDebounce(handleButtonClick, 500);

  return (
    <button
      type='submit'
      tabIndex={tabIndex}
      data-testid={`btn-${id}`}
      onClick={debouncedClick}
      disabled={disabled}
      className={cn(
        className,
        wrapperAlignmentClass,
        wrapperSizeClass,
        BUTTON_DEFAULT_STYLES,
        BUTTON_STATE_STYLES[type as BUTTON_TYPES]?.[state as BUTTON_STATE_TYPES],
        BUTTON_STATE_STYLES[type as BUTTON_TYPES]?.[BUTTON_STATE_TYPES.COMMON],
        BUTTON_STATE_STYLES[type as BUTTON_TYPES]?.[BUTTON_STATE_TYPES.HOVER],
        BUTTON_STATE_STYLES[type as BUTTON_TYPES]?.[BUTTON_STATE_TYPES.PRESSED],
        BUTTON_STATE_STYLES[type as BUTTON_TYPES]?.[BUTTON_STATE_TYPES.DISABLED],
        isLoading ? BUTTON_STATE_STYLES[type]?.[BUTTON_STATE_TYPES.LOADING] : '',
      )}
      onMouseEnter={onMouseEnter}
      onMouseLeave={onMouseLeave}
      {...customAttributes}
      {...rest}
    >
      {isLoading ? (
        loader ? (
          loader
        ) : (
          <div className='w-full flex items-center justify-center'>
            <ProgressBar
              trackColor={COLORS.BLACK}
              indicatorColor={COLORS.WHITE}
              indicatorWidth={2.5}
              trackWidth={2.5}
              size={ICON_SIZE_BY_TYPE[size] + 5}
              className='animate-spin'
              progress={30}
            />
          </div>
        )
      ) : (
        <>
          {((iconPosition === ICON_POSITION_TYPES.LEFT && iconProps) || customLeadingIcon) && (
            <div>
              {customLeadingIcon ? (
                customLeadingIcon
              ) : (
                <SvgSpriteLoader
                  className={`${iconProps?.className} `}
                  id={iconProps?.id ?? ''}
                  height={iconProps?.height ?? iconSize}
                  width={iconProps?.width ?? iconSize}
                  iconCategory={iconProps?.iconCategory}
                  color={iconProps?.color}
                  size={iconProps?.size}
                />
              )}
            </div>
          )}

          {children && <div className={`${textSizeClass} ${childrenClassName}`}>{children}</div>}

          {((iconPosition === ICON_POSITION_TYPES.RIGHT && iconProps?.id) || customTrailingIcon) && (
            <div>
              {customTrailingIcon ? (
                customTrailingIcon
              ) : (
                <SvgSpriteLoader
                  className={`${iconProps?.className} mt-0.5`}
                  id={iconProps?.id ?? ''}
                  height={iconProps?.height ?? iconSize}
                  width={iconProps?.width ?? iconSize}
                  iconCategory={iconProps?.iconCategory}
                  color={iconProps?.color}
                />
              )}
            </div>
          )}
        </>
      )}
    </button>
  );
};
