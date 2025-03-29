import { useEffect, useMemo, useState } from 'react';
import { Responsive, WidthProvider } from 'react-grid-layout';
import { useGetSheetDetailsQuery } from 'apis/pages';
import { ZAMP_LOGO_LOADER } from 'constants/lottie/zamp-logo-loader';
import { PAGE_CURRENCY_OPTIONS } from 'modules/page/pages.constants';
import InitializeSheetsFilters from 'modules/sheets/InitializeSheetsFilters';
import SingleSelectFilter from 'modules/widgets/components/SingleSelectFilter';
import WidgetSwitcher from 'modules/widgets/components/widgetSwitcher';
import { ROW_HEIGHT, SCREEN_BREAKPOINTS, WIDGETS_LAYOUT_MARGIN } from 'modules/widgets/widgets.constant';
import { WIDGET_TYPES } from 'types/api/widgets.types';
import CommonWrapper from 'components/commonWrapper';
import { SkeletonTypes } from 'components/commonWrapper/commonWrapper.types';
import DynamicLottiePlayer from 'components/DynamicLottiePlayer';
import FiltersWrapper from 'components/filter/filterMenu/FiltersWrapper';
import { useFiltersContextStore, withFiltersContext } from 'components/filter/filters.context';
import 'react-grid-layout/css/styles.css'; // Include default styles
import 'react-resizable/css/styles.css'; // Include resizable styles

interface SheetsProps {
  pageId: string;
  sheetId: string;
  isPageLoading: boolean;
}

const ResponsiveGridLayout = WidthProvider(Responsive);

const Sheets = ({ pageId, sheetId, isPageLoading }: SheetsProps) => {
  const {
    state: { filtersConfig, isFilterInitialized },
  } = useFiltersContextStore();
  const [currency, setCurrency] = useState<string[]>(['USD']);
  const [widgetDetails, setWidgetDetails] = useState<{
    height: number;
    isSingleHeader: boolean;
  }>({
    height: 0,
    isSingleHeader: true,
  });

  const handleWidgetHeightChange = (height: number, isSingleHeader: boolean) => {
    setWidgetDetails({
      height,
      isSingleHeader,
    });
  };

  const {
    data: sheetDetails,
    isFetching: isSheetLoading,
    isError: isSheetDetailsError,
    refetch: refetchSheetDetails,
  } = useGetSheetDetailsQuery(
    { pageId: pageId as string, sheetId: sheetId as string },
    { skip: !pageId || !sheetId, refetchOnMountOrArgChange: false },
  );

  //converts the pixel height from pivot-table into grid-layout height
  const getHfromWidgetHeight = (widgetHeight: number): number => {
    return (widgetHeight + 20) / 76;
  };

  const sheetLayout = useMemo(() => {
    return sheetDetails?.sheet_config?.sheet_layout?.map((widgetConfig) => {
      const widgetType = sheetDetails?.widget_instances?.find(
        (widget) => widget?.widget_instance_id === widgetConfig?.default_widget,
      )?.widget_type;

      return {
        i: widgetConfig?.default_widget,
        ...widgetConfig?.layout,
        h:
          widgetType === WIDGET_TYPES.PIVOT_TABLE && widgetDetails?.height > 0
            ? Math.min(getHfromWidgetHeight(widgetDetails?.height), widgetConfig?.layout?.h)
            : widgetConfig?.layout?.h,
      };
    });
  }, [sheetDetails?.sheet_config?.sheet_layout, widgetDetails, pageId, sheetId]);

  //Add the max-height on pivot table based on sheet layout height for pivot and current actual height of the grid
  useEffect(() => {
    if (typeof document !== 'undefined' && sheetLayout) {
      if (sheetLayout && sheetDetails?.sheet_config?.sheet_layout[0]?.layout?.h == sheetLayout[0].h) {
        const substractedHeight = widgetDetails?.isSingleHeader ? 93 : 135;

        document.documentElement.style.setProperty(
          '--pivot-max-height',
          `${sheetLayout[0]?.h * 56 + (sheetLayout[0]?.h - 1) * 20 - substractedHeight}px`,
        );
      }
    }
  }, [sheetDetails, sheetLayout]);

  //To reset the widget details on switch of page or sheet
  useEffect(() => {
    setWidgetDetails({
      height: 0,
      isSingleHeader: true,
    });
  }, [pageId, sheetId]);

  return (
    <InitializeSheetsFilters pageId={pageId} sheetId={sheetId}>
      <div className='relative h-[calc(100vh-94px)] overflow-scroll py-6 pl-3 pr-0'>
        <CommonWrapper
          isLoading={isSheetLoading || isPageLoading}
          skeletonType={SkeletonTypes.CUSTOM}
          isError={isSheetDetailsError}
          className='h-full'
          refetchFunction={refetchSheetDetails}
          loader={
            <div className='flex justify-center items-center w-full h-full z-1000 bg-white'>
              <DynamicLottiePlayer
                src={ZAMP_LOGO_LOADER}
                className='lottie-player h-[140px]'
                autoplay
                loop
                keepLastFrame
              />
            </div>
          }
        >
          <div className='flex justify-between items-center z-100 px-5'>
            <div className='f-24-450 text-GRAY_950'>{sheetDetails?.name}</div>
            <div className='flex items-center gap-2'>
              <FiltersWrapper
                allowClear={false}
                label='Filter'
                className='px-0'
                allowActions={false}
                filterConfig={filtersConfig ?? []}
                isPeriodicityEnabled
                isRightAligned
              />
              {isFilterInitialized && !sheetDetails?.sheet_config?.currency?.hide_currency_filter && currency && (
                <div className='flex items-center gap-2'>
                  {!!filtersConfig?.length && <div className='border-r border-GRAY_400 h-7'></div>}
                  <SingleSelectFilter
                    filterKey='currency'
                    options={PAGE_CURRENCY_OPTIONS.filter((option) => option !== 'local')}
                    onFilterChange={(value) => setCurrency(value)}
                    value={currency}
                    label='Currency'
                  />
                </div>
              )}
            </div>
          </div>

          {sheetDetails && (
            <ResponsiveGridLayout
              className='layout'
              layout={sheetLayout}
              cols={{ lg: 16, md: 16, sm: 16, xs: 16 }}
              breakpoints={SCREEN_BREAKPOINTS}
              rowHeight={ROW_HEIGHT}
              width={1200} // Adjust grid width as per container
              margin={WIDGETS_LAYOUT_MARGIN}
              isResizable={false}
              isDraggable={false}
              useCSSTransforms={false}
            >
              {sheetDetails?.sheet_config?.sheet_layout?.map((widgetConfig) => (
                <div
                  key={`widget-${widgetConfig?.default_widget}`}
                  data-grid={sheetLayout?.find((layout) => layout.i === widgetConfig?.default_widget)}
                  className='bg-white'
                >
                  <div key={widgetConfig?.default_widget} className='h-full w-full'>
                    <WidgetSwitcher
                      widgetConfig={widgetConfig}
                      currency={sheetDetails?.sheet_config?.currency?.hide_currency_filter ? [] : currency}
                      defaultCurrency={sheetDetails?.sheet_config?.currency?.default_currency}
                      widgetInstances={sheetDetails?.widget_instances ?? []}
                      handleWidgetHeightChange={handleWidgetHeightChange}
                    />
                  </div>
                </div>
              ))}
            </ResponsiveGridLayout>
          )}
        </CommonWrapper>
      </div>
    </InitializeSheetsFilters>
  );
};

export default withFiltersContext(Sheets);
