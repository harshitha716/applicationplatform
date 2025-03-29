import { FC, memo } from 'react';
import { PIVOT_HEADER_BG } from 'constants/icons';
import { formatColGroupHeaderDisplayName } from 'modules/widgets/Pivot/pivot.utils';
import Image from 'next/image';
import { cn } from 'utils/common';

type PivotAutoGroupHeaderProps = {
  displayName: string;
  isSingleHeader: boolean;
};

const PivotColGroupHeader: FC<PivotAutoGroupHeaderProps> = ({ displayName, isSingleHeader }) => {
  const { mainText, suffix } = formatColGroupHeaderDisplayName(displayName);

  return (
    <div
      className={cn(
        'w-full h-full p-3 flex flex-col items-center justify-center bg-white break-words whitespace-normal overflow-hidden',
        isSingleHeader && 'relative flex items-end justify-end',
      )}
    >
      {isSingleHeader && (
        <Image
          src={PIVOT_HEADER_BG}
          alt='Pivot Header Background'
          fill
          priority
          className='shrink-0 object-cover object-center'
        />
      )}
      <div className='relative flex gap-2 justify-center items-center z-10 text-center'>
        <span className={cn(isSingleHeader ? 'f-13-550' : 'f-13-450')}> {mainText}</span>
        {suffix && (
          <span className='p-1.5 py-1 rounded border border-GRAY_400 bg-white f-12-450 text-GRAY_900'>{suffix}</span>
        )}
      </div>
    </div>
  );
};

export default memo(PivotColGroupHeader);
