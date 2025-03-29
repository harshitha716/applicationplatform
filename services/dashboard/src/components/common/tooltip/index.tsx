import React, { FC, useCallback, useEffect, useRef, useState } from 'react';
import { COLORS } from 'constants/colors';
import { cn } from 'utils/common';

const TOOLTIP_COLORS: Record<string, string> = {
  TEXT_PRIMARY: COLORS.BLACK,
  SECONDARY: COLORS.WHITE,
};

export enum TooltipPositions {
  TOP = 'top',
  RIGHT = 'right',
  RIGHT_TOP = 'right-top',
  BOTTOM = 'bottom',
  LEFT = 'left',
  LEFT_TOP = 'left-top',
  TOP_LEFT = 'top-left',
  TOP_RIGHT = 'top-right',
  BOTTOM_LEFT = 'bottom-left',
  BOTTOM_RIGHT = 'bottom-right',
}

export interface TooltipProps {
  children: React.ReactNode;
  tooltipBody?: string | React.ReactNode;
  position?: TooltipPositions;
  style?: Record<string, string | number>;
  wrapperClassName?: string;
  tooltipBodyClassName?: string;
  wrapperStyle?: string;
  tooltipBodystyle?: string;
  caratClassName?: string;
  className?: string;
  disabled?: boolean;
  caratOverrideClassName?: string;
  color?: string;
  id?: string;
}

const TOOLTIP_BODY_COLOR_CLASSNAME_MAP = {
  [TOOLTIP_COLORS.TEXT_PRIMARY]: 'bg-black text-white',
  [TOOLTIP_COLORS.SECONDARY]: 'text-TEXT_SECONDARY bg-ZAMP_SECONDARY',
};

const CARET_COLOR_CLASSNAME_MAP = {
  [TOOLTIP_COLORS.TEXT_PRIMARY]: {
    [TooltipPositions.RIGHT]: 'border-r-TEXT_PRIMARY',
    [TooltipPositions.RIGHT_TOP]: 'border-r-TEXT_PRIMARY',
    [TooltipPositions.LEFT]: 'border-l-TEXT_PRIMARY',
    [TooltipPositions.LEFT_TOP]: 'border-l-TEXT_PRIMARY',
    [TooltipPositions.BOTTOM]: 'border-b-TEXT_PRIMARY',
    [TooltipPositions.TOP]: 'border-t-TEXT_PRIMARY',
    [TooltipPositions.BOTTOM_RIGHT]: 'border-b-TEXT_PRIMARY',
    [TooltipPositions.BOTTOM_LEFT]: 'border-b-TEXT_PRIMARY',
    [TooltipPositions.TOP_RIGHT]: 'border-t-TEXT_PRIMARY',
    [TooltipPositions.TOP_LEFT]: 'border-t-TEXT_PRIMARY',
  },
  [TOOLTIP_COLORS.SECONDARY]: {
    [TooltipPositions.RIGHT]: 'border-r-ZAMP_SECONDARY',
    [TooltipPositions.RIGHT_TOP]: 'border-r-ZAMP_SECONDARY',
    [TooltipPositions.LEFT]: 'border-l-ZAMP_SECONDARY',
    [TooltipPositions.LEFT_TOP]: 'border-l-ZAMP_SECONDARY',
    [TooltipPositions.BOTTOM]: 'border-b-ZAMP_SECONDARY',
    [TooltipPositions.TOP]: 'border-t-ZAMP_SECONDARY',
    [TooltipPositions.BOTTOM_RIGHT]: 'border-b-ZAMP_SECONDARY',
    [TooltipPositions.BOTTOM_LEFT]: 'border-b-ZAMP_SECONDARY',
    [TooltipPositions.TOP_RIGHT]: 'border-t-ZAMP_SECONDARY',
    [TooltipPositions.TOP_LEFT]: 'border-t-ZAMP_SECONDARY',
  },
};

export const Tooltip: FC<TooltipProps> = ({
  position = TooltipPositions.BOTTOM,
  wrapperClassName = 'absolute ease-in-out duration-300 opacity-0 pointer-events-none',
  tooltipBodyClassName = 'normal-case w-full relative rounded-2.5 py-4 px-6',
  caratOverrideClassName = 'absolute border-8 border-solid border-transparent w-0 h-0 z-10',
  tooltipBody = '',
  style,
  children,
  wrapperStyle = '',
  tooltipBodystyle = '',
  className = '',
  disabled = false,
  caratClassName = '',
  color = COLORS.BLACK,
  id,
}) => {
  const ref = useRef<HTMLInputElement>(null);
  const parentRef = useRef<HTMLInputElement>(null);
  const [tooltipPosition, setTooltipPosition] = useState({});

  const caratPostionStyle = {
    [TooltipPositions.RIGHT]: 'm-auto -left-4 top-[calc(50%-8px)]',
    [TooltipPositions.LEFT]: 'm-auto -right-4 top-[calc(50%-8px)] ',
    [TooltipPositions.BOTTOM]: 'm-auto -top-4 left-[calc(50%-8px)]',
    [TooltipPositions.TOP]: 'm-auto -bottom-4 left-[calc(50%-8px)]',
    [TooltipPositions.RIGHT_TOP]: 'top-2.5 -left-4',
    [TooltipPositions.LEFT_TOP]: '-right-4 top-2.5',
    [TooltipPositions.BOTTOM_RIGHT]: '-top-4 right-2.5',
    [TooltipPositions.BOTTOM_LEFT]: '-top-4 left-2.5',
    [TooltipPositions.TOP_RIGHT]: '-bottom-4 right-2.5',
    [TooltipPositions.TOP_LEFT]: '-bottom-4 left-2.5',
  };

  const wrapperPositionStyle = {
    [TooltipPositions.RIGHT]: ` ml-2.5 z-[500] left-full`,
    [TooltipPositions.RIGHT_TOP]: `ml-2.5 z-[500] -top-2.5 left-full`,
    [TooltipPositions.LEFT]: 'mr-2 z-[500] right-full',
    [TooltipPositions.LEFT_TOP]: '-top-2.5 mr-2 z-[500] right-full',
    [TooltipPositions.TOP]: ' mb-2.5 z-[500] bottom-full',
    [TooltipPositions.BOTTOM]: 'mt-2.5 z-[500] top-full',
    [TooltipPositions.TOP_RIGHT]: '-right-2.5 mb-2.5 z-[500] bottom-full',
    [TooltipPositions.TOP_LEFT]: '-left-2.5 mb-2.5 z-[500] bottom-full',
    [TooltipPositions.BOTTOM_RIGHT]: '-right-2.5 mt-2.5 z-[500] top-full',
    [TooltipPositions.BOTTOM_LEFT]: '-left-2.5 mt-2.5 z-[500] top-full',
  };

  const calculatePositionStyle = useCallback((position: TooltipPositions) => {
    let style = {};
    let height = ref?.current?.clientHeight;
    let width = ref?.current?.clientWidth;
    const parentHeight = parentRef.current?.clientHeight;
    const parentWidth = parentRef.current?.clientWidth;

    if ([TooltipPositions.LEFT, TooltipPositions.RIGHT].includes(position)) {
      if (height && parentHeight) {
        height = height / 2;
        height = height - parentHeight / 2;
        style = { top: -height + 'px' };
      }
    } else {
      if (width && height && parentWidth && parentHeight) {
        switch (position) {
          case TooltipPositions.BOTTOM:
          case TooltipPositions.TOP: {
            width = width / 2;
            width = width - parentWidth / 2 - 2;

            return { left: -width + 'px' };
          }
          case TooltipPositions.LEFT:
          case TooltipPositions.RIGHT: {
            height = height / 2;
            height = height - parentHeight / 2;

            return { top: -height + 'px' };
          }
        }
      }
    }

    return style;
  }, []);

  useEffect(() => {
    setTooltipPosition(calculatePositionStyle(position));
  }, [position, calculatePositionStyle]);

  return (
    <div
      ref={parentRef}
      className={`relative group/tooltip cursor-pointer ${className}`}
      data-testid={`tooltip-wrapper-${id}`}
    >
      {children}
      {!!tooltipBody && (
        <div
          ref={ref}
          className={cn(
            wrapperClassName,
            wrapperStyle,
            wrapperPositionStyle[position],
            !disabled &&
              'z-1000 group-hover/tooltip:opacity-100 group-hover/tooltip:pointer-events-auto f-12-450 px-3 py-1.5 rounded-md whitespace-nowrap z-999 bg-black text-GRAY_200',
          )}
          style={style ? style : { ...tooltipPosition }}
          data-testid={`tooltip-${id}`}
        >
          <div
            className={`${tooltipBodyClassName} ${TOOLTIP_BODY_COLOR_CLASSNAME_MAP[color]} ${tooltipBodystyle}`}
            data-testid={`tooltip-body-wrapper-${id}`}
          >
            <div
              className={`${caratOverrideClassName} ${caratPostionStyle[position]} ${CARET_COLOR_CLASSNAME_MAP[color]?.[position]} ${caratClassName}`}
              data-testid={`tooltip-carat-${id}`}
            ></div>
            <div className='cursor-text' data-testid={`tooltip-body-${id}`}>
              {tooltipBody}
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
