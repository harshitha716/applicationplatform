import { FC, useEffect, useMemo, useRef, useState } from 'react';
import { useGetWidgetDataQuery } from 'apis/widgets';
import { PERIODICITY_TYPES } from 'constants/date.constants';
import { useWindowDimensions } from 'hooks/useWindowDimensions';
import { CURRENCY_SYMBOLS } from 'modules/page/pages.constants';
import { WIDGET_TYPES, WidgetInstanceType } from 'types/api/widgets.types';
import { getCommaSeparatedNumber } from 'utils/common';
import { Tooltip, TooltipPositions } from 'components/common/tooltip';
import CommonWrapper from 'components/commonWrapper';
import { SkeletonTypes } from 'components/commonWrapper/commonWrapper.types';
import SkeletonElement from 'components/skeletons/SkeletonElement';

interface KpiTagProps {
  widgetDetails: Extract<WidgetInstanceType, { widget_type: WIDGET_TYPES.KPI }>;
  currentPageFilters: string;
  isFilterInitialized?: boolean;
  periodicity: string;
  timeColumns: string;
  isFilterLoading?: boolean;
  currency?: string;
  defaultCurrency: string;
}

const KpiTag: FC<KpiTagProps> = ({
  widgetDetails,
  currentPageFilters,
  isFilterInitialized,
  periodicity,
  timeColumns,
  isFilterLoading,
  currency,
  defaultCurrency,
}) => {
  const valueContainerRef = useRef<HTMLDivElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const windowWidth = useWindowDimensions().width;
  const [showTooltip, setShowTooltip] = useState(false);

  const { data: widgetData, isFetching } = useGetWidgetDataQuery(
    {
      widgetId: widgetDetails?.widget_instance_id,
      payload: {
        filters: currentPageFilters,
        time_columns: timeColumns,
        periodicity: periodicity as PERIODICITY_TYPES,
        currency: currency,
      },
    },
    { refetchOnMountOrArgChange: false, skip: !isFilterInitialized },
  );

  const value: string = useMemo(() => {
    const key = widgetDetails?.data_mappings?.mappings?.[0]?.fields?.primary_value?.[0]?.column;
    const data = widgetData?.result?.[0]?.data[0] as Record<string, any>;

    const currency =
      (defaultCurrency || widgetData?.currency) &&
      (CURRENCY_SYMBOLS[defaultCurrency as keyof typeof CURRENCY_SYMBOLS] ??
        CURRENCY_SYMBOLS[widgetData?.currency as keyof typeof CURRENCY_SYMBOLS] ??
        widgetData?.currency);

    if (isNaN(Number(data?.[key]))) return data?.[key];

    return currency ? `${currency} ${getCommaSeparatedNumber(Number(data?.[key]), 2)}` : Number(data?.[key]);
  }, [widgetData]);

  useEffect(() => {
    const callback = () => {
      if (!containerRef.current || !valueContainerRef.current) return;

      const contentWidth = containerRef.current?.offsetWidth - 48;
      const valueWidth = valueContainerRef.current?.scrollWidth;

      if (valueWidth && valueWidth > contentWidth) {
        setShowTooltip(true);
      } else {
        setShowTooltip(false);
      }
    };

    const timer = setTimeout(callback, 500);

    return () => {
      if (timer) clearTimeout(timer);
    };
  }, [containerRef, valueContainerRef, widgetData, isFetching, windowWidth]);

  return (
    <div className='bg-white h-full border border-GRAY_400 rounded-xl px-6 pt-4.5 pb-5' ref={containerRef}>
      <div className='f-13-450 text-GRAY_900 mb-2 truncate'>{widgetDetails?.title}</div>
      <CommonWrapper
        skeletonType={SkeletonTypes.CUSTOM}
        isLoading={isFetching || isFilterLoading}
        loader={<SkeletonElement className='max-w-[250px]' />}
      >
        <Tooltip
          tooltipBody={value}
          disabled={isFetching || !showTooltip}
          tooltipBodyClassName='f-14-300 px-3 ml-2 py-1.5 rounded-md whitespace-nowrap z-999 bg-black text-white'
          position={TooltipPositions.BOTTOM_LEFT}
          className='!cursor-text'
        >
          <div className='f-24-450 text-GRAY_950 truncate sensitive' ref={valueContainerRef}>
            {value}
          </div>
        </Tooltip>
      </CommonWrapper>
    </div>
  );
};

export default KpiTag;
