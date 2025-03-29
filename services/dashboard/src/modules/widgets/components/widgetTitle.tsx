import { useRef, useState } from 'react';
import ReactDOM from 'react-dom';
import { useOnClickOutside } from 'hooks';
import { WidgetOptionDropdown } from 'modules/widgets/components/WidgetOptionDropdown';
import { WIDGET_TYPES } from 'types/api/widgets.types';
import { OptionsType } from 'types/commonTypes';
import { cn } from 'utils/common';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

interface WidgetTitleProps {
  title: string;
  groupWidgetsOptions: OptionsType[];
  onWidgetChange: (widgetId: string) => void;
  widgetType: WIDGET_TYPES;
  isSingleValue?: boolean;
  activeWidget: string;
  className?: string;
  isPortalNeeded?: boolean;
}

const WidgetTitle = ({
  title,
  groupWidgetsOptions,
  onWidgetChange,
  widgetType,
  activeWidget,
  className,
  isPortalNeeded = false,
}: WidgetTitleProps) => {
  const dropdownRef = useRef<HTMLDivElement>(null);
  const titleRef = useRef<HTMLDivElement>(null);
  const [isGroupWidgetOptionsOpen, setIsGroupWidgetOptionsOpen] = useState(false);
  const isGroupWidgetOptions = groupWidgetsOptions.length > 1;
  const isPivotTable = widgetType === WIDGET_TYPES.PIVOT_TABLE;

  useOnClickOutside(dropdownRef, (event) => {
    if (titleRef?.current && titleRef.current.contains(event?.target as Node)) return;
    setIsGroupWidgetOptionsOpen(false);
  });

  const handleToggle = (e: React.MouseEvent<HTMLDivElement>) => {
    e.stopPropagation();
    if (!isGroupWidgetOptions) return;
    setIsGroupWidgetOptionsOpen((prev) => !prev);
  };

  const handleWidgetChange = (widgetId: string) => {
    onWidgetChange(widgetId);
    setIsGroupWidgetOptionsOpen(false);
  };

  return (
    <div className={className}>
      <div
        ref={titleRef}
        className={cn(
          'px-6 flex flex-col items-start w-fit select-none cursor-pointer',
          ![WIDGET_TYPES.DONUT_CHART, WIDGET_TYPES.PIE_CHART].includes(widgetType) && 'mb-10',
          isPivotTable && isGroupWidgetOptions && 'px-0 gap-y-2 items-start justify-center mb-0',
          isPivotTable && !isGroupWidgetOptions && 'mb-0 px-0 justify-center cursor-default',
        )}
        onClick={handleToggle}
      >
        <div className='flex items-center gap-1'>
          <span className='f-18-450 text-GRAY_1000'>{title}</span>
          {isGroupWidgetOptions && (
            <SvgSpriteLoader
              id='chevron-down'
              width={18}
              height={18}
              className={cn(
                'text-GRAY_900 transition-transform duration-300',
                isGroupWidgetOptionsOpen ? 'rotate-180' : 'rotate-0',
              )}
            />
          )}
        </div>

        {isGroupWidgetOptions && (
          <span className='f-12-450 text-GRAY_700 opacity-0 group-hover:opacity-100 transition-opacity duration-200'>{`${groupWidgetsOptions.length} Variants`}</span>
        )}
      </div>
      {isGroupWidgetOptionsOpen &&
        (isPivotTable && isPortalNeeded ? (
          ReactDOM?.createPortal(
            <WidgetOptionDropdown
              options={groupWidgetsOptions}
              onSelect={handleWidgetChange}
              activeWidget={activeWidget}
              className='top-14 left-5'
              dropdownRef={dropdownRef}
            />,
            document?.querySelector('.pivot') as HTMLElement,
          )
        ) : (
          <WidgetOptionDropdown
            options={groupWidgetsOptions}
            onSelect={handleWidgetChange}
            activeWidget={activeWidget}
            className='top-12 left-6'
            dropdownRef={dropdownRef}
          />
        ))}
    </div>
  );
};

export default WidgetTitle;
