import { OptionsType } from 'types/commonTypes';
import { cn } from 'utils/common';

interface DropdownProps {
  options: OptionsType[];
  onSelect: (widgetId: string) => void;
  activeWidget: string;
  className?: string;
  dropdownRef: React.RefObject<HTMLDivElement>;
}

export const WidgetOptionDropdown = ({ options, onSelect, activeWidget, className, dropdownRef }: DropdownProps) => {
  return (
    <div
      ref={dropdownRef}
      className={cn(
        'absolute z-40 bg-white flex flex-col gap-2 pt-2 pb-1 border border-GRAY_400 rounded-md shadow-tableFilterMenu max-h-[330px] w-[200px]',
        className,
      )}
    >
      <div className='flex flex-col h-full overflow-y-auto custom-scroll-bar-common px-1 select-none'>
        {options.map((option) => (
          <div
            key={option.value}
            onClick={() => onSelect(option.value as string)}
            className={cn('py-2 px-2.5 cursor-pointer select-none rounded hover:bg-GRAY_100', {
              'bg-GRAY_100': activeWidget === option.value,
            })}
          >
            <div className='f-12-400 text-GRAY_1000'>{option.label}</div>
          </div>
        ))}
      </div>
    </div>
  );
};
