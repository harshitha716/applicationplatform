import React, { FC, memo } from 'react';
import { MapAny } from 'types/commonTypes';

const AgGridTooltipComponent: FC<MapAny> = ({ value }) => {
  return <div className='p-2 text-TEXT_WHITE bg-black f-12-300'>{value}</div>;
};

export default memo(AgGridTooltipComponent);
