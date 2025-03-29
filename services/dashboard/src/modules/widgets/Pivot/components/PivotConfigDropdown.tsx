import { FC, useRef, useState } from 'react';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { useOnClickOutside } from 'hooks';
import { defaultFnType } from 'types/commonTypes';
import { cn } from 'utils/common';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

interface PivotConfigDropdownProps {
  handleExportAgGridData: defaultFnType;
}

const PivotConfigDropdown: FC<PivotConfigDropdownProps> = ({ handleExportAgGridData }) => {
  const ref = useRef<HTMLDivElement>(null);
  const [showDisplayConfig, setShowDisplayConfig] = useState(false);

  useOnClickOutside(ref, () => setShowDisplayConfig(false));

  return (
    <>
      <div
        className='h-[38px] px-[2px] py-1.5 flex items-center rounded-full border border-GRAY_200 absolute w-fit top-[29px] -left-[11px] cursor-pointer z-1000 overflow-hidden bg-[#fafafa] opacity-0 group-hover:opacity-100 transition-opacity duration-200'
        onClick={() => setShowDisplayConfig(!showDisplayConfig)}
      >
        <SvgSpriteLoader id='dots-vertical' width={16} height={16} iconCategory={ICON_SPRITE_TYPES.GENERAL} />
      </div>
      {showDisplayConfig && (
        <div
          className={cn(
            'absolute z-[9999] top-[72px] -left-2.5 bg-white flex flex-col gap-2 p-1.5 border border-GRAY_400 rounded-md shadow-tableFilterMenu max-h-[330px] min-w-[200px] ',
          )}
          ref={ref}
        >
          <div className='flex flex-col h-full overflow-y-auto custom-scroll-bar-common select-none'>
            <div
              className='flex items-center justify-between gap-2.5 py-1.5 px-2 rounded cursor-pointer hover:bg-GRAY_100'
              onClick={handleExportAgGridData}
            >
              <span className='f-13-450 text-GRAY_1000'>Export Data</span>
            </div>
          </div>
        </div>
      )}
    </>
  );
};

export default PivotConfigDropdown;
