import { FC, memo } from 'react';
import { ColDef } from 'ag-grid-community';
import { PIVOT_HEADER_BG } from 'constants/icons';
import Image from 'next/image';
import { snakeCaseToSentenceCase } from 'utils/common';

interface PivotColHeaderProps {
  column: {
    colDef: ColDef;
  };
  displayName: string;
}

const PivotColHeader: FC<PivotColHeaderProps> = ({ column, displayName }) => {
  const contextFieldName = snakeCaseToSentenceCase(column.colDef?.context?.name || '');

  return (
    <div className='relative w-full h-full flex items-end justify-end p-3 border-r-0.5 border-b-0.5 border-GRAY_400 break-words whitespace-normal bg-white overflow-hidden'>
      <Image
        src={PIVOT_HEADER_BG}
        alt='Pivot Header Background'
        fill
        priority
        className='shrink-0 object-cover object-center'
      />
      <span className='relative z-10 f-13-550'>{contextFieldName || displayName}</span>
    </div>
  );
};

export default memo(PivotColHeader);
