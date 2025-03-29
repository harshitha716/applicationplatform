import React, { FC } from 'react';
import { ValueContainerContentProps } from 'types/common/components/dropdown/dropdown.types';
import SelectedCountTooltip from 'components/common/dropdown/SelectedCountTooltip';

const ValueContainerContent: FC<ValueContainerContentProps> = ({
  labelProps = {},
  value,
  showCountOfSelected,
  tooltipBodyClassName,
}) => {
  return (
    <div className='flex justify-between items-center flex-1 mr-2.5'>
      <div>{labelProps.title}</div>
      {showCountOfSelected && <SelectedCountTooltip value={value} tooltipBodyClassName={tooltipBodyClassName} />}
    </div>
  );
};

export default ValueContainerContent;
