import { FC } from 'react';
import { COLORS } from 'constants/colors';
import { ICON_SPRITE_TYPES } from 'constants/icons';
import { cn } from 'utils/common';
import SvgSpriteLoader from 'components/SvgSpriteLoader';

interface NoWidgetDataProps {
  className?: string;
}

const NoWidgetData: FC<NoWidgetDataProps> = ({ className }) => {
  return (
    <div
      className={cn(
        'top-0 right-0 w-full h-[calc(100%-100px)] flex justify-center items-center z-1000 bg-white',
        className,
      )}
    >
      <div className='flex items-center flex-col gap-3'>
        <SvgSpriteLoader
          id='coins-stacked-03'
          iconCategory={ICON_SPRITE_TYPES.FINANCE_AND_ECOMMERCE}
          width={24}
          height={24}
          color={COLORS.GRAY_700}
        />
        <div className='text-GRAY_700 f-12-450'>No data available, try again with different filters</div>
      </div>
    </div>
  );
};

export default NoWidgetData;
