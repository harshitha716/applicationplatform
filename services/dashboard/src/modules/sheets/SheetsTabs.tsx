import React, { FC } from 'react';
import { useAppSelector } from 'hooks/toolkit';
import { useRouter } from 'next/router';
import { RootState } from 'store';
import { MenuItem, SIZE_TYPES } from 'types/common/components';
import { BUTTON_TYPES } from 'types/components/button.type';
import { cn } from 'utils/common';
import { LOCAL_STORAGE_KEYS, setToLocalStorage } from 'utils/localstorage';
import { Button } from 'components/common/button/Button';
import { Tooltip, TooltipPositions } from 'components/common/tooltip';
import CommonWrapper from 'components/commonWrapper';
import { SkeletonTypes } from 'components/commonWrapper/commonWrapper.types';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

interface SheetsTabsProps {
  tabs: MenuItem[];
  currentSheetId: string;
  isPageLoading: boolean;
}

const SheetsTabs: FC<SheetsTabsProps> = ({ tabs, currentSheetId, isPageLoading }) => {
  const router = useRouter();
  const { id } = router.query;
  const { isSidebarOpen } = useAppSelector((state: RootState) => state.layoutConfig);

  const handleTabSelect = (selected?: MenuItem) => {
    if (!selected?.value) return;
    router.push(`#${selected?.value}`);
    setToLocalStorage(LOCAL_STORAGE_KEYS.DATA_SHEET_ID, JSON.stringify({ [id as string]: selected?.value }));
  };

  return (
    <div
      className={cn(
        'flex items-center fixed  z-1000 bottom-0 right-0 border-t border-l border-border-GRAY_400 h-[57px] bg-white shadow-pageBottomBar px-8 gap-3 transition-all duration-300',
        !isSidebarOpen ? 'w-full' : 'w-[calc(100%-240px)]',
      )}
    >
      <CommonWrapper
        skeletonType={SkeletonTypes.CUSTOM}
        isLoading={isPageLoading}
        className='flex items-center gap-3'
        loader={<div className='w-25 rounded-md block animate-pulse bg-GRAY_50 h-8' />}
      >
        {tabs?.map((tab) => (
          <Button
            key={tab?.value}
            id='sheets-tabs'
            onClick={() => handleTabSelect(tab)}
            type={BUTTON_TYPES.SECONDARY}
            className={cn(
              'w-fit !rounded-lg',
              currentSheetId === tab?.value ? '!bg-BG_GRAY_2 !border-GRAY_500' : 'bg-white !border-GRAY_400',
            )}
            size={SIZE_TYPES.MEDIUM}
          >
            <div className={`transition-all duration-100 f-12-450 whitespace-nowrap`}>{tab?.label}</div>
          </Button>
        ))}
      </CommonWrapper>

      <Tooltip
        tooltipBody='Coming soon'
        color='{TMS_COLORS.GRAY_200}'
        tooltipBodyClassName='f-12-300 px-3 py-1.5 rounded-md whitespace-nowrap z-999 bg-black text-GRAY_200'
        className='z-1'
        position={TooltipPositions.TOP}
      >
        <div className='flex items-center gap-1 f-12-450 text-GRAY_700 cursor-not-allowed select-none'>
          <SvgSpriteLoader id='plus' width={16} height={16} />
          <div className='whitespace-nowrap'>New sheet</div>
        </div>
      </Tooltip>
    </div>
  );
};

export default SheetsTabs;
