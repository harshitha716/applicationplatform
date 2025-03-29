import React, { FC } from 'react';
import { SelectedCountTooltipPropsType } from 'types/common/components/dropdown/dropdown.types';
import { Tooltip, TooltipPositions } from 'components/common/tooltip';

const SelectedCountTooltip: FC<SelectedCountTooltipPropsType> = ({ value, tooltipBodyClassName }) => {
  return (
    <div>
      <Tooltip
        style={{ right: '-34px' }}
        caratClassName='border-b-GRAY_700 left-[calc(100%-50px)]'
        position={TooltipPositions.BOTTOM}
        disabled={value.length === 0}
        tooltipBody={
          <div className='flex flex-col gap-2'>
            {value?.map((item) => (
              <div key={item?.value} className='whitespace-nowrap'>
                {item?.label}
              </div>
            ))}
          </div>
        }
        tooltipBodystyle={`bg-GRAY_700 text-white f-12-300 !px-3 !py-2  ${tooltipBodyClassName}`}
      >
        {!!value.length && (
          <div className='bg-BASE_PRIMARY f-10-600 h-4 w-5 flex items-center justify-center'>{value.length}</div>
        )}
      </Tooltip>
    </div>
  );
};

export default SelectedCountTooltip;
