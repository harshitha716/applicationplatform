import { FC } from 'react';
import NoWidgetData from 'modules/widgets/components/NoWidgetData';
import WidgetTitle from 'modules/widgets/components/widgetTitle';
import { WIDGET_TYPES } from 'types/api/widgets.types';
import { OptionsType } from 'types/commonTypes';

interface NoPivotDataProps {
  groupWidgetsOptions: OptionsType[];
  onWidgetChange: (widgetId: string) => void;
  title: string;
  activeWidget: string;
}

const NoPivotData: FC<NoPivotDataProps> = ({ groupWidgetsOptions, onWidgetChange, title, activeWidget }) => {
  return (
    <div className='overflow-x-auto flex flex-col w-full h-full border border-GRAY_400 rounded-xl overflow-hidden group'>
      <div className='bg-white w-full flex justify-between items-start h-[110px] p-6 border-b-0.5'>
        <WidgetTitle
          groupWidgetsOptions={groupWidgetsOptions}
          onWidgetChange={onWidgetChange}
          title={title}
          widgetType={WIDGET_TYPES.PIVOT_TABLE}
          activeWidget={activeWidget}
        />
      </div>
      <div className='w-full h-full flex items-center justify-center z-0'>
        <NoWidgetData />
      </div>
    </div>
  );
};

export default NoPivotData;
